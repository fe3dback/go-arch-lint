package warnparser

import (
	"reflect"
	"testing"
)

func Test_parseSourceError(t *testing.T) {
	type args struct {
		sourceText string
	}
	tests := []struct {
		name string
		args args
		want sourceMarker
	}{
		{
			name: "valid 183",
			args: args{
				sourceText: `
  181 |   game_component:
  182 |     mayDependOn:
> 183 |       - engine
                ^
  184 |       - engine_entity
  185 |       - game_component
  186 |       - game_utils
  187 |     
`,
			},
			want: sourceMarker{
				valid:      true,
				sourceLine: 183,
				sourcePos:  9,
			},
		},
		{
			name: "valid 1",
			args: args{
				sourceText: `>  1 | version: 2
                ^
   2 | allow:
   3 |   depOnAnyVendor: false
   4 |`,
			},
			want: sourceMarker{
				valid:      true,
				sourceLine: 1,
				sourcePos:  10,
			},
		},
		{
			name: "invalid",
			args: args{
				sourceText: "",
			},
			want: sourceMarker{
				valid:      false,
				sourceLine: 0,
				sourcePos:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSourceWarning(tt.args.sourceText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSourceError() = %v, want %v", got, tt.want)
			}
		})
	}
}
