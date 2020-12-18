package checker

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

const (
	testModulePath = "github.com/fe3dback/go-arch-lint/checker/test"
)

func makeTestProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename) + "/test"
}

func makeTestAbsPath(localPath string) string {
	return makeTestProjectRoot() + "/" + localPath
}

func makeBool(b bool) speca.ReferableBool {
	return speca.NewReferableBool(b, speca.NewEmptyReference())
}

func makeTestResolvedPath(localPath string) speca.ReferableResolvedPath {
	return speca.NewReferableResolvedPath(
		models.ResolvedPath{
			ImportPath: testModulePath + "/" + localPath,
			AbsPath:    makeTestAbsPath(localPath),
			LocalPath:  localPath,
		},
		speca.NewEmptyReference(),
	)
}

func makeTestResolvedProjectImport(localPath string) models.ResolvedImport {
	return models.ResolvedImport{
		Name:       testModulePath + "/" + localPath,
		ImportType: models.ImportTypeProject,
	}
}

func makeTestResolvedVendorImport(localPath string) models.ResolvedImport {
	return models.ResolvedImport{
		Name:       "github.com/vendor/lib/" + localPath,
		ImportType: models.ImportTypeVendor,
	}
}

func makeTestResolvedStdlibImport() models.ResolvedImport {
	return models.ResolvedImport{
		Name:       "fmt",
		ImportType: models.ImportTypeStdLib,
	}
}

func Test_checkImportPath(t *testing.T) {
	type args struct {
		componentImports []speca.ReferableResolvedPath
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
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("needle"),
					makeTestResolvedPath("checker"),
					makeTestResolvedPath("path"),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "not have needed import",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("checker"),
					makeTestResolvedPath("path"),
					makeTestResolvedPath("some"),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "empty",
			args: args{
				componentImports: []speca.ReferableResolvedPath{},
				resolvedImport:   makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "only needed",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("needle"),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := speca.Component{
				Name:           speca.NewReferableString("component", speca.NewEmptyReference()),
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
		componentImports []speca.ReferableResolvedPath
		componentFlags   speca.SpecialFlags
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
				componentImports: []speca.ReferableResolvedPath{},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(false),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "project flag, empty list",
			args: args{
				componentImports: []speca.ReferableResolvedPath{},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(true),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "vendor flag, empty list",
			args: args{
				componentImports: []speca.ReferableResolvedPath{},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(false),
					AllowAllVendorDeps:  makeBool(true),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "flag + list exactly same",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("needle"),
				},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(true),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "flag + list have needle",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("needle"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(true),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "flag + list not have needle",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(true),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "no flag + list have needle",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("needle"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(false),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "no flag + list not have needle",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("some"),
					makeTestResolvedPath("module12"),
				},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(false),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := speca.Component{
				Name:           speca.NewReferableString("component", speca.NewEmptyReference()),
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
		componentImports []speca.ReferableResolvedPath
		componentFlags   speca.SpecialFlags
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
				componentImports: []speca.ReferableResolvedPath{},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(false),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "vendor flag, empty list",
			args: args{
				componentImports: []speca.ReferableResolvedPath{},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(false),
					AllowAllVendorDeps:  makeBool(true),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
		{
			name: "project flag, empty list",
			args: args{
				componentImports: []speca.ReferableResolvedPath{},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(true),
					AllowAllVendorDeps:  makeBool(false),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: false,
		},
		{
			name: "flag + list exactly same",
			args: args{
				componentImports: []speca.ReferableResolvedPath{
					makeTestResolvedPath("needle"),
				},
				componentFlags: speca.SpecialFlags{
					AllowAllProjectDeps: makeBool(false),
					AllowAllVendorDeps:  makeBool(true),
				},
				resolvedImport: makeTestResolvedProjectImport("needle"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := speca.Component{
				Name:           speca.NewReferableString("component", speca.NewEmptyReference()),
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
				resolvedImport:    makeTestResolvedVendorImport("needle"),
				dependOnAnyVendor: true,
			},
			want: true,
		},
		{
			name: "not allow any vendor / vendor",
			args: args{
				resolvedImport:    makeTestResolvedVendorImport("needle"),
				dependOnAnyVendor: false,
			},
			want: false,
		},
		{
			name: "allow any vendor / project dep",
			args: args{
				resolvedImport:    makeTestResolvedProjectImport("needle"),
				dependOnAnyVendor: true,
			},
			want: false,
		},
		{
			name: "stdlib always ok",
			args: args{
				resolvedImport:    makeTestResolvedStdlibImport(),
				dependOnAnyVendor: false,
			},
			want: true,
		},
		{
			name: "project reject",
			args: args{
				resolvedImport:    makeTestResolvedProjectImport("needle"),
				dependOnAnyVendor: false,
			},
			want: false,
		},
		{
			name: "vendor reject",
			args: args{
				resolvedImport:    makeTestResolvedVendorImport("needle"),
				dependOnAnyVendor: false,
			},
			want: false,
		},
	}

	cmp := speca.Component{
		Name: speca.NewReferableString("component", speca.NewEmptyReference()),
		SpecialFlags: speca.SpecialFlags{
			AllowAllProjectDeps: makeBool(false),
			AllowAllVendorDeps:  makeBool(false),
		},
		AllowedImports: []speca.ReferableResolvedPath{},
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
