package main

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
	"github.com/ricochhet/mhwarchivemanager/pkg/win32"
)

func main() {
	if len(os.Args) > 1 {
		win32.AttachConsoleW()
		logger.Stdout = os.Stdout
		manager.A_InitializeCommandLine()
	} else {
		_, stdout, _ := win32.AllocConsole()
		logger.Stdout = stdout
		manager.A_InitializeUI()
	}
}
