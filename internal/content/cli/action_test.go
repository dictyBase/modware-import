package cli

import (
	"context"
	"io"
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "frontpage art", expected: "frontpage-art"},
		{input: "already-slugified", expected: "already-slugified"},
		{
			input:    "  leading and trailing spaces  ",
			expected: "leading-and-trailing-spaces",
		},
		{input: "special!@#characters", expected: "special-characters"},
		{input: "multiple   spaces", expected: "multiple-spaces"},
		{
			input:    "news 256cd371-8710-462a-9f7c-d34774526c8f",
			expected: "news-256cd371-8710-462a-9f7c-d34774526c8f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			require.Equal(t, tt.expected, Slugify(tt.input))
		})
	}
}

func TestParseSourceKey(t *testing.T) {
	tests := []struct {
		name              string
		key               string
		expectedName      string
		expectedNamespace string
	}{
		{
			name:              "dfp nested key",
			key:               "some/path/dfp-about.json",
			expectedName:      "about",
			expectedNamespace: "frontpage",
		},
		{
			name:              "dfp flat key",
			key:               "dfp-contact.json",
			expectedName:      "contact",
			expectedNamespace: "frontpage",
		},
		{
			name:              "dsc nested key",
			key:               "some/path/dsc-information.json",
			expectedName:      "information",
			expectedNamespace: "stockcenter",
		},
		{
			name:              "news nested hyphenated name",
			key:               "some/path/news-256cd371-8710-462a-9f7c-d34774526c8f.json",
			expectedName:      "256cd371-8710-462a-9f7c-d34774526c8f",
			expectedNamespace: "news",
		},
		{
			name:              "hyphenated dfp name",
			key:               "dfp-about-us.json",
			expectedName:      "about-us",
			expectedNamespace: "frontpage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, namespace, err := parseSourceKey(tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.expectedName, name)
			require.Equal(t, tt.expectedNamespace, namespace)
		})
	}
}

