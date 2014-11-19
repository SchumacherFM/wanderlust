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
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/rucksack"
)

const (
	// no slashes for this PrePath in the route
	UrlRoutePrefix string = "provisioners"
)

type (
	// https://restful-api-design.readthedocs.org/en/latest/methods.html#standard-methods

	// ColdCutter methods will be used to query the provisioner instance, set values
	ColdCutter interface {
		// Route returns the endpoint of the route and is also abused as the database name
		Route() string
		// FormHandler returns a JSON object with a key called data which contains a key/value object
		// key is the form field name/id and value the value
		Config() []string

		// PrepareSave prepares the POSTed data from the input fields for the saving the rucksack and also checks if
		// the entered data is valid
		PrepareSave(pd *PostData) ([]byte, error)

		// idea: maybe for output
		// ProgressInfo() string

		// ConfigComplete checks if all config values
		// have been successfully entered by the user. if so brotzeit can start automatically fetching URLs
		// This func will also be used to check if we can add FetchUrls to the crond. A user can have configured a
		// cron schedule but not fully provided all config details
		ConfigComplete(rucksack.Backpacker) (bool, error)

		// FetchURLs downloads all the URLs and stores them in the brotzeit DB.
		// Process can take quite a long time or nanoseconds.
		// @todo implement http://talks.golang.org/2013/advconc.slide
		FetchURLs(rucksack.Backpacker, *log.Logger) func()
	}

	// Implements encoding/json.Marshaler interface is mainly used for the route
	Config struct {
		// This name appears in the frontend
		Name string
		// REST path
		Url string
		// can be a fa-* icon or path to an image
		Icon string
		// internal handler
		Api ColdCutter
	}

	PostData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	// Base struct is for embedding/composition in other provisioner structs aka Parent Class ;-)
	Base struct {
		// route is the public name for the resource access and acts also as internal identifier
		TheRoute string
		// config contains all the input field names which are use in the HTML partials
		TheConfig []string
	}
)

// NewProvisioner returns a new Config with the API of them
func NewProvisioner(n, i string, a ColdCutter) *Config {
	return &Config{
		Name: n,
		Url:  "/" + UrlRoutePrefix + "/" + a.Route(),
		Icon: i,
		Api:  a,
	}
}

func (b *Base) Route() string {
	return b.TheRoute
}

func (b *Base) Config() []string {
	return b.TheConfig
}
