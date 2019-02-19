package capibaribe

import (
	"bytes"
	"errors"
	"github.com/helmutkemper/seelog"
	"github.com/helmutkemper/yaml"
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
	"strings"
	"time"
)

const (
	kIgnorePortRegExp      = "^(.*?):[0-9]+$"
	kLoadBalanceRoundRobin = "roundRobin"
	kLoadBalanceRandom     = "random"
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

						proxy.ErrorLog = log.New(DebugLogger{}, "", 0)

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
							return
							if el.Pygocentrus.Enabled == true {
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
	ChangeContent changeContent   `yaml:"changeContent"`
	DeleteContent float64         `yaml:"deleteContent"`
	ChangeHeaders []changeHeaders `yaml:"changeHeaders"`
}

type changeContent struct {
	ChangeRateMin  float64 `yaml:"changeRateMin"`
	ChangeRateMax  float64 `yaml:"changeRateMax"`
	ChangeBytesMin int     `yaml:"changeBytesMin"`
	ChangeBytesMax int     `yaml:"changeBytesMax"`
	Rate           float64 `yaml:"rate"`
}

func (el *changeContent) prepare() error {
	if el.Rate == 0.0 {
		return nil
	}

	if el.ChangeRateMin == el.ChangeRateMax && el.ChangeBytesMin == el.ChangeBytesMax && el.ChangeRateMin == 0.0 {
		el.Rate = 0.0
		return errors.New("pygocentrus attack > changeContent > rate set to zero")
	}

	if el.ChangeRateMin > el.ChangeRateMax {
		return errors.New("pygocentrus attack > changeContent > rate > the minimum value is greater than the maximum value")
	}

	if el.ChangeBytesMin > el.ChangeBytesMax {
		return errors.New("pygocentrus attack > changeContent > bytes > the minimum value is greater than the maximum value")
	}

	if (el.ChangeRateMin != 0.0 || el.ChangeRateMax != 0.0) && (el.ChangeBytesMin != 0.0 || el.ChangeBytesMax != 0.0) {
		return errors.New("pygocentrus attack > changeContent > you must choose option rate change or option bytes change")
	}

	return nil
}

func (el *changeContent) GetRandomByMaxMin(length int) int {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))

	if el.ChangeRateMin != 0.0 || el.ChangeRateMax != 0.0 {
		var changeMin = int(float64(length) * el.ChangeRateMin)
		var changeMax = int(float64(length) * el.ChangeRateMax)

		return r1.Intn(changeMax-changeMin) + changeMin
	}

	return r1.Intn(el.ChangeBytesMax-el.ChangeBytesMin) + el.ChangeBytesMin
}

func (el *changeContent) GetRandomByLength(length int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(length)
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

	keyProxy  int
	keyServer int

	lastError      error
	lastRoundError bool

	MaxAttemptToRescueLoop int         `yaml:"maxAttemptToRescueLoop"`
	IgnorePort             bool        `yaml:"ignorePort"`
	Host                   string      `yaml:"host"`
	Bind                   []bind      `yaml:"bind"`
	LoadBalancing          string      `yaml:"loadBalancing"`
	Path                   string      `yaml:"path"`
	HealthCheck            healthCheck `yaml:"healthCheck"`
	Servers                []servers   `yaml:"servers"`
}

func (el *proxy) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {

	//w.WriteHeader(500)
	el.consecutiveErrors += 1
	el.consecutiveSuccess = 0
	el.consecutiveErrors += 1
	el.consecutiveSuccess = 0
	el.errors += 1
	el.lastRoundError = true
	el.lastError = err

	//seelog.Criticalf("1 server host %v error - %v", hostServer, err.Error())
}

func (el *proxy) SuccessHandler(w http.ResponseWriter, r *http.Request, err error) {

	//seelog.Criticalf("1 server host %v error - %v", hostServer, err.Error())
}

func (el *proxy) ModifyResponse(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

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

type pygocentrusFunc func(req *http.Request) (resp *http.Response, err error)

func (el *transport) roundTripReadBody(req *http.Request) (*http.Response, []byte, error) {
	var resp *http.Response
	var err error
	var inBody []byte

	resp, err = el.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, nil, err
	}

	inBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return resp, nil, err
	}

	return resp, inBody, err
}

func (el *transport) roundTripCopyBody(inBody []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(inBody))
}

func (el *transport) PygocentrusDontRespond(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: dont respond", req.RemoteAddr, req.RequestURI)
	return nil, nil

}

func (el *transport) PygocentrusDeleteContent(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: delete content", req.RemoteAddr, req.RequestURI)

	var inBody []byte
	resp, inBody, err = el.roundTripReadBody(req)
	if err != nil {
		return nil, err
	}

	inBody = make([]byte, len(inBody))

	resp.Body = el.roundTripCopyBody(inBody)
	return resp, nil

}

func (el *transport) PygocentrusChangeContent(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: change content", req.RemoteAddr, req.RequestURI)

	var inBody []byte
	resp, inBody, err = el.roundTripReadBody(req)
	if err != nil {
		return nil, err
	}

	length := len(inBody)
	forLength := el.Project.Pygocentrus.ChangeContent.GetRandomByMaxMin(length)
	for i := 0; i != forLength; i += 1 {
		indexChange := el.Project.Pygocentrus.ChangeContent.GetRandomByLength(length)
		inBody = append(append(inBody[:indexChange], byte(rand.Intn(255))), inBody[indexChange+1:]...)
	}

	resp.Body = el.roundTripCopyBody(inBody)
	return resp, nil

}

func (el *transport) PygocentrusChangeLength(req *http.Request) (resp *http.Response, err error) {

	seelog.Debugf("%v%v were eaten by a pygocentrus attack: change length", req.RemoteAddr, req.RequestURI)

	var inBody []byte
	resp, inBody, err = el.roundTripReadBody(req)
	if err != nil {
		return nil, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(inBody))

	randLength := rand.Intn(len(inBody))

	resp.ContentLength = int64(randLength)
	//resp.Header.Set("Content-Length", strconv.Itoa(randLength))
	return resp, nil

}

// todo: fazer
//func (el *transport) PygocentrusChangeHeaders(req *http.Request) (resp *http.Response, err error) {}

func (el *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	if el.Project.Pygocentrus.Enabled == true {

		var randAttack int

		var list = make([]pygocentrusFunc, 0)

		if el.Project.Pygocentrus.DontRespond != 0.0 {

			list = append(list, el.PygocentrusDontRespond)

		}

		if el.Project.Pygocentrus.DeleteContent != 0.0 {

			list = append(list, el.PygocentrusDeleteContent)

		}

		if el.Project.Pygocentrus.ChangeContent.Rate != 0.0 {

			list = append(list, el.PygocentrusChangeContent)

		}

		if el.Project.Pygocentrus.ChangeLength != 0.0 {

			list = append(list, el.PygocentrusChangeLength)

		}

		/* todo: fazer
		if el.Project.Pygocentrus.ChangeHeaders[0].Rate != 0.0 {}
		*/

		randAttack = inLineRand().Intn(len(list))
		return list[randAttack](req)

	}

	return el.RoundTripper.RoundTrip(req)
}

type DebugLogger struct{}

func (d DebugLogger) Write(p []byte) (n int, err error) {
	s := string(p)
	if strings.Contains(s, "multiple response.WriteHeader") {
		debug.PrintStack()
	}
	return os.Stderr.Write(p)
}

func inLineRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
