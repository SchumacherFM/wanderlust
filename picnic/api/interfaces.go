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

package api

import (
	"github.com/SchumacherFM/wanderlust/helpers"
	"net/http"
	"time"
)

// And here we go ... ;-(
// https://medium.com/@rakyll/interface-pollution-in-go-7d58bccec275
// Interface Pollution in Go

type (
	// our custom handler
	HandlerFunc func(rc RequestContextIf, w http.ResponseWriter, r *http.Request) error

	SessionManagerIf interface {
		ReadToken(*http.Request) (string, time.Duration, error)
		WriteToken(http.ResponseWriter, string) error
	}

	PicnicAppIf interface {
		GetSessionManager() SessionManagerIf
		Execute() error
		GetListenAddress() string
	}

	RequestContextIf interface {
		GetApp() PicnicAppIf
		GetParamString(string) string
		GetParamInt64(string) int64
		GetUser() UserSessionIf
	}

	// UserSessionIf is special interface only used when requesting the session in a handler
	UserSessionIf interface {
		GetEmail() string
		GetName() string
		GetUserName() string
		IsLoggedIn() bool
		IsActive() bool
		IsAdministrator() bool
		IsValidForSession() bool
		SetAuthenticated(bool) error
		SetSessionExpiresIn(time.Duration) error
		GetSessionExpiresIn() int
		CheckPassword(string) bool
	}

	UserGetterIf interface {
		GetId() int
		GetEmail() string
		GetName() string
		GetUserName() string
		GetSessionExpiresIn() int
		ToStringInterface() map[string]interface{}
		FindMe() (bool, error)
		helpers.FfjsonIf
	}

	UserSetterIf interface {
		SetEmail(string) error
		SetName(string) error
		SetUserName(string) error

		// validate(ctx *context, r *http.Request, errors map[string]string) error
		GenerateRecoveryCode() (string, error)
		ResetRecoveryCode()
		GeneratePassword() error
		ChangePassword(string) error
		EncryptPassword() error
		UnsetPassword()
	}
)
