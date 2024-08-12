package xpath

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Helper struct {
	matchers      map[string]typeMatcher
	queryCtx      queryContext
	cachedRegExps map[string]*regexp.Regexp
}

// todo: fix glob matcher
// todo: tests table (lookup)

func NewHelper(
	projectDirectory string,
	matcherRelative typeMatcher,
	matcherAbsolute typeMatcher,
	matcherGlobRelative typeMatcher,
	matcherGlobAbsolute typeMatcher,
) *Helper {
	rootDirectory, err := filepath.Abs(projectDirectory)
	if err != nil {
		panic(fmt.Errorf("failed get working directory: %w", err))
	}

	srv := &Helper{
		queryCtx: newQueryContext(models.PathAbsolute(rootDirectory)),
		matchers: map[string]typeMatcher{
			getType(models.PathRelative("/")):     matcherRelative,
			getType(models.PathAbsolute("/")):     matcherAbsolute,
			getType(models.PathRelativeGlob("/")): matcherGlobRelative,
			getType(models.PathAbsoluteGlob("/")): matcherGlobAbsolute,
		},
	}

	err = srv.indexProjectFiles()
	if err != nil {
		panic(fmt.Errorf("failed build project files index from rootDirectory '%s': %w", rootDirectory, err))
	}

	return srv
}

func (h *Helper) indexProjectFiles() error {
	return filepath.Walk(string(h.queryCtx.projectDirectory), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed walking %q: %w", path, err)
		}

		relativePathStr, err := filepath.Rel(string(h.queryCtx.projectDirectory), path)
		if err != nil {
			return fmt.Errorf("failed getting relative path '%q' from '%q': %w", path, h.queryCtx.projectDirectory, err)
		}

		relativePath := models.PathRelative(relativePathStr)
		extLower := strings.ToLower(filepath.Ext(path))
		extLower = strings.TrimLeft(extLower, ".")

		h.queryCtx.index.appendToIndex(relativePath, models.FileDescriptor{
			PathRel:   relativePath,
			PathAbs:   models.PathAbsolute(path),
			IsDir:     info.IsDir(),
			Extension: extLower,
		})

		return nil
	})
}

func (h *Helper) FindProjectFiles(query models.FileQuery) ([]models.FileDescriptor, error) {
	pathType := getType(query.Path)
	matcher, exist := h.matchers[pathType]
	if !exist {
		return nil, fmt.Errorf("unknown matcher type %s", pathType)
	}

	if matcher == nil {
		return nil, fmt.Errorf("NIL matcher registered for type %s", pathType)
	}

	// match
	found, err := matcher.match(&h.queryCtx, query)
	if err != nil {
		return nil, fmt.Errorf("failed match files by path '%s/%s': %w", query.WorkingDirectory, query.Path, err)
	}

	// filter
	result := make([]models.FileDescriptor, 0, len(found))
	for _, dst := range found {
		suitable, err := h.isSuitable(dst, &query)
		if err != nil {
			return nil, fmt.Errorf("failed check file name '%s': %w", dst.PathRel, err)
		}

		if !suitable {
			continue
		}

		result = append(result, dst)
	}

	return result, nil
}

func (h *Helper) isSuitable(dst models.FileDescriptor, query *models.FileQuery) (bool, error) {
	// only directories
	if dst.IsDir && !(query.Type == models.FileMatchQueryTypeAll || query.Type == models.FileMatchQueryTypeOnlyDirectories) {
		return false, nil
	}

	// only files
	if !dst.IsDir && !(query.Type == models.FileMatchQueryTypeAll || query.Type == models.FileMatchQueryTypeOnlyFiles) {
		return false, nil
	}

	// find dir
	dstDirectory := dst.PathRel
	if !dst.IsDir {
		dstDirectory = models.PathRelative(filepath.Dir(string(dst.PathRel)))
	}

	// exclude by directory
	if len(query.ExcludeDirectories) > 0 {
		if slices.Contains(query.ExcludeDirectories, dstDirectory) {
			return false, nil
		}
	}

	// exclude by file path
	if len(query.ExcludeFiles) > 0 {
		if !dst.IsDir && slices.Contains(query.ExcludeFiles, dst.PathRel) {
			return false, nil
		}
	}

	// exclude by regexp
	for _, pattern := range query.ExcludeRegexp {
		reg, err := h.compileRegexp(string(pattern))
		if err != nil {
			return false, fmt.Errorf("failed compile regular expression '%s': %w", pattern, err)
		}

		if reg.MatchString(string(dst.PathRel)) {
			return false, nil
		}
	}

	// exclude by file ext
	if len(query.Extensions) > 0 {
		if !dst.IsDir && !slices.Contains(query.Extensions, dst.Extension) {
			return false, nil
		}
	}

	// ok
	return true, nil
}

func (h *Helper) compileRegexp(expr string) (*regexp.Regexp, error) {
	if _, ok := h.cachedRegExps[expr]; !ok {
		regular, err := regexp.Compile(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid regular expression '%s': %w", expr, err)
		}

		h.cachedRegExps[expr] = regular
	}

	return h.cachedRegExps[expr], nil
}
