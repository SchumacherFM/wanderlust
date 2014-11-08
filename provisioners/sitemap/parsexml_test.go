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

var siteMapIndexCollection = []string{
	`<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"/>`,
	`<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	`,
	`<?xml version="1.0" encoding="UTF-8"?>
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
	`<?xml version="1.0" encoding="UTF-8"?>
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
</urlset>
	`,
}

func TestParseSiteMapIndex(t *testing.T) {

	for _, sitemapIndex := range siteMapIndexCollection {
		si := helpers.NewReadCloser(sitemapIndex)
		sc, err := parseSiteMapIndex(si)
		if nil != err {
			t.Error(err)
		}
		if 11 != len(sc) {
			t.Errorf("\nExpected: 11\nActual: %d", len(sc))
		}
		for _, sv := range sc {
			if false == strings.Contains(sv, ".xml") {
				t.Error("Not a sitemap.xml", sv)
			}
		}
	}
}

// MBA Mid 2012 1.8 GHz Intel Core i5
// BenchmarkParseSiteMapIndex	   10000	    189198 ns/op
//func BenchmarkParseSiteMapIndex(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		si := helpers.NewReadCloser(mockSitemapIndex)
//		parseSiteMapIndex(si)
//	}
//}
