package capibaribe

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
)

func ConfigCertificates(config ssl, server *http.ServeMux) {
	var err error

	if config.Enabled == true {

		var certificatesList = tls.Certificate{}

		if config.X509.Certificate != "" && config.X509.CertificateKey != "" {

			if _, err = os.Stat(config.X509.Certificate); os.IsNotExist(err) {
				log.Fatalf("sll x509 certificate error: %v\n", err.Error())
			}

			if _, err = os.Stat(config.X509.CertificateKey); os.IsNotExist(err) {
				log.Fatalf("sll x509 certificate key error: %v\n", err.Error())
			}

			certificatesList, err = tls.LoadX509KeyPair(config.X509.Certificate, config.X509.CertificateKey)
			if err != nil {
				log.Fatalf("sll x509 certificate load pair error: %v\n", err.Error())
			}

		}

		var tlsMinVersion uint16 = 0
		if config.Version.Min == 10 {
			tlsMinVersion = tls.VersionTLS10
		} else if config.Version.Min == 11 {
			tlsMinVersion = tls.VersionTLS11
		} else if config.Version.Min == 12 {
			tlsMinVersion = tls.VersionTLS12
		} else if config.Version.Min == 30 {
			tlsMinVersion = tls.VersionSSL30
		}

		var tlsMaxVersion uint16 = 0
		if config.Version.Max == 10 {
			tlsMaxVersion = tls.VersionTLS10
		} else if config.Version.Max == 11 {
			tlsMaxVersion = tls.VersionTLS11
		} else if config.Version.Max == 12 {
			tlsMaxVersion = tls.VersionTLS12
		} else if config.Version.Max == 30 {
			tlsMaxVersion = tls.VersionSSL30
		}

		var curveIdList = make([]tls.CurveID, len(config.CurvePreferences.([]string)))
		for k, v := range config.CurvePreferences.([]string) {
			if v == "P256" {
				curveIdList[k] = tls.CurveP256
			} else if v == "P384" {
				curveIdList[k] = tls.CurveP384
			} else if v == "P521" {
				curveIdList[k] = tls.CurveP521
			} else if v == "X25519" {
				curveIdList[k] = tls.X25519
			}
		}

		var cipherSuitesList = make([]uint16, len(config.CipherSuites.([]string)))
		for k, v := range config.CurvePreferences.([]string) {
			if v == "TLS_RSA_WITH_RC4_128_SHA" {
				cipherSuitesList[k] = tls.TLS_RSA_WITH_RC4_128_SHA
			} else if v == "TLS_RSA_WITH_3DES_EDE_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA
			} else if v == "TLS_RSA_WITH_AES_128_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_128_CBC_SHA
			} else if v == "TLS_RSA_WITH_AES_256_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_256_CBC_SHA
			} else if v == "TLS_RSA_WITH_AES_128_CBC_SHA256" {
				cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_128_CBC_SHA256
			} else if v == "TLS_RSA_WITH_AES_128_GCM_SHA256" {
				cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_128_GCM_SHA256
			} else if v == "TLS_RSA_WITH_AES_256_GCM_SHA384" {
				cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_256_GCM_SHA384
			} else if v == "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA" {
				cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA
			} else if v == "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
			} else if v == "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
			} else if v == "TLS_ECDHE_RSA_WITH_RC4_128_SHA" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA
			} else if v == "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
			} else if v == "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
			} else if v == "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
			} else if v == "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256" {
				cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
			} else if v == "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
			} else if v == "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
			} else if v == "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256" {
				cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
			} else if v == "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
			} else if v == "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384" {
				cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
			} else if v == "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305" {
				cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
			} else if v == "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305" {
				cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
			} else if v == "TLS_FALLBACK_SCSV" {
				cipherSuitesList[k] = tls.TLS_FALLBACK_SCSV
			}
		}

		newServer := &http.Server{
			TLSConfig: &tls.Config{
				MinVersion:               tlsMinVersion,
				MaxVersion:               tlsMaxVersion,
				CurvePreferences:         curveIdList,
				PreferServerCipherSuites: config.PreferServerCipherSuites,
				CipherSuites:             cipherSuitesList,
				Certificates:             []tls.Certificate{certificatesList},
			},
			//TLSNextProto:               make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
			Addr:    config.Listen,
			Handler: server,
		}

		if config.DebugServerEnable == true {
			newServer.ErrorLog = log.New(capib.DebugLogger{}, "", 0)
		}

		if config.Certificate != "" && config.CertificateKey != "" {

			log.Fatal(newServer.ListenAndServeTLS(config.Certificate, config.CertificateKey))

		} else {

			log.Fatal(newServer.ListenAndServe())

		}

	} else {

		newServer := &http.Server{
			Addr:    config.Listen,
			Handler: server,
		}

		if config.DebugServerEnable == true {
			newServer.ErrorLog = log.New(capib.DebugLogger{}, "", 0)
		}

		log.Fatal(newServer.ListenAndServe())

	}

}
