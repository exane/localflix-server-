package main

import "time"

const (
	LIMIT_REQUEST = 40
	LIMIT_RESET   = 15 //seconds
)

var rlc *RequestLimitCheck

type RequestLimitCheck struct {
	started  time.Time
	requests int
}

func (this *RequestLimitCheck) time() time.Duration {
	return time.Duration(LIMIT_RESET)*time.Second - time.Since(this.started)
}

func (this *RequestLimitCheck) reset() {
	this.requests = 0
	this.started = time.Now()
}

func (this *RequestLimitCheck) wait() {
	time.Sleep(this.time())
}

func (this *RequestLimitCheck) checkRequest() {
	if this.requests >= LIMIT_REQUEST {
		println("TMDb Request Limit Wait: ", this.time().String())
		this.wait()
		println("TMDb Request Limit Continue")
		this.reset()
	}
	this.requests++
}
