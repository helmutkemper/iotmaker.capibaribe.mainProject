package capibaribe

import (
	"bytes"
	"github.com/helmutkemper/seelog"
	"github.com/helmutkemper/yaml"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	kIgnorePortRegExp      = "^(.*?):[0-9]+$"
	kLoadBalanceRoundRobin = "roundRobin"
	kLoadBalanceRandom     = "random"

	kPygocentrusDontRespond   = 0
	kPygocentrusChangeLength  = 1
	kPygocentrusChangeContent = 2
	kPygocentrusDeleteContent = 3
	kPygocentrusChangeHeaders = 4
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

type Project struct {
	Listen      string      `yaml:"listen"`
	Sll         ssl         `yaml:"ssl"`
	Pygocentrus pygocentrus `yaml:"pygocentrus"`
	Proxy       []proxy     `yaml:"proxy"`
	Static      []static    `yaml:"static"`
}

func (el *Project) HandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	var host = r.Host
	var remoteAddr string
	var re *regexp.Regexp
	var hostServer string
	var serverKey int
	var loopCounter = 0

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
					if loopCounter > 10 {
						// fixme: colocar o que fazer no erro de todas as rotas
						// fixme: o valor 10 deve ser configurado no sistema
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

						proxy.ErrorLog = log.New(DebugLogger{}, "", 0)

						//if el.Pygocentrus.Enabled == true {
						proxy.Transport = &transport{RoundTripper: http.DefaultTransport, Project: el}
						//}
						proxy.ModifyResponse = proxyData.ModifyResponse

						proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {

							//w.WriteHeader(500)
							el.Proxy[proxyKey].consecutiveErrors += 1
							el.Proxy[proxyKey].consecutiveSuccess = 0
							el.Proxy[proxyKey].Servers[serverKey].consecutiveErrors += 1
							el.Proxy[proxyKey].Servers[serverKey].consecutiveSuccess = 0
							el.Proxy[proxyKey].Servers[serverKey].errors += 1
							el.Proxy[proxyKey].Servers[serverKey].lastRoundError = true

							seelog.Criticalf("1 server host %v error - %v", hostServer, err.Error())
						}

						el.Proxy[proxyKey].Servers[serverKey].lastRoundError = false

						proxy.ServeHTTP(w, r)

						if el.Proxy[proxyKey].Servers[serverKey].lastRoundError == true {

							el.Proxy[proxyKey].consecutiveErrors = 0
							el.Proxy[proxyKey].consecutiveSuccess += 1
							el.Proxy[proxyKey].Servers[serverKey].consecutiveErrors = 0
							el.Proxy[proxyKey].Servers[serverKey].consecutiveSuccess += 1
							seelog.Critical("continue")
							continue
						}

						seelog.Critical("return")
						return

					} else {

						//fixme: colocar um log aqui

					}

				}

			}

		}

	}

}

type static struct {
	FilePath   string `yaml:"filePath"`
	ServerPath string `yaml:"serverPath"`
}

type ssl struct {
	Enabled        bool     `yaml:"enabled"`
	Certificate    string   `yaml:"certificate"`
	CertificateKey string   `yaml:"certificateKey"`
	SllProtocols   []string `yaml:"sslProtocols"`
}

type pygocentrus struct {
	Enabled       bool            `yaml:"enabled"`
	DontRespond   float64         `yaml:"dontRespond"`
	ChangeLength  float64         `yaml:"changeLength"`
	ChangeContent float64         `yaml:"changeContent"`
	DeleteContent float64         `yaml:"deleteContent"`
	ChangeHeaders []changeHeaders `yaml:"changeHeaders"`
}

type changeHeaders struct {
	Number int      `yaml:"number"`
	Header []header `yaml:"header"`
	Rate   float64  `yaml:"rate"`
}

type header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type proxy struct {
	consecutiveErrors  int
	consecutiveSuccess int
	errors             int
	success            int

	lastRoundError bool

	IgnorePort    bool        `yaml:"ignorePort"`
	Host          string      `yaml:"host"`
	Bind          []bind      `yaml:"bind"`
	LoadBalancing string      `yaml:"loadBalancing"`
	Path          string      `yaml:"path"`
	HealthCheck   healthCheck `yaml:"healthCheck"`
	Servers       []servers   `yaml:"servers"`
}

