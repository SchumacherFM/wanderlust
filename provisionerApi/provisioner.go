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
	"github.com/SchumacherFM/wanderlust/rucksack"
)

const (
	// no slashes for this PrePath in the route
	UrlRoutePrefix string = "provisioners"
)

type (
	// https://restful-api-design.readthedocs.org/en/latest/methods.html#standard-methods

	// ColdCuts methods will be used to query the provisioner instance, set values
	ColdCuts interface {
		// Route returns the endpoint of the route and is also abused as the database name
		Route() string
		// FormHandler returns a JSON object with a key called data which contains a key/value object
		// key is the form field name/id and value the value
		Config() []string

		// PrepareSave prepares the POSTed data from the input fields for the saving the rucksack and also checks if
		// the entered data is valid
		PrepareSave(pd *PostData) ([]byte, error)

		// ConfigComplete checks if all config values
		// have been successfully entered by the user. if so brotzeit can start automatically fetching URLs
		ConfigComplete(rucksack.Backpacker) (bool, error)

		// FetchUrls process can take quite a long time or nanoseconds
		FetchUrls(rucksack.Backpacker) []string

		// idea: maybe for output
		// ProgressInfo() string
	}

	// Implements encoding/json.Marshaler interface is mainly used for the route
	Config struct {
		// This name appears in the frontend
		name string
		// REST path
		url string
		// can be a fa-* icon or path to an image
		icon string
		// internal handler
		Api ColdCuts
	}

	PostData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	// Base struct is for embedding/composition in other provisioner structs aka Parent Class ;-)
	Base struct {
		// route is the public name for the resource access
		TheRoute string
		// config contains all the input field names which are use in the HTML partials
		TheConfig []string
	}
)

// NewProvisioner returns a new Config with the API of them
func NewProvisioner(n, i string, a ColdCuts) *Config {
	return &Config{
		name: n,
		url:  "/" + UrlRoutePrefix + "/" + a.Route(),
		icon: i,
		Api:  a,
	}
}

func (b *Base) Route() string {
	return b.TheRoute
}

func (b *Base) Config() []string {
	return b.TheConfig
}