func TestParseSourceKeyErrors(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "unknown prefix", key: "xyz-about.json", want: "unknown namespace prefix"},
		{name: "missing name", key: "dfp-.json", want: "missing name"},
		{name: "missing extension", key: "dfp-about", want: "file extension"},
		{name: "empty extension", key: "dfp-about.", want: "file extension"},
		{name: "missing namespace prefix", key: "about.json", want: "missing namespace prefix"},
		{name: "empty stem", key: ".json", want: "file extension"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseSourceKey(tt.key)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestValidatePayload(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		wantErr string
	}{
		{name: "empty payload", payload: []byte(""), wantErr: "empty content payload"},
		{name: "whitespace payload", payload: []byte("  \n\t "), wantErr: "empty content payload"},
		{name: "invalid json", payload: []byte("{not json"), wantErr: "invalid JSON payload"},
		{name: "valid object", payload: []byte(`{"a":1}`)},
		{name: "valid array", payload: []byte(`[1,2,3]`)},
		{name: "valid scalar string", payload: []byte(`"hello"`)},
		{name: "valid scalar number", payload: []byte(`42`)},
		{name: "valid null", payload: []byte(`null`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePayload("some-key.json", tt.payload)
			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

// fakeObjectSource implements objectSource with in-memory data.
type fakeObjectSource struct {
	keys      []string
	listErrs  map[string]error
	data      map[string][]byte
	getErrs   map[string]error
	readErrs  map[string]error
	closeErrs map[string]error
	closed    []string
	getCalls  []string
}

func (f *fakeObjectSource) list(
	_, _ string, doneCh <-chan struct{},
) <-chan minio.ObjectInfo {
	ch := make(chan minio.ObjectInfo)
	go func() {
		defer close(ch)
		for _, key := range f.keys {
			select {
			case <-doneCh:
				return
			case ch <- minio.ObjectInfo{Key: key, Err: f.listErrs[key]}:
			}
		}
	}()
	return ch
}

func (f *fakeObjectSource) get(_, key string) (io.ReadCloser, error) {
	f.getCalls = append(f.getCalls, key)
	if err := f.getErrs[key]; err != nil {
		return nil, err
	}
	return &fakeReadCloser{
		data:     f.data[key],
		readErr:  f.readErrs[key],
		closeErr: f.closeErrs[key],
		onClose: func() {
			f.closed = append(f.closed, key)
		},
	}, nil
}

type fakeReadCloser struct {
	data     []byte
	readErr  error
	closeErr error
	onClose  func()
}

func (r *fakeReadCloser) Read(p []byte) (int, error) {
	if r.readErr != nil {
		return 0, r.readErr
	}
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}

func (r *fakeReadCloser) Close() error {
	r.onClose()
	return r.closeErr
}

// fakeContentClient implements content.ContentServiceClient.
type fakeContentClient struct {
	getBySlugFn func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error)
	storeFn     func(context.Context, *content.StoreContentRequest, ...grpc.CallOption) (*content.Content, error)
	updateFn    func(context.Context, *content.UpdateContentRequest, ...grpc.CallOption) (*content.Content, error)

	getCalls    int
	storeCalls  int
	updateCalls int
}

func (f *fakeContentClient) GetContentBySlug(
	ctx context.Context, in *content.ContentRequest, opts ...grpc.CallOption,
) (*content.Content, error) {
	f.getCalls++
	return f.getBySlugFn(ctx, in, opts...)
}

func (f *fakeContentClient) GetContent(
	context.Context, *content.ContentIdRequest, ...grpc.CallOption,
) (*content.Content, error) {
	return nil, status.Error(codes.Unimplemented, "unused")
}

func (f *fakeContentClient) StoreContent(
	ctx context.Context, in *content.StoreContentRequest, opts ...grpc.CallOption,
) (*content.Content, error) {
	f.storeCalls++
	return f.storeFn(ctx, in, opts...)
}

func (f *fakeContentClient) UpdateContent(
	ctx context.Context, in *content.UpdateContentRequest, opts ...grpc.CallOption,
) (*content.Content, error) {
	f.updateCalls++
	return f.updateFn(ctx, in, opts...)
}

func (f *fakeContentClient) DeleteContent(
	context.Context, *content.ContentIdRequest, ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "unused")
}

func (f *fakeContentClient) ListContents(
	context.Context, *content.ListParameters, ...grpc.CallOption,
) (*content.ContentCollection, error) {
	return nil, status.Error(codes.Unimplemented, "unused")
}

func testLogger() *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logrus.NewEntry(logger)
}

func testLoggerWithHook() (*logrus.Entry, *test.Hook) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	hook := test.NewLocal(logger)
	return logrus.NewEntry(logger), hook
}

func testItem() sourceItem {
	return sourceItem{
		key:       "dfp-about.json",
		name:      "about",
		namespace: "frontpage",
		slug:      "frontpage-about",
		payload:   `{"a":1}`,
	}
}

func storeResp() *content.Content {
	return &content.Content{
		Data: &content.ContentData{
			Id: 1,
			Attributes: &content.ContentAttributes{
				Name:      "about",
				Namespace: "frontpage",
				Content:   `{"a":1}`,
			},
		},
	}
}

