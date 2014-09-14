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

	auth := router.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", p.handler(sessionInfoHandler, AUTH_LEVEL_CHECK)).Methods("GET")
	auth.HandleFunc("/", p.handler(loginHandler, AUTH_LEVEL_IGNORE)).Methods("POST")
	auth.HandleFunc("/", p.handler(logoutHandler, AUTH_LEVEL_LOGIN)).Methods("DELETE")
	//	auth.HandleFunc("/signup", p.handler(signup, AUTH_LEVEL_IGNORE)).Methods("POST")
	//	auth.HandleFunc("/recoverpass", p.handler(recoverPassword, AUTH_LEVEL_IGNORE)).Methods("PUT")
	//	auth.HandleFunc("/changepass", p.handler(changePassword, AUTH_LEVEL_IGNORE)).Methods("PUT")

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

	// a provisioner can be:
	// ga (Google Analytics), pw (Piwik), sm (URL to sitemap.xml), url (any URL), json (our special JSON format)
	provisionerApi := router.PathPrefix("/provisioners/").Subrouter()
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")        // get account
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("DELETE")     // delete account
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}/save", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("PATCH") // save account data
	provisionerApi.HandleFunc("/{provisioner}/{id:[0-9]+}/urls", p.handler(noopHandler, AUTH_LEVEL_LOGIN)).Methods("GET")   // retrieve all urls associated

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
	return renderString(w, 200, fmt.Sprintf("Found route \n%#v\n %#v\n", r, rc))
}

func sessionInfoHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {
	return renderJSON(w, newSessionInfo(rc.getUser()), http.StatusOK)
}

func loginHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {

	var invalidLogin = httpError{
		Status:      http.StatusBadRequest,
		Description: "Invalid username or password",
	}

	s := &struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}

	if s.UserName == "" || s.Password == "" {
		return invalidLogin
	}

	// find user and login ...
	user := NewUserModel(s.UserName)
	userFound, userErr := user.findMe()
	if nil != userErr {
		return userErr
	}
	if false == userFound {
		logger.Debug("loginHandler 148: user not found %#v", userFound)
		return invalidLogin
	}
	if false == user.checkPassword(s.Password) {
		logger.Debug("loginHandler 152: password incorrect %#v", userFound)
		return invalidLogin
	}

	if err := rc.getApp().getSessionManager().writeToken(w, user.UserName); err != nil {
		return err
	}

	user.setAuthenticated(true)
	// @todo use websocket to send message
	// @todo still a bug here
	return renderJSON(w, newSessionInfo(user), http.StatusOK)
}

func logoutHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {

	if err := rc.getApp().getSessionManager().writeToken(w, ""); err != nil {
		return err
	}
	// @todo use websocket to send message
	return renderJSON(w, newSessionInfo(nil), http.StatusOK)
}
