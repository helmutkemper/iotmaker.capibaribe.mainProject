package capibaribe

import (
	"errors"
	"github.com/helmutkemper/yaml"
	"io/ioutil"
)

type MainConfig struct {
	Version       float64            `yaml:"version"`
	AffluentRiver map[string]Project `yaml:"capibaribe"`
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

			if proxyData.LoadBalancing == kLoadBalanceRoundRobin || proxyData.LoadBalancing == "" {

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
					WeightsSum += float64(serversData.Weight)
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
