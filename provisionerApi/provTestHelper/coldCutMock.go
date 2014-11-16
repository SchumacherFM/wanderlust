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

package provTestHelper

import (
	"github.com/SchumacherFM/wanderlust/provisionerApi"
	"github.com/SchumacherFM/wanderlust/rucksack"
)

type (
	ColdCutMock struct {
		RouteMock  string
		ConfigMock []string
	}
)

var _ provisionerApi.ColdCutter = &ColdCutMock{} // check if struct implements interface

func (c *ColdCutMock) Route() string                                           { return c.RouteMock }
func (c *ColdCutMock) Config() []string                                        { return c.ConfigMock }
func (c *ColdCutMock) PrepareSave(pd *provisionerApi.PostData) ([]byte, error) { return nil, nil }
func (c *ColdCutMock) ConfigComplete(bp rucksack.Backpacker) (bool, error)     { return false, nil }
func (c *ColdCutMock) FetchUrls(bp rucksack.Backpacker) []string               { return nil }
