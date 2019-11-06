package capibaribe

import (
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

/*
pt_br: Recebe a lista de containers e endpoints para redirecionar cada endpoint de entrada
*/
type proxy struct {
	// pt_br: Quantidade máxima de testes antes de de uma falha ser aceita
	MaxAttemptToRescueLoop int `yaml:"maxAttemptToRescueLoop" json:"maxAttemptToRescueLoop"`

	// pt_br: ignora a porta de entrada de dados //todo: é isto mesmo?
	IgnorePort bool `yaml:"ignorePort" json:"ignorePort"`

	// pt_br: host do servidor para servidores com vários domínios
	Host string `yaml:"host" json:"host"`

	// escolha do tipo de load balancing
	LoadBalancing string `yaml:"loadBalancing" json:"loadBalancing"`

	// pt_br: path dentro do domínio.
	// quando definido, redireciona o path para um endereço específico
	Path string `yaml:"path" json:"path"`

	Header []headerMonitor `yaml:"header" json:"header"`

	// healthCheck //todo: fazer
	HealthCheck healthCheck `yaml:"healthCheck" json:"healthCheck"`

	// pt_br: lista de servidores secundários
	Servers []servers `yaml:"servers" json:"servers"`

	// quantidade de erros consecutivos
	// zerado quando há um sucesso
	consecutiveErrors int

	// quantidade de sucessos consecutivos
	// zerado quando há um erro
	consecutiveSuccess int

	// total de erros
	errors int

	// total de sucessos
	success int

	keyProxy       int
	keyServer      int
	lastError      error
	lastRoundError bool
}

func (el *proxy) SelectLoadBalance() (string, int) {
	if el.LoadBalancing == KLoadBalanceRandom {
		return el.random()
	} else if el.LoadBalancing == KLoadBalanceExecutionTime {
		return el.executionTime()
	} else if el.LoadBalancing == KLoadBalanceExecutionTimeAverage {
		return el.executionTimeAverage()
	}

	//if el.LoadBalancing == KLoadBalanceRoundRobin || el.LoadBalancing == ""
	return el.roundRobin()

}

func (el *proxy) VerifyHeaderMatchValueToRoute(w http.ResponseWriter, r *http.Request) bool {
	pass := false
	for _, headerData := range el.Header {

		if headerData.Type == KHeaderTypeString && r.Header.Get(headerData.Key) == headerData.Value {
			pass = true
			break
		} else if headerData.Type == KHeaderTypeRegExp {
			re := regexp.MustCompile(headerData.Value)
			if re.MatchString(r.Header.Get(headerData.Key)) == true {
				pass = true
				break
			}
		}

	}

	return pass
}

func (el *proxy) WriteHealthCheckDataToOutputEndpoint(w http.ResponseWriter, r *http.Request) {
	for _, headerData := range el.HealthCheck.Header {
		w.Header().Add(headerData.Key, headerData.Value)
	}

	w.Write([]byte(el.HealthCheck.Body))
}

func (el *proxy) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	el.consecutiveErrors += 1
	el.consecutiveSuccess = 0
	el.errors += 1
	el.lastRoundError = true
	el.lastError = err
}

func (el *proxy) SuccessHandler(w http.ResponseWriter, r *http.Request) {
	el.consecutiveErrors = 0
	el.consecutiveSuccess += 1
	el.success += 1
	el.lastRoundError = false
	el.lastError = nil
}

func (el *proxy) ModifyResponse(resp *http.Response) error {
	//b, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (el *proxy) executionTimeAverage() (string, int) {

	minTime := time.Duration(math.MaxInt64)
	keyToReturn := 0

	for serverKey, serverData := range el.Servers {
		if minTime > serverData.executionDurationAverage {

			keyToReturn = serverKey
			minTime = serverData.executionDurationAverage

		}
	}

	return el.Servers[keyToReturn].Host, keyToReturn
}

func (el *proxy) executionTime() (string, int) {

	minTime := time.Duration(math.MaxInt64)
	keyToReturn := 0

	for serverKey, serverData := range el.Servers {
		if minTime > serverData.executionDurationMin {

			keyToReturn = serverKey
			minTime = serverData.executionDurationMin

		}
	}

	return el.Servers[keyToReturn].Host, keyToReturn
}

func (el *proxy) roundRobin() (string, int) {

	for {

		randNumber := rand.Float64()
		for serverKey, serverData := range el.Servers {
			if randNumber <= serverData.Weight {
				return serverData.Host, serverKey
			}
		}

	}

}

func (el *proxy) random() (string, int) {
	randNumber := rand.Intn(len(el.Servers))
	return el.Servers[randNumber].Host, randNumber
}
