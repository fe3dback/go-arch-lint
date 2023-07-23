package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	deepscan2 "github.com/fe3dback/go-arch-lint/internal/services/checker/deepscan"
	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	// assemble
	_, callerDir, _, _ := runtime.Caller(0)
	fmt.Println("called from: " + callerDir)

	projectDir := filepath.Join(filepath.Dir(callerDir), "project")
	fmt.Println("project root dir: " + projectDir)

	searcher := deepscan2.NewSearcher()
	criteria, err := deepscan2.NewCriteria(
		deepscan2.WithPackagePath(filepath.Join(projectDir, "internal", "operations")),
		deepscan2.WithAnalyseScope(filepath.Join(projectDir, "internal")),
	)

	// act
	expected := test1Expected()
	actual, err := searcher.Usages(criteria)

	// assert
	assert.NoError(t, err)

	// todo: write tests for this test project
	_, _ = expected, actual
	// assert.Equal(t, expected, actual)
}

func test1Expected() []deepscan2.InjectionMethod {
	return []deepscan2.InjectionMethod{
		{
			Name:       "hello",
			Definition: deepscan2.Source{},
			Gates:      nil,
		},
	}
}
