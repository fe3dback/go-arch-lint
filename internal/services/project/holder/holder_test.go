package holder

import (
	"reflect"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/common"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
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
				packagePath: "/app/internal/glue/project/holder",
				path:        "/app/internal/glue/project/holder",
			},
			want: true,
		},
		{
			name: "subfolder",
			args: args{
				packagePath: "/app/internal/glue/project/holder",
				path:        "/app/internal/glue/project/holder/sub",
			},
			want: false,
		},
		{
			name: "subfolder 2",
			args: args{
				packagePath: "/app/internal/glue/project/holder",
				path:        "/app/internal/glue/project/holder/sub/b",
			},
			want: false,
		},
		{
			name: "lower 1",
			args: args{
				packagePath: "/app/internal/glue/project/holder",
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
		component   speca.Component
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "match",
			args: args{
				packagePath: "/app/internal/glue/project/package",
				component: speca.Component{
					ResolvedPaths: []common.Referable[models.ResolvedPath]{
						common.NewReferable(
							models.ResolvedPath{AbsPath: "/app/internal/glue/project/package"},
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
				packagePath: "/app/internal/glue/project/package",
				component: speca.Component{
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
				packagePath: "/app/internal/glue/project/package",
				component: speca.Component{
					ResolvedPaths: []common.Referable[models.ResolvedPath]{
						common.NewReferable(
							models.ResolvedPath{AbsPath: "/app/internal/glue/project/package/sub"},
							common.NewEmptyReference(),
						),
						common.NewReferable(
							models.ResolvedPath{AbsPath: "/app/internal/glue/project/package"},
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
		components []speca.Component
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
				components: []speca.Component{
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
			if got := componentsMatchesFile(tt.args.filePath, tt.args.components); !reflect.DeepEqual(got, tt.want) {
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
