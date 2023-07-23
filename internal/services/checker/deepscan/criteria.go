package deepscan

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/mod/modfile"
)

type Criteria struct {
	// required
	// package to check
	// example: /home/neo/go/src/example.com/neo/project/internal/operation/hello
	packagePath string

	// optional (if not set, moduleRootPath will be used)
	// directory for recursive usage analyse
	// example: /home/neo/go/src/example.com/neo/project/internal
	analyseScope string

	// optional (will be auto located, where the closest go.mod placed)
	// full path to project, where go.mod is located
	// example: /home/neo/go/src/example.com/neo/project
	moduleRootPath string

	// optional (will be parsed from go.mod)
	// exactly as in go.mod
	// example: example.com/neo/project
	moduleName string

	// optional (default empty)
	// excluded absolute path's from analyseScope
	excludePaths []string

	// optional (default empty)
	// exclude regexp matchers, for each file in analyseScope
	// all rejected files will not be parsed
	excludeFileMatchers []*regexp.Regexp
}

type CriteriaArg = func(*Criteria)

// NewCriteria build search Criteria for analyse
func NewCriteria(args ...CriteriaArg) (Criteria, error) {
	// init
	criteria := Criteria{}
	for _, builder := range args {
		builder(&criteria)
	}

	// check required fields
	if criteria.packagePath == "" {
		return Criteria{}, fmt.Errorf("criteria packagePath should be set")
	}

	// set up optional fields
	err := fillDefaultCriteriaFields(&criteria)
	if err != nil {
		return Criteria{}, fmt.Errorf("failed fill optional criteria fields: %w", err)
	}

	// ok
	return criteria, nil
}

// WithPackagePath set full abs path to go package
// who will be analysed
func WithPackagePath(path string) CriteriaArg {
	return func(criteria *Criteria) {
		criteria.packagePath = path
	}
}

// WithAnalyseScope set full abs path to some project directory
// can be project root, or some child directory, like 'internal'
// only this directory and it`s child will be analysed for params
func WithAnalyseScope(scope string) CriteriaArg {
	return func(criteria *Criteria) {
		criteria.analyseScope = scope
	}
}

// WithExcludedPath define list of abs path directories
// for exclude from analyse scope
func WithExcludedPath(paths []string) CriteriaArg {
	return func(criteria *Criteria) {
		criteria.excludePaths = paths
	}
}

// WithExcludedFileMatchers define list of regexp matchers
// that will match each file name in analyse scope
// if regexp match file name, it will be analysed for params
func WithExcludedFileMatchers(matchers []*regexp.Regexp) CriteriaArg {
	return func(criteria *Criteria) {
		criteria.excludeFileMatchers = matchers
	}
}

func fillDefaultCriteriaFields(criteria *Criteria) error {
	if criteria.moduleName == "" || criteria.moduleRootPath == "" {
		moduleName, root, err := findRootPath(criteria.packagePath)
		if err != nil {
			return fmt.Errorf("failed find root path of '%s': %w", criteria.packagePath, err)
		}

		criteria.moduleName = moduleName
		criteria.moduleRootPath = root
	}

	if criteria.analyseScope == "" {
		criteria.analyseScope = criteria.moduleRootPath
	}

	return nil
}

func findRootPath(packagePath string) (moduleName string, rootPath string, err error) {
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("packagePath directory not exist")
	}

	goModPath := filepath.Join(packagePath, "go.mod")
	_, err = os.Stat(goModPath)
	if err != nil {
		if os.IsNotExist(err) {
			// try find one level upper
			upperPath := filepath.Dir(packagePath)
			if upperPath == string(filepath.Separator) {
				return "", "", fmt.Errorf("go.mod not found on all parent levels")
			}

			return findRootPath(upperPath)
		}

		return "", "", fmt.Errorf("failed stat '%s': %w", goModPath, err)
	}

	file, err := os.ReadFile(goModPath)
	if err != nil {
		return "", "", fmt.Errorf("failed read '%s': %w", goModPath, err)
	}

	mod, err := modfile.ParseLax(goModPath, file, nil)
	if err != nil {
		return "", "", fmt.Errorf("modfile parse failed '%s': %w", goModPath, err)
	}

	return mod.Module.Mod.Path, packagePath, nil
}
