package xpath_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
	"github.com/fe3dback/go-arch-lint/v4/internal/services/project/xpath"
)

const (
	fakeFilesProjectDirectory = "/home/u/project/"
)

// first block of *.gold files
// is yaml document that parsed with this struct
type goldenMeta struct {
	Query              string   `json:"query"`
	QueryType          string   `json:"type"`
	WorkDirectory      string   `json:"workdir"`
	FilterExtensions   []string `json:"ext"`
	ExcludeDirs        []string `json:"exclude_dirs"`
	ExcludeFilesRegexp []string `json:"exclude_files_regexp"`
	ExcludeFiles       []string `json:"exclude_file"`
}

// nolint
func TestHelper_FindProjectFiles(t *testing.T) {
	tests := []struct {
		gold string
	}{
		{gold: "relative_file_pick"},
		{gold: "relative_file_pick_glob"},
		{gold: "relative_directory"},
		{gold: "relative_files_in_dir_with_wd_not_found"},
		{gold: "relative_files_glob_in_dir_with_wd"},
		{gold: "relative_files_glob_in_dir_with_wd_filter_ext"},
		{gold: "relative_files_glob_in_any_sub_dir_with_wd"},
		{gold: "relative_files_glob_in_any_sub_dir_with_wd_with_filter_ext"},
		{gold: "relative_files_glob_sub_dir"},
		{gold: "relative_files_glob_sub_dir_any_include_self"},
		{gold: "relative_ddd_repos"},
	}
	for _, tt := range tests {
		t.Run(tt.gold, func(t *testing.T) {
			// load golden
			goldBytes, err := os.ReadFile(fmt.Sprintf("./tests/%s.gold", tt.gold))
			require.NoError(t, err)

			gold := string(goldBytes)
			documents := strings.Split(gold, "#### #### #### #### #### #### #### ####")
			require.Len(t, documents, 2, "*.gold files should have 2 sections, separated by '#'")

			sectionMeta := documents[0]
			sectionFakeFiles := documents[1]

			meta := goldenMeta{}
			err = yaml.Unmarshal([]byte(sectionMeta), &meta)
			require.NoError(t, err)

			// stub and mocks
			stubContent := createStubContent(sectionFakeFiles)
			scanner := newFileScannerTestStub(stubContent)

			// real matchers
			matcherRelative := xpath.NewMatcherRelative()
			matcherAbsolute := xpath.NewMatcherAbsolute(matcherRelative)
			matcherRelativeGlob := xpath.NewMatcherRelativeGlob()

			hlp := xpath.NewHelper(
				fakeFilesProjectDirectory,
				scanner,
				matcherRelative,
				matcherAbsolute,
				matcherRelativeGlob,
				nil, // todo:
			)

			// act
			got, err := hlp.FindProjectFiles(createQueryFromDocMeta(meta))
			require.NoError(t, err)

			// assert
			wantContent := make(map[models.PathAbsolute]any)
			for _, stub := range stubContent {
				if !stub.expected {
					continue
				}

				wantContent[stub.path] = struct{}{}
			}

			gotContent := make(map[models.PathAbsolute]any)
			for _, dst := range got {
				gotContent[dst.PathAbs] = struct{}{}
			}

			// check that all wanted to exist in gotten
			hasWarnings := false

			for wantPath := range wantContent {
				_, exists := gotContent[wantPath]

				if !assert.True(t, exists, fmt.Sprintf("path '%s' is not found by query (but expected to be)", wantPath)) {
					hasWarnings = true
				}
			}

			// check that all gotten is expected
			for gotPath := range gotContent {
				_, exists := wantContent[gotPath]

				if !assert.True(t, exists, fmt.Sprintf("found unexpected path '%s' by query", gotPath)) {
					hasWarnings = true
				}
			}

			if hasWarnings {
				fmt.Println("found:")
				for foundPath := range gotContent {
					fmt.Printf("- %s\n", foundPath)
				}
				fmt.Println("")
			}

			_ = got
		})
	}
}

func createStubContent(testFilesSection string) []filePathStub {
	lines := strings.Split(testFilesSection, "\n")

	content := make([]filePathStub, 0, len(lines))
	for _, line := range lines {
		isExpected := false

		if strings.HasPrefix(line, "> ") {
			isExpected = true
			line = strings.TrimPrefix(line, "> ")
		}

		if strings.Contains(line, "//") {
			parts := strings.Split(line, "//")
			line = parts[0]
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// this is wrong detector, but good enough for tests
		var isDir bool

		if strings.Contains(line, ".") {
			isDir = false
		} else {
			isDir = true
		}

		content = append(content, filePathStub{
			path:     models.PathAbsolute(line),
			isDir:    isDir,
			expected: isExpected,
		})
	}

	return content
}

// nolint
func createQueryFromDocMeta(meta goldenMeta) models.FileQuery {
	var queryType models.FileMatchQueryType

	switch meta.QueryType {
	case "files":
		queryType = models.FileMatchQueryTypeOnlyFiles
	case "directories":
		queryType = models.FileMatchQueryTypeOnlyDirectories
	case "all":
		queryType = models.FileMatchQueryTypeAll
	default:
		panic(fmt.Errorf("unknown query type %s", meta.QueryType))
	}

	var path any
	isGlob := strings.Contains(meta.Query, "*")

	if filepath.IsAbs(meta.Query) {
		if isGlob {
			path = models.PathAbsoluteGlob(meta.Query)
		} else {
			path = models.PathAbsolute(meta.Query)
		}
	} else {
		if isGlob {
			path = models.PathRelativeGlob(meta.Query)
		} else {
			path = models.PathRelative(meta.Query)
		}
	}

	var exts []string
	if len(meta.FilterExtensions) > 0 {
		exts = meta.FilterExtensions
	}

	var excludeDirs []models.PathRelative
	if len(meta.ExcludeDirs) > 0 {
		for _, dir := range meta.ExcludeDirs {
			excludeDirs = append(excludeDirs, models.PathRelative(dir))
		}
	}

	var excludeRegexp []models.PathRelativeRegExp
	if len(meta.ExcludeFilesRegexp) > 0 {
		for _, fileRegexp := range meta.ExcludeFilesRegexp {
			excludeRegexp = append(excludeRegexp, models.PathRelativeRegExp(fileRegexp))
		}
	}

	var excludeFiles []models.PathRelative
	if len(meta.ExcludeFiles) > 0 {
		for _, file := range meta.ExcludeFiles {
			excludeFiles = append(excludeFiles, models.PathRelative(file))
		}
	}

	return models.FileQuery{
		Path:               path,
		WorkingDirectory:   models.PathRelative(meta.WorkDirectory),
		Type:               queryType,
		ExcludeDirectories: excludeDirs,
		ExcludeFiles:       excludeFiles,
		ExcludeRegexp:      excludeRegexp,
		Extensions:         exts,
	}
}
