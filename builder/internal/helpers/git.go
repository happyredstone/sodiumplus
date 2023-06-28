package helpers

import (
	"os"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

var IgnoreDefaults = []string{
	".git/**",
	".gitattributes",
	".gitignore",
	".DS_Store",
	"/*.zip",
	"*.mrpack",
	"packwiz.exe",
	"packwiz",
}

func ReadGitignore(path string) (*gitignore.GitIgnore, bool) {
	data, err := os.ReadFile(path)

	if err != nil {
		return gitignore.CompileIgnoreLines(IgnoreDefaults...), false
	}

	s := strings.Split(string(data), "\n")

	var lines []string

	lines = append(lines, IgnoreDefaults...)
	lines = append(lines, s...)

	return gitignore.CompileIgnoreLines(lines...), true
}
