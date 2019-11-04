//todo: tem que ter um log avisando quando foi a ultima requisição
//todo: microsservico deveria ter um id inicial que permita seguir uma requisição em caso de erro
package main

import (
	capib "./capibaribe"

	"flag"
	//"fmt"
	//"github.com/etcd-io/etcd/clientv3"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

var config capib.MainConfig

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	//mutexTerminateServer = make( map[string]*sync.WaitGroup )

	filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
	flag.Parse()

	loadConf(*filePath)

	//etcdTest()

	wg.Wait()
}

func loadConf(filePath string) {
	var err error
	err = config.Unmarshal(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for projectName, projectConfig := range config.AffluentRiver {

		// fixme: isto tem lógica se ficar em branco?
		if projectConfig.ListenAndServer.InAddress != "" {

			go func(name string, config capib.Project) {

				server := http.NewServeMux()

				capib.ConfigStatic(&config.Static, server)

				server.HandleFunc("/", config.HandleFunc)

				newServer := &http.Server{
					//TLSNextProto:               make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
					Addr:    config.ListenAndServer.InAddress,
					Handler: server,
				}

				if config.DebugServerEnable == true {
					newServer.ErrorLog = log.New(capib.DebugLogger{}, "", 0)
				}

				capib.ConfigCertificates(config.Sll, newServer)

				// fixme: melhorar isto
				// enabled == true e sem certificados é um erro
				if config.Sll.Enabled == true && config.Sll.Certificate != "" && config.Sll.CertificateKey != "" {

					log.Fatal(newServer.ListenAndServeTLS(config.Sll.Certificate, config.Sll.CertificateKey))

				} else {

					log.Fatal(newServer.ListenAndServe())

				}

			}(projectName, projectConfig)

			// todo: pygocentrus não deveria está aqui
		} else if projectConfig.Listen.InAddress != "" {

			//projectConfig.Listen.Pygocentrus = projectConfig.Pygocentrus
			log.Fatal(projectConfig.Listen.Listen())

		}
	}
}
