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
	. "github.com/SchumacherFM/wanderlust/provisioners/api"
	"strings"
)

func GetProvisioner() *Provisioner {
	sitemap := &sm{
		myRoute: "sitemap",
		config:  []string{"SiteMapUrl1", "SiteMapUrl2"}, // used in the html input field names
	}
	p := NewProvisioner("Sitemap", "fa-sitemap", sitemap)
	return p
}

type (
	sm struct {
		myRoute string
		config  []string
	}
)

var (
	ErrValidate = errors.New("Failed to validate the value")
)

func (s *sm) Route() string {
	return s.myRoute
}

func (s *sm) FormHandler() picnicApi.HandlerFunc {
	return FormGenerate(s.Route(), s.config)
}

// use this instead of the the SaveHandler()
func (s *sm) IsValid(p *PostData) error {

	return ErrValidate

	if "" == p.Value {
		return nil
	}

	val := strings.ToLower(strings.TrimSpace(p.Value))

	if false == strings.HasPrefix(val, "http") {
		return ErrValidate
	}

	if false == strings.HasSuffix(val, "/sitemap.xml") {
		return ErrValidate
	}
	return nil
}

func valueModifier(pd *PostData) []byte {
	return []byte(strings.TrimSpace(pd.Value))
}

// https://restful-api-design.readthedocs.org/en/latest/methods.html#standard-methods
func (s *sm) SaveHandler() picnicApi.HandlerFunc {
	return FormSave(s, valueModifier)
}
