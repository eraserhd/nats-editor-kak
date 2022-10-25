package service

import (
        "os"
        "path"
)

func BinPath() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return path.Dir(exe)
}
