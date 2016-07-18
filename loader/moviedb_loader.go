package loader

import (
	"fmt"

	"github.com/exane/localflix-server-/config"
	"github.com/exane/localflix-server-/database"
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

func Import(series []*database.Serie) {
	ImportData(series)
	ImportTmdb(getTmdb(), series)
}

func Update(series []*database.Serie) {
	UpdateDB(getTmdb(), series)
}

type tmdbInterface interface {
	SearchTv(name string, options map[string]string) (*tmdb.TvSearchResults, error)
	GetTvInfo(showid int, options map[string]string) (*tmdb.TV, error)
	GetTvSeasonInfo(showid, seasonid int, options map[string]string) (*tmdb.TvSeason, error)
	GetTvEpisodeInfo(showID, seasonNum, episodeNum int, options map[string]string) (*tmdb.TvEpisode, error)
}

func ImportData(series []*database.Serie) error {
	db := database.DB
	for _, val := range series {
		db.FirstOrCreate(val, "name = ? ", val.Name)
	}
	return nil
}

func ImportTmdb(t tmdbInterface, series []*database.Serie) {
	for _, serie := range series {
		if !IsTesting {
			fmt.Printf("\nTMDb load serie %s\n", serie.Name)
		}
		tvInfo := loadSerie(t, serie.Name)
		applySerieData(serie, tvInfo)
		loadSeasons(t, tvInfo, serie)

		if !IsTesting {
			fmt.Printf("\nTMDb finished loading serie %s\n", serie.Name)
		}
		database.DB.Save(serie)
	}
}

func UpdateDB(t tmdbInterface, series []*database.Serie) {
	db := database.DB

	for _, serie := range series {
		db.Find(&serie, "name = ?", serie.Name)

		if !IsTesting {
			fmt.Printf("\nTMDb load serie %s\n", serie.Name)
		}
		tvInfo := loadSerie(t, serie.Name)
		if serie.TmdbId == 0 {
			applySerieData(serie, tvInfo)
		}

		loadSeasons(t, tvInfo, serie)
		if !IsTesting {
			fmt.Printf("\nTMDb finished loading serie %s\n", serie.Name)
		}
		database.DB.Save(serie)
	}
}

func loadSerie(t tmdbInterface, name string) *tmdb.TV {
	result := searchTv(t, name)

	if len(result.Results) == 0 {
		return &tmdb.TV{}
	}

	tvInfo := getTvInfo(t, result)

	return tvInfo
}

func fetchCurrentSeason(tmdbid int, seasons []*database.Season) *database.Season {
	for _, val := range seasons {
		if tmdbid == val.TmdbId {
			return val
		}
	}
	return nil
}

func debugRecover(args ...interface{}) {
	if r := recover(); r != nil {
		_ = "breakpoint"
	}
}

func loadSeasons(t tmdbInterface, tvInfo *tmdb.TV, serie *database.Serie) {
	if tvInfo == nil {
		return
	}
	for index := range tvInfo.Seasons {
		seasonInfo := getTvSeasonInfo(t, serie, tvInfo, index)
		season := applySeasonData(seasonInfo, &serie.Seasons)

		loadEpisodes(t, serie.TmdbId, seasonInfo, season)
	}
	sortSeasons(serie.Seasons)
}

func loadEpisodes(t tmdbInterface, showID int, seasonInfo *tmdb.TvSeason, season *database.Season) {
	if season.Missing {
		return
	}
	for _, episode := range seasonInfo.Episodes {
		episodeInfo := getTvEpisodeInfo(t, showID, seasonInfo, episode)
		applyEpisodeData(&season.Episodes, episodeInfo)
	}
	sortEpisodes(season.Episodes)
}

func applySerieData(serie *database.Serie, info *tmdb.TV) {
	if info == nil {
		return
	}
	serie.TmdbId = info.ID
	serie.OriginalName = info.OriginalName
	serie.Description = info.Overview
	serie.PosterPath = info.PosterPath
	serie.VoteAverage = info.VoteAverage
	serie.VoteCount = info.VoteCount
	serie.FirstAirDate = info.FirstAirDate
}

func applySeasonData(info *tmdb.TvSeason, seasons *[]*database.Season) *database.Season {
	var season *database.Season

	for _, val := range *seasons {
		if fetchNumber(val.Name) == info.SeasonNumber {
			season = val
			break
		}
	}
	if season == nil {
		season = &database.Season{Missing: true, Name: fmt.Sprintf("S%d", info.SeasonNumber)}
		*seasons = append(*seasons, season)
	}
	season.TmdbId = info.ID
	season.AirDate = info.AirDate
	season.OriginalName = info.Name
	season.Description = info.Overview
	season.PosterPath = info.PosterPath
	season.SeasonNumber = info.SeasonNumber

	return season
}

func applyEpisodeData(episodes *[]*database.Episode, info *tmdb.TvEpisode) {
	if info == nil {
		return
	}
	var episode *database.Episode
	for _, val := range *episodes {
		if val.EpisodeNumber == 0 {
			val.EpisodeNumber = fetchNumber(val.Name)
		}
		if info.EpisodeNumber == fetchNumber(val.Name) {
			episode = val
			break
		}
	}
	if episode == nil {
		episode = &database.Episode{Missing: true}
		*episodes = append(*episodes, episode)
	}

	episode.TmdbId = info.ID
	episode.OriginalName = info.Name
	episode.AirDate = info.AirDate
	episode.EpisodeNumber = info.EpisodeNumber
	episode.StillPath = info.StillPath
	episode.Description = info.Overview
}
