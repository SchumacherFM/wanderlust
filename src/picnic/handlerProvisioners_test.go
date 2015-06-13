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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"picnicApi"
	"provisioners"
	"rucksack"
)

type testRequestContext struct {
}

func (rc *testRequestContext) SessionManager() picnicApi.SessionManager { return nil }
func (rc *testRequestContext) Backpacker() rucksack.Backpacker          { return nil }
func (rc *testRequestContext) GetParamString(s string) string           { return s }
func (rc *testRequestContext) GetParamInt64(s string) int64             { return 0 }
func (rc *testRequestContext) User() picnicApi.UserSessionIf {
	um := NewUserModel(nil, "testUser")
	return um
}

func TestAvailableProvisionersHandler(t *testing.T) {

	expected := []byte(`{"Collection":[{"Icon":"fa-sitemap","Name":"Sitemap","Url":"/provisioners/sitemap"},{"Icon":"fa-file-text-o","Name":"Textarea","Url":"/provisioners/textarea"}]}`)

	rc := &testRequestContext{}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost"+provisioners.GetRoutePathPrefix(), nil)
	assert.NoError(t, err)
	actualErr := availableProvisionersHandler(rc, w, req)
	assert.NoError(t, actualErr)
	assert.Exactly(t, w.Body.Bytes(), expected)
}
