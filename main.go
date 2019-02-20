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

				newServer := &http.Server{
					TLSConfig: &tls.Config{
						MinVersion: tls.VersionTLS10,
						MaxVersion: tls.VersionSSL30,
						/*CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
						PreferServerCipherSuites: true,
						CipherSuites: []uint16{
							tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
							tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
							tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
							tls.TLS_RSA_WITH_AES_256_CBC_SHA,
						},*/
					},
					Addr:     config.Listen,
					Handler:  server,
					ErrorLog: log.New(capib.DebugLogger{}, "", 0),
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
					ErrorLog: log.New(capib.DebugLogger{}, "", 0),
				}
				log.Fatal(newServer.ListenAndServe())

			}

		}(projectConfig)

	}

	wg.Wait()
}
