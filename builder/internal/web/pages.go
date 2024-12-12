package web

import (
	"errors"
	"html/template"
	"os"
	"strings"

	_ "embed"
	"path/filepath"
)

type FileData struct {
	Name string
	Type string
}

type PageData struct {
	Path  string
	Files []FileData
}

//go:embed page.tmpl
var PageTemplate string

func CreatePage(root string, path string) (string, error) {
	real_root, err := filepath.EvalSymlinks(root)

	if err != nil {
		return "", err
	}

	real_file, err := filepath.EvalSymlinks(path)

	if err != nil {
		return "", err
	}

	ok := strings.HasPrefix(real_file, real_root)

	if !ok {
		return "", errors.New("path out of bounds")
	}

	real_path, err := filepath.Rel(real_root, real_file)

	if err != nil {
		return "", err
	}

	files, err := os.ReadDir(real_file)

	if err != nil {
		return "", err
	}

	items := []FileData{}

	for _, file := range files {
		kind := "file"

		if file.IsDir() {
			kind = "dir"
		}

		items = append(items, FileData{
			Name: file.Name(),
			Type: kind,
		})
	}

	data := PageData{
		Path:  real_path,
		Files: items,
	}

	buf := &strings.Builder{}
	tmpl := template.Must(template.New("html-tmpl").Parse(PageTemplate))
	err = tmpl.Execute(buf, data)

	if err != nil {
		return "", err
	}

	templateData := buf.String()

	return templateData, nil
}
