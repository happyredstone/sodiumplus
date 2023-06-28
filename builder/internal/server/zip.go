package server

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/magefile/mage/target"
)

var zipOutput = "server.zip"

func Zip() error {
	zipInputGlob := []string{
		GetServerFolder() + "/**/*",
		GetServerFolder() + "/**",
		GetServerFolder(),
	}

	CopyServerFiles()

	newer, err := target.Glob(zipOutput, zipInputGlob...)

	if err != nil {
		return err
	}

	if !newer {
		return nil
	}

	fmt.Println("Creating server zip bundle...")

	archive, err := os.Create("server.zip")

	if err != nil {
		return err
	}

	defer archive.Close()

	writer := zip.NewWriter(archive)

	defer writer.Close()

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithFinalMSG("✓ Done!\n"))

	s.Start()

	walker := func(path string, info os.FileInfo, err error) error {
		s.Suffix = " Walk: " + path

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		out, err := writer.Create(path)

		if err != nil {
			return err
		}

		_, err = io.Copy(out, file)

		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(GetServerFolder(), walker)

	s.Stop()

	if err != nil {
		return err
	}

	return nil
}
