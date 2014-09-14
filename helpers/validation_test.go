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

package helpers

import (
	"errors"
	"testing"
)

// only use such a construct (predefined map returned by a func) in a test, never in prod code as this is slow ...
// emailAddress => isValid?
func getEmailTestData() map[string]bool {
	// http://en.wikipedia.org/wiki/Email_address#Valid_email_addresses
	return map[string]bool{
		"a@b":                                            true,
		"a@":                                             false,
		"hello":                                          false,
		"hello@world":                                    true,
		"niceandsimple@example.com":                      true,
		"very.common@example.com":                        true,
		"very+common@gmail.com":                          true,
		"a.little.lengthy.but.fine@dept.example.com":     true,
		"disposable.style.email.with+symbol@example.com": true,
		"other.email-with-dash@example.com":              true,
		"user@localserver":                               true,
		"Abc.example.com":                                false, // (an @ character must separate the local and domain parts)
		"A@b@c@example.com":                              false, // (only one @ is allowed outside quotation marks)
		"a\"b(c)d,e:f;g<h>i[j\\k]l@example.com":          false, // (none of the special characters in this local part is allowed outside quotation marks)
		"just\"not\"right@example.com":                   false, // (quoted strings must be dot separated or the only element making up the local-part)
		"this is\"not\\allowed@example.com":              false, // (spaces, quotes, and backslashes may only exist when within quoted strings and preceded by a backslash)
		"this\\ still\\\"not\\allowed@example.com":       false, // (even if escaped (preceded by a backslash), spaces, quotes, and backslashes must still be contained by quotes)
		"john..doe@example.com":                          false, // (double dot before @)
		"john.doe@example..com":                          false, // (double dot after @)
	}
}

func TestValidateEmail(t *testing.T) {
	for email, expected := range getEmailTestData() {
		if actual := ValidateEmail(email); expected != actual {
			t.Errorf("Expected: %t Actual: %t Email: %s", expected, actual, email)
		}
	}
}

// old: BenchmarkValidateEmail	   10000	    186637 ns/op	Mid  2012; 1.8 GHz Intel Core i5; OS X 10.9.4 (13E28)
// new: BenchmarkValidateEmail	   10000	     27534 ns/op	"
// new: BenchmarkValidateEmail	  100000	     26940 ns/op	"
// diff between old/new: old code was the rune maps were returned by functions like in getEmailTestData()
func BenchmarkValidateEmail(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for email, _ := range getEmailTestData() {
			ValidateEmail(email)
		}
	}
}

type listenAddressResult struct {
	host string
	port string
	err  error
}

func getListenAddressData() map[string]listenAddressResult {
	return map[string]listenAddressResult{
		":3109": {
			"127.0.0.1",
			"3109",
			nil,
		},
		"246.245.145.235:3108": {
			"246.245.145.235",
			"3108",
			nil,
		},
		"246.245.145.235": {
			"246.245.145.235",
			"",
			errors.New("Missing : separator or too many"),
		},
	}
}

func TestValidateListenAddress(t *testing.T) {
	var host, port string
	var err error

	for input, res := range getListenAddressData() {
		host, port, err = ValidateListenAddress(input)

		if nil == err && res.host != host {
			t.Errorf("Expected %s got %s", res.host, host)
		}
		if nil == err && res.port != port {
			t.Errorf("Expected %s got %s", res.port, port)
		}
		if err != nil && res.err.Error() != err.Error() {
			t.Errorf("Expected %s got %s", res.err, err)
		}
	}
}

/**
05. Sept 2014; run with `go test -v -bench=.`
BenchmarkValidateListenAddress	 1000000	      1594 ns/op  Late 2013; 2.4 GHz Intel Core i5; OS X 10.9.4 (13E28)
BenchmarkValidateListenAddress	 1000000	      2003 ns/op  Mid  2012; 1.8 GHz Intel Core i5; OS X 10.9.4 (13E28)
go version go1.3.1 darwin/amd64
*/
func BenchmarkValidateListenAddress(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for input := range getListenAddressData() {
			ValidateListenAddress(input)
		}
	}
}
