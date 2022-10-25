package tls

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"time"
)

// CipherSuites without known attacks or extreme CPU usage
// https://golang.org/src/crypto/tls/cipher_suites.go#L75
var CipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

	// Go 1.8 only
	// tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	// tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,

	// Best disabled, as they don't provide Forward Secrecy,
	// but might be necessary for some clients
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256,

	// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_RSA_WITH_AES_256_CBC_SHA,
}

// Curves without known attacks or extreme CPU usage
// https://golang.org/src/crypto/tls/common.go#L542
var Curves = []tls.CurveID{
	// Only use curves which have assembly implementations
	tls.CurveP256,
	// tls.X25519, // Go 1.8 only
	// tls.CurveP384,
	// tls.CurveP521,
}

// TLSConfig for including autocert manager
func TLSConfig(servername string) *tls.Config {
	certs, err := GetCertificate()
	cas := GetCertificateAuthorities()
	if err != nil {
		panic(err.(any))
	}

	return &tls.Config{
		MaxVersion:       tls.VersionTLS12,
		Certificates:     []tls.Certificate{*certs},
		RootCAs:          cas,
		ServerName:       servername,
		CurvePreferences: Curves,
		CipherSuites:     CipherSuites,
	}
}

func GetCertificateAuthorities() *x509.CertPool {
	caCert, err := os.ReadFile("./certs/cas.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool
}

// GetCertificate using autocert
func GetCertificate() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair("./certs/tls.crt", "./certs/tls.key")
	if err != nil {
		return &tls.Certificate{}, err
	}
	return &cert, nil
}

// GetHTTPSServer fully secured
func GetHTTPSServer(servername string, addr string) (s *http.Server) {
	tlsConfig := TLSConfig(servername)
	s = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second, // go 1.8

		Addr:      addr,
		TLSConfig: tlsConfig,
	}

	return s
}
