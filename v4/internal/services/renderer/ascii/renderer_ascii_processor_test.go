package ascii

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_preprocessTemplateLine(t *testing.T) {
	tests := []struct {
		row  string
		want string
	}{
		{
			row:  "\t\thello world",
			want: "hello world",
		},
	}
	for ind, tt := range tests {
		t.Run(fmt.Sprintf("testcase-%d", ind), func(t *testing.T) {
			got := preprocessTemplateLine(tt.want)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_preprocessRawASCIITemplate(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{
			text: `
Simple Line
	Tabbed line
	Line  With  Multiple spaces
Another   spaced line
	  
  
 
^ empty lines here, we skip this
	and another line
`,
			want: `Simple Line
Tabbed line
Line  With  Multiple spaces
Another   spaced line
^ empty lines here, we skip this
and another line`,
		},
	}
	for ind, tt := range tests {
		t.Run(fmt.Sprintf("testcase-%d", ind), func(t *testing.T) {
			got := preprocessRawASCIITemplate(tt.text)
			assert.Equal(t, tt.want, got)
		})
	}
}
