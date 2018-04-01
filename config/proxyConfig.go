package config

import (
  "encoding/json"
  "net/url"
  "os"
  "fmt"
)

type ProxyConfig struct {
  /*
  Configuração do seelog
  @see https://github.com/cihub/seelog
  */
  SeeLogConfig                    string                  `json:"seeLogConfig"`

  /*
  Configuração das chaves {{incompleto!}}
  */
  Rsa                             ProxyRsaStt             `json:"rsa"`

  /*
  Expressão regular que identifica o domínio do site
  */
  DomainExpReg                    string                  `json:"domainExpReg"`

  /*
  Função de erro genérica, caso a função do domínio não seja definida
  */
  ErrorHandle                     ProxyHandlerFunc        `json:"-"`

  /*
  Função de page not found genérica, caso a função do domínio não seja definida
  */
  NotFoundHandle                  ProxyHandlerFunc        `json:"-"`

  /*
  Tamanho de caracteres do token de segurança
  */
  UniqueIdLength                  int                     `json:"uniqueIdLength"`

  /*
  URL do servidor principal
  */
  ListenAndServe                  string                  `json:"listenAndServe"`

  /*
  Quantidade máxima de loop quando todas as rotas do proxy falham
  */
  MaxLoopTry                      int

  /*
  Quantidades de erros consecutivos para desabilitar uma rota do proxy.
  A ideia é que uma rota do proxy possa está dando erro temporário, assim, o código desabilita a rota por um tempo e
  depois habilita de novo para testar se a mesma voltou.
  Caso haja apenas uma instabilidade, a rota continua.
  */
  ConsecutiveErrorsToDisable      int64

  /*
  Tempo para manter uma rota do proxy desabilitada antes de testar novamente
  */
  TimeToKeepDisabled              time.Duration

  /*
  Há uma função em loop infinito e a cada x período de tempo, ela verifica se alguma rota está desabilitada e reabilita
  caso o tempo de espera tenha sido excedido
  */
  TimeToVerifyDisabled            time.Duration

  /*
  Rotas do servidor proxy
  */
  Routes                          []ProxyRoute
}

/*
Esta função adiciona novas rotas ao proxy
{
    "name": "news",
    "domain": {
      "subDomain": "news",
      "domain": "localhost",
      "port": "8888"
    },
    "proxyEnable": true,
    "proxyServers": [
    {
          "name": "docker 1 - ok",
          "url": "http://localhost:2368"
    },
    {
          "name": "docker 2 - ok",
          "url": "http://localhost:2368"
    },
    {
      "name": "docker 3 - ok",
          "url": "http://localhost:2368"
    }
  ]
}
*/
func(el *ProxyConfig)RouteAdd(w ProxyResponseWriter, r *ProxyRequest) {
  // Esta função coloca a rota nova em 'ProxyNewRootConfig' e espera uma nova chamada em uma rota qualquer para que a
  // nova rota tenha efeito. Isso é transparente para o usuário final, mas, a rota não pode entrar em vigor durante o
  // processamento da rota anterior, ou o sistema trava, devido a mudança dos 'ponteiros'

  var newRoute ProxyRoute
  var output = JSonOutStt{}

  if len(ProxyNewRootConfig) != 0 {
    ProxyRootConfig.Routes = ProxyNewRootConfig
  }

  err := json.NewDecoder(r.Body).Decode(&newRoute)

  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  if newRoute.ProxyEnable == false {
    output.ToOutput(0, errors.New("this function only adds new routes that can be used in conjunction with the reverse proxy"), []int{}, w)
    return
  }

  if len( newRoute.ProxyServers ) == 0 {
    output.ToOutput(0, errors.New("this function must receive at least one route that can be used in conjunction with the reverse proxy"), []int{}, w)
    return
  }

  for _, route := range newRoute.ProxyServers {
    if route.Name == "" {
      output.ToOutput(0, errors.New("every route must have a name assigned to it"), []int{}, w)
      return
    }

    _, err := url.Parse(route.Url)
    if err != nil {
      output.ToOutput(0, errors.New("the route of name '" + route.Name + "' presented the following error: " + err.Error()), []int{}, w)
      return
    }
  }

  // Habilita todas as rotas, pois, o padrão do go é false
  for urlKey := range newRoute.ProxyServers {
    newRoute.ProxyServers[ urlKey ].Enabled = true
  }

  ProxyNewRootConfig = append(ProxyRootConfig.Routes, newRoute)
  output.ToOutput(int64( len( ProxyNewRootConfig ) ), nil, ProxyNewRootConfig, w)
}

