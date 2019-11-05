package capibaribe

import (
	"errors"
	"github.com/helmutkemper/yaml"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type MainConfig struct {
	Version       float64            `yaml:"version"           json:"version"`
	AffluentRiver map[string]Project `yaml:"capibaribe"        json:"capibaribe"`
	waitGroup     sync.WaitGroup     `yaml:"-"                 json:"-"`
}

func (el *MainConfig) WaitAddDelta() {
	el.waitGroup.Add(1)
}

func (el *MainConfig) WaitDone() {
	el.waitGroup.Done()
}

func (el *MainConfig) Wait() {
	el.waitGroup.Wait()
}

func (el *MainConfig) LoadConfAndStart(filePath string) {
	var err error
	err = el.Unmarshal(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for projectName, projectConfig := range el.AffluentRiver {

		go func(name string, config Project) {

			server := http.NewServeMux()

			server.HandleFunc("/", config.HandleFunc)

			newServer := &http.Server{
				//TLSNextProto:               make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
				Addr:    config.ListenAndServer.InAddress,
				Handler: server,
			}

			if config.DebugServerEnable == true {
				newServer.ErrorLog = log.New(DebugLogger{}, "", 0)
			}

			ConfigCertificates(config.Sll, newServer)

			if config.Sll.Enabled == true {

				if config.Sll.Certificate != "" && config.Sll.CertificateKey != "" {

					log.Fatal(newServer.ListenAndServeTLS(config.Sll.Certificate, config.Sll.CertificateKey))

				} else {
					//fixme: log de erro
				}

			} else {

				log.Fatal(newServer.ListenAndServe())

			}

		}(projectName, projectConfig)

	}
}

func (el *MainConfig) Unmarshal(filePath string) error {
	var fileContent []byte
	var err error

	fileContent, err = ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContent, el)
	if err != nil {
		return err
	}

	return el.prepare()
}

func (el *MainConfig) UnmarshalByte(fileContent []byte) error {
	var err error

	err = yaml.Unmarshal(fileContent, el)
	if err != nil {
		return err
	}

	return el.prepare()
}

func (el *MainConfig) prepare() error {

	var WeightsSum = 0.0

	if el.Version != 1.0 {
		return errors.New("this project accepts only version 1.0")
	}

	for affluentKey := range el.AffluentRiver {

		for proxyKey, proxyData := range el.AffluentRiver[affluentKey].Proxy {

			if el.AffluentRiver[affluentKey].Proxy[proxyKey].MaxAttemptToRescueLoop == 0 {
				el.AffluentRiver[affluentKey].Proxy[proxyKey].MaxAttemptToRescueLoop = 10
			}

			// fixme: por que o tipo de loadbalance está sendo levado em conta aqui?
			if proxyData.LoadBalancing == KLoadBalanceRoundRobin || proxyData.LoadBalancing == "" {

				pass := false
				for serverKey := range el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers {

					if el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers[serverKey].Weight != 0 {
						pass = true
					}

				}

				if pass == false {

					for serverKey := range el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers {
						el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers[serverKey].Weight = 1
					}

				}

				for _, serversData := range el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers {
					WeightsSum += serversData.Weight
				}

				for serversKey, serversData := range el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers {
					if serversKey == 0 {
						el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers[serversKey].Weight = serversData.Weight / WeightsSum
					} else {
						el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers[serversKey].Weight = (serversData.Weight / WeightsSum) + el.AffluentRiver[affluentKey].Proxy[proxyKey].Servers[serversKey-1].Weight
					}

				}

			}

		}

	}

	return nil
}
