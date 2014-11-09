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

// provisionerApi acts also as an external Api to create custom provisioners
// and also to avoid import cycles
package provisionerApi

import (
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/picnicApi"
	"github.com/SchumacherFM/wanderlust/rucksack"
	"net/http"
)

const (
	// no slashes for this PrePath in the route
	UrlRoutePrefix string = "provisioners"
)

type (
	// this methods will be used to query the provisioner instance and set values
	ProvisionerApi interface {
		// Route returns the endpoint of the route and is also abused as the database name
		Route() string
		// FormHandler returns a JSON object with a key called data which contains a key/value object
		// key is the form field name/id and value the value
		FormHandler() picnicApi.HandlerFunc
		// SaveHandler saves the POSTed data from the input fields into the rucksack
		SaveHandler() picnicApi.HandlerFunc
		// IsValid validates the post data and returns different errors which will pop up on the front end
		IsValid(p *PostData) error

		// ConfigComplete checks if all config values
		// have been successfully entered by the user. if so brotzeit can start automatically fetching URLs
		ConfigComplete(rucksack.Backpacker) (bool, error)

		// FetchUrls process can take quite a long time or nanoseconds
		FetchUrls(rucksack.Backpacker) []string

		// idea: maybe for output
		// ProgressInfo() string
	}

	// Implements encoding/json.Marshaler interface
	Config struct {
		// This name appears in the frontend
		name string
		// REST path
		url string
		// can be a fa-* icon or path to an image
		icon string
		// internal handler
		Api ProvisionerApi
	}

	PostData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	// modifier func which returns the value for saving in the rucksack.Backpacker
	ValueCallBack func(pd *PostData) []byte

	// internal struct for returning JSON with its data slice
	formData struct {
		Data []string `json:"data"`
	}
)

// NewProvisioner returns a new Config with the API of them
func NewProvisioner(n, i string, a ProvisionerApi) *Config {
	return &Config{
		name: n,
		url:  "/" + UrlRoutePrefix + "/" + a.Route(),
		icon: i,
		Api:  a,
	}
}

// FormGenerate prepares the JSON object for AngularJS to fill the input fields of the HTML partials
// with its values. The input field names are hardcoded in the HTML partial as we would like to avoid
// dynamic rendered forms ...
func FormGenerate(dbName string, config []string) picnicApi.HandlerFunc {
	return func(c picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
		var jsonData = make([]string, 2*len(config))
		var i = 0
		for _, cfg := range config {
			jsonData[i] = cfg
			cVal, _ := c.Backpacker().FindOne(dbName, cfg)
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
func FormSave(p ProvisionerApi, cb ValueCallBack) picnicApi.HandlerFunc {
	return func(c picnicApi.Context, w http.ResponseWriter, r *http.Request) error {
		status := http.StatusOK

		pd := &PostData{}
		err := helpers.DecodeJSON(r, pd)
		if nil != err {
			he := &picnicApi.HttpError{
				Status:      http.StatusBadRequest,
				Description: err.Error(),
			}
			return he
		}

		if errV := p.IsValid(pd); nil != errV {
			he := &picnicApi.HttpError{
				Status:      http.StatusExpectationFailed,
				Description: errV.Error(),
			}
			return he
		}

		err = c.Backpacker().Insert(p.Route(), pd.Key, cb(pd))
		pd = nil
		if nil != err {
			status = http.StatusBadRequest
		}
		return helpers.RenderString(w, status, "")
	}
}
