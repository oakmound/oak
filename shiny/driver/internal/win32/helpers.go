package win32

func Minimize(hwnd HWND) bool {
	return ShowWindow(hwnd, _SW_MINIMIZE)
}

func Maximize(hwnd HWND) bool {
	return ShowWindow(hwnd, _SW_SHOWMAXIMIZED)
}

func Normalize(hwnd HWND) bool {
	return ShowWindow(hwnd, _SW_SHOWNORMAL)
}
