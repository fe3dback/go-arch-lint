package validator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint/v4/internal/models"
)

func Test_CommonComponentsValidator(t *testing.T) {
	tests := []struct {
		name    string
		mutator func(config *models.Config)
		out     []string
	}{
		{
			name: "happy",
			mutator: func(config *models.Config) {
				config.CommonComponents = append(config.CommonComponents, nf(models.ComponentName("my-cmp")))
			},
			out: []string{
				"Common component 'my-cmp' is not known",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vld := NewCommonComponentsValidator()

			in := createValidatorIn()
			vCtx := &validationContext{
				conf: in.conf,
			}

			tt.mutator(&vCtx.conf)
			vld.Validate(vCtx)

			wantNotices := make([]models.Notice, 0, len(tt.out))
			for _, wantNoticeText := range tt.out {
				wantNotices = append(wantNotices, models.Notice{
					Message:   wantNoticeText,
					Reference: models.NewInvalidReference(),
				})
			}

			require.Equal(t, wantNotices, append(vCtx.notices, vCtx.missUsage...))
		})
	}
}
