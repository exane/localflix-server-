package main

import (
  _ "time"
  "net/http"
  "time"
  "os"
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

var(
  INSTALL = os.Getenv("INSTALL")
)

func main() {

  init_db()

  if INSTALL == "true" {
    go func() {
      create_tables()
      dump_import()
      load_tmdb()
    }()
  }


  defer DB.Close()
  go fileserver()
  go server()
  for {
    time.Sleep(10 * time.Second)
  }

  return
}