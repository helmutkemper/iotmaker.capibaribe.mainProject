package main

import (
  "flag"
  "fmt"
  sProxy "github.com/helmutkemper/SimpleReverseProxy"
  "github.com/helmutkemper/dockerManager/image"
  "github.com/helmutkemper/yaml"
  "github.com/pkg/errors"
  "io/ioutil"
  "log"
  "net/http"
  "reflect"
  "strconv"
)



func containerListHtml(w sProxy.ProxyResponseWriter, r *sProxy.ProxyRequest) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  page, err := ioutil.ReadFile("containerStatus.html")
  if err != nil {
    w.Write( []byte( err.Error() ) )
    return
  }

  w.Write( page )
}

const vCodeVersion = "0.1 alpha"
const kVersionMinimum = 1.0
const kVersionMaximum   = 1.0
const kVersionMinimumString = "1.0"
const kVersionMaximumString = "1.0"
const kSiteErrorInformation = " Please, se manual at kemper.com.br for more information."


type ConfigProxyServerNameAndHost struct {
  Name string `yaml:"name"`
  Host string `yaml:"host"`
}

type ConfigProxy struct {
  Host string `yaml:"host"`
  Server []ConfigProxyServerNameAndHost `yaml:"server"`
}

type ConfigServer struct {
  ListenAndServer string `yaml:"listenAndServer"`
  OutputConfig bool `yaml:"outputConfig"`
  StaticServer bool `yaml:"staticServer"`
  StaticFolder []string `yaml:"staticFolder"`
}

type ConfigReverseProxy struct {
  Config ConfigServer `yaml:"config"`
  Proxy map[string]ConfigProxy `yaml:"proxy"`
}

type ConfigMainServer struct {
  Version string `yaml:"version"`
  ReverseProxy ConfigReverseProxy `yaml:"reverseProxy"`
}

func (el *ConfigMainServer) Unmarshal(filePath string) error {
  var fileContent []byte
  var err error
  var version float64
  
  fileContent, err = ioutil.ReadFile(filePath)
  if err != nil {
    return err
  }
  
  err = yaml.Unmarshal(fileContent, el)
  if err != nil {
    return err
  }
  
  version, err = strconv.ParseFloat( el.Version, 64 )
  
  if version == 0.0 {
    return errors.New("you must inform the version of config file as numeric value. Example: version: '1.0'" + kSiteErrorInformation)
  }
  
  if version < kVersionMinimum || version > kVersionMaximum {
    return errors.New("this project version accept only configs between versions " + kVersionMaximumString + " and " + kVersionMinimumString + "." + kSiteErrorInformation)
  }
  
  if reflect.DeepEqual( el.ReverseProxy, ConfigReverseProxy{} ) {
    return errors.New("reverse proxy config not found. " + kSiteErrorInformation)
  }
  
  if el.ReverseProxy.Config.ListenAndServer == "" {
    return errors.New("reverseProxy > config > listenAndServer config not found. " + kSiteErrorInformation)
  }
  
  if el.ReverseProxy.Config.StaticServer == true && len( el.ReverseProxy.Config.StaticFolder ) == 0 {
    return errors.New("reverseProxy > config > staticFolder config not found. " + kSiteErrorInformation)
  }
  
  for _, folder := range el.ReverseProxy.Config.StaticFolder {
    _, err = ioutil.ReadDir( folder )
    if err != nil {
      return errors.New("reverseProxy > config > staticFolder error: " + err.Error())
    }
  }
  
  for proxyName, proxyConfig := range el.ReverseProxy.Proxy {
    if proxyConfig.Host == "" {
      return errors.New("reverseProxy > proxy > " + proxyName + " > host not found. " + kSiteErrorInformation)
    }
    
    for _, proxyServerConfig := range proxyConfig.Server {
      if proxyServerConfig.Host == "" {
        return errors.New("reverseProxy > proxy > " + proxyName + " > server > host not found. " + kSiteErrorInformation)
      }
      
      if proxyServerConfig.Name == "" {
        return errors.New("reverseProxy > proxy > " + proxyName + " > server > name not found. " + kSiteErrorInformation)
      }
    }
  }
  
  return nil
}


// fixme: quando o usuário for criar um banco de dados ou parecido, ele tem que ser avisado de exportar o diretório de dados para a máquina
func main() {
  var err error
  
  filePath := flag.String("f", "./reverseProxy-config.yml", "./reverseProxy-config.yml")
  flag.Parse()
  
  fmt.Printf("reverseProxy version: %v\n", vCodeVersion)
  
  configServer := ConfigMainServer{}
  if err = configServer.Unmarshal( *filePath ); err != nil {
    log.Fatalf("file %v parser error: %v\n", *filePath, err.Error())
  }
  
  fmt.Print("\nServer configs:\n\n")
  
  if configServer.ReverseProxy.Config.OutputConfig == true {
    fmt.Printf("listen and server: %v\n", configServer.ReverseProxy.Config.ListenAndServer)
    
    if configServer.ReverseProxy.Config.StaticServer == true {
      fmt.Print("static server enabled at folders:\n" )
      for _, folder := range configServer.ReverseProxy.Config.StaticFolder {
        fmt.Printf("  %v\n", folder)
      }
    }
    
    fmt.Print("\nProxy configs:\n\n")
    
    for proxyServiceName, proxyServiceConfig := range configServer.ReverseProxy.Proxy {
      
      fmt.Printf( "proxy service name: %v\n", proxyServiceName )
      fmt.Printf( "proxy income host: %v\n", proxyServiceConfig.Host )
      fmt.Printf( "proxy alternative hosts: %v\n", len( proxyServiceConfig.Server ) )
      
      for _, toServerConfig := range proxyServiceConfig.Server {
        fmt.Printf( "  alternative host name: %v\n", toServerConfig.Name )
        fmt.Printf( "  alternative host addr: %v\n", toServerConfig.Host )
      }
      
      fmt.Println()
    }
  }
  
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyNotFound )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyError )
  sProxy.FuncMap.Add( containerListHtml )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.RouteAdd )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.RouteDelete )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyStatistics )
  sProxy.FuncMap.Add( image.WebList )
  
  
  
  
  sProxy.ProxyRootConfig = sProxy.ProxyConfig{
    ListenAndServe: configServer.ReverseProxy.Config.ListenAndServer,
  }
  
  for proxyConfigName, proxyConfig := range configServer.ReverseProxy.Proxy {
    
    servers := make( []sProxy.ProxyUrl, len( proxyConfig.Server ) )
    for k, v := range proxyConfig.Server {
      servers[ k ] = sProxy.ProxyUrl{
        Name: v.Name,
        Url: v.Host,
      }
    }
    
    err = sProxy.ProxyRootConfig.AddRouteToProxyStt(
      sProxy.ProxyRoute{
        // docker run -d --name ghost-blog-demo -p 2368:2368 ghost
        Name: proxyConfigName,
        Domain: sProxy.ProxyDomain{
          Host: proxyConfig.Host,
        },
        ProxyEnable: true,
        ProxyServers: servers,
      },
    )
  }

  sProxy.ProxyRootConfig.Prepare()

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer( http.Dir( "static" ) ) ) )
  http.HandleFunc("/", sProxy.ProxyFunc)
  http.ListenAndServe(sProxy.ProxyRootConfig.ListenAndServe, nil)
}