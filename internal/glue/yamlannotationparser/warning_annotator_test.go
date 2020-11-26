package yamlannotationparser

import (
	"reflect"
	"testing"
)

func Test_parseSourceError(t *testing.T) {
	type args struct {
		sourceText string
	}
	tests := []struct {
		name      string
		args      args
		want      sourceMarker
		wantError bool
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
				sourceLine: 183,
				sourcePos:  9,
			},
			wantError: false,
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
				sourceLine: 1,
				sourcePos:  10,
			},
			wantError: false,
		},
		{
			name: "valid, no line marker",
			args: args{
				sourceText: `32 |        - 3rd-cobra
   33 |   cmd:    canUse:
   			      ^
   36 |       - go-modfile
   37 |
   38 |   a:`,
			},
			want: sourceMarker{
				sourceLine: 33,
				sourcePos:  5,
			},
			wantError: false,
		},
		{
			name: "invalid",
			args: args{
				sourceText: "",
			},
			want: sourceMarker{
				sourceLine: 0,
				sourcePos:  0,
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.sourceText)

			if err == nil && tt.wantError {
				t.Errorf("parse() = expected error, but is not, got = %v", got)
				return
			}

			if err != nil && !tt.wantError {
				t.Errorf("parse() = unexpected err: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
