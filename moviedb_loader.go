package main

import (
  "github.com/ryanbradynd05/go-tmdb"
  "regexp"
  "strconv"
  "strings"
  "github.com/jinzhu/gorm"
)

var seriesDump []Serie

type MoviedbLoaderType interface {
  Tmdb() *tmdb.TMDb
  Rlc() *RequestLimitCheck
  DB() gorm.DB
}

type MoviedbLoader struct {
  RequestLimitcheck *RequestLimitCheck
  Tmdn *tmdb.TMDb
}

var rlc *RequestLimitCheck
func (this *MoviedbLoader) Rlc() *RequestLimitCheck {
  if rlc == nil {
    rlc = &RequestLimitCheck{}
    rlc.Reset()
  }
  if this.RequestLimitcheck == nil {
    this.RequestLimitcheck = rlc
  }
  return this.RequestLimitcheck
}
func (this *MoviedbLoader) Tmdb() *tmdb.TMDb {
  if this.Tmdb == nil {
    this.Tmdb = tmdb.Init(loadConfig().TMDb.API_KEY)
  }
  this.Rlc().CheckRequest()
  return this.Tmdb
}
func (this *MoviedbLoader) DB() gorm.DB {
  return DB
}

func loadTmdb(loader MoviedbLoaderType) {
  seriesDump = loadDump("./fetch/DATA_DUMP.json")

  loadSeries(loader)
}

func loadSeries(loader MoviedbLoaderType) {
  println("TMDb Series Import Start")

  for _, val := range seriesDump {
    loadSerie(loader, val.Name)
  }
  println("TMDb Series Import Done")
}

func loadSerie(loader MoviedbLoaderType, title string) {
  //tmdn := getTmdb()

  serie := Serie{Name: title}
  loader.DB().Where("Name = ?", title).First(&serie)

  //loader.Rlc().CheckRequest()
  searchTv, err := loader.Tmdb().SearchTv(title, nil)

  if err != nil {
    println("@@@@@@@@@@@@@@@")
    println(err.Error())
    println("@@@@@@@@@@@@@@@")
  }

  if len(searchTv.Results) == 0 {
    return
  }
  //loader.Rlc().CheckRequest()
  tv, _ := loader.Tmdb().GetTvInfo(searchTv.Results[0].ID, nil)
  applySerie(&serie, *tv)

  //fetch season data
  loader.DB().Model(serie).Related(&serie.Seasons)
  loadSeasons(loader, &serie, tv)

  loader.DB().Save(&serie)
}
func loadSeasons(loader MoviedbLoaderType, serie *Serie, tv *tmdb.TV) {
  println("TMDb Seasons Import Start", tv.Name)

  for _, tvSeason := range tv.Seasons {
    hasSeason := false

    //loader.Rlc().CheckRequest()
    seasonInfo, err := loader.Tmdb().GetTvSeasonInfo(tv.ID, tvSeason.SeasonNumber, nil)

    for _, season := range serie.Seasons {
      nr := fetchNumber(season.Name)

      if tvSeason.SeasonNumber != nr {
        continue
      }

      hasSeason = true
      if err != nil {
        println(err.Error())
      }
      //load episodes
      loader.DB().Model(season).Related(&season.Episodes)
      loadEpisodes(loader, season, seasonInfo, tv)
      applySeason(season, seasonInfo)
    }
    if !hasSeason {
      s := Season{}
      applySeason(&s, seasonInfo)
      serie.Seasons = append(serie.Seasons, &s)
      s.Missing = true
    }
  }
  println("TMDb Seasons Import Done", tv.Name)
}
func loadEpisodes(loader MoviedbLoaderType, season *Season, seasonInfo *tmdb.TvSeason, tv *tmdb.TV) {
  for _, tvEpisode := range seasonInfo.Episodes {
    hasEpisode := false

    episodeInfo, err := loader.Tmdb().GetTvEpisodeInfo(tv.ID, tvEpisode.SeasonNumber, tvEpisode.EpisodeNumber, nil)

    if err != nil {
      println(err.Error())
      println(tv.Name, tvEpisode.SeasonNumber, tvEpisode.EpisodeNumber)
    }

    for _, episode := range season.Episodes {
      nr := fetchNumber(episode.Name)
      if nr != tvEpisode.EpisodeNumber {
        continue
      }

      applyEpisode(episode, episodeInfo)

      hasEpisode = true
    }
    if !hasEpisode {
      e := Episode{}
      applyEpisode(&e, episodeInfo)
      e.Missing = true
      season.Episodes = append(season.Episodes, &e)
    }
  }
}

