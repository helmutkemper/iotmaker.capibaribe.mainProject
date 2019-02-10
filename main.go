package main

import (
  "flag"
  "fmt"
  sProxy "github.com/helmutkemper/SimpleReverseProxy"
  "log"
  "net/http"
)

func main() {
  var err error
  
  filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
  flag.Parse()
  
  fmt.Printf("reverseProxy version: %v\n", sProxy.KCodeVersion)
  
  configServer := sProxy.NewConfig()
  if err = configServer.Unmarshal( *filePath ); err != nil {
    log.Fatalf("file %v parser error: %v\n", *filePath, err.Error())
  }
  
  if configServer.ReverseProxy.Config.OutputConfig == true {
    fmt.Print("\nServer configs:\n\n")
    
    fmt.Printf("listen and server: %v\n", configServer.ReverseProxy.Config.ListenAndServer)
    
    if configServer.ReverseProxy.Config.StaticServer == true {
      fmt.Print("static server enabled at folders:\n" )
      for _, folder := range configServer.ReverseProxy.Config.StaticFolder {
        fmt.Printf("  [%v]: %v \n", folder.ServerPath, folder.Folder)
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
  
  fmt.Print("stating server...\n\n")
  
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
  
  if configServer.ReverseProxy.Config.StaticServer {
    for _, staticPath := range configServer.ReverseProxy.Config.StaticFolder {
      http.Handle("/" + staticPath.ServerPath + "/", http.StripPrefix("/" + staticPath.ServerPath + "/", http.FileServer( http.Dir( staticPath.Folder ) ) ) )
    }
  }
  
  http.HandleFunc("/", sProxy.ProxyFunc)
  
  if err = http.ListenAndServe(configServer.ReverseProxy.Config.ListenAndServer, nil); err != nil {
    log.Fatalf( err.Error() )
  }
  
  /*
  mux := http.NewServeMux()
  if configServer.ReverseProxy.Config.StaticServer {
    for _, staticPath := range configServer.ReverseProxy.Config.StaticFolder {
      mux.Handle("/" + staticPath.ServerPath + "/", http.StripPrefix("/" + staticPath.ServerPath + "/", http.FileServer( http.Dir( staticPath.Folder ) ) ) )
    }
  }
  mux.HandleFunc("/", sProxy.ProxyFunc)
  
  if err = certmagic.HTTPS([]string{configServer.ReverseProxy.Config.ListenAndServer}, mux); err != nil {
    log.Fatalf( err.Error() )
  }
  */
  
}