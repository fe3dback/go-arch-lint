package info

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"

	"golang.org/x/mod/modfile"
)

type Assembler struct {
}

func NewAssembler() *Assembler {
	return &Assembler{}
}

func (a *Assembler) ProjectInfo(rootDirectory string, archFilePath string) (common.Project, error) {
	projectPath, err := filepath.Abs(rootDirectory)
	if err != nil {
		return common.Project{}, fmt.Errorf("failed to resolve abs path '%s'", rootDirectory)
	}

	// check arch file
	goArchFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, archFilePath))
	_, err = os.Stat(goArchFilePath)
	if os.IsNotExist(err) {
		return common.Project{}, fmt.Errorf("not found archfile in '%s'", goArchFilePath)
	}

	// check go.mod
	goModFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, models.DefaultGoModFileName))
	_, err = os.Stat(goModFilePath)
	if os.IsNotExist(err) {
		return common.Project{}, fmt.Errorf("not found project '%s' in '%s'",
			models.DefaultGoModFileName,
			goModFilePath,
		)
	}

	// parse go.mod
	moduleName, err := checkCmdExtractModuleName(goModFilePath)
	if err != nil {
		return common.Project{}, fmt.Errorf("failed get module name: %w", err)
	}

	return common.Project{
		Directory:      projectPath,
		GoArchFilePath: goArchFilePath,
		GoModFilePath:  goModFilePath,
		ModuleName:     moduleName,
	}, nil
}

func checkCmdExtractModuleName(goModPath string) (string, error) {
	goModFile, err := checkCmdParseGoModFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("can`t parse gomod: %w", err)
	}

	moduleName := goModFile.Module.Mod.Path
	if moduleName == "" {
		return "", fmt.Errorf("%s should contain module name in 'module'", models.DefaultGoModFileName)
	}

	return moduleName, nil
}

func checkCmdParseGoModFile(path string) (*modfile.File, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %w", path, err)
	}

	mod, err := modfile.ParseLax(path, file, nil)
	if err != nil {
		return nil, fmt.Errorf("modfile parseLax failed '%s': %w", path, err)
	}

	return mod, nil
}
