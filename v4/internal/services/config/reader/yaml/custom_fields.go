package yaml

import "fmt"

type stringList []string

func (s *stringList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []string
	var lastErr error

	if err := unmarshal(&list); err == nil {
		*s = list
		return nil
	} else {
		lastErr = err
	}

	var value string
	if err := unmarshal(&value); err == nil {
		*s = []string{value}
		return nil
	} else {
		lastErr = fmt.Errorf("%v: %w", lastErr, err)
	}

	return fmt.Errorf("failed decode yaml stringsList: %w", lastErr)
}
