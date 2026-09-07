package assembler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type resolver struct {
	pathResolver  pathResolver
	rootDirectory string
	moduleName    string
}

func newResolver(
	pathResolver pathResolver,
	rootDirectory string,
	moduleName string,
) *resolver {
	return &resolver{
		pathResolver:  pathResolver,
		rootDirectory: rootDirectory,
		moduleName:    moduleName,
	}
}

func (r *resolver) resolveLocalGlobPath(localGlobPath string) ([]models.ResolvedPath, error) {
	list := make([]models.ResolvedPath, 0)

	// localGlobPath always comes from archfile in unix form ('a/b/**'),
	// but rootDirectory is os specific ('/app' or 'C:\app'), so glob pattern
	// should be joined with os specific separator
	absPath := filepath.Join(r.rootDirectory, filepath.FromSlash(localGlobPath))
	resolved, err := r.pathResolver.Resolve(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path '%s'", absPath)
	}

	for _, absResolvedPath := range resolved {
		// resolved paths are os specific, but LocalPath/ImportPath is a part of
		// go import path, that is always in unix form. Otherwise on windows we will
		// get import like 'module/C:\app\internal\pkg' and no one import will match it
		localPath, err := localPathOf(r.rootDirectory, absResolvedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve local path of '%s': %w", absResolvedPath, err)
		}

		importPath := r.moduleName
		if localPath != "" {
			importPath = fmt.Sprintf("%s/%s", r.moduleName, localPath)
		}

		list = append(list, models.ResolvedPath{
			ImportPath: importPath,
			LocalPath:  localPath,
			AbsPath:    filepath.Clean(absResolvedPath),
		})
	}

	return list, nil
}

// localPathOf returns absPath relative to rootDirectory in unix form ('a/b'),
// empty string is returned when absPath is rootDirectory itself.
func localPathOf(rootDirectory string, absPath string) (string, error) {
	localPath, err := filepath.Rel(rootDirectory, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to make '%s' relative to '%s': %w", absPath, rootDirectory, err)
	}

	localPath = strings.Trim(filepath.ToSlash(localPath), "/")
	if localPath == "." {
		return "", nil
	}

	return localPath, nil
}
