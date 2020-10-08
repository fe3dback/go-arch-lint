package checker

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/fe3dback/go-arch-lint/path"
	"github.com/fe3dback/go-arch-lint/spec/archfile"

	"github.com/fe3dback/go-arch-lint/models"

	"github.com/fe3dback/go-arch-lint/spec"
)

const testModulePath = "github.com/fe3dback/go-arch-lint/checker/test"
const testArchFileV1 = ".go-arch-lint-v1.yml"

func makeTestProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename) + "/test"
}

func makeTestAbsPath(localPath string) string {
	return makeTestProjectRoot() + "/" + localPath
}

func makeTestResolvedPath(localPath string) *spec.ResolvedPath {
	return &spec.ResolvedPath{
		ImportPath: testModulePath + "/" + localPath,
		AbsPath:    makeTestAbsPath(localPath),
		LocalPath:  localPath,
	}
}

func makeTestResolvedProjectImport(localPath string) *models.ResolvedImport {
	return &models.ResolvedImport{
		Name:       testModulePath + "/" + localPath,
		ImportType: models.ImportTypeProject,
	}
}

func makeTestResolvedVendorImport(localPath string) *models.ResolvedImport {
	return &models.ResolvedImport{
		Name:       "github.com/vendor/lib/" + localPath,
		ImportType: models.ImportTypeVendor,
	}
}

func makeTestResolvedStdlibImport() *models.ResolvedImport {
	return &models.ResolvedImport{
		Name:       "fmt",
		ImportType: models.ImportTypeStdLib,
	}
}

