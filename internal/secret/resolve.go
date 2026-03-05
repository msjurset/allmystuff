package secret

import (
	"fmt"
	"os/exec"
	"strings"
)

// Resolve returns the given string as-is unless it starts with "op://",
// in which case it invokes the 1Password CLI to read the secret reference.
func Resolve(s string) (string, error) {
	if !strings.HasPrefix(s, "op://") {
		return s, nil
	}
	out, err := exec.Command("op", "read", s, "--no-newline").Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("resolving secret: op read failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("resolving secret: %w", err)
	}
	return string(out), nil
}
