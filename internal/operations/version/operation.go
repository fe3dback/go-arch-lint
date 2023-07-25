package version

import (
	"fmt"
	"runtime/debug"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Operation struct {
	version    string
	buildTime  string
	commitHash string
}

func NewOperation(
	version string,
	buildTime string,
	commitHash string,
) *Operation {
	return &Operation{
		version:    version,
		buildTime:  buildTime,
		commitHash: commitHash,
	}
}

func (o *Operation) Behave() (models.CmdVersionOut, error) {
	if o.version == models.UnknownVersion {
		if data, err := o.fromCompiledMeta(); err == nil {
			return data, nil
		}
	}

	return o.fromLdFlags()
}

func (o *Operation) fromCompiledMeta() (models.CmdVersionOut, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return models.CmdVersionOut{}, fmt.Errorf("nothing to read")
	}

	if bi.Main.Version == "" {
		return models.CmdVersionOut{}, fmt.Errorf("nothing to read")
	}

	vcsHash := "unknown"
	vcsTime := "unknown"

	for _, setting := range bi.Settings {
		if setting.Key == "vcs.revision" {
			vcsHash = setting.Value
			continue
		}

		if setting.Key == "vcs.time" {
			vcsTime = setting.Value
			continue
		}
	}

	return models.CmdVersionOut{
		GoArchFileSupported: o.supportedSchemas(),
		LinterVersion:       bi.Main.Version,
		BuildTime:           vcsTime,
		CommitHash:          vcsHash,
	}, nil
}

func (o *Operation) fromLdFlags() (models.CmdVersionOut, error) {
	return models.CmdVersionOut{
		GoArchFileSupported: o.supportedSchemas(),
		LinterVersion:       o.version,
		BuildTime:           o.buildTime,
		CommitHash:          o.commitHash,
	}, nil
}

func (o *Operation) supportedSchemas() string {
	return fmt.Sprintf("%d .. %d", models.SupportedVersionMin, models.SupportedVersionMax)
}
