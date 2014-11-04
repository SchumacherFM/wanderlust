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

	storage struct {
		SiteMapUrl string
	}

	// @todo urgent temp implementation. needs to be more global.
	postData struct {
		Key, Prov, Value string
	}
)

func (s *sm) Route() string {
	return s.myRoute
}

func (s *sm) FormHandler() picnicApi.HandlerFunc {

	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {

		// @todo also that App().Backpacker() is way tooooooo long
		url, err := rc.App().Backpacker().FindOne(s.Route(), "SiteMapUrl")
		if nil != err && rucksack.ErrBreadNotFound != err {
			return err
		}
		s.Data.SiteMapUrl = string(url)
		return helpers.RenderJSON(w, s, 200)
	}
}

// https://restful-api-design.readthedocs.org/en/latest/methods.html#standard-methods
func (s *sm) SaveHandler() picnicApi.HandlerFunc {
	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {
		// status 200 is ok, and
		//		status := http.StatusBadRequest

		jsonData := &postData{}
		err := helpers.DecodeJSON(r, jsonData)
		if nil != err {
			return err
		}

		rc.App().Backpacker().Insert(s.Route(), jsonData.Key, []byte(jsonData.Value))
		status := http.StatusOK
		return helpers.RenderString(w, status, "")
	}
}

func (s *sm) DeleteHandler() picnicApi.HandlerFunc {
	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {
		return helpers.RenderString(w, 200, "[\"Deleted Data\"]")
	}
}
