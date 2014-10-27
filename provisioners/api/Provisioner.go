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
package api

import (
	. "github.com/SchumacherFM/wanderlust/picnic/api"
)

const (
	// no slashes for this PrePath in the route
	URL_PRE_ROUTE string = "provisioners"
)

type (
	// this methods will be used to query the provisioner instance and set values
	ProvisionerApi interface {
		// GetRoutes returns the endpoint of the route
		Route() string
		// GetRouteHandler returns a handler which must manage GET, POST and DELETE methods
		// GET returns the config for the <form> fields and also the saved data for the inputs
		// POST saves the data from the input field into the rucksack
		RouteHandler() HandlerFunc

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
)

func NewProvisioner(n, i string, a ProvisionerApi) *Provisioner {
	return &Provisioner{
		Name: n,
		Url:  "/" + URL_PRE_ROUTE + "/" + a.Route(),
		Icon: i,
		Api:  a,
	}
}
