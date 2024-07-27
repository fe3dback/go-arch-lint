package xpath

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Helper struct {
	workingDirectory models.PathAbsolute
	index            map[models.PathRelative]models.FileRef
}

func NewHelper(workingDirectory string) *Helper {
	workDir, err := filepath.Abs(workingDirectory)
	if err != nil {
		panic(fmt.Errorf("failed get working directory: %w", err))
	}

	h := &Helper{
		workingDirectory: models.PathAbsolute(workDir),
		index:            make(map[models.PathRelative]models.FileRef, 255),
	}

	err = h.indexFiles()
	if err != nil {
		panic(fmt.Errorf("failed build project files index from workDir '%s': %w", workingDirectory, err))
	}

	return h
}

func (h *Helper) indexFiles() error {
	return filepath.Walk(string(h.workingDirectory), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed walking %q: %w", path, err)
		}

		absPath := models.PathAbsolute(path)
		relPath, err := filepath.Rel(string(h.workingDirectory), path)
		if err != nil {
			return fmt.Errorf("failed getting relative path '%q' from '%q': %w", path, h.workingDirectory, err)
		}

		h.index[models.PathRelative(relPath)] = models.FileRef{
			IsDir: info.IsDir(),
			Path:  absPath,
		}

		return nil
	})
}

// MatchProjectFiles will find all exist files in project by provided path or glob
// with any type. But this will try to found files only inside project directory
func (h *Helper) MatchProjectFiles(somePath any, onlyDirectories bool) ([]models.FileRef, error) {
	switch actual := somePath.(type) {
	case models.PathRelative:
		return h.matchFileExact(actual, onlyDirectories), nil
	case models.PathAbsolute:
		rel, err := filepath.Rel(string(h.workingDirectory), string(actual))
		if err != nil {
			return nil, fmt.Errorf("failed getting relative path from '%q': %w", string(actual), err)
		}

		return h.matchFileExact(models.PathRelative(rel), onlyDirectories), nil
	case models.PathRelativeGlob:
		return h.matchFileGlob(actual, onlyDirectories)
	}

	return []models.FileRef{}, fmt.Errorf("failed match files by pattern, unknown type %T", somePath)
}

func (h *Helper) prependWorkdir(relPath models.PathRelative) models.PathAbsolute {
	return models.PathAbsolute(path.Join(string(h.workingDirectory), string(relPath)))
}

func (h *Helper) matchFileExact(path models.PathRelative, onlyDirectories bool) []models.FileRef {
	found, exist := h.index[path]
	if !exist {
		return nil
	}

	if onlyDirectories && !found.IsDir {
		return nil
	}

	return []models.FileRef{
		found,
	}
}

func (h *Helper) matchFileGlob(path models.PathRelativeGlob, onlyDirectories bool) ([]models.FileRef, error) {
	patternNormal, err := glob.Compile(string(path), '/')
	if err != nil {
		return nil, fmt.Errorf("failed compile glob matcher '%s': %w", path, err)
	}

	var patternSelf glob.Glob

	if strings.HasSuffix(string(path), "/**") {
		pathSelf := strings.TrimSuffix(string(path), "/**")
		patternSelf, err = glob.Compile(pathSelf, '/')
		if err != nil {
			return nil, fmt.Errorf("failed compile glob matcher '%s': %w", pathSelf, err)
		}
	}

	results := make([]models.FileRef, 0, 16)

	for relative, ref := range h.index {
		matchSelf := false

		if patternSelf != nil {
			matchSelf = patternSelf.Match(string(relative))
		}

		if !(patternNormal.Match(string(relative)) || matchSelf) {
			continue
		}

		if onlyDirectories && !ref.IsDir {
			continue
		}

		results = append(results, ref)
	}

	return results, nil
}
