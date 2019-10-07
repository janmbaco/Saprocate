package tlsdial

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/janmbaco/Saprocate/config"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
)

func GetConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair(config.Config.ServerCertificate.PublicKeyFile, config.Config.ServerCertificate.PrivateKeyFile)
	cross.TryPanic(err)
	certpool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(config.Config.AuthorizeCertificate.PublicKeyFile)
	cross.TryPanic(err)
	certpool.AppendCertsFromPEM(pem)

	return &tls.Config{
		Rand:                  rand.Reader,
		Time:                  nil,
		Certificates:          []tls.Certificate{cert},
		NameToCertificate:     nil,
		GetCertificate:        nil,
		GetClientCertificate:  nil,
		GetConfigForClient:    nil,
		VerifyPeerCertificate: nil,
		RootCAs:               nil,
		NextProtos:            nil,
		ServerName:            "",
		ClientAuth:            tls.RequireAndVerifyClientCert,
		ClientCAs:             certpool,
		InsecureSkipVerify:    false,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites:    true,
		SessionTicketsDisabled:      false,
		SessionTicketKey:            [32]byte{},
		ClientSessionCache:          nil,
		MinVersion:                  tls.VersionTLS12,
		MaxVersion:                  0,
		CurvePreferences:            []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		DynamicRecordSizingDisabled: false,
		Renegotiation:               0,
		KeyLogWriter:                nil,
	}
}
