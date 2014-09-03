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
	"os"
	"testing"
)

func TestGetTempDir(t *testing.T) {
	dir := GetTempDir()
	if "" == dir {
		t.Error("No tempdir found!")
	}

	if string(os.PathSeparator) != dir[len(dir)-1:] {
		t.Errorf("Trailing %s not found in %s", os.PathSeparator, dir)
	}
}
func TestPathExists(t *testing.T) {
	dir := GetTempDir()
	isPath, err := PathExists(dir)
	if nil != err {
		t.Error(err)
	}
	if false == isPath {
		t.Errorf("Path %s not found", dir)
	}
}
func TestCreateDirectoryIfNotExists(t *testing.T) {
	dir := GetTempDir() + "testing_" + RandomString(10)
	CreateDirectoryIfNotExists(dir)
	defer os.Remove(dir)
	isPath, err := PathExists(dir)
	if nil != err {
		t.Error(err)
	}
	if false == isPath {
		t.Errorf("Path %s not found", dir)
	}
}
