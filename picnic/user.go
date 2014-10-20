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
	. "github.com/SchumacherFM/wanderlust/picnic/api"
	"time"
)

const (
	USER_PASSWORD_LENGTH      = 24
	USER_RECOVERY_CODE_LENGTH = 30
	USER_DB_COLLECTION_NAME   = "users"
	USER_ROOT                 = "administrator"
)

type UserModelCollection struct {
	Users []UserIf
}

type UserModel struct {
	CreatedAt        time.Time
	UserName         string
	Name             string
	Email            string
	Password         string
	IsAdmin          bool
	IsActivated      bool
	RecoveryCode     string
	IsAuthenticated  bool
	SessionExpiresIn time.Duration // not exported in JSON
}

func (u *UserModel) GetId() int               { return helpers.StringHash(u.UserName) }
func (u *UserModel) GetEmail() string         { return u.Email }
func (u *UserModel) GetUserName() string      { return u.UserName }
func (u *UserModel) GetName() string          { return u.Name }
func (u *UserModel) GetSessionExpiresIn() int { return int(u.SessionExpiresIn.Seconds()) }

func (u *UserModel) SetEmail(e string) error                    { u.Email = e; return nil }
func (u *UserModel) SetName(n string) error                     { u.Name = n; return nil }
func (u *UserModel) SetUserName(n string) error                 { u.UserName = n; return nil }
func (u *UserModel) SetAuthenticated(auth bool) error           { u.IsAuthenticated = auth; return nil }
func (u *UserModel) SetSessionExpiresIn(ei time.Duration) error { u.SessionExpiresIn = ei; return nil }

func (u *UserModel) IsLoggedIn() bool      { return u.IsAuthenticated }
func (u *UserModel) IsAdministrator() bool { return u.IsAdmin }
func (u *UserModel) IsActive() bool        { return u.IsActivated }

// PreInsert hook for new users
func (u *UserModel) prepareNew() error {
	u.IsActivated = true
	u.CreatedAt = time.Now()
	return u.EncryptPassword()
}

// IsValidForSession() is only used in newSessionInfo()
func (u *UserModel) IsValidForSession() bool {
	return true == helpers.ValidateEmail(u.GetEmail()) && "" != u.GetUserName() && true == u.IsLoggedIn()
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

func (u *UserModel) GenerateRecoveryCode() (string, error) {
	code := helpers.RandomString(USER_RECOVERY_CODE_LENGTH)
	u.RecoveryCode = code
	return code, nil
}

func (u *UserModel) ResetRecoveryCode() {
	u.RecoveryCode = ""
}

// generates an unencrypted password
func (u *UserModel) GeneratePassword() error {
	var err error
	u.Password, err = helpers.NewPassword(USER_PASSWORD_LENGTH)
	return err
}

func (u *UserModel) ChangePassword(password string) error {
	u.Password = password
	return u.EncryptPassword()
}

func (u *UserModel) UnsetPassword() {
	u.Password = ""
}

func (u *UserModel) EncryptPassword() error {
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
func (u *UserModel) CheckPassword(password string) bool {
	if "" == u.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *UserModel) ToStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"CreatedAt":       u.CreatedAt.Unix(),
		"UserName":        u.UserName,
		"Name":            u.Name,
		"Email":           u.Email,
		"Password":        u.Password,
		"IsAdmin":         u.IsAdmin,
		"IsActivated":     u.IsActivated,
		"IsAuthenticated": u.IsAuthenticated,
	}
}

// finds a user in the database and fills the struct
func (u *UserModel) FindMe() (bool, error) {
	searchedUser, _ := rsdb.FindOne(USER_DB_COLLECTION_NAME, u.GetId())
	if nil == searchedUser {
		return false, nil
	}
	u.ApplyDbData(searchedUser)

	return true, nil
}

func (u *UserModel) ApplyDbData(d map[string]interface{}) error {
	// panic free type conversion
	tIsAdmin, _ := d["IsAdmin"].(bool)
	tIsActivated, _ := d["IsActivated"].(bool)
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
	u.IsActivated = tIsActivated
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
		newUser.ApplyDbData(u)
		newUser.UnsetPassword()
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
		UserName:    USER_ROOT,
		Name:        "Default Root User",
		Email:       USER_ROOT + "@localhost",
		Password:    pwd,
		IsAdmin:     true,
		IsActivated: true,
	}
	u.GeneratePassword()
	root, _ = rsdb.FindOne(USER_DB_COLLECTION_NAME, u.GetId())

	if nil == root {
		logger.Emergency("Created new user %s with password: %s", u.UserName, u.Password)
		u.prepareNew()
		rsdb.InsertRecovery(USER_DB_COLLECTION_NAME, u.GetId(), u.ToStringInterface())
	} else {
		logger.Emergency("Root user %s already exists!", USER_ROOT)
	}
	return errgo.Mask(err)
}
