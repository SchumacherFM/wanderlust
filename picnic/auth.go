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

import "net/http"

type authLevel int

const (
	AUTH_LEVEL_IGNORE authLevel = iota + 1 // we don't need the user in this handler
	AUTH_LEVEL_CHECK                       // prefetch user, doesn't matter if not logged in
	AUTH_LEVEL_LOGIN                       // user required, 401 if not available
	AUTH_LEVEL_ADMIN                       // admin required, 401 if no user, 403 if not admin
)

func checkAuthLevel(level authLevel, user userIf) error {
	var errLoginRequired = httpError{
		Status:      http.StatusUnauthorized,
		Description: "You must be logged in!",
	}

	switch level {
	case AUTH_LEVEL_LOGIN:
		if nil == user || false == user.isAuthenticated() {
			logger.Debug("L46: user is %#v", user)
			return errLoginRequired
		}
		break
	case AUTH_LEVEL_ADMIN:
		if nil == user || false == user.isAuthenticated() {
			logger.Debug("L52: user %#v", user)
			return errLoginRequired
		}
		if false == user.isAdmin() {
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
func (p *PicnicApp) authenticate(r *http.Request, level authLevel) (userIf, error) {

	if level == AUTH_LEVEL_IGNORE {
		return nil, nil
	}

	userID, err := p.session.readToken(r)
	if err != nil {
		return nil, err
	}
	user := NewUserModel(userID)
	if "" == userID {
		return nil, checkAuthLevel(level, nil)
	}
	var found bool
	found, err = user.findMe()
	if false == found || err != nil {
		return nil, checkAuthLevel(level, nil)
	}
	user.setAuthenticated(true)

	return user, checkAuthLevel(level, user)

}
