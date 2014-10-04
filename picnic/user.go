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

//type permissions struct {
//	Edit   bool `json:"edit"`
//	Delete bool `json:"delete"`
//}

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

func (um *UserModel) getId() int               { return helpers.StringHash(um.UserName) }
func (um *UserModel) getEmail() string         { return um.Email }
func (um *UserModel) getUserName() string      { return um.UserName }
func (um *UserModel) getName() string          { return um.Name }
func (um *UserModel) getSessionExpiresIn() int { return int(um.SessionExpiresIn.Seconds()) }

func (um *UserModel) setEmail(e string) error                    { um.Email = e; return nil }
func (um *UserModel) setName(n string) error                     { um.Name = n; return nil }
func (um *UserModel) setUserName(u string) error                 { um.UserName = u; return nil }
func (um *UserModel) setAuthenticated(auth bool) error           { um.IsAuthenticated = auth; return nil }
func (um *UserModel) setSessionExpiresIn(ei time.Duration) error { um.SessionExpiresIn = ei; return nil }

func (um *UserModel) isAuthenticated() bool { return um.IsAuthenticated }
func (um *UserModel) isAdmin() bool         { return um.IsAdmin }
func (um *UserModel) isActive() bool        { return um.IsActive }

// PreInsert hook for new users
func (um *UserModel) prepareNew() error {
	um.IsActive = true
	um.CreatedAt = time.Now()
	return um.encryptPassword()
}

// isValidForSession() is only used in newSessionInfo()
func (um *UserModel) isValidForSession() bool {
	return true == helpers.ValidateEmail(um.getEmail()) && "" != um.getUserName() && true == um.isAuthenticated()
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

func (um *UserModel) generateRecoveryCode() (string, error) {
	code := helpers.RandomString(USER_RECOVERY_CODE_LENGTH)
	um.RecoveryCode = code
	return code, nil
}

func (um *UserModel) resetRecoveryCode() {
	um.RecoveryCode = ""
}

// generates an unencrypted password
func (um *UserModel) generatePassword() error {
	var err error
	um.Password, err = helpers.NewPassword(USER_PASSWORD_LENGTH)
	return err
}

func (um *UserModel) changePassword(password string) error {
	um.Password = password
	return um.encryptPassword()
}

func (um *UserModel) unsetPassword() {
	um.Password = ""
}

func (um *UserModel) encryptPassword() error {
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

// not sure if it is a good idea to carry the whole time the bcrypted password with the UserModel object ...
func (um *UserModel) checkPassword(password string) bool {
	if "" == um.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(um.Password), []byte(password))
	return err == nil
}

func (um *UserModel) toStringInterface() map[string]interface{} {
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
func (um *UserModel) findMe() (bool, error) {
	searchedUser, _ := rsdb.FindOne(USER_DB_COLLECTION_NAME, um.getId())
	if nil == searchedUser {
		return false, nil
	}
	um.applyDbData(searchedUser)

	return true, nil
}

func (um *UserModel) applyDbData(data map[string]interface{}) error {
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
	return nil
}

// needed in auth when user tries to login
func NewUserModel(userName string) *UserModel {
	um := &UserModel{
		UserName: userName,
	}
	return um
}

// GetAllUsers returns a user collection with empty passwords
func GetAllUsers() (*UserModelCollection, error) {
	dbUserCollection, err := rsdb.FindAll(USER_DB_COLLECTION_NAME)
	umc := &UserModelCollection{}
	for _, rawUser := range dbUserCollection {
		newUser := NewUserModel("")
		newUser.applyDbData(rawUser)
		newUser.unsetPassword()
		umc.Users = append(umc.Users, newUser)
	}
	return umc, err
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

	um := UserModel{
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
