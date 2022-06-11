package deepscan

func mapStrToSlice(src map[string]struct{}) []string {
	r := make([]string, 0, len(src))

	for value := range src {
		r = append(r, value)
	}

	return r
}

func sliceStrToMap(src []string) map[string]struct{} {
	r := make(map[string]struct{}, len(src))

	for _, value := range src {
		r[value] = struct{}{}
	}

	return r
}
