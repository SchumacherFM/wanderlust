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
	log "github.com/SchumacherFM/wanderlust/github.com/segmentio/go-log"
	"github.com/SchumacherFM/wanderlust/helpers"
	. "github.com/SchumacherFM/wanderlust/picnic/api"
	rdb "github.com/SchumacherFM/wanderlust/rucksack/api"
	"net/http"
	"sync"
)

const (
	PEM_CERT                           string = "cert.pem"
	PEM_KEY                            string = "key.pem"
	RD_DIST_DIR                        string = "responsive-dashboard/dist/"
	DEFAULT_TLS_SESSION_CACHE_CAPACITY int    = 128
)

var (
	rsdb   rdb.RDBIF
	logger *log.Logger
)

type PicnicApp struct {
	ListenAddress string
	PemDir        string
	session       SessionManagerIf
	certFile      string
	keyFile       string
	httpRunning   sync.Once
}

// la = listen address, pd = pemDir, lo = logger
func NewPicnicApp(la, pd string, lo *log.Logger, theDb rdb.RDBIF) (*PicnicApp, error) {
	var err error
	rsdb = theDb
	logger = lo
	pa := &PicnicApp{
		ListenAddress: la,
		PemDir:        pd,
	}
	pa.certFile, pa.keyFile, err = pa.generatePems()
	if nil != err {
		return nil, err
	}
	pa.session, err = newSessionManager(pa.certFile, pa.keyFile)
	if nil != err {
		return nil, err
	}
	err = initUsers(rsdb)
	if nil != err {
		return nil, err
	}
	return pa, nil
}

func (p *PicnicApp) GetSessionManager() SessionManagerIf {
	return p.session
}

func (p *PicnicApp) getServer() *http.Server {
	s := &http.Server{
		Addr:      p.GetListenAddress(),
		Handler:   p.getHandler(),
		TLSConfig: helpers.GetTlsConfig(),
	}
	return s
}

func (p *PicnicApp) getPemDir() string {
	return p.PemDir
}

func (p *PicnicApp) generatePems() (certFile, keyFile string, err error) {
	// PemDir can be empty then it will generate a random one
	dir, err := helpers.GeneratePems(p.GetListenAddress(), p.getPemDir(), PEM_CERT, PEM_KEY)
	if nil != err {
		return "", "", err
	}
	if "" != dir {
		logger.Notice("PEM certificate temp directory is %s", dir)
	}
	p.PemDir = dir
	certFile = dir + PEM_CERT
	keyFile = dir + PEM_KEY
	return
}

// make sure to execute only once
func (p *PicnicApp) Execute() error {
	p.httpRunning.Do(func() {
		logger.Check(p.getServer().ListenAndServeTLS(p.certFile, p.keyFile))
	})
	return nil
}

func (p *PicnicApp) GetListenAddress() string {
	address, port, err := helpers.ValidateListenAddress(p.ListenAddress)
	if nil != err {
		logger.Check(err)
	}
	return address + ":" + port
}
