package root

import (
	"fmt"
)

func (c *CommandAssembler) failedToGetFlag(err error, flagName string) error {
	return fmt.Errorf("can`t get flag '%s': %w", flagName, err)
}
