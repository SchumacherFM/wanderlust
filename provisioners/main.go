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
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	. "github.com/SchumacherFM/wanderlust/provisioners/api"
	"github.com/SchumacherFM/wanderlust/provisioners/sitemap"
	"github.com/SchumacherFM/wanderlust/provisioners/textarea"
)

var (
	provisionerCollection = NewProvisioners()
	ErrCollectionEmpty    = errgo.New("Provisioner Collection is empty")
)

// initializes all the build-in provisioners, every custom provisioner will be added via build tag
func init() {
	sm := sitemap.GetProvisioner()
	AddProvisioner(sm)
	ta := textarea.GetProvisioner()
	AddProvisioner(ta)
}

func AddProvisioner(p *Provisioner) {
	provisionerCollection.Add(p)
}

func GetAvailable() (*Provisioners, error) {
	if 0 == provisionerCollection.Length() {
		return nil, ErrCollectionEmpty
	}
	return provisionerCollection, nil
}

func GetRoutePathPrefix() string {
	return URL_PRE_ROUTE
}
