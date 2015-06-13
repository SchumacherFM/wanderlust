// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.
// Cyrill: created a struct with a method to use it internally
package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/juju/errgo"
)

// GeneratePems creates new PEM files if they not yet exists
// returns the name of the PEM directory and error|nil
func GeneratePems(listenAddress, pemDir, certFileName, keyFileName string) (string, error) {
	dir := pemDir
	pathSep := string(os.PathSeparator)
	if "" == dir {
		dir = GetTempDir() + "wlpem_" + RandomString(10)
	}
	CreateDirectoryIfNotExists(dir)

	if len(dir) > 0 && pathSep != dir[len(dir)-1:] {
		dir = dir + pathSep
	}

	certFile := dir + certFileName
	keyFile := dir + keyFileName
	address, _, vlaErr := ValidateListenAddress(listenAddress)
	if nil != vlaErr {
		return "", errgo.Mask(vlaErr)
	}
	duration := 180 * 24 * time.Hour                      // half a year
	validFrom := time.Now().Format("Jan 2 15:04:05 2006") // any other year number results in a weird result :-?

	isDir1, _ := PathExists(certFile)
	isDir2, _ := PathExists(keyFile)
	if isDir1 && isDir2 {
		return dir, nil
	}
	certGenerator := &generateCert{
		Host:         address,   // "Comma-separated hostnames and IPs to generate a certificate for"
		ValidFrom:    validFrom, // "Creation date formatted as Jan 1 15:04:05 2011"
		ValidFor:     duration,  // "duration", 365*24*time.Hour, "Duration that certificate is valid for"
		IsCA:         true,      // "ca", false, "whether this cert should be its own Certificate Authority"
		RsaBits:      2048,      // "rsa-bits", 2048, "Size of RSA key to generate"
		CertFileName: certFile,
		KeyFileName:  keyFile,
	}
	return dir, errgo.Mask(certGenerator.Generate())
}

type generateCert struct {
	Host         string        // flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	ValidFrom    string        // flag.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	ValidFor     time.Duration // flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
	IsCA         bool          // flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
	RsaBits      int           // flag.Int("rsa-bits", 2048, "Size of RSA key to generate")
	CertFileName string
	KeyFileName  string
}

func (gc *generateCert) Generate() error {

	if len(gc.Host) == 0 {
		return errors.New("Missing required host parameter")
	}

	priv, err := rsa.GenerateKey(rand.Reader, gc.RsaBits)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to generate private key: %s", err))
	}

	var notBefore time.Time
	if len(gc.ValidFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", gc.ValidFrom)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to parse creation date: %s\n", err))
		}
	}

	notAfter := notBefore.Add(gc.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to generate serial number: %s", err))
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Wanderlust Co"},
			Country:      []string{"AU"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(gc.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if gc.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create certificate: %s", err))
	}

	certOut, err := os.Create(gc.CertFileName)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to open "+gc.CertFileName+" for writing: %s", err))
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(gc.KeyFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to open "+gc.KeyFileName+" for writing:", err))
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	return nil
}
