package check

import (
	"fmt"
	"sort"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

const highlightPreviewCodeLinesYAML = 1

type (
	Service struct {
		specAssembler        SpecAssembler
		specChecker          SpecChecker
		referenceRender      ReferenceRender
		highlightCodePreview bool
	}

	limiterResult struct {
		results      models.CheckResult
		omittedCount int
	}
)

func NewService(
	specAssembler SpecAssembler,
	specChecker SpecChecker,
	referenceRender ReferenceRender,
	highlightCodePreview bool,
) *Service {
	return &Service{
		specAssembler:        specAssembler,
		specChecker:          specChecker,
		referenceRender:      referenceRender,
		highlightCodePreview: highlightCodePreview,
	}
}

func (s *Service) Behave(maxWarnings int) (models.Check, error) {
	spec, err := s.specAssembler.Assemble()
	if err != nil {
		return models.Check{}, fmt.Errorf("failed to assemble spec: %w", err)
	}

	result, err := s.specChecker.Check(spec)
	if err != nil {
		return models.Check{}, fmt.Errorf("failed to check project deps: %w", err)
	}

	limitedResult := s.limitResults(result, maxWarnings)

	model := models.Check{
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

func (s *Service) limitResults(result models.CheckResult, maxWarnings int) limiterResult {
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
	const maxDeepScan = 10
	deepScanCount := 0

	for _, notice := range result.DeepscanWarnings {
		if passCount >= maxWarnings {
			break
		}
		if deepScanCount >= maxDeepScan {
			break
		}

		limitedResults.DeepscanWarnings = append(limitedResults.DeepscanWarnings, notice)
		passCount++
		deepScanCount++
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

func (s *Service) resultsHasWarnings(result models.CheckResult) bool {
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

func (s *Service) assembleNotice(integrity speca.Integrity) []models.CheckNotice {
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
				notice.Ref,
				highlightPreviewCodeLinesYAML,
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
