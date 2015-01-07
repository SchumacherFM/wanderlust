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
	"runtime"
	"testing"

	"github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func compareSystemInfoJsonString(t *testing.T, json string) {
	expected := [4]string{"Brotzeit", "Goroutines", "Wanderers", "SessionExpires"}
	for _, e := range expected {
		assert.Contains(t, json, e)
	}
}

func TestSystemInfoHandler(t *testing.T) {
	rc := &testRequestContext{}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost/systeminfo", nil)
	assert.NoError(t, err)
	actualErr := systemInfoHandler(rc, w, req)
	assert.NoError(t, actualErr)
	compareSystemInfoJsonString(t, w.Body.String())
}

func TestNewSystemInfo(t *testing.T) {

	si := newSystemInfo()
	assert.Exactly(t, si.Goroutines, runtime.NumGoroutine()) // ;-)
	assert.Exactly(t, 0, si.SessionExpires)

	sij, err := si.MarshalJSON()
	assert.NoError(t, err)
	compareSystemInfoJsonString(t, string(sij))
}
