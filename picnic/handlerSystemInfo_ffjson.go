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
// source: handlerSystemInfo.go

package picnic

import (
	"bytes"

	"github.com/SchumacherFM/wanderlust/helpers"
)

func (mj *SystemInfo) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(1024)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *SystemInfo) MarshalJSONBuf(buf *bytes.Buffer) error {
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
	buf.WriteString(`"Brotzeit":`)
	helpers.Ffjson_FormatBits(buf, uint64(mj.Brotzeit), 10, mj.Brotzeit < 0)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"CronJobs":`)
	if mj.CronJobs != nil {
		buf.WriteString(`[`)
		for i, v := range mj.CronJobs {
			if i != 0 {
				buf.WriteString(`,`)
			}
			helpers.Ffjson_WriteJsonString(buf, v)
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Goroutines":`)
	helpers.Ffjson_FormatBits(buf, uint64(mj.Goroutines), 10, mj.Goroutines < 0)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"SessionExpires":`)
	helpers.Ffjson_FormatBits(buf, uint64(mj.SessionExpires), 10, mj.SessionExpires < 0)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Wanderers":`)
	helpers.Ffjson_FormatBits(buf, uint64(mj.Wanderers), 10, mj.Wanderers < 0)
	buf.WriteString(`}`)
	return nil
}
