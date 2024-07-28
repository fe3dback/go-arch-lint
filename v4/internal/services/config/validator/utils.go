package validator

func xpathOr(xpath string, defaultPath string) string {
	if xpath != "" {
		return xpath
	}

	return defaultPath
}
