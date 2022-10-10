package kakoune

import (
	"fmt"
	"os"
)

func Quote(s string) string {
	result := "'"
	for _, ch := range s {
		if ch == '\'' {
			result += "'"
		}
		result += string(ch)
	}
	return result + "'"
}

func Fail(msg string) {
	fmt.Printf("fail " + Quote(msg) + "\n")
	os.Exit(1)
}
