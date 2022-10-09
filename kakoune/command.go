package kakoune

import (
	"fmt"
	"os/exec"
)

type Command struct {
	Session string
	Script  fmt.Stringer
}

func Run(o Command) error {
	cmd := exec.Command("kak", "-p", o.Session)
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting kak: %w", err)
	}
	if _, err := in.Write([]byte(o.Script.String())); err != nil {
		return fmt.Errorf("error writing script: %w", err)
	}
	if err := in.Close(); err != nil {
		return fmt.Errorf("closing pipe: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error responding: %w", err)
	}
	return nil
}
