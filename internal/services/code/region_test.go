package code

import (
	"math"
	"reflect"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

func Test_calculateCodeRegion(t *testing.T) {
	type args struct {
		line         int
		regionHeight int
		maxLines     int
	}
	tests := []struct {
		name string
		args args
		want codeRegion
	}{
		{
			name: "6 even",
			args: args{
				line:         22,
				regionHeight: 6,
				maxLines:     37,
			},
			want: codeRegion{
				lineFirst: 19,
				lineMain:  22,
				lineLast:  25,
			},
		},
		{
			name: "5 odd",
			args: args{
				line:         22,
				regionHeight: 5,
				maxLines:     37,
			},
			want: codeRegion{
				lineFirst: 19,
				lineMain:  22,
				lineLast:  24,
			},
		},
		{
			name: "1",
			args: args{
				line:         1,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 1,
				lineMain:  1,
				lineLast:  4,
			},
		},
		{
			name: "2",
			args: args{
				line:         2,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 1,
				lineMain:  2,
				lineLast:  5,
			},
		},
		{
			name: "3",
			args: args{
				line:         3,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 1,
				lineMain:  3,
				lineLast:  6,
			},
		},
		{
			name: "4",
			args: args{
				line:         4,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 1,
				lineMain:  4,
				lineLast:  7,
			},
		},
		{
			name: "5",
			args: args{
				line:         5,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 2,
				lineMain:  5,
				lineLast:  8,
			},
		},
		{
			name: "-0",
			args: args{
				line:         100,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 97,
				lineMain:  100,
				lineLast:  100,
			},
		},
		{
			name: "-1",
			args: args{
				line:         99,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 96,
				lineMain:  99,
				lineLast:  100,
			},
		},
		{
			name: "-2",
			args: args{
				line:         98,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 95,
				lineMain:  98,
				lineLast:  100,
			},
		},
		{
			name: "-4",
			args: args{
				line:         96,
				regionHeight: 6,
				maxLines:     100,
			},
			want: codeRegion{
				lineFirst: 93,
				lineMain:  96,
				lineLast:  99,
			},
		},
		{
			name: "zero",
			args: args{
				line:         3,
				regionHeight: 0,
				maxLines:     10,
			},
			want: codeRegion{
				lineFirst: 3,
				lineMain:  3,
				lineLast:  3,
			},
		},
		{
			name: "one",
			args: args{
				line:         3,
				regionHeight: 1,
				maxLines:     10,
			},
			want: codeRegion{
				lineFirst: 2,
				lineMain:  3,
				lineLast:  3,
			},
		},
		{
			name: "two",
			args: args{
				line:         3,
				regionHeight: 2,
				maxLines:     10,
			},
			want: codeRegion{
				lineFirst: 2,
				lineMain:  3,
				lineLast:  4,
			},
		},
		{
			name: "small height",
			args: args{
				line:         2,
				regionHeight: 6,
				maxLines:     5,
			},
			want: codeRegion{
				lineFirst: 1,
				lineMain:  2,
				lineLast:  5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref := models.NewCodeReferenceRelative(
				models.Reference{
					Valid:  true,
					File:   "/tmp/dev",
					Line:   tt.args.line,
					Offset: 0,
				},
				int(math.Ceil(float64(tt.args.regionHeight)/2)),
				int(math.Floor(float64(tt.args.regionHeight)/2)),
			)

			if got := calculateCodeRegion(ref, tt.args.maxLines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateCodeRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}
