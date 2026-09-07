package assembler

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

// fakePathResolver returns paths in os specific form, exactly like
// filepath.Glob/filepath.Walk do inside real path.Resolver.
type fakePathResolver struct {
	matches []string
}

func (f fakePathResolver) Resolve(_ string) ([]string, error) {
	resolved := make([]string, 0, len(f.matches))
	for _, match := range f.matches {
		resolved = append(resolved, filepath.FromSlash(match))
	}

	return resolved, nil
}

// see: https://github.com/fe3dback/go-arch-lint/issues/79
// glob matches are os specific ('C:\app\internal\pkg' on windows), but
// ImportPath/LocalPath is a part of go import path, and always use '/'.
func Test_resolveLocalGlobPath(t *testing.T) {
	const moduleName = "example.com/app"

	tests := []struct {
		name          string
		rootDirectory string
		matches       []string
		want          []models.ResolvedPath
	}{
		{
			name:          "nested glob match",
			rootDirectory: "/app",
			matches:       []string{"/app/internal/services/checker"},
			want: []models.ResolvedPath{
				{
					ImportPath: "example.com/app/internal/services/checker",
					LocalPath:  "internal/services/checker",
					AbsPath:    filepath.FromSlash("/app/internal/services/checker"),
				},
			},
		},
		{
			name:          "root itself",
			rootDirectory: "/app",
			matches:       []string{"/app"},
			want: []models.ResolvedPath{
				{
					ImportPath: "example.com/app",
					LocalPath:  "",
					AbsPath:    filepath.FromSlash("/app"),
				},
			},
		},
		{
			name:          "trailing separator in match",
			rootDirectory: "/app",
			matches:       []string{"/app/internal/"},
			want: []models.ResolvedPath{
				{
					ImportPath: "example.com/app/internal",
					LocalPath:  "internal",
					AbsPath:    filepath.FromSlash("/app/internal"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newResolver(
				fakePathResolver{matches: tt.matches},
				filepath.FromSlash(tt.rootDirectory),
				moduleName,
			)

			got, err := r.resolveLocalGlobPath("internal/**")
			if err != nil {
				t.Fatalf("resolveLocalGlobPath() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resolveLocalGlobPath() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
