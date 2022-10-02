package fragment

import (
	"testing"
)

func Test_line_fragments(t *testing.T) {
	for _, testCase := range []struct {
		desc string
		in   string
		out  LineFragment
		err  error
	}{
		{
			desc: "single line",
			in:   "line=42",
			out: LineFragment{
				Start: LinePosition{
					Line: 42,
				},
				End: LinePosition{
					Line: 42,
				},
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			out, err := Parse(testCase.in)
			if err != nil {
				t.Error(err)
				return
			}
			if out != testCase.out {
				t.Errorf("want %v, got %v", testCase.out, out)
			}
		})
	}
}
