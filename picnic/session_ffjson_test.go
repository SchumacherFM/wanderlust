package picnic

import (
	"bytes"
	"encoding/json"
	"testing"
)

type testUser struct {
}

// satisfy interface userSessionIf
func (tu *testUser) getEmail() string        { return "root@localhost.dev" }
func (tu *testUser) getName() string         { return "Joanna Gopher" }
func (tu *testUser) getUserName() string     { return "gopher" }
func (tu *testUser) isAdmin() bool           { return true }
func (tu *testUser) isValidForSession() bool { return true }

var expected = []byte(`{"email":"root@localhost.dev","isAdmin":true,"loggedIn":true,"name":"Joanna Gopher","username":"gopher"}`)

func TestFFMarshalJSON(t *testing.T) {
	tu := &testUser{}
	si := newSessionInfo(tu)

	actual, err := si.MarshalJSON()
	if nil != err {
		t.Error(err)
	}

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("\nExpected: %s\nActual: %s\n", expected, actual)
	}
}

func TestGoMarshalJSON(t *testing.T) {
	tu := &testUser{}
	si := newSessionInfo(tu)
	var jBufActual bytes.Buffer
	err := json.NewEncoder(&jBufActual).Encode(si)
	if nil != err {
		t.Error(err)
	}
	es := string(expected)
	jBufActualBytes := jBufActual.Bytes()
	jBufActualBytes = jBufActualBytes[:len(jBufActualBytes)-1] // remove last \n
	as := string(jBufActualBytes)
	if as != es {
		t.Errorf("\nExpected: %s\nActual:   %s\n", es, as)
	}
}

// BenchmarkFFMarshalJSON	 1,000,000	      1459 ns/op
func BenchmarkFFMarshalJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tu := &testUser{}
		si := newSessionInfo(tu)
		si.MarshalJSON()
	}
}

// BenchmarkGoMarshalJSON	  500,000	      3526 ns/op
func BenchmarkGoMarshalJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tu := &testUser{}
		si := newSessionInfo(tu)
		var jBufActual bytes.Buffer
		json.NewEncoder(&jBufActual).Encode(si)
	}
}
