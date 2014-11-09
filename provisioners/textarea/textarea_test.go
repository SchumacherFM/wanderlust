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

package textarea

import (
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	. "github.com/SchumacherFM/wanderlust/rucksack/rsTestHelper"
	"strings"
	"testing"
)

func TestIsValid(t *testing.T) {
	p := GetProvisioner()

	pd := &provisionerApi.PostData{}

	err := p.Api.IsValid(pd) // must succeed
	if nil != err {
		t.Errorf("%#v is valid 31!", pd)
	}

	pd.Value = " http://www.golang.org/tour"
	err = p.Api.IsValid(pd) // must succeed
	if nil != err {
		t.Errorf("%#v is not valid 37!", pd)
	}

	pd.Value = "htp://www.golang.org/siteap.xml"
	err = p.Api.IsValid(pd) // must fail
	if nil == err {
		t.Errorf("%#v is not valid 43!", pd)
	}

	pd.Value = "hTtp://www.golang.org/siteMap.xml\nhttps://www.google.com/search?q=golang\n"
	err = p.Api.IsValid(pd) // must succeed
	if nil != err {
		t.Errorf("%#v must be valid 49!", pd)
	}

	pd.Value = strings.Repeat("hTtp://www.golang.org/siteMap.xml\n", 21)
	err = p.Api.IsValid(pd) // must fail
	if nil == err {
		t.Errorf("%#v must invalid 57!", pd)
	}

}

// MBA Mid 2012 1.8 GHz Intel Core i5
// BenchmarkIsValid	 2.000.000	      1008 ns/op
func BenchmarkIsValid(b *testing.B) {
	p := GetProvisioner()
	pd := &provisionerApi.PostData{}
	pd.Value = "hTtp://www.golang.org/siteMap.xml"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		p.Api.IsValid(pd) // must succeed
	}
}

func TestConfigComplete(t *testing.T) {
	p := GetProvisioner()
	db := &DbMock{
		FindOneData: []byte(`http://www.golang.org/sitemap.xml`),
	}
	ok, err := p.Api.ConfigComplete(db)
	if false == ok {
		t.Error("Config must be complete!")
	}
	if nil != err {
		t.Error(err)
	}
	db.FindOneData = []byte(``)
	ok, err = p.Api.ConfigComplete(db)
	if true == ok {
		t.Error("Config is not complete!")
	}
	if nil != err {
		t.Error(err)
	}
}

func TestFetchUrls(t *testing.T) {
	t.Skip("@todo")
}
