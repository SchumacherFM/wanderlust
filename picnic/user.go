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
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"github.com/SchumacherFM/wanderlust/helpers"
	"time"
)

const (
	USER_PASSWORD_LENGTH      = 14
	USER_RECOVERY_CODE_LENGTH = 30
	USER_DB_COLLECTION_NAME   = "users"
	USER_ROOT                 = "root@localhost"
)

type permissions struct {
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
}

type userIf interface {
	getId() int
	getEmail() string
	getName() string
	isAuthenticated() bool
	isAdmin() bool
	prepareNew() error
	// validate(ctx *context, r *http.Request, errors map[string]string) error
	generateRecoveryCode() (string, error)
	resetRecoveryCode()
	generatePassword()
	changePassword(password string) error
	encryptPassword() error
	checkPassword(password string) bool
}

//type userModelCollection struct {
//	userModel []userIf
//}

// not sure if the json tag is needed
type userModel struct {
	CreatedAt       time.Time `json:"createdAt"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	IsAdmin         bool      `json:"isAdmin"`
	IsActive        bool      `json:"isActive"`
	RecoveryCode    string    `json:""`
	IsAuthenticated bool      `json:"isAuthenticated"`
}

func (p *userModel) getId() int            { return helpers.StringHash(p.Email) }
func (p *userModel) getEmail() string      { return p.Email }
func (p *userModel) getName() string       { return p.Name }
func (p *userModel) isAuthenticated() bool { return p.IsAuthenticated }
func (p *userModel) isAdmin() bool         { return p.IsAdmin }

// PreInsert hook
func (p *userModel) prepareNew() error {
	p.IsActive = true
	p.CreatedAt = time.Now()
	p.encryptPassword()
	return nil
}

//func (userModel *userModel) validate(ctx *context, r *http.Request, errors map[string]string) error {
//
//	if userModel.Name == "" {
//		errors["name"] = "Name is missing"
//	} else {
//		ok, err := ctx.datamapper.isUserNameAvailable(userModel)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["name"] = "Name already taken"
//		}
//	}
//
//	if userModel.Email == "" {
//		errors["email"] = "Email is missing"
//	} else if !validateEmail(userModel.Email) {
//		errors["email"] = "Invalid email address"
//	} else {
//		ok, err := ctx.datamapper.isUserEmailAvailable(userModel)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["email"] = "Email already taken"
//		}
//
//	}
//
//	// tbd: we need flag userModel is third-party
//	if userModel.Password == "" {
//		errors["password"] = "Password is missing"
//	}
//
//	return nil
//}

func (p *userModel) generateRecoveryCode() (string, error) {
	code := helpers.RandomString(USER_RECOVERY_CODE_LENGTH)
	p.RecoveryCode = code
	return code, nil
}

func (p *userModel) resetRecoveryCode() {
	p.RecoveryCode = ""
}

// generates an unencrypted password
func (p *userModel) generatePassword() error {
	var err error
	p.Password, err = helpers.NewPassword(USER_PASSWORD_LENGTH)
	return err
}

func (p *userModel) changePassword(password string) error {
	p.Password = password
	return p.encryptPassword()
}

func (p *userModel) encryptPassword() error {
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

func (p *userModel) checkPassword(password string) bool {
	if "" == p.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password))
	return err == nil
}

func (p *userModel) toStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"CreatedAt":       p.CreatedAt.Unix(),
		"Name":            p.Name,
		"Email":           p.Email,
		"Password":        p.Password,
		"IsAdmin":         p.IsAdmin,
		"IsActive":        p.IsActive,
		"IsAuthenticated": p.IsAuthenticated,
	}
}

func StringInterfaceToUser(data map[string]interface{}) userIf {
	return nil
}

// initUsers() runs in NewPicnicApp() function
func initUsers() error {
	var err error
	var rootUser map[string]interface{}
	var password string
	err = rsdb.CreateDatabase(USER_DB_COLLECTION_NAME)
	if nil != err {
		return errgo.Mask(err)
	}

	pn := userModel{
		Name:     "Default Root User",
		Email:    USER_ROOT,
		Password: password,
		IsAdmin:  true,
		IsActive: true,
	}
	pn.generatePassword()
	rootUser, _ = rsdb.FindOne(USER_DB_COLLECTION_NAME, pn.getId())

	if nil == rootUser {
		logger.Printf("Created new user %s with password: %s", pn.Email, pn.Password)
		pn.prepareNew()
		rsdb.InsertRecovery(USER_DB_COLLECTION_NAME, pn.getId(), pn.toStringInterface())
	} else {
		logger.Printf("Root user %s already exists! I don't know the password :-)", USER_ROOT)
	}

	return err
}
