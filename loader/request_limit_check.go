package loader

import "time"

const (
	LIMIT_REQUEST = 40
	LIMIT_RESET   = 15 //seconds
)

var (
	started   time.Time
	requests  int
	IsTesting bool = false
)

func Time() time.Duration {
	return time.Duration(LIMIT_RESET)*time.Second - time.Since(started)
}

func Reset() {
	requests = 0
	started = time.Now()
}

func Wait() {
	if IsTesting {
		return
	}
	time.Sleep(Time())
}

func LimitReached() bool {
	return requests >= LIMIT_REQUEST
}

func CheckRequest() {
	if LimitReached() {
		if !IsTesting {
			println("TMDb Request Limit Wait: ", Time().String())
		}
		Wait()
		if !IsTesting {
			println("TMDb Request Limit Continue")
		}
		Reset()
	}
	requests++
}

func Requests() int {
	return requests
}

func init() {
	Reset()
}
