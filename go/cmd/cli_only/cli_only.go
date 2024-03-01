package main

import (
	"os"

	"github.com/ricochhet/mhwarchivemanager/pkg/core"
	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

func main() {
	logger.Stdout = os.Stdout
	core.A_InitializeCommandLine()
}
