package models

type FileRef struct {
	IsDir bool
	Path  PathAbsolute
}

type FileMatchQueryType string

const (
	FileMatchQueryTypeAll             FileMatchQueryType = "all"
	FileMatchQueryTypeOnlyFiles       FileMatchQueryType = "files"
	FileMatchQueryTypeOnlyDirectories FileMatchQueryType = "directories"
)
