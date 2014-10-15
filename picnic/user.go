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
	USER_ROOT                 = "administrator"
)

// now that's a hell lot of interfaces ...

type userIf interface {
	userGetterIf
	userSetterIf
	userPermissionsIf
}

// userGetPermIf is for GetterPermissions Interface
type userGetPermIf interface {
	userGetterIf
	userPermissionsIf
}

// userSessionIf is special interface only used when requesting the session in a handler
type userSessionIf interface {
	getEmail() string
	getName() string
	getUserName() string
	isAdmin() bool
	isValidForSession() bool
}

type userGetterIf interface {
	getId() int
	getEmail() string
	getName() string
	getUserName() string
	getSessionExpiresIn() int
	toStringInterface() map[string]interface{}
	findMe() (bool, error)
	helpers.FfjsonIf
}

type userSetterIf interface {
	setEmail(string) error
	setName(string) error
	setUserName(string) error
	setAuthenticated(bool) error
	setSessionExpiresIn(time.Duration) error
	prepareNew() error
	applyDbData(map[string]interface{}) error
	// validate(ctx *context, r *http.Request, errors map[string]string) error
	generateRecoveryCode() (string, error)
	resetRecoveryCode()
	generatePassword() error
	changePassword(string) error
	encryptPassword() error
	unsetPassword()
}

type userPermissionsIf interface {
	isAuthenticated() bool
	isAdmin() bool
	isActive() bool
	checkPassword(string) bool
	isValidForSession() bool
}

type UserModelCollection struct {
	Users []userIf
}

type UserModel struct {
	CreatedAt        time.Time
	UserName         string
	Name             string
	Email            string
	Password         string
	IsAdmin          bool
	IsActive         bool
	RecoveryCode     string
	IsAuthenticated  bool
	SessionExpiresIn time.Duration // not exported in JSON
}

func (u *UserModel) getId() int               { return helpers.StringHash(u.UserName) }
func (u *UserModel) getEmail() string         { return u.Email }
func (u *UserModel) getUserName() string      { return u.UserName }
func (u *UserModel) getName() string          { return u.Name }
func (u *UserModel) getSessionExpiresIn() int { return int(u.SessionExpiresIn.Seconds()) }

func (u *UserModel) setEmail(e string) error                    { u.Email = e; return nil }
func (u *UserModel) setName(n string) error                     { u.Name = n; return nil }
func (u *UserModel) setUserName(n string) error                 { u.UserName = n; return nil }
func (u *UserModel) setAuthenticated(auth bool) error           { u.IsAuthenticated = auth; return nil }
func (u *UserModel) setSessionExpiresIn(ei time.Duration) error { u.SessionExpiresIn = ei; return nil }

func (u *UserModel) isAuthenticated() bool { return u.IsAuthenticated }
func (u *UserModel) isAdmin() bool         { return u.IsAdmin }
func (u *UserModel) isActive() bool        { return u.IsActive }

// PreInsert hook for new users
func (u *UserModel) prepareNew() error {
	u.IsActive = true
	u.CreatedAt = time.Now()
	return u.encryptPassword()
}

// isValidForSession() is only used in newSessionInfo()
func (u *UserModel) isValidForSession() bool {
	return true == helpers.ValidateEmail(u.getEmail()) && "" != u.getUserName() && true == u.isAuthenticated()
}

//func (UserModel *UserModel) validate(ctx *context, r *http.Request, errors map[string]string) error {
//
//	if UserModel.Name == "" {
//		errors["name"] = "Name is missing"
//	} else {
//		ok, err := ctx.datamapper.isUserNameAvailable(UserModel)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["name"] = "Name already taken"
//		}
//	}
//
//	if UserModel.Email == "" {
//		errors["email"] = "Email is missing"
//	} else if !validateEmail(UserModel.Email) {
//		errors["email"] = "Invalid email address"
//	} else {
//		ok, err := ctx.datamapper.isUserEmailAvailable(UserModel)
//		if err != nil {
//			return err
//		}
//		if !ok {
//			errors["email"] = "Email already taken"
//		}
//
//	}
//
//	// tbd: we need flag UserModel is third-party
//	if UserModel.Password == "" {
//		errors["password"] = "Password is missing"
//	}
//
//	return nil
//}

