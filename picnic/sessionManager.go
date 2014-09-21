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
	jwt "github.com/SchumacherFM/wanderlust/github.com/dgrijalva/jwt-go"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"github.com/SchumacherFM/wanderlust/picnic/middleware"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	TOKEN_EXPIRY = 60 // minutes
)

type sessionManagerI interface {
	readToken(*http.Request) (string, error)
	_createToken(string) (string, error)
	writeToken(http.ResponseWriter, string) error
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
	verifyKey []byte
	signKey   []byte
}

func (m *defaultSessionManager) readToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get(middleware.HEADER_X_AUTH_TOKEN)
	if tokenString == "" {
		return "", nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.verifyKey, nil
	})
	switch err.(type) {
	case nil:
		if !token.Valid {
			return "", nil
		}
		token := token.Claims["uid"].(string)
		if "" == token {
			return "", nil
		}
		return token, nil
	case *jwt.ValidationError:
		return "", nil
	default:
		return "", errgo.Mask(err)
	}
}

// creates and signs a token, private method
func (m *defaultSessionManager) _createToken(userID string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["uid"] = userID
	token.Claims["exp"] = time.Now().Add(time.Minute * TOKEN_EXPIRY).Unix()
	tokenString, err := token.SignedString(m.signKey)
	if err != nil {
		return tokenString, errgo.Mask(err)
	}
	return tokenString, nil
}

func (m *defaultSessionManager) writeToken(w http.ResponseWriter, userID string) error {
	tokenString, err := m._createToken(userID)
	if err != nil {
		return err
	}
	w.Header().Set(middleware.HEADER_X_AUTH_TOKEN, tokenString)
	return nil
}
