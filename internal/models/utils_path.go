package models

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	Glob string

	ResolvedPath struct {
		ImportPath string
		LocalPath  string
		AbsPath    string
	}
)

// Match check if path is subset of glob, for example:
//   - github.com/**/library/*/abc
//
// will match:
//   - github.com/a/b/c/library/any/abc
//   - github.com/test/library/another/awesome
//
// and not match:
//   - github.com/a/b/c/library/any
//   - github.com/library/another/awesome
func (glob Glob) Match(testedPath string) (bool, error) {
	regGlob := string(glob)
	regGlob = strings.ReplaceAll(regGlob, ".", "\\.") // safe dots
	regGlob = strings.ReplaceAll(regGlob, "/", "\\/") // safe slash

	regGlob = strings.ReplaceAll(regGlob, "**", "<M_ALL>") // super glob tmp
	regGlob = strings.ReplaceAll(regGlob, "*", "[^\\/]+")  // single glob
	regGlob = strings.ReplaceAll(regGlob, "<M_ALL>", ".*") // super glob
	regGlob = fmt.Sprintf("^%s$", regGlob)

	matcher, err := regexp.Compile(regGlob)
	if err != nil {
		return false, fmt.Errorf("failed compile glob '%s' as regexp '%s': %w",
			glob,
			regGlob,
			err,
		)
	}

	return matcher.MatchString(testedPath), nil
}
