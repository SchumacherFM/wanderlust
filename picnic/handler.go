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

// Wanderlust uses go.rice package for serving static web content
//

package picnic

import (
	"fmt"
	gzrice "github.com/SchumacherFM/wanderlust/github.com/SchumacherFM/go.gzrice"
	"github.com/SchumacherFM/wanderlust/github.com/codegangsta/negroni"
	"github.com/SchumacherFM/wanderlust/github.com/gorilla/mux"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/picnic/middleware"
	"net/http"
)

// our custom handler
type handlerFunc func(rc requestContextI, w http.ResponseWriter, r *http.Request) error

// the handler should create a new context on each request, and handle any returned
// errors appropriately.
func (p *PicnicApp) handler(h handlerFunc, level authLevel) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		doAuthentication := func() error {
			user, err := p.authenticate(r, level)
			if err != nil {
				return err
			}
			return h(newRequestContext(p, r, user), w, r)
		}
		p.handleError(w, r, doAuthentication())
	}
}

func (p *PicnicApp) getHandler() *negroni.Negroni {
	router := mux.NewRouter()

	p.initRoutesAuth(router)
	p.initRoutesUsers(router)
	p.initRoutesSystemInfo(router)
	p.initRoutesProvisioners(router)

	brotzeitApi := router.PathPrefix("/brotzeit/").Subrouter()
	brotzeitApi.HandleFunc("/start", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")
	brotzeitApi.HandleFunc("/stop", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")
	brotzeitApi.HandleFunc("/purge", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET") // purges all collected URLs
	brotzeitApi.HandleFunc("/concurrency", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("PUT")
	brotzeitApi.HandleFunc("/collections", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET") // retrieves running processes

	wandererApi := router.PathPrefix("/wanderer/").Subrouter()
	wandererApi.HandleFunc("/start", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")
	wandererApi.HandleFunc("/stop", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")
	wandererApi.HandleFunc("/concurrency", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("PUT")
	wandererApi.HandleFunc("/current", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")

	// start stop the database web interface
	rucksackApi := router.PathPrefix("/rucksack/").Subrouter()
	rucksackApi.HandleFunc("/start", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")
	rucksackApi.HandleFunc("/stop", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")

	dashboardApi := router.PathPrefix("/dashboard/").Subrouter()

	// loads automatically the index.html
	dashboardApi.PathPrefix("/").Handler(http.StripPrefix("/dashboard", http.FileServer(gzrice.MustFindBox("rd/dist/").HTTPBox())))
	router.HandleFunc("/", handlerRedirectDashboard)
	router.HandleFunc("/dashboard", handlerRedirectDashboard)
	router.HandleFunc("/favicon.ico", handlerFavicon)

	n := negroni.New(
		negroni.HandlerFunc(middleware.CorsMiddleware),
		negroni.HandlerFunc(middleware.GzipContentTypeMiddleware),
	)
	n.UseHandler(router)
	return n
}

func handlerRedirectDashboard(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard/", 301)
}
func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	w.Write(gzrice.MustFindBox("rd/dist/").MustBytes("img/favicon.ico"))
}

func noopHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {
	return helpers.RenderString(w, 200, fmt.Sprintf("Found route \n%#v\n %#v\n", r, rc))
}