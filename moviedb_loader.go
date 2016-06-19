package main

import (
	"strings"

	"github.com/ryanbradynd05/go-tmdb"
)

var tmdn *tmdb.TMDb

var seriesDump []Serie

func getTmdb() *tmdb.TMDb {
	if tmdn == nil {
		tmdn = tmdb.Init(loadConfig().TMDb.API_KEY)
		rlc = &RequestLimitCheck{}
		rlc.reset()
	}
	return tmdn
}

func LoadTmdb() error {
	return nil
}

//func loadTmdb() {
//tmdn = getTmdb()
//seriesDump = loadDump("./fetch/DATA_DUMP.json")

//loadSeries()
//}

//func loadSeries() {
//println("TMDb Series Import Start")

//for _, val := range seriesDump {
//loadSerie(val.Name)
//}
//println("TMDb Series Import Done")
//}

//func loadSerie(title string) {
//tmdn := getTmdb()

//serie := Serie{Name: title}
//DB.Where("Name = ?", title).First(&serie)

//rlc.checkRequest()
//searchTv, err := tmdn.SearchTv(title, nil)

//if err != nil {
//println("@@@@@@@@@@@@@@@")
//println(err.Error())
//println("@@@@@@@@@@@@@@@")
//}

//if len(searchTv.Results) == 0 {
//return
//}
//rlc.checkRequest()
//tv, _ := tmdn.GetTvInfo(searchTv.Results[0].ID, nil)
//applySerie(&serie, *tv)

////fetch season data
//DB.Model(serie).Related(&serie.Seasons)
//loadSeasons(&serie, tv)

//DB.Save(&serie)
//}

//func findSerie(name string) *tmdb.TvSearchResults {
//tmdn := getTmdb()

//rlc.checkRequest()
//result, err := tmdn.SearchTv(name, nil)
//if err != nil {
//println(err)
//}
//return result
//}

//func fetchNumber(name string) int {
//name = strings.Trim(name, " ")
//regex := regexp.MustCompile("[Ss]?(\\d+)")
//ret, _ := strconv.Atoi(regex.ReplaceAllString(name, "$1"))
//return ret
//}

//func applySerie(serie *Serie, info tmdb.TV) {
//if len(info.PosterPath) > 0 {
//serie.PosterPath = info.PosterPath
//}
//if len(info.FirstAirDate) > 0 {
//serie.FirstAirDate = info.FirstAirDate
//}
//if info.VoteAverage > 0 {
//serie.VoteAverage = info.VoteAverage
//}
//if info.VoteCount > 0 {
//serie.VoteCount = info.VoteCount
//}
//if len(info.OriginalName) > 0 {
//serie.OriginalName = info.OriginalName
//}
//if len(info.Overview) > 0 {
//serie.Description = info.Overview
//}
//if info.ID > 0 {
//serie.Tmdb_id = info.ID
//}
//}

//func applySeason(season *Season, seasonInfo *tmdb.TvSeason) {
//if len(seasonInfo.Overview) > 0 {
//season.Description = seasonInfo.Overview
//}

//if len(seasonInfo.PosterPath) > 0 {
//season.PosterPath = seasonInfo.PosterPath
//}

//if seasonInfo.SeasonNumber > 0 {
//season.SeasonNumber = seasonInfo.SeasonNumber
//}

//if len(seasonInfo.AirDate) > 0 {
//season.AirDate = seasonInfo.AirDate
//}

//if len(seasonInfo.Name) > 0 {
//season.OriginalName = seasonInfo.Name
//}

//if seasonInfo.ID > 0 {
//season.Tmdb_id = seasonInfo.ID
//}
//}

//func loadSeasons(serie *Serie, tv *tmdb.TV) {
//println("TMDb Seasons Import Start", tv.Name)

//for _, tvSeason := range tv.Seasons {
//hasSeason := false

//rlc.checkRequest()
//seasonInfo, err := tmdn.GetTvSeasonInfo(tv.ID, tvSeason.SeasonNumber, nil)

//for _, season := range serie.Seasons {
//nr := fetchNumber(season.Name)

//if tvSeason.SeasonNumber != nr {
//continue
//}

//hasSeason = true
//if err != nil {
//println(err.Error())
//}
////load episodes
//DB.Model(season).Related(&season.Episodes)
//loadEpisodes(season, seasonInfo, tv)
//applySeason(season, seasonInfo)
//}
//if !hasSeason {
//s := Season{}
//applySeason(&s, seasonInfo)
//serie.Seasons = append(serie.Seasons, &s)
//s.Missing = true
//}
//}
//println("TMDb Seasons Import Done", tv.Name)
//}

//func applyEpisode(e *Episode, i *tmdb.TvEpisode) {
//if len(i.AirDate) > 0 {
//e.AirDate = i.AirDate
//}
//if i.EpisodeNumber > 0 {
//e.EpisodeNumber = i.EpisodeNumber
//}
//if len(i.Name) > 0 {
//e.OriginalName = i.Name
//}
//if len(i.StillPath) > 0 {
//e.StillPath = i.StillPath
//}
//if i.ID > 0 {
//e.Tmdb_id = i.ID
//}
//}

//func loadEpisodes(season *Season, seasonInfo *tmdb.TvSeason, tv *tmdb.TV) {
//for _, tvEpisode := range seasonInfo.Episodes {
//hasEpisode := false

//rlc.checkRequest()
//episodeInfo, err := tmdn.GetTvEpisodeInfo(tv.ID, tvEpisode.SeasonNumber, tvEpisode.EpisodeNumber, nil)

//if err != nil {
//println(err.Error())
//println(tv.Name, tvEpisode.SeasonNumber, tvEpisode.EpisodeNumber)
//}

//for _, episode := range season.Episodes {
//nr := fetchNumber(episode.Name)
//if nr != tvEpisode.EpisodeNumber {
//continue
//}

//applyEpisode(episode, episodeInfo)
////episode.Missing = false

//hasEpisode = true
//}
//if !hasEpisode {
//e := Episode{}
//applyEpisode(&e, episodeInfo)
//e.Missing = true
//season.Episodes = append(season.Episodes, &e)
//}
//}
//}

//func loadEpisode(serieId, seasonNr, episodeNr int) *tmdb.TvEpisode {
//tmdn := getTmdb()

//rlc.checkRequest()
//episodeInfo, err := tmdn.GetTvEpisodeInfo(serieId, seasonNr, episodeNr, nil)

//if err != nil {
//println(err)
//}

//return episodeInfo
//}

//func loadSeason(serieId, seasonNr int) (*tmdb.TvSeason, error) {
//tmdn := getTmdb()

//rlc.checkRequest()
//info, err := tmdn.GetTvSeasonInfo(serieId, seasonNr, nil)

//return info, err
//}

func validTitle(title string) bool {
	name := strings.Trim(title, "? ")
	if len(name) > 0 {
		return true
	}
	return false
}
