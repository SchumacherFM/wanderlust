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
	"log"
	"net/http"
)

const (
	LOCALHOST_IP4 = "127.0.0.1"
)

type PicnicApp struct {
	Port int
	Ip   string
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
	err := server.ListenAndServeTLS(certFile, keyFile)
	if nil != err {
		log.Fatal("Picnic ListenAndServe: ", err)
	}
}

func (p *PicnicApp) GetListenAddress() string {
	ip := LOCALHOST_IP4
	if "" != p.Ip {
		ip = p.Ip
	}
	return fmt.Sprintf("%s:%d", ip, p.Port)
}

func dashBoardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there! This is the DashBoard %#v!", r.URL)
}

func testDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test data Hi there! %#v!", r.URL)
}
