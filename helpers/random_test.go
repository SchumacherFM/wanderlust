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
	"testing"
)

func TestRandomString(t *testing.T) {
	rs := RandomString(10)
	if 10 != len(rs) {
		t.Errorf("String %s is not 10 char long", rs)
	}
	// test that string is in range
}

/**
05. Sept 2014; run with `go test -v -bench=.`
BenchmarkRandomString	  200000	     11956 ns/op Late 2013; 2.4 GHz Intel Core i5; OS X 10.9.4 (13E28)
BenchmarkRandomString	  200000	     14001 ns/op Mid  2012; 1.8 GHz Intel Core i5; OS X 10.9.4 (13E28)
go version go1.3.1 darwin/amd64
*/
func BenchmarkRandomString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomString(10)
	}
}
