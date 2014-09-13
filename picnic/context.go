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
	"github.com/SchumacherFM/wanderlust/github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type requestParamsI interface {
	get(name string) string
	getInt(name string) int64
}

// contains route parameters in a map
type requestParams struct {
	vars map[string]string
}

func (r *requestParams) get(name string) string {
	return r.vars[name]
}

func (r *requestParams) getInt(name string) int64 {
	value, _ := strconv.ParseInt(r.vars[name], 10, 0)
	return value
}

type requestContextI interface {
	getApp() PicnicAppI
	getParamString(name string) string
	getParamInt64(name string) int64
	getPicnicer() userIf
}

func newRequestContext(app PicnicAppI, r *http.Request, user userIf) *requestContext {
	ctx := &requestContext{
		app: app,
		par: &requestParams{
			vars: mux.Vars(r),
		},
		pic: user,
	}
	return ctx
}

// request-specific requestContext
// contains the app config so we have access to all the objects we need
type requestContext struct {
	app PicnicAppI
	par requestParamsI
	pic userIf
}

func (rc *requestContext) getApp() PicnicAppI {
	return rc.app
}
func (rc *requestContext) getParamString(name string) string {
	return rc.par.get(name)
}
func (rc *requestContext) getParamInt64(name string) int64 {
	return rc.par.getInt(name)
}
func (rc *requestContext) getPicnicer() userIf {
	return rc.pic
}

//func (ctx *requestContext) validate(v validator, r *http.Request) error {
//	errors := make(map[string]string)
//	if err := v.validate(ctx, r, errors); err != nil {
//		return err
//	}
//	if len(errors) > 0 {
//		return validationFailure{errors}
//	}
//	return nil
//}
