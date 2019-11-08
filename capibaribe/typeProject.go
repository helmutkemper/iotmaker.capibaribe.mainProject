package capibaribe

import (
	"github.com/helmutkemper/seelog"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
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
	var path = r.URL.Path
	var re *regexp.Regexp
	var hostServer string
	var serverKey int
	var loopCounter = 0

	el.WaitAddDelta()

	defer el.WaitDone()

	if el.Proxy != nil {

		for proxyKey, proxyData := range el.Proxy {

			if proxyData.IgnorePort == true {
				if re, err = regexp.Compile(KIgnorePortRegExp); err != nil {
					HandleCriticalError(err)
					return
				}

				host = re.ReplaceAllString(host, "$1")
			}

			if proxyData.VerifyHostPathToValidateRoute(host) == true {

				if proxyData.VerifyHealthCheckPathToValidatePathIntoHost(path) == true {
					proxyData.WriteHealthCheckDataToOutputEndpoint(w, r)
					return
				}

				if proxyData.VerifyRouteDataPathToValidatePathIntoHost(path) == true {
					proxyData.WriteProxyDataToOutputJSonEndpoint(w, r)
					return
				}

				A := proxyData.VerifyPathAndHeaderInformationToValidateRoute(path, w, r)
				B := proxyData.VerifyPathWithoutVerifyHeaderInformationToValidateRoute(path)
				C := proxyData.VerifyHeaderInformationWithoutVerifyPathToValidateRoute(w, r)

				// simplified true table
				// | A | B | C | S |
				// |---|---|---|---|
				// | X | X | 1 | 1 |
				// | X | 1 | X | 1 |
				// | 1 | X | X | 1 |
				// | 0 | 0 | 0 | 0 |
				if !(A || B || C) {
					continue
				}

				for {

					// Check the maximum number of interactions of the route from the proxy to prevent an infinite loop
					loopCounter += 1
					if loopCounter > el.Proxy[proxyKey].MaxAttemptToRescueLoop {
						_ = seelog.Critical("todas as rotas deram erro. final.")
						// fixme: colocar o que fazer no erro de todas as rotas
						return
					}

					hostServer, serverKey = proxyData.SelectLoadBalance()

					// Prepare the reverse proxy
					rpURL, err := url.Parse(hostServer)
					if err != nil {
						HandleCriticalError(err)
					}

					proxy := httputil.NewSingleHostReverseProxy(rpURL)
					proxy.ErrorLog = log.New(DebugLogger{}, "", 0)
					proxy.Transport = &transport{RoundTripper: http.DefaultTransport, Project: el}
					// Prepare the statistics of the TotalErrorsCounter and successes of the route in the reverse proxy
					proxy.ErrorHandler = el.Proxy[proxyKey].OnErrorHandlerEvent

					el.Proxy[proxyKey].lastRoundError = false

					//todo: implementar
					//proxy.ModifyResponse = proxyData.ModifyResponse

					// Run the route and measure execution time
					el.Proxy[proxyKey].Servers[serverKey].OnExecutionStartEvent()
					el.Proxy[proxyKey].OnExecutionStartEvent()
					proxy.ServeHTTP(w, r)

					// Verify error and continue to select a new route in case of error
					if el.Proxy[proxyKey].lastRoundError == true {
						el.Proxy[proxyKey].Servers[serverKey].OnExecutionEndWithErrorEvent()
						_ = seelog.Critical("todas as rotas deram erro. testando novamente")
						continue
					}

					// Statistics of successes of the route
					el.Proxy[proxyKey].OnSuccessHandlerEvent()
					el.Proxy[proxyKey].Servers[serverKey].OnExecutionEndWithSuccessEvent()
					_ = seelog.Critical("rota ok")
					return

				}
			}
		}
	}
}
