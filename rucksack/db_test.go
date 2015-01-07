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
	"os"
	"strconv"
	"strings"
	"testing"

	log "github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/SchumacherFM/wanderlust/helpers"
)

var (
	bucketName  = "RBucket"
	keyPrefix   = "RKey"
	benchMarkDb = helpers.GetTempDir() + "wldbBENCH_" + helpers.RandomString(10) + ".db"
	benchDb     *Rucksack
)

func init() {
	benchDb, _, _ = setUpDb(benchMarkDb)
	go benchDb.Writer()

	isBench := false
	for _, arg := range os.Args {
		if false == isBench && strings.Contains(arg, ".bench") {
			isBench = true
			break
		}
	}
	// remove created bench.db when no benchmark is required ... :-\
	if false == isBench {
		defer os.Remove(benchMarkDb)
	}
}

func TestGetBucketByte(t *testing.T) {
	expectedB := []byte(`Bucket`)
	be := &bEntity{
		bucket: "Bucket",
		key:    "Key",
		data:   []byte(`Data`),
	}
	assert.Exactly(t, expectedB, be.getBucketByte())
}

func TestGetKeyByte(t *testing.T) {
	expectedK := []byte(`Key`)
	be := &bEntity{
		bucket: "Bucket",
		key:    "Key",
		data:   []byte(`Data`),
	}
	assert.Exactly(t, expectedK, be.getKeyByte())
}

func setUpDb(f string) (*Rucksack, *log.Logger, error) {
	l := log.New(os.Stdout, log.DEBUG, "Testing")
	db, err := New(f, l)
	return db, l, err
}

func TestNewRDB(t *testing.T) {
	f := helpers.GetTempDir() + "wldbTEST1_" + helpers.RandomString(10) + ".db"
	db, l, err := setUpDb(f)
	defer func() {
		err := os.Remove(f)
		assert.NoError(t, err)
	}()
	assert.NoError(t, err)
	assert.NotNil(t, db.writerChan, "db.writerChan is nil")
	assert.Equal(t, l, db.logger)

	if _, ferr := os.Stat(f); os.IsNotExist(ferr) {
		t.Errorf("File not created: %s\n%s", f, ferr)
	}
}

func TestWriterIntegration(t *testing.T) {

	f := helpers.GetTempDir() + "wldbTEST2_" + helpers.RandomString(10) + ".db"
	db, _, err := setUpDb(f)
	defer func() {
		err := os.Remove(f)
		assert.NoError(t, err)
	}()
	assert.NoError(t, err)

	// Replace WriterDone with a closure that will tell us when the writer is
	// exiting.
	done := make(chan bool)
	WriterDone = func() {
		close(db.writerChan)
		done <- true
	}

	// Put things as they were when the test finishes.
	defer func() {
		WriterDone = func() {}
	}()

	go db.Writer()

	// here we create 10 entries in the boltdb by writing into the channel
	testData := [10][]byte{} // []byte array with 10 entries
	for i := 0; i < 10; i++ {
		testData[i] = []byte(helpers.RandomString(i + 1*10))
		err = db.Insert(bucketName, keyPrefix+strconv.Itoa(i), testData[i])
		assert.NoError(t, err)
	}

	// I think that is lame with time.Sleep, waiting for the end of the write
	//	time.Sleep(1 * time.Millisecond)

	<-done // Wait for Writer() to empty the channel

	for i := 0; i < 10; i++ {
		foundData, err := db.FindOne(bucketName, keyPrefix+strconv.Itoa(i))
		assert.NoError(t, err)
		assert.Exactly(t, testData[i], foundData)
	}

	// test FindAll()
	result, err := db.FindAll(bucketName)
	assert.NoError(t, err)

	// test Count()
	kc, err := db.Count(bucketName)
	assert.NoError(t, err)
	assert.Equal(t, len(result)/2, kc)

	ti := 0
	for j := 0; j < len(result); j = j + 2 {
		actualKey := result[j]
		actualVal := result[j+1]
		assert.Exactly(t, []byte(keyPrefix+strconv.Itoa(ti)), actualKey)
		assert.Exactly(t, testData[ti], actualVal)
		if 0 == j%2 {
			ti++
		}
	}

	// test Delete
	err = db.Delete(bucketName, keyPrefix+"1")
	assert.NoError(t, err)
	notFoundData, err := db.FindOne(bucketName, keyPrefix+"1")
	assert.Error(t, err)
	assert.Empty(t, notFoundData)

	err = db.Delete(bucketName, keyPrefix+"NonExistent")
	assert.NoError(t, err)
}

// BufferSize = 10
// BenchmarkInsert	    5000	    432939 ns/op	   38337 B/op	      63 allocs/op 2.218s
// BufferSize = 100
// BenchmarkInsert	   10000	    437301 ns/op	   42151 B/op	      67 allocs/op 4.392s
func BenchmarkInsert(b *testing.B) {
	// would be nice to get the current iteration of calling BenchmarkInsert to be more unique for the key ;-)
	for i := 0; i < b.N; i++ {
		key := keyPrefix + strconv.Itoa(i)
		err := benchDb.Insert(bucketName, key, []byte(`http://www.youtube.com/watch?v=LJvEIjRBSDA`))
		if nil != err {
			b.Error(err)
		}
	}
}

// BenchmarkFindOne	 1000000	      6728 ns/op	    1243 B/op	      20 allocs/op
func BenchmarkFindOne(b *testing.B) {

	// hmm no tear down .... so place that always in the last benchmark to run :-(
	defer os.Remove(benchMarkDb)

	// this benchmark will find entries which have not been created in the func BenchmarkInsert :-( @todo refactor
	for i := 0; i < b.N; i++ {
		key := keyPrefix + strconv.Itoa(i)
		_, err := benchDb.FindOne(bucketName, key)
		if nil != err {
			b.Log(key, err)
		}
	}
}
