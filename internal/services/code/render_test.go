package code

import (
	"bytes"
	"io"
	"testing"
)

func Test_readLines(t *testing.T) {
	type args struct {
		r      io.Reader
		region codeRegion
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "simple",
			args: args{
				r: bytes.NewReader([]byte("Line1\nLine2\nLine3\nLine4\nLine5\nLine6")),
				region: codeRegion{
					lineFirst: 2,
					lineMain:  3,
					lineLast:  4,
				},
			},
			want: []byte("Line2\nLine3\nLine4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readLines(tt.args.r, tt.args.region); string(got) != string(tt.want) {
				t.Errorf("readLines() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
