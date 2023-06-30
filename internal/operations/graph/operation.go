package graph

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"

	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
)

type Operation struct {
	specAssembler SpecAssembler
}

func NewOperation(
	specAssembler SpecAssembler,
) *Operation {
	return &Operation{
		specAssembler: specAssembler,
	}
}

func (s *Operation) Behave(
	ctx context.Context,
	in models.FlagsGraph,
) (models.Graph, error) {
	spec, err := s.specAssembler.Assemble()
	if err != nil {
		return models.Graph{}, fmt.Errorf("failed to assemble spec: %w", err)
	}

	fmt.Println(in.IncludeVendors)
	graphCode := s.buildGraph(spec, in)
	svg, err := s.compileGraph(ctx, graphCode)
	if err != nil {
		return models.Graph{}, fmt.Errorf("failed to compile graph: %w", err)
	}

	err = os.WriteFile(in.OutFile, svg, os.ModePerm)
	if err != nil {
		return models.Graph{}, fmt.Errorf("failed write graph into '%s' file: %w", in.OutFile, err)
	}

	return models.Graph{
		ProjectDirectory: spec.RootDirectory.Value(),
		ModuleName:       spec.ModuleName.Value(),
		OutFile:          in.OutFile,
	}, nil
}

func (s *Operation) buildGraph(spec speca.Spec, opts models.FlagsGraph) string {
	var buff bytes.Buffer

	whiteList := map[string]struct{}{
		opts.Focus: {},
	}

	isVisible := func(cmp string) bool {
		if opts.Focus == "" {
			// all nodes visible
			return true
		}

		if _, ok := whiteList[cmp]; ok {
			return true
		}

		return false
	}

	applyToVisibleList := func(cmp string) {
		whiteList[cmp] = struct{}{}
	}

	for _, cmp := range spec.Components {
		if !isVisible(cmp.Name.Value()) {
			continue
		}

		for _, dep := range cmp.MayDependOn {
			applyToVisibleList(dep.Value())
			buff.WriteString(fmt.Sprintf("%s -> %s\n", cmp.Name.Value(), dep.Value()))
		}

		if opts.IncludeVendors {
			for _, vnd := range cmp.CanUse {
				vars := map[string]string{
					"vnd": vnd.Value(),
					"cmp": cmp.Name.Value(),
				}

				tpl := `
				{{vnd}}.style.font-size: 12
				{{cmp}} <- {{vnd}} {
				source-arrowhead: {
					shape: diamond
					style.filled: false
				  }
				}
				`

				for name, value := range vars {
					tpl = strings.ReplaceAll(tpl, fmt.Sprintf("{{%s}}", name), value)
				}
				buff.WriteString(tpl)
			}
		}
	}

	return buff.String()
}

func (s *Operation) compileGraph(ctx context.Context, graphCode string) ([]byte, error) {
	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return nil, fmt.Errorf("failed create ruler: %w", err)
	}

	diagram, _, err := d2lib.Compile(ctx, graphCode, &d2lib.CompileOptions{
		Layout: func(ctx context.Context, g *d2graph.Graph) error {
			return d2dagrelayout.Layout(ctx, g, nil)
		},
		Ruler: ruler,
	})
	if err != nil {
		return nil, fmt.Errorf("failed compile d2 graph: %w", err)
	}

	out, err := d2svg.Render(diagram, &d2svg.RenderOpts{
		Pad:     d2svg.DEFAULT_PADDING,
		Sketch:  true,
		ThemeID: d2themescatalog.NeutralDefault.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("svg render failed: %w", err)
	}

	return out, nil
}
