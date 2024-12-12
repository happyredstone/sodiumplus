package list

import (
	"fmt"
	"strings"

	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/helpers"
)

const MarkdownPrefix = "# Mod List\n\n> %s v%s\n\n"

func Markdown() (string, error) {
	pack, _, err := helpers.GetPack()

	if err != nil {
		return "", err
	}

	urls, err := CreateModList()

	if err != nil {
		return "", err
	}

	out := ""

	for _, url := range urls {
		out += "- [" + url.Mod + "](" + url.Url + ")\n"
	}

	out = strings.Trim(out, "\n")
	prefix := fmt.Sprintf(MarkdownPrefix, pack.Name, pack.Version)

	return prefix + out, nil
}
