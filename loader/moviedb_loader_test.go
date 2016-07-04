package loader_test

import (
	"errors"
	"reflect"

	"github.com/exane/localflix-server-/database"
	"github.com/exane/localflix-server-/loader"
	"github.com/exane/localflix-server-/loader/loaderfakes"
	"github.com/exane/localflix-server-/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanbradynd05/go-tmdb"
)

var _ = Describe("MoviedbLoader", func() {
	var db *mock.DbMock
	var series []*database.Serie

	BeforeEach(func() {
		db = &mock.DbMock{}
		db.NewRecordCall.Returns = true
		db.NewRecordCall.GotCalled = 0

		series = []*database.Serie{
			&database.Serie{Name: "got"},
			&database.Serie{Name: "vikings"},
		}
	})

	Describe("ImportData", func() {
		It("should not throw an error", func() {
			err := loader.ImportData(db, series)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should call db.NewRecord", func() {
			loader.ImportData(db, series)
			Expect(db.NewRecordCall.GotCalled).ToNot(Equal(0))
		})

		It("should create 2 series", func() {
			loader.ImportData(db, series)
			Expect(db.NewRecordCall.GotCalled).To(Equal(2))
			Expect(len(db.NewRecordCall.Received)).To(Equal(2))
			Expect(reflect.DeepEqual(db.NewRecordCall.Received[0], series[0])).To(Equal(true))
			Expect(reflect.DeepEqual(db.NewRecordCall.Received[1], series[1])).To(Equal(true))
		})
	})

	Describe("ImportTmdb", func() {
		var tmdbMock *loaderfakes.FakeTmdbInterface

		BeforeEach(func() {
			loader.IsTesting = true
			tmdbMock = &loaderfakes.FakeTmdbInterface{}
			tmdbMock.GetTvInfoStub = func(id int, options map[string]string) (*tmdb.TV, error) {
				if loader.Requests() > loader.LIMIT_REQUEST {
					return nil, errors.New("tmdb limit reached")
				}
				var result1 *tmdb.TV
				if id == 1 {
					result1 = &tmdb.TV{
						Overview:     "got desc",
						OriginalName: "got tmdb",
						Name:         "got name",
						ID:           1,
						PosterPath:   "xyz",
						VoteAverage:  10.0,
						VoteCount:    1000,
						FirstAirDate: "1.1.2010",
					}
				}
				if id == 2 {
					result1 = &tmdb.TV{
						Overview:     "",
						OriginalName: "",
						Name:         "",
						ID:           2,
						PosterPath:   "",
						VoteAverage:  0,
						VoteCount:    0,
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

		It("should not panic", func() {
			series = nil
			for i := 0; i < 10; i++ {
				series = append(series, &database.Serie{Name: "got"})
			}
			Expect(func() {
				loader.ImportTmdb(tmdbMock, series)
			}).ToNot(Panic())

			series = nil
			for i := 0; i < 50; i++ {
				series = append(series, &database.Serie{Name: "got"})
			}
			Expect(func() {
				loader.ImportTmdb(tmdbMock, series)
			}).ToNot(Panic())
		})

		It("should fetch entities from tmdb", func() {
			loader.ImportTmdb(tmdbMock, series)

			Expect(tmdbMock.SearchTvCallCount()).To(Equal(2))
			got, _ := tmdbMock.SearchTvArgsForCall(0)
			vikings, _ := tmdbMock.SearchTvArgsForCall(1)
			Expect(got).To(Equal("got"))
			Expect(vikings).To(Equal("vikings"))
		})

		It("should apply on series", func() {
			loader.ImportTmdb(tmdbMock, series)

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
	})
})
