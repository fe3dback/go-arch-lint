package module

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/mod/modfile"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type Fetcher struct {
	rootDirectory models.PathAbsolute
	configPath    models.PathRelative
}

func NewFetcher(
	rootDirectory models.PathAbsolute,
	configPath models.PathRelative,
) *Fetcher {
	return &Fetcher{
		rootDirectory: rootDirectory,
		configPath:    configPath,
	}
}

func (f *Fetcher) Fetch() (models.ProjectInfo, error) {
	modPath := models.PathAbsolute(path.Join(string(f.rootDirectory), "go.mod"))
	modContent, err := os.ReadFile(string(modPath))
	if err != nil {
		return models.ProjectInfo{}, fmt.Errorf("failed to read config file '%s': %w", modPath, err)
	}

	modData, err := modfile.ParseLax(string(modPath), modContent, nil)
	if err != nil {
		return models.ProjectInfo{}, fmt.Errorf("failed parse go.mod file '%s': %w", modPath, err)
	}

	return models.ProjectInfo{
		Directory:  f.rootDirectory,
		ConfigPath: models.PathAbsolute(path.Join(string(f.rootDirectory), string(f.configPath))),
		Module:     models.GoModule(modData.Module.Mod.Path),
	}, nil
}
