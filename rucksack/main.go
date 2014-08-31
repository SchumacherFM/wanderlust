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
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/httpapi"
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/db"
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/webcp"
	"fmt"
	"os"
)

const (
	LOCALHOST_IP4 = "127.0.0.1"
)

type RucksackApp struct {
	Port      int
	Ip        string
	Db        *db.DB
	DbDir     string
}

func (p *RucksackApp) InitDb() {
	dbDir := p.DbDir
	var err error
	if "" == p.DbDir {
		dbDir = os.TempDir()+"/wldb_"+helpers.RandomString(10)
	}
	isDir, _ := helpers.DirectoryExists(dbDir)
	if false == isDir {
		err = os.Mkdir(dbDir, 0700)
		if nil != err {
			panic("Cannot create database directory: " + dbDir)
		}
	}
	p.Db, err = db.OpenDB(dbDir)
	if nil != err {
		panic(err)
	}
}

func (p *RucksackApp) GetDb() *db.DB {
	return p.Db
}

func (p *RucksackApp) StartHttp() {
	webcp.WebCp = "webcp"
	httpapi.Start(p.Db, p.GetListenAddress())
}

func (p *RucksackApp) GetListenAddress() string {
	ip := LOCALHOST_IP4
	if "" != p.Ip {
		ip = p.Ip
	}
	return fmt.Sprintf("%s:%d", ip, p.Port)
}
