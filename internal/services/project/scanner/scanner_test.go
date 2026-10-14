package scanner_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/services/project/scanner"
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

func TestNewScanner_DoesNotPanicWithoutGo(t *testing.T) {
	origPath := os.Getenv("PATH")
	origGoroot := os.Getenv("GOROOT")
	t.Setenv("PATH", t.TempDir())
	t.Setenv("GOROOT", t.TempDir())
	defer func() {
		os.Setenv("PATH", origPath)
		os.Setenv("GOROOT", origGoroot)
	}()

	var sc *scanner.Scanner
	requireNoPanic := func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("NewScanner panicked without go toolchain: %v", r)
			}
		}()
		sc = scanner.NewScanner()
	}
	requireNoPanic()

	if sc == nil {
		t.Fatal("expected scanner instance, got nil")
	}

	projectDir := t.TempDir()
	source := `package main

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"example.com/proj/pkg"
	"github.com/stretchr/testify/assert"
)
`
	filePath := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(filePath, []byte(source), 0o600); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	files, err := sc.Scan(context.Background(), projectDir, "example.com/proj", nil, nil)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	types := make(map[string]models.ImportType)
	for _, imp := range files[0].Imports {
		types[imp.Name] = imp.ImportType
	}

	expectedTypes := map[string]models.ImportType{
		"fmt":                               models.ImportTypeStdLib,
		"iter":                              models.ImportTypeStdLib,
		"maps":                              models.ImportTypeStdLib,
		"slices":                            models.ImportTypeStdLib,
		"example.com/proj/pkg":              models.ImportTypeProject,
		"github.com/stretchr/testify/assert": models.ImportTypeVendor,
	}

	for pkg, expected := range expectedTypes {
		actual, ok := types[pkg]
		if !ok {
			t.Errorf("expected import %q not found in scan results", pkg)
			continue
		}
		if actual != expected {
			t.Errorf("package %q: expected type %d, got %d", pkg, expected, actual)
		}
	}
}

func TestNewScanner_LoadsFromGoroot(t *testing.T) {
	fakeGoroot := t.TempDir()
	customPkgDir := filepath.Join(fakeGoroot, "src", "custompkg")
	if err := os.MkdirAll(customPkgDir, 0o755); err != nil {
		t.Fatalf("failed to create fake goroot pkg dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(customPkgDir, "pkg.go"), []byte("package custompkg\n"), 0o600); err != nil {
		t.Fatalf("failed to write fake pkg: %v", err)
	}

	origGoroot := os.Getenv("GOROOT")
	t.Setenv("GOROOT", fakeGoroot)
	defer os.Setenv("GOROOT", origGoroot)

	sc := scanner.NewScanner()

	projectDir := t.TempDir()
	source := `package main

import (
	"custompkg"
	"fmt"
	"example.com/proj/pkg"
)
`
	filePath := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(filePath, []byte(source), 0o600); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	files, err := sc.Scan(context.Background(), projectDir, "example.com/proj", nil, nil)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	types := make(map[string]models.ImportType)
	for _, imp := range files[0].Imports {
		types[imp.Name] = imp.ImportType
	}

	if types["custompkg"] != models.ImportTypeStdLib {
		t.Errorf("custompkg from GOROOT: expected type %d, got %d", models.ImportTypeStdLib, types["custompkg"])
	}
	if types["fmt"] != models.ImportTypeStdLib {
		t.Errorf("fmt: expected type %d, got %d", models.ImportTypeStdLib, types["fmt"])
	}
	if types["example.com/proj/pkg"] != models.ImportTypeProject {
		t.Errorf("example.com/proj/pkg: expected type %d, got %d", models.ImportTypeProject, types["example.com/proj/pkg"])
	}
}

func TestScan_ProjectWithoutDomainInModuleName(t *testing.T) {
	sc := scanner.NewScanner()
	projectDir := t.TempDir()
	source := `package main

import (
	"fmt"
	"myapp/pkg"
	"github.com/foo/bar"
)
`
	filePath := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(filePath, []byte(source), 0o600); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	files, err := sc.Scan(context.Background(), projectDir, "myapp", nil, nil)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	types := make(map[string]models.ImportType)
	for _, imp := range files[0].Imports {
		types[imp.Name] = imp.ImportType
	}

	if types["fmt"] != models.ImportTypeStdLib {
		t.Errorf("fmt: expected type %d, got %d", models.ImportTypeStdLib, types["fmt"])
	}
	if types["myapp/pkg"] != models.ImportTypeProject {
		t.Errorf("myapp/pkg: expected type %d, got %d", models.ImportTypeProject, types["myapp/pkg"])
	}
	if types["github.com/foo/bar"] != models.ImportTypeVendor {
		t.Errorf("github.com/foo/bar: expected type %d, got %d", models.ImportTypeVendor, types["github.com/foo/bar"])
	}
}
