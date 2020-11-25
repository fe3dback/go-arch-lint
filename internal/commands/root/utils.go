package root

import (
	"fmt"

	"github.com/spf13/cobra"
)

func failedToGetFlag(err error, flagName string) error {
	return fmt.Errorf("can`t get flag '%s': %w", flagName, err)
}

func assertCommandIsValid(cmd *cobra.Command) error {
	if cmd.PersistentPreRun != nil {
		return fmt.Errorf("root sub command '%s' should not have 'PersistentPreRun', "+
			"use 'PreRun' instead", cmd.Name(),
		)
	}

	if cmd.PersistentPreRunE != nil {
		return fmt.Errorf("root sub command '%s' should not have 'PersistentPreRunE', "+
			"use 'PreRunE' instead", cmd.Name(),
		)
	}

	return nil
}