func Test_checkImportPath(t *testing.T) {
	type args struct {
		componentImports []*spec.ResolvedPath
		resolvedImport   models.ResolvedImport
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "have needed import",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("needle"),
					makeTestResolvedPath("checker"),
					makeTestResolvedPath("path"),
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "not have needed import",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("checker"),
					makeTestResolvedPath("path"),
					makeTestResolvedPath("some"),
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "empty",
			args: args{
				componentImports: []*spec.ResolvedPath{},
				resolvedImport:   *makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "only needed",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("needle"),
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := &spec.Component{
				Name:           "component",
				AllowedImports: tt.args.componentImports,
			}

			if got := checkImportPath(cmp, tt.args.resolvedImport); got != tt.want {
				t.Errorf("checkImportPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkProjectImport(t *testing.T) {
	type args struct {
		componentImports []*spec.ResolvedPath
		componentFlags   *spec.SpecialFlags
		resolvedImport   models.ResolvedImport
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no flags, empty list",
			args: args{
				componentImports: []*spec.ResolvedPath{},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: false,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "project flag, empty list",
			args: args{
				componentImports: []*spec.ResolvedPath{},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: true,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "vendor flag, empty list",
			args: args{
				componentImports: []*spec.ResolvedPath{},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: false,
					AllowAllVendorDeps:  true,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "flag + list exactly same",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("needle"),
				},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: true,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "flag + list have needle",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("needle"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: true,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "flag + list not have needle",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: true,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "no flag + list have needle",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("needle"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: false,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "no flag + list not have needle",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: false,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := &spec.Component{
				Name:           "component",
				SpecialFlags:   tt.args.componentFlags,
				AllowedImports: tt.args.componentImports,
			}

			if got := checkProjectImport(cmp, tt.args.resolvedImport); got != tt.want {
				t.Errorf("checkProjectImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkVendorImport(t *testing.T) {
	type args struct {
		componentImports []*spec.ResolvedPath
		componentFlags   *spec.SpecialFlags
		resolvedImport   models.ResolvedImport
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no flags, empty list",
			args: args{
				componentImports: []*spec.ResolvedPath{},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: false,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "vendor flag, empty list",
			args: args{
				componentImports: []*spec.ResolvedPath{},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: false,
					AllowAllVendorDeps:  true,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "project flag, empty list",
			args: args{
				componentImports: []*spec.ResolvedPath{},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: true,
					AllowAllVendorDeps:  false,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "flag + list exactly same",
			args: args{
				componentImports: []*spec.ResolvedPath{
					makeTestResolvedPath("needle"),
				},
				componentFlags: &spec.SpecialFlags{
					AllowAllProjectDeps: false,
					AllowAllVendorDeps:  true,
				},
				resolvedImport: *makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := &spec.Component{
				Name:           "component",
				SpecialFlags:   tt.args.componentFlags,
				AllowedImports: tt.args.componentImports,
			}

			if got := checkVendorImport(cmp, tt.args.resolvedImport); got != tt.want {
				t.Errorf("checkVendorImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChecker_checkImport(t *testing.T) {
	type args struct {
		resolvedImport    models.ResolvedImport
		dependOnAnyVendor bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "allow any vendor / vendor",
			args: args{
				resolvedImport:    *makeTestResolvedVendorImport("needle"),
				dependOnAnyVendor: true,
			},
			want: true,
		},
		{
			name: "not allow any vendor / vendor",
			args: args{
				resolvedImport:    *makeTestResolvedVendorImport("needle"),
				dependOnAnyVendor: false,
			},
			want: false,
		},
		{
			name: "allow any vendor / project dep",
			args: args{
				resolvedImport:    *makeTestResolvedProjectImport("needle"),
				dependOnAnyVendor: true,
			},
			want: false,
		},
		{
			name: "stdlib always ok",
			args: args{
				resolvedImport:    *makeTestResolvedStdlibImport(),
				dependOnAnyVendor: false,
			},
			want: true,
		},
		{
			name: "project reject",
			args: args{
				resolvedImport:    *makeTestResolvedProjectImport("needle"),
				dependOnAnyVendor: false,
			},
			want: false,
		},
		{
			name: "vendor reject",
			args: args{
				resolvedImport:    *makeTestResolvedVendorImport("needle"),
				dependOnAnyVendor: false,
			},
			want: false,
		},
	}

	cmp := &spec.Component{
		Name: "component",
		SpecialFlags: &spec.SpecialFlags{
			AllowAllProjectDeps: false,
			AllowAllVendorDeps:  false,
		},
		AllowedImports: []*spec.ResolvedPath{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkImport(cmp, tt.args.resolvedImport, tt.args.dependOnAnyVendor); got != tt.want {
				t.Errorf("checkImport() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("check unknown import type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		resolvedImport := models.ResolvedImport{
			Name:       "something",
			ImportType: 100,
		}

		_ = checkImport(cmp, resolvedImport, false)
	})
}

func Test_longestPathComponent(t *testing.T) {
	type args struct {
		matched map[string]string
	}
	tests := []struct {
		name      string
		args      args
		wantId    string
		expectNil bool
	}{
		{
			name:      "empty list",
			args:      args{},
			wantId:    "",
			expectNil: true,
		},
		{
			name: "one element",
			args: args{
				matched: map[string]string{
					makeTestAbsPath("cat"): "cat",
				},
			},
			wantId:    "cat",
			expectNil: false,
		},
		{
			name: "two elements, same len, expect first (because cat < dog in alpha order)",
			args: args{
				matched: map[string]string{
					makeTestAbsPath("cat"): "cat",
					makeTestAbsPath("dog"): "dog",
				},
			},
			wantId:    "cat",
			expectNil: false,
		},
		{
			name: "two elements, same len, expect second (because cat < dog in alpha order)",
			args: args{
				matched: map[string]string{
					makeTestAbsPath("dog"): "dog",
					makeTestAbsPath("cat"): "cat",
				},
			},
			wantId:    "cat",
			expectNil: false,
		},
		{
			name: "two elements, second len > first len, expect second",
			args: args{
				matched: map[string]string{
					makeTestAbsPath("cat"):   "cat",
					makeTestAbsPath("doggy"): "doggy",
				},
			},
			wantId:    "doggy",
			expectNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := make(map[string]*spec.Component, 0)

			for archPath, id := range tt.args.matched {
				list[archPath] = &spec.Component{
					Name: id,
				}
			}

			got := longestPathComponent(list)

			if got == nil {
				if !tt.expectNil {
					t.Errorf("longestPathComponent() = %v, want nil", got)
					return
				}

				return
			}

			want := tt.wantId
			if !reflect.DeepEqual(got.Name, want) {
				t.Errorf("longestPathComponent() = %v, want %v", got, want)
			}
		})
	}
}

func TestChecker_Check(t *testing.T) {
	projectFiles := []*models.ResolvedFile{
		{
			Path: makeTestAbsPath("a/file.go"),
			Imports: []models.ResolvedImport{
				*makeTestResolvedProjectImport("b/subA/target"), // accepted
				*makeTestResolvedProjectImport("b/subB/target"), // accepted
				*makeTestResolvedProjectImport("a"),             // rejected
				*makeTestResolvedProjectImport("utils/foo"),     // accepted
				*makeTestResolvedProjectImport("utils/bar"),     // accepted
			},
		},
		{
			Path: makeTestAbsPath("b/subA/target/file.go"),
			Imports: []models.ResolvedImport{
				*makeTestResolvedProjectImport("a"),                  // rejected
				*makeTestResolvedProjectImport("b/subA/target"),      // rejected
				*makeTestResolvedProjectImport("b/subB/target"),      // rejected
				*makeTestResolvedProjectImport("c/any/path/foo/bar"), // accepted
				*makeTestResolvedProjectImport("utils/foo"),          // accepted
				*makeTestResolvedProjectImport("utils/bar/var/any"),  // accepted
			},
		},
		{
			Path: makeTestAbsPath("b/subB/target/file.go"),
			Imports: []models.ResolvedImport{
				*makeTestResolvedProjectImport("a"),                  // rejected
				*makeTestResolvedProjectImport("b/subA/target"),      // rejected
				*makeTestResolvedProjectImport("b/subB/target"),      // rejected
				*makeTestResolvedProjectImport("c/any/path/foo/bar"), // accepted
				*makeTestResolvedProjectImport("utils/foo"),          // accepted
				*makeTestResolvedProjectImport("utils/bar/var/any"),  // accepted
			},
		},
		{
			Path: makeTestAbsPath("c/file.go"),
			Imports: []models.ResolvedImport{
				*makeTestResolvedProjectImport("a"),                  // rejected
				*makeTestResolvedProjectImport("b/subA/target"),      // rejected
				*makeTestResolvedProjectImport("b/subB/target"),      // rejected
				*makeTestResolvedProjectImport("c/any/path/foo/bar"), // rejected
				*makeTestResolvedProjectImport("utils/foo"),          // accepted
				*makeTestResolvedProjectImport("utils/bar"),          // accepted
			},
		},
		{
			Path:    makeTestAbsPath("d/unknown.go"), // rejected
			Imports: []models.ResolvedImport{},
		},
		{
			Path:    makeTestAbsPath("a/sub-unknown/unknown.go"), // rejected
			Imports: []models.ResolvedImport{},
		},
	}

	root := makeTestProjectRoot()

	archFilePath := root + "/" + testArchFileV1
	sourceCode, err := ioutil.ReadFile(archFilePath)
	if err != nil {
		t.Fatalf("failed read archfile '%s'", archFilePath)
	}

	yamlSpec, err := archfile.NewYamlSpec(sourceCode)

	arch, err := spec.NewArch(
		path.NewResolver(),
		yamlSpec,
		testModulePath,
		root,
	)
	if err != nil {
		t.Errorf(fmt.Sprintf("failed to make arch: %v", err))
		return
	}

	checker := NewChecker(
		root,
		arch,
		projectFiles,
	)

	// -----------------

	fmt.Printf("Check result's on virtual arch FS:\n")

	result := checker.Check()
	for _, warn := range result.NotMatchedWarnings() {
		fmt.Printf("    - virtual fs info: notMatched: '%+v'\n", warn.FileRelativePath)
	}
	for _, warn := range result.DependencyWarnings() {
		fmt.Printf("    - virtual fs info: dep: '%s' should depend on '%s'\n",
			warn.FileRelativePath,
			warn.ResolvedImportName,
		)
	}

	if result.TotalCount() != 13 {
		t.Errorf("Expected 13 errors, got %d", result.TotalCount())
	}
}
