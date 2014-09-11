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
	"crypto/tls"
)

const (
	DEFAULT_TLS_SESSION_CACHE_CAPACITY int = 128
)

func GetTlsConfig() *tls.Config {
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

	// https://twitter.com/karlseguin/status/508531717011820544
	tlsConfig.ClientSessionCache = tls.NewLRUClientSessionCache(DEFAULT_TLS_SESSION_CACHE_CAPACITY)

	return tlsConfig
}
