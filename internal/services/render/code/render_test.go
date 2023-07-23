package code

import (
	"bytes"
	"io"
	"testing"

	"github.com/fe3dback/go-arch-lint/internal/models/common"
)

func Test_readLines(t *testing.T) {
	type args struct {
		r   io.Reader
		ref common.Reference
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "simple",
			args: args{
				r:   bytes.NewReader([]byte("Line1\nLine2\nLine3\nLine4\nLine5\nLine6")),
				ref: common.NewReferenceRange("/", 2, 3, 4),
			},
			want: []byte("Line2\nLine3\nLine4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readLines(tt.args.r, tt.args.ref); string(got) != string(tt.want) {
				t.Errorf("readLines() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
