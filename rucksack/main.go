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
	rdb "github.com/SchumacherFM/wanderlust/rucksack/api"
)

type RucksackApp struct {
	rdb    rdb.RDBIF
	logger *log.Logger
}

func NewRucksackApp(dbFileName string, l *log.Logger) (*RucksackApp, error) {
	rucksackApp := &RucksackApp{
		logger: l,
	}
	rucksackApp.initDb(dbFileName)
	return rucksackApp, nil
}

func (r *RucksackApp) initDb(dbFileName string) error {
	var err error
	if "" == dbFileName {
		dbFileName = helpers.GetTempDir() + "wldb_" + helpers.RandomString(10) + ".db"
		r.logger.Notice("Database temp directory is %s", dbFileName)
	}
	r.rdb, err = rdb.NewRDB(dbFileName, r.logger)
	return err
}

func (r *RucksackApp) GetDb() rdb.RDBIF {
	return r.rdb
}
