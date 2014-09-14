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
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"github.com/SchumacherFM/wanderlust/helpers"
	"github.com/SchumacherFM/wanderlust/rucksack/rucksackdb"
	"log"
	"net/http"
)

const (
	PEM_CERT                           string = "cert.pem"
	PEM_KEY                            string = "key.pem"
	RD_DIST_DIR                        string = "responsive-dashboard/dist/"
	DEFAULT_TLS_SESSION_CACHE_CAPACITY int    = 128
)

var (
	rsdb   rucksackdb.RDBI
	logger *log.Logger
)

type PicnicAppI interface {
	getSessionManager() sessionManagerI
	getServer() *http.Server
	generatePems() (certFile, keyFile string, err error)
	Execute() error
	GetListenAddress() string
	getPemDir() string
}

type PicnicApp struct {
	ListenAddress string
	PemDir        string
	session       sessionManagerI
	certFile      string
	keyFile       string
}

func NewPicnicApp(listenAddress, pemDir string, theLogger *log.Logger, theDb rucksackdb.RDBI) (PicnicAppI, error) {
	var err error
	rsdb = theDb
	logger = theLogger
	picnicApp := &PicnicApp{
		ListenAddress: listenAddress,
		PemDir:        pemDir,
	}
	picnicApp.certFile, picnicApp.keyFile, err = picnicApp.generatePems()
	if nil != err {
		return nil, err
	}
	picnicApp.session, err = newSessionManager(picnicApp.certFile, picnicApp.keyFile)
	if nil != err {
		return nil, err
	}
	err = initUsers()
	if nil != err {
		return nil, err
	}
	return picnicApp, nil
}

func (p *PicnicApp) getSessionManager() sessionManagerI {
	return p.session
}

func (p *PicnicApp) getServer() *http.Server {
	server := &http.Server{
		Addr:      p.GetListenAddress(),
		Handler:   p.getHandler(),
		TLSConfig: helpers.GetTlsConfig(),
	}
	return server
}

func (p *PicnicApp) getPemDir() string {
	return p.PemDir
}

func (p *PicnicApp) generatePems() (certFile, keyFile string, err error) {
	// PemDir can be empty then it will generate a random one
	pemDir, err := helpers.GeneratePems(p.GetListenAddress(), p.getPemDir(), PEM_CERT, PEM_KEY)
	if nil != err {
		return "", "", err
	}
	if "" != pemDir {
		logger.Printf("PEM certificate temp directory is %s", pemDir)
	}
	p.PemDir = pemDir
	certFile = pemDir + PEM_CERT
	keyFile = pemDir + PEM_KEY
	return
}

func (p *PicnicApp) Execute() error {
	return errgo.Mask(p.getServer().ListenAndServeTLS(p.certFile, p.keyFile))
}

func (p *PicnicApp) GetListenAddress() string {
	address, port, err := helpers.ValidateListenAddress(p.ListenAddress)
	if nil != err {
		logger.Fatal(err, p.ListenAddress)
	}
	return address + ":" + port
}