func (el *proxy) ModifyResponse(resp *http.Response) error {
	//seelog.Criticalf("header code %v", resp.StatusCode)
	return nil
}

func (el *proxy) roundRobin() (string, int) {
	randNumber := rand.Float64()

	for serverKey, serverData := range el.Servers {

		if randNumber <= serverData.Weight {
			return serverData.Host, serverKey
		}

	}

	return "", -1
}

func (el *proxy) random() (string, int) {
	randNumber := rand.Intn(len(el.Servers))
	return el.Servers[randNumber].Host, randNumber
}

type bind struct {
	Host       string `yaml:"host"`
	IgnorePort bool   `yaml:"ignorePort"`
}

type servers struct {
	consecutiveErrors  int
	consecutiveSuccess int
	errors             int
	success            int

	lastRoundError bool

	Host     string  `yaml:"host"`
	Weight   float64 `yaml:"weight"`
	OverLoad int     `yaml:"overLoad"`
}

type healthCheck struct {
	Enabled         bool  `yaml:"enabled"`
	Interval        int   `yaml:"interval"`
	Fails           int   `yaml:"fails"`
	Passes          int   `yaml:"passes"`
	Uri             int   `yaml:"rui"`
	SuspendInterval int   `yaml:"suspendInterval"`
	Match           match `yaml:"match"`
}

type match struct {
	Status []status `yaml:"status"`
	Header []header `yaml:"header"`
	Body   []string `yaml:"body"`
}

type status struct {
	ExpReg string   `yaml:"expReg"`
	Value  int      `yaml:"value"`
	In     []maxMin `yaml:"in"`
	NotIn  []maxMin `yaml:"notIn"`
}

type maxMin struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type transport struct {
	RoundTripper http.RoundTripper
	Project      *Project
}

func (el *transport) roundTripReadBody(req *http.Request) ([]byte, error) {
	var resp *http.Response
	var err error
	var inBody []byte

	resp, err = el.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	inBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return inBody, err
}

func (el *transport) roundTripCopyBody(inBody []byte) io.ReadCloser {
	inBody = make([]byte, len(inBody))
	return ioutil.NopCloser(bytes.NewReader(inBody))
}

func (el *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	if el.Project.Pygocentrus.Enabled == true {

		randAttack := rand.Intn(4)

		if randAttack == kPygocentrusDontRespond {
			if el.Project.Pygocentrus.DontRespond >= rand.Float64() {

				return nil, errors.New("this data were eaten by a pygocentrus attack")

			}
		}

		if randAttack == kPygocentrusDeleteContent {
			if el.Project.Pygocentrus.DeleteContent >= rand.Float64() {

				var inBody []byte
				inBody, err = el.roundTripReadBody(req)
				if err != nil {
					return nil, err
				}

				inBody = make([]byte, len(inBody))

				resp.Body = el.roundTripCopyBody(inBody)
				return resp, nil

			}
		}

		if randAttack == kPygocentrusChangeContent {
			if el.Project.Pygocentrus.ChangeContent >= rand.Float64() {

				var inBody []byte
				inBody, err = el.roundTripReadBody(req)
				if err != nil {
					return nil, err
				}

				for i := 0; i != rand.Intn(len(inBody)); i += 1 {
					inBody = append(append(inBody[:i], byte(rand.Intn(255))), inBody[i+1:]...)
				}

				resp.Body = el.roundTripCopyBody(inBody)
				return resp, nil

			}
		}

		if randAttack == kPygocentrusChangeLength {
			if el.Project.Pygocentrus.ChangeLength >= rand.Float64() {

				var inBody []byte
				inBody, err = el.roundTripReadBody(req)
				if err != nil {
					return nil, err
				}
				resp.Body = el.roundTripCopyBody(inBody)

				randLength := rand.Intn(len(inBody))

				resp.ContentLength = int64(randLength)
				resp.Header.Set("Content-Length", strconv.Itoa(randLength))
				return resp, nil

			}
		}
	}

	resp, err = el.RoundTripper.RoundTrip(req)
	return resp, err
}

type DebugLogger struct{}

func (d DebugLogger) Write(p []byte) (n int, err error) {
	s := string(p)
	if strings.Contains(s, "multiple response.WriteHeader") {
		debug.PrintStack()
	}
	return os.Stderr.Write(p)
}
