package util

import (
	"os"
	"os/exec"
)

// NewExecutableDockerCmd new docker cmd
func NewExecutableDockerCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
