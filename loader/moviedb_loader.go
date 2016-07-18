package loader

import (
	"fmt"
	"regexp"
	"sort"
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
	Create(value interface{}) *gorm.DB
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
		db.Create(val)
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

		loadSeasons(t, tvInfo, serie)
		if !IsTesting {
			fmt.Printf("\nTMDb finished loading serie %s\n", serie.Name)
		}
		db.Save(serie)
	}
}

func UpdateDB(db databaseInterface, series []*database.Serie) {
	for _, val := range series {
		_ = "breakpoint"
		created := db.NewRecord(val)
		if !created {
			db.Save(val)
		}
	}
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

func loadSeasons(t tmdbInterface, tvInfo *tmdb.TV, serie *database.Serie) {
	for _, season := range tvInfo.Seasons {
		CheckRequest("GetTvSeasonInfo")
		seasonInfo, err := t.GetTvSeasonInfo(serie.TmdbId, season.SeasonNumber, nil)
		if err != nil {
			fmt.Printf("\nError: %s\n%v\n%v\n", err.Error(), serie, season)
		}

		s := applySeasonData(seasonInfo, &serie.Seasons)
		loadEpisodes(t, serie.TmdbId, seasonInfo, s)
		sortSeasons(serie.Seasons)
	}
}

func sortSeasons(seasons []*database.Season) {
	sort.Sort(seasonSort(seasons))
}

type seasonSort []*database.Season

func (s seasonSort) Len() int {
	return len(s)
}

func (s seasonSort) Less(i, j int) bool {
	return s[i].SeasonNumber < s[j].SeasonNumber
}

func (s seasonSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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
	season.AirDate = info.AirDate
	season.OriginalName = info.Name
	season.Description = info.Overview
	season.PosterPath = info.PosterPath
	season.SeasonNumber = info.SeasonNumber

	return season
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

func loadEpisodes(t tmdbInterface, showID int, seasonInfo *tmdb.TvSeason, season *database.Season) {
	if season.Missing {
		return
	}
	for _, episode := range seasonInfo.Episodes {
		CheckRequest("GetTvEpisodeInfo")
		episodeInfo, err := t.GetTvEpisodeInfo(showID, season.SeasonNumber, episode.EpisodeNumber, nil)
		if err != nil {
			fmt.Printf("\nError: %s\n%v\n%v\n%v\n", err.Error(), showID, season, episode)
			continue
		}
		applyEpisodeData(&season.Episodes, episodeInfo)
	}
	sortEpisodes(season.Episodes)
}

func sortEpisodes(episodes []*database.Episode) {
	sort.Sort(episodeSort(episodes))
}

type episodeSort []*database.Episode

func (s episodeSort) Len() int {
	return len(s)
}

func (s episodeSort) Less(i, j int) bool {
	return s[i].EpisodeNumber < s[j].EpisodeNumber
}

func (s episodeSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func ValidTitle(title string) bool {
	name := strings.Trim(title, "? ")
	if len(name) > 0 {
		return true
	}
	return false
}
