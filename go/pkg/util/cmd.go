package util

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ricochhet/mhwarchivemanager/pkg/logger"
)

const BIN_SIZE int = 2

func Cmd(args []string, name string, argCount int) ([]string, error) {
	var nameArgs []string
	if slices.Contains(args, name) {
		ind := slices.Index(args, name)

		if ind != -1 {
			if len(args) >= argCount+BIN_SIZE {
				nameArgs = append(nameArgs, args[ind:]...)
				return nameArgs, nil
			}

			logger.SharedLogger.Error("Expected " + strconv.Itoa(argCount) + " arguments, but got " + strconv.Itoa(len(args)-BIN_SIZE))
			return nil, fmt.Errorf("length of args is less than argCount")
		}

		return nil, fmt.Errorf("%s not found", name)
	}

	return nil, fmt.Errorf("unhandled error")
}
