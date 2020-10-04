package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint/models"
	"golang.org/x/mod/modfile"
)

type checkCmdInput struct {
	settingsProjectDirectory string
	settingsGoArchFilePath   string
	settingsGoModFilePath    string
	settingsModuleName       string
}

func checkCmdAssembleCommandInput(flags *rootInput) checkCmdInput {
	if flags.projectRootDirectory == "" {
		panic(models.NewUserSpaceError(fmt.Sprintf("flag '%s' should by set", flagNameProjectPath)))
	}

	projectPath := filepath.Clean(flags.projectRootDirectory)

	// check arch file
	settingsGoArchFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, flags.archFile))
	_, err := os.Stat(settingsGoArchFilePath)
	if os.IsNotExist(err) {
		panic(models.NewUserSpaceError(fmt.Sprintf("not found archfile in '%s'", settingsGoArchFilePath)))
	}

	// check go.mod
	settingsGoModFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, goModFileName))
	_, err = os.Stat(settingsGoModFilePath)
	if os.IsNotExist(err) {
		panic(models.NewUserSpaceError(fmt.Sprintf("not found project '%s' in '%s'", goModFileName, settingsGoModFilePath)))
	}

	// parse go.mod
	moduleName, err := checkCmdExtractModuleName(settingsGoModFilePath)
	if err != nil {
		panic(models.NewUserSpaceError(fmt.Sprintf("failed get module name: %s", err)))
	}

	// assemble command input
	return checkCmdInput{
		settingsProjectDirectory: projectPath,
		settingsGoArchFilePath:   settingsGoArchFilePath,
		settingsGoModFilePath:    settingsGoModFilePath,
		settingsModuleName:       moduleName,
	}
}

func checkCmdExtractModuleName(goModPath string) (string, error) {
	goModFile, err := checkCmdParseGoModFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("can`t parse gomod: %v", err)
	}

	moduleName := goModFile.Module.Mod.Path
	if moduleName == "" {
		return "", fmt.Errorf("%s should contain module name in 'module'", goModFileName)
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
