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

// Wanderlust uses go.rice package for serving static web content
//

package picnic

// @todo add kr/secureheader

import (
	"fmt"
	gzrice "github.com/SchumacherFM/wanderlust/github.com/SchumacherFM/go.gzrice"
	"github.com/SchumacherFM/wanderlust/github.com/gorilla/mux"
	"net/http"
)

func getHandler() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", dashBoardHandler).Methods("GET")
	router.HandleFunc("/test", testDataHandler).Methods("GET")

	router.HandleFunc("/favicon.ico", handlerFavicon)

	// due to the rice box regex when building embedded files we must use the full path in the MustFindBox method
	router.PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(gzrice.MustFindBox("rd/dist/css").HTTPBox())),
	)
	router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(gzrice.MustFindBox("rd/dist/fonts").HTTPBox())))
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(gzrice.MustFindBox("rd/dist/img").HTTPBox())))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(gzrice.MustFindBox("rd/dist/js").HTTPBox())))
	router.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(gzrice.MustFindBox("rd/dist/lib").HTTPBox())))

	return router
}

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	w.Write(gzrice.MustFindBox("rd/dist/img").MustBytes("favicon.ico"))
}

func dashBoardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there! This is the DashBoard %#v!", r.URL)
}

func testDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test data Hi there! %#v!", r.URL)
}
