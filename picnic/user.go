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
	"github.com/SchumacherFM/wanderlust/code.google.com/p/go.crypto/bcrypt"
	"github.com/SchumacherFM/wanderlust/helpers"
	"time"
)

const (
	RECOVERY_CODE_LENGTH = 30
)

type permissions struct {
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
}

// User represents users in database
type user struct {
	ID              int     `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	Name            string    `json:"name"`
	Password        string    `json:""`
	Email           string    `json:"email"`
	Votes           string    `json:""`
	IsAdmin         bool      `json:"isAdmin"`
	IsActive        bool      `json:"isActive"`
	RecoveryCode    string    `json:""`
	IsAuthenticated bool      `json:"isAuthenticated"`
}

// PreInsert hook
func (user *user) PrepareNewUser() error {
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.encryptPassword()
	return nil
}

//func (user *user) validate(ctx *context, r *http.Request, errors map[string]string) error {
//
//	if user.Name == "" {
//		errors["name"] = "Name is missing"
//	} else {
//		ok, err := ctx.datamapper.isUserNameAvailable(user)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["name"] = "Name already taken"
//		}
//	}
//
//	if user.Email == "" {
//		errors["email"] = "Email is missing"
//	} else if !validateEmail(user.Email) {
//		errors["email"] = "Invalid email address"
//	} else {
//		ok, err := ctx.datamapper.isUserEmailAvailable(user)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["email"] = "Email already taken"
//		}
//
//	}
//
//	// tbd: we need flag user is third-party
//	if user.Password == "" {
//		errors["password"] = "Password is missing"
//	}
//
//	return nil
//}

func (user *user) generateRecoveryCode() (string, error) {
	code := helpers.RandomString(RECOVERY_CODE_LENGTH)
	user.RecoveryCode = code
	return code, nil
}

func (user *user) resetRecoveryCode() {
	user.RecoveryCode = ""
}

func (user *user) changePassword(password string) error {
	user.Password = password
	return user.encryptPassword()
}

func (user *user) encryptPassword() error {
	if "" == user.Password {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if nil != err {
		return err
	}
	user.Password = string(hashed)
	return nil
}

func (user *user) checkPassword(password string) bool {
	if "" == user.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
