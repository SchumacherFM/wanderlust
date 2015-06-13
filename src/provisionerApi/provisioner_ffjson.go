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

	"helpers"
)

// MarshalJSON implements encoding/json.Marshaler interface
func (mj *Config) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(1024)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (mj *Config) MarshalJSONBuf(buf *bytes.Buffer) error {
	var err error
	var obj []byte
	var first bool = true
	_ = obj
	_ = err
	_ = first
	buf.WriteString(`{`)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Icon":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Icon)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Name":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Name)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Url":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Url)
	buf.WriteString(`}`)
	return nil
}
