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
	"github.com/SchumacherFM/wanderlust/picnicApi"
	"net/http"
	"runtime"
)

type SystemInfo struct {
	Goroutines int
	// Number of total URLs in the rucksack
	Brotzeit int
	// Number of wanderers running
	Wanderers int
	// Session expires in
	SessionExpires int
	// Which cronjobs are running
	CronJobs []string
}

func (p *PicnicApp) initRoutesSystemInfo(r *mux.Router) error {
	sr := r.PathPrefix("/sysinfo/").Subrouter()
	sr.HandleFunc("/", p.handler(systemInfoHandler, AUTH_LEVEL_LOGIN_WAIT)).Methods("GET")
	return nil
}

func systemInfoHandler(rc picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
	d := newSystemInfo()
	d.SessionExpires = rc.User().GetSessionExpiresIn()
	return helpers.RenderFFJSON(w, d, http.StatusOK)
}

func newSystemInfo() *SystemInfo {
	si := &SystemInfo{
		Goroutines:     runtime.NumGoroutine(),
		Brotzeit:       helpers.RandomInt(6),  // @todo
		Wanderers:      helpers.RandomInt(20), // @todo
		SessionExpires: 0,
		CronJobs:       []string{"Sitemap 23m5s", "G'Analytics 33m10s"},
	}
	return si
}
