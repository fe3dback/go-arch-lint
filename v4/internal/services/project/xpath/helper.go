package xpath

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Helper struct {
	projectDirectory models.PathAbsolute
	indexFiles       map[models.PathRelative]models.FileRef
	indexDirectories map[models.PathRelative][]models.FileRef
}

func NewHelper(workingDirectory string) *Helper {
	// todo: refactor for composite + functions
	// todo: exclude files/directories by input context
	// todo: input context(ext filter, type filter, excludes, etc..)
	// todo: tests

	workDir, err := filepath.Abs(workingDirectory)
	if err != nil {
		panic(fmt.Errorf("failed get working directory: %w", err))
	}

	h := &Helper{
		projectDirectory: models.PathAbsolute(workDir),
		indexFiles:       make(map[models.PathRelative]models.FileRef, 255),
		indexDirectories: make(map[models.PathRelative][]models.FileRef, 64),
	}

	err = h.scanProjectFilesToIndex()
	if err != nil {
		panic(fmt.Errorf("failed build project files index from workDir '%s': %w", workingDirectory, err))
	}

	return h
}

func (h *Helper) scanProjectFilesToIndex() error {
	return filepath.Walk(string(h.projectDirectory), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed walking %q: %w", path, err)
		}

		isInteresting := info.IsDir() || strings.ToLower(filepath.Ext(info.Name())) == ".go"
		if !isInteresting {
			return nil
		}

		relPathStr, err := filepath.Rel(string(h.projectDirectory), path)
		if err != nil {
			return fmt.Errorf("failed getting relative path '%q' from '%q': %w", path, h.projectDirectory, err)
		}

		relPath := models.PathRelative(relPathStr)

		// add to files index
		h.indexFiles[relPath] = models.FileRef{
			IsDir: info.IsDir(),
			Path:  models.PathAbsolute(path),
		}

		// add to directories index
		if _, exists := h.indexDirectories[relPath]; !exists && info.IsDir() {
			h.indexDirectories[relPath] = make([]models.FileRef, 0, 8)
		}

		if !info.IsDir() {
			dirRelPathStr, err := filepath.Rel(string(h.projectDirectory), filepath.Dir(path))
			if err != nil {
				return fmt.Errorf("failed getting relative path '%q' from '%q': %w", path, h.projectDirectory, err)
			}

			dirRelPath := models.PathRelative(dirRelPathStr)
			h.indexDirectories[dirRelPath] = append(h.indexDirectories[dirRelPath], models.FileRef{
				IsDir: info.IsDir(),
				Path:  models.PathAbsolute(path),
			})
		}

		return nil
	})
}

// MatchProjectFiles will find all exist files in project by provided path or glob
// with any type. But this will try to found files only inside project directory
func (h *Helper) MatchProjectFiles(somePath any, queryType models.FileMatchQueryType) ([]models.FileRef, error) {
	switch actual := somePath.(type) {
	case models.PathRelative:
		return h.matchFileExact(actual, queryType), nil
	case models.PathAbsolute:
		rel, err := filepath.Rel(string(h.projectDirectory), string(actual))
		if err != nil {
			return nil, fmt.Errorf("failed getting relative path from '%q': %w", string(actual), err)
		}

		return h.matchFileExact(models.PathRelative(rel), queryType), nil
	case models.PathRelativeGlob:
		return h.matchFileGlob(actual, queryType)
	case models.PathAbsoluteGlob:
		relGlob, err := filepath.Rel(string(h.projectDirectory), string(actual))
		if err != nil {
			return nil, fmt.Errorf("failed getting relative path from '%q': %w", string(actual), err)
		}

		return h.matchFileGlob(models.PathRelativeGlob(relGlob), queryType)
	}

	return []models.FileRef{}, fmt.Errorf("failed match files by pattern, unknown type %T", somePath)
}

func (h *Helper) matchFileExact(path models.PathRelative, queryType models.FileMatchQueryType) []models.FileRef {
	found, exist := h.indexFiles[path]
	if !exist {
		return nil
	}

	if !isMatchedByQueryType(found, queryType) {
		return nil
	}

	return []models.FileRef{
		found,
	}
}

func (h *Helper) matchFileGlob(path models.PathRelativeGlob, queryType models.FileMatchQueryType) ([]models.FileRef, error) {
	patternNormal, err := glob.Compile(string(path), '/')
	if err != nil {
		return nil, fmt.Errorf("failed compile glob matcher '%s': %w", path, err)
	}

	var patternLast glob.Glob

	if strings.HasSuffix(string(path), "/**") {
		pathLast := strings.TrimSuffix(string(path), "/**")
		patternLast, err = glob.Compile(pathLast, '/')
		if err != nil {
			return nil, fmt.Errorf("failed compile glob matcher '%s': %w", pathLast, err)
		}
	}

	results := make([]models.FileRef, 0, 16)

	// glob working only with directories right now
	for relative, refs := range h.indexDirectories {
		matchedNormal := patternNormal.Match(string(relative))
		matchedLast := false

		if patternLast != nil {
			matchedLast = patternLast.Match(string(relative))
		}

		if !(matchedNormal || matchedLast) {
			continue
		}

		if queryType == models.FileMatchQueryTypeAll || queryType == models.FileMatchQueryTypeOnlyDirectories {
			results = append(results, h.indexFiles[relative])
		}

		if queryType == models.FileMatchQueryTypeAll || queryType == models.FileMatchQueryTypeOnlyFiles {
			results = append(results, refs...)
		}
	}

	return results, nil
}

func isMatchedByQueryType(ref models.FileRef, queryType models.FileMatchQueryType) bool {
	switch queryType {
	case models.FileMatchQueryTypeAll:
		return true
	case models.FileMatchQueryTypeOnlyFiles:
		return !ref.IsDir
	case models.FileMatchQueryTypeOnlyDirectories:
		return ref.IsDir
	default:
		panic(fmt.Sprintf("unexpected query type '%v'", queryType))
	}
}
