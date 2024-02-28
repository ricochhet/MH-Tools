package main

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/core"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
	"github.com/ricochhet/mhwarchivemanager/pkg/win32"
)

func main() {
	if len(os.Args) > 1 {
		win32.AttachConsoleW()
		logger.Stdout = os.Stdout
		core.A_InitializeCommandLine()
	} else {
		// Alloc console for development builds
		// _, stdout, _ := win32.AllocConsole()
		logger.Stdout = os.Stdout
		core.A_InitializeUI()
	}
}
