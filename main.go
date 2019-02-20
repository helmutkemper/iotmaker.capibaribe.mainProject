package main

import (
	capib "./capibaribe"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	var err error
	var wg sync.WaitGroup

	filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
	flag.Parse()

	config := capib.MainConfig{}
	err = config.Unmarshal(*filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for projectName, projectConfig := range config.AffluentRiver {

		_ = projectName

		wg.Add(1)
		go func(config capib.Project) {
			var err error
			server := http.NewServeMux()

			defer wg.Done()

			for _, staticPath := range config.Static {

				if _, err = os.Stat(staticPath.FilePath); os.IsNotExist(err) {
					log.Fatalf("static dir error: %v\n", err.Error())
				}

				server.Handle("/"+staticPath.ServerPath+"/", http.StripPrefix("/"+staticPath.ServerPath+"/", http.FileServer(http.Dir(staticPath.FilePath))))
			}

			server.HandleFunc("/", config.HandleFunc)

			if config.Sll.Enabled == true {

				/*
					cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
					if err != nil {
							log.Println(err)
							return
					}

					config := &tls.Config{Certificates: []tls.Certificate{cer}}
				*/

				var tlsMinVersion uint16 = 0
				if config.Sll.Version.Min == 10 {
					tlsMinVersion = tls.VersionTLS10
				} else if config.Sll.Version.Min == 11 {
					tlsMinVersion = tls.VersionTLS11
				} else if config.Sll.Version.Min == 12 {
					tlsMinVersion = tls.VersionTLS12
				} else if config.Sll.Version.Min == 30 {
					tlsMinVersion = tls.VersionSSL30
				}

				var tlsMaxVersion uint16 = 0
				if config.Sll.Version.Max == 10 {
					tlsMaxVersion = tls.VersionTLS10
				} else if config.Sll.Version.Max == 11 {
					tlsMaxVersion = tls.VersionTLS11
				} else if config.Sll.Version.Max == 12 {
					tlsMaxVersion = tls.VersionTLS12
				} else if config.Sll.Version.Max == 30 {
					tlsMaxVersion = tls.VersionSSL30
				}

				var curveIdList = make([]tls.CurveID, len(config.Sll.CurvePreferences.([]string)))
				for k, v := range config.Sll.CurvePreferences.([]string) {
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

				var cipherSuitesList = make([]uint16, len(config.Sll.CipherSuites.([]string)))
				for k, v := range config.Sll.CurvePreferences.([]string) {
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
						PreferServerCipherSuites: config.Sll.PreferServerCipherSuites,
						CipherSuites:             cipherSuitesList,
					},
					Addr:     config.Listen,
					Handler:  server,
					ErrorLog: log.New(capib.DebugLogger{}, "", 0), //fixme: tem que ser opcional
				}

				/*
					srv := &http.Server{
						Addr:         ":443",
						Handler:      mux,
						TLSConfig:    cfg,
						TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
					}
					log.Fatal(srv.ListenAndServeTLS("tls.crt", "tls.key"))
				*/

				log.Fatal(newServer.ListenAndServeTLS(config.Sll.Certificate, config.Sll.CertificateKey))

			} else {

				newServer := &http.Server{
					Addr:     config.Listen,
					Handler:  server,
					ErrorLog: log.New(capib.DebugLogger{}, "", 0), //fixme: tem que ser opcional
				}
				log.Fatal(newServer.ListenAndServe())

			}

		}(projectConfig)

	}

	wg.Wait()
}
