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
	"fmt"
	"net/http"

	"github.com/juju/errgo"
	"picnicApi"
)

// final method in the handler chain
func (p *PicnicApp) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if err, ok := err.(*picnicApi.HttpError); ok {
		http.Error(w, err.Error(), err.Status)
		return
	}
	if err, ok := err.(picnicApi.HttpError); ok {
		http.Error(w, err.Error(), err.Status)
		return
	}

	//	e.g. any other error
	//	if err, ok := err.(validationFailure); ok {
	//		renderJSON(w, err, http.StatusBadRequest)
	//		return
	//	}

	s := fmt.Sprintf("Error:%s", err)
	if err, ok := err.(errgo.Locationer); ok { // type casting
		s += fmt.Sprintf(" %s", err.Location())
	}
	logger.Debug(s)

	http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)
}
