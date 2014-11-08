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
	"github.com/SchumacherFM/wanderlust/picnicApi"
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"strings"
)

func GetProvisioner() *provisionerApi.Config {
	sitemap := &sm{
		myRoute: "sitemap",
		config:  []string{"SiteMapUrl1", "SiteMapUrl2"}, // used in the html input field names
	}
	p := provisionerApi.NewProvisioner("Sitemap", "fa-sitemap", sitemap)
	return p
}

type (
	sm struct {
		// myRoute is the public name for the resource access
		myRoute string
		// config contains all the input field names which are use in the HTML partials
		config []string
	}
)

var (
	ErrValidate = errors.New("Invalid sitemap URL")
)

func (s *sm) Route() string {
	return s.myRoute
}

func (s *sm) FormHandler() picnicApi.HandlerFunc {
	return provisionerApi.FormGenerate(s.Route(), s.config)
}

// IsValid checks if the Value of PostData is valid sitemap URL
func (s *sm) IsValid(p *provisionerApi.PostData) error {
	return isValidSitemapUrl(p.Value)
}

// isValid checks if the Value of PostData is valid sitemap URL
func isValidSitemapUrl(v string) error {
	if "" == v {
		return nil
	}

	val := strings.ToLower(v)

	if false == strings.HasPrefix(val, "http") {
		return ErrValidate
	}

	if false == strings.HasSuffix(val, ".xml") && false == strings.HasSuffix(val, ".xml.gz") {
		return ErrValidate
	}
	return nil
}

func valueModifier(pd *provisionerApi.PostData) []byte {
	return []byte(strings.TrimSpace(pd.Value))
}

// https://restful-api-design.readthedocs.org/en/latest/methods.html#standard-methods
func (s *sm) SaveHandler() picnicApi.HandlerFunc {
	return provisionerApi.FormSave(s, valueModifier)
}
