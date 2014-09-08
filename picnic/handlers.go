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

	brotzeitApi := router.PathPrefix("/brotzeit/").Subrouter()
	brotzeitApi.HandleFunc("/start", p.handler(noopHandler, AUTH_LEVEL_IGNORE)).Methods("GET")
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

	// a provisioner can be:
	// ga (Google Analytics), pw (Piwik), sm (URL to sitemap.xml), url (any URL), json (our special JSON format)
	provisionerApi := router.PathPrefix("/provisioners/").Subrouter()
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")        // get account
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("DELETE")     // delete account
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}/save", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("PATCH") // save account data
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}/urls", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")   // retrieve all urls associated

	router.HandleFunc("/", dashBoardHandler).Methods("GET")
	router.HandleFunc("/favicon.ico", handlerFavicon).Methods("GET")

	// due to the rice box regex when building embedded files we must use the full path in the MustFindBox method
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(gzrice.MustFindBox("rd/dist/css").HTTPBox())))
	router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(gzrice.MustFindBox("rd/dist/fonts").HTTPBox())))
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(gzrice.MustFindBox("rd/dist/img").HTTPBox())))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(gzrice.MustFindBox("rd/dist/js").HTTPBox())))
	router.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(gzrice.MustFindBox("rd/dist/lib").HTTPBox())))

	n := negroni.New(
		negroni.HandlerFunc(corsMiddleware),
		negroni.HandlerFunc(GzipContentTypeMiddleware),
	)
	n.UseHandler(router)
	return n
}

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	w.Write(gzrice.MustFindBox("rd/dist/img").MustBytes("favicon.ico"))
}

func dashBoardHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(gzrice.MustFindBox("rd/dist").MustBytes("index.html"))
}

func noopHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {
	return renderString(w, 200, fmt.Sprintf("Found route \n%#v\n %#v\n", r,rc))
}
