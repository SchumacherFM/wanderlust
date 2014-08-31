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
	"time"
	"math/rand"
)

// randomString generates a pseudo-random alpha-numeric string with given length.
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	k := make([]rune, length)
	for i := 0; i < length; i++ {
		c := rand.Intn(35)
		if c < 10 {
			c += 48 // numbers (0-9) (0+48 == 48 == '0', 9+48 == 57 == '9')
		} else {
			c += 87 // lower case alphabets (a-z) (10+87 == 97 == 'a', 35+87 == 122 = 'z')
		}
		k[i] = rune(c)
	}
	return string(k)
}
