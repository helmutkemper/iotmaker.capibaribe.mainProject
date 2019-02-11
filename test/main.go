package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

const connections = 20000

func main() {

	var a, b uint64
	/*
	  var err error
	  wg.Add( 1 )
	  go func(){
	    filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
	    flag.Parse()

	    fmt.Printf("reverseProxy version: %v\n", sProxy.KCodeVersion)

	    configServer := sProxy.NewConfig()
	    if err = configServer.Unmarshal( *filePath ); err != nil {
	      log.Fatalf("file %v parser error: %v\n", *filePath, err.Error())
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
	      if err != nil {
	        log.Fatal( err.Error() )
	      }
	    }

	    err = sProxy.ProxyRootConfig.AddRouteFromFuncStt(
	      sProxy.ProxyRoute{
	        Name: "process end",
	        Handle: sProxy.ProxyHandle{
	          Handle: func(w sProxy.ProxyResponseWriter, r *sProxy.ProxyRequest){
	            wg.Done()
	          },
	        },
	        Path: sProxy.ProxyPath{
	          Path: "end",
	          Method: "GET",
	        },
	      },
	    )
	    if err != nil {
	      log.Fatal( err.Error() )
	    }

	    err = sProxy.ProxyRootConfig.AddRouteFromFuncStt(
	      sProxy.ProxyRoute{
	        Name: "process ready",
	        Handle: sProxy.ProxyHandle{
	          Handle: func(w sProxy.ProxyResponseWriter, r *sProxy.ProxyRequest){
	            w.Header().Set( "Content-Type", "application/json; charset=utf-8" )
	            w.Write([]byte(`{"started": true}`))
	          },
	        },
	        Path: sProxy.ProxyPath{
	          Path: "end",
	          Method: "GET",
	        },
	      },
	    )
	    if err != nil {
	      log.Fatal( err.Error() )
	    }

	    sProxy.ProxyRootConfig.Prepare()

	    if configServer.ReverseProxy.Config.StaticServer {
	      for _, staticPath := range configServer.ReverseProxy.Config.StaticFolder {
	        http.Handle("/" + staticPath.ServerPath + "/", http.StripPrefix("/" + staticPath.ServerPath + "/", http.FileServer( http.Dir( staticPath.Folder ) ) ) )
	      }
	    }

	    http.HandleFunc("/", sProxy.ProxyFunc)

	    //if err = http.ListenAndServe(configServer.ReverseProxy.Config.ListenAndServer, nil); err != nil {
	      //log.Fatalf( err.Error() )
	    //}
	  }()
	*/
	/*
	  for{
	    res, err := http.Get("http://localhost:8080/end")
	    if err != nil {
	      time.Sleep( 1 * time.Second )
	      continue
	    }
	    responseBody, err := ioutil.ReadAll(res.Body)

	    res.Body.Close()

	    if err != nil {
	      log.Fatal(err)
	    }

	    pass := true
	    for k, v := range []byte(`{"started": true}`){
	      if responseBody[k] != v {
	        pass = false
	        break
	      }
	    }

	    if pass == true {
	      break
	    }
	  }
	*/

	for i := 0; i != connections; i += 1 {
		continue
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()

			res, err := http.Get("http://127.0.0.1:3000")
			if err != nil {
				log.Fatal(err)
			}
			_, err = ioutil.ReadAll(res.Body)

			res.Body.Close()

			if err != nil {
				log.Fatal(err)
			}

			elapsed := time.Since(start)
			a += (uint64)(elapsed.Nanoseconds())
			//log.Printf("%v\n", elapsed.Nanoseconds()/1000)
		}()

	}

	wg.Wait()
	fmt.Printf("média static: %v\n\n\n", a/connections)

	for i := 0; i != connections; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			start := time.Now()

			req, err := http.NewRequest(http.MethodGet, "http://blog.localhost:8080", nil)
			if err != nil {
				log.Fatal(err)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			res.Body.Close()

			elapsed := time.Since(start)
			b += (uint64)(elapsed.Nanoseconds())
			//log.Printf("%v\n", elapsed.Nanoseconds()/1000)

		}()

	}

	wg.Wait()
	fmt.Printf("média proxy: %v\n\n\n", b/connections)
	log.Fatal("process end")

	/*
	  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    fmt.Fprintln(w, "Hello, client")
	  }))
	  defer ts.Close()

	  res, err := http.Get(ts.URL)
	  if err != nil {
	    log.Fatal(err)
	  }
	  greeting, err := ioutil.ReadAll(res.Body)
	  res.Body.Close()
	  if err != nil {
	    log.Fatal(err)
	  }

	  fmt.Printf("%s", greeting)
	*/

	wg.Wait()
}
