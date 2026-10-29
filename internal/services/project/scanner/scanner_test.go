package scanner_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/tools/go/packages"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/services/project/scanner"
)

const (
	cmpPkg    = "cmp"
	fmtPkg    = "fmt"
	osPkg     = "os"
	iterPkg   = "iter"
	slicesPkg = "slices"
)

// TestScan_IgnoresUnstatableOutOfScopeFile covers a directory that can be listed
// but whose entries cannot be stat'ed (mode 0o400: read, no search). The files
// inside are out of scope anyway (not ".go"), so the scan must not care about
// them. Walking with filepath.Walk stat'ed every entry before the scope check
// and aborted the whole scan on the first permission error; walking with
// filepath.WalkDir never stats out-of-scope entries.
func TestScan_IgnoresUnstatableOutOfScopeFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions do not restrict stat on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions; cannot simulate an unstatable entry")
	}

	projectDir := t.TempDir()

	goFile := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main\n"), 0o600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	// Not excluded in the config: the walk descends into it and lists it fine,
	// but stat'ing what it contains is denied.
	unsearchableDir := filepath.Join(projectDir, "assets")
	if err := os.Mkdir(unsearchableDir, 0o755); err != nil {
		t.Fatalf("mkdir dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(unsearchableDir, ".env"), []byte("SECRET=1\n"), 0o600); err != nil {
		t.Fatalf("write out of scope file: %v", err)
	}
	if err := os.Chmod(unsearchableDir, 0o400); err != nil {
		t.Fatalf("chmod dir: %v", err)
	}
	// Restore perms so t.TempDir cleanup can remove it.
	t.Cleanup(func() { _ = os.Chmod(unsearchableDir, 0o755) })

	files, err := scanner.NewScanner().Scan(context.Background(), projectDir, "example.com/proj", nil, nil)
	if err != nil {
		t.Fatalf("scan should ignore the unstatable out of scope file, got error: %v", err)
	}

	if len(files) != 1 || files[0].Path != goFile {
		t.Fatalf("expected only %q to be scanned, got %+v", goFile, files)
	}
}

// TestScan_SkipsUnreadableExcludedDir reproduces the case where an excluded
// directory inside the project tree is unreadable (e.g. a root-owned local
// docker volume). Before the fix the walk descended into it, hit a
// permission-denied readdir error and aborted the whole scan. Excluded dirs
// must now be skipped before descending, so the scan succeeds and still
// returns the in-scope source files.
func TestScan_SkipsUnreadableExcludedDir(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions; cannot simulate an unreadable dir")
	}

	projectDir := t.TempDir()

	goFile := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main\n"), 0o600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	excludedDir := filepath.Join(projectDir, ".data")
	if err := os.Mkdir(excludedDir, 0o755); err != nil {
		t.Fatalf("mkdir excluded dir: %v", err)
	}
	// A subdir the walker would try to readdir into; made unreadable so that
	// descending into it fails, mirroring the root-owned postgres data dir.
	unreadable := filepath.Join(excludedDir, "pgdata")
	if err := os.Mkdir(unreadable, 0o000); err != nil {
		t.Fatalf("mkdir unreadable dir: %v", err)
	}
	// Restore perms so t.TempDir cleanup can remove it.
	t.Cleanup(func() { _ = os.Chmod(unreadable, 0o755) })

	excludePaths := []models.ResolvedPath{{
		LocalPath: ".data",
		AbsPath:   excludedDir,
	}}

	files, err := scanner.NewScanner().Scan(context.Background(), projectDir, "example.com/proj", excludePaths, nil)
	if err != nil {
		t.Fatalf("scan should skip the excluded unreadable dir, got error: %v", err)
	}

	if len(files) != 1 || files[0].Path != goFile {
		t.Fatalf("expected only %q to be scanned, got %+v", goFile, files)
	}
}

