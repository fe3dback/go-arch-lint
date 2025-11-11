package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fe3dback/go-arch-lint/internal/services/checker/deepscan"
)

func Test1(t *testing.T) {
	// assemble
	_, callerDir, _, _ := runtime.Caller(0)
	fmt.Println("called from: " + callerDir)

	projectDir := filepath.Join(filepath.Dir(callerDir), "project")
	fmt.Println("project root dir: " + projectDir)

	searcher := deepscan.NewSearcher()
	criteria, err := deepscan.NewCriteria(
		deepscan.WithPackagePath(filepath.Join(projectDir, "internal", "operations")),
		deepscan.WithAnalyseScope(filepath.Join(projectDir, "internal")),
	)
	assert.NoError(t, err)

	// act
	expected := test1Expected()
	actual, err := searcher.Usages(criteria)

	// assert
	assert.NoError(t, err)

	// todo: write tests for this test project
	_, _ = expected, actual
	// assert.Equal(t, expected, actual)
}

func test1Expected() []deepscan.InjectionMethod {
	return []deepscan.InjectionMethod{
		{
			Name:       "hello",
			Definition: deepscan.Source{},
			Gates:      nil,
		},
	}
}
