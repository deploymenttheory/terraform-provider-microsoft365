package powershell

import (
	"bytes"
	"fmt"
	"os/exec"
)

// RunPowerShell executes a PowerShell command and returns its output or error
func RunPowerShell(cmd string) (string, error) {
	c := exec.Command("pwsh", "-Command", cmd)
	var out, stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.String(), nil
}