func TestPreflightListing(t *testing.T) {
	t.Run("preserves listing order", func(t *testing.T) {
		src := &fakeObjectSource{
			keys: []string{"a/dfp-first.json", "m/news-item.json", "z/dsc-later.json"},
			data: map[string][]byte{
				"a/dfp-first.json": []byte(`{"a":1}`),
				"m/news-item.json": []byte(`{"m":1}`),
				"z/dsc-later.json": []byte(`{"z":1}`),
			},
		}

		items, err := preflight(src, testLogger(), "bucket", "prefix")

		require.NoError(t, err)
		require.Len(t, items, 3)
		require.Equal(t, "a/dfp-first.json", items[0].key)
		require.Equal(t, "m/news-item.json", items[1].key)
		require.Equal(t, "z/dsc-later.json", items[2].key)
	})

	t.Run("retains exact payload bytes", func(t *testing.T) {
		payload := "{\"b\":  2}\n\t"
		src := &fakeObjectSource{
			keys: []string{"dfp-about.json"},
			data: map[string][]byte{"dfp-about.json": []byte(payload)},
		}
		items, err := preflight(src, testLogger(), "bucket", "prefix")
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, payload, items[0].payload)
	})

	t.Run("skips duplicate slug without fetching its body", func(t *testing.T) {
		firstKey := "a/dfp-about.json"
		duplicateKey := "b/dfp-about.json"
		src := &fakeObjectSource{
			// listing order decides the winner; the fake emits keys in slice order
			keys: []string{firstKey, duplicateKey},
			data: map[string][]byte{
				firstKey: []byte(`{"a":1}`),
			},
			getErrs: map[string]error{
				duplicateKey: io.ErrUnexpectedEOF,
			},
		}
		logger, hook := testLoggerWithHook()

		items, err := preflight(src, logger, "bucket", "prefix")

		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, firstKey, items[0].key)
		require.Equal(t, []string{firstKey}, src.getCalls)

		require.Len(t, hook.Entries, 1)
		require.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
		require.Contains(t, hook.LastEntry().Message, "duplicate slug")
		require.Contains(t, hook.LastEntry().Message, "frontpage-about")
		require.Contains(t, hook.LastEntry().Message, duplicateKey)
	})
}

func TestPreflightErrors(t *testing.T) {
	t.Run("surfaces listing error with source key", func(t *testing.T) {
		src := &fakeObjectSource{
			keys:     []string{"dfp-about.json"},
			listErrs: map[string]error{"dfp-about.json": io.ErrUnexpectedEOF},
		}
		_, err := preflight(src, testLogger(), "bucket", "prefix")
		require.Error(t, err)
		require.Contains(t, err.Error(), "dfp-about.json")
	})

	t.Run("surfaces get error with source key", func(t *testing.T) {
		src := &fakeObjectSource{
			keys:    []string{"dfp-about.json"},
			getErrs: map[string]error{"dfp-about.json": io.ErrUnexpectedEOF},
		}
		_, err := preflight(src, testLogger(), "bucket", "prefix")
		require.Error(t, err)
		require.Contains(t, err.Error(), "dfp-about.json")
	})

	t.Run("closes body on successful read", func(t *testing.T) {
		src := &fakeObjectSource{
			keys: []string{"dfp-about.json"},
			data: map[string][]byte{"dfp-about.json": []byte(`{"a":1}`)},
		}
		_, err := preflight(src, testLogger(), "bucket", "prefix")
		require.NoError(t, err)
		require.Contains(t, src.closed, "dfp-about.json")
	})

	t.Run("closes body on read error", func(t *testing.T) {
		src := &fakeObjectSource{
			keys:     []string{"dfp-about.json"},
			data:     map[string][]byte{"dfp-about.json": []byte(`{"a":1}`)},
			readErrs: map[string]error{"dfp-about.json": io.ErrUnexpectedEOF},
		}
		_, err := preflight(src, testLogger(), "bucket", "prefix")
		require.Error(t, err)
		require.Contains(t, src.closed, "dfp-about.json")
	})

	t.Run("surfaces close error when read succeeds", func(t *testing.T) {
		src := &fakeObjectSource{
			keys:      []string{"dfp-about.json"},
			data:      map[string][]byte{"dfp-about.json": []byte(`{"a":1}`)},
			closeErrs: map[string]error{"dfp-about.json": io.ErrUnexpectedEOF},
		}
		_, err := preflight(src, testLogger(), "bucket", "prefix")
		require.Error(t, err)
		require.Contains(t, src.closed, "dfp-about.json")
	})
}

