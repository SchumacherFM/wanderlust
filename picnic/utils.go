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
	"encoding/json"
	"fmt"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"github.com/SchumacherFM/wanderlust/helpers"
	"net/http"
	"strconv"
)

func writeBody(w http.ResponseWriter, body []byte, status int, contentType string) error {
	w.Header().Set("Content-Type", contentType+"; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	_, err := w.Write(body)
	return errgo.Mask(err)
}

func renderFFJSON(w http.ResponseWriter, value helpers.FfjsonIf, status int) error {
	body, err := value.MarshalJSON()
	if nil != err {
		return errgo.Mask(err)
	}
	return writeBody(w, body, status, "application/json")
}

func renderJSON(w http.ResponseWriter, value interface{}, status int) error {
	body, err := json.Marshal(value)
	if nil != err {
		return errgo.Mask(err)
	}
	return writeBody(w, body, status, "application/json")
}

func renderString(w http.ResponseWriter, status int, msg string) error {
	return writeBody(w, []byte(msg), status, "text/plain")
}

func getScheme(r *http.Request) string {
	if nil == r.TLS {
		return "http"
	}
	return "https"
}

func getBaseURL(r *http.Request) string {
	return fmt.Sprintf("%s://%s", getScheme(r), r.Host)
}

func decodeJSON(r *http.Request, value interface{}) error {
	return errgo.Mask(json.NewDecoder(r.Body).Decode(value))
}
