package fragment

import (
	"testing"
)

func Test_line_fragments(t *testing.T) {
	for _, testCase := range []struct {
		desc string
		in   string
		out  Selection
		err  error
	}{
		{
			desc: "single line",
			in:   "line=42",
			out: LineFragment{
				Start: LinePosition{Line: 42},
				End:   LinePosition{Line: 42},
			},
		},
		{
			desc: "start and end lines",
			in:   "line=12,16",
			out: LineFragment{
				Start: LinePosition{Line: 12},
				End:   LinePosition{Line: 16},
			},
		},
		{
			desc: "lines with columns",
			in:   "line=17.6,19.3",
			out: LineFragment{
				Start: LinePosition{Line: 17, Column: 6},
				End:   LinePosition{Line: 19, Column: 3},
			},
		},
		{
			desc: "parse error",
			in:   "3ka/3:--",
			err:  CannotParse,
		},
		{
			desc: "parse char offset",
			in:   "char=167",
			out: CharFragment{
				Start: 167,
				End:   167,
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			out, err := Parse(testCase.in)
			if err != testCase.err {
				if err == nil {
					t.Errorf("want %v, got no error", testCase.err)
				} else {
					t.Errorf("want no error, got %v", err)
				}
				return
			}
			if out != testCase.out {
				t.Errorf("want %v, got %v", testCase.out, out)
			}
		})
	}
}
