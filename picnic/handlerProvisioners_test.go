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
	"bytes"
	. "github.com/SchumacherFM/wanderlust/picnic/api"
	"github.com/SchumacherFM/wanderlust/provisioners"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRequestContext struct {
}

func (rc *testRequestContext) App() PicnicAppIf               { return nil }
func (rc *testRequestContext) GetParamString(s string) string { return s }
func (rc *testRequestContext) GetParamInt64(s string) int64   { return 0 }
func (rc *testRequestContext) User() UserSessionIf {
	um := NewUserModel(nil, "testUser")
	return um
}

func TestAvailableProvisionersHandler(t *testing.T) {

	expected := []byte(`{"Collection":[{"Icon":"fa-sitemap","Name":"Sitemap","Url":"/provisioners/sitemap"},{"Icon":"fa-file-text-o","Name":"Textarea","Url":"/provisioners/textarea"}]}`)

	rc := &testRequestContext{}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost"+provisioners.GetRoutePathPrefix(), nil)
	if nil != err {
		t.Error(err)
	}
	actualErr := availableProvisionersHandler(rc, w, req)
	if nil != actualErr {
		t.Error(actualErr)
	}
	if bytes.Compare(w.Body.Bytes(), expected) != 0 {
		t.Errorf("\nExpected: %s\nActual:   %s\n", expected, w.Body.Bytes())
	}
}
