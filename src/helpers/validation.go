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
	"errors"
	"strconv"
	"strings"
)

var (
	validRunesUser = make(map[rune]bool)
	validRunesHost = map[rune]bool{
		'a': true,
		'b': true,
		'c': true,
		'd': true,
		'e': true,
		'f': true,
		'g': true,
		'h': true,
		'i': true,
		'j': true,
		'k': true,
		'l': true,
		'm': true,
		'n': true,
		'o': true,
		'p': true,
		'q': true,
		'r': true,
		's': true,
		't': true,
		'u': true,
		'v': true,
		'w': true,
		'x': true,
		'y': true,
		'z': true,
		'0': true,
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
		'-': true,
		'.': true,
	}
)

func init() {
	for k, v := range validRunesHost {
		validRunesUser[k] = v
	}
	validRunesUser['_'] = true
	validRunesUser['+'] = true
	validRunesUser[''] = true // ;-)
}

func checkRune(validMap map[rune]bool, checkString string) bool {
	checkRunes := []rune(checkString)
	for i := 0; i < len(checkRunes); i++ {
		_, isSet := validMap[checkRunes[i]]
		if false == isSet {
			return false
		}
	}
	return true
}

// ValidateEmail does what is says in a simple way and without regex. everybody is doing it with regex ;-)
func ValidateEmail(email string) bool {
	email = strings.ToLower(email)

	if len(email) < 3 {
		return false
	}
	emailParts := strings.Split(email, "@")
	if 2 != len(emailParts) {
		return false // e.g. hello@w@rld or hello
	}
	user := emailParts[0]
	host := emailParts[1]

	isValid := checkRune(validRunesUser, user) && checkRune(validRunesHost, host)
	if true == isValid {
		isValid = false == strings.Contains(user, "..") && false == strings.Contains(host, "..")
	}
	return isValid
}

func ValidateListenAddress(address string) (string, string, error) {
	var host string
	var port int
	var err error
	parts := strings.Split(address, ":")
	if 2 != len(parts) {
		return "", "", errors.New("Missing : separator or too many")
	}
	host = parts[0]

	port, err = strconv.Atoi(parts[1])
	if nil != err {
		return "", "", err
	}
	if "" == host {
		host = "127.0.0.1"
	}
	if port < 1 {
		return "", "", errors.New("Port is zero!")
	}
	return host, strconv.Itoa(port), nil
}
