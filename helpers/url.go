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

func ValidateListenAddress(address string) (string, string, error) {
	var host string
	var port int
	var err error
	parts := strings.Split(address, ":")
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
