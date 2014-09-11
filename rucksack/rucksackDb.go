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

package rucksack

import (
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/db"
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/httpapi"
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/webcp"
)

type RucksackDbI interface {
	StartHttp(listenAddress string) error
	Close() error
	Query()
}

type RucksackDb struct {
	db            *db.DB
	isHttpRunning bool
}

func NewRucksackDb(dbDir string) (RucksackDbI, error) {
	rdb := &RucksackDb{}
	var err error
	rdb.db, err = db.OpenDB(dbDir)
	return rdb, err
}

func (rdb *RucksackDb) Close() error {
	return rdb.db.Close()
}

func (rdb *RucksackDb) Query() {

}

func (rdb *RucksackDb) StartHttp(listenAddress string) error {
	webcp.WebCp = "webcp"
	if false == rdb.isHttpRunning {
		rdb.isHttpRunning = true
		return httpapi.Start(rdb.db, listenAddress)
	}
	return nil
}
