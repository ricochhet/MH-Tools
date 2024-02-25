package process

import (
	"os/exec"
	"runtime"
	"syscall"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func DoesExecutableExist(name string) bool {
	if _, err := exec.LookPath(name); err != nil {
		logger.SharedLogger.Debug(name + " does not exist in PATH")
		return false
	}

	return true
}

func RunExecutable(name string, hideWindow bool, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if runtime.GOOS == "windows" {
		// cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: hideWindow}
	}
	return cmd.Run()
}
