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

	// todo: max-levels options when "focus" option is used

	graphCode, err := s.buildGraph(spec, in)
	if err != nil {
		return models.Graph{}, fmt.Errorf("failed build graph: %w", err)
	}

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

func (s *Operation) buildGraph(spec speca.Spec, opts models.FlagsGraph) (string, error) {
	var buff bytes.Buffer
	whiteList, err := s.populateGraphWhitelist(spec, opts)
	if err != nil {
		return "", err
	}

	for _, cmp := range spec.Components {
		if _, visible := whiteList[cmp.Name.Value()]; !visible {
			continue
		}

		for _, dep := range cmp.MayDependOn {
			if _, visible := whiteList[dep.Value()]; !visible {
				continue
			}

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

	return buff.String(), nil
}

func (s *Operation) populateGraphWhitelist(spec speca.Spec, opts models.FlagsGraph) (map[string]struct{}, error) {
	if opts.Focus == "" {
		return s.populateGraphWhitelistAll(spec)
	}

	return s.populateGraphWhitelistFocused(spec, opts.Focus)
}

func (s *Operation) populateGraphWhitelistAll(spec speca.Spec) (map[string]struct{}, error) {
	whiteList := make(map[string]struct{}, len(spec.Components))

	for _, cmp := range spec.Components {
		whiteList[cmp.Name.Value()] = struct{}{}
	}

	return whiteList, nil
}

func (s *Operation) populateGraphWhitelistFocused(spec speca.Spec, focusCmpName string) (map[string]struct{}, error) {
	cmpMap := make(map[string]speca.Component)
	rootExist := false

	for _, cmp := range spec.Components {
		cmpMap[cmp.Name.Value()] = cmp

		if focusCmpName == cmp.Name.Value() {
			rootExist = true
		}
	}

	if !rootExist {
		return nil, fmt.Errorf("focused cmp %s is not defined", focusCmpName)
	}

	whiteList := make(map[string]struct{}, len(spec.Components))
	resolved := make(map[string]struct{}, 64)
	resolveList := make([]string, 0, 64)
	resolveList = append(resolveList, focusCmpName)

	for len(resolveList) > 0 {
		cmp := cmpMap[resolveList[0]]
		resolveList = resolveList[1:]
		fmt.Printf("resolve %s, left: %d\n", cmp.Name.Value(), len(resolveList))

		if _, alreadyResolved := resolved[cmp.Name.Value()]; alreadyResolved {
			continue
		}

		// cmp itself
		whiteList[cmp.Name.Value()] = struct{}{}

		// cmp deps
		for _, dep := range cmp.MayDependOn {
			whiteList[dep.Value()] = struct{}{}
			resolveList = append(resolveList, dep.Value())
		}

		// mark as resolved (for recursion check)
		resolved[cmp.Name.Value()] = struct{}{}
	}

	return whiteList, nil
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
