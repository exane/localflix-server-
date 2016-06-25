package loader_test

import (
	"reflect"

	"github.com/exane/localflix-server-/database"
	"github.com/exane/localflix-server-/loader"
	"github.com/exane/localflix-server-/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MoviedbLoader", func() {
	var db *mock.DbMock
	var series []database.Serie

	BeforeEach(func() {
		db = &mock.DbMock{}
		db.NewRecordCall.Returns = true
		db.NewRecordCall.GotCalled = 0

		series = []database.Serie{
			database.Serie{Name: "got"},
			database.Serie{Name: "vikings"},
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
		var tmdbMock *mock.TmdbMock

		BeforeEach(func() {
			tmdbMock = &mock.TmdbMock{}
			tmdbMock.SearchTvCall.GotCalled = 0
			tmdbMock.GetTvInfoCall.GotCalled = 0

			tmdbMock.GetTvInfoCall.Returns.TV = make(map[int]interface{})
			tmdbMock.GetTvInfoCall.Returns.TV[1] = &mock.TV{
				Overview:     "got desc",
				OriginalName: "got tmdb",
				Name:         "got name",
				ID:           1,
				PosterPath:   "xyz",
				VoteAverage:  10.0,
				VoteCount:    1000,
				FirstAirDate: "1.1.2010",
			}

			tmdbMock.GetTvInfoCall.Returns.TV[2] = &mock.TV{
				Overview:     "vikings desc",
				OriginalName: "vikings tmdb",
				Name:         "vikings name",
				ID:           2,
				PosterPath:   "abc",
				VoteAverage:  9.5,
				VoteCount:    2000,
				FirstAirDate: "1.1.2011",
			}

			tmdbMock.SearchTvCall.Returns.TvSearchResults = make(map[string]interface{})
			tmdbMock.SearchTvCall.Returns.TvSearchResults["got"] = &mock.TvSearchResults{
				Results: []mock.SearchTVCallResult{
					{
						ID: 1,
					},
				},
			}
			tmdbMock.SearchTvCall.Returns.TvSearchResults["vikings"] = &mock.TvSearchResults{
				Results: []mock.SearchTVCallResult{
					{
						ID: 2,
					},
				},
			}
		})

		It("should fetch entities from tmdb", func() {
			loader.ImportTmdb(tmdbMock, series)
			Expect(tmdbMock.SearchTvCall.GotCalled).To(Equal(2))
			Expect(tmdbMock.SearchTvCall.Received[0].Name).To(Equal("got"))
			Expect(tmdbMock.SearchTvCall.Received[1].Name).To(Equal("vikings"))
		})

		It("should apply on series", func() {
			loader.ImportTmdb(tmdbMock, series)
			Expect(series[0].Name).To(Equal("got"))
			Expect(series[0].TmdbId).To(Equal(1))
			Expect(series[0].OriginalName).To(Equal("got tmdb"))
			Expect(series[0].Description).To(Equal("got desc"))
			Expect(series[0].PosterPath).To(Equal("xyz"))
			Expect(series[0].VoteAverage).To(Equal(10.0))
			Expect(series[0].VoteCount).To(Equal(1000))
			Expect(series[0].FirstAirDate).To(Equal("1.1.2010"))

			Expect(series[1].Name).To(Equal("vikings"))
			Expect(series[1].TmdbId).To(Equal(2))
			Expect(series[1].OriginalName).To(Equal("vikings tmdb"))
			Expect(series[1].Description).To(Equal("vikings desc"))
			Expect(series[1].PosterPath).To(Equal("abc"))
			Expect(series[1].VoteAverage).To(Equal(9.5))
			Expect(series[1].VoteCount).To(Equal(2000))
			Expect(series[1].FirstAirDate).To(Equal("1.1.2011"))
		})
	})
})
