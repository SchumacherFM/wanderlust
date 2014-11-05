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
	picnicApi "github.com/SchumacherFM/wanderlust/picnic/api"
	. "github.com/SchumacherFM/wanderlust/provisioners/api"
	"github.com/SchumacherFM/wanderlust/rucksack"
	"net/http"
)

func GetProvisioner() *Provisioner {
	sitemap := &sm{
		myRoute: "sitemap",
		Data:    &storage{},
	}
	p := NewProvisioner("Sitemap", "fa-sitemap", sitemap)
	return p
}

type (
	sm struct {
		myRoute string
		Data    *storage `json:"data"`
	}

	// define storage keys for global use and move that into api package
	// so then no any handler func will be needed !?
	storage struct {
		SiteMapUrl1 string // that field name is the same as the html input field name in a partial
		SiteMapUrl2 string
	}
)

func (s *sm) Route() string {
	return s.myRoute
}

func (s *sm) FormHandler() picnicApi.HandlerFunc {

	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {

		url, err := rc.Backpacker().FindOne(s.Route(), "SiteMapUrl1")
		if nil != err && rucksack.ErrBreadNotFound != err {
			return err
		}
		s.Data.SiteMapUrl1 = string(url)

		url2, err := rc.Backpacker().FindOne(s.Route(), "SiteMapUrl2")
		if nil != err && rucksack.ErrBreadNotFound != err {
			return err
		}
		s.Data.SiteMapUrl2 = string(url2)

		return helpers.RenderJSON(w, s, 200)
	}
}

// use this instead of the the SaveHandler()
//func (s *sm) IsValid(p *PostData) bool {
//
//}

// https://restful-api-design.readthedocs.org/en/latest/methods.html#standard-methods
func (s *sm) SaveHandler() picnicApi.HandlerFunc {
	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {
		status := http.StatusOK
		err := SavePostData(r, rc.Backpacker(), s.Route())
		if nil != err {
			status = http.StatusBadRequest
		}
		return helpers.RenderString(w, status, "")
	}
}

func (s *sm) DeleteHandler() picnicApi.HandlerFunc {
	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {
		return helpers.RenderString(w, 200, "[\"Deleted Data\"]")
	}
}
