package capibaribe

import (
	"github.com/helmutkemper/seelog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sync"
)

type Project struct {
	ListenAndServer   ListenAndServer `yaml:"listenAndServer"`
	Sll               ssl             `yaml:"ssl"`
	Protocol          string          `yaml:"protocol"`
	Pygocentrus       pygocentrus     `yaml:"pygocentrus"`
	Proxy             []proxy         `yaml:"proxy"`
	Static            []static        `yaml:"static"`
	DebugServerEnable bool            `yaml:"debugServerEnable"`
	Listen            Listen          `yaml:"-"`
	waitGroup         sync.WaitGroup  `yaml:"-"`
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
	var remoteAddr string
	var re *regexp.Regexp
	var hostServer string
	var serverKey int
	var loopCounter = 0

	//el.waitGroup.Add(1)

	//defer el.waitGroup.Done()

	if el.Proxy != nil {

		for proxyKey, proxyData := range el.Proxy {

			pass := len(proxyData.Bind) == 0
			for _, bind := range proxyData.Bind {
				if bind.IgnorePort == true {
					if re, err = regexp.Compile(kIgnorePortRegExp); err != nil {
						HandleCriticalError(err)
					}

					remoteAddr = re.ReplaceAllString(r.RemoteAddr, "$1")
				}

				if remoteAddr == bind.Host {
					pass = true
					break
				}
			}
			if pass == false {
				seelog.Debugf("remote address ( %v ) not in bind list\n", r.RemoteAddr)
				return
			}

			if proxyData.IgnorePort == true {
				if re, err = regexp.Compile(kIgnorePortRegExp); err != nil {
					HandleCriticalError(err)
				}

				host = re.ReplaceAllString(host, "$1")
			}

			if proxyData.Host == host {

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

							if el.Pygocentrus.GetAttack() == true {
								el.Pygocentrus.ClearAttack()
								//seelog.Critical("return after a pygocentrus attack")
								return
							}

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
