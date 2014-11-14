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
	"sort"
	"strings"
)

const (
	// http://en.wikipedia.org/wiki/Sitemaps
	maxUrlsPerSitemap int    = 50000
	sortKeySeparator  string = "€"
)

var (
	// @todo investigate if a []string is faster instead of a map ...
	changefreqMapper = map[string]string{
		"always":  "z",
		"hourly":  "y",
		"daily":   "x",
		"weekly":  "w",
		"monthly": "v",
		"yearly":  "u",
		"never":   "t",
		"":        "s",
	}
)

// UrlNode is used in the xml decode tokenizer
type UrlNode struct {
	Loc        string      `xml:"loc"`
	Lastmod    string      `xml:"lastmod"`
	Changefreq string      `xml:"changefreq"`
	Priority   string      `xml:"priority"`
	XhtmlLink  []XhtmlLink `xml:"http://www.w3.org/1999/xhtml link"`
}

type XhtmlLink struct {
	Rel      string `xml:"rel,attr"`
	Hreflang string `xml:"hreflang,attr"`
	Href     string `xml:"href,attr"`
}

// sortKey generates the sort key if Priority, Lastmod & Changefreq aren't empty
func (u *UrlNode) sortKey(url string) string {
	cf, _ := changefreqMapper[u.Changefreq]
	if "" == url {
		url = u.Loc
	}
	if "" == u.Priority && "" == u.Lastmod && "" == u.Changefreq {
		return url
	}
	return u.Priority + u.Lastmod + cf + sortKeySeparator + url
}

// sortKey generates the sort key if Priority, Lastmod & Changefreq aren't empty
func (x *XhtmlLink) sortKey(u *UrlNode) string {
	return u.sortKey(x.Href)
}

// parseSiteMap parses sitemapindex and sitemap XML files. Returns a slice with all available
// URLs in reverse sorted order. Sort criteria is: Priority, LastMod, ChangeFreqeuency and the URL.
// The XML is parsed continuously and not all at once.
func parseSiteMap(r io.ReadCloser) ([]string, bool, error) {
	maxUrls := make([]string, maxUrlsPerSitemap)
	urlCount := 0
	totalErr := 0
	isSiteMapIndex := false
	decoder := xml.NewDecoder(r)
	defer r.Close()

	var inElement string

	for {
		// Read tokens from the XML document in a stream.
		t, dtErr := decoder.Token()
		if t == nil {
			break
		}
		if nil != dtErr {
			return nil, isSiteMapIndex, dtErr
		}

		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "url"
			if "sitemap" == inElement || "url" == inElement {
				var un UrlNode
				// decode a whole chunk of following XML into the
				// variable un which is a UrlNode (see above); decErr will be ignored ...
				decErr := decoder.DecodeElement(&un, &se)
				if nil != decErr {
					totalErr++
				}
				isSiteMapUrl := isValidSitemapUrl(un.Loc)
				if true == isSiteMapUrl && false == isSiteMapIndex {
					isSiteMapIndex = true
				}
				if true == isSiteMapUrl || true == isValidUrl(un.Loc) {
					maxUrls[urlCount] = un.sortKey("")
					urlCount++
				}
				if urlCount >= maxUrlsPerSitemap {
					return sortUrls(maxUrls, isSiteMapIndex), isSiteMapIndex, nil
				}
				// the following if block is only for valid endpoints
				if nil != un.XhtmlLink && len(un.XhtmlLink) > 0 {
					for _, xHref := range un.XhtmlLink {
						if true == isValidUrl(un.Loc) {
							maxUrls[urlCount] = xHref.sortKey(&un)
							urlCount++
						}
						if urlCount >= maxUrlsPerSitemap {
							return sortUrls(maxUrls, isSiteMapIndex), isSiteMapIndex, nil
						}
					}
				}
			}
		default:
		}
	}

	// shrink the previous created slice and free memory for maxUrls
	urls := make([]string, urlCount)
	copy(urls, maxUrls[:urlCount])
	maxUrls = nil
	return sortUrls(urls, isSiteMapIndex), isSiteMapIndex, nil
}

// sortUrls sorts a string slice in reverse order and removes the sort index prefix
func sortUrls(urls []string, isSiteMapIndex bool) []string {
	if true == isSiteMapIndex {
		return urls
	}
	sl := len(sortKeySeparator)
	sort.Sort(sort.Reverse(sort.StringSlice(urls)))
	for i, u := range urls {
		cut := strings.Index(u, sortKeySeparator) + sl
		if cut > sl {
			urls[i] = u[cut:]
		}
	}
	return urls
}
