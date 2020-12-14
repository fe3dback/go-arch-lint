package check

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

const (
	flagMaxWarnings = "max-warnings"
	flagProjectPath = "project-path"
	flagArchFile    = "arch-file"
)

const (
	goModFileName       = "go.mod"
	defaultArchFileName = ".go-arch-lint.yml"
)

type (
	processorFn = func(models.FlagsCheck) error

	CommandAssembler struct {
		processorFn processorFn
		localFlags  *localFlags
	}

	localFlags struct {
		MaxWarnings int
		ProjectPath string
		ArchFile    string
	}
)

func NewCheckCommandAssembler(processorFn processorFn) *CommandAssembler {
	return &CommandAssembler{
		processorFn: processorFn,
		localFlags: &localFlags{
			MaxWarnings: 512,
			ProjectPath: "",
			ArchFile:    defaultArchFileName,
		},
	}
}

func (c *CommandAssembler) Assemble() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check",
		Short:   "check project architecture by yaml file",
		Long:    "compare project *.go files with arch defined in spec file",
		PreRunE: c.prePersist,
		RunE:    c.invoke,
	}

	c.assembleFlags(cmd)

	return cmd
}

func (c *CommandAssembler) invoke(_ *cobra.Command, _ []string) error {
	input, err := c.assembleInput()
	if err != nil {
		return fmt.Errorf("failed to assemble input params: %w", err)
	}

	return c.processorFn(input)
}

func (c *CommandAssembler) prePersist(cmd *cobra.Command, _ []string) error {
	rootDirectory, err := cmd.Flags().GetString(flagProjectPath)
	if err != nil {
		return failedToGetFlag(err, flagProjectPath)
	}

	archFile, err := cmd.Flags().GetString(flagArchFile)
	if err != nil {
		return failedToGetFlag(err, flagArchFile)
	}

	maxWarnings, err := cmd.Flags().GetInt(flagMaxWarnings)
	if err != nil {
		return failedToGetFlag(err, flagMaxWarnings)
	}

	const warningsRangeMin = 1
	const warningsRangeMax = 32768

	if maxWarnings < warningsRangeMin || maxWarnings > warningsRangeMax {
		return fmt.Errorf(
			"flag '%s' should by in range [%d .. %d]",
			flagMaxWarnings,
			warningsRangeMin,
			warningsRangeMax,
		)
	}

	// assemble localFlags
	c.localFlags.ProjectPath = rootDirectory
	c.localFlags.ArchFile = archFile
	c.localFlags.MaxWarnings = maxWarnings

	return nil
}

func (c *CommandAssembler) assembleInput() (models.FlagsCheck, error) {
	if c.localFlags.ProjectPath == "" {
		return models.FlagsCheck{}, fmt.Errorf("flag '%s' should by set", flagProjectPath)
	}

	projectPath, err := filepath.Abs(c.localFlags.ProjectPath)
	if err != nil {
		return models.FlagsCheck{}, fmt.Errorf("failed to resolve abs path '%s'", c.localFlags.ProjectPath)
	}

	// check arch file
	settingsGoArchFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, c.localFlags.ArchFile))
	_, err = os.Stat(settingsGoArchFilePath)
	if os.IsNotExist(err) {
		return models.FlagsCheck{}, fmt.Errorf("not found archfile in '%s'", settingsGoArchFilePath)
	}

	// check go.mod
	settingsGoModFilePath := filepath.Clean(fmt.Sprintf("%s/%s", projectPath, goModFileName))
	_, err = os.Stat(settingsGoModFilePath)
	if os.IsNotExist(err) {
		return models.FlagsCheck{}, fmt.Errorf("not found project '%s' in '%s'", goModFileName, settingsGoModFilePath)
	}

	// parse go.mod
	moduleName, err := checkCmdExtractModuleName(settingsGoModFilePath)
	if err != nil {
		return models.FlagsCheck{}, fmt.Errorf("failed get module name: %s", err)
	}

	return models.FlagsCheck{
		ProjectDirectory: projectPath,
		GoArchFilePath:   settingsGoArchFilePath,
		GoModFilePath:    settingsGoModFilePath,
		ModuleName:       moduleName,
		MaxWarnings:      c.localFlags.MaxWarnings,
	}, nil
}
