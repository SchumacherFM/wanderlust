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

package rucksackdb

// tiedot is pretty slow database. but lets replace it later

import (
	"github.com/SchumacherFM/wanderlust/github.com/HouzuoGuo/tiedot/db"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
)

type RDBIF interface {
	CreateDatabase(name string) error
	Close() error
	UseDatabase(name string) *db.Col
	FindOne(dbName string, documentId int) (doc map[string]interface{}, err error)
	FindAll(dbName string) (doc []map[string]interface{}, err error)
	Insert(dbName string, doc map[string]interface{}) (id int, err error)
	InsertRecovery(dbName string, id int, doc map[string]interface{}) (err error)
}

type RDB struct {
	db *db.DB
}

func NewRDB(dbDir string) (RDBIF, error) {
	rdb := &RDB{}
	var err error
	rdb.db, err = db.OpenDB(dbDir)
	return rdb, err
}

func (rdb *RDB) CreateDatabase(name string) error {
	if nil == rdb.UseDatabase(name) {
		return rdb.db.Create(name)
	}
	return nil
}

func (rdb *RDB) UseDatabase(name string) *db.Col {
	return rdb.db.Use(name)
}

func (rdb *RDB) Close() error {
	return rdb.db.Close()
}

func (rdb *RDB) FindOne(dbName string, documentId int) (doc map[string]interface{}, err error) {
	doc, err = rdb.UseDatabase(dbName).Read(documentId)
	return
}

func (rdb *RDB) FindAll(dbName string) (docs []map[string]interface{}, err error) {
	var currentDb *db.Col
	err = nil
	currentDb = rdb.UseDatabase(dbName)
	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys
	err = db.EvalAllIDs(currentDb, &queryResult)
	if nil != err {
		return nil, errgo.Mask(err)
	}

	docs = make([]map[string]interface{}, len(queryResult))

	// Query result are document IDs
	c := 0
	for id, _ := range queryResult {
		d, err := currentDb.Read(id)
		if nil != err {
			return nil, err
		}
		docs[c] = d
		c++
	}
	return
}

func (rdb *RDB) Insert(dbName string, doc map[string]interface{}) (id int, err error) {
	id, err = rdb.UseDatabase(dbName).Insert(doc)
	return
}

func (rdb *RDB) InsertRecovery(dbName string, id int, doc map[string]interface{}) (err error) {
	err = rdb.UseDatabase(dbName).InsertRecovery(id, doc)
	return
}
