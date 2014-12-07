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
	"encoding/json"
	"time"

	"github.com/SchumacherFM/wanderlust/code.google.com/p/go.crypto/bcrypt"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/picnicApi"
	"github.com/SchumacherFM/wanderlust/rucksack"
)

var (
	// check for interfaces
	_ helpers.FfjsonIf        = (*userModel)(nil)
	_ picnicApi.UserSessionIf = (*userModel)(nil)
	_ picnicApi.UserGetterIf  = (*userModel)(nil)
	_ picnicApi.UserSetterIf  = (*userModel)(nil) // @todo revise: not sure if needed or incorrect named
)

const (
	USER_PASSWORD_LENGTH      = 24
	USER_RECOVERY_CODE_LENGTH = 30
	USER_DB_COLLECTION_NAME   = "users"
	USER_ROOT                 = "administrator"
)

type userModelCollection struct {
	db    rucksack.Backpacker
	Users []*userModel
}

type userModel struct {
	db rucksack.Backpacker
	// @todo field names to lower case as they are private ?
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

// needed in auth when user tries to login
func NewUserModel(db rucksack.Backpacker, userName string) *userModel {
	u := &userModel{
		db:        db,
		UserName:  userName,
		CreatedAt: time.Now(),
	}
	return u
}

func (u *userModel) GetId() string            { return helpers.StringHashString(u.UserName) }
func (u *userModel) GetEmail() string         { return u.Email }
func (u *userModel) GetUserName() string      { return u.UserName }
func (u *userModel) GetName() string          { return u.Name }
func (u *userModel) GetSessionExpiresIn() int { return int(u.SessionExpiresIn.Seconds()) }

func (u *userModel) SetAuthenticated(auth bool) error           { u.IsAuthenticated = auth; return nil }
func (u *userModel) SetSessionExpiresIn(ei time.Duration) error { u.SessionExpiresIn = ei; return nil }

func (u *userModel) IsLoggedIn() bool      { return u.IsAuthenticated }
func (u *userModel) IsAdministrator() bool { return u.IsAdmin }
func (u *userModel) IsActive() bool        { return u.IsActivated }

// PreInsert hook for new users
func (u *userModel) prepareNew() error {
	u.IsActivated = true
	u.CreatedAt = time.Now()
	return u.EncryptPassword()
}

func (u *userModel) Decode(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *userModel) Encode() ([]byte, error) {
	return u.MarshalJSON()
}

// IsValidForSession() is only used in newSessionInfo()
func (u *userModel) IsValidForSession() bool {
	return true == helpers.ValidateEmail(u.GetEmail()) && "" != u.GetUserName() && true == u.IsLoggedIn()
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

func (u *userModel) GenerateRecoveryCode() (string, error) {
	code := helpers.RandomString(USER_RECOVERY_CODE_LENGTH)
	u.RecoveryCode = code
	return code, nil
}

func (u *userModel) ResetRecoveryCode() {
	u.RecoveryCode = ""
}

// generates an unencrypted password
func (u *userModel) GeneratePassword() error {
	var err error
	u.Password, err = helpers.NewPassword(USER_PASSWORD_LENGTH)
	return err
}

func (u *userModel) ChangePassword(password string) error {
	u.Password = password
	return u.EncryptPassword()
}

func (u *userModel) UnsetPassword() {
	u.Password = ""
}

func (u *userModel) EncryptPassword() error {
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

// not sure if it is a good idea to carry the whole time the bcrypted password with the userModel object ...
func (u *userModel) CheckPassword(password string) bool {
	if "" == u.Password {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// finds a user in the database and fills the struct
func (u *userModel) FindMe() (bool, error) {
	b, err := u.db.FindOne(USER_DB_COLLECTION_NAME, u.GetId())
	if nil == b || nil != err {
		return false, err
	}
	err = u.Decode(b)
	return true, err
}

// needed in auth when user tries to login
func NewUserModelCollection(db rucksack.Backpacker) *userModelCollection {
	uc := &userModelCollection{
		db: db,
	}
	return uc
}

// FindAllUsers populates the internal User slice.
// Returns errors like Database, tx failed or userIds are not matching
func (uc *userModelCollection) FindAllUsers() error {
	d, err := uc.db.FindAll(USER_DB_COLLECTION_NAME)
	if nil != err {
		return err
	}

	for i := 0; i < len(d); i = i + 2 {
		id := string(d[i])
		newUser := NewUserModel(nil, "") // no database connection needed
		newUser.Decode(d[i+1])
		newUser.UnsetPassword()

		if newUser.GetId() != id {
			return errgo.Newf("UserID corrupt!\nExpected: %s\nActual: %s\n", newUser.GetId(), id)
		}

		uc.Users = append(uc.Users, newUser)
	}
	return nil
}

// initUsers() runs in NewPicnicApp() function
func initUsers(db rucksack.Backpacker) error {
	var pwd string

	u := &userModel{
		UserName:    USER_ROOT,
		Name:        "Default Root User",
		Email:       USER_ROOT + "@localhost",
		Password:    pwd,
		IsAdmin:     true,
		IsActivated: true,
	}
	u.GeneratePassword()
	userByte, err := db.FindOne(USER_DB_COLLECTION_NAME, u.GetId())

	switch err {
	case nil:
		break
	case rucksack.ErrBreadNotFound:
		break
	default:
		return errgo.Mask(err)
	}

	if nil == userByte {
		logger.Emergency("Created new user %s with password: %s", u.UserName, u.Password)
		u.prepareNew()
		ue, err := u.Encode()
		logger.Check(err)
		db.Insert(USER_DB_COLLECTION_NAME, u.GetId(), ue)
	} else {
		logger.Emergency("Root user %s already exists!", USER_ROOT)
	}
	return nil
}
