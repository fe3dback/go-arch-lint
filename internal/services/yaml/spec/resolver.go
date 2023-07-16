package spec

import (
	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

type yamlDocumentPathResolver = func(yamlPath string) common.Reference