/*
Esta função elimina rotas do proxy
{
    "name": "name_of_route"
}
*/
func(el *ProxyConfig)RouteDelete(w ProxyResponseWriter, r *ProxyRequest) {
  // Esta função coloca a rota nova em 'ProxyNewRootConfig' e espera uma nova chamada em uma rota qualquer para que a
  // nova rota tenha efeito. Isso é transparente para o usuário final, mas, a rota não pode entrar em vigor durante o
  // processamento da rota anterior, ou o sistema trava, devido a mudança dos 'ponteiros'

  var newRoute ProxyRoute
  var output = JSonOutStt{}

  err := json.NewDecoder(r.Body).Decode(&newRoute)

  if err != nil {
    output.ToOutput(0, err, []int{}, w)
    return
  }

  var i int
  nameFound := false
  for i = 0; i != len( ProxyRootConfig.Routes ); i += 1 {
    if ProxyRootConfig.Routes[i].Name == newRoute.Name {
      nameFound = true
      break
    }
  }

  if nameFound == true {
    if i == 0 {
      ProxyNewRootConfig = ProxyRootConfig.Routes[1:]
    } else if i == len(ProxyRootConfig.Routes) - 1 {
      ProxyNewRootConfig = ProxyRootConfig.Routes[:len(ProxyRootConfig.Routes)-1]
    } else {
      ProxyNewRootConfig = append( ProxyRootConfig.Routes[0:i], ProxyRootConfig.Routes[i+1:]... )
    }
  }

  if ProxyRootConfig.Routes[i].ProxyEnable == false {
    output.ToOutput(0, errors.New("this function can only remove the routes used with the reverse proxy, not being able to remove other types of routes"), []int{}, w)
    return
  }

  b, e := json.Marshal(ProxyNewRootConfig)
  if e != nil {
    w.Write([]byte(e.Error()))
    return
  }
  w.Write(b)

  output.ToOutput(int64( len( ProxyNewRootConfig ) ), nil, ProxyNewRootConfig, w)
}

// Inicializa algumas variáveis
func(el *ProxyConfig)Prepare(){

  // Cria a pasta de log
  logPath := "./log"
  if _, err := os.Stat( logPath ); os.IsNotExist( err ) {
    err = os.Mkdir( logPath, 0777 )
  }

  // Configura o log como arquivos com tamanho limitado. Um arquivo, info.log para coisas simples e um arquivo warn.log
  // para coisas que devem ser observadas pelo administrador
  if el.SeeLogConfig == "" {
    el.SeeLogConfig = `<seelog minlevel="warn" maxlevel="critical" type="sync">
  <outputs formatid="all">
    <filter levels="trace">
      <rollingfile type="size" filename="` + logPath + `/info.log" maxrolls="2" maxsize="10000" />
    </filter>
    <filter levels="debug">
      <rollingfile type="size" filename="` + logPath + `/info.log" maxrolls="2" maxsize="10000" />
    </filter>
    <filter levels="info">
      <rollingfile type="size" filename="` + logPath + `/info.log" maxrolls="2" maxsize="10000" />
    </filter>
    <filter levels="warn">
      <rollingfile type="size" filename="` + logPath + `/warn.log" maxrolls="2" maxsize="10000" />
      <console/>
    </filter>
    <filter levels="error">
      <rollingfile type="size" filename="` + logPath + `/warn.log" maxrolls="2" maxsize="10000" />
      <console/>
    </filter>
    <filter levels="critical">
      <rollingfile type="size" filename="` + logPath + `/warn.log" maxrolls="2" maxsize="10000" />
      <console/>
    </filter>
  </outputs>
  <formats>
    <format id="all" format="[%Level::%Date %Time] %Msg%n"/>
  </formats>
</seelog>`
  }

  logger, err := log.LoggerFromConfigAsBytes([]byte(el.SeeLogConfig))
  if err != nil {
    fmt.Printf( "Erro na configuração do log. Error: %v", err.Error() )
  }
  log.UseLogger(logger)

  // Define o tamanho do token como sendo 30 caracteres
  if el.UniqueIdLength == 0 {
    el.UniqueIdLength = 30
  }

  // Expressão regular do domínio
  if el.DomainExpReg == "" {
    el.DomainExpReg = `^(?P<subDomain>[a-zA-Z0-9]??|[a-zA-Z0-9]?[a-zA-Z0-9-.]*?[a-zA-Z0-9]*)[.]*(?P<domain>[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9]):*(?P<port>[0-9]*)$`
  }

  // Após 20 tentativas de acessar todos os containers da rota, uma mensagem de erro é exibida
  if el.MaxLoopTry == 0 {
    el.MaxLoopTry = 20
  }

  // Caso um container apresente mais de 10 erros consecutivos, o mesmo é desabilitado
  if el.ConsecutiveErrorsToDisable == 0 {
    el.ConsecutiveErrorsToDisable = 10
  }

  // Deixa um container desabilitado por 90 segundos após vários erros consecutivos
  if el.TimeToKeepDisabled == 0 {
    el.TimeToKeepDisabled = time.Second * 90
  }

  // Faz um teste a cada 30 segundos para saber se há containers desabilitados além do tempo
  if el.TimeToVerifyDisabled == 0 {
    el.TimeToVerifyDisabled = time.Second * 30
  }

  // Função de erro padrão do sistema
  if el.ErrorHandle == nil {
    el.ErrorHandle = el.ProxyError
  }

  // Função de page not found padrão do sistema
  if el.NotFoundHandle == nil {
    el.NotFoundHandle = el.ProxyNotFound
  }

  // Habilita todas as rotas do proxy, pois, o padrão do go é false
  for routesKey := range el.Routes{
    for urlKey := range el.Routes[ routesKey ].ProxyServers {
      el.Routes[ routesKey ].ProxyServers[ urlKey ].Enabled = true
    }
  }
}

