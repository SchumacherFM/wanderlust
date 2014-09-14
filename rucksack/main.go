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
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/rucksack/rucksackdb"
)

// @todo create maybe custom REST API instead of using the tiedot unsecure API

type RucksackApp struct {
	rdb           rucksackdb.RDBI
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
		r.Logger.Notice("Database temp directory is %s", dbDir)
	}
	helpers.CreateDirectoryIfNotExists(dbDir)
	r.rdb, err = rucksackdb.NewRDB(dbDir)
	return err
}

func (r *RucksackApp) GetDb() rucksackdb.RDBI {
	return r.rdb
}

// listens on the DefaultServeMux and runs in a goroutine
func (r *RucksackApp) StartHttp() {
	r.Logger.Notice("Database webinterface running: http://%s", r.GetListenAddress())
	err := r.rdb.StartHttp(r.GetListenAddress())
	if nil != err {
		r.Logger.Check(err)
	}
}

// @todo How to implement a stopable http server
// http://www.hydrogen18.com/blog/stop-listening-http-server-go.html
func (r *RucksackApp) StopHttp() error {
	return nil
}

func (r *RucksackApp) GetListenAddress() string {
	address, port, err := helpers.ValidateListenAddress(r.ListenAddress)
	if nil != err {
		r.Logger.Check(err)
	}
	return address + ":" + port
}
