package assembler

import (
	"strings"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

type comparableComponent struct {
	name               models.ComponentName // unique name
	matchedCount       int                  // how much files is matched by this component
	concatMatchPattern string               // concat of all `component.in` in one string
}

func (a *Assembler) calculateFilesOwnage(components []*models.SpecComponent) {
	// a.go     -> [cmpA, cmpB]
	// dir/b.go -> [cmpA]
	// dir/c.go -> [cmpB, cmpA, cmpC]
	filePotentialOwners := make(map[models.PathRelative][]models.ComponentName)

	// Just map from input components array
	componentsMap := make(map[models.ComponentName]*models.SpecComponent)

	// populate utils maps
	for _, component := range components {
		componentsMap[component.Name.Value] = component

		for _, file := range component.MatchedFiles {
			if _, exist := filePotentialOwners[file]; !exist {
				filePotentialOwners[file] = make([]models.ComponentName, 0, 3)
			}

			filePotentialOwners[file] = append(filePotentialOwners[file], component.Name.Value)
		}
	}

	for file, potentialOwners := range filePotentialOwners {
		// find best owner
		owner := a.calculateFileOwner(potentialOwners, componentsMap)

		// add this file to owner
		componentsMap[owner].OwnedFiles = append(componentsMap[owner].OwnedFiles, file)
	}
}

func (a *Assembler) calculateFileOwner(
	potentialOwners []models.ComponentName,
	components map[models.ComponentName]*models.SpecComponent,
) models.ComponentName {
	// this file matcher only by one component
	if len(potentialOwners) == 1 {
		return potentialOwners[0]
	}

	cmpBestName := potentialOwners[0]
	cmpBest := comparableComponent{
		name:               cmpBestName,
		matchedCount:       len(components[cmpBestName].MatchedFiles),
		concatMatchPattern: concatMatchPatterns(components[cmpBestName]),
	}

	for i := 1; i < len(potentialOwners); i++ {
		candidateName := potentialOwners[i]
		candidate := comparableComponent{
			name:               candidateName,
			matchedCount:       len(components[candidateName].MatchedFiles),
			concatMatchPattern: concatMatchPatterns(components[candidateName]),
		}

		if compare(cmpBest, candidate) {
			// candidate is better
			cmpBestName = candidateName
			cmpBest = candidate
		}
	}

	return cmpBestName
}

func concatMatchPatterns(cmp *models.SpecComponent) string {
	var patters string

	for _, pattern := range cmp.MatchPatterns {
		patters += string(pattern.Value)
	}

	return patters
}

// should return true if B better than A
func compare(a, b comparableComponent) bool {
	if a.name == b.name {
		return false
	}

	// smallest files match count
	if b.matchedCount != a.matchedCount {
		return b.matchedCount < a.matchedCount
	}

	// has more specified directory
	aLen := strings.Count(a.concatMatchPattern, "/")
	bLen := strings.Count(b.concatMatchPattern, "/")
	if bLen != aLen {
		return bLen > aLen
	}

	// longest name
	if len(b.concatMatchPattern) != len(a.concatMatchPattern) {
		return len(b.concatMatchPattern) > len(a.concatMatchPattern)
	}

	// stable sort for equal priority path's
	return b.name < a.name
}
