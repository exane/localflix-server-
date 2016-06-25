package mock

import "github.com/ryanbradynd05/go-tmdb"

type TmdbMock struct {
	SearchTvCall struct {
		GotCalled int
		Returns   struct {
			TvSearchResults map[string]interface{}
			Error           map[string]error
		}
		Received []searchTvCallReceived
	}
	GetTvInfoCall struct {
		GotCalled int
		Returns   struct {
			TV    map[int]interface{}
			Error map[int]error
		}
		Received []getTvInfoReceived
	}
}

func (t *TmdbMock) GetTvInfo(id int, options map[string]string) (*tmdb.TV, error) {
	t.GetTvInfoCall.GotCalled++
	t.GetTvInfoCall.Received = append(t.GetTvInfoCall.Received,
		getTvInfoReceived{Id: id, Options: options})
	return t.GetTvInfoCall.Returns.TV[id].(*tmdb.TV), t.GetTvInfoCall.Returns.Error[id]
}

type TV struct {
	Name         string
	Overview     string
	OriginalName string
	PosterPath   string
	VoteAverage  float32
	VoteCount    uint32
	FirstAirDate string
	ID           int
}

type getTvInfoReceived struct {
	Id      int
	Options map[string]string
}

type searchTvCallReceived struct {
	Name    string
	Options map[string]string
}

type SearchTVCallResult struct {
	BackdropPath  string
	ID            int
	OriginalName  string
	FirstAirDate  string
	OriginCountry []string
	PosterPath    string
	Popularity    float32
	Name          string
	VoteAverage   float32
	VoteCount     uint32
}

type TvSearchResults struct {
	Results []SearchTVCallResult
}

func (t *TmdbMock) SearchTv(name string, options map[string]string) (*tmdb.TvSearchResults, error) {
	t.SearchTvCall.GotCalled++
	t.SearchTvCall.Received = append(t.SearchTvCall.Received, searchTvCallReceived{Name: name, Options: options})
	result := t.SearchTvCall.Returns.TvSearchResults[name]
	//PANIC interface conversion: interface {} is *mock.TvSearchResults, not *tmdb.TvSearchResults
	return result.(*tmdb.TvSearchResults), t.SearchTvCall.Returns.Error[name]
}
