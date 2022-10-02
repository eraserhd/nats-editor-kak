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
}

type LinePosition struct {
	Line, Column Offset
}

type LineAndColumnSelection struct {
	Start, End LinePosition
}

func (_ LineAndColumnSelection) isSelection() {}

type CharSelection struct {
	Start, End Offset
}

func (_ CharSelection) isSelection() {}

var (
	charPattern = regexp.MustCompile(`^char=(\d+)(?:,(\d+))?$`)
	linePattern = regexp.MustCompile(`^line=(\d+)(?:\.(\d+))?(?:,(\d+)(?:\.(\d+))?)?$`)

	CannotParse = errors.New("cannot parse fragment identifier")
	noMatch     = errors.New("did not match")
)

func matchAndExtractOffset(pattern *regexp.Regexp, s string) ([]*Offset, error) {
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
	if parts, err := matchAndExtractOffset(charPattern, fragment); err == nil {
		sel := CharSelection{Start: *parts[0], End: *parts[0]}
		if parts[1] != nil {
			sel.End = *parts[1]
		}
		return sel, nil
	}

	if match := linePattern.FindStringSubmatch(fragment); match != nil {
		var parts [4]int64
		for i := 0; i < 4; i++ {
			if match[i+1] == "" {
				continue
			}
			var err error
			parts[i], err = strconv.ParseInt(match[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %w", parts[i], err)
			}
		}

		var result LineAndColumnSelection
		result.Start.Line = int(parts[0])
		result.Start.Column = int(parts[1])
		if parts[2] != 0 {
			result.End.Line = int(parts[2])
			result.End.Column = int(parts[3])
		} else {
			result.End = result.Start
		}

		return result, nil
	}

	return nil, CannotParse
}
