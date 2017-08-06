package apkg

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// File is the opened apkg.
type File struct {
	path  string
	zip   *zip.ReadCloser
	close []func() error
}

func (f *File) Path() string {
	return f.path
}

// Config parses config.json into Config.
func (f *File) Config() (*Config, error) {
	c := new(Config)
	r := newTgzReader(f.find("control.tar.gz"), "config.json")
	defer r.Close()

	err := json.NewDecoder(r).Decode(c)
	return c, err
}

func (f *File) reader(tgz string, name string) *tgzReader {
	r := newTgzReader(f.find(tgz), name)

	// Clean up all readers when File is closed.
	f.close = append(f.close, r.Close)

	return r
}

// Changelog returns
func (f *File) Changelog() io.ReadCloser {
	return f.reader("control.tar.gz", "changelog.txt")
}

func (f *File) Description() io.ReadCloser {
	return f.reader("control.tar.gz", "description.txt")
}

func Open(name string) (*File, error) {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return nil, fmt.Errorf("apk: could not read zip: %s: %v", name, err)
	}

	f := &File{zip: rc, path: name}
	err = f.validate()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *File) find(name string) *zip.File {
	for _, ff := range f.zip.File {
		if ff.Name == name {
			return ff
		}
	}
	return nil
}

func (f *File) validate() error {
	v := f.find("apkg-version")
	if v == nil {
		return errors.New("apkg: could not find apkg-version inside apk")
	}
	rc, err := v.Open()
	if err != nil {
		return err
	}
	b := make([]byte, 3) // Fits 2.0
	_, err = rc.Read(b)
	if err != nil {
		return err
	}
	if !bytes.Equal(b, []byte{'2', '.', '0'}) {
		return fmt.Errorf("apkg: incompatible apk version %s, expected 2.0", b)
	}
	for _, ff := range []string{"control.tar.gz", "data.tar.gz"} {
		if f.find(ff) == nil {
			return fmt.Errorf("apkg: could not find %s in apk", ff)
		}
	}

	return nil
}

func (f *File) Close() error {
	var err1 error
	for _, c := range f.close {
		err2 := c()
		if err1 == nil {
			err1 = err2
		}
	}
	f.close = nil
	err2 := f.zip.Close()
	if err2 != nil {
		return err2
	}
	return err1
}
