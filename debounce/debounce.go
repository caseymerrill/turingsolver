package debounce

import "time"

func Debounce(f func(), delay time.Duration) func() {
	var timer *time.Timer
	return func() {
		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(delay, f)
	}
}