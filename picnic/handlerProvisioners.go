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
	. "github.com/SchumacherFM/wanderlust/picnic/api"
	"github.com/SchumacherFM/wanderlust/provisioners"
	"net/http"
)

// a provisioner can be:
// ga (Google Analytics), pw (Piwik), sm (URL to sitemap.xml), url (any URL), json (our special JSON format)
//provisionerApi := router.PathPrefix("/provisioners/").Subrouter()
//provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")        // get account
//provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("DELETE")     // delete account
//provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}/save", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("PATCH") // save account data
//provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}/urls", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")   // retrieve all urls associated

func (p *PicnicApp) initRoutesProvisioners(r *mux.Router) error {
	sr := r.PathPrefix("/" + provisioners.GetRoutePathPrefix() + "/").Subrouter()
	sr.HandleFunc("/", p.handler(availableProvisionersHandler, AUTH_LEVEL_LOGIN_WAIT)).Methods("GET")

	pc, err := provisioners.GetAvailable()
	if nil != err {
		return err
	}
	for _, prov := range pc.Collection {
		sr.HandleFunc("/"+prov.Api.Route(), p.handler(prov.Api.FormHandler(), AUTH_LEVEL_IGNORE)).Methods("GET")
		sr.HandleFunc("/"+prov.Api.Route()+"/save", p.handler(prov.Api.SaveHandler(), AUTH_LEVEL_LOGIN)).Methods("POST")
		sr.HandleFunc("/"+prov.Api.Route(), p.handler(prov.Api.DeleteHandler(), AUTH_LEVEL_LOGIN)).Methods("DELETE")
	}
	return nil
}

func availableProvisionersHandler(rc RequestContextIf, w http.ResponseWriter, r *http.Request) error {
	p, err := provisioners.GetAvailable()
	if nil != err {
		return httpError{
			Status:      http.StatusNotFound,
			Description: "",
		}
	}
	return helpers.RenderFFJSON(w, p, http.StatusOK)
}
