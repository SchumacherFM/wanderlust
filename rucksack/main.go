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
	"github.com/SchumacherFM/wanderlust/helpers"
	"log"
)

type RucksackApp struct {
	ListenAddress     string
	db     *db.DB
	DbDir             string
	Logger *log.Logger
}

func (r *RucksackApp) InitDb() {
	dbDir := r.DbDir
	var err error
	if "" == r.DbDir {
		dbDir = helpers.GetTempDir()+"wldb_"+helpers.RandomString(10)
		r.Logger.Printf("Database temp directory is %s", dbDir)
	}
	helpers.CreateDirectoryIfNotExists(dbDir)
	r.db, err = db.OpenDB(dbDir)
	if nil != err {
		r.Logger.Fatalln("RucksackApp InitDB: ", err)
	}
}

func (r *RucksackApp) GetDb() *db.DB {
	return r.db
}

// listens on the DefaultServeMux
func (r *RucksackApp) StartHttp() {
	webcp.WebCp = "webcp"
	httpapi.Start(r.db, r.GetListenAddress())
}

func (r *RucksackApp) GetListenAddress() string {
	address, port, err := helpers.ValidateListenAddress(r.ListenAddress)
	if nil != err {
		r.Logger.Fatal(err, r.ListenAddress)
	}
	return address + ":" + port
}
