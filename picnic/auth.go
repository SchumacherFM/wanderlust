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
	AUTH_LEVEL_IGNORE authLevel = iota // we don't need the user in this handler
	AUTH_LEVEL_CHECK                   // prefetch user, doesn't matter if not logged in
	AUTH_LEVEL_LOGIN                   // user required, 401 if not available
	AUTH_LEVEL_ADMIN                   // admin required, 401 if no user, 403 if not admin
)

// lazily fetches the current session user
func (p *PicnicApp) authenticate(r *http.Request, level authLevel) (*user, error) {

}
