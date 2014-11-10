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

package textarea

import (
	"bytes"
	"errors"
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"github.com/SchumacherFM/wanderlust/rucksack"
	"strings"
)

func GetProvisioner() *provisionerApi.Config {
	t := &ta{
		Base: provisionerApi.Base{
			TheRoute:  "textarea",
			TheConfig: []string{"TextAreaData"}, // used in the html input field names
		},
	}
	p := provisionerApi.NewProvisioner("Textarea", "fa-file-text-o", t)
	return p
}

type (
	ta struct {
		provisionerApi.Base
		// other fields ...
	}
)

var (
	ErrValidate    = errors.New("Failed to validate the value")
	ErrTooManyURLs = errors.New("Too many URLs detected! Maximum is 20.")
)

// PrepareSave removes duplicated lines and too much white space. creates a unique URL set
// @todo extract URLs from pasted text and ignore invalid characters. nice and convenient.
func (t *ta) PrepareSave(p *provisionerApi.PostData) ([]byte, error) {
	if "" == p.Value {
		return nil, nil
	}

	valueSlice := strings.Split(strings.TrimSpace(p.Value), "\n")
	if len(valueSlice) > 20 {
		return nil, ErrTooManyURLs
	}
	var ret bytes.Buffer
	defer ret.Reset() // free memory
	unique := make(map[string]bool)
	defer func() { unique = nil }() // free memory
	for _, v := range valueSlice {
		v = strings.TrimSpace(v)
		vl := strings.ToLower(v)

		if "" == vl {
			continue
		}

		if false == strings.HasPrefix(vl, "http") {
			return nil, ErrValidate
		}

		if _, isSet := unique[v]; false == isSet {
			ret.WriteString(strings.TrimSpace(v) + "\n")
			unique[v] = true
		}
	}
	if ret.Len() > 1 {
		ret.Truncate(ret.Len() - 1) // remove last \n
	}

	return ret.Bytes(), nil
}

// ConfigComplete implements the brotzeit.Fetcher interface to check if all config values
// have been successfully entered by the user. if so brotzeit can start automatically fetching URLs
func (t *ta) ConfigComplete(bp rucksack.Backpacker) (bool, error) {
	tad, err := bp.FindOne(t.Route(), "TextAreaData")
	if nil != err {
		return false, err
	}
	return len(tad) > 5, nil
}

// FetchUrls implements the brotzeit.Fetcher interface
func (s *ta) FetchUrls(bp rucksack.Backpacker) []string {
	return nil
}
