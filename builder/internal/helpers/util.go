package helpers

import (
	"io"
	"os"
)

func Exists(file string) bool {
	_, err := os.Stat(file)

	return !os.IsNotExist(err)
}

func Copy(srcpath string, dstpath string) error {
	r, err := os.Open(srcpath)

	if err != nil {
		return err
	}

	defer r.Close()

	w, err := os.Create(dstpath)
	if err != nil {
		return err
	}

	defer func() {
		if c := w.Close(); err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)

	return err
}
