package loop

import "syscall"

func configureTimer() func() {
	winmm := syscall.NewLazyDLL("winmm.dll")
	winmm.NewProc("timeBeginPeriod").Call(uintptr(1))
	return func() { winmm.NewProc("timeEndPeriod").Call(uintptr(1)) }
}
