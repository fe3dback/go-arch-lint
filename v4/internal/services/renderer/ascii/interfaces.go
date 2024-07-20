package ascii

import "github.com/fe3dback/go-arch-lint/v4/internal/models"

//go:generate ../../../../bin/mockgen -source=$GOFILE -destination=mocks/mocks.go -package=mocks

type asciiColorizer interface {
	Colorize(color models.ColorName, text string) string
}
