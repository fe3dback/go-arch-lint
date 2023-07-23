package decoder

import (
	"testing"
)

func Test_jsonSchemeTransformJsonPathToYamlPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "s1",
			args: args{
				path: "(root)",
			},
			want: "$",
		},
		{
			name: "s2",
			args: args{
				path: "(root).simple",
			},
			want: "$.simple",
		},
		{
			name: "s3",
			args: args{
				path: "(root).tree.path",
			},
			want: "$.tree.path",
		},
		{
			name: "with index",
			args: args{
				path: "(root).tree.path.5",
			},
			want: "$.tree.path[5]",
		},
		{
			name: "with index 2",
			args: args{
				path: "(root).tree.path.100",
			},
			want: "$.tree.path[100]",
		},
		{
			name: "with index tree",
			args: args{
				path: "(root).tree.path.100.anotherList.5.anotherItem",
			},
			want: "$.tree.path[100].anotherList[5].anotherItem",
		},
		{
			name: "inv 1",
			args: args{
				path: "(root).123-hello.b",
			},
			want: "$.123-hello.b",
		},
		{
			name: "inv 1",
			args: args{
				path: "(root).3rd-hello.13.13b",
			},
			want: "$.3rd-hello[13].13b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonSchemeTransformJSONPathToYamlPath(tt.args.path); got != tt.want {
				t.Errorf("jsonSchemeTransformJsonPathToYamlPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
