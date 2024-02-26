package win32

import "syscall"

func Console(show bool) {
	GetConsoleWindow := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	ShowWindow := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	hwnd, _, _ := GetConsoleWindow.Call()
	if hwnd == 0 {
		return
	}
	if show {
		var SW_RESTORE uintptr = 9
		ShowWindow.Call(hwnd, SW_RESTORE)
	} else {
		var SW_HIDE uintptr = 0
		ShowWindow.Call(hwnd, SW_HIDE)
	}
}
