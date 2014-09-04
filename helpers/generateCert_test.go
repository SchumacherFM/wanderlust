// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.
// Cyrill: created a struct with a method to use it internally
package helpers

import (
	//	"crypto/x509"
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {

	duration := 180 * 24 * time.Hour                      // half a year
	validFrom := time.Now().Format("Jan 2 15:04:05 2006") // any other year number results in a weird result :-?

	pathSep := string(os.PathSeparator)
	dir := GetTempDir() + "wlpemTEST_" + RandomString(10)

	CreateDirectoryIfNotExists(dir)
	dir = dir + pathSep
	defer os.RemoveAll(dir)
	certFile := dir + "cert.pem"
	keyFile := dir + "key.pem"

	certGenerator := &GenerateCert{
		Host:         "127.0.0.1",
		ValidFrom:    validFrom,
		ValidFor:     duration,
		IsCA:         false,
		RsaBits:      2048,
		CertFileName: certFile,
		KeyFileName:  keyFile,
	}
	t.Log(certFile)
	t.Log(keyFile)

	err := certGenerator.Generate()
	if nil != err {
		t.Error(err)
	}

	mustContain := []byte(`--END CERTIFICATE--`)
	certContent, err := ioutil.ReadFile(certFile)
	if false == bytes.Contains(certContent, mustContain) {
		t.Errorf("%s does not contain %s\n%s", certFile, mustContain, certContent)
	}

	mustContain = []byte(`--END RSA PRIVATE KEY--`)
	keyContent, err := ioutil.ReadFile(keyFile)
	if false == bytes.Contains(keyContent, mustContain) {
		t.Errorf("%s does not contain %s\n%s", keyFile, mustContain, keyContent)
	}
	// @todo use maybe x509 verify method
}
