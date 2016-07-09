package loader

import (
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

type databaseInterface interface {
	NewRecord(interface{}) bool
	Save(interface{}) *gorm.DB
}

type tmdbInterface interface {
	SearchTv(name string, options map[string]string) (*tmdb.TvSearchResults, error)
	GetTvInfo(showid int, options map[string]string) (*tmdb.TV, error)
	GetTvSeasonInfo(showid, seasonid int, options map[string]string) (*tmdb.TvSeason, error)
}

func ImportData(db databaseInterface, series []*database.Serie) error {
	for _, val := range series {
		db.NewRecord(val)
	}
	return nil
}

func ImportTmdb(t tmdbInterface, series []*database.Serie) {
	for _, serie := range series {
		tvInfo := loadSerie(t, serie.Name)
		applyTmdbIds(serie, tvInfo)
		applySerieData(serie, tvInfo)

		loadSeasons(t, serie)
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
		panic("error searchtv")
	}
	CheckRequest("GetTvInfo")
	tvInfo, err := t.GetTvInfo(result.Results[0].ID, nil)
	if err != nil {
		panic("error getvinfo")
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
		seasonInfo, err := t.GetTvSeasonInfo(serie.TmdbId, season.TmdbId, nil)
		if err != nil {
			panic("GetTvSeasonInfo error")
		}
		applySeasonData(seasonInfo, season)
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

func loadSeries() {
	//println("TMDb Series Import Start")

	//for _, val := range seriesDump {
	//loadSerie(val.Name)
	//}
	//println("TMDb Series Import Done")
}

func findSerie(name string) *tmdb.TvSearchResults {
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

func applyEpisode(e *database.Episode, i *tmdb.TvEpisode) {
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
}

func loadEpisodes(season *database.Season, seasonInfo *tmdb.TvSeason, tv *tmdb.TV) {
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
	//e := database.Episode{}
	//applyEpisode(&e, episodeInfo)
	//e.Missing = true
	//season.Episodes = append(season.Episodes, &e)
	//}
	//}
}

func loadEpisode(serieId, seasonNr, episodeNr int) *tmdb.TvEpisode {
	return nil
	//tmdn := getTmdb()

	//rlc.checkRequest()
	//episodeInfo, err := tmdn.GetTvEpisodeInfo(serieId, seasonNr, episodeNr, nil)

	//if err != nil {
	//println(err)
	//}

	//return episodeInfo
}

func loadSeason(serieId, seasonNr int) (*tmdb.TvSeason, error) {
	return nil, nil
	//tmdn := getTmdb()

	//rlc.checkRequest()
	//info, err := tmdn.GetTvSeasonInfo(serieId, seasonNr, nil)

	//return info, err
}

func validTitle(title string) bool {
	name := strings.Trim(title, "? ")
	if len(name) > 0 {
		return true
	}
	return false
}
