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
	DB_COLLECTION_NAME = "users"
)

type permissions struct {
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
}

type picnicerI interface {
	getNextId() int
	PreparePicnicer() error
	// validate(ctx *context, r *http.Request, errors map[string]string) error
	generateRecoveryCode() (string, error)
	resetRecoveryCode()
	changePassword(password string) error
	encryptPassword() error
	checkPassword(password string) bool
}

type picnicerCollection struct {
	picnicer []picnicerI
}

// Picnicer represents users
type picnicer struct {
	ID              int       `json:"id"`
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

func newPicnicer() *picnicer {
	return &picnicer{}
}

func (p *picnicer) getNextId() int {
	return 1
}

// PreInsert hook
func (p *picnicer) PreparePicnicer() error {
	p.IsActive = true
	p.CreatedAt = time.Now()
	p.encryptPassword()
	return nil
}

//func (picnicer *picnicer) validate(ctx *context, r *http.Request, errors map[string]string) error {
//
//	if picnicer.Name == "" {
//		errors["name"] = "Name is missing"
//	} else {
//		ok, err := ctx.datamapper.isUserNameAvailable(picnicer)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["name"] = "Name already taken"
//		}
//	}
//
//	if picnicer.Email == "" {
//		errors["email"] = "Email is missing"
//	} else if !validateEmail(picnicer.Email) {
//		errors["email"] = "Invalid email address"
//	} else {
//		ok, err := ctx.datamapper.isUserEmailAvailable(picnicer)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["email"] = "Email already taken"
//		}
//
//	}
//
//	// tbd: we need flag picnicer is third-party
//	if picnicer.Password == "" {
//		errors["password"] = "Password is missing"
//	}
//
//	return nil
//}

func (p *picnicer) generateRecoveryCode() (string, error) {
	code := helpers.RandomString(RECOVERY_CODE_LENGTH)
	p.RecoveryCode = code
	return code, nil
}

func (p *picnicer) resetRecoveryCode() {
	p.RecoveryCode = ""
}

func (p *picnicer) changePassword(password string) error {
	p.Password = password
	return p.encryptPassword()
}

func (p *picnicer) encryptPassword() error {
	if "" == p.Password {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if nil != err {
		return err
	}
	p.Password = string(hashed)
	return nil
}

func (p *picnicer) checkPassword(password string) bool {
	if "" == p.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	return err == nil
}
