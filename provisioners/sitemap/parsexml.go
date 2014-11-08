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
)

const (
	// http://en.wikipedia.org/wiki/Sitemaps
	maxUrlsPerSitemap int = 50000
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

// r is equal to res, err := http.Get(url)
func parseSiteMapIndex(r io.ReadCloser) ([]string, error) {
	return privateParser(r, "sitemap", isValidSitemapUrl)
}

// Use a sitemap to indicate alternate language pages
// https://support.google.com/webmasters/answer/2620865?hl=en
func parseSiteMap(r io.ReadCloser) ([]string, error) {
	return privateParser(r, "url", isValidUrl)
}

// privateParser parses sitemapindex and sitemap XML files. The file size does not matter as this
// parser can parse TB huge files with using only ~5MB of memory
func privateParser(r io.ReadCloser, startElement string, urlValidator func(string) bool) ([]string, error) {
	maxUrls := make([]string, maxUrlsPerSitemap)
	urlCount := 0
	totalErr := 0

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
			return nil, dtErr
		}

		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "url"
			if inElement == startElement {
				var un UrlNode
				// decode a whole chunk of following XML into the
				// variable un which is a UrlNode (see above); decErr will be ignored ...
				decErr := decoder.DecodeElement(&un, &se)
				if nil != decErr {
					totalErr++
				}
				if true == urlValidator(un.Loc) {
					maxUrls[urlCount] = un.Loc
					urlCount++
				}
				if urlCount >= maxUrlsPerSitemap {
					return maxUrls, nil
				}
				// the following if block is only for valid endpoints
				if nil != un.XhtmlLink && len(un.XhtmlLink) > 0 {
					for _, xHref := range un.XhtmlLink {
						if true == isValidUrl(un.Loc) {
							maxUrls[urlCount] = xHref.Href
							urlCount++
						}
						if urlCount >= maxUrlsPerSitemap {
							return maxUrls, nil
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
	return urls, nil
}
