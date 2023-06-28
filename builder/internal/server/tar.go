package server

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/magefile/mage/target"
)

var tarOutput = "server.tar.gz"

func Tar() error {
	tarInputGlob := []string{
		GetServerFolder() + "/**/*",
		GetServerFolder() + "/**",
		GetServerFolder(),
	}

	CopyServerFiles()

	newer, err := target.Glob(tarOutput, tarInputGlob...)

	if err != nil {
		return err
	}

	if !newer {
		return nil
	}

	fmt.Println("Creating server tar bundle...")

	archive, err := os.Create("server.tar.gz")

	if err != nil {
		return err
	}

	defer archive.Close()

	var buf bytes.Buffer

	zr := gzip.NewWriter(&buf)
	tw := tar.NewWriter(zr)

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithFinalMSG("✓ Done!\n"))

	s.Start()

	walker := func(path string, info os.FileInfo, err error) error {
		s.Suffix = " Walk: " + path

		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, path)

		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(path)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.Open(path)

			if err != nil {
				return err
			}

			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}

		return nil
	}

	err = filepath.Walk(GetServerFolder(), walker)

	if err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	if err := zr.Close(); err != nil {
		return err
	}

	s.Suffix = " Copying data..."

	_, err = io.Copy(archive, &buf)

	s.Stop()

	if err != nil {
		return err
	}

	return nil
}
