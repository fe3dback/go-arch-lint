package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

const (
	archFileName  = ".go-arch-lint.yml"
	goModFileName = "go.mod"
)

const (
	flagNamePath              = "project-path"
	flagNameMaxWarningsOutput = "max-warnings"
	flagNameColorOutput       = "color-output"
)

var flagPath = flag.String(flagNamePath, "", "Absolute path to go project, like '~/go/src/github.com/google/project'")
var flagMaxWarningsOutput = flag.Int(flagNameMaxWarningsOutput, 512, "Max number of warnings to output")
var flagColorOutput = flag.Bool(flagNameColorOutput, true, "Use colors stdout/stderr output")

type flags struct {
	pathProjectDirectory string
	pathGoModFile        string
	pathArchFile         string
	moduleName           string
	maxWarningsOutput    int
	useColorsOutput      bool
}

func flagsParse() (flags, error) {
	flag.Parse()
	flags := flags{}

	// collectors - path's
	flags.collectArchFilePath()

	// collectors - params
	flags.collectParams()

	// collectors - go mod
	err := flags.collectGoMod()
	if err != nil {
		return flags, fmt.Errorf("collect gomod error: %v", err)
	}

	// log
	flags.logCollected()
	return flags, nil
}

func (f *flags) collectArchFilePath() {
	f.pathProjectDirectory = *flagPath
	if f.pathProjectDirectory == "" {
		panic(fmt.Sprintf("flag '%s' should by set", flagNamePath))
	}

	f.pathProjectDirectory = filepath.Clean(f.pathProjectDirectory)

	// check arch file
	f.pathArchFile = filepath.Clean(fmt.Sprintf("%s/%s", f.pathProjectDirectory, archFileName))
	dd, err := os.Stat(f.pathArchFile)
	_ = dd
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("not found archfile in '%s'", f.pathArchFile))
	}

	// check go.mod
	f.pathGoModFile = filepath.Clean(fmt.Sprintf("%s/%s", f.pathProjectDirectory, goModFileName))
	_, err = os.Stat(f.pathGoModFile)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("not found project '%s' in '%s'", goModFileName, f.pathGoModFile))
	}
}

func (f *flags) collectParams() {
	const warningsRangeMin = 1
	const warningsRangeMax = 32768

	maxWarnings := *flagMaxWarningsOutput
	if maxWarnings < warningsRangeMin || maxWarnings > warningsRangeMax {
		panic(fmt.Sprintf("flag '%s' should by in range [%d .. %d]",
			flagNameMaxWarningsOutput,
			warningsRangeMin,
			warningsRangeMax,
		))
	}

	f.maxWarningsOutput = maxWarnings

	//
	f.useColorsOutput = *flagColorOutput
}

func (f *flags) collectGoMod() error {
	gomod, err := f.parseGoMod(f.pathGoModFile)
	if err != nil {
		return fmt.Errorf("can`t parse gomod: %v", err)
	}

	f.moduleName = gomod.Module.Mod.Path
	if f.moduleName == "" {
		return fmt.Errorf("%s should contain module name in 'module'", goModFileName)
	}

	return nil
}

func (f *flags) parseGoMod(path string) (*modfile.File, error) {
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

func (f *flags) logCollected() {
	fmt.Printf("used arch file: %s\n", f.pathArchFile)
	fmt.Printf("        module: %s\n", f.moduleName)
	fmt.Println("")
}
