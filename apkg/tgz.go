package apkg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"

	"github.com/pkg/errors"
)

type tgzReader struct {
	zip  *zip.File
	name string
	err  error

	rc io.ReadCloser
	gz *gzip.Reader
	t  *tar.Reader
}

func (r *tgzReader) Read(b []byte) (int, error) {
	if r.err == nil && r.t == nil {
		r.open()
	}
	if r.err != nil {
		return 0, r.err
	}

	var n int
	n, r.err = r.t.Read(b)
	return n, r.err
}

func (r *tgzReader) open() {
	r.rc, r.err = r.zip.Open()
	if r.err != nil {
		return
	}

	r.gz, r.err = gzip.NewReader(r.rc)
	if r.err != nil {
		r.rc.Close()
		return
	}

	r.t = tar.NewReader(r.gz)

	for {
		h, err := r.t.Next()
		if err != nil {
			r.err = errors.Wrapf(err, "apkg: could not find %q in %q", r.name, r.zip.Name)
			return
		}

		if h.FileInfo().Name() == r.name {
			break
		}
	}
}

func (r *tgzReader) Close() error {
	var err error
	if r.gz != nil {
		err = r.gz.Close()
	}
	if r.rc != nil {
		err2 := r.rc.Close()
		if err == nil {
			err = err2
		}
	}
	if r.err == nil {
		r.err = err
	}
	return err
}

func newTgzReader(zf *zip.File, name string) *tgzReader {
	return &tgzReader{zip: zf, name: name}
}
