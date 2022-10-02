package fragment

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Offset = int

type Selection interface {
	isSelection()
	Fragment() string
}

type LinePosition struct {
	Line, Column Offset
}

func (lp LinePosition) fragmentString() string {
	switch true {
	case lp.Column == 0:
		return fmt.Sprintf("%d", lp.Line)
	default:
		return fmt.Sprintf("%d.%d", lp.Line, lp.Column)
	}
}

type LineAndColumnSelection struct {
	Start, End LinePosition
}

func (_ LineAndColumnSelection) isSelection() {}

func (lc LineAndColumnSelection) Fragment() string {
	switch true {
	case lc.Start == lc.End:
		return fmt.Sprintf("line=%s", lc.Start.fragmentString())
	default:
		return fmt.Sprintf("line=%s,%s", lc.Start.fragmentString(), lc.End.fragmentString())
	}
}

type CharSelection struct {
	Start, End Offset
}

func (_ CharSelection) isSelection() {}

func (cs CharSelection) Fragment() string {
	return ""
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

func Parse(fragment string) (Selection, error) {
	if parts, err := matchAndParseInts(charPattern, fragment); err == nil {
		sel := CharSelection{Start: *parts[0], End: *parts[0]}
		if parts[1] != nil {
			sel.End = *parts[1]
		}
		return sel, nil
	}

	if parts, err := matchAndParseInts(linePattern, fragment); err == nil {
		sel := LineAndColumnSelection{
			Start: LinePosition{Line: *parts[0]},
			End:   LinePosition{Line: *parts[0]},
		}
		if parts[1] != nil {
			sel.Start.Column = *parts[1]
			sel.End.Column = *parts[1]
		}
		if parts[2] != nil {
			sel.End = LinePosition{Line: *parts[2]}
			if parts[3] != nil {
				sel.End.Column = *parts[3]
			}
		}
		return sel, nil
	}

	return nil, CannotParse
}
