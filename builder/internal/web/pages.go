package web

import (
	"html/template"
	"os"
	"strings"

	_ "embed"
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

func CreatePage(path string) (string, error) {
	files, err := os.ReadDir(path)

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
		Path:  path,
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
