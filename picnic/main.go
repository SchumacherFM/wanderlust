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
	"github.com/SchumacherFM/wanderlust/github.com/juju/errgo"
	"github.com/SchumacherFM/wanderlust/helpers"
	"log"
	"net/http"
)

const (
	PEM_CERT                           string = "cert.pem"
	PEM_KEY                            string = "key.pem"
	RD_DIST_DIR                        string = "responsive-dashboard/dist/"
	DEFAULT_TLS_SESSION_CACHE_CAPACITY int    = 128
)

type PicnicAppI interface {
	getServer() *http.Server
	generatePems() (certFile, keyFile string, err error)
	Execute() error
	getTlsConfig() *tls.Config
	GetListenAddress() string
	getPemDir() string
}

type PicnicApp struct {
	ListenAddress string
	PemDir        string
	Logger        *log.Logger
	session       sessionManagerI
	certFile      string
	keyFile       string
}

func NewPicnicApp(listenAddress, pemDir string, logger *log.Logger) (PicnicAppI, error) {
	var err error
	picnicApp := &PicnicApp{
		ListenAddress: listenAddress,
		PemDir:        pemDir,
		Logger:        logger,
	}
	picnicApp.certFile, picnicApp.keyFile, err = picnicApp.generatePems()
	if nil != err {
		return nil, err
	}
	picnicApp.session, err = newSessionManager(picnicApp.certFile, picnicApp.keyFile)
	if nil != err {
		return nil, err
	}
	return picnicApp, nil
}

func (p *PicnicApp) getServer() *http.Server {
	server := &http.Server{
		Addr:      p.GetListenAddress(),
		Handler:   p.getHandler(),
		TLSConfig: p.getTlsConfig(),
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
		p.Logger.Printf("PEM certificate temp directory is %s", pemDir)
	}
	p.PemDir = pemDir
	certFile = pemDir + PEM_CERT
	keyFile = pemDir + PEM_KEY
	return
}

func (p *PicnicApp) Execute() error {
	return errgo.Mask(p.getServer().ListenAndServeTLS(p.certFile, p.keyFile))
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
	//tlsConfig.MinVersion = tls.VersionTLS12
	// no need to disable session resumption http://chimera.labs.oreilly.com/books/1230000000545/ch04.html#TLS_RESUME

	// https://twitter.com/karlseguin/status/508531717011820544
	tlsConfig.ClientSessionCache = tls.NewLRUClientSessionCache(DEFAULT_TLS_SESSION_CACHE_CAPACITY)

	return tlsConfig
}

func (p *PicnicApp) GetListenAddress() string {
	address, port, err := helpers.ValidateListenAddress(p.ListenAddress)
	if nil != err {
		p.Logger.Fatal(err, p.ListenAddress)
	}
	return address + ":" + port
}
