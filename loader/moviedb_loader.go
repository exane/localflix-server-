package loader

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/exane/localflix-server-/config"
	"github.com/exane/localflix-server-/database"
	"github.com/jinzhu/gorm"
	"github.com/ryanbradynd05/go-tmdb"
)

var tmdn *tmdb.TMDb

var seriesDump []database.Serie

func getTmdb() *tmdb.TMDb {
	if tmdn == nil {
		tmdn = tmdb.Init(config.LoadConfig().TMDb.ApiKey)
		Reset()
	}
	return tmdn
}

func Import(db databaseInterface, series []*database.Serie) {
	ImportData(db, series)
	ImportTmdb(db, getTmdb(), series)
}

type databaseInterface interface {
	NewRecord(interface{}) bool
	Save(interface{}) *gorm.DB
}

type tmdbInterface interface {
	SearchTv(name string, options map[string]string) (*tmdb.TvSearchResults, error)
	GetTvInfo(showid int, options map[string]string) (*tmdb.TV, error)
	GetTvSeasonInfo(showid, seasonid int, options map[string]string) (*tmdb.TvSeason, error)
	GetTvEpisodeInfo(showID, seasonNum, episodeNum int, options map[string]string) (*tmdb.TvEpisode, error)
}

func ImportData(db databaseInterface, series []*database.Serie) error {
	for _, val := range series {
		db.NewRecord(val)
	}
	return nil
}

func ImportTmdb(db databaseInterface, t tmdbInterface, series []*database.Serie) {
	for _, serie := range series {
		if !IsTesting {
			fmt.Printf("\nTMDb load serie %s\n", serie.Name)
		}
		tvInfo := loadSerie(t, serie.Name)
		applyTmdbIds(serie, tvInfo)
		applySerieData(serie, tvInfo)

		loadSeasons(t, serie)
		if !IsTesting {
			fmt.Printf("\nTMDb finished loading serie %s\n", serie.Name)
		}
		db.Save(serie)
	}
}

func UpdateDB(db databaseInterface, series []*database.Serie) {
	db.Save(series)
}

func applyTmdbIds(serie *database.Serie, info *tmdb.TV) {
	serie.TmdbId = info.ID
	applyTmdbIdsToSeasons(serie, info)
}

func applyTmdbIdsToSeasons(serie *database.Serie, info *tmdb.TV) {
	for _, season := range serie.Seasons {
		seasonNr := fetchNumber(season.Name)
		season.TmdbId = getTmdbIdFromSeasons(seasonNr, info)
	}
}

func loadSerie(t tmdbInterface, name string) *tmdb.TV {
	CheckRequest("SearchTv")
	result, err := t.SearchTv(name, nil)
	if err != nil {
		fmt.Printf("\nError: %s\n", err.Error())
	}

	if len(result.Results) == 0 {
		return &tmdb.TV{}
	}

	CheckRequest("GetTvInfo")
	tvInfo, err := t.GetTvInfo(result.Results[0].ID, nil)
	if err != nil {
		fmt.Printf("\nError: %s\n", err.Error())
	}
	return tvInfo
}

func applySerieData(serie *database.Serie, info *tmdb.TV) {
	serie.OriginalName = info.OriginalName
	serie.Description = info.Overview
	serie.PosterPath = info.PosterPath
	serie.VoteAverage = info.VoteAverage
	serie.VoteCount = info.VoteCount
	serie.FirstAirDate = info.FirstAirDate
}

func loadSeasons(t tmdbInterface, serie *database.Serie) {
	for _, season := range serie.Seasons {
		CheckRequest("GetTvSeasonInfo")
		seasonInfo, err := t.GetTvSeasonInfo(serie.TmdbId, fetchNumber(season.Name), nil)
		if err != nil {
			fmt.Printf("\nError: %s\n%v\n%v\n", err.Error(), serie, season)
		}
		applySeasonData(seasonInfo, season)

		loadEpisodes(t, serie.TmdbId, season)
	}
}

func applySeasonData(info *tmdb.TvSeason, season *database.Season) {
	if info == nil {
		return
	}
	season.AirDate = info.AirDate
	season.OriginalName = info.Name
	season.Description = info.Overview
	season.PosterPath = info.PosterPath
	season.SeasonNumber = info.SeasonNumber
}

func getTmdbIdFromSeasons(seasonNumber int, info *tmdb.TV) int {
	for _, val := range info.Seasons {
		if val.SeasonNumber == seasonNumber {
			return val.ID
		}
	}
	return -1
}

func FindSerie(name string) *tmdb.TvSearchResults {
	//tmdn := getTmdb()

	//rlc.checkRequest()
	//result, err := tmdn.SearchTv(name, nil)
	//if err != nil {
	//println(err)
	//}
	//return result
	return nil
}

func fetchNumber(name string) int {
	name = strings.Trim(name, " ")
	regex := regexp.MustCompile("[Ss]?(\\d+)")
	ret, _ := strconv.Atoi(regex.ReplaceAllString(name, "$1"))
	return ret
}

func applyEpisodeData(episode *database.Episode, info *tmdb.TvEpisode) {
	if info == nil {
		return
	}
	episode.TmdbId = info.ID
	episode.OriginalName = info.Name
	episode.AirDate = info.AirDate
	episode.EpisodeNumber = info.EpisodeNumber
	episode.StillPath = info.StillPath
	episode.Description = info.Overview
}

func loadEpisodes(t tmdbInterface, showID int, season *database.Season) {
	for _, episode := range season.Episodes {
		episodeNum := fetchNumber(episode.Name)
		CheckRequest("GetTvEpisodeInfo")
		episodeInfo, err := t.GetTvEpisodeInfo(showID, season.SeasonNumber, episodeNum, nil)

		if err != nil {
			fmt.Printf("\nError: %s\n%v\n%v\n%v\n", err.Error(), showID, season, episode)
			return
		}

		applyEpisodeData(episode, episodeInfo)
	}
}

func ValidTitle(title string) bool {
	name := strings.Trim(title, "? ")
	if len(name) > 0 {
		return true
	}
	return false
}
