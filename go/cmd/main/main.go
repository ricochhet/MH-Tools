package main

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/manager"
	"github.com/ricochhet/mhwarchivemanager/pkg/win32"
)

func main() {
	if len(os.Args) > 1 {
		win32.Console(true)
		manager.A_InitializeCommandLine()
	} else {
		win32.Console(false)
		manager.A_InitializeUI()
	}
}
