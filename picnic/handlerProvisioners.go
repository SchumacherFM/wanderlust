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

	"github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/picnicApi"
	"github.com/SchumacherFM/wanderlust/provisioners"
)

func (p *PicnicApp) initRoutesProvisioners(r *httprouter.Router) error {
	prefix := "/" + provisioners.GetRoutePathPrefix() + "/"
	r.HandlerFunc(
		"GET",
		prefix,
		p.handler(availableProvisionersHandler, AUTH_LEVEL_LOGIN_WAIT),
	)

	pc, err := provisioners.GetAvailable()
	if nil != err {
		return err
	}
	for _, prov := range pc.Collection() {
		route := prefix + prov.Api.Route()
		r.HandlerFunc(
			"GET",
			route,
			p.handler(provisioners.FormGenerate(prov.Api), AUTH_LEVEL_LOGIN),
		)
		r.HandlerFunc(
			"POST",
			route,
			p.handler(provisioners.FormSave(prov.Api), AUTH_LEVEL_LOGIN),
		)
	}
	return nil
}

func availableProvisionersHandler(rc picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
	p, err := provisioners.GetAvailable()
	if nil != err {
		return picnicApi.HttpError{
			Status:      http.StatusNotFound,
			Description: "",
		}
	}
	return helpers.RenderFFJSON(w, p, http.StatusOK)
}
