package schema

import (
	"fmt"
)

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Provide(version int) (string, error) {
	switch version {
	case 3:
		return V3, nil
	case 2:
		return V2, nil
	case 1:
		return V1, nil
	default:
		return "", fmt.Errorf("unknown version: %d", version)
	}
}
