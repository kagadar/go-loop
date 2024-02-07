package loop

import "syscall"

// https://learn.microsoft.com/en-us/windows/win32/api/timeapi/nf-timeapi-timebeginperiod
func configureTimer() func() {
	winmm := syscall.NewLazyDLL("winmm.dll")
	winmm.NewProc("timeBeginPeriod").Call(uintptr(1))
	return func() { winmm.NewProc("timeEndPeriod").Call(uintptr(1)) }
}
