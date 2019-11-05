package capibaribe

import (
	"github.com/helmutkemper/seelog"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"
)

type Project struct {
	ListenAndServer   ListenAndServer `yaml:"listenAndServer"   json:"listenAndServer"`
	Sll               ssl             `yaml:"ssl"               json:"ssl"`
	Proxy             []proxy         `yaml:"proxy"             json:"proxy"`
	DebugServerEnable bool            `yaml:"debugServerEnable" json:"debugServerEnable"`
	Listen            Listen          `yaml:"-"                 json:"-"`
	waitGroup         int             `yaml:"-"                 json:"-"`
}

func (el *Project) WaitAddDelta() {
	el.waitGroup += 1
}

func (el *Project) WaitDone() {
	el.waitGroup -= 1
}

func (el *Project) HandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	var host = r.Host
	var re *regexp.Regexp
	var hostServer string
	var serverKey int
	var loopCounter = 0

	el.WaitAddDelta()

	defer el.WaitDone()

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

					// Check the maximum number of interactions of the route from the proxy to prevent an infinite loop
					loopCounter += 1
					if loopCounter > el.Proxy[proxyKey].MaxAttemptToRescueLoop {
						_ = seelog.Critical("todas as rotas deram erro. final.")
						// fixme: colocar o que fazer no erro de todas as rotas
						return
					}

					// Select a type of load balancing to be applied on the proxy route
					if proxyData.LoadBalancing == kLoadBalanceRoundRobin || proxyData.LoadBalancing == "" {
						hostServer, serverKey = proxyData.roundRobin()
					} else if proxyData.LoadBalancing == kLoadBalanceRandom {
						hostServer, serverKey = proxyData.random()
					}

					// Prepare the reverse proxy
					rpURL, err := url.Parse(hostServer)
					if err != nil {
						HandleCriticalError(err)
					}

					proxy := httputil.NewSingleHostReverseProxy(rpURL)

					proxy.ErrorLog = log.New(DebugLogger{}, "", 0)

					proxy.Transport = &transport{RoundTripper: http.DefaultTransport, Project: el}

					//todo: implementar
					//proxy.ModifyResponse = proxyData.ModifyResponse

					// Prepare the statistics of the errors and successes of the route in the reverse proxy
					proxy.ErrorHandler = el.Proxy[proxyKey].ErrorHandler

					el.Proxy[proxyKey].lastRoundError = false
					el.Proxy[proxyKey].Servers[serverKey].lastRoundError = false

					startTime := time.Now()

					// Run the route
					proxy.ServeHTTP(w, r)

					elapsedTime := time.Since(startTime)
					log.Printf("execution time: %s", elapsedTime)

					// Verify error
					if el.Proxy[proxyKey].lastRoundError == true {

						el.Proxy[proxyKey].Servers[serverKey].lastRoundError = true
						el.Proxy[proxyKey].Servers[serverKey].consecutiveErrors = 0
						el.Proxy[proxyKey].Servers[serverKey].consecutiveSuccess += 1

						_ = seelog.Critical("todas as rotas deram erro. testando novamente")

						// Prepare to select a new route after error
						continue
					}

					// Statistics of successes of the route
					el.Proxy[proxyKey].SuccessHandler(w, r)
					_ = seelog.Critical("rota ok")
					return

				}
			}
		}
	}
}
