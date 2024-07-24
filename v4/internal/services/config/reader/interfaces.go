package reader

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	FileReader interface {
		ReadFile(path models.PathAbsolute) (models.Config, error)
	}
)
