package yaml

import (
	"errors"
	"fmt"
)

type stringList []string

func (s *stringList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var list []string
	var lastErr error

	err := unmarshal(&list)
	if err != nil {
		lastErr = err
	} else {
		*s = list
		return nil
	}

	var value string
	err = unmarshal(&value)
	if err != nil {
		lastErr = errors.Join(lastErr, err)
	} else {
		*s = []string{value}
		return nil
	}

	return fmt.Errorf("failed decode yaml stringsList: %w", lastErr)
}

type optional[T any] struct {
	value   T
	defined bool
}

func (s *optional[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s.defined = true

	var value T
	err := unmarshal(&value)
	if err != nil {
		return err
	}

	s.value = value
	return nil
}
