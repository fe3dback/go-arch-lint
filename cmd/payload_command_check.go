package cmd

import (
	"github.com/fe3dback/go-arch-lint/checker"
	"github.com/fe3dback/go-arch-lint/spec/annotated_validator"
)

type payloadTypeCommandCheck struct {
	ExecutionWarnings      []annotated_validator.YamlAnnotatedWarning
	ExecutionError         string
	ArchHasWarnings        bool
	ArchWarningsDeps       []checker.WarningDep
	ArchWarningsNotMatched []checker.WarningNotMatched
}