func TestLoadContentPreflightFailureNoMutations(t *testing.T) {
	src := &fakeObjectSource{
		keys: []string{"dfp-about.json", "dfp-broken.json"},
		data: map[string][]byte{
			"dfp-about.json":  []byte(`{"a":1}`),
			"dfp-broken.json": []byte(`{not json`),
		},
	}
	client := &fakeContentClient{
		getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return nil, status.Error(codes.NotFound, "not found")
		},
		storeFn: func(context.Context, *content.StoreContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return storeResp(), nil
		},
		updateFn: func(context.Context, *content.UpdateContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return storeResp(), nil
		},
	}

	err := loadContent(client, testLogger(), src, "bucket", "prefix")
	require.Error(t, err)
	require.Zero(t, client.getCalls)
	require.Zero(t, client.storeCalls)
	require.Zero(t, client.updateCalls)
}

func TestUpsertContent_CreateOnNotFound(t *testing.T) {
	var stored *content.StoreContentRequest
	client := &fakeContentClient{
		getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return nil, status.Error(codes.NotFound, "not found")
		},
		storeFn: func(_ context.Context, in *content.StoreContentRequest, _ ...grpc.CallOption) (*content.Content, error) {
			stored = in
			return storeResp(), nil
		},
	}

	item := testItem()
	require.NoError(t, upsertContent(client, testLogger(), item))
	require.Equal(t, 1, client.storeCalls)
	require.Zero(t, client.updateCalls)
	require.NotNil(t, stored)
	attrs := stored.Data.Attributes
	require.Equal(t, "about", attrs.Name)
	require.Equal(t, "frontpage", attrs.Namespace)
	require.Equal(t, "frontpage-about", attrs.Slug)
	require.Equal(t, item.payload, attrs.Content)
	require.Equal(t, "pfey@northwestern.edu", attrs.CreatedBy)
}

func TestUpsertContent_NoOpOnIdentical(t *testing.T) {
	client := &fakeContentClient{
		getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return &content.Content{
				Data: &content.ContentData{
					Id:         7,
					Attributes: &content.ContentAttributes{Content: testItem().payload},
				},
			}, nil
		},
	}

	require.NoError(t, upsertContent(client, testLogger(), testItem()))
	require.Equal(t, 1, client.getCalls)
	require.Zero(t, client.storeCalls)
	require.Zero(t, client.updateCalls)
}

func TestUpsertContent_UpdateOnDifferent(t *testing.T) {
	var updated *content.UpdateContentRequest
	client := &fakeContentClient{
		getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return &content.Content{
				Data: &content.ContentData{
					Id:         7,
					Attributes: &content.ContentAttributes{Content: "old"},
				},
			}, nil
		},
		updateFn: func(_ context.Context, in *content.UpdateContentRequest, _ ...grpc.CallOption) (*content.Content, error) {
			updated = in
			return storeResp(), nil
		},
	}

	item := testItem()
	require.NoError(t, upsertContent(client, testLogger(), item))
	require.Equal(t, 1, client.updateCalls)
	require.Zero(t, client.storeCalls)
	require.NotNil(t, updated)
	require.Equal(t, int64(7), updated.Id)
	require.Equal(t, item.payload, updated.Data.Attributes.Content)
	require.Equal(t, "pfey@northwestern.edu", updated.Data.Attributes.UpdatedBy)
	require.Equal(t, []string{"content"}, updated.UpdateMask.Paths)
}

func TestUpsertContent_InvalidGetResponses(t *testing.T) {
	tests := []struct {
		name string
		resp *content.Content
		want string
	}{
		{name: "missing data", resp: &content.Content{}, want: "missing data"},
		{
			name: "missing attributes",
			resp: &content.Content{Data: &content.ContentData{Id: 7}},
			want: "missing attributes",
		},
		{
			name: "invalid record id",
			resp: &content.Content{
				Data: &content.ContentData{
					Attributes: &content.ContentAttributes{Content: `{"a":1}`},
				},
			},
			want: "invalid record id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fakeContentClient{
				getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
					return tt.resp, nil
				},
			}
			err := upsertContent(client, testLogger(), testItem())
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
			require.Zero(t, client.storeCalls)
			require.Zero(t, client.updateCalls)
		})
	}
}

