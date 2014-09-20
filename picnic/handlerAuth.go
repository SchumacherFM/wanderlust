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
	"github.com/SchumacherFM/wanderlust/github.com/gorilla/mux"
	"net/http"
)

type loginPostData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (p *PicnicApp) initRoutesAuth(router *mux.Router) error {
	auth := router.PathPrefix("/auth/").Subrouter()
	auth.HandleFunc("/", p.handler(sessionInfoHandler, AUTH_LEVEL_CHECK)).Methods("GET")
	auth.HandleFunc("/", p.handler(loginHandler, AUTH_LEVEL_IGNORE)).Methods("POST")
	auth.HandleFunc("/", p.handler(logoutHandler, AUTH_LEVEL_LOGIN)).Methods("DELETE")
	return nil
}

func sessionInfoHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {
	return renderFFJSON(w, newSessionInfo(rc.getUser()), http.StatusOK)
}

func loginHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {

	var invalidLogin = httpError{
		Status:      http.StatusBadRequest,
		Description: "Invalid username or password",
	}

	lpd := &loginPostData{}

	if err := decodeJSON(r, lpd); nil != err {
		return err
	}

	if "" == lpd.UserName || "" == lpd.Password {
		return invalidLogin
	}

	// find user and login ...
	user := NewUserModel(lpd.UserName)
	userFound, userErr := user.findMe()
	if nil != userErr {
		return userErr
	}
	if false == userFound {
		logger.Debug("loginHandler 148: user not found %#v", userFound)
		return invalidLogin
	}
	if false == user.checkPassword(lpd.Password) {
		logger.Debug("loginHandler 152: password incorrect %#v", userFound)
		return invalidLogin
	}

	if err := rc.getApp().getSessionManager().writeToken(w, user.getUserName()); nil != err {
		return err
	}

	user.setAuthenticated(true)
	// @todo use websocket to send message
	return renderFFJSON(w, newSessionInfo(user), http.StatusOK)
}

func logoutHandler(rc requestContextI, w http.ResponseWriter, r *http.Request) error {

	if err := rc.getApp().getSessionManager().writeToken(w, ""); err != nil {
		return err
	}
	// @todo use websocket to send message
	return renderFFJSON(w, newSessionInfo(nil), http.StatusOK)
}
