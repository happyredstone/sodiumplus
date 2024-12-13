package mcmeta

func FetchVersionList() ([]string, error) {
	vers, err := FetchVersionManifest()

	if err != nil {
		return []string{}, err
	}

	res := []string{}

	for _, ver := range vers.Versions {
		res = append(res, ver.Id)
	}

	return res, nil
}
