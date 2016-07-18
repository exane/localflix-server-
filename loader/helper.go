package loader

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/exane/localflix-server-/database"
	"github.com/ryanbradynd05/go-tmdb"
)

func hasSeasonLoaded(seasonNr int, seasons []*database.Season) bool {
	for _, season := range seasons {
		if fetchNumber(season.Name) == seasonNr {
			return season.TmdbId != 0
		}
	}
	return false
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

func getTvSeasonInfo(t tmdbInterface, serie *database.Serie, tvInfo *tmdb.TV, seasonIndex int) *tmdb.TvSeason {
	season := tvInfo.Seasons[seasonIndex]
	CheckRequest("GetTvSeasonInfo")
	seasonInfo, err := t.GetTvSeasonInfo(serie.TmdbId, season.SeasonNumber, nil)
	if err != nil {
		fmt.Printf("\nError: %s\n%v\n%v\n", err.Error(), serie, season)
	}
	return seasonInfo
}

func searchTv(t tmdbInterface, name string) *tmdb.TvSearchResults {
	CheckRequest("SearchTv")
	result, err := t.SearchTv(name, nil)
	if err != nil {
		fmt.Printf("\nError: %s\n", err.Error())
	}
	return result
}

func getTvInfo(t tmdbInterface, result *tmdb.TvSearchResults) *tmdb.TV {
	CheckRequest("GetTvInfo")
	tvInfo, err := t.GetTvInfo(result.Results[0].ID, nil)
	if err != nil {
		fmt.Printf("\nError: %s\n", err.Error())
	}
	return tvInfo
}

func getTvEpisodeInfo(t tmdbInterface, showID int, season *tmdb.TvSeason, episode tmdb.TvEpisode) *tmdb.TvEpisode {
	CheckRequest("GetTvEpisodeInfo")
	episodeInfo, err := t.GetTvEpisodeInfo(showID, season.SeasonNumber, episode.EpisodeNumber, nil)
	if err != nil {
		fmt.Printf("\nError: %s\n%v\n%v\n%v\n", err.Error(), showID, season, episode)
	}
	return episodeInfo
}
