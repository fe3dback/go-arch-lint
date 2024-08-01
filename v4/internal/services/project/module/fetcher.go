package module

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type Fetcher struct {
}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) Fetch() (models.ProjectInfo, error) {
	// todo:
	return models.ProjectInfo{
		Directory:  "/home/neo/code/fe3dback/go-arch-lint/v4",
		ConfigPath: "/home/neo/code/fe3dback/go-arch-lint/v4/.go-arch-lint.yml",
		Module:     "github.com/fe3dback/go-arch-lint/v4",
	}, nil
}
