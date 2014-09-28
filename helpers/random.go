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

package helpers

import (
	crand "crypto/rand"
	"errors"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"io"
	mrand "math/rand"
	"time"
)

// Hash a string using sdbm algorithm.
func StringHash(str string) int {
	var hash int
	for _, c := range str {
		hash = int(c) + (hash << 6) + (hash << 16) - hash
	}
	if hash < 0 {
		return -hash
	}
	return hash
}

func RandomInt(n int) int {
	mrand.Seed(time.Now().UnixNano())
	return mrand.Intn(n)
}

// randomString generates a pseudo-random alpha-numeric string with given length.
func RandomString(length int) string {
	mrand.Seed(time.Now().UnixNano())
	k := make([]rune, length)
	for i := 0; i < length; i++ {
		c := mrand.Intn(35)
		if c < 10 {
			c += 48 // numbers (0-9) (0+48 == 48 == '0', 9+48 == 57 == '9')
		} else {
			c += 87 // lower case alphabets (a-z) (10+87 == 97 == 'a', 35+87 == 122 = 'z')
		}
		k[i] = rune(c)
	}
	return string(k)
}

var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'\"!@#%^&*()-_=+,.?/:;{}[]`~")

// generates a random string with the crypto package
func NewPassword(length int) (string, error) {
	return randChar(length, stdChars)
}

func randChar(length int, chars []byte) (string, error) {
	newPword := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(crand.Reader, randomData); err != nil {
			return "", err
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}
			newPword[i] = chars[c%clen]
			i++
			if i == length {
				return string(newPword), nil
			}
		}
	}
}
