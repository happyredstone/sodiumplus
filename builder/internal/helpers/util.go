package helpers

import "os"

func Exists(file string) bool {
	_, err := os.Stat(file)

	return !os.IsNotExist(err)
}
