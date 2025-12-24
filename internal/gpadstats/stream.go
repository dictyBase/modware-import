package gpadstats

import (
	"bufio"
	"compress/gzip"
	"io"
	"net/http"

	F "github.com/IBM/fp-go/v2/function"
	H "github.com/IBM/fp-go/v2/http"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEH "github.com/IBM/fp-go/v2/ioeither/http"
)

const header = "DB_Object_ID\tNegation\tRelation\tGO_ID\tDB_Reference\tEvidence_Code\tWith_or_From\t" +
	"Interacting_Taxon_ID\tDate\tAssigned_By\tAnnotation_Extensions\tAnnotation_Properties\n"

type wrappedReadCloser struct {
	io.Reader
	closer func() error
}

func (w *wrappedReadCloser) Close() error {
	return w.closer()
}

func httpGet(url string) IOE.IOEither[error, io.ReadCloser] {
	return F.Pipe4(
		url,
		IOEH.MakeGetRequest,
		IOEH.MakeClient(http.DefaultClient).Do,
		IOE.ChainEitherK(H.ValidateResponse),
		IOE.Map[error](H.GetBody),
	)
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
