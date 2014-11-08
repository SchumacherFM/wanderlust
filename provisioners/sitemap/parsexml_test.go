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
	"github.com/SchumacherFM/wanderlust/helpers"
	"strings"
	"testing"
)

// run with $ go test -v --bench=. -test.benchmem .

type smc struct {
	isSiteMapIndex bool
	isSiteMap      bool
	loc            int
	data           string
}

var sitemapCollection = []smc{
	smc{
		isSiteMapIndex: true,
		isSiteMap:      false,
		loc:            0,
		data: `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"/>`,
	},
	smc{
		isSiteMapIndex: false,
		isSiteMap:      false,
		loc:            0,
		data: `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`,
	},
	smc{
		isSiteMapIndex: true,
		isSiteMap:      false,
		loc:            11,
		data: `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
<sitemap>
		<loc>http://www.golang.com/de-de/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/de-fr/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/de-en/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/de-it/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/ch-de/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/ch-fr/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/ch-en/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/ch-it/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/it-de/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/it-it/sitemap.xml</loc>
</sitemap>
<sitemap>
		<loc>http://www.golang.com/fr-fr/sitemap.xml</loc>
</sitemap>
</sitemapindex>`,
	},
	smc{
		isSiteMapIndex: false,
		isSiteMap:      true,
		loc:            36,
		data: `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
	<url>
			<loc>http://www.golang.com/de-de/damen</loc>
			<xhtml:link rel="alternate" hreflang="de" href="http://www.golang.com/de-de/damen"/>
			<xhtml:link rel="alternate" hreflang="de-ch" href="http://www.golang.com/ch-de/damen"/>
			<xhtml:link rel="alternate" hreflang="fr" href="http://www.golang.com/fr-fr/femme"/>
			<xhtml:link rel="alternate" hreflang="it" href="http://www.golang.com/it-it/donna"/>
			<xhtml:link rel="alternate" hreflang="de" href="http://www.golang.com/it-de/damen"/>
			<xhtml:link rel="alternate" hreflang="fr" href="http://www.golang.com/de-fr/femme"/>
			<xhtml:link rel="alternate" hreflang="fr-ch" href="http://www.golang.com/ch-fr/femme"/>
			<xhtml:link rel="alternate" hreflang="it" href="http://www.golang.com/de-it/donna"/>
			<xhtml:link rel="alternate" hreflang="it-ch" href="http://www.golang.com/ch-it/donna"/>
			<xhtml:link rel="alternate" hreflang="en" href="http://www.golang.com/de-en/women"/>
			<xhtml:link rel="alternate" hreflang="en-ch" href="http://www.golang.com/ch-en/women"/>
			<lastmod>2014-11-08</lastmod>
			<changefreq>weekly</changefreq>
			<priority>0.5</priority>
	</url>
	<url>
			<loc>http://www.golang.com/de-de/damen/bhs</loc>
			<xhtml:link rel="alternate" hreflang="de" href="http://www.golang.com/de-de/damen/bhs"/>
			<xhtml:link rel="alternate" hreflang="de-ch" href="http://www.golang.com/ch-de/damen/bhs"/>
			<xhtml:link rel="alternate" hreflang="fr" href="http://www.golang.com/fr-fr/femme/soutiens-gorge"/>
			<xhtml:link rel="alternate" hreflang="it" href="http://www.golang.com/it-it/donna/reggiseni"/>
			<xhtml:link rel="alternate" hreflang="de" href="http://www.golang.com/it-de/damen/bhs"/>
			<xhtml:link rel="alternate" hreflang="fr" href="http://www.golang.com/de-fr/femme/soutiens-gorge"/>
			<xhtml:link rel="alternate" hreflang="fr-ch" href="http://www.golang.com/ch-fr/femme/soutiens-gorge"/>
			<xhtml:link rel="alternate" hreflang="it" href="http://www.golang.com/de-it/donna/reggiseni"/>
			<xhtml:link rel="alternate" hreflang="it-ch" href="http://www.golang.com/ch-it/donna/reggiseni"/>
			<xhtml:link rel="alternate" hreflang="en" href="http://www.golang.com/de-en/women/bras"/>
			<xhtml:link rel="alternate" hreflang="en-ch" href="http://www.golang.com/ch-en/women/bras"/>
			<lastmod>2014-11-08</lastmod>
			<changefreq>weekly</changefreq>
			<priority>0.5</priority>
	</url>
    <url>
        <loc>http://www.wltest.com/ch-it/elastic-02138</loc>
        <xhtml:link rel="alternate" hreflang="de" href="http://www.wltest.com/de-de/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="de-ch" href="http://www.wltest.com/ch-de/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="fr" href="http://www.wltest.com/fr-fr/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="it" href="http://www.wltest.com/it-it/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="de" href="http://www.wltest.com/it-de/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="fr" href="http://www.wltest.com/de-fr/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="fr-ch" href="http://www.wltest.com/ch-fr/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="it" href="http://www.wltest.com/de-it/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="it-ch" href="http://www.wltest.com/ch-it/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="en" href="http://www.wltest.com/de-en/elastic-02138"/>
        <xhtml:link rel="alternate" hreflang="en-ch" href="http://www.wltest.com/ch-en/elastic-02138"/>
        <lastmod>2014-11-08</lastmod>
        <changefreq>weekly</changefreq>
        <priority>1.0</priority>
    </url>
</urlset>`,
	},
}

