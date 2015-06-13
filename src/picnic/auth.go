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

	"github.com/SchumacherFM/wanderlust/picnicApi"
)

type authLevel int

const (
	AUTH_LEVEL_IGNORE     authLevel = iota + 1 // we don't need the user in this handler
	AUTH_LEVEL_CHECK                           // prefetch user, doesn't matter if not logged in
	AUTH_LEVEL_LOGIN_WAIT                      // user required, 412 if not available
	AUTH_LEVEL_LOGIN                           // user required, 401 if not available
	AUTH_LEVEL_ADMIN                           // admin required, 401 if no user, 403 if not admin
)

func checkAuthLevel(l authLevel, u picnicApi.UserSessionIf) error {
	var (
		errLoginRequired = picnicApi.HttpError{
			Status:      http.StatusUnauthorized,
			Description: "You must be logged in!",
		}
		errWaitingForLogin = picnicApi.HttpError{
			Status:      http.StatusPreconditionFailed,
			Description: "Waiting for login ...",
		}
		falseUser = nil == u || false == u.IsLoggedIn()
	)

	switch l {
	case AUTH_LEVEL_LOGIN_WAIT:
		if falseUser {
			logger.Debug("checkAuthLevel 46: user %#v", u)
			return errWaitingForLogin
		}
		break
	case AUTH_LEVEL_LOGIN:
		if falseUser {
			logger.Debug("checkAuthLevel 52: user %#v", u)
			return errLoginRequired
		}
		break
	case AUTH_LEVEL_ADMIN:
		if falseUser {
			logger.Debug("checkAuthLevel 59: user %#v", u)
			return errLoginRequired
		}
		if false == u.IsAdministrator() {
			return picnicApi.HttpError{
				Status:      http.StatusForbidden,
				Description: "You must be an admin!",
			}
		}
	}
	return nil
}

// lazily fetches the current session user
// check also JWT
func (p *PicnicApp) authenticate(r *http.Request, l authLevel) (*userModel, error) {

	if l == AUTH_LEVEL_IGNORE {
		return nil, nil
	}

	uid, expiresIn, err := p.session.ReadToken(r)
	if err != nil {
		return nil, err
	}

	u := NewUserModel(p.backpacker, uid)
	if "" == uid {
		logger.Debug("p.authenticate: userID from token is empty")
		return u, checkAuthLevel(l, nil)
	}
	var f bool
	f, err = u.FindMe()
	if false == f || err != nil {
		logger.Debug("p.authenticate: user not found in DB %#v", u)
		return nil, checkAuthLevel(l, nil)
	}
	u.SetAuthenticated(true)
	u.SetSessionExpiresIn(expiresIn)
	return u, checkAuthLevel(l, u)

}
