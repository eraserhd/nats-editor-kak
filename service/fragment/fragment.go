package fragment

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type TextPlainFragmentIdentifier struct {
	StartLine int
	EndLine   int
}

var fragmentRegexp = regexp.MustCompile(`^line=(\d+)(?:,(\d+))?`)

func Parse(fragment string) (TextPlainFragmentIdentifier, int, error) {
	var result TextPlainFragmentIdentifier
	match := fragmentRegexp.FindStringSubmatch(fragment)
	if match == nil {
		return result, 0, errors.New("cannot parse fragment identifier")
	}
	line, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return result, 0, fmt.Errorf("parsing fragment identifer: %w", err)
	}
	result.StartLine = int(line)
	result.EndLine = int(line)
	if match[2] != "" {
		endLine, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return result, 0, fmt.Errorf("parsing fragment identifier: %w", err)
		}
		result.EndLine = int(endLine)
	}
	return result, int(result.EndLine), nil
}
