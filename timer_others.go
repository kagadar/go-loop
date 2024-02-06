//go:build !windows

package loop

func configureTimer() func() {
	return func() {}
}
