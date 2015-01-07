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

	"github.com/SchumacherFM/wanderlust/brotzeit"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/picnicApi"
	"github.com/SchumacherFM/wanderlust/provisioners"
	"github.com/julienschmidt/httprouter"
)

func (p *PicnicApp) initRoutesBrotzeit(r *httprouter.Router) error {

	r.HandlerFunc("GET", "/brotzeit/", p.handler(brotzeitCollectionHandler, AUTH_LEVEL_LOGIN))
	r.HandlerFunc("POST", "/brotzeit/", p.handler(brotzeitSaveHandler, AUTH_LEVEL_LOGIN))

	//	// @todo
	//	brotzeitApi.HandleFunc("/start", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")
	//	brotzeitApi.HandleFunc("/stop", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")
	//	brotzeitApi.HandleFunc("/purge", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET") // purges all collected URLs
	//	brotzeitApi.HandleFunc("/concurrency", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("PUT")
	//	brotzeitApi.HandleFunc("/collections", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET") // retrieves running processes
	//

	return nil
}

// brotzeitCollectionHandler returns JSON containing all provisioners with name and icon, their URL
// count and cron settings which includes cron schedule and if cron is available
func brotzeitCollectionHandler(rc picnicApi.Context, w http.ResponseWriter, r *http.Request) error {

	pc, _ := provisioners.GetAvailable()
	p, err := brotzeit.GetCollection(pc, rc.Backpacker())
	if nil != err {
		return picnicApi.HttpError{
			Status:      http.StatusNotFound,
			Description: "",
		}
	}
	return helpers.RenderFFJSON(w, p, http.StatusOK)

}

// brotzeitSaveHandler saves the cron config
func brotzeitSaveHandler(rc picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
	err := brotzeit.SaveConfig(rc.Backpacker(), r)
	if nil != err {
		return picnicApi.HttpError{
			Status:      http.StatusBadRequest,
			Description: err.Error(),
		}
	}
	return helpers.RenderString(w, http.StatusOK, "")
}
