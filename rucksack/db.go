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
	"errors"
	"github.com/SchumacherFM/wanderlust/github.com/boltdb/bolt"
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/helpers"
	"time"
)

const (
	WRITER_CHANNEL_BUFFER_SIZE = 100
)

var (
	ErrBreadbasketNotFound = errors.New("Breadbasket / Database not found")
	ErrBreadNotFound       = errors.New("Bread / DB Entity not found")
	// Hook that may be overridden for integration tests.
	writerDone            = func() {}
	_          Backpacker = &Rucksack{}
)

type (
	Backpacker interface {

		// FindOne searches for bucketName and keyName to return the value or an error
		FindOne(string, string) ([]byte, error)

		// FindAll uses the bucketName to search for all keys/values in a database
		// Returns an array i = key, i+1 = value
		FindAll(string) ([][]byte, error)

		// Insert Data: bucketName, keyName, data
		Insert(string, string, []byte) error

		// Delete removes a keyName from a bucketName
		Delete(string, string) error

		// Count returns the number of keys in a bucketName or an error and key count = 0
		Count(string) (int, error)
	}

	bEntity struct {
		bucket string
		key    string
		data   []byte
	}

	Rucksack struct {
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

func New(dbFileName string, l *log.Logger) (*Rucksack, error) {

	if "" == dbFileName {
		dbFileName = helpers.GetTempDir() + "wldb_" + helpers.RandomString(10) + ".db"
		l.Notice("Database created: %s", dbFileName)
	}

	// @see idea from http://paulosuzart.github.io/blog/2014/07/07/going-back-to-go/ regarding channel
	boltOpt := &bolt.Options{
		Timeout: 1 * time.Second,
	}
	w := make(chan *bEntity, WRITER_CHANNEL_BUFFER_SIZE)
	db, err := bolt.Open(dbFileName, 0600, boltOpt)
	rdb := &Rucksack{
		db:         db,
		writerChan: w,
		logger:     l.New("RS"),
	}
	return rdb, err
}

// Writer method runs in a goroutine and waits to data in the channel to write into the boltdb
func (this *Rucksack) Writer() {
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

// Close closes the database during app shutdown sequence. implements io.Closer
func (this *Rucksack) Close() error {
	return this.db.Close()
}

// FindOne returns a value for a bucketName and keyName
func (this *Rucksack) FindOne(b, k string) ([]byte, error) {
	data := newbEntity(b, k, nil)
	this.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(data.getBucketByte())
		if nil == b {
			return ErrBreadbasketNotFound
		}
		data.data = b.Get(data.getKeyByte())
		return nil
	})
	if nil == data.data {
		return nil, ErrBreadNotFound
	}
	return data.data, nil
}

// FindAll finds all keys belonging to a bucketName. Returns an array i = key, i+1 = value
func (this *Rucksack) FindAll(bn string) ([][]byte, error) {

	tx, err := this.db.Begin(false)
	if nil != err {
		return nil, err
	}
	b := tx.Bucket([]byte(bn))
	if nil == b {
		return nil, ErrBreadbasketNotFound
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

func (this *Rucksack) Insert(b, k string, d []byte) error {
	be := &bEntity{
		bucket: b,
		key:    k,
		data:   d,
	}
	this.writerChan <- be
	return nil
}

func (this *Rucksack) Delete(b, k string) error {
	be := &bEntity{
		bucket: b,
		key:    k,
	}
	err := this.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(be.getBucketByte())
		return b.Delete(be.getKeyByte())
	})
	if nil != err {
		this.logger.Emergency("DB delete failed: %s", err)
	}
	return nil
}

// Count returns the total number of keys within that bucketName
func (this *Rucksack) Count(bn string) (int, error) {

	tx, err := this.db.Begin(false)
	if nil != err {
		return 0, err
	}
	b := tx.Bucket([]byte(bn))
	if nil == b {
		return 0, ErrBreadbasketNotFound
	}
	s := b.Stats()
	tx.Commit()
	return s.KeyN, nil
}
