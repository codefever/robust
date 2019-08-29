// +build linux

package subprocess

import (
	"os/exec"
	"syscall"
)

// maybeSetUpPdeathsig@linux kills the child process when its parent exits.
func maybeSetUpPdeathsig(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}
}
