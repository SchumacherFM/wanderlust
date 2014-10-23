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

import (
	"github.com/SchumacherFM/wanderlust/github.com/boltdb/bolt"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"time"
)

type (
	RDBIF interface {
		GoRoutineWriter() error
		Close() error
		//		CreateDatabase(name string) error
		//		UseDatabase(name string) *db.Col
		FindOne(dbName string, documentId int) (doc map[string]interface{}, err error)
		FindAll(dbName string) (doc []map[string]interface{}, err error)
		// Inserts Data: b = bucket Name, k = key, d = data
		Insert(b, k string, d []byte) error
	}
	bEntity struct {
		bucket string
		key    string
		data   []byte
	}

	RDB struct {
		db         *bolt.DB
		writerChan chan *bEntity
		logger     *log.Logger
	}
)

func NewRDB(dbFileName string, l *log.Logger) (*RDB, error) {
	var err error
	var db *bolt.DB
	// @see idea from http://paulosuzart.github.io/blog/2014/07/07/going-back-to-go/
	boltOpt := &bolt.Options{
		Timeout: 1 * time.Second,
	}
	w := make(chan *bEntity, 10)
	db, err = bolt.Open(dbFileName, 0600, boltOpt)
	rdb := &RDB{
		db:         db,
		writerChan: w,
		logger:     l,
	}
	return rdb, errgo.Mask(err)
}

func (this *RDB) GoRoutineWriter() error {
	for data := range this.writerChan {
		err := this.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(data.Bucket))
			if err != nil {
				return err
			}
			return b.Put([]byte(data.Key), data.Data)
		})
		if nil != err {
			this.logger.Emergency("DB update failed: %s", err)
		}
	}
	return nil
}

//func (rdb *RDB) CreateDatabase(name string) error {
//	if nil == rdb.UseDatabase(name) {
//		return rdb.db.Create(name)
//	}
//	return nil
//}
//
//func (rdb *RDB) UseDatabase(name string) *db.Col {
//	return rdb.db.Use(name)
//}

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

func (rdb *RDB) Insert(b, k string, d []byte) error {
	be := &bEntity{
		bucket: b,
		key:    k,
		data:   d,
	}
	rdb.writerChan <- be
	return nil
}
