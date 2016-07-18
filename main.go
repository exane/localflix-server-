package main

import (
	"net/http"
	"os"
	"time"

	"github.com/exane/localflix-server-/config"
	"github.com/exane/localflix-server-/database"
	"github.com/exane/localflix-server-/fetch"
	"github.com/exane/localflix-server-/loader"
)

func fileserver() {
	cfg := config.LoadConfig()
	http.ListenAndServe(
		cfg.Fileserver.URL+":"+cfg.Fileserver.Port,
		http.FileServer(http.Dir(cfg.Fileserver.RootDirectory)),
	)
}

func server() {
	router()
}

var (
	INSTALL = os.Getenv("INSTALL")
)

func main() {
	os.Setenv("ENV", "development")
	database.InitDb()

	go func() {
		fetch.Fetch()
		series := database.DumpImport()
		if INSTALL == "true" {
			database.CreateTables()
			loader.Import(series)
		} else {
			loader.Update(series)
		}
	}()

	defer database.DB.Close()
	go fileserver()
	go server()
	for {
		time.Sleep(10 * time.Second)
	}
}
