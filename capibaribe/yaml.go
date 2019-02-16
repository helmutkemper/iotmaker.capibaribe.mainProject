package capibaribe

import (
	"bytes"
	"github.com/helmutkemper/seelog"
	"github.com/helmutkemper/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
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
	Listen       string       `yaml:"listen"`
	Sll          ssl          `yaml:"ssl"`
	CrazyMonkeys crazyMonkeys `yaml:"crazyMonkeys"`
	Proxy        []proxy      `yaml:"proxy"`
	Static       []static     `yaml:"static"`
}

func (el *Project) HandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	var host = r.Host
	var remoteAddr string
	var re *regexp.Regexp
	var hostServer string

	if el.Proxy != nil {

		for _, proxyData := range el.Proxy {

			pass := len(proxyData.Bind) == 0
			for _, bind := range proxyData.Bind {
				if bind.IgnorePort == true {
					if re, err = regexp.Compile(kIgnorePortRegExp); err != nil {
						if err = seelog.Errorf("reg exp error: %v\n", err.Error()); err != nil {
							log.Fatalf("seelog is miss configured. Error: %v\n", err.Error())
						}
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
					if err = seelog.Errorf("reg exp error: %v\n", err.Error()); err != nil {
						log.Fatalf("seelog is miss configured. Error: %v\n", err.Error())
					}
				}

				host = re.ReplaceAllString(host, "$1")
			}

			if proxyData.Host == host {

				if proxyData.LoadBalancing == kLoadBalanceRoundRobin || proxyData.LoadBalancing == "" {

					hostServer = proxyData.roundRobin()

				} else if proxyData.LoadBalancing == kLoadBalanceRandom {

					hostServer = proxyData.random()

				}

				if hostServer != "" {

					rpURL, err := url.Parse(hostServer)
					if err != nil {
						log.Fatal(err)
						// fixme: seelog aqui
					}

					proxy := httputil.NewSingleHostReverseProxy(rpURL)
					proxy.Transport = &transport{http.DefaultTransport}
					proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
						// fixme: o que fazer quando tem erro?
					}

					proxy.ServeHTTP(w, r)
					return

				} else {

					//fixme: colocar um log aqui

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

type crazyMonkeys struct {
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
	IgnorePort    bool        `yaml:"ignorePort"`
	Host          string      `yaml:"host"`
	Bind          []bind      `yaml:"bind"`
	LoadBalancing string      `yaml:"loadBalancing"`
	Path          string      `yaml:"path"`
	HealthCheck   healthCheck `yaml:"healthCheck"`
	Servers       []servers   `yaml:"servers"`
}

func (el *proxy) roundRobin() string {
	randNumber := rand.Float64()

	for _, serverData := range el.Servers {

		if randNumber <= serverData.Weight {
			return serverData.Host
		}

	}

	return ""
}

func (el *proxy) random() string {
	randNumber := rand.Intn(len(el.Servers))
	return el.Servers[randNumber].Host
}

type bind struct {
	Host       string `yaml:"host"`
	IgnorePort bool   `yaml:"ignorePort"`
}

type servers struct {
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
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("reverse"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}
