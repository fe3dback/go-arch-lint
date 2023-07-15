package spec

import "github.com/fe3dback/go-arch-lint/internal/models"

type yamlDocumentPathResolver = func(yamlPath string) models.Reference
