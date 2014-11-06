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

// this package acts also as an external api to create custom provisioners
// and also to avoid import cycles
package api

import (
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/picnicApi"
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
		// DeleteHandler removes data from the rucksack and clears everything out
		//		DeleteHandler() HandlerFunc

		IsValid(p *PostData) error
	}

	// Implements encoding/json.Marshaler interface
	Provisioner struct {
		// This name appears in the frontend
		Name string
		// REST path
		Url string
		// can be a fa-* icon or path to an image
		Icon string
		// internal handler
		Api ProvisionerApi
	}

	PostData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	// modifier func which returns the value for saving in the rucksack
	ValueCallBack func(pd *PostData) []byte

	formData struct {
		Data []string `json:"data"`
	}
)

func NewProvisioner(n, i string, a ProvisionerApi) *Provisioner {
	return &Provisioner{
		Name: n,
		Url:  "/" + UrlRoutePrefix + "/" + a.Route(),
		Icon: i,
		Api:  a,
	}
}

// FormGenerate prepares the JSON object for AngularJS to fill the input fields of partial
// with its values. The input field names are hardcoded in the html partial as we won't to avoid
// dynamic rendered forms ...
func FormGenerate(dbName string, config []string) picnicApi.HandlerFunc {
	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {
		var jsonData = make([]string, 2*len(config))
		var i = 0
		for _, c := range config {
			jsonData[i] = c
			cVal, _ := rc.Backpacker().FindOne(dbName, c)
			jsonData[i+1] = string(cVal)
			i = i + 2
		}
		fd := &formData{
			Data: jsonData,
		}
		return helpers.RenderJSON(w, fd, 200)
	}
}

func FormSave(p ProvisionerApi, cb ValueCallBack) picnicApi.HandlerFunc {
	return func(rc picnicApi.RequestContextIf, w http.ResponseWriter, r *http.Request) error {
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

		err = rc.Backpacker().Insert(p.Route(), pd.Key, cb(pd))
		pd = nil
		if nil != err {
			status = http.StatusBadRequest
		}
		return helpers.RenderString(w, status, "")
	}
}
