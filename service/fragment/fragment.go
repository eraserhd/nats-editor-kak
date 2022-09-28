package fragment

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var fragmentRegexp = regexp.MustCompile(`^line=(\d+)(?:,(\d+))?`)

func Parse(fragment string) (int, int, error) {
	match := fragmentRegexp.FindStringSubmatch(fragment)
	if match == nil {
		return 0, 0, errors.New("cannot parse fragment identifier")
	}
	line, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing fragment identifer: %w", err)
	}
	endLine := line
	if match[2] != "" {
		endLine, err = strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("parsing fragment identifier: %w", err)
		}
	}
	return int(line), int(endLine), nil
}
