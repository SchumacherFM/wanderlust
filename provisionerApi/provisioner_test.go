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

package provisionerApi

import (
	"github.com/SchumacherFM/wanderlust/rucksack"
	"testing"
)

type testPapi struct {
}

func (this *testPapi) Route() string    { return "TestRoute" }
func (this *testPapi) Config() []string { return nil }

func (this *testPapi) PrepareSave(pd *PostData) ([]byte, error) { return nil, nil }

func (this *testPapi) ConfigComplete(bp rucksack.Backpacker) (bool, error) { return false, nil }

func (this *testPapi) FetchUrls(bp rucksack.Backpacker) []string { return nil }

func TestNewProvisioner(t *testing.T) {
	papi := &testPapi{}
	p := NewProvisioner("TestProv", "TestIcon", papi)
	if e := "/" + UrlRoutePrefix + "/TestRoute"; p.url != e {
		t.Errorf("\nActual\t\t%s\nExpected\t%s\n", p.url, e)
	}
}
