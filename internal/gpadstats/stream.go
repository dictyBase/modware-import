package gpadstats

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	IOE "github.com/IBM/fp-go/v2/ioeither"
)

const header = "DB_Object_ID\tNegation\tRelation\tGO_ID\tDB_Reference\tEvidence_Code\tWith_or_From\tInteracting_Taxon_ID\tDate\tAssigned_By\tAnnotation_Extensions\tAnnotation_Properties\n"

type wrappedReadCloser struct {
	io.Reader
	closer func() error
}

func (w *wrappedReadCloser) Close() error {
	return w.closer()
}

func httpGet(url string) IOE.IOEither[error, io.ReadCloser] {
	return IOE.TryCatchError(func() (io.ReadCloser, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("bad status: %s", resp.Status)
		}
		return resp.Body, nil
	})
}

func gzipReader(r io.ReadCloser) IOE.IOEither[error, io.ReadCloser] {
	return IOE.TryCatchError(func() (io.ReadCloser, error) {
		gr, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		return &wrappedReadCloser{
			Reader: gr,
			closer: func() error {
				e1 := gr.Close()
				e2 := r.Close()
				if e1 != nil {
					return e1
				}
				return e2
			},
		}, nil
	})
}

func transformStream(r io.ReadCloser) IOE.IOEither[error, io.Reader] {
	return IOE.TryCatchError(func() (io.Reader, error) {
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()
			defer r.Close()
			br := bufio.NewReader(r)
			// Skip 3 lines
			for range 3 {
				if _, err := br.ReadBytes('\n'); err != nil {
					_ = pw.CloseWithError(err)
					return
				}
			}
			// Write Header
			if _, err := pw.Write([]byte(header)); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
			// Copy rest
			if _, err := io.Copy(pw, br); err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}()
		return pr, nil
	})
}
