package loader

import "time"

//X-RateLimit-Limit: 40
//X-RateLimit-Remaining: 19
//X-RateLimit-Reset: 1468249347

const (
	LIMIT_REQUEST = 40
	LIMIT_RESET   = 20 //seconds
)

var (
	started   time.Time
	requests  int
	IsTesting bool = false
	Requested map[string]int
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

func CheckRequest(name string) {
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

	Requested[name]++
	requests++
}

func Requests() int {
	return requests
}

func init() {
	Requested = make(map[string]int)
	Reset()
}
