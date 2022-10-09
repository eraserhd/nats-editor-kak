package kakoune

import (
        "fmt"
)

type Command struct {
	Session string
	Script  fmt.Stringer
}
