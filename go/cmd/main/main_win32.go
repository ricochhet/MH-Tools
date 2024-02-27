package main

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
	"github.com/ricochhet/mhwarchivemanager/pkg/win32"
)

func main() {
	if len(os.Args) > 1 {
		win32.AttachConsoleW()
		manager.A_InitializeCommandLine()
	} else {
		manager.A_InitializeUI()
	}
}
