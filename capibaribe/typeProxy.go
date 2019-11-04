package capibaribe

import (
	"math/rand"
	"net/http"
)

/*
pt_br: Recebe a lista de containers e endpoints para redirecionar cada endpoint de entrada
*/
type proxy struct {
	// pt_br: Quantidade máxima de testes antes de de uma falha ser aceita
	MaxAttemptToRescueLoop int `yaml:"maxAttemptToRescueLoop" json:"maxAttemptToRescueLoop"`

	// pt_br: ignora a porta de entrada de dados //todo: é isto mesmo?
	IgnorePort bool `yaml:"ignorePort"             json:"ignorePort"`

	// pt_br: host local, ex.: 127.0.0.1 //todo: é isto mesmo?
	Host string `yaml:"host"                   json:"host"`

	// escolha do tipo de load balancing
	LoadBalancing string `yaml:"loadBalancing"          json:"loadBalancing"`

	// todo: o que é isto?
	Path string `yaml:"path"                   json:"path"`

	// healthCheck //todo: fazer
	HealthCheck healthCheck `yaml:"healthCheck"            json:"healthCheck"`

	// pt_br: lista de servidores secundários
	Servers []servers `yaml:"servers"                json:"servers"`

	consecutiveErrors  int
	consecutiveSuccess int
	errors             int
	success            int
	keyProxy           int
	keyServer          int
	lastError          error
	lastRoundError     bool
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
	//b, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}

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