func (u *UserModel) generateRecoveryCode() (string, error) {
	code := helpers.RandomString(USER_RECOVERY_CODE_LENGTH)
	u.RecoveryCode = code
	return code, nil
}

func (u *UserModel) resetRecoveryCode() {
	u.RecoveryCode = ""
}

// generates an unencrypted password
func (u *UserModel) generatePassword() error {
	var err error
	u.Password, err = helpers.NewPassword(USER_PASSWORD_LENGTH)
	return err
}

func (u *UserModel) changePassword(password string) error {
	u.Password = password
	return u.encryptPassword()
}

func (u *UserModel) unsetPassword() {
	u.Password = ""
}

func (u *UserModel) encryptPassword() error {
	if "" == u.Password {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if nil != err {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// not sure if it is a good idea to carry the whole time the bcrypted password with the UserModel object ...
func (u *UserModel) checkPassword(password string) bool {
	if "" == u.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *UserModel) toStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"CreatedAt":       u.CreatedAt.Unix(),
		"UserName":        u.UserName,
		"Name":            u.Name,
		"Email":           u.Email,
		"Password":        u.Password,
		"IsAdmin":         u.IsAdmin,
		"IsActive":        u.IsActive,
		"IsAuthenticated": u.IsAuthenticated,
	}
}

// finds a user in the database and fills the struct
func (u *UserModel) findMe() (bool, error) {
	searchedUser, _ := rsdb.FindOne(USER_DB_COLLECTION_NAME, u.getId())
	if nil == searchedUser {
		return false, nil
	}
	u.applyDbData(searchedUser)

	return true, nil
}

func (u *UserModel) applyDbData(d map[string]interface{}) error {
	// panic free type conversion
	tIsAdmin, _ := d["IsAdmin"].(bool)
	tIsActive, _ := d["IsActive"].(bool)
	tIsAuthenticated, _ := d["IsAuthenticated"].(bool)
	tCreatedAt, _ := d["CreatedAt"].(int64)
	tUserName, _ := d["UserName"].(string)
	tName, _ := d["Name"].(string)
	tEmail, _ := d["Email"].(string)
	tPassword, _ := d["Password"].(string)

	u.CreatedAt = time.Unix(tCreatedAt, 0)
	u.UserName = tUserName
	u.Name = tName
	u.Email = tEmail
	u.Password = tPassword
	u.IsAdmin = tIsAdmin
	u.IsActive = tIsActive
	u.IsAuthenticated = tIsAuthenticated
	return nil
}

// needed in auth when user tries to login
func NewUserModel(userName string) *UserModel {
	u := &UserModel{
		UserName: userName,
	}
	return u
}

// GetAllUsers returns a user collection with empty passwords
func GetAllUsers() (*UserModelCollection, error) {
	col, err := rsdb.FindAll(USER_DB_COLLECTION_NAME)
	umc := &UserModelCollection{}
	for _, u := range col {
		newUser := NewUserModel("")
		newUser.applyDbData(u)
		newUser.unsetPassword()
		umc.Users = append(umc.Users, newUser)
	}
	return umc, err
}

// initUsers() runs in NewPicnicApp() function
func initUsers() error {
	var err error
	var root map[string]interface{}
	var pwd string
	err = rsdb.CreateDatabase(USER_DB_COLLECTION_NAME)
	if nil != err {
		return errgo.Mask(err)
	}

	u := UserModel{
		UserName: USER_ROOT,
		Name:     "Default Root User",
		Email:    USER_ROOT + "@localhost",
		Password: pwd,
		IsAdmin:  true,
		IsActive: true,
	}
	u.generatePassword()
	root, _ = rsdb.FindOne(USER_DB_COLLECTION_NAME, u.getId())

	if nil == root {
		logger.Emergency("Created new user %s with password: %s", u.UserName, u.Password)
		u.prepareNew()
		rsdb.InsertRecovery(USER_DB_COLLECTION_NAME, u.getId(), u.toStringInterface())
	} else {
		logger.Emergency("Root user %s already exists!", USER_ROOT)
	}
	return errgo.Mask(err)
}