func fetchNumber(name string) int {
  name = strings.Trim(name, " ")
  regex := regexp.MustCompile("[Ss]?(\\d+)")
  ret, _ := strconv.Atoi(regex.ReplaceAllString(name, "$1"))
  return ret
}

func applySerie(serie *Serie, info tmdb.TV) {
  if len(info.PosterPath) > 0 {
    serie.PosterPath = info.PosterPath
  }
  if len(info.FirstAirDate) > 0 {
    serie.FirstAirDate = info.FirstAirDate
  }
  if info.VoteAverage > 0 {
    serie.VoteAverage = info.VoteAverage
  }
  if info.VoteCount > 0 {
    serie.VoteCount = info.VoteCount
  }
  if len(info.OriginalName) > 0 {
    serie.OriginalName = info.OriginalName
  }
  if len(info.Overview) > 0 {
    serie.Description = info.Overview
  }
  if info.ID > 0 {
    serie.Tmdb_id = info.ID
  }
}
func applySeason(season *Season, seasonInfo *tmdb.TvSeason) {
  if len(seasonInfo.Overview) > 0 {
    season.Description = seasonInfo.Overview
  }

  if len(seasonInfo.PosterPath) > 0 {
    season.PosterPath = seasonInfo.PosterPath
  }

  if seasonInfo.SeasonNumber > 0 {
    season.SeasonNumber = seasonInfo.SeasonNumber
  }

  if len(seasonInfo.AirDate) > 0 {
    season.AirDate = seasonInfo.AirDate
  }

  if len(seasonInfo.Name) > 0 {
    season.OriginalName = seasonInfo.Name
  }

  if seasonInfo.ID > 0 {
    season.Tmdb_id = seasonInfo.ID
  }
}
func applyEpisode(e *Episode, i *tmdb.TvEpisode) {
  if len(i.AirDate) > 0 {
    e.AirDate = i.AirDate
  }
  if i.EpisodeNumber > 0 {
    e.EpisodeNumber = i.EpisodeNumber
  }
  if len(i.Name) > 0 {
    e.OriginalName = i.Name
  }
  if len(i.StillPath) > 0 {
    e.StillPath = i.StillPath
  }
  if i.ID > 0 {
    e.Tmdb_id = i.ID
  }
}

func loadSeasonFromTMDB(loader MoviedbLoaderType, serieId, seasonNr int) (*tmdb.TvSeason, error) {
  info, err := loader.Tmdb().GetTvSeasonInfo(serieId, seasonNr, nil)
  return info, err
}
func loadEpisodeFromTMDB(loader MoviedbLoaderType, serieId, seasonNr, episodeNr int) *tmdb.TvEpisode {
  episodeInfo, err := loader.Tmdb().GetTvEpisodeInfo(serieId, seasonNr, episodeNr, nil)

  if err != nil {
    println(err)
  }

  return episodeInfo
}

func findSerie(loader MoviedbLoader, name string) *tmdb.TvSearchResults {
  result, err := loader.Tmdb().SearchTv(name, nil)
  if err != nil {
    println(err)
  }
  return result
}
func validTitle(title string) bool {
  name := strings.Trim(title, "? ")
  if len(name) > 0 {
    return true
  }
  return false
}