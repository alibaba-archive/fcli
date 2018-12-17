package util

import (
	"os/exec"
)

// NewExecutableSubCmd in windows, powershell installed by default in windows
func NewExecutableSubCmd(cmdString string) *exec.Cmd {
	return GetBindingCmd("powershell", "-Command", cmdString)
}
