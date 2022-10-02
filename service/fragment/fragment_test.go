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
				Start: LineAndColumn{Line: 42},
				End:   LineAndColumn{Line: 42},
			},
		},
		{
			desc: "single line and column",
			in:   "line=42.7",
			out: LineAndColumnSelection{
				Start: LineAndColumn{Line: 42, Column: 7},
				End:   LineAndColumn{Line: 42, Column: 7},
			},
		},
		{
			desc: "start and end lines",
			in:   "line=12,16",
			out: LineAndColumnSelection{
				Start: LineAndColumn{Line: 12},
				End:   LineAndColumn{Line: 16},
			},
		},
		{
			desc: "lines with columns",
			in:   "line=17.6,19.3",
			out: LineAndColumnSelection{
				Start: LineAndColumn{Line: 17, Column: 6},
				End:   LineAndColumn{Line: 19, Column: 3},
			},
		},
		{
			desc: "line range with only start column",
			in:   "line=77.4,78",
			out: LineAndColumnSelection{
				Start: LineAndColumn{Line: 77, Column: 4},
				End:   LineAndColumn{Line: 78, Column: 0},
			},
		},
		{
			desc: "parse error",
			in:   "3ka/3:--",
			err:  CannotParse,
		},
		{
			desc: "integrity checks are ignored",
			in:   "line=13;md5=d41d8cd98f00b204e9800998ecf8427e;length=0",
			out: LineAndColumnSelection{
				Start: LineAndColumn{Line: 13},
				End:   LineAndColumn{Line: 13},
			},
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
		{
			desc: "char integrity checks are ignored",
			in:   "char=13;md5=d41d8cd98f00b204e9800998ecf8427e;length=0",
			out: CharSelection{
				Start: 13,
				End:   13,
			},
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			out, err := ParseRFC5147FragmentIdentifier(testCase.in)
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
				Start: LineAndColumn{Line: 42},
				End:   LineAndColumn{Line: 42},
			},
			out: "line=42",
		},
		{
			desc: "zero-width, line and column selection",
			in: LineAndColumnSelection{
				Start: LineAndColumn{Line: 12, Column: 7},
				End:   LineAndColumn{Line: 12, Column: 7},
			},
			out: "line=12.7",
		},
		{
			desc: "nonzero-width multi-line selection",
			in: LineAndColumnSelection{
				Start: LineAndColumn{Line: 16},
				End:   LineAndColumn{Line: 22},
			},
			out: "line=16,22",
		},
		{
			desc: "zero-width char offset",
			in:   CharSelection{Start: 167, End: 167},
			out:  "char=167",
		},
		{
			desc: "char range",
			in:   CharSelection{Start: 49, End: 77},
			out:  "char=49,77",
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			out := testCase.in.RFC5147FragmentIdentifier()
			if out != testCase.out {
				t.Errorf("want out = %q, got %q", testCase.out, out)
			}
			roundTrip, err := ParseRFC5147FragmentIdentifier(out)
			if err != nil {
				t.Errorf("round trip got error %v", err)
			}
			if roundTrip != testCase.in {
				t.Errorf("round trip got %v, but wanted %v", roundTrip, testCase.in)
			}
		})
	}
}
