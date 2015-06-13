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

// Code generated by ffjson <https://github.com/pquerna/ffjson>
// source: brotzeit.go

package brotzeit

import (
	"bytes"

	"helpers"
)

func (mj *BzConfig) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(336)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *BzConfig) MarshalJSONBuf(buf *bytes.Buffer) error {
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
	buf.WriteString(`"Route":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Route)
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
	buf.WriteString(`"Schedule":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Schedule)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"UrlCount":`)
	helpers.Ffjson_FormatBits(buf, uint64(mj.UrlCount), 10, mj.UrlCount < 0)
	buf.WriteString(`}`)
	return nil
}

func (mj *BzConfigs) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(144)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *BzConfigs) MarshalJSONBuf(buf *bytes.Buffer) error {
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
	buf.WriteString(`"Collection":`)
	if mj.Collection != nil {
		buf.WriteString(`[`)
		for i, v := range mj.Collection {
			if i != 0 {
				buf.WriteString(`,`)
			}
			if v != nil {
				err = v.MarshalJSONBuf(buf)
				if err != nil {
					return err
				}
			} else {
				buf.WriteString(`null`)
			}
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	buf.WriteString(`}`)
	return nil
}