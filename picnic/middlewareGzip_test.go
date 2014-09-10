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

// https://github.com/phyber/negroni-gzip/blob/master/LICENSE
// The MIT License (MIT)
// Copyright (c) 2013 Jeremy Saenz; 2014 David O'Rourke
// and modified by Cyrill

package picnic

import (
	"testing"
)

var (
	prepareRequestUriTests = map[string]string{
	"": "/index.html",
	"/": "/index.html",
	"/lib/bootstrap/css/bootstrap.css": "/lib/bootstrap/css/bootstrap.css",
	"/lib/font-awesome/fonts/fontawesome-webfont.woff?v=4.1.0": "/lib/font-awesome/fonts/fontawesome-webfont.woff",
	"/lib/angular/angular.js?v=2.1&a=b": "/lib/angular/angular.js",
}
)

func TestPrepareRequestUri(t *testing.T) {
	for input, expected := range prepareRequestUriTests {
		if expected != prepareRequestUri(input) {
			t.Error("No equal: ", input, expected)
		}
	}
}
