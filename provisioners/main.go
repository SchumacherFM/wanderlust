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

const (
	// no slashes for this PrePath in the route
	URL_PRE_ROUTE string = "provisioners"
)

type (
	Provisioner struct {
		// This name appears in the frontend
		Name string
		// REST path
		Url string
		// can be a fa-* icon or path to an image
		Icon string
		// internal handler
		method func() error
	}

	Provisioners struct {
		Collection []*Provisioner
	}
)

var (
	provisionerCollection *Provisioners
)

func init() {
	provisionerCollection = &Provisioners{}
}

func (p *Provisioners) add(prov *Provisioner) {
	p.Collection = append(p.Collection, prov)
}

func GetAvailable() (*Provisioners, error) {
	return provisionerCollection, nil
}

// @todo remove this as temp implementation
func GetRoutePathPrefix() string {
	return "/provisioners"
}
