// Copyright (C) Cyrill@Schumacher.fm @SchumacherFM Twitter/GitHub
// Wanderlust - a cache warmer for your web app with priorities
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package brotzeit

import (
	"net/http"
	"testing"

	"github.com/SchumacherFM/wanderlust/github.com/stretchr/testify/assert"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"github.com/SchumacherFM/wanderlust/provisionerApi/provTestHelper"
	"github.com/SchumacherFM/wanderlust/rucksack"
	"github.com/SchumacherFM/wanderlust/rucksack/rsTestHelper"
)

func TestGetCollection(t *testing.T) {

	bp := &rsTestHelper.DbMock{
		CountValue: 29,
	}

	pc := provisionerApi.NewProvisioners()
	p1c := &provTestHelper.ColdCutMock{
		RouteMockFnc: func() string { return "testprov1" },
	}
	p1 := provisionerApi.NewProvisioner("TestProv1", "no-icon", p1c)
	pc.Add(p1)

	c, err := GetCollection(pc, bp)
	assert.NoError(t, err)
	assert.Equal(t, "testprov1", c.Collection[0].Route)
	assert.Equal(t, 29, c.Collection[0].UrlCount)

}

func TestSaveConfigErrorEmptyRoute(t *testing.T) {
	bp := &rsTestHelper.DbMock{
		CountValue: 52,
	}

	body := helpers.NewReadCloser(`{"Route":"","Schedule":"Lord of the tests"}`)
	req, err := http.NewRequest("GET", "/unimportant", body)
	assert.NoError(t, err)
	saveErr := SaveConfig(bp, req)
	assert.Error(t, saveErr)
	assert.EqualError(t, saveErr, ErrCronScheduleEmpty.Error())
}

func TestSaveConfigErrorEmptySchedule(t *testing.T) {
	bp := &rsTestHelper.DbMock{
		CountValue: 65,
	}

	body := helpers.NewReadCloser(`{"Route":"testProv","Schedule":""}`)
	req, err := http.NewRequest("GET", "/unimportant", body)
	assert.NoError(t, err)
	saveErr := SaveConfig(bp, req)
	assert.NoError(t, saveErr)
}

func TestSaveConfigErrorCron(t *testing.T) {
	bp := &rsTestHelper.DbMock{
		CountValue: 67,
	}

	body := helpers.NewReadCloser(`{"Route":"testProv","Schedule":"This is sparta"}`)
	req, err := http.NewRequest("GET", "/unimportant", body)
	assert.NoError(t, err)
	saveErr := SaveConfig(bp, req)
	assert.Error(t, saveErr)
	assert.EqualError(t, saveErr, "Expected 5 or 6 fields, found 3: This is sparta")
}

type dbMock2 struct {
	rsTestHelper.DbMock
	t *testing.T
}

// Composition: we're now "overloading" the Insert method but the parent method is still available thru DbMock.Insert
func (db *dbMock2) Insert(b, k string, d []byte) error {
	assert.Equal(db.t, []byte(`1 3 * * *`), d)
	assert.Equal(db.t, dbCronKeyPrefix+"testProvisioner", k)
	return nil
}

var _ rucksack.Backpacker = &dbMock2{} // check that we implement the interface

func TestSaveConfigSuccess(t *testing.T) {

	bp := &dbMock2{
		DbMock: rsTestHelper.DbMock{
			CountValue: 80,
		},
		t: t,
	}

	body := helpers.NewReadCloser(`{"Route":"testProvisioner","Schedule":"1 3 * * *"}`)
	req, err := http.NewRequest("GET", "/unimportant", body)
	assert.NoError(t, err)
	saveErr := SaveConfig(bp, req)
	assert.Nil(t, saveErr)
}

func TestBootCronNotifier(t *testing.T) {
	t.Skip("@todo implement")
}
