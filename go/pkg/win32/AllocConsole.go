package win32

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/ricochhet/mhwarchivemanager/pkg/thirdparty/ansi"
)

func AllocConsole() (aIn, aOut, aErr io.Writer) {
	kernal23 := syscall.NewLazyDLL("kernel32.dll")
	allocConsole := kernal23.NewProc("AllocConsole")
	r0, _, err0 := syscall.SyscallN(allocConsole.Addr(), 0, 0, 0, 0)
	if r0 == 0 {
		fmt.Printf("Could not allocate console: %s. Check build flags.", err0)
		os.Exit(1)
	}

	hIn, err1 := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	hOut, err2 := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	hErr, err3 := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)

	if err1 != nil || err2 != nil || err3 != nil {
		os.Exit(2)
	}

	stdinFile := os.NewFile(uintptr(hIn), "/dev/stdin")
	stdoutFile := os.NewFile(uintptr(hOut), "/dev/stdout")
	stderrFile := os.NewFile(uintptr(hErr), "/dev/stderr")

	aStdin := ansi.NewAnsiStdoutW(stdinFile)
	aStdout := ansi.NewAnsiStdoutW(stdoutFile)
	aStderr := ansi.NewAnsiStdoutW(stderrFile)

	os.Stdin = stdinFile
	os.Stdout = stdoutFile
	os.Stderr = stderrFile

	return aStdin, aStdout, aStderr
}
