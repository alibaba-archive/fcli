// +build !windows

package util

import (
	"os/exec"
)

// NewExecutableSubCmd in unix
func NewExecutableSubCmd(cmdString string) *exec.Cmd {
	return GetBindingCmd("sh", "-c", cmdString)
}
