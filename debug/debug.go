package main

import (
  "io/ioutil"
  "encoding/json"
  "github.com/davecgh/go-spew/spew"
  "github.com/jinzhu/gorm"
  "time"
)

func main() {
  LIMIT_REQUEST := 40
  LIMIT_RESET := 10 //seconds

  rlc_reset := func(rlc *RequestLimitCheck) {
    rlc.requests = 0
    rlc.started = time.Now()
  }

  rlc_wait := func(rlc *RequestLimitCheck) time.Duration {
    return time.Duration(LIMIT_RESET) * time.Second - time.Since(rlc.started)
  }
  rlc := &RequestLimitCheck{}
  rlc_reset(rlc)
  for i := 0; i < 60; i++ {
    if rlc.requests >= LIMIT_REQUEST {
      println("RLC WAIT: %d", rlc_wait(rlc).String())
      time.Sleep(rlc_wait(rlc))
      rlc_reset(rlc)
    }
    time.Sleep(time.Millisecond * 100)
    println("Do request....", i)
    rlc.requests++
  }

}

type RequestLimitCheck struct {
  started  time.Time
  requests int
}

func dump_import() {
  js, _ := ioutil.ReadFile("Y:/golangWorkspace/src/github.com/exane/localflix/debug/DATA_DUMP.json")
  data := []Serie{}
  json.Unmarshal(js, &data)
  spew.Dump(data)
}

type Serie struct {
  gorm.Model
  Name        string `json:"Name"`
  Description string
  Seasons     []Season
}

type Season struct {
  gorm.Model
  Name        string
  Nr          string
  Description string
  Episodes    []Episode
  SerieID     int
}

type Episode struct {
  gorm.Model
  Nr          string
  Name        string
  Description string
  Src         string
  Ext         string
  SeasonID    int
  Extension   string
  Subtitles   []string
}
