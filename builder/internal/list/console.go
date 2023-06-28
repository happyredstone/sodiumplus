package list

import "strings"

func Console() (string, error) {
	urls, err := CreateModList()

	if err != nil {
		return "", err
	}

	out := ""

	for _, url := range urls {
		out += url.Mod + ": " + url.Url + "\n"
	}

	out = strings.Trim(out, "\n")

	return out, nil
}
