package main

import (
  mktp "github.com/helmutkemper/marketPlaceProxy"
  "net/http"
  "io/ioutil"
)

func hello(w mktp.ProxyResponseWriter, r *mktp.ProxyRequest) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  w.Write( []byte( "controller: " ) )
  w.Write( []byte( r.ExpRegMatches[ "controller" ] ) )
  w.Write( []byte( "<br>" ) )

  w.Write( []byte( "module: " ) )
  w.Write( []byte( r.ExpRegMatches[ "module" ] ) )
  w.Write( []byte( "<br>" ) )

  w.Write( []byte( "site: " ) )
  w.Write( []byte( r.ExpRegMatches[ "site" ] ) )
  w.Write( []byte( "<br>" ) )
}

func containerListHtml(w mktp.ProxyResponseWriter, r *mktp.ProxyRequest) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  page, err := ioutil.ReadFile("containerStatus.html")
  if err != nil {
    w.Write( []byte( err.Error() ) )
    return
  }

  w.Write( page )
}

func main() {
  mktp.ProxyRootConfig = mktp.ProxyConfig{
    ListenAndServe: ":8888",
    Routes: []mktp.ProxyRoute{
      {
        // docker run -d --name ghost-blog-demo -p 2368:2368 ghost
        Name: "blog",
        Domain: mktp.ProxyDomain{
          SubDomain: "blog",
          Domain: "localhost",
          Port: "8888",
        },
        ProxyEnable: true,
        ProxyServers: []mktp.ProxyUrl{
          {
            Name: "docker 1 - ok",
            Url: "http://localhost:2368",
          },
          {
            Name: "docker 2 - error",
            Url: "http://localhost:2367",
          },
          {
            Name: "docker 3 - error",
            Url: "http://localhost:2367",
          },
        },
      },
      {
        Name: "hello",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: hello,
        },
      },
      {
        Name: "image_list",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/imageList",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ImageWebList,
        },
      },
      {
        Name: "container_list_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/restContainerList",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerWebList,
        },
      },
      {
        Name: "container_stop_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStop/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerStopById,
        },
      },
      {
        Name: "container_remove_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerRemove/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerRemove,
        },
      },
      {
        Name: "container_kill_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerKill/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerRemove,
        },
      },
      {
        Name: "container_stats_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStats/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerStatsById,
        },
      },
      {
        Name: "container_stats_all_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path: "/restContainerStatsAll",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerWebStatsLog,
        },
      },
      {
        Name: "container_start_rest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStart/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerStart,
        },
      },
      {
        Name: "container_stats_logById",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: "^/restContainerStatsLog/(?P<id>[a-fA-F0-9-]{64})$",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerWebStatsLogById,
        },
      },
      {
        Name: "container_list_html",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/htmlContainerList",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: containerListHtml,
        },
      },
      {
        Name: "container_log",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "panel",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          ExpReg: `^/containerLog/(?P<id>[a-fA-F0-9-]{64})$`,
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ContainerLogsById,
        },
      },







      {
        Name: "addTest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/add",
          Method: "POST",
          //ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ProxyRootConfig.RouteAdd,
        },
      },
      {
        Name: "removeTest",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path : "/remove",
          Method: "POST",
          //ExpReg: `^/(?P<controller>[a-z0-9-]+)/(?P<module>[a-z0-9-]+)/(?P<site>[a-z0-9]+.(htm|html))$`,
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ProxyRootConfig.RouteDelete,
        },
      },
      {
        Name: "panel",
        Domain: mktp.ProxyDomain{
          NotFoundHandle: mktp.ProxyRootConfig.ProxyNotFound,
          ErrorHandle: mktp.ProxyRootConfig.ProxyError,
          SubDomain: "root",
          Domain: "localhost",
          Port: "8888",
        },
        Path: mktp.ProxyPath{
          Path: "/statistics",
          Method: "GET",
        },
        ProxyEnable: false,
        Handle: mktp.ProxyHandle{
          Handle: mktp.ProxyRootConfig.ProxyStatistics,
        },
      },
    },
  }
  mktp.ProxyRootConfig.Prepare()
  go mktp.ProxyRootConfig.VerifyDisabled()

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer( http.Dir( "static" ) ) ) )
  http.HandleFunc("/", mktp.ProxyFunc)
  http.ListenAndServe(mktp.ProxyRootConfig.ListenAndServe, nil)
}