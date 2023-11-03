package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	return nil
}
