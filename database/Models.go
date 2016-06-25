package database

import "github.com/jinzhu/gorm"

type Serie struct {
	gorm.Model
	Name         string `gorm:"unique"`
	Description  string `gorm:"type:text"`
	Seasons      []*Season
	OriginalName string
	PosterPath   string
	VoteAverage  float32
	VoteCount    uint32
	FirstAirDate string
	TmdbId      int
}

type Season struct {
	gorm.Model
	Name         string
	Description  string `gorm:"type:text"`
	PosterPath   string
	AirDate      string
	OriginalName string
	Episodes     []*Episode
	SeasonNumber int
	SerieID      int
	Missing      bool `sql:"DEFAULT:false"`
	TmdbId      int
}

type Episode struct {
	gorm.Model
	Name          string
	Description   string `gorm:"type:text"`
	Src           string
	SeasonID      int
	Extension     string
	Subtitles     []*Subtitle
	Missing       bool `sql:"DEFAULT:false"`
	AirDate       string
	EpisodeNumber int
	OriginalName  string
	StillPath     string
	TmdbId       int
}

type Subtitle struct {
	gorm.Model
	Name      string
	EpisodeID int
}

func (s *Serie) seasons() *Serie {
	season := &Season{}
	serie := &Serie{}
	DB.Find(&season)
	DB.Select("seasons.*").Joins("JOIN seasons ON seasons.serie_id = series.id").Find(&serie)
	return serie
}

func allSeries() []Serie {
	series := []Serie{}
	DB.Find(&series)
	return series
}