// TestScan_DoesNotDescendIntoExcludedDir asserts that .go files living inside
// an excluded directory are never scanned.
func TestScan_DoesNotDescendIntoExcludedDir(t *testing.T) {
	projectDir := t.TempDir()

	kept := filepath.Join(projectDir, "keep.go")
	if err := os.WriteFile(kept, []byte("package main\n"), 0o600); err != nil {
		t.Fatalf("write kept file: %v", err)
	}

	excludedDir := filepath.Join(projectDir, "vendor")
	if err := os.Mkdir(excludedDir, 0o755); err != nil {
		t.Fatalf("mkdir excluded dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(excludedDir, "dep.go"), []byte("package vendor\n"), 0o600); err != nil {
		t.Fatalf("write excluded file: %v", err)
	}

	excludePaths := []models.ResolvedPath{{
		LocalPath: "vendor",
		AbsPath:   excludedDir,
	}}

	files, err := scanner.NewScanner().Scan(context.Background(), projectDir, "example.com/proj", excludePaths, nil)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if len(files) != 1 || files[0].Path != kept {
		t.Fatalf("expected only %q to be scanned, got %+v", kept, files)
	}
}

func TestScanner_NewScanner(t *testing.T) {
	s := scanner.NewScanner()
	if s.StdPackagesCount() == 0 {
		t.Fatal("expected standard packages to be populated")
	}

	for _, expected := range []string{fmtPkg, osPkg, "strings", iterPkg, slicesPkg, "maps", cmpPkg} {
		if !s.HasStdPackage(expected) {
			t.Errorf("expected std package %q to be registered", expected)
		}
	}
}

func TestScanner_ToleratesMissingOptionalStdPackages(t *testing.T) {
	// Simulate an older Go runtime (e.g. Go 1.20) whose "std" expansion returns
	// only legacy packages and omits packages introduced in newer releases (e.g. "iter").
	legacyLoader := func(_ *packages.Config, _ ...string) ([]*packages.Package, error) {
		return []*packages.Package{
			{ID: fmtPkg, PkgPath: fmtPkg},
			{ID: osPkg, PkgPath: osPkg},
		}, nil
	}

	s := scanner.NewScannerWithLoader(legacyLoader)
	for _, expected := range []string{fmtPkg, osPkg, iterPkg, slicesPkg, "maps", cmpPkg, "unique", "weak"} {
		if !s.HasStdPackage(expected) {
			t.Errorf("expected std package %q to be resolved despite missing from loader output", expected)
		}
	}
}

func TestScanner_ToleratesLoaderError(t *testing.T) {
	// Simulate an environment where the go toolchain cannot be invoked
	// (e.g. minimal container missing the go binary, or broken GOROOT).
	failingLoader := func(_ *packages.Config, _ ...string) ([]*packages.Package, error) {
		return nil, errors.New("exec: \"go\": executable file not found in $PATH")
	}

	s := scanner.NewScannerWithLoader(failingLoader)
	for _, expected := range []string{fmtPkg, osPkg, iterPkg, slicesPkg, cmpPkg} {
		if !s.HasStdPackage(expected) {
			t.Errorf("expected fallback std package %q to be resolved after loader error", expected)
		}
	}
}

func TestScanner_ToleratesPackageErrors(t *testing.T) {
	// Simulate loader output containing package-level errors for an optional or missing package.
	errorReportingLoader := func(_ *packages.Config, _ ...string) ([]*packages.Package, error) {
		return []*packages.Package{
			{ID: fmtPkg, PkgPath: fmtPkg},
			{
				ID:      iterPkg,
				PkgPath: iterPkg,
				Errors: []packages.Error{
					{Msg: "cannot find package iter in any of: /usr/local/go/src/iter (from $GOROOT)"},
				},
			},
		}, nil
	}

	s := scanner.NewScannerWithLoader(errorReportingLoader)
	if !s.HasStdPackage(fmtPkg) {
		t.Error("expected fmt to be registered")
	}
	if !s.HasStdPackage(iterPkg) {
		t.Error("expected iter to be registered despite package error")
	}
}

func TestScan_ClassifiesOptionalAndNewerStdPackages(t *testing.T) {
	projectDir := t.TempDir()
	sourceFile := filepath.Join(projectDir, "main.go")

	sourceCode := `package main

import (
	"cmp"
	"fmt"
	"iter"
	"slices"

	"example.com/proj/internal/foo"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Println
var _ = iter.Seq[int](nil)
var _ = slices.Sort[[]int]
var _ = cmp.Compare(1, 2)
`

	if err := os.WriteFile(sourceFile, []byte(sourceCode), 0o600); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	files, err := scanner.NewScanner().Scan(context.Background(), projectDir, "example.com/proj", nil, nil)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	expectedTypes := map[string]models.ImportType{
		cmpPkg:                               models.ImportTypeStdLib,
		fmtPkg:                               models.ImportTypeStdLib,
		iterPkg:                              models.ImportTypeStdLib,
		slicesPkg:                            models.ImportTypeStdLib,
		"example.com/proj/internal/foo":      models.ImportTypeProject,
		"github.com/stretchr/testify/assert": models.ImportTypeVendor,
	}

	for _, imp := range files[0].Imports {
		expected, exists := expectedTypes[imp.Name]
		if !exists {
			t.Errorf("unexpected import %q in scan results", imp.Name)
			continue
		}
		if imp.ImportType != expected {
			t.Errorf("import %q classified as %d, expected %d", imp.Name, imp.ImportType, expected)
		}
	}
}
