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

// https://github.com/phyber/negroni-gzip/blob/master/LICENSE
// The MIT License (MIT)
// Copyright (c) 2013 Jeremy Saenz; 2014 David O'Rourke
// and modified by Cyrill

package picnic

import (
	"compress/gzip"
	gzrice "github.com/SchumacherFM/wanderlust/github.com/SchumacherFM/go.gzrice"
	"github.com/SchumacherFM/wanderlust/github.com/codegangsta/negroni"
	"net/http"
	"strings"
)

// These compression constants are copied from the compress/gzip package.
const (
	encodingGzip = "gzip"

	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentLength   = "Content-Length"
	headerContentType     = "Content-Type"
	headerVary            = "Vary"

	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression

	indexPage = "index.html"
)

func prepareRequestUri(uri string) string {
	if "" == uri || "/" == uri {
		return "/" + indexPage
	}
	parts := strings.Split(uri, "?")
	lastChar := ""
	if len(parts[0]) > 1 {
		lastChar = uri[len(parts[0])-1:]
	}
	if "/" == lastChar {
		return parts[0] + indexPage
	}
	return parts[0]
}

// checks if files in gzricebox are already compressed
// skips compression for binary files
// compresses html content if client supports
func GzipContentTypeMiddleware(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	isAllowed := gzrice.IsCompressingAllowed(prepareRequestUri(req.RequestURI))
	embeddedBoxesExists := gzrice.IsEmbedded() // only for dev ...
	if false == embeddedBoxesExists && 1 == isAllowed {
		isAllowed = 0
	}

	if 1 == isAllowed { // file in ricebox is already gzcompressed
		res.Header().Set(headerContentEncoding, encodingGzip)
		res.Header().Set(headerVary, headerAcceptEncoding)
		next(res, req)
		return
	}
	if -1 == isAllowed { // .jpg .png and other binary files ...
		next(res, req)
		return
	}

	// 0 == isAllowed
	h := newGzipHandler(BestSpeed)
	h.ServeHTTP(res, req, next)
}

// gzipHandler struct contains the ServeHTTP method and the compressionLevel to be
// used.
type gzipHandler struct {
	compressionLevel int
}

// Gzip returns a gzipHandler which will handle the Gzip compression in ServeHTTP.
// Valid values for level are identical to those in the compress/gzip package.
func newGzipHandler(level int) *gzipHandler {
	return &gzipHandler{
		compressionLevel: level,
	}
}

// ServeHTTP wraps the http.ResponseWriter with a gzip.Writer.
func (h *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Skip compression if the client doesn't accept gzip encoding.
	// @todo needs to be implemented in GzipContentTypeMiddleware
	if !strings.Contains(r.Header.Get(headerAcceptEncoding), encodingGzip) {
		next(w, r)
		return
	}

	// Create new gzip Writer. Skip compression if an invalid compression
	// level was set.
	gz, err := gzip.NewWriterLevel(w, h.compressionLevel)
	if err != nil {
		next(w, r)
		return
	}
	defer gz.Close()

	// Set the appropriate gzip headers.
	w.Header().Set(headerContentEncoding, encodingGzip)
	w.Header().Set(headerVary, headerAcceptEncoding)

	// Wrap the original http.ResponseWriter with negroni.ResponseWriter
	// and create the gzipResponseWriter.
	nrw := negroni.NewResponseWriter(w)
	grw := gzipResponseWriter{
		gz,
		nrw,
	}

	// Call the next gzipHandler supplying the gzipResponseWriter instead of
	// the original.
	next(grw, r)

	// Delete the content length after we know we have been written to.
	grw.Header().Del(headerContentLength)
}

// gzipResponseWriter is the ResponseWriter that negroni.ResponseWriter is
// wrapped in.
type gzipResponseWriter struct {
	w *gzip.Writer
	negroni.ResponseWriter
}

// Write writes bytes to the gzip.Writer. It will also set the Content-Type
// header using the net/http library content type detection if the Content-Type
// header was not set yet.
func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	if len(grw.Header().Get(headerContentType)) == 0 {
		grw.Header().Set(headerContentType, http.DetectContentType(b))
	}
	return grw.w.Write(b)
}
