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

package provisioners

import (
	"errors"
	"net/http"

	"helpers"
	"picnicApi"
	"provisionerApi"
	"provisioners/sitemap"
	"provisioners/textarea"
)

var (
	provisionerCollection = provisionerApi.NewProvisioners()
	ErrCollectionEmpty    = errors.New("Provisioner Collection is empty")
)

type (
	// internal struct for returning JSON with its data slice
	formData struct {
		Data []string `json:"data"`
	}
)

// initializes all the build-in provisioners, every custom provisioner will be added via build tag
func init() {
	AddProvisioner(sitemap.GetProvisioner())
	AddProvisioner(textarea.GetProvisioner())
}

func AddProvisioner(p *provisionerApi.Config) {
	provisionerCollection.Add(p)
}

func GetAvailable() (*provisionerApi.Provisioners, error) {
	return provisionerCollection, nil
}

func GetRoutePathPrefix() string {
	return provisionerApi.UrlRoutePrefix
}

// FormGenerate prepares the JSON object for AngularJS to fill the input fields of the HTML partials
// with its values. The input field names are hardcoded in the HTML partial as we would like to avoid
// dynamic rendered forms ...
func FormGenerate(p provisionerApi.ColdCutter) picnicApi.HandlerFunc {
	return func(c picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
		var jsonData = make([]string, 2*len(p.Config()))
		var i = 0
		for _, cfg := range p.Config() {
			jsonData[i] = cfg
			cVal, _ := c.Backpacker().FindOne(p.Route(), cfg)
			jsonData[i+1] = string(cVal)
			i = i + 2
		}
		fd := &formData{
			Data: jsonData,
		}
		return helpers.RenderJSON(w, fd, 200)
	}
}

// FormSave saves the key/value pair in the rucksack.Backpacker. Does also some validation provided by the
// Config API
func FormSave(p provisionerApi.ColdCutter) picnicApi.HandlerFunc {
	return func(c picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
		status := http.StatusOK

		pd := &provisionerApi.PostData{}
		err := helpers.DecodeJSON(r, pd)
		if nil != err {
			he := &picnicApi.HttpError{
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			}
			return he
		}

		data, errV := p.PrepareSave(pd)
		if nil != errV {
			he := &picnicApi.HttpError{
				Status:      http.StatusExpectationFailed,
				Description: errV.Error(),
			}
			return he
		}

		err = c.Backpacker().Insert(p.Route(), pd.Key, data)
		pd = nil
		if nil != err {
			status = http.StatusBadRequest
		}
		return helpers.RenderString(w, status, "")
	}
}
