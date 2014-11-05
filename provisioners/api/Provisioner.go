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
	. "github.com/SchumacherFM/wanderlust/picnic/api"
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
		// GetRoutes returns the endpoint of the route
		Route() string
		// FormHandler returns a JSON object with a key called data which contains a key/value object
		// key is the form field name/id and value the value
		FormHandler() HandlerFunc
		// SaveHandler saves the POSTed data from the input fields into the rucksack
		SaveHandler() HandlerFunc
		// DeleteHandler removes data from the rucksack and clears everything out
		DeleteHandler() HandlerFunc

		// maybe more methods to add ...
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
)

func NewProvisioner(n, i string, a ProvisionerApi) *Provisioner {
	return &Provisioner{
		Name: n,
		Url:  "/" + UrlRoutePrefix + "/" + a.Route(),
		Icon: i,
		Api:  a,
	}
}

func SavePostData(r *http.Request, bp rucksack.Backpacker, dbName string) error {
	p := &PostData{}
	err := helpers.DecodeJSON(r, p)
	if nil != err {
		return err
	}
	return bp.Insert(dbName, p.Key, []byte(p.Value))
}
