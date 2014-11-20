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
	"github.com/SchumacherFM/wanderlust/picnicApi"
	"github.com/SchumacherFM/wanderlust/rucksack"
	"net/http"
	"strconv"
)

type (

	// request-specific requestContext
	// contains the app config so we have access to all the objects we need
	requestContext struct {
		sm picnicApi.SessionManager
		//		vars map[string]string
		user picnicApi.UserSessionIf
		bp   rucksack.Backpacker
	}
)

//SessionManager() SessionManager
//GetParamString(string) string
//GetParamInt64(string) int64
//User() UserSessionIf
//Backpacker() rucksack.Backpacker

// invoked in (p *PicnicApp) handler()
// per request on context
func newRequestContext(s picnicApi.SessionManager, r *http.Request, u picnicApi.UserSessionIf, b rucksack.Backpacker) *requestContext {
	ctx := &requestContext{
		sm: s,
		//		vars: mux.Vars(r),
		user: u,
		bp:   b,
	}
	return ctx
}

func (rc *requestContext) SessionManager() picnicApi.SessionManager {
	return rc.sm
}

func (rc *requestContext) Backpacker() rucksack.Backpacker {
	return rc.bp
}

func (rc *requestContext) GetParamString(name string) string {
	if val, ok := rc.vars[name]; ok {
		return val
	}
	logger.Debug("%s not found request vars %#v", name, rc.vars)
	return ""
}

func (rc *requestContext) GetParamInt64(name string) int64 {
	v, e := strconv.ParseInt(rc.vars[name], 10, 0)
	if nil == e {
		return v
	}
	logger.Debug("%s not found request vars %#v", name, rc.vars)
	return 0
}

func (rc *requestContext) User() picnicApi.UserSessionIf {
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
