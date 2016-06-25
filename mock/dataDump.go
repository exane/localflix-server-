package mock

import (
	"github.com/exane/localflix-server-/database"
)

type DataMock struct {
	ParseCall struct {
		GotCalled bool
		Returns   struct {
			Series []database.Serie
		}
	}
}

func (d *DataMock) Parse(json []byte) []database.Serie {
	d.ParseCall.GotCalled = true
	return d.ParseCall.Returns.Series
}
