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

// @todo add kr/secureheader

import (
	"net/http"
)

func corsMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Return CORS Headers and end
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "X-AUTH-TOKEN, X-API-VERSION, X-Requested-With, Content-Type, Accept, Origin")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	res.Header().Set("Access-Control-Max-Age", "1728000")

	if req.Method == "OPTIONS" {
		res.WriteHeader(200)
		return
	}

	next(res, req)
}
