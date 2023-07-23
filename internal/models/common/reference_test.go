package common_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

func TestReference_ClampWithRealLinesCount(t *testing.T) {
	type args struct {
		line         int
		regionHeight int
		maxLines     int
	}

	tests := []struct {
		name         string
		args         args
		linesCount   int
		wantLineFrom int
		wantLineMain int
		wantLineTo   int
	}{
		{
			name: "6 even",
			args: args{
				line:         22,
				regionHeight: 6,
				maxLines:     37,
			},
			wantLineFrom: 19,
			wantLineMain: 22,
			wantLineTo:   25,
		},
		{
			name: "5 odd",
			args: args{
				line:         22,
				regionHeight: 5,
				maxLines:     37,
			},
			wantLineFrom: 19,
			wantLineMain: 22,
			wantLineTo:   24,
		},
		{
			name: "1",
			args: args{
				line:         1,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 1,
			wantLineMain: 1,
			wantLineTo:   4,
		},
		{
			name: "2",
			args: args{
				line:         2,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 1,
			wantLineMain: 2,
			wantLineTo:   5,
		},
		{
			name: "3",
			args: args{
				line:         3,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 1,
			wantLineMain: 3,
			wantLineTo:   6,
		},
		{
			name: "4",
			args: args{
				line:         4,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 1,
			wantLineMain: 4,
			wantLineTo:   7,
		},
		{
			name: "5",
			args: args{
				line:         5,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 2,
			wantLineMain: 5,
			wantLineTo:   8,
		},
		{
			name: "-0",
			args: args{
				line:         100,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 97,
			wantLineMain: 100,
			wantLineTo:   100,
		},
		{
			name: "-1",
			args: args{
				line:         99,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 96,
			wantLineMain: 99,
			wantLineTo:   100,
		},
		{
			name: "-2",
			args: args{
				line:         98,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 95,
			wantLineMain: 98,
			wantLineTo:   100,
		},
		{
			name: "-4",
			args: args{
				line:         96,
				regionHeight: 6,
				maxLines:     100,
			},
			wantLineFrom: 93,
			wantLineMain: 96,
			wantLineTo:   99,
		},
		{
			name: "zero",
			args: args{
				line:         3,
				regionHeight: 0,
				maxLines:     10,
			},
			wantLineFrom: 3,
			wantLineMain: 3,
			wantLineTo:   3,
		},
		{
			name: "one",
			args: args{
				line:         3,
				regionHeight: 1,
				maxLines:     10,
			},
			wantLineFrom: 2,
			wantLineMain: 3,
			wantLineTo:   3,
		},
		{
			name: "two",
			args: args{
				line:         3,
				regionHeight: 2,
				maxLines:     10,
			},
			wantLineFrom: 2,
			wantLineMain: 3,
			wantLineTo:   4,
		},
		{
			name: "small height",
			args: args{
				line:         2,
				regionHeight: 6,
				maxLines:     5,
			},
			wantLineFrom: 1,
			wantLineMain: 2,
			wantLineTo:   5,
		},
		{
			name: "main upper bound",
			args: args{
				line:         15,
				regionHeight: 6,
				maxLines:     13,
			},
			wantLineFrom: 12,
			wantLineMain: 13,
			wantLineTo:   13,
		},
		{
			name: "main all bound",
			args: args{
				line:         15,
				regionHeight: 6,
				maxLines:     8,
			},
			wantLineFrom: 8,
			wantLineMain: 8,
			wantLineTo:   8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := common.NewReferenceSingleLine("/tmp/dev", tt.args.line, 0)
			ref = ref.ExtendRange(
				int(math.Ceil(float64(tt.args.regionHeight)/2)),
				int(math.Floor(float64(tt.args.regionHeight)/2)),
			)

			want := common.NewReferenceRange(
				"/tmp/dev",
				tt.wantLineFrom,
				tt.wantLineMain,
				tt.wantLineTo,
			)

			if got := ref.ClampWithRealLinesCount(tt.args.maxLines); !reflect.DeepEqual(got, want) {
				t.Errorf("ClampWithRealLinesCount() = %v, want %v", got, want)
			}
		})
	}
}
