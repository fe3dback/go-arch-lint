package validator

import "github.com/fe3dback/go-arch-lint/spec/archfile"

type (
	checkerRegistry interface {
		applyChecker(path string, fn ArchFileValidatorFn)
		utils() *utils
		spec() archfile.YamlSpec
	}

	archFileCheckerRegistry struct {
		createdCheckers []ArchFileRuleChecker
		validatorUtils  *utils
		archFileSpec    archfile.YamlSpec
	}

	checkerFactoryFn func(checkerRegistry)
	factories        []checkerFactoryFn
)

func newArchFileCheckerRegistry(spec archfile.YamlSpec, validatorUtils *utils) *archFileCheckerRegistry {
	factories := factories{
		withCheckerCommonComponents,
		withCheckerCommonVendors,
		withCheckerComponents,
		withCheckerDependencies,
		withCheckerExcludedFiles,
		withCheckerVendors,
		withCheckerVersion,
	}

	registry := &archFileCheckerRegistry{
		createdCheckers: make([]ArchFileRuleChecker, 0),
		validatorUtils:  validatorUtils,
		archFileSpec:    spec,
	}

	for _, factory := range factories {
		factory(registry)
	}

	return registry
}

func (v *archFileCheckerRegistry) applyChecker(path string, fn ArchFileValidatorFn) {
	v.createdCheckers = append(v.createdCheckers, ArchFileRuleChecker{
		path:    path,
		checker: fn,
	})
}

func (v *archFileCheckerRegistry) utils() *utils {
	return v.validatorUtils
}

func (v *archFileCheckerRegistry) spec() archfile.YamlSpec {
	return v.archFileSpec
}
