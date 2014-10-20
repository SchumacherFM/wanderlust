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

type RucksackApp struct {
	rdb    rucksackdb.RDBIF
	Logger *log.Logger
}

func NewRucksackApp(dbDir string, logger *log.Logger) (*RucksackApp, error) {
	rucksackApp := &RucksackApp{
		Logger: logger,
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

func (r *RucksackApp) GetDb() rucksackdb.RDBIF {
	return r.rdb
}
