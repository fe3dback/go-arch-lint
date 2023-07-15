package schema

import (
	_ "embed"
	"fmt"
)

//go:embed v1.json
var v1 []byte

//go:embed v2.json
var v2 []byte

//go:embed v3.json
var v3 []byte

type Provider struct {
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) Provide(version int) ([]byte, error) {
	switch version {
	case 3:
		return v3, nil
	case 2:
		return v2, nil
	case 1:
		return v1, nil
	default:
		return nil, fmt.Errorf("unknown version: %d", version)
	}
}
