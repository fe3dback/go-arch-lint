package spec

import (
	"fmt"
)

type stringsList struct {
	list          []string
	definedAsList bool
}

func (s *stringsList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []string
	var lastErr error

	if err := unmarshal(&list); err == nil {
		s.list = list
		s.definedAsList = true

		return nil
	} else {
		lastErr = err
	}

	var value string
	if err := unmarshal(&value); err == nil {
		s.list = []string{value}
		s.definedAsList = false

		return nil
	} else {
		lastErr = fmt.Errorf("%v: %w", lastErr, err)
	}

	return fmt.Errorf("failed decode yaml stringsList: %w", lastErr)
}
