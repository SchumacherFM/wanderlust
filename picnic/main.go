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

package picnic

import (
	"crypto/tls"
	"github.com/SchumacherFM/wanderlust/helpers"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	PEM_CERT    = "cert.pem"
	PEM_KEY     = "key.pem"
	RD_DIST_DIR = "responsive-dashboard/dist/"
)

type PicnicApp struct {
	ListenAddress string
	PemDir        string
	Logger        *log.Logger
}

func (p *PicnicApp) Execute() {

	server := &http.Server{
		Addr:      p.GetListenAddress(),
		Handler:   getRoutes(),
		TLSConfig: p.getTlsConfig(),
	}

	err := server.ListenAndServeTLS(p.generatePems())
	if nil != err {
		p.Logger.Fatal("Picnic ListenAndServe: ", err)
	}
}

func (p *PicnicApp) getTlsConfig() *tls.Config {
	tlsConfig := &tls.Config{}
	// @see http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html
	tlsConfig.CipherSuites = []uint16{
		//		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		//		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		//		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		//		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}
	tlsConfig.MinVersion = tls.VersionTLS12
	// no need to disable session resumption http://chimera.labs.oreilly.com/books/1230000000545/ch04.html#TLS_RESUME
	return tlsConfig
}

func (p *PicnicApp) GetListenAddress() string {
	address, port, err := helpers.ValidateListenAddress(p.ListenAddress)
	if nil != err {
		p.Logger.Fatal(err, p.ListenAddress)
	}
	return address + ":" + port
}

func (p *PicnicApp) generatePems() (certFile, keyFile string) {
	var dir string
	dir = p.PemDir
	pathSep := string(os.PathSeparator)
	if "" == dir {
		dir = helpers.GetTempDir() + "wlpem_" + helpers.RandomString(10)
		p.Logger.Printf("PEM certificate temp directory is %s", dir)
	}
	helpers.CreateDirectoryIfNotExists(dir)
	dir = dir + pathSep
	certFile = dir + PEM_CERT
	keyFile = dir + PEM_KEY
	address, _, _ := helpers.ValidateListenAddress(p.ListenAddress)
	duration := 180 * 24 * time.Hour                      // half a year
	validFrom := time.Now().Format("Jan 2 15:04:05 2006") // any other year number results in a weird result :-?

	isDir1, _ := helpers.PathExists(certFile)
	isDir2, _ := helpers.PathExists(keyFile)
	if isDir1 && isDir2 {
		return
	}
	// @todo check if pems exists if so do nothing
	certGenerator := &helpers.GenerateCert{
		Host:         address,   // "Comma-separated hostnames and IPs to generate a certificate for"
		ValidFrom:    validFrom, // "Creation date formatted as Jan 1 15:04:05 2011"
		ValidFor:     duration,  // "duration", 365*24*time.Hour, "Duration that certificate is valid for"
		IsCA:         true,      // "ca", false, "whether this cert should be its own Certificate Authority"
		RsaBits:      2048,      // "rsa-bits", 2048, "Size of RSA key to generate"
		CertFileName: certFile,
		KeyFileName:  keyFile,
	}
	certGenerator.Generate()
	return
}
