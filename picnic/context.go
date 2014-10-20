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
	. "github.com/SchumacherFM/wanderlust/picnic/api"
	"net/http"
	"strconv"
)

type (

	// request-specific requestContext
	// contains the app config so we have access to all the objects we need
	requestContext struct {
		app  PicnicAppI
		vars map[string]string
		user UserGetPermIf
	}
)

// invoked in (p *PicnicApp) handler()
// per request on context
func newRequestContext(app PicnicAppI, r *http.Request, theUser UserGetPermIf) *requestContext {
	ctx := &requestContext{
		app: app,
		vars: mux.Vars(r),
		user: theUser,
	}
	return ctx
}

func (rc *requestContext) GetApp() PicnicAppI {
	return rc.app
}

func (rc *requestContext) GetParamString(name string) string {
	if val, ok := rc.vars[name]; ok {
		return val
	}
	logger.Debug("%s not found request vars %#v", name, rc.vars)
	return ""
}

func (rc *requestContext) GetParamInt64(name string) int64 {
	v,e := strconv.ParseInt(rc.vars[name], 10, 0);
	if nil == e {
		return v
	}
	logger.Debug("%s not found request vars %#v", name, rc.vars)
	return 0
}

func (rc *requestContext) GetUser() UserGetPermIf {
	return rc.user
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
