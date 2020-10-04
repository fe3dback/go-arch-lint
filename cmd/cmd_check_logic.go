package cmd

import (
	"github.com/fe3dback/go-arch-lint/checker"
	"github.com/fe3dback/go-arch-lint/cmd/container"
	"github.com/fe3dback/go-arch-lint/spec/annotated_validator"
)

func checkCmdProcess(flags *rootInput, input checkCmdInput) payloadTypeCommandCheck {
	payload := payloadTypeCommandCheck{
		ExecutionWarnings:      []annotated_validator.YamlAnnotatedWarning{},
		ExecutionError:         "",
		ArchHasWarnings:        false,
		ArchWarningsDeps:       []checker.WarningDep{},
		ArchWarningsNotMatched: []checker.WarningNotMatched{},
	}

	di := container.NewContainer(
		input.settingsGoArchFilePath,
		input.settingsProjectDirectory,
		input.settingsModuleName,
		flags.useColors,
	)

	annotatedValidator := di.ProvideSpecAnnotatedValidator()
	warnings, err := annotatedValidator.Validate()
	if err != nil {
		payload.ExecutionError = err.Error()
		return payload
	}

	if len(warnings) > 0 {
		payload.ExecutionError = "arch file invalid syntax"
		payload.ExecutionWarnings = warnings
		return payload
	}

	archChecker := di.ProvideChecker()
	result := archChecker.Check()
	if result.IsOk() {
		return payload
	}

	cmdCheckAssembleWarnings(&payload, result, flags.maxWarnings)

	return payload
}

func cmdCheckAssembleWarnings(payload *payloadTypeCommandCheck, res checker.CheckResult, maxWarnings int) {
	outputCount := 0

	// append deps
	for _, dep := range res.DependencyWarnings() {
		if outputCount >= maxWarnings {
			break
		}

		payload.ArchWarningsDeps = append(payload.ArchWarningsDeps, dep)
		outputCount++
	}

	// append not matched
	for _, notMatched := range res.NotMatchedWarnings() {
		if outputCount >= maxWarnings {
			break
		}

		payload.ArchWarningsNotMatched = append(payload.ArchWarningsNotMatched, notMatched)
		outputCount++
	}

	if outputCount > 0 {
		payload.ArchHasWarnings = true
	}
}
