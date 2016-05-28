package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
  "encoding/json"
)

func toJSON(v interface{}) []byte {
  js, _ := json.Marshal(v)
  return js
}

func router() {
  router := mux.NewRouter().StrictSlash(true)

  router.HandleFunc("/", index)
  router.HandleFunc("/video/{id}", video)
  router.HandleFunc("/test", test)

  router.HandleFunc("/tmdb/search/{name}", findSeries).Methods("GET")

  router.HandleFunc("/series", series).Methods("GET")
  router.HandleFunc("/serie/{serie_id}", serie).Methods("GET")
  router.HandleFunc("/season/{season_id}", season).Methods("GET")
  router.HandleFunc("/episodes/{season_id}", episodes).Methods("GET")
  router.HandleFunc("/episode/{episode_id}", episode).Methods("GET")

  http.ListenAndServe(config.Server.Url + ":" + config.Server.Port, router)
}

func findSeries(w http.ResponseWriter, r *http.Request) {
  title := mux.Vars(r)["name"]
  res := findSerie(title)
  w.Header().Set("Content-Type", "text/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(toJSON(res.Results)))
}

func episode(w http.ResponseWriter, r *http.Request) {
  episode_id := mux.Vars(r)["episode_id"]
  episode := Episode{}

  DB.Find(&episode, episode_id)

  season := Season{}
  DB.Find(&season, episode.SeasonID)
  season_name := season.Name

  serie := Serie{}
  DB.Find(&serie, season.SerieID)
  serie_name := serie.Name

  seasonTitle := season.OriginalName
  if !validTitle(seasonTitle) {
    seasonTitle = season.Name
  }

  serieTitle := serie.OriginalName
  if !validTitle(serieTitle) {
    serieTitle = serie.Name
  }

  result := struct {
    Episode
    SeasonName         string
    SeasonOriginalName string
    SerieName          string
    SerieOriginalName  string
    SerieID            int
  }{
    episode,
    season_name,
    seasonTitle,
    serie_name,
    serieTitle,
    season.SerieID,
  }

  w.Header().Set("Content-Type", "text/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(toJSON(result)))
}

func episodes(w http.ResponseWriter, r *http.Request) {
  season_id := mux.Vars(r)["season_id"]
  type episode struct {
    Episode
    SeasonNumber      int
    SeasonName        string
    SeasonDescription string
  }
  episodes := []episode{}

  //DB.Where("season_id = ?", season_id).Find(&episodes)
  DB.
  Select("episodes.*, s.season_number, s.original_name season_name, s.description season_description").
  //Table("episodes").
  Joins("left join seasons s on episodes.season_id = ?", season_id).
  Where("s.id = ?", season_id).
  Find(&episodes)

  w.Header().Set("Content-Type", "text/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(toJSON(episodes)))
}

func season(w http.ResponseWriter, r *http.Request) {
  season_id := mux.Vars(r)["season_id"]

  season := Season{}

  DB.Find(&season, season_id).
  Related(&season.Episodes)

  serie := Serie{}
  DB.Find(&serie, season.SerieID)

  title := serie.OriginalName
  if !validTitle(title) {
    title = serie.Name
  }

  result := struct {
    Season
    SerieName string
  }{
    season,
    title,
  }

  w.Header().Set("Content-Type", "text/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(toJSON(result)))
}

func serie(w http.ResponseWriter, r *http.Request) {
  serie_id := mux.Vars(r)["serie_id"]

  serie := Serie{}

  DB.Find(&serie, serie_id)
  DB.Model(serie).Related(&serie.Seasons)

  w.Header().Set("Content-Type", "text/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(toJSON(serie)))
}

func series(w http.ResponseWriter, r *http.Request) {
  /*series := []Serie{}
  DB.Find(&series)

  for i, _ := range series {
    DB.Model(series[i]).Related(&series[i].Seasons)
  }*/
  //spew.Dump(series)
  w.Header().Set("Content-Type", "text/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Write([]byte(toJSON(allSeries())))
}

func index(rw http.ResponseWriter, req *http.Request) {
  fmt.Println("yo index")
}

func video(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "text/json")
  /*json, _ := json.Marshal(struct {
    Name string
    Url string
  }{result.Name, "http://localhost:3001/" + result.Name + "/" + "S4" + "/" + result.Nr + "." + result.Ext})*/
  //ep := DB.First(&Episode{}, 1)
  episode := &Episode{}
  season := &Season{}
  subtitle := &Subtitle{}

  DB.Find(episode, mux.Vars(req)["id"])
  DB.Find(season)
  DB.Find(subtitle)
  DB.Model(season).Related(episode)
  DB.Model(subtitle).Related(episode)
  //video.Related(&season)
  json, _ := json.Marshal(episode)
  w.Write([]byte(json))
}

func test(w http.ResponseWriter, r *http.Request) {
  serie := &Serie{}
  a := serie.seasons()
  fmt.Println(a)
}