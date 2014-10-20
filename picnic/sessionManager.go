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

func newSessionManager(publicKeyFilePath, privateKeyFilePath string) (*defaultSessionManager, error) {
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

// readToken, reads the token, validates it and returns the uid, validity in seconds and err|nil
func (m *defaultSessionManager) ReadToken(r *http.Request) (string, time.Duration, error) {
	tokenString := r.Header.Get(middleware.HEADER_X_AUTH_TOKEN)
	if tokenString == "" {
		return "", 0, nil
	}
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return m.verifyKey, nil
	})
	switch err.(type) {
	case nil:
		if !t.Valid {
			return "", 0, nil
		}
		tokenUid := t.Claims["uid"].(string)
		if "" == tokenUid {
			return "", 0, nil
		}
		exp, expOK := t.Claims["exp"].(float64)
		if false == expOK {
			return "", 0, nil
		}
		expTm := time.Unix(int64(exp), 0).Unix()
		dur := time.Duration((expTm - time.Now().Unix()))
		dur = dur * time.Second
		return tokenUid, dur, nil
	case *jwt.ValidationError:
		return "", 0, nil
	default:
		return "", 0, errgo.Mask(err)
	}
}

// creates and signs a token, private method
func (m *defaultSessionManager) createToken(userID string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims["uid"] = userID
	t.Claims["exp"] = time.Now().Add(time.Minute * TOKEN_EXPIRY).Unix()
	tokenString, err := t.SignedString(m.signKey)
	if err != nil {
		return tokenString, errgo.Mask(err)
	}
	return tokenString, nil
}

func (m *defaultSessionManager) WriteToken(w http.ResponseWriter, userID string) error {
	tokenString, err := m.createToken(userID)
	if err != nil {
		return err
	}
	w.Header().Set(middleware.HEADER_X_AUTH_TOKEN, tokenString)
	return nil
}
