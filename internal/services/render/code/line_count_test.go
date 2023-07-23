package code

import (
	"bytes"
	"io"
	"testing"
)

func Test_lineCounter(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				r: bytes.NewReader([]byte("Hello world\nThis buffer has three lines\n")),
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lineCounter(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("lineCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lineCounter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
