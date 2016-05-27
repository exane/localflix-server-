package main

import (
  _ "time"
  "net/http"
  "time"
  "os"
  "github.com/exane/localflix/fetch"
)

func fileserver() {
  http.ListenAndServe(
    config.Fileserver.Url + ":" + config.Fileserver.Port,
    http.FileServer(http.Dir(config.Fileserver.Root_directory)),
  )
}

func server() {
  router()
}

var (
  INSTALL = os.Getenv("INSTALL")
)

func main() {

  initDb()

  go func() {
    fetch.Fetch()
    if INSTALL == "true" {
      createTables()
      dumpImport()
      loadTmdb()
    } else {
      updateDb()
    }
  }()

  defer DB.Close()
  go fileserver()
  go server()
  for {
    time.Sleep(10 * time.Second)
  }

  return
}