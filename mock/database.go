package mock

import (
	"github.com/exane/localflix-server-/database"
)

type DbMock struct {
	NewRecordCall struct {
		GotCalled int
		Received  []interface{}
		Returns   bool
		Receives  []database.Serie
	}
	SaveCall struct {
		GotCalled int
		Returns   *DbMock
		Received  []interface{}
	}
}

func (d *DbMock) Save(v interface{}) *DbMock {
	d.SaveCall.GotCalled++
	d.SaveCall.Received = append(d.SaveCall.Received, v)
	return d.SaveCall.Returns
}

func (d *DbMock) NewRecord(v interface{}) bool {
	d.NewRecordCall.GotCalled++
	d.NewRecordCall.Received = append(d.NewRecordCall.Received, v)
	return d.NewRecordCall.Returns
}
