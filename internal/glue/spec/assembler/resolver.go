package assembler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type resolver struct {
	pathResolver  PathResolver
	rootDirectory string
	moduleName    string
}

func newResolver(
	pathResolver PathResolver,
	rootDirectory string,
	moduleName string,
) *resolver {
	return &resolver{
		pathResolver:  pathResolver,
		rootDirectory: rootDirectory,
		moduleName:    moduleName,
	}
}

func (r *resolver) resolveLocalPath(localPathMask string) ([]models.ResolvedPath, error) {
	list := make([]models.ResolvedPath, 0)

	absPath := fmt.Sprintf("%s/%s", r.rootDirectory, localPathMask)
	resolved, err := r.pathResolver.Resolve(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path '%s'", absPath)
	}

	for _, absResolvedPath := range resolved {
		localPath := strings.TrimPrefix(absResolvedPath, fmt.Sprintf("%s/", r.rootDirectory))
		localPath = strings.TrimRight(localPath, "/")
		importPath := fmt.Sprintf("%s/%s", r.moduleName, localPath)

		list = append(list, models.ResolvedPath{
			ImportPath: strings.TrimRight(importPath, "/"),
			LocalPath:  strings.TrimRight(localPath, "/"),
			AbsPath:    filepath.Clean(strings.TrimRight(absResolvedPath, "/")),
		})
	}

	return list, nil
}

func (r *resolver) resolveVendorPath(localPathMask string) ([]models.ResolvedPath, error) {
	list, err := r.resolveLocalPath(localPathMask)
	if err != nil {
		return nil, err
	}

	transformedList := make([]models.ResolvedPath, len(list))
	for ind, globalFormat := range list {
		transformedList[ind] = r.transformPathGlobalToVendor(globalFormat)
	}

	return transformedList, nil
}

func (r *resolver) transformPathGlobalToVendor(global models.ResolvedPath) models.ResolvedPath {
	return models.ResolvedPath{
		ImportPath: strings.TrimPrefix(global.LocalPath, "vendor/"),
		LocalPath:  global.LocalPath,
		AbsPath:    global.AbsPath,
	}
}
