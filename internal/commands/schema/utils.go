package schema

import (
	"fmt"
)

func failedToGetFlag(err error, flagName string) error {
	return fmt.Errorf("can`t get flag '%s': %w", flagName, err)
}
