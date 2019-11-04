package capibaribe

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sync"
)

type Project struct {
	ListenAndServer   ListenAndServer `yaml:"listenAndServer"   json:"listenAndServer"`
	Sll               ssl             `yaml:"ssl"               json:"ssl"`
	Proxy             []proxy         `yaml:"proxy"             json:"proxy"`
	DebugServerEnable bool            `yaml:"debugServerEnable" json:"debugServerEnable"`
	Listen            Listen          `yaml:"-"                 json:"-"`
	waitGroup         sync.WaitGroup  `yaml:"-"                 json:"-"`
}

func (el *Project) WaitAddDelta() {
	el.waitGroup.Add(1)
}

func (el *Project) WaitDone() {
	el.waitGroup.Done()
}

func (el *Project) Wait() {
	el.waitGroup.Wait()
}

func (el *Project) HandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	var host = r.Host
	var re *regexp.Regexp
	var hostServer string
	var serverKey int
	var loopCounter = 0

	//el.waitGroup.Add(1)

	//defer el.waitGroup.Done()

	if el.Proxy != nil {

		for proxyKey, proxyData := range el.Proxy {

			if proxyData.IgnorePort == true {
				if re, err = regexp.Compile(kIgnorePortRegExp); err != nil {
					HandleCriticalError(err)
				}

				host = re.ReplaceAllString(host, "$1")
			}

			if proxyData.Host == host || proxyData.Host == "" {

				for {

					loopCounter += 1
					if loopCounter > el.Proxy[proxyKey].MaxAttemptToRescueLoop {
						// fixme: colocar o que fazer no erro de todas as rotas
						return
					}

					if proxyData.LoadBalancing == kLoadBalanceRoundRobin || proxyData.LoadBalancing == "" {

						hostServer, serverKey = proxyData.roundRobin()

					} else if proxyData.LoadBalancing == kLoadBalanceRandom {

						hostServer, serverKey = proxyData.random()

					}

					if hostServer != "" {

						rpURL, err := url.Parse(hostServer)
						if err != nil {
							HandleCriticalError(err)
						}

						proxy := httputil.NewSingleHostReverseProxy(rpURL)

						//						proxy.ErrorLog = log.New(DebugLogger{}, "", 0)

						proxy.Transport = &transport{RoundTripper: http.DefaultTransport, Project: el}

						//todo: implementar
						//proxy.ModifyResponse = proxyData.ModifyResponse

						proxy.ErrorHandler = el.Proxy[proxyKey].ErrorHandler

						el.Proxy[proxyKey].Servers[serverKey].lastRoundError = false

						proxy.ServeHTTP(w, r)

						if el.Proxy[proxyKey].Servers[serverKey].lastRoundError == true {

							el.Proxy[proxyKey].consecutiveErrors = 0
							el.Proxy[proxyKey].consecutiveSuccess += 1
							el.Proxy[proxyKey].Servers[serverKey].consecutiveErrors = 0
							el.Proxy[proxyKey].Servers[serverKey].consecutiveSuccess += 1

							//seelog.Critical("continue")
							continue
						}

						//seelog.Critical("return")
						return

					} else {

						//fixme: colocar um log aqui

					}
				}
			}
		}
	}
}
