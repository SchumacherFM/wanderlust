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
	"github.com/SchumacherFM/wanderlust/github.com/dgrijalva/jwt"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	TOKEN_EXPIRY = 60 // minutes
)

type sessionManagerI interface {
	readToken(*http.Request) (int, error)
	createToken(int) (string, error)
	writeToken(http.ResponseWriter, int) error
}

// Basic user session info
type sessionInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
	LoggedIn bool   `json:"loggedIn"`
}

func newSessionInfo(user *picnicer) *sessionInfo {
	if user == nil || user.ID == 0 || !user.IsAuthenticated {
		return &sessionInfo{}
	}

	return &sessionInfo{user.ID, user.Name, user.Email, user.IsAdmin, true}
}

func newSessionManager(publicKeyFilePath, privateKeyFilePath string) (sessionManagerI, error) {
	mgr := &defaultSessionManager{}
	var err error
	mgr.signKey, err = ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return mgr, errgo.Mask(err)
	}
	mgr.verifyKey, err = ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return mgr, errgo.Mask(err)
	}
	return mgr, nil
}

type defaultSessionManager struct {
	verifyKey, signKey []byte
}

func (m *defaultSessionManager) readToken(r *http.Request) (int, error) {
	tokenString := r.Header.Get(HEADER_X_AUTH_TOKEN)
	if tokenString == "" {
		return 0, nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.verifyKey, nil
	})
	switch err.(type) {
	case nil:
		if !token.Valid {
			return 0, nil
		}
		token := token.Claims["uid"].(string)
		userID, err := strconv.Atoi(token)
		if err != nil {
			return 0, nil
		}
		return userID, nil
	case *jwt.ValidationError:
		return 0, nil
	default:
		return 0, errgo.Mask(err)
	}
}

func (m *defaultSessionManager) createToken(userID int) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["uid"] = strconv.Itoa(userID)
	token.Claims["exp"] = time.Now().Add(time.Minute * TOKEN_EXPIRY).Unix()
	tokenString, err := token.SignedString(m.signKey)
	if err != nil {
		return tokenString, errgo.Mask(err)
	}
	return tokenString, nil
}

func (m *defaultSessionManager) writeToken(w http.ResponseWriter, userID int) error {
	tokenString, err := m.createToken(userID)
	if err != nil {
		return err
	}
	w.Header().Set(HEADER_X_AUTH_TOKEN, tokenString)
	return nil
}