func TestParseSiteMapIndex(t *testing.T) {
	for _, s := range sitemapCollection {
		si := helpers.NewReadCloser(s.data)
		sc, err := parseSiteMapIndex(si)
		if nil != err {
			t.Error(err)
		}
		if nil != sc && len(sc) > 0 && false == s.isSiteMapIndex {
			t.Errorf("Should be not a siteMapIndex\n%#v\n%#v", sc, s)
		}
		if true == s.isSiteMapIndex && s.loc != len(sc) {
			t.Errorf("\nExpected: %d\nActual: %d", s.loc, len(sc))
		}
		for _, sv := range sc {
			if false == strings.Contains(sv, ".xml") {
				t.Error("Not a sitemap.xml", sv)
			}
		}
	}
}

// MBA Mid 2012 1.8 GHz Intel Core i5
// BenchmarkParseSiteMapIndex	   10.000	    172.082 ns/op with 30 maxUrls
// BenchmarkParseSiteMapIndex	    2000	    739977 ns/op	  823752 B/op	     351 allocs/op with 50.000 maxUrls
func BenchmarkParseSiteMapIndex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		si := helpers.NewReadCloser(sitemapCollection[2].data)
		parseSiteMapIndex(si)
	}
}

func TestParseSiteMap(t *testing.T) {

	for _, s := range sitemapCollection {
		si := helpers.NewReadCloser(s.data)
		sm, err := parseSiteMap(si)
		if nil != err {
			t.Error(err)
		}

		if 0 == len(sm) && true == s.isSiteMap {
			t.Errorf("Should be not a siteMap\n%#v\n%#v\n", sm, s)
		}
		if true == s.isSiteMap && s.loc != len(sm) {
			t.Errorf("\nExpected: %d\nActual: %d", s.loc, len(sm))
		}
		for _, sLoc := range sm {
			if false == isValidUrl(sLoc) {
				t.Errorf("Invalid URL %s", sLoc)
			}
		}
	}
}

// MBA Mid 2012 1.8 GHz Intel Core i5
// BenchmarkParseSiteMap	    2000	   1115713 ns/op	  849352 B/op	    1101 allocs/op with 50.000 maxUrls
func BenchmarkParseSiteMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		si := helpers.NewReadCloser(sitemapCollection[3].data)
		parseSiteMap(si)
	}
}

func TestIsValidUrl(t *testing.T) {
	expected := map[string]bool{
		"http://golang.org":  true,
		"https://golang.org": true,
		"htps://golang.org":  false,
		"ftp://golang.org":   false,
		"onion://golang.org": false,
		"http:/golang.org":   false,
	}

	for url, res := range expected {
		if act := isValidUrl(url); res != act {
			t.Errorf("Expected: %t got %t for %s ", res, act, url)
		}
	}
}
