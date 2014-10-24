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
	"bytes"
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/helpers"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestGetBucketByte(t *testing.T) {
	expectedB := []byte(`Bucket`)
	be := &bEntity{
		bucket: "Bucket",
		key:    "Key",
		data:   []byte(`Data`),
	}
	if 0 != bytes.Compare(expectedB, be.getBucketByte()) {
		t.Errorf("\nExpected\t%s\nActual\t\t%s\n", expectedB, be.getBucketByte())
	}
}

func TestGetKeyByte(t *testing.T) {
	expectedK := []byte(`Key`)
	be := &bEntity{
		bucket: "Bucket",
		key:    "Key",
		data:   []byte(`Data`),
	}
	if 0 != bytes.Compare(expectedK, be.getKeyByte()) {
		t.Errorf("\nExpected\t%s\nActual\t\t%s\n", expectedK, be.getKeyByte())
	}
}

func setUpDb(f string) (*RDB, *log.Logger, error) {
	l := log.New(os.Stdout, log.DEBUG, "Testing")
	db, err := NewRDB(f, l)
	return db, l, err
}

func TestNewRDB(t *testing.T) {
	f := helpers.GetTempDir() + "wldbTEST1_" + helpers.RandomString(10) + ".db"
	db, l, err := setUpDb(f)
	defer func() {
		err := os.Remove(f)
		if nil != err {
			t.Error(err)
		}
	}()

	if nil != err {
		t.Error(err)
	}
	if nil == db.writerChan {
		t.Error("db.writerChan is nil")
	}
	if false == reflect.DeepEqual(l, db.logger) {
		t.Errorf("Loggers are different!\nE: %#v\nA: %#v\n", l, db.logger)
	}
	if _, ferr := os.Stat(f); os.IsNotExist(ferr) {
		t.Errorf("File not created: %s\n%s", f, ferr)
	}
}

func TestWriterIntegration(t *testing.T) {

	f := helpers.GetTempDir() + "wldbTEST2_" + helpers.RandomString(10) + ".db"
	db, _, err := setUpDb(f)
	defer func() {
		err := os.Remove(f)
		if nil != err {
			t.Error(err)
		}
	}()

	if nil != err {
		t.Error(err)
	}

	// Replace writerDone with a closure that will tell us when the writer is
	// exiting.
	done := make(chan bool)
	writerDone = func() {
		close(db.writerChan)
		done <- true
	}

	// Put things as they were when the test finishes.
	defer func() {
		writerDone = func() {}
	}()

	go db.Writer()

	// here we create 10 entries in the boltdb by writing into the channel
	testData := [10][]byte{} // []byte array with 10 entries
	for i := 0; i < 10; i++ {
		testData[i] = []byte(helpers.RandomString(i + 1*10))
		err = db.Insert("TestBucket", "TestKey"+strconv.Itoa(i), testData[i])
		if nil != err {
			t.Error(err)
		}
	}

	// I think that is lame with time.Sleep, waiting for the end of the write
	//	time.Sleep(1 * time.Millisecond)

	<-done // Wait for Writer() to empty the channel

	for i := 0; i < 10; i++ {
		foundData, err := db.FindOne("TestBucket", "TestKey"+strconv.Itoa(i))
		if nil != err {
			t.Error(err)
		}
		if 0 != bytes.Compare(testData[i], foundData) {
			t.Errorf("Data NOT found!\nE: %#v\nA: %#v\n", testData, foundData)
		}
	}
}
