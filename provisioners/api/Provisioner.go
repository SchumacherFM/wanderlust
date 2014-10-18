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

// this package acts also as an external api to create
package api

import "bytes"

const (
	// no slashes for this PrePath in the route
	URL_PRE_ROUTE string = "provisioners"
)

type (
	// this methods will be used to query the provisioner instance and set values
	ProvisionerMethod interface {
		MethodA()
		MethodB()
	}

	Provisioner struct {
		// This name appears in the frontend
		Name string
		// REST path
		Url string
		// can be a fa-* icon or path to an image
		Icon string
		// internal handler
		Method ProvisionerMethod
	}

	ProvisionerJsonIf interface {
		MarshalJSON() ([]byte, error)
		MarshalJSONBuf(buf *bytes.Buffer) error
	}
)

func NewProvisioner(n, u, i string, m ProvisionerMethod) *Provisioner {
	return &Provisioner{
		Name:   n,
		Url:    "/" + URL_PRE_ROUTE + u,
		Icon:   i,
		Method: m,
	}
}
