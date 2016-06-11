package main

import "time"

const (
  LIMIT_REQUEST = 40
  LIMIT_RESET = 15 //seconds
)

type RequestLimitCheck struct {
  Started  time.Time
  Requests int
}

func (this *RequestLimitCheck) Time() time.Duration {
  return time.Duration(LIMIT_RESET) * time.Second - time.Since(this.Started)
}

func (this *RequestLimitCheck) Reset() {
  this.Requests = 0
  this.Started = time.Now()
}

func (this *RequestLimitCheck) Wait() {
  time.Sleep(this.Time())
}

func (this *RequestLimitCheck) CheckRequest() {
  if this.Requests >= LIMIT_REQUEST {
    println("TMDb Request Limit Wait: ", this.Time().String())
    this.Wait()
    println("TMDb Request Limit Continue")
    this.Reset()
  }
  this.Requests++
}