func(el *ProxyConfig)ProxyError(w ProxyResponseWriter, r *ProxyRequest) {
  w.Write( []byte( `<html><header><style>body{height:100%; position:relative}div{margin:auto;height: 100%;width: 100%;position:fixed;top:0;bottom:0;left:0;right:0;background:blue;}div.center{margin:auto;height: 70%;width: 70%;}</style></header><body><div><div style="color:#ffff;" class="center"><p style="text-align: center; background-color: #888888;">There is something very wrong!</p><p>&nbsp;</p>The address is correct, but no server has responded correctly. The system administrator will be informed about this.<p>&nbsp;</p>Mussum Ipsum, cacilds vidis litro abertis. Interagi no mé, cursus quis, vehicula ac nisi. Viva Forevis aptent taciti sociosqu ad litora torquent. Atirei o pau no gatis, per gatis num morreus. Quem num gosta di mim que vai caçá sua turmis!</div></div></body></html>` ) )
}

func(el *ProxyConfig)ProxyNotFound(w ProxyResponseWriter, r *ProxyRequest) {
  w.Write( []byte( `<html><header><style>body{height:100%; position:relative}div{margin:auto;height: 100%;width: 100%;position:fixed;top:0;bottom:0;left:0;right:0;background:blue;}div.center{margin:auto;height: 70%;width: 70%;}</style></header><body><div><div style="color:#ffff;" class="center"><p style="text-align: center; background-color: #888888;">Page Not Found!</p><p>&nbsp;</p>Mussum Ipsum, cacilds vidis litro abertis. Interagi no mé, cursus quis, vehicula ac nisi. Viva Forevis aptent taciti sociosqu ad litora torquent. Atirei o pau no gatis, per gatis num morreus. Quem num gosta di mim que vai caçá sua turmis!<p>&nbsp;</p>Mussum Ipsum, cacilds vidis litro abertis. Interagi no mé, cursus quis, vehicula ac nisi. Viva Forevis aptent taciti sociosqu ad litora torquent. Atirei o pau no gatis, per gatis num morreus. Quem num gosta di mim que vai caçá sua turmis!</div></div></body></html>` ) )
}

// Verifica se há urls do proxy desabilitadas e as habilita depois de um tempo
// A ideia é que o servidor possa está fora do ar por um tempo, por isto, ele remove a rota por algum tempo, para evitar
// chamadas desnecessárias ao servidor
func(el *ProxyConfig)VerifyDisabled(){
  for {
    for routesKey := range el.Routes {
      for urlKey := range el.Routes[ routesKey ].ProxyServers {
        if time.Since(el.Routes[ routesKey ].ProxyServers[ urlKey ].DisabledSince) >= el.TimeToKeepDisabled && el.Routes[ routesKey ].ProxyServers[ urlKey ].Enabled == false && el.Routes[ routesKey ].ProxyServers[ urlKey ].Forever == false {
          el.Routes[ routesKey ].ProxyServers[ urlKey ].ErrorConsecutiveCounter = 0
          el.Routes[ routesKey ].ProxyServers[ urlKey ].Enabled = true
        }
      }
    }

    time.Sleep(el.TimeToVerifyDisabled)
  }
}

func(el *ProxyConfig)ProxyStatistics(w ProxyResponseWriter, r *ProxyRequest) {

  byteJSon, err := json.Marshal( ProxyRootConfig.Routes )
  if err != nil {
    w.Write( []byte( err.Error() ) )
    return
  }

  w.Write( byteJSon )
}

