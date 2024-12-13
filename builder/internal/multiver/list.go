package multiver

import (
	"os"
	"slices"
	"strings"
)

func FindVersions(dir string) (map[string][]string, error) {
	vers := map[string][]string{}

	dirs, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	for _, item := range dirs {
		if item.IsDir() {
			items := []string{}
			loaders, err := os.ReadDir(dir + "/" + item.Name())

			if err != nil {
				return nil, err
			}

			for _, loader := range loaders {
				if slices.Contains(LowerLoaders, strings.ToLower(loader.Name())) {
					realLoader := MappedLoaders[strings.ToLower(loader.Name())]

					items = append(items, realLoader)
				}
			}

			vers[item.Name()] = items
		}
	}

	return vers, nil
}
