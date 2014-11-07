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

package sitemap

import (
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"testing"
)

func TestIsValid(t *testing.T) {
	p := GetProvisioner()

	expected := map[string]error{
		"": nil,
		" http://www.golang.org/sitemap.xml":                                   nil,
		"http://www.amazon.de/sitemap-manual-index.xml":                        nil,
		"http://www.amazon.de/sitemaps.f9053414d236e84.SitemapIndex_1.xml.gz":  nil,
		"htstp://www.amazon.de/sitemaps.f3051414d236e84.SitemapIndex_2.xml.gz": ErrValidate,
		"http://www.amazon.de/sitemaps.f30v3414d236e84.SitemapIndex_3.html":    ErrValidate,
		"http://www.amazon.de/sitemaps.f30g3414d236e84.SitemapIndex_4.html.gz": ErrValidate,
		"http://www.amazon.de/sitemaps.f30g3414d236e84.SitemapIndex_4xml.gz":   ErrValidate,
	}

	pd := &provisionerApi.PostData{}
	for url, eErr := range expected {
		pd.Value = url
		actual := p.Api.IsValid(pd)
		if eErr != actual {
			t.Errorf("\nExpected: %#v\nActual:\t%#v\nURL: %s\n", eErr, actual, url)
		}

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
