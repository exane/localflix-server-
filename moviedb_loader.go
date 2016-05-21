package main

import (
  "github.com/ryanbradynd05/go-tmdb"
  "time"
  "regexp"
  "strconv"
)

var tmdn *tmdb.TMDb

const (
  LIMIT_REQUEST = 40
  LIMIT_RESET = 15 //seconds
)

var series_dump []Serie

func getTmdb() *tmdb.TMDb {
  if tmdn == nil {
    tmdn = tmdb.Init(load_config().TMDb.API_KEY)
  }
  return tmdn
}

func load_tmdb() {
  tmdn = tmdb.Init(load_config().TMDb.API_KEY)
  series_dump = loadDump("./fetch/DATA_DUMP.json")

  rlc := &RequestLimitCheck{}
  rlc.reset()

  load_series(rlc)
}

type RequestLimitCheck struct {
  started  time.Time
  requests int
}

func (rlc *RequestLimitCheck) time() time.Duration {
  return time.Duration(LIMIT_RESET) * time.Second - time.Since(rlc.started)
}

func (rlc *RequestLimitCheck) reset() {
  rlc.requests = 0
  rlc.started = time.Now()
}

func (rlc *RequestLimitCheck) wait() {
  time.Sleep(rlc.time())
}

func (rlc *RequestLimitCheck) checkRequest() {
  if rlc.requests >= LIMIT_REQUEST {
    println("TMDb Request Limit Wait: ", rlc.time().String())
    rlc.wait()
    println("TMDb Request Limit Continue")
    rlc.reset()
  }
  rlc.requests++
}

func load_series(rlc *RequestLimitCheck) {
  println("TMDb Series Import Start")

  for _, val := range series_dump {
    /*if rlc.requests >= LIMIT_REQUEST {
      println("TMDb Request Limit Wait: ", rlc.time().String())
      rlc.wait()
      println("TMDb Request Limit Continue")
      rlc.reset()
    }*/
    serie := Serie{}
    DB.Where("Name = ?", val.Name).First(&serie)

    rlc.checkRequest()
    search_tv, err := tmdn.SearchTv(val.Name, nil)

    if err != nil {
      println("@@@@@@@@@@@@@@@")
      println(err.Error())
      println("@@@@@@@@@@@@@@@")
    }

    if len(search_tv.Results) == 0 {
      continue
    }
    rlc.checkRequest()
    tv, _ := tmdn.GetTvInfo(search_tv.Results[0].ID, nil)
    serie.PosterPath = tv.PosterPath
    serie.FirstAirDate = tv.FirstAirDate
    serie.VoteAverage = tv.VoteAverage
    serie.VoteCount = tv.VoteCount
    serie.OriginalName = tv.OriginalName
    serie.Description = tv.Overview

    //fetch season data
    DB.Model(serie).Related(&serie.Seasons)
    load_seasons(rlc, &serie, tv)

    DB.Save(&serie)
  }
  println("TMDb Series Import Done")
}

func fetch_number(name string) int {
  regex := regexp.MustCompile("[Ss]?(\\d+)")
  ret, _ := strconv.Atoi(regex.ReplaceAllString(name, "$1"))
  return ret
}

func fetch_season(season *Season, seasonInfo *tmdb.TvSeason) {
  season.Description = seasonInfo.Overview
  season.PosterPath = seasonInfo.PosterPath
  season.SeasonNumber = seasonInfo.SeasonNumber
  season.AirDate = seasonInfo.AirDate
  season.OriginalName = seasonInfo.Name
}

func load_seasons(rlc *RequestLimitCheck, serie *Serie, tv *tmdb.TV) {
  println("TMDb Seasons Import Start", tv.Name)

  for _, tv_season := range tv.Seasons {
    hasSeason := false

    rlc.checkRequest()
    seasonInfo, err := tmdn.GetTvSeasonInfo(tv.ID, tv_season.SeasonNumber, nil)

    for _, season := range serie.Seasons {
      nr := fetch_number(season.Name)

      if tv_season.SeasonNumber != nr {
        continue
      }

      hasSeason = true
      if err != nil {
        println(err.Error())
      }
      //load episodes
      DB.Model(season).Related(&season.Episodes)
      load_episodes(rlc, season, seasonInfo, tv)
      fetch_season(season, seasonInfo)
    }
    if !hasSeason {
      s := Season{}
      fetch_season(&s, seasonInfo)
      serie.Seasons = append(serie.Seasons, &s)
      s.Missing = true
    }
  }
  println("TMDb Seasons Import Done", tv.Name)
}

func fetch_episode(e *Episode, i *tmdb.TvEpisode) {
  e.AirDate = i.AirDate
  e.EpisodeNumber = i.EpisodeNumber
  e.OriginalName = i.Name
  e.StillPath = i.StillPath
}

func load_episodes(rlc *RequestLimitCheck, season *Season, seasonInfo *tmdb.TvSeason, tv *tmdb.TV) {
  for _, tv_episode := range seasonInfo.Episodes {
    hasEpisode := false

    rlc.checkRequest()
    episodeInfo, err := tmdn.GetTvEpisodeInfo(tv.ID, tv_episode.SeasonNumber, tv_episode.EpisodeNumber, nil)

    if err != nil {
      println(err.Error())
      println(tv.Name, tv_episode.SeasonNumber, tv_episode.EpisodeNumber)
    }

    for _, episode := range season.Episodes {
      nr := fetch_number(episode.Name)
      if nr != tv_episode.EpisodeNumber {
        continue
      }

      fetch_episode(episode, episodeInfo)
      //episode.Missing = false

      hasEpisode = true
    }
    if !hasEpisode {
      e := Episode{}
      fetch_episode(&e, episodeInfo)
      e.Missing = true
      season.Episodes = append(season.Episodes, &e)
    }
  }
}

func findSerie(name string) *tmdb.TvSearchResults {
  tmdn := getTmdb()
  result, err := tmdn.SearchTv(name, nil)
  if err != nil {
    println(err)
  }
  return result
}

