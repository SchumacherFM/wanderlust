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
	USER_PASSWORD_LENGTH      = 24
	USER_RECOVERY_CODE_LENGTH = 30
	USER_DB_COLLECTION_NAME   = "users"
	USER_ROOT                 = "adiministrator"
)

type permissions struct {
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
}

type userIf interface {
	getId() int
	getEmail() string
	getName() string
	getUserName() string
	isAuthenticated() bool
	setAuthenticated(bool)
	isAdmin() bool
	prepareNew() error
	isValidForSession() bool
	// validate(ctx *context, r *http.Request, errors map[string]string) error
	generateRecoveryCode() (string, error)
	resetRecoveryCode()
	generatePassword() error
	changePassword(string) error
	encryptPassword() error
	checkPassword(string) bool
	findMe() (bool, error)
	applyDbData(map[string]interface{})
}

type userGetterSetter interface {
	getId() int
	getEmail() string
	getName() string
	getUserName() string
	isAuthenticated() bool
	setAuthenticated(bool)
}

//type userModelCollection struct {
//	userModel []userIf
//}

type userModel struct {
	CreatedAt       time.Time
	UserName        string
	Name            string
	Email           string
	Password        string
	IsAdmin         bool
	IsActive        bool
	RecoveryCode    string
	IsAuthenticated bool
}

func (um *userModel) getId() int                 { return helpers.StringHash(um.UserName) }
func (um *userModel) getEmail() string           { return um.Email }
func (um *userModel) getUserName() string        { return um.UserName }
func (um *userModel) getName() string            { return um.Name }
func (um *userModel) isAuthenticated() bool      { return um.IsAuthenticated }
func (um *userModel) setAuthenticated(auth bool) { um.IsAuthenticated = auth }
func (um *userModel) isAdmin() bool              { return um.IsAdmin }

// PreInsert hook for new users
func (um *userModel) prepareNew() error {
	um.IsActive = true
	um.CreatedAt = time.Now()
	return um.encryptPassword()
}

// isValidForSession() is only used in newSessionInfo()
func (um *userModel) isValidForSession() bool {
	return true == helpers.ValidateEmail(um.getEmail()) && "" != um.getUserName() && true == um.isAuthenticated()
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

func (um *userModel) generateRecoveryCode() (string, error) {
	code := helpers.RandomString(USER_RECOVERY_CODE_LENGTH)
	um.RecoveryCode = code
	return code, nil
}

func (um *userModel) resetRecoveryCode() {
	um.RecoveryCode = ""
}

// generates an unencrypted password
func (um *userModel) generatePassword() error {
	var err error
	um.Password, err = helpers.NewPassword(USER_PASSWORD_LENGTH)
	return err
}

func (um *userModel) changePassword(password string) error {
	um.Password = password
	return um.encryptPassword()
}

func (um *userModel) encryptPassword() error {
	if "" == um.Password {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(um.Password), bcrypt.DefaultCost)
	if nil != err {
		return err
	}
	um.Password = string(hashed)
	return nil
}

// not sure if it is a good idea to carry the whole time the bcrypted password with the userModel object ...
func (um *userModel) checkPassword(password string) bool {
	if "" == um.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(um.Password), []byte(password))
	return err == nil
}

func (um *userModel) toStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"CreatedAt":       um.CreatedAt.Unix(),
		"UserName":        um.UserName,
		"Name":            um.Name,
		"Email":           um.Email,
		"Password":        um.Password,
		"IsAdmin":         um.IsAdmin,
		"IsActive":        um.IsActive,
		"IsAuthenticated": um.IsAuthenticated,
	}
}

// finds a user in the database and fills the struct
func (um *userModel) findMe() (bool, error) {
	searchedUser, _ := rsdb.FindOne(USER_DB_COLLECTION_NAME, um.getId())
	if nil == searchedUser {
		return false, nil
	}
	um.applyDbData(searchedUser)

	return true, nil
}

func (um *userModel) applyDbData(data map[string]interface{}) {
	// panic free type conversion
	tIsAdmin, _ := data["IsAdmin"].(bool)
	tIsActive, _ := data["IsActive"].(bool)
	tIsAuthenticated, _ := data["IsAuthenticated"].(bool)
	tCreatedAt, _ := data["CreatedAt"].(int64)
	tUserName, _ := data["UserName"].(string)
	tName, _ := data["Name"].(string)
	tEmail, _ := data["Email"].(string)
	tPassword, _ := data["Password"].(string)

	um.CreatedAt = time.Unix(tCreatedAt, 0)
	um.UserName = tUserName
	um.Name = tName
	um.Email = tEmail
	um.Password = tPassword
	um.IsAdmin = tIsAdmin
	um.IsActive = tIsActive
	um.IsAuthenticated = tIsAuthenticated
}

// needed in auth when user tries to login
func NewUserModel(userName string) *userModel {
	um := &userModel{
		UserName: userName,
	}
	return um
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

	um := userModel{
		UserName: USER_ROOT,
		Name:     "Default Root User",
		Email:    USER_ROOT + "@localhost",
		Password: password,
		IsAdmin:  true,
		IsActive: true,
	}
	um.generatePassword()
	rootUser, _ = rsdb.FindOne(USER_DB_COLLECTION_NAME, um.getId())

	if nil == rootUser {
		logger.Emergency("Created new user %s with password: %s", um.UserName, um.Password)
		um.prepareNew()
		rsdb.InsertRecovery(USER_DB_COLLECTION_NAME, um.getId(), um.toStringInterface())
	} else {
		logger.Emergency("Root user %s already exists!", USER_ROOT)
	}
	return errgo.Mask(err)
}
