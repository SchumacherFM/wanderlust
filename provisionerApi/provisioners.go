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

package provisionerApi

import (
	"bytes"
)

type (
	// Provisioners is a collection of configurations
	Provisioners struct {
		Col []*Config
	}

	ProvisionersJsonIf interface { // crap name and superfluous
		MarshalJSON() ([]byte, error)
		MarshalJSONBuf(buf *bytes.Buffer) error
	}

	Collectioner interface {
		Collection() []*Config
	}
)

var (
	_ Collectioner = &Provisioners{}
)

func (p *Provisioners) Add(prov *Config) {
	p.Col = append(p.Col, prov)
}

func (p *Provisioners) Collection() []*Config {
	return p.Col
}

func NewProvisioners() *Provisioners {
	return &Provisioners{}
}