func TestUpsertContent_AlreadyExistsRace(t *testing.T) {
	t.Run("refetches and no-ops when equal", func(t *testing.T) {
		calls := 0
		client := &fakeContentClient{
			getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
				calls++
				if calls == 1 {
					return nil, status.Error(codes.NotFound, "not found")
				}
				return &content.Content{
					Data: &content.ContentData{
						Id:         9,
						Attributes: &content.ContentAttributes{Content: testItem().payload},
					},
				}, nil
			},
			storeFn: func(context.Context, *content.StoreContentRequest, ...grpc.CallOption) (*content.Content, error) {
				return nil, status.Error(codes.AlreadyExists, "already exists")
			},
		}

		require.NoError(t, upsertContent(client, testLogger(), testItem()))
		require.Equal(t, 2, client.getCalls)
		require.Equal(t, 1, client.storeCalls)
		require.Zero(t, client.updateCalls)
	})

	t.Run("refetches and updates when different", func(t *testing.T) {
		calls := 0
		var updated *content.UpdateContentRequest
		client := &fakeContentClient{
			getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
				calls++
				if calls == 1 {
					return nil, status.Error(codes.NotFound, "not found")
				}
				return &content.Content{
					Data: &content.ContentData{
						Id:         9,
						Attributes: &content.ContentAttributes{Content: "other"},
					},
				}, nil
			},
			storeFn: func(context.Context, *content.StoreContentRequest, ...grpc.CallOption) (*content.Content, error) {
				return nil, status.Error(codes.AlreadyExists, "already exists")
			},
			updateFn: func(_ context.Context, in *content.UpdateContentRequest, _ ...grpc.CallOption) (*content.Content, error) {
				updated = in
				return storeResp(), nil
			},
		}

		item := testItem()
		require.NoError(t, upsertContent(client, testLogger(), item))
		require.Equal(t, 2, client.getCalls)
		require.Equal(t, 1, client.storeCalls)
		require.Equal(t, 1, client.updateCalls)
		require.Equal(t, int64(9), updated.Id)
		require.Equal(t, item.payload, updated.Data.Attributes.Content)
	})
}

func TestUpsertContent_GetError(t *testing.T) {
	client := &fakeContentClient{
		getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return nil, status.Error(codes.Internal, "boom")
		},
	}
	item := testItem()
	err := upsertContent(client, testLogger(), item)
	require.Error(t, err)
	require.Contains(t, err.Error(), item.key)
	require.Contains(t, err.Error(), item.slug)
	require.Zero(t, client.storeCalls)
	require.Zero(t, client.updateCalls)
}

func TestUpsertContent_CreateError(t *testing.T) {
	client := &fakeContentClient{
		getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return nil, status.Error(codes.NotFound, "not found")
		},
		storeFn: func(context.Context, *content.StoreContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return nil, status.Error(codes.Internal, "boom")
		},
	}
	item := testItem()
	err := upsertContent(client, testLogger(), item)
	require.Error(t, err)
	require.Contains(t, err.Error(), item.key)
	require.Contains(t, err.Error(), item.slug)
	require.Zero(t, client.updateCalls)
}

func TestUpsertContent_UpdateError(t *testing.T) {
	client := &fakeContentClient{
		getBySlugFn: func(context.Context, *content.ContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return &content.Content{
				Data: &content.ContentData{
					Id:         7,
					Attributes: &content.ContentAttributes{Content: "old"},
				},
			}, nil
		},
		updateFn: func(context.Context, *content.UpdateContentRequest, ...grpc.CallOption) (*content.Content, error) {
			return nil, status.Error(codes.Internal, "boom")
		},
	}
	item := testItem()
	err := upsertContent(client, testLogger(), item)
	require.Error(t, err)
	require.Contains(t, err.Error(), item.key)
	require.Contains(t, err.Error(), item.slug)
	require.Zero(t, client.storeCalls)
}
