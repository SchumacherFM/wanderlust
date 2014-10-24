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

package api

import (
	"github.com/SchumacherFM/wanderlust/github.com/boltdb/bolt"
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"time"
)

const (
	WRITER_CHANNEL_BUFFER_SIZE = 10
)

var (
	ErrEntityNotFound = errgo.New("Entity not found")
	// Hook that may be overridden for integration tests.
	writerDone = func() {}
)

type (
	RDBIF interface {
		// Writer runs in a goroutine and waits for data coming through the channel
		Writer()
		// Close closes the database when terminating the app
		Close() error
		//		CreateDatabase(name string) error
		//		UseDatabase(name string) *db.Col
		FindOne(b, k string) ([]byte, error)
		//		FindAll(dbName string) (doc []map[string]interface{}, err error)
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

func newbEntity(b, k string, d []byte) *bEntity {
	be := &bEntity{
		bucket: b,
		key:    k,
		data:   d,
	}
	return be
}

func (this *bEntity) getBucketByte() []byte {
	return []byte(this.bucket)
}
func (this *bEntity) getKeyByte() []byte {
	return []byte(this.key)
}

func NewRDB(dbFileName string, l *log.Logger) (*RDB, error) {
	var err error
	var db *bolt.DB
	// @see idea from http://paulosuzart.github.io/blog/2014/07/07/going-back-to-go/
	boltOpt := &bolt.Options{
		Timeout: 1 * time.Second,
	}
	w := make(chan *bEntity, WRITER_CHANNEL_BUFFER_SIZE)
	db, err = bolt.Open(dbFileName, 0600, boltOpt)
	rdb := &RDB{
		db:         db,
		writerChan: w,
		logger:     l,
	}
	return rdb, errgo.Mask(err)
}

// Writer method runs in a goroutine and waits to data in the channel to write into the boltdb
func (this *RDB) Writer() {
	for data := range this.writerChan {
		err := this.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(data.getBucketByte())
			if err != nil {
				return err
			}
			return b.Put(data.getKeyByte(), data.data)
		})
		if nil != err {
			this.logger.Emergency("DB update failed: %s", err)
		}
		if 0 == len(this.writerChan) {
			writerDone()
		}
	}
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

func (this *RDB) Close() error {
	return this.db.Close()
}

func (this *RDB) FindOne(b, k string) ([]byte, error) {
	data := newbEntity(b, k, nil)
	this.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(data.getBucketByte())
		if nil == b {
			return errgo.Newf("Bucket %s not found", data.bucket)
		}
		data.data = b.Get(data.getKeyByte())
		return nil
	})
	if nil == data.data {
		return nil, ErrEntityNotFound
	}
	return data.data, nil
}

//
//func (rdb *RDB) FindAll(dbName string) (docs []map[string]interface{}, err error) {
//	var currentDb *db.Col
//	err = nil
//	currentDb = rdb.UseDatabase(dbName)
//	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys
//	err = db.EvalAllIDs(currentDb, &queryResult)
//	if nil != err {
//		return nil, errgo.Mask(err)
//	}
//
//	docs = make([]map[string]interface{}, len(queryResult))
//
//	// Query result are document IDs
//	c := 0
//	for id, _ := range queryResult {
//		d, err := currentDb.Read(id)
//		if nil != err {
//			return nil, err
//		}
//		docs[c] = d
//		c++
//	}
//	return
//}

func (rdb *RDB) Insert(b, k string, d []byte) error {
	be := &bEntity{
		bucket: b,
		key:    k,
		data:   d,
	}
	rdb.writerChan <- be
	return nil
}
