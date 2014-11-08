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
	"encoding/xml"
	"io"
	"io/ioutil"
)

type SitemapIndex struct {
	Sitemap []UrlNode `xml:"sitemap"`
}

//
//type ItemNode struct{
//	Loc       string   `xml:"loc"`
//	Lastmod   string `xml:"lastmod"`
//}
//type Result struct {
//	XMLName  xml.Name `xml:"sitemapindex"`
//	ItemNode []ItemNode `xml:"sitemap"`
//}
//
type UrlNode struct {
	Loc        string `xml:"loc"`
	Lastmod    string `xml:"lastmod"`
	Changefreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

//type ItemDetail struct {
//
//	XMLName     xml.Name `xml:"urlset"`
//	UrlNodeList [] UrlNode `xml:"url"`
//}

// r is equal to res, err := http.Get(url)
func parseSiteMapIndex(r io.ReadCloser) ([]string, error) {
	data, err := ioutil.ReadAll(r)
	defer r.Close()

	if nil != err {
		return nil, err
	}

	si := &SitemapIndex{}
	xml.Unmarshal(data, si)
	data = nil

	if 0 == len(si.Sitemap) {
		return nil, nil
	}

	ret := make([]string, len(si.Sitemap))
	for k, un := range si.Sitemap {
		if nil == isValid(un.Loc) {
			ret[k] = un.Loc
		}
	}

	return ret, nil
}
