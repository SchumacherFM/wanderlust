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
	"net/http"

	"github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/picnicApi"
)

type loginPostData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (p *PicnicApp) initRoutesAuth(r *httprouter.Router) error {

	r.HandlerFunc("GET", "/auth/", p.handler(sessionInfoHandler, AUTH_LEVEL_CHECK))
	r.HandlerFunc("POST", "/auth/", p.handler(loginHandler, AUTH_LEVEL_IGNORE))
	r.HandlerFunc("DELETE", "/auth/", p.handler(logoutHandler, AUTH_LEVEL_LOGIN))
	return nil
}

func sessionInfoHandler(rc picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
	return helpers.RenderFFJSON(w, newSessionInfo(rc.User()), http.StatusOK)
}

func loginHandler(rc picnicApi.Context, w http.ResponseWriter, r *http.Request) error {

	var errLogin = picnicApi.HttpError{
		Status:      http.StatusBadRequest,
		Description: "Invalid username or password",
	}

	lpd := &loginPostData{}

	if err := helpers.DecodeJSON(r, lpd); nil != err {
		return err
	}

	if "" == lpd.UserName || "" == lpd.Password {
		return errLogin
	}

	// find user and login ...
	u := NewUserModel(rc.Backpacker(), lpd.UserName)
	uFound, uErr := u.FindMe()
	if nil != uErr {
		return uErr
	}
	if false == uFound {
		logger.Debug("loginHandler 148: user not found %#v", uFound)
		return errLogin
	}
	if false == u.CheckPassword(lpd.Password) {
		logger.Debug("loginHandler 152: password incorrect %#v", uFound)
		return errLogin
	}

	if err := rc.SessionManager().WriteToken(w, u.GetUserName()); nil != err {
		return err
	}

	u.SetAuthenticated(true)
	// @todo use websocket to send message
	return helpers.RenderFFJSON(w, newSessionInfo(u), http.StatusOK)
}

func logoutHandler(rc picnicApi.Context, w http.ResponseWriter, r *http.Request) error {

	if err := rc.SessionManager().WriteToken(w, ""); err != nil {
		return err
	}
	// @todo use websocket to send message
	return helpers.RenderFFJSON(w, newSessionInfo(nil), http.StatusOK)
}
