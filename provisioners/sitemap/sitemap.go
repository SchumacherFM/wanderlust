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
	"errors"
	"strings"

	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"github.com/SchumacherFM/wanderlust/rucksack"
	log "github.com/segmentio/go-log"
)

func GetProvisioner() *provisionerApi.Config {
	s := &sm{
		Base: provisionerApi.Base{
			TheRoute:  "sitemap",
			TheConfig: []string{"SiteMapUrl1", "SiteMapUrl2"}, // used in the html input field names
		},
	}
	p := provisionerApi.NewProvisioner("Sitemap", "fa-sitemap", s)
	return p
}

type (
	sm struct {
		provisionerApi.Base
		// other fields ...
	}
)

var (
	ErrValidate = errors.New("Invalid sitemap URL")
	// provide a static (compile time) check that sm satisfies the ColdCutter interface
	_ provisionerApi.ColdCutter = &sm{}
)

func (s *sm) PrepareSave(pd *provisionerApi.PostData) ([]byte, error) {
	if false == isValidSitemapUrl(pd.Value) {
		return nil, ErrValidate
	}
	return []byte(strings.TrimSpace(pd.Value)), nil
}

// ConfigComplete implements the brotzeit.Fetcher interface to check if all config values
// have been successfully entered by the user. if so brotzeit can start automatically fetching URLs
func (s *sm) ConfigComplete(bp rucksack.Backpacker) (bool, error) {
	sm1, err := bp.FindOne(s.Route(), "SiteMapUrl1")
	if nil != err {
		return false, err
	}
	sm2, err := bp.FindOne(s.Route(), "SiteMapUrl2")
	if nil != err {
		return false, err
	}
	return (len(sm1) > 5 && true == isValidSitemapUrl(string(sm1))) ||
		(len(sm2) > 5 && true == isValidSitemapUrl(string(sm2))), nil
}

// Fetch implements the brotzeit.Fetcher interface
func (s *sm) FetchURLs(bp rucksack.Backpacker, l *log.Logger) func() {
	return func() {
		l.Debug("SM CRON %s", s.Route())
	}
}

// isValid checks if the Value of PostData is valid sitemap URL
func isValidSitemapUrl(v string) bool {
	if "" == v {
		return true
	}
	val := strings.ToLower(v)
	if false == strings.HasPrefix(val, "http") {
		return false
	}
	if false == strings.HasSuffix(val, ".xml") && false == strings.HasSuffix(val, ".xml.gz") {
		return false
	}
	return true
}

func isValidUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
