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
	"strings"
	"testing"
)

func compareSystemInfoJsonString(t *testing.T, json string) {
	expected := [4]string{"Brotzeit", "Goroutines", "Wanderers", "SessionExpires"}
	for _, e := range expected {
		if false == strings.Contains(json, e) {
			t.Errorf("\nExpected: %s in Actual:   %s\n", e, json)
		}
	}
}

func TestSystemInfoHandler(t *testing.T) {
	rc := &testRequestContext{}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost/systeminfo", nil)
	if nil != err {
		t.Error(err)
	}
	actualErr := systemInfoHandler(rc, w, req)
	if nil != actualErr {
		t.Error(actualErr)
	}
	compareSystemInfoJsonString(t, w.Body.String())
}

func TestNewSystemInfo(t *testing.T) {

	si := newSystemInfo()
	if si.Goroutines != runtime.NumGoroutine() {
		t.Errorf("Number of Goroutines changed ;-) from %d to %d", si.Goroutines, runtime.NumGoroutine())
	}
	if 0 != si.SessionExpires {
		t.Errorf("SessionExpires must be zero but it is %d", si.SessionExpires)
	}
	sij, err := si.MarshalJSON()
	if nil != err {
		t.Error(err)
	}
	compareSystemInfoJsonString(t, string(sij))
}
