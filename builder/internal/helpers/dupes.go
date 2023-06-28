package helpers

func Dedupe[K comparable](slice []K) []K {
	keys := make(map[K]bool)
	list := []K{}

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
