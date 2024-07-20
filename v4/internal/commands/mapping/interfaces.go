package mapping

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

type (
	operation interface {
		Mapping(in models.CmdMappingIn) (models.CmdMappingOut, error)
	}
)
