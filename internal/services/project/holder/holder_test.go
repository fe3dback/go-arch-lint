package holder

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

const (
	testPathHolder  = "/app/internal/glue/project/holder"
	testPathPackage = "/app/internal/glue/project/package"
)

func Test_packageMathPath(t *testing.T) {
	type args struct {
		packagePath string
		path        string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "exactly",
			args: args{
				packagePath: testPathHolder,
				path:        testPathHolder,
			},
			want: true,
		},
		{
			name: "subfolder",
			args: args{
				packagePath: testPathHolder,
				path:        "/app/internal/glue/project/holder/sub",
			},
			want: false,
		},
		{
			name: "subfolder 2",
			args: args{
				packagePath: testPathHolder,
				path:        "/app/internal/glue/project/holder/sub/b",
			},
			want: false,
		},
		{
			name: "lower 1",
			args: args{
				packagePath: testPathHolder,
				path:        "/app/internal/glue/project",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := packageMathPath(tt.args.packagePath, tt.args.path); got != tt.want {
				t.Errorf("packageMathPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_componentMatchPackage(t *testing.T) {
	type args struct {
		packagePath string
		component   arch.Component
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "match",
			args: args{
				packagePath: testPathPackage,
				component: arch.Component{
					ResolvedPaths: []common.Referable[models.ResolvedPath]{
						common.NewReferable(
							models.ResolvedPath{AbsPath: testPathPackage},
							common.NewEmptyReference(),
						),
					},
				},
			},
			want: true,
		},
		{
			name: "not match",
			args: args{
				packagePath: testPathPackage,
				component: arch.Component{
					ResolvedPaths: []common.Referable[models.ResolvedPath]{
						common.NewReferable(
							models.ResolvedPath{AbsPath: "/app/internal/glue/project/package/sub"},
							common.NewEmptyReference(),
						),
					},
				},
			},
			want: false,
		},
		{
			name: "any match",
			args: args{
				packagePath: testPathPackage,
				component: arch.Component{
					ResolvedPaths: []common.Referable[models.ResolvedPath]{
						common.NewReferable(
							models.ResolvedPath{AbsPath: "/app/internal/glue/project/package/sub"},
							common.NewEmptyReference(),
						),
						common.NewReferable(
							models.ResolvedPath{AbsPath: testPathPackage},
							common.NewEmptyReference(),
						),
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := componentMatchPackage(tt.args.packagePath, tt.args.component); got != tt.want {
				t.Errorf("componentMatchPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_componentsMatchesFile(t *testing.T) {
	type args struct {
		filePath   string
		components []arch.Component
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "s1",
			args: args{
				filePath: "/app/file.go",
				components: []arch.Component{
					{
						Name: common.NewReferable("A", common.NewEmptyReference()),
						ResolvedPaths: []common.Referable[models.ResolvedPath]{
							common.NewReferable(
								models.ResolvedPath{AbsPath: "/app"},
								common.NewEmptyReference(),
							),
						},
					},
					{
						Name: common.NewReferable("C", common.NewEmptyReference()),
						ResolvedPaths: []common.Referable[models.ResolvedPath]{
							common.NewReferable(
								models.ResolvedPath{AbsPath: "/app/sub"},
								common.NewEmptyReference(),
							),
						},
					},
					{
						Name: common.NewReferable("D", common.NewEmptyReference()),
						ResolvedPaths: []common.Referable[models.ResolvedPath]{
							common.NewReferable(
								models.ResolvedPath{AbsPath: "/"},
								common.NewEmptyReference(),
							),
						},
					},
					{
						Name: common.NewReferable("B", common.NewEmptyReference()),
						ResolvedPaths: []common.Referable[models.ResolvedPath]{
							common.NewReferable(
								models.ResolvedPath{AbsPath: "/app"},
								common.NewEmptyReference(),
							),
						},
					},
				},
			},
			want: []string{"A", "B"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// componentsMatchesFile use filepath.Dir inside, so both file and component
			// paths should be in os specific form (test cases are written in unix form)
			filePath := filepath.FromSlash(tt.args.filePath)
			components := osSpecificComponents(tt.args.components)

			if got := componentsMatchesFile(filePath, components); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("componentsMatchesFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compare(t *testing.T) {
	type args struct {
		a matchedComponent
		b matchedComponent
	}
	tests := []struct {
		name      string
		args      args
		bIsBetter bool
	}{
		{
			name: "count better A",
			args: args{
				a: matchedComponent{id: "A", filesCount: 3},
				b: matchedComponent{id: "B", filesCount: 4},
			},
			bIsBetter: false,
		},
		{
			name: "count better B",
			args: args{
				a: matchedComponent{id: "A", filesCount: 4},
				b: matchedComponent{id: "B", filesCount: 3},
			},
			bIsBetter: true,
		},
		{
			name: "more specified, better A",
			args: args{
				a: matchedComponent{id: "/a/b/c/d", filesCount: 3},
				b: matchedComponent{id: "/a/b/c", filesCount: 3},
			},
			bIsBetter: false,
		},
		{
			name: "more specified, better B",
			args: args{
				a: matchedComponent{id: "/a/b/c", filesCount: 3},
				b: matchedComponent{id: "/a/b/c/d", filesCount: 3},
			},
			bIsBetter: true,
		},
		{
			name: "longer name, better A",
			args: args{
				a: matchedComponent{id: "/a/b/aaaa", filesCount: 3},
				b: matchedComponent{id: "/a/b/bbb", filesCount: 3},
			},
			bIsBetter: false,
		},
		{
			name: "longer name, better B",
			args: args{
				a: matchedComponent{id: "/a/b/bbb", filesCount: 3},
				b: matchedComponent{id: "/a/b/aaaa", filesCount: 3},
			},
			bIsBetter: true,
		},
		{
			name: "stable sort, better A",
			args: args{
				a: matchedComponent{id: "/aaa", filesCount: 3},
				b: matchedComponent{id: "/bbb", filesCount: 3},
			},
			bIsBetter: false,
		},
		{
			name: "stable sort, better B",
			args: args{
				a: matchedComponent{id: "/bbb", filesCount: 3},
				b: matchedComponent{id: "/aaa", filesCount: 3},
			},
			bIsBetter: true,
		},
		{
			name: "equal, better always A",
			args: args{
				a: matchedComponent{id: "/file/src.go", filesCount: 3},
				b: matchedComponent{id: "/file/src.go", filesCount: 3},
			},
			bIsBetter: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.args.a, tt.args.b); got != tt.bIsBetter {
				t.Errorf("compare() = %v, want %v", got, tt.bIsBetter)
			}
		})
	}
}

func osSpecificComponents(components []arch.Component) []arch.Component {
	result := make([]arch.Component, 0, len(components))

	for _, component := range components {
		paths := make([]common.Referable[models.ResolvedPath], 0, len(component.ResolvedPaths))

		for _, resolvedPath := range component.ResolvedPaths {
			resolvedPath.Value.AbsPath = filepath.FromSlash(resolvedPath.Value.AbsPath)
			paths = append(paths, resolvedPath)
		}

		component.ResolvedPaths = paths
		result = append(result, component)
	}

	return result
}
