package main

import (
  "flag"
  "fmt"
  sProxy "github.com/helmutkemper/SimpleReverseProxy"
  "github.com/helmutkemper/dockerManager/image"
  "io/ioutil"
  "log"
  "net/http"
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



func main() {
  var err error
  
  filePath := flag.String("f", "./reverseProxy-config.yml", "./reverseProxy-config.yml")
  flag.Parse()
  
  fmt.Printf("reverseProxy version: %v\n", sProxy.KCodeVersion)
  
  configServer := sProxy.NewConfig()
  if err = configServer.Unmarshal( *filePath ); err != nil {
    log.Fatalf("file %v parser error: %v\n", *filePath, err.Error())
  }
  
  fmt.Print("\nServer configs:\n\n")
  
  if configServer.ReverseProxy.Config.OutputConfig == true {
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
  
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyNotFound )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyError )
  sProxy.FuncMap.Add( containerListHtml )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.RouteAdd )
  sProxy.FuncMap.Add( sProxy.ProxyRootConfig.RouteDelete )
  //sProxy.FuncMap.Add( sProxy.ProxyRootConfig.ProxyStatistics )
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
  if err = http.ListenAndServe(sProxy.ProxyRootConfig.ListenAndServe, nil); err != nil {
    log.Fatalf( err.Error() )
  }
  
}