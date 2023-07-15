package check

import (
	"context"
	"fmt"
	"sort"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
	terminal "github.com/fe3dback/span-terminal"
)

type (
	Operation struct {
		projectInfoAssembler projectInfoAssembler
		specAssembler        specAssembler
		specChecker          specChecker
		referenceRender      referenceRender
		highlightCodePreview bool
	}

	limiterResult struct {
		results      models.CheckResult
		omittedCount int
	}
)

func NewOperation(
	projectInfoAssembler projectInfoAssembler,
	specAssembler specAssembler,
	specChecker specChecker,
	referenceRender referenceRender,
	highlightCodePreview bool,
) *Operation {
	return &Operation{
		projectInfoAssembler: projectInfoAssembler,
		specAssembler:        specAssembler,
		specChecker:          specChecker,
		referenceRender:      referenceRender,
		highlightCodePreview: highlightCodePreview,
	}
}

func (s *Operation) Behave(ctx context.Context, in models.CmdCheckIn) (models.CmdCheckOut, error) {
	// track progress of this command, with stdout drawing
	terminal.CaptureOutput()
	defer terminal.ReleaseOutput()

	projectInfo, err := s.projectInfoAssembler.ProjectInfo(in.ProjectPath, in.ArchFile)
	if err != nil {
		return models.CmdCheckOut{}, fmt.Errorf("failed to assemble project info: %w", err)
	}

	spec, err := s.specAssembler.Assemble(projectInfo)
	if err != nil {
		return models.CmdCheckOut{}, fmt.Errorf("failed to assemble spec: %w", err)
	}

	result, err := s.specChecker.Check(ctx, spec)
	if err != nil {
		return models.CmdCheckOut{}, fmt.Errorf("failed to check project deps: %w", err)
	}

	limitedResult := s.limitResults(result, in.MaxWarnings)

	model := models.CmdCheckOut{
		ModuleName:             spec.ModuleName.Value(),
		DocumentNotices:        s.assembleNotice(spec.Integrity),
		ArchHasWarnings:        s.resultsHasWarnings(limitedResult.results),
		ArchWarningsDependency: limitedResult.results.DependencyWarnings,
		ArchWarningsMatch:      limitedResult.results.MatchWarnings,
		ArchWarningsDeepScan:   limitedResult.results.DeepscanWarnings,
		OmittedCount:           limitedResult.omittedCount,
	}

	if model.ArchHasWarnings || len(model.DocumentNotices) > 0 {
		// normal output with exit code 1
		return model, models.NewUserSpaceError("check not successful")
	}

	return model, nil
}

func (s *Operation) limitResults(result models.CheckResult, maxWarnings int) limiterResult {
	passCount := 0
	limitedResults := models.CheckResult{
		DependencyWarnings: []models.CheckArchWarningDependency{},
		MatchWarnings:      []models.CheckArchWarningMatch{},
		DeepscanWarnings:   []models.CheckArchWarningDeepscan{},
	}

	// append deps
	for _, notice := range result.DependencyWarnings {
		if passCount >= maxWarnings {
			break
		}

		limitedResults.DependencyWarnings = append(limitedResults.DependencyWarnings, notice)
		passCount++
	}

	// append not matched
	for _, notice := range result.MatchWarnings {
		if passCount >= maxWarnings {
			break
		}

		limitedResults.MatchWarnings = append(limitedResults.MatchWarnings, notice)
		passCount++
	}

	// append deep scan
	for _, notice := range result.DeepscanWarnings {
		if passCount >= maxWarnings {
			break
		}

		limitedResults.DeepscanWarnings = append(limitedResults.DeepscanWarnings, notice)
		passCount++
	}

	totalCount := 0 +
		len(result.DeepscanWarnings) +
		len(result.DependencyWarnings) +
		len(result.MatchWarnings)

	return limiterResult{
		results:      limitedResults,
		omittedCount: totalCount - passCount,
	}
}

func (s *Operation) resultsHasWarnings(result models.CheckResult) bool {
	if len(result.DependencyWarnings) > 0 {
		return true
	}

	if len(result.MatchWarnings) > 0 {
		return true
	}

	if len(result.DeepscanWarnings) > 0 {
		return true
	}

	return false
}

func (s *Operation) assembleNotice(integrity speca.Integrity) []models.CheckNotice {
	notices := make([]speca.Notice, 0)
	notices = append(notices, integrity.DocumentNotices...)

	results := make([]models.CheckNotice, 0)
	for _, notice := range notices {
		results = append(results, models.CheckNotice{
			Text:   fmt.Sprintf("%s", notice.Notice),
			File:   notice.Ref.File,
			Line:   notice.Ref.Line,
			Offset: notice.Ref.Offset,
			SourceCodePreview: s.referenceRender.SourceCode(
				models.NewCodeReferenceRelative(notice.Ref, 1, 1),
				s.highlightCodePreview,
			),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		sI := results[i]
		sJ := results[j]

		if sI.File == sJ.File {
			return sI.Line < sJ.Line
		}

		return sI.File < sJ.File
	})

	return results
}
