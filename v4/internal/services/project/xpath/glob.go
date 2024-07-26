package xpath

// example: /hello/*/world/**/anything
// matched
// - /hello/1/world/sub/anything
// - /hello/1/world/sub/2/anything
// - /hello/1/world/sub/2/3/anything
func expandGlob(index []string, pattern string) ([]string, error) {
	//todo:
	// matcher, err := regexp.Compile()
	return []string{}, nil
}
