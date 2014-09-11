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
	"github.com/SchumacherFM/wanderlust/helpers"
	"log"
)

// @todo create maybe custom REST API instead of using the tiedot unsecure API

type RucksackApp struct {
	rdb           RucksackDbI
	Logger        *log.Logger
	ListenAddress string
}

func NewRucksackApp(listenAddress, dbDir string, logger *log.Logger) (*RucksackApp, error) {
	rucksackApp := &RucksackApp{
		ListenAddress: listenAddress,
		Logger:        logger,
	}
	rucksackApp.initDb(dbDir)
	return rucksackApp, nil
}

func (r *RucksackApp) initDb(dbDir string) error {
	var err error
	if "" == dbDir {
		dbDir = helpers.GetTempDir() + "wldb_" + helpers.RandomString(10)
		r.Logger.Printf("Database temp directory is %s", dbDir)
	}
	helpers.CreateDirectoryIfNotExists(dbDir)
	r.rdb, err = NewRucksackDb(dbDir)
	return err
}

func (r *RucksackApp) GetDb() RucksackDbI {
	return r.rdb
}

// listens on the DefaultServeMux and runs in a goroutine
func (r *RucksackApp) StartHttp() {
	r.Logger.Printf("Database webinterface running: %s", r.GetListenAddress())
	err := r.rdb.StartHttp(r.GetListenAddress())
	if nil != err {
		r.Logger.Fatal(err)
	}
}

func (r *RucksackApp) GetListenAddress() string {
	address, port, err := helpers.ValidateListenAddress(r.ListenAddress)
	if nil != err {
		r.Logger.Fatal(err, r.ListenAddress)
	}
	return address + ":" + port
}
