package main

import (
  "encoding/json"
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "io/ioutil"
)

var DB gorm.DB

type Config struct {
  Database   struct {
               Root     string
               Password string
               Db       string
             }
  Server     struct {
               Url  string
               Port string
             }
  Fileserver struct {
               Url            string
               Root_directory string
               Port           string
             }
  TMDb       struct {
               API_KEY string
             }
}

var config *Config

func loadConfig() *Config {
  if config != nil {
    return config
  }
  config = &Config{}
  js, _ := ioutil.ReadFile("./config.json")
  json.Unmarshal(js, config)
  return config
}

func initDb() {
  loadConfig()
  db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", config.Database.Root, config.Database.Password, config.Database.Db))
  DB = *db
  DB.LogMode(true)

  if err != nil {
    panic("failed to connect database")
  }
}

func createTables() {
  DB.DropTableIfExists(Episode{}, Serie{}, Season{}, Subtitle{})
  DB.CreateTable(Episode{}, Serie{}, Season{}, Subtitle{})
  //DB.AutoMigrate(Episode{}, Serie{}, Season{})
}

func dumpImport() {
  data := loadDump("./fetch/DATA_DUMP.json")

  for _, val := range data {
    DB.Create(&val)
  }
}

func loadDump(file string) []Serie {
  js, _ := ioutil.ReadFile(file)
  data := []Serie{}
  json.Unmarshal(js, &data)
  return data
}

func updateDb() {
  data := loadDump("./fetch/DATA_DUMP.json")

  for _, serie_data := range data {
    serie := Serie{}
    notFound := DB.Where("name = ?", serie_data.Name).Find(&serie).RecordNotFound()

    if notFound {
      println("new serie! create:", serie_data.Name)
      loadSerie(serie_data.Name)
    }

    updateSeasons(&serie, serie_data)
  }
}

func updateSeasons(serie *Serie, serie_data Serie) {
  for _, season_data := range serie_data.Seasons {
    season := Season{}
    notFound := DB.Where("serie_id = ? and name = ?", serie.ID, season_data.Name).Find(&season).RecordNotFound()

    if notFound {
      println("new season! create:", serie.Name, season_data.Name)
      season.SeasonNumber = fetchNumber(season_data.Name)
      season.Name = season_data.Name
      info, _ := loadSeasonFromTMDB(serie.Tmdb_id, season.SeasonNumber)
      applySeason(&season, info)
      DB.Model(&serie).Association("Seasons").Append(&season)
    }

    updateEpisodes(serie, &season, season_data)
  }
}

func updateEpisodes(serie *Serie, season *Season, season_data *Season) {
  for _, episode_data := range season_data.Episodes {
    //src + name as unique key
    episode := Episode{}
    notFound := DB.Where("name = ? AND src = ? AND missing = ?", episode_data.Name, episode_data.Src, 0).Find(&episode).RecordNotFound()
    if notFound {
      //create new episode
      println("new episode! create:", episode_data.Src, episode_data.Name)
      //find tmdb entry
      episodeNumber := fetchNumber(episode_data.Name)
      notExist := DB.Where("missing = ? and season_id = ? and episode_number = ?", 1, season.ID, episodeNumber).
      Find(&episode).RecordNotFound()

      episode.Missing = false
      episode.Name = episode_data.Name
      episode.Src = episode_data.Src
      episode.Extension = episode_data.Extension
      episode.Subtitles = episode_data.Subtitles
      episode.EpisodeNumber = episodeNumber

      if notExist {
        //tmdb entry does not exist
        //load tmdb
        episodeInfo := loadEpisodeFromTMDB(serie.Tmdb_id, season.SeasonNumber, episode.EpisodeNumber)
        applyEpisode(&episode, episodeInfo)
        DB.Model(&season).Association("Episodes").Append(&episode)
      } else {
        //tmdb entry exist, but no actual video file is available
        DB.Save(&episode)
      }

      println("updated", episode.ID, episode_data.Src, episode_data.Name)
    }
  }
}