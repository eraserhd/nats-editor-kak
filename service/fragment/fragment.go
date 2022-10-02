package fragment

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Offset = int

type LinePosition struct {
	Line, Column Offset
}

type LineFragment struct {
	Start, End LinePosition
}

var (
	fragmentRegexp = regexp.MustCompile(`^line=(\d+)(?:\.(\d+))?(?:,(\d+)(?:\.(\d+))?)?$`)

	CannotParse = errors.New("cannot parse fragment identifier")
)

func Parse(fragment string) (LineFragment, error) {
	var result LineFragment
	match := fragmentRegexp.FindStringSubmatch(fragment)
	if match == nil {
		return result, CannotParse
	}
	var parts [4]int64
	for i := 0; i < 4; i++ {
		if match[i+1] == "" {
			continue
		}
		var err error
		parts[i], err = strconv.ParseInt(match[i+1], 10, 64)
		if err != nil {
			return result, fmt.Errorf("parsing %q: %w", parts[i], err)
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
