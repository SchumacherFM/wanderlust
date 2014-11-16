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

package picnic

import (
	"github.com/SchumacherFM/wanderlust/github.com/stretchr/testify/assert"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/rucksack/rsTestHelper"
	"testing"
)

var (
	fixture0 = []byte(`{"CreatedAt":"2014-10-26T11:11:03.753110065+11:00","Email":"test@tester.com","IsActivated":false,"IsAdmin":false,"IsAuthenticated":false,"Name":"Test User","Password":"p/aH4*VfiXka7{3sB1PGOQ!p","RecoveryCode":"","UserName":"testuser"}`)
	// password had been removed because FindAll removes PW for security reasons
	fixture1 = []byte(`{"CreatedAt":"2014-10-26T11:11:03.753110065+11:00","Email":"test@tester.com","IsActivated":false,"IsAdmin":false,"IsAuthenticated":false,"Name":"Test User","Password":"","RecoveryCode":"","UserName":"testuser"}`)
	fixture2 = []byte(`{"CreatedAt":"2014-10-26T11:11:03.753110066+11:00","Email":"test@tester2.com","IsActivated":false,"IsAdmin":false,"IsAuthenticated":false,"Name":"Test User2","Password":"","RecoveryCode":"","UserName":"testuser2"}`)
	db       = &rsTestHelper.DbMock{
		FindOneData: fixture0,
		FindAllData: [][]byte{
			[]byte(helpers.StringHashString("testuser")),
			fixture1,
			[]byte(helpers.StringHashString("testuser2")),
			fixture2,
		},
	}
)

func TestFindMe(t *testing.T) {

	u := NewUserModel(db, "testuser")
	found, err := u.FindMe()
	if false == found || nil != err {
		t.Error(found, err)
	}

	actual, _ := u.MarshalJSON()
	assert.Exactly(t, fixture0, actual)
}

func TestFindAllUsers(t *testing.T) {

	uc := NewUserModelCollection(db)
	err := uc.FindAllUsers()
	if nil != err {
		t.Error(err)
	}
	a1, _ := uc.Users[0].MarshalJSON()
	assert.Exactly(t, fixture1, a1)
	a2, _ := uc.Users[1].MarshalJSON()
	assert.Exactly(t, fixture2, a2)
}
