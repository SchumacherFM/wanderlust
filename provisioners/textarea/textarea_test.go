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

func TestPrepareSave(t *testing.T) {
	p := GetProvisioner()
	var ret []byte
	pd := &provisionerApi.PostData{}

	type eData struct {
		in  string
		out string
		err error
	}

	var ed = make([]eData, 7)
	ed[0].in = ""
	ed[0].out = ""
	ed[0].err = nil
	ed[1].in = " http://www.golang.org/tour "
	ed[1].out = "http://www.golang.org/tour"
	ed[1].err = nil
	ed[2].in = "htp://www.golang.org/siteap.xml"
	ed[2].out = "htp://www.golang.org/siteap.xml"
	ed[2].err = nil
	ed[3].in = "hTtp://www.golang.org/siteMap.xml \nhttps://www.google.com/search?q=golang\n"
	ed[3].out = "hTtp://www.golang.org/siteMap.xml\nhttps://www.google.com/search?q=golang"
	ed[3].err = nil
	ed[4].in = strings.Repeat("hTtp://www.golang.org/siteMap.xml\n", 21)
	ed[4].out = ""
	ed[4].err = ErrTooManyURLs
	ed[5].in = "htp://www.golang.org/siteMap.xml \nhttps://www.google.com/search?q=golang\n"
	ed[5].out = ""
	ed[5].err = ErrValidate
	ed[6].in = "http://www.golang.org/siteMap.xml \n\n\nhttps://www.google.com/search?q=golang\n\n\nhttp://localhost"
	ed[6].out = "http://www.golang.org/siteMap.xml\nhttps://www.google.com/search?q=golang\nhttp://localhost"
	ed[6].err = nil
	ed[7].in = strings.Repeat("hTtp://www.golang.org/siteMap.xml\n", 18)
	ed[7].out = "hTtp://www.golang.org/siteMap.xml"
	ed[7].err = ErrTooManyURLs

	_, err := p.Api.PrepareSave(pd) // must succeed
	if nil != err {
		t.Errorf("%#v is valid 31!", pd)
	}

	pd.Value = " http://www.golang.org/tour"
	ret, err = p.Api.PrepareSave(pd) // must succeed
	if nil != err {
		t.Errorf("%#v is not valid 37!", pd)
	}
	if "http://www.golang.org/tour" != string(ret) {
		t.Errorf("Expected: http://www.golang.org/tour Got: %s", ret)
	}

	pd.Value = "htp://www.golang.org/siteap.xml"
	_, err = p.Api.PrepareSave(pd) // must fail
	if nil == err {
		t.Errorf("%#v is not valid 43!", pd)
	}

	pd.Value = "hTtp://www.golang.org/siteMap.xml\nhttps://www.google.com/search?q=golang\n"
	_, err = p.Api.PrepareSave(pd) // must succeed
	if nil != err {
		t.Errorf("%#v must be valid 49!", pd)
	}

	pd.Value = "hTtp://www.golang.org/siteMap.xml\n\nhttps://www.google.com/search?q=golang\t\n\n"
	_, err = p.Api.PrepareSave(pd) // must succeed
	if nil != err {
		t.Errorf("%#v must be valid 49!", pd)
	}

	pd.Value = strings.Repeat("hTtp://www.golang.org/siteMap.xml\n", 21)
	_, err = p.Api.PrepareSave(pd) // must fail
	if nil == err {
		t.Errorf("%#v must invalid 57!", pd)
	}

}

// MBA Mid 2012 1.8 GHz Intel Core i5
// BenchmarkPrepareSave	 1.000.000	      1520 ns/op
func BenchmarkPrepareSave(b *testing.B) {
	p := GetProvisioner()
	pd := &provisionerApi.PostData{}
	pd.Value = "hTtp://www.golang.org/siteMap.xml"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		p.Api.PrepareSave(pd) // must succeed
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
