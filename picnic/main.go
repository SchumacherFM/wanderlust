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
	"fmt"
	"github.com/SchumacherFM/wanderlust/github.com/gorilla/mux"
	"github.com/SchumacherFM/wanderlust/helpers"
	"log"
	"net/http"
	"os"
	"time"
)

type PicnicApp struct {
	ListenAddress string
	PemDir        string
	Logger *log.Logger
}

func (p *PicnicApp) Execute() {
	r := mux.NewRouter()
	r.HandleFunc("/", dashBoardHandler).Methods("GET")
	r.HandleFunc("/test", testDataHandler).Methods("GET")
	//	http.Handle("/", r)

	server := &http.Server{
		Addr:    p.GetListenAddress(),
		Handler: r,
	}
	err := server.ListenAndServeTLS(p.generatePems())
	if nil != err {
		p.Logger.Fatal("Picnic ListenAndServe: ", err)
	}
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
	address, _, _ := helpers.ValidateListenAddress(p.ListenAddress)
	duration := 180 * 24 * time.Hour // half a year
	validFrom := time.Now().Format("Jan 2 15:04:05 2006") // any other year number results in a weird result :-?

	// @todo check if pems exists if so do nothing
	certGenerator := &helpers.GenerateCert{
		Host:         address,   // "Comma-separated hostnames and IPs to generate a certificate for"
		ValidFrom:    validFrom, // "Creation date formatted as Jan 1 15:04:05 2011"
		ValidFor:     duration,  // "duration", 365*24*time.Hour, "Duration that certificate is valid for"
		IsCA:         true,     // "ca", false, "whether this cert should be its own Certificate Authority"
		RsaBits:      2048,      // "rsa-bits", 2048, "Size of RSA key to generate"
		CertFileName: dir + "cert.pem",
		KeyFileName:  dir + "key.pem",
	}

	certGenerator.Generate()

	return certGenerator.CertFileName, certGenerator.KeyFileName
}

func dashBoardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there! This is the DashBoard %#v!", r.URL)
}

func testDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test data Hi there! %#v!", r.URL)
}
