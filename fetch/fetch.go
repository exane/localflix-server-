package main

import (
  "io/ioutil"
  "encoding/json"
  "os"
  "regexp"
)

const (
  PATH = "Z:/serien"
  OUT_PATH = "Y:/golangWorkspace/src/github.com/exane/localflix/fetch"
  OUT = "DATA_DUMP"
  VIDEO = ".*[.](avi|web|mkv|mp4)$"
  SUBTITLE = ".*[.](srt|idx|sub|sfv)$"
)

type Serie struct {
  Name    string
  Seasons []*Season
}

type Season struct {
  Name     string
  Episodes []*Episode
}

type Episode struct {
  Name      string
  Extension string
  Subtitles []Subtitle
  Src       string
}

type Subtitle struct {
  Name string
}

func (season *Season) findEpisode(name string) (*Episode, bool) {
  for _, val := range season.Episodes {
    if val.Name == name {
      return val, true
    }
  }
  return &Episode{}, false
}

func main() {
  series := []Serie{}

  files, _ := ioutil.ReadDir(PATH)
  for i, val := range files {
    println(val.Name())
    series = append(series, Serie{Name: val.Name()})
    series[i].fetchSeasons(PATH)
  }
  js, _ := json.Marshal(series)
  ioutil.WriteFile(OUT_PATH + "/" + OUT + ".json", []byte(js), 0666)
}

func (serie *Serie) fetchSeasons(path string) {
  newPath := path + "/" + serie.Name
  files, _ := ioutil.ReadDir(newPath)

  for i, val := range files {
    serie.Seasons = append(serie.Seasons, &Season{Name: val.Name()})
    serie.Seasons[i].fetchEpisodes(newPath)
  }
}

func (season *Season) fetchEpisodes(path string) {
  newPath := path + "/" + season.Name
  files, _ := ioutil.ReadDir(newPath)

  for _, val := range files {
    episode := &Episode{
      Name: filename(val.Name()),
      Extension: "",
      Subtitles: []Subtitle{},
    }

    if isValid(VIDEO, val) {
      ep, exist := season.findEpisode(filename(val.Name()))
      if !exist {
        season.Episodes = append(season.Episodes, episode)
        ep = episode
      }
      if len(ep.Extension) == 0 {
        ep.Extension = extension(val.Name())
      }
      ep.Src = newPath
    }

    if isValid(SUBTITLE, val) {
      name := filename(val.Name())

      ep, exist := season.findEpisode(name)
      if !exist {
        season.Episodes = append(season.Episodes, episode)
        ep = episode
      }
      ep.Subtitles = append(ep.Subtitles, Subtitle{Name: val.Name()})
    }
  }
}

func isValid(TYPE string, file os.FileInfo) bool {
  matched, _ := regexp.MatchString(TYPE, file.Name())
  return matched
}

func filename(name string) string {
  regex := regexp.MustCompile("^(.*)[.].*")

  return regex.ReplaceAllString(name, "$1")
}

func extension(name string) string {
  regex := regexp.MustCompile(".*[.](.*)$")

  return regex.ReplaceAllString(name, "$1")
}