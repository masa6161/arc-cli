// Package git provides git operations including diff generation, branch management, and remote handling.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetRoot returns the root directory of the current git repository.
func GetRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not inside a git repository: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
