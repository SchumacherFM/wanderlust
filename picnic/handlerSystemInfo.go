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
	"github.com/SchumacherFM/wanderlust/github.com/gorilla/mux"
	"github.com/SchumacherFM/wanderlust/helpers"
	"net/http"
	"runtime"
)

type SystemInfo struct {
	Goroutines   int
	Brotzeit     int
	Wanderers    int
	Provisioners int
}

func (p *PicnicApp) initRoutesSystemInfo(router *mux.Router) error {
	sysinfoRoute := router.PathPrefix("/sysinfo/").Subrouter()
	sysinfoRoute.HandleFunc("/", p.handler(systemInfoHandler, AUTH_LEVEL_LOGIN_WAIT)).Methods("GET")

	// @todo
	sysinfoRoute.HandleFunc("/provisioners", p.handler(systemInfoHandler, AUTH_LEVEL_IGNORE)).Methods("GET")
	return nil
}

func systemInfoHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {
	data := newSystemInfo()
	return renderFFJSON(w, data, http.StatusOK)
}

func newSystemInfo() *SystemInfo {
	si := &SystemInfo{
		Goroutines:   runtime.NumGoroutine(),
		Brotzeit:     helpers.RandomInt(6),  // @todo
		Wanderers:    helpers.RandomInt(20), // @todo
		Provisioners: 3,
	}
	return si
}
