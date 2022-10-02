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

type LineFragment struct {
	Start, End LinePosition
}

func (_ LineFragment) isSelection() {}

type CharFragment struct {
	Start, End Offset
}

func (_ CharFragment) isSelection() {}

var (
	charPattern = regexp.MustCompile(`^char=(\d+)$`)
	linePattern = regexp.MustCompile(`^line=(\d+)(?:\.(\d+))?(?:,(\d+)(?:\.(\d+))?)?$`)

	CannotParse = errors.New("cannot parse fragment identifier")
)

func Parse(fragment string) (Selection, error) {
	if match := charPattern.FindStringSubmatch(fragment); match != nil {
		offset, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %w", match[1], err)
		}
		return CharFragment{
			Start: int(offset),
			End:   int(offset),
		}, nil
	}

	var result LineFragment
	match := linePattern.FindStringSubmatch(fragment)
	if match == nil {
		return nil, CannotParse
	}
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
