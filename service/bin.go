package service

import (
        "os"
)

func PluggoBin() string {
	exe, err := os.Executable()
	if err != nil {
		return "kakoune-pluggo"
	}
	return exe
}
