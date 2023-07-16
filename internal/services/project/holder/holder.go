package holder

import (
	"path/filepath"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/arch"
)

type (
	Holder struct {
	}

	matchedComponent struct {
		id         string
		filesCount int
	}
)

func NewHolder() *Holder {
	return &Holder{}
}

func (h *Holder) HoldProjectFiles(files []models.ProjectFile, components []arch.Component) []models.FileHold {
	matchedCount := make(map[string]int)
	// example:
	// /** 			= 100
	// /a/** 		= 5
	// /a/b/** 		= 1

	mapping := make(map[string][]string)
	// example:
	// /main.go			= ["/"]
	// /a/src.go		= ["/", "/a"]
	// /a/b/src.go		= ["/", "/a", "/b"]

	backMapping := make(map[string]models.ProjectFile)
	for _, file := range files {
		backMapping[file.Path] = file

		if _, ok := mapping[file.Path]; !ok {
			mapping[file.Path] = make([]string, 0)
		}

		components := componentsMatchesFile(file.Path, components)
		for _, component := range components {
			if _, ok := matchedCount[component]; !ok {
				matchedCount[component] = 0
			}

			matchedCount[component]++
			mapping[file.Path] = append(mapping[file.Path], component)
		}
	}

	results := make([]models.FileHold, 0)
	for filePath, componentIDs := range mapping {
		if len(componentIDs) == 0 {
			results = append(results, models.FileHold{
				File:        backMapping[filePath],
				ComponentID: nil,
			})

			continue
		}

		defComponent := componentIDs[0]
		holder := matchedComponent{
			id:         defComponent,
			filesCount: matchedCount[defComponent],
		}

		if len(componentIDs) > 1 {
			for _, componentID := range componentIDs {
				variant := matchedComponent{
					id:         componentID,
					filesCount: matchedCount[componentID],
				}

				if compare(holder, variant) {
					holder = variant
				}
			}
		}

		results = append(results, models.FileHold{
			File:        backMapping[filePath],
			ComponentID: &holder.id,
		})
	}

	return results
}

// should return true if B better than A
func compare(a, b matchedComponent) bool {
	if a.id == b.id {
		return false
	}

	// smallest files match count
	if b.filesCount != a.filesCount {
		return b.filesCount < a.filesCount
	}

	// has more specified directory
	aLen := strings.Count(a.id, "/")
	bLen := strings.Count(b.id, "/")
	if bLen != aLen {
		return bLen > aLen
	}

	// longest name
	if len(b.id) != len(a.id) {
		return len(b.id) > len(a.id)
	}

	// stable sort for equal priority path's
	return b.id < a.id
}

func componentsMatchesFile(filePath string, components []arch.Component) []string {
	matched := make([]string, 0)
	packagePath := filepath.Dir(filePath)

	for _, component := range components {
		if componentMatchPackage(packagePath, component) {
			matched = append(matched, component.Name.Value)
		}
	}

	return matched
}

func componentMatchPackage(packagePath string, component arch.Component) bool {
	for _, componentDirectoryRef := range component.ResolvedPaths {
		resolvedPackagePath := componentDirectoryRef.Value.AbsPath
		if packageMathPath(packagePath, resolvedPackagePath) {
			return true
		}
	}

	return false
}

func packageMathPath(packagePath string, resolvedPackagePath string) bool {
	return packagePath == resolvedPackagePath
}
