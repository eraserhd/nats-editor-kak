package fragment

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Offset is the basic component for coordinates.
type Offset = int

// Selection is a possibly zero-width character range in a text document.
type Selection interface {
	isSelection()

	// RFC5147FragmentIdentifier returns the fragment identifier (without the #) that identifies this selection.
	RFC5147FragmentIdentifier() string
}

type LineAndColumn struct {
	Line, Column Offset
}

func (lp LineAndColumn) fragmentString() string {
	switch true {
	case lp.Column == 0:
		return fmt.Sprintf("%d", lp.Line)
	default:
		return fmt.Sprintf("%d.%d", lp.Line, lp.Column)
	}
}

// LineAndColumnSelection represents a selection whose ends are specified in lines and columns.
type LineAndColumnSelection struct {
	Start, End LineAndColumn
}

func (_ LineAndColumnSelection) isSelection() {}

func (lc LineAndColumnSelection) RFC5147FragmentIdentifier() string {
	switch true {
	case lc.Start == lc.End:
		return fmt.Sprintf("line=%s", lc.Start.fragmentString())
	default:
		return fmt.Sprintf("line=%s,%s", lc.Start.fragmentString(), lc.End.fragmentString())
	}
}

// CharSelection represents a selection whose ends are specified in codepoint offsets.
type CharSelection struct {
	Start, End Offset
}

func (_ CharSelection) isSelection() {}

func (cs CharSelection) RFC5147FragmentIdentifier() string {
	switch true {
	case cs.Start == cs.End:
		return fmt.Sprintf("char=%d", cs.Start)
	default:
		return fmt.Sprintf("char=%d,%d", cs.Start, cs.End)
	}
}

var (
	charPattern = regexp.MustCompile(`^char=(\d+)(?:,(\d+))?$`)
	linePattern = regexp.MustCompile(`^line=(\d+)(?:\.(\d+))?(?:,(\d+)(?:\.(\d+))?)?$`)

	CannotParse = errors.New("cannot parse fragment identifier")
	noMatch     = errors.New("did not match")
)

func matchAndParseInts(pattern *regexp.Regexp, s string) ([]*Offset, error) {
	match := pattern.FindStringSubmatch(s)
	if match == nil {
		return nil, noMatch
	}
	result := make([]*Offset, len(match)-1)
	for i, s := range match[1:] {
		if s == "" {
			continue
		}
		var err error
		val, err := strconv.ParseInt(s, 10, strconv.IntSize)
		if err != nil {
			return nil, err
		}
		slot := int(val)
		result[i] = &slot
	}
	return result, nil
}

// ParseRFC5147FragmentIdentifier parses fragment into a Selection or returns an error.
// Selection will have concrete type LineAndColumnSelection or CharSelection, depending on
// what was parsed.
func ParseRFC5147FragmentIdentifier(fragment string) (Selection, error) {
	if parts, err := matchAndParseInts(charPattern, fragment); err == nil {
		sel := CharSelection{Start: *parts[0], End: *parts[0]}
		if parts[1] != nil {
			sel.End = *parts[1]
		}
		return sel, nil
	}

	if parts, err := matchAndParseInts(linePattern, fragment); err == nil {
		sel := LineAndColumnSelection{
			Start: LineAndColumn{Line: *parts[0]},
			End:   LineAndColumn{Line: *parts[0]},
		}
		if parts[1] != nil {
			sel.Start.Column = *parts[1]
			sel.End.Column = *parts[1]
		}
		if parts[2] != nil {
			sel.End = LineAndColumn{Line: *parts[2]}
			if parts[3] != nil {
				sel.End.Column = *parts[3]
			}
		}
		return sel, nil
	}

	return nil, CannotParse
}
