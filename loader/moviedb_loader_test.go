package loader_test

import (
	"errors"
	"reflect"

	"github.com/exane/localflix-server-/database"
	"github.com/exane/localflix-server-/loader"
	"github.com/exane/localflix-server-/loader/loaderfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanbradynd05/go-tmdb"
)

var _ = Describe("MoviedbLoader", func() {
	var db *loaderfakes.FakeDatabaseInterface
	var series []*database.Serie
	var tmdbMock *loaderfakes.FakeTmdbInterface

	BeforeEach(func() {
		db = &loaderfakes.FakeDatabaseInterface{}
		db.NewRecordReturns(true)

		series = []*database.Serie{
			&database.Serie{Name: "got", Seasons: []*database.Season{
				{Name: "s1", Episodes: []*database.Episode{{Name: "01"}, {Name: "02"}}},
				{Name: "s2", Episodes: []*database.Episode{{Name: "01"}}},
			}},
			&database.Serie{Name: "vikings", Seasons: []*database.Season{
				{Name: "s1", Episodes: []*database.Episode{{Name: "01"}}},
				{Name: "s2", Episodes: []*database.Episode{{Name: "01"}}},
			}},
		}

		loader.IsTesting = true
		loader.Requested = make(map[string]int)
		tmdbMock = &loaderfakes.FakeTmdbInterface{}
		tmdbMock.GetTvInfoStub = func(id int, options map[string]string) (*tmdb.TV, error) {
			if loader.Requests() > loader.LIMIT_REQUEST {
				return nil, errors.New("tmdb limit reached")
			}
			var result1 *tmdb.TV
			if id == 1 {
				result1 = &tmdb.TV{
					Overview:        "got desc",
					OriginalName:    "got tmdb",
					Name:            "got name",
					ID:              1,
					PosterPath:      "xyz",
					VoteAverage:     10.0,
					VoteCount:       1000,
					FirstAirDate:    "1.1.2010",
					NumberOfSeasons: 3,
					Seasons: []struct {
						AirDate      string `json:"air_date"`
						EpisodeCount int    `json:"episode_count"`
						ID           int
						PosterPath   string `json:"poster_path"`
						SeasonNumber int    `json:"season_number"`
					}{
						{ID: 100, SeasonNumber: 1},
						{ID: 101, SeasonNumber: 2},
						{ID: 102, SeasonNumber: 3},
					},
				}
			}
			if id == 2 {
				result1 = &tmdb.TV{
					Overview:        "",
					OriginalName:    "",
					Name:            "",
					ID:              2,
					PosterPath:      "",
					VoteAverage:     0,
					VoteCount:       0,
					NumberOfSeasons: 2,
					Seasons: []struct {
						AirDate      string `json:"air_date"`
						EpisodeCount int    `json:"episode_count"`
						ID           int
						PosterPath   string `json:"poster_path"`
						SeasonNumber int    `json:"season_number"`
					}{
						{ID: 102, SeasonNumber: 1},
						{ID: 103, SeasonNumber: 2},
					},
				}
			}
			return result1, nil
		}

		tmdbMock.SearchTvStub = func(name string, options map[string]string) (*tmdb.TvSearchResults, error) {
			id := 0
			if loader.Requests() > loader.LIMIT_REQUEST {
				return nil, errors.New("tmdb limit reached")
			}

			if name == "got" {
				id = 1
			}
			if name == "vikings" {
				id = 2
			}

			result_empty := &tmdb.TvSearchResults{
				Results: []struct {
					BackdropPath  string `json:"backdrop_path"`
					ID            int
					OriginalName  string   `json:"original_name"`
					FirstAirDate  string   `json:"first_air_date"`
					OriginCountry []string `json:"origin_country"`
					PosterPath    string   `json:"poster_path"`
					Popularity    float32
					Name          string
					VoteAverage   float32 `json:"vote_average"`
					VoteCount     uint32  `json:"vote_count"`
				}{},
			}

			if name == "empty" {
				return result_empty, nil
			}

			result1 := &tmdb.TvSearchResults{
				Results: []struct {
					BackdropPath  string `json:"backdrop_path"`
					ID            int
					OriginalName  string   `json:"original_name"`
					FirstAirDate  string   `json:"first_air_date"`
					OriginCountry []string `json:"origin_country"`
					PosterPath    string   `json:"poster_path"`
					Popularity    float32
					Name          string
					VoteAverage   float32 `json:"vote_average"`
					VoteCount     uint32  `json:"vote_count"`
				}{
					{ID: id},
				},
			}
			return result1, nil
		}
	})

	Describe("ImportData", func() {
		It("should not throw an error", func() {
			err := loader.ImportData(db, series)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should call db.NewRecord", func() {
			loader.ImportData(db, series)
			Expect(db.NewRecordCallCount()).ToNot(Equal(0))
		})

		It("should create 2 series", func() {
			loader.ImportData(db, series)
			Expect(db.NewRecordCallCount()).To(Equal(2))
			Expect(reflect.DeepEqual(db.NewRecordArgsForCall(0), series[0])).To(Equal(true))
			Expect(reflect.DeepEqual(db.NewRecordArgsForCall(1), series[1])).To(Equal(true))
		})
	})

	Describe("ImportTmdb", func() {
		It("should not panic", func() {
			series = nil
			for i := 0; i < 10; i++ {
				series = append(series, &database.Serie{Name: "got"})
			}
			Expect(func() {
				loader.ImportTmdb(db, tmdbMock, series)
			}).ToNot(Panic())

			series = nil
			for i := 0; i < 50; i++ {
				series = append(series, &database.Serie{Name: "got"})
			}
			Expect(func() {
				loader.ImportTmdb(db, tmdbMock, series)
			}).ToNot(Panic())
		})

		It("should call check rlc for SearchTv", func() {
			loader.ImportTmdb(db, tmdbMock, series)

			Expect(loader.Requested["SearchTv"]).To(Equal(2))
		})

		It("should call check rlc for GetTvInfo", func() {
			loader.ImportTmdb(db, tmdbMock, series)

			Expect(loader.Requested["GetTvInfo"]).To(Equal(2))
		})

		It("should ignore empty results", func() {
			series[0].Name = "empty"
			Expect(func() {
				loader.ImportTmdb(db, tmdbMock, series)
			}).ToNot(Panic())
		})

		It("should fetch entities from tmdb", func() {
			loader.ImportTmdb(db, tmdbMock, series)

			Expect(tmdbMock.SearchTvCallCount()).To(Equal(2))
			got, _ := tmdbMock.SearchTvArgsForCall(0)
			vikings, _ := tmdbMock.SearchTvArgsForCall(1)
			Expect(got).To(Equal("got"))
			Expect(vikings).To(Equal("vikings"))
		})

		It("should apply on series", func() {
			loader.ImportTmdb(db, tmdbMock, series)

			got := series[0]
			vikings := series[1]
			Expect(got.Name).To(Equal("got"))
			Expect(got.TmdbId).To(Equal(1))
			Expect(got.OriginalName).To(Equal("got tmdb"))
			Expect(got.Description).To(Equal("got desc"))
			Expect(got.PosterPath).To(Equal("xyz"))
			Expect(got.VoteAverage).To(Equal(float32(10.0)))
			Expect(got.VoteCount).To(Equal(uint32(1000)))
			Expect(got.FirstAirDate).To(Equal("1.1.2010"))

			Expect(vikings.Name).To(Equal("vikings"))
			Expect(vikings.TmdbId).To(Equal(2))
			Expect(vikings.OriginalName).To(Equal(""))
			Expect(vikings.Description).To(Equal(""))
			Expect(vikings.PosterPath).To(Equal(""))
			Expect(vikings.VoteAverage).To(Equal(float32(0)))
			Expect(vikings.VoteCount).To(Equal(uint32(0)))
			Expect(vikings.FirstAirDate).To(Equal(""))
		})

		It("should save each serie after fetching all seasons and episodes", func() {
			tmdbMock.GetTvEpisodeInfoReturns(&tmdb.TvEpisode{
				AirDate:       "1.1.2010",
				EpisodeNumber: 1,
				Name:          "Episode 1",
				Overview:      "ep1 desc",
				ID:            1000,
				SeasonNumber:  1,
				StillPath:     "stillpath",
				VoteAverage:   1,
				VoteCount:     1,
			}, nil)

			loader.ImportTmdb(db, tmdbMock, series)
			Expect(db.SaveCallCount()).To(Equal(2))

			got := series[0]
			got_s1 := got.Seasons[0]
			got_s1_e1 := got_s1.Episodes[0]

			Expect(got.TmdbId).ToNot(BeZero())
			Expect(got_s1.TmdbId).ToNot(BeZero())
			Expect(got_s1_e1.TmdbId).ToNot(BeZero())
		})

		Context("Seasons", func() {
			BeforeEach(func() {
				tmdbMock.GetTvSeasonInfoStub = func(showid, seasonid int, options map[string]string) (*tmdb.TvSeason, error) {
					result := make(map[int](map[int]*tmdb.TvSeason))
					result_got := make(map[int]*tmdb.TvSeason)
					result_vikings := make(map[int]*tmdb.TvSeason)

					result_got[1] = &tmdb.TvSeason{
						ID:           100,
						Name:         "Season 1",
						AirDate:      "1.1.2010",
						Overview:     "got desc",
						PosterPath:   "got posterpath",
						SeasonNumber: 1,
					}
					result_got[2] = &tmdb.TvSeason{
						ID:           101,
						Name:         "Season 2",
						SeasonNumber: 2,
					}
					result_got[3] = &tmdb.TvSeason{
						ID:           102,
						Name:         "Season 3",
						SeasonNumber: 3,
					}
					result_vikings[1] = &tmdb.TvSeason{
						ID:           110,
						Name:         "Season 1",
						SeasonNumber: 1,
					}
					result_vikings[2] = &tmdb.TvSeason{
						ID:           111,
						Name:         "Season 2",
						SeasonNumber: 2,
					}

					result[1] = result_got
					result[2] = result_vikings
					return result[showid][seasonid], nil
				}
			})

			It("should call TvSeason", func() {
				loader.ImportTmdb(db, tmdbMock, series)

				Expect(tmdbMock.GetTvSeasonInfoCallCount()).To(Equal(4))
				got_s1_show_id, _, _ := tmdbMock.GetTvSeasonInfoArgsForCall(0)
				got_s2_show_id, _, _ := tmdbMock.GetTvSeasonInfoArgsForCall(1)
				vikings_s1_show_id, _, _ := tmdbMock.GetTvSeasonInfoArgsForCall(2)
				vikings_s2_show_id, _, _ := tmdbMock.GetTvSeasonInfoArgsForCall(3)

				Expect(got_s1_show_id).To(Equal(1))
				Expect(got_s2_show_id).To(Equal(1))
				Expect(vikings_s1_show_id).To(Equal(2))
				Expect(vikings_s2_show_id).To(Equal(2))
			})

			It("should call CheckRequest GetTvSeasonInfo", func() {
				loader.ImportTmdb(db, tmdbMock, series)

				Expect(loader.Requested["GetTvSeasonInfo"]).To(Equal(4))
			})

			It("should apply tmdb season ids to seasons", func() {
				loader.ImportTmdb(db, tmdbMock, series)

				_, got_s1_season_id, _ := tmdbMock.GetTvSeasonInfoArgsForCall(0)
				_, got_s2_season_id, _ := tmdbMock.GetTvSeasonInfoArgsForCall(1)
				_, vikings_s1_season_id, _ := tmdbMock.GetTvSeasonInfoArgsForCall(2)
				_, vikings_s2_season_id, _ := tmdbMock.GetTvSeasonInfoArgsForCall(3)

				Expect(got_s1_season_id).To(Equal(1))
				Expect(got_s2_season_id).To(Equal(2))
				Expect(vikings_s1_season_id).To(Equal(1))
				Expect(vikings_s2_season_id).To(Equal(2))
			})

			It("should load tmdb season infos and apply them on seasons", func() {
				loader.ImportTmdb(db, tmdbMock, series)

				got := series[0]
				got_s1 := got.Seasons[0]
				got_s2 := got.Seasons[1]
				vikings := series[1]
				vikings_s1 := vikings.Seasons[0]
				vikings_s2 := vikings.Seasons[1]

				Expect(got_s1.TmdbId).To(Equal(100))
				Expect(got_s2.TmdbId).To(Equal(101))
				Expect(vikings_s1.TmdbId).To(Equal(102))
				Expect(vikings_s2.TmdbId).To(Equal(103))

				Expect(got_s1.AirDate).To(Equal("1.1.2010"))
				Expect(got_s1.OriginalName).To(Equal("Season 1"))
				Expect(got_s1.Name).To(Equal("s1"))
				Expect(got_s1.Description).To(Equal("got desc"))
				Expect(got_s1.PosterPath).To(Equal("got posterpath"))
				Expect(got_s1.SeasonNumber).To(Equal(1))
			})

			It("should load missing seasons", func() {
				series = []*database.Serie{
					&database.Serie{Name: "got", Seasons: []*database.Season{
						{Name: "s2", Episodes: []*database.Episode{{Name: "01"}}},
					}},
					&database.Serie{Name: "vikings", Seasons: []*database.Season{
						{Name: "s1", Episodes: []*database.Episode{{Name: "01"}}},
						{Name: "s2", Episodes: []*database.Episode{{Name: "01"}}},
					}},
				}
				loader.ImportTmdb(db, tmdbMock, series)
				got := series[0].Seasons
				Expect(len(got)).To(Equal(3))
				got_s1 := got[0]
				got_s2 := got[1]
				got_s3 := got[2]
				Expect(got_s1.SeasonNumber).To(Equal(1))
				Expect(got_s2.SeasonNumber).To(Equal(2))
				Expect(got_s3.SeasonNumber).To(Equal(3))

			})

			Context("Episodes", func() {
				It("has episodes", func() {
					loader.ImportTmdb(db, tmdbMock, series)

					got := series[0]
					got_s1 := got.Seasons[0]
					vikings := series[1]
					vikings_s1 := vikings.Seasons[0]

					Expect(len(got_s1.Episodes)).To(Equal(2))
					Expect(len(vikings_s1.Episodes)).To(Equal(1))
				})

				It("should call GetTvEpisodeInfo", func() {
					loader.ImportTmdb(db, tmdbMock, series)

					Expect(tmdbMock.GetTvEpisodeInfoCallCount()).To(Equal(5))
					showid, seasonNum, episodeNum, opt := tmdbMock.GetTvEpisodeInfoArgsForCall(0)
					Expect(showid).To(Equal(1))
					Expect(seasonNum).To(Equal(1))
					Expect(episodeNum).To(Equal(1))
					Expect(len(opt)).To(Equal(0))
				})

				It("should call rlc CheckRequest GetTvEpsiodeInfo", func() {
					loader.ImportTmdb(db, tmdbMock, series)
					Expect(loader.Requested["GetTvEpisodeInfo"]).To(Equal(5))
				})

				It("should apply tmdb episode data to episodes", func() {
					tmdbMock.GetTvEpisodeInfoReturns(&tmdb.TvEpisode{
						AirDate:       "1.1.2010",
						EpisodeNumber: 1,
						Name:          "Episode 1",
						Overview:      "ep1 desc",
						ID:            1000,
						SeasonNumber:  1,
						StillPath:     "stillpath",
						VoteAverage:   1,
						VoteCount:     1,
					}, nil)

					loader.ImportTmdb(db, tmdbMock, series)

					got := series[0]
					got_s1 := got.Seasons[0]
					got_s1_e1 := got_s1.Episodes[0]

					Expect(got_s1_e1.TmdbId).To(Equal(1000))
					Expect(got_s1_e1.Name).To(Equal("01"))
					Expect(got_s1_e1.OriginalName).To(Equal("Episode 1"))
					Expect(got_s1_e1.AirDate).To(Equal("1.1.2010"))
					Expect(got_s1_e1.EpisodeNumber).To(Equal(1))
					Expect(got_s1_e1.StillPath).To(Equal("stillpath"))
					Expect(got_s1_e1.Description).To(Equal("ep1 desc"))
				})
			})
		})
	})

	Describe("UpdateDB", func() {
		It("updated all entries", func() {
			loader.UpdateDB(db, series)
			Expect(db.SaveCallCount()).To(Equal(1))
			Expect(db.SaveArgsForCall(0)).To(Equal(series))
		})
	})
})
