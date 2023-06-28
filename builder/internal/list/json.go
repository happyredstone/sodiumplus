package list

import "encoding/json"

func Json() (string, error) {
	urls, err := CreateModList()

	if err != nil {
		return "", err
	}

	val, err := json.Marshal(urls)

	if err != nil {
		return "", err
	}

	return string(val), nil
}
