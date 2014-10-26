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
	"errors"
	"github.com/SchumacherFM/wanderlust/github.com/boltdb/bolt"
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"time"
)

const (
	WRITER_CHANNEL_BUFFER_SIZE = 100
)

var (
	ErrDatabaseNotFound = errors.New("Database not found")
	ErrEntityNotFound   = errors.New("DB Entity not found")
	// Hook that may be overridden for integration tests.
	writerDone = func() {}
)

type (
	// @todo use encoding BinaryMarshaler and BinaryUnmarshaler interface
	// This interface can have various implementation for saving struct in database. JSON is only one option.
	UserEncoding interface {
		// Decode decodes the data which is coming from the database
		Decode(data []byte) error
		// Encode encodes the data for saving in the database
		Encode() ([]byte, error)
	}

	RDBIF interface {
		// Writer runs in a goroutine and waits for data coming through the channel
		Writer()

		// WriterOnce() @todo checks if the key already exists, if so returns something that key exits

		// Close closes the database when terminating the app
		Close() error

		FindOne(b, k string) ([]byte, error)

		FindAll(bn string) ([][]byte, error)
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
	return rdb, err
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

// Close closes the database during app shutdown sequence
func (this *RDB) Close() error {
	return this.db.Close()
}

// FindOne returns a value for a bucketName and keyName
func (this *RDB) FindOne(b, k string) ([]byte, error) {
	data := newbEntity(b, k, nil)
	this.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(data.getBucketByte())
		if nil == b {
			return ErrDatabaseNotFound
		}
		data.data = b.Get(data.getKeyByte())
		return nil
	})
	if nil == data.data {
		return nil, ErrEntityNotFound
	}
	return data.data, nil
}

// FindAll finds all keys belonging to a bucketName. Returns an array i = key, i+1 = value
func (this *RDB) FindAll(bn string) ([][]byte, error) {

	tx, err := this.db.Begin(false)
	if nil != err {
		return nil, err
	}
	b := tx.Bucket([]byte(bn))
	if nil == b {
		return nil, ErrDatabaseNotFound
	}
	bStat := b.Stats()
	ret := make([][]byte, 2*bStat.KeyN)
	var i = 0
	b.ForEach(func(k, v []byte) error {
		ret[i] = k
		ret[i+1] = v
		i = i + 2
		return nil
	})
	tx.Commit()
	return ret, nil
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
