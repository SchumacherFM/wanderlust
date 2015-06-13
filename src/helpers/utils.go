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

package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/juju/errgo"
)

func WriteBody(w http.ResponseWriter, body []byte, status int, contentType string) error {
	w.Header().Set("Content-Type", contentType+"; charset=UTF8")
	//w.Header().Set("Content-Length", strconv.Itoa(len(body))) Length will be set via GZIP middleware
	w.WriteHeader(status)
	_, err := w.Write(body)
	return errgo.Mask(err)
}

func RenderFFJSON(w http.ResponseWriter, value FfjsonIf, status int) error {
	body, err := value.MarshalJSON()
	if nil != err {
		return errgo.Mask(err)
	}
	return WriteBody(w, body, status, "application/json")
}

func RenderJSON(w http.ResponseWriter, value interface{}, status int) error {
	body, err := json.Marshal(value)
	if nil != err {
		return errgo.Mask(err)
	}
	return WriteBody(w, body, status, "application/json")
}

func RenderString(w http.ResponseWriter, status int, msg string) error {
	return WriteBody(w, []byte(msg), status, "text/plain")
}

func RenderHTML(w http.ResponseWriter, status int, msg string) error {
	return WriteBody(w, []byte(msg), status, "text/html")
}

func DecodeJSON(r *http.Request, value interface{}) error {
	return errgo.Mask(json.NewDecoder(r.Body).Decode(value))
}
