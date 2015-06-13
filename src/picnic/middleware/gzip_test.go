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

package middleware

import (
	"testing"
)

// run via $ go test -v --bench=. -test.benchmem .

var (
	prepareRequestUriMap = map[string]string{
		"":                                                         "/index.html",
		"/":                                                        "/index.html",
		"dashboard/":                                               "dashboard/index.html",
		"dashboard":                                                "dashboard",
		"/lib/bootstrap/css/bootstrap.css":                         "/lib/bootstrap/css/bootstrap.css",
		"/lib/font-awesome/fonts/fontawesome-webfont.woff?v=4.1.0": "/lib/font-awesome/fonts/fontawesome-webfont.woff",
		"/lib/angular/angular.js?v=2.1&a=b":                        "/lib/angular/angular.js",
	}
	prepareRequestUriSlice = []string{
		"",            // input
		"/index.html", // expected
		"/",           // input
		"/index.html", // expected
		"dashboard/",
		"dashboard/index.html",
		"dashboard",
		"dashboard",
		"/lib/bootstrap/css/bootstrap.css",
		"/lib/bootstrap/css/bootstrap.css",
		"/lib/font-awesome/fonts/fontawesome-webfont.woff?v=4.1.0",
		"/lib/font-awesome/fonts/fontawesome-webfont.woff",
		"/lib/angular/angular.js?v=2.1&a=b",
		"/lib/angular/angular.js",
	}
)

func TestPrepareRequestUri(t *testing.T) {
	for input, expected := range prepareRequestUriMap {
		if expected != prepareRequestUri(input) {
			t.Error("No equal: ", input, expected)
		}
	}
}

// BenchmarkPrepareRequestUriMap	 1000000	      1230 ns/op	     150 B/op	       6 allocs/op
func BenchmarkPrepareRequestUriMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for input, expected := range prepareRequestUriMap {
			if expected != prepareRequestUri(input) {
				b.Log(expected)
			}
		}
	}
}

// Using an array instead of a slice is around 1180 ns/op
// BenchmarkPrepareRequestUriSlice	 1000000	      1075 ns/op	     150 B/op	       6 allocs/op
func BenchmarkPrepareRequestUriSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(prepareRequestUriSlice); i += 2 {
			input := prepareRequestUriSlice[i]
			expected := prepareRequestUriSlice[i+1]
			//b.Log(input,expected)
			if expected != prepareRequestUri(input) {
				b.Error(expected)
			}
		}
	}
}

// @todo
func TestGzipContentTypeMiddleware(t *testing.T) {

}
