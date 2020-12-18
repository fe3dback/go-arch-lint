package info

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Assembler struct {
}

func NewAssembler() *Assembler {
	return &Assembler{}
}

func (a *Assembler) ProjectInfo(rootDirectory string, archFilePath string) (models.ProjectInfo, error) {
	projectPath, err := filepath.Abs(rootDirectory)
	if err != nil {
		return models.ProjectInfo{}, fmt.Errorf("failed to resolve abs path '%s'", rootDirectory)
	}

	// check arch file
	goArchFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, archFilePath))
	_, err = os.Stat(goArchFilePath)
	if os.IsNotExist(err) {
		return models.ProjectInfo{}, fmt.Errorf("not found archfile in '%s'", goArchFilePath)
	}

	// check go.mod
	goModFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, models.ProjectInfoDefaultGoModFileName))
	_, err = os.Stat(goModFilePath)
	if os.IsNotExist(err) {
		return models.ProjectInfo{}, fmt.Errorf("not found project '%s' in '%s'", models.ProjectInfoDefaultGoModFileName, goModFilePath)
	}

	// parse go.mod
	moduleName, err := checkCmdExtractModuleName(goModFilePath)
	if err != nil {
		return models.ProjectInfo{}, fmt.Errorf("failed get module name: %s", err)
	}

	return models.ProjectInfo{
		Directory:      projectPath,
		GoArchFilePath: goArchFilePath,
		GoModFilePath:  goModFilePath,
		ModuleName:     moduleName,
	}, nil
}

func checkCmdExtractModuleName(goModPath string) (string, error) {
	goModFile, err := checkCmdParseGoModFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("can`t parse gomod: %v", err)
	}

	moduleName := goModFile.Module.Mod.Path
	if moduleName == "" {
		return "", fmt.Errorf("%s should contain module name in 'module'", models.ProjectInfoDefaultGoModFileName)
	}

	return moduleName, nil
}

func checkCmdParseGoModFile(path string) (*modfile.File, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read '%s': %v", path, err)
	}

	mod, err := modfile.ParseLax(path, file, nil)
	if err != nil {
		return nil, fmt.Errorf("modfile parseLax failed '%s': %v", path, err)
	}

	return mod, nil
}
