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
			out: LineAndColumnSelection{
				Start: LinePosition{Line: 42},
				End:   LinePosition{Line: 42},
			},
		},
		{
			desc: "single line and column",
			in:   "line=42.7",
			out: LineAndColumnSelection{
				Start: LinePosition{Line: 42, Column: 7},
				End:   LinePosition{Line: 42, Column: 7},
			},
		},
		{
			desc: "start and end lines",
			in:   "line=12,16",
			out: LineAndColumnSelection{
				Start: LinePosition{Line: 12},
				End:   LinePosition{Line: 16},
			},
		},
		{
			desc: "lines with columns",
			in:   "line=17.6,19.3",
			out: LineAndColumnSelection{
				Start: LinePosition{Line: 17, Column: 6},
				End:   LinePosition{Line: 19, Column: 3},
			},
		},
		{
			desc: "line range with only start column",
			in:   "line=77.4,78",
			out: LineAndColumnSelection{
				Start: LinePosition{Line: 77, Column: 4},
				End:   LinePosition{Line: 78, Column: 0},
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
			out: CharSelection{
				Start: 167,
				End:   167,
			},
		},
		{
			desc: "parse char range",
			in:   "char=96,107",
			out: CharSelection{
				Start: 96,
				End:   107,
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

func Test_Selection_to_fragment(t *testing.T) {
	for _, testCase := range []struct {
		desc string
		in   Selection
		out  string
	}{
		{
			desc: "zero-width, single line selection",
			in: LineAndColumnSelection{
				Start: LinePosition{Line: 42},
				End:   LinePosition{Line: 42},
			},
			out: "line=42",
		},
		{
			desc: "zero-width, line and column selection",
			in: LineAndColumnSelection{
				Start: LinePosition{Line: 12, Column: 7},
				End:   LinePosition{Line: 12, Column: 7},
			},
			out: "line=12.7",
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			out := testCase.in.Fragment()
			if out != testCase.out {
				t.Errorf("want out = %q, got %q", testCase.out, out)
			}
		})
	}
}
