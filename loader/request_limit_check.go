package loader

import "time"

const (
	LIMIT_REQUEST = 40
	LIMIT_RESET   = 15 //seconds
)

var (
	started  time.Time
	requests int
)

func Time() time.Duration {
	return time.Duration(LIMIT_RESET)*time.Second - time.Since(started)
}

func Reset() {
	requests = 0
	started = time.Now()
}

func Wait() {
	time.Sleep(Time())
}

func CheckRequest() {
	if requests >= LIMIT_REQUEST {
		println("TMDb Request Limit Wait: ", Time().String())
		Wait()
		println("TMDb Request Limit Continue")
		Reset()
	}
	requests++
}

func init() {
	Reset()
}
