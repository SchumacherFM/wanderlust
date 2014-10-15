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
)

type authLevel int

const (
	AUTH_LEVEL_IGNORE     authLevel = iota + 1 // we don't need the user in this handler
	AUTH_LEVEL_CHECK                           // prefetch user, doesn't matter if not logged in
	AUTH_LEVEL_LOGIN_WAIT                      // user required, 412 if not available
	AUTH_LEVEL_LOGIN                           // user required, 401 if not available
	AUTH_LEVEL_ADMIN                           // admin required, 401 if no user, 403 if not admin
)

func checkAuthLevel(l authLevel, u userPermissionsIf) error {
	var (
		errLoginRequired = httpError{
			Status:      http.StatusUnauthorized,
			Description: "You must be logged in!",
		}
		errWaitingForLogin = httpError{
			Status:      http.StatusPreconditionFailed,
			Description: "Waiting for login ...",
		}
		falseUser = nil == u || false == u.isAuthenticated()
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
		if false == u.isAdmin() {
			return httpError{
				Status:      http.StatusForbidden,
				Description: "You must be an admin!",
			}
		}
	}
	return nil
}

// lazily fetches the current session user
// check also JWT
func (p *PicnicApp) authenticate(r *http.Request, l authLevel) (userIf, error) {

	if l == AUTH_LEVEL_IGNORE {
		return nil, nil
	}

	uid, expiresIn, err := p.session.readToken(r)
	if err != nil {
		return nil, err
	}
	u := NewUserModel(uid)
	if "" == uid {
		logger.Debug("p.authenticate: userID from token is empty")
		return nil, checkAuthLevel(l, nil)
	}
	var f bool
	f, err = u.findMe()
	if false == f || err != nil {
		logger.Debug("p.authenticate: user not found in DB %#v", u)
		return nil, checkAuthLevel(l, nil)
	}
	u.setAuthenticated(true)
	u.setSessionExpiresIn(expiresIn)
	return u, checkAuthLevel(l, u)

}
