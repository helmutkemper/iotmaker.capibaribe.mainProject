//todo: tem que ter um log avisando quando foi a ultima requisição

package main

import (
	capib "./capibaribe"
	//"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	//"fmt"
	//"github.com/etcd-io/etcd/clientv3"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

/*func etcdTest() {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	cli, _ := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"127.0.0.1:2379"},
	})
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	wt := clientv3.NewWatcher(cli)

	GetSingleValueDemo(ctx, kv, wt)
}*/

/*func GetSingleValueDemo(ctx context.Context, kv clientv3.KV, wt clientv3.Watcher) {
	go func() {

		rch := wt.Watch(context.Background(), "key")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf(">> %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}

	}()

	fmt.Println("*** GetSingleValueDemo()")
	// Delete all keys
	kv.Delete(ctx, "key", clientv3.WithPrefix())

	// Insert a key value
	pr, _ := kv.Put(ctx, "key", "444")
	rev := pr.Header.Revision
	fmt.Println("Revision:", rev)

	gr, _ := kv.Get(ctx, "key")
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	// Modify the value of an existing key (create new revision)
	kv.Put(ctx, "key", "555")

	gr, _ = kv.Get(ctx, "key")
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	// Get the value of the previous revision
	gr, _ = kv.Get(ctx, "key", clientv3.WithRev(rev))
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	kv.Delete(ctx, "key", clientv3.WithPrefix())
}*/

/*var mutexTerminateServer map[string]*sync.WaitGroup

func terminateAllServers() {

	for affluentRiverName := range mutexTerminateServer {
		mutexTerminateServer[ affluentRiverName ].Done()
	}

}

func terminateServerByName (affluentRiverName string) {

	if len( mutexTerminateServer ) == 0 {
		return
	}

	mutexTerminateServer[ affluentRiverName ].Done()
}*/

var config capib.MainConfig

func main() {
	var wg sync.WaitGroup

	wg.Add(1)

	//mutexTerminateServer = make( map[string]*sync.WaitGroup )

	filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
	howToConfig := flag.String("howToConf", "file", "file|etcd")
	etcdConn := flag.String("etcdConn", "['127.0.0.1:2379']", "['127.0.0.1:2379']")
	etcdKey := flag.String("etcdKey", "capibaribe-config-yaml-file", "capibaribe-config-yaml-file")
	flag.Parse()

	loadConf(*filePath, *howToConfig, *etcdConn, *etcdKey)

	//etcdTest()

	wg.Wait()
}

func loadConf(filePath, howToConfig, etcdConn, etcdKey string) {
	var err error
	//var etcdKeyConfig *clientv3.GetResponse

	//config = capib.MainConfig{}

	//for _, AffluentRiverConfig := range config.AffluentRiver {
	//AffluentRiverConfig.Wait()
	//}

	switch howToConfig {

	case "etcd":

		err = json.Unmarshal([]byte(etcdConn), &config.Etcd.Connection)
		if err != nil {
			log.Fatal(err.Error())
		}

		config.Etcd.ConfigKey = etcdKey

	case "file":

		err = config.Unmarshal(filePath)
		if err != nil {
			log.Fatal(err.Error())
		}

	default:
		log.Fatal("-howToConf must be 'file' or 'etcd'")
	}

	config.Etcd.Prepare()

	/*ctx, _ := context.WithTimeout(context.Background(), config.Etcd.RequestTimeout)
	cli, _ := clientv3.New(clientv3.Config{
		DialTimeout: config.Etcd.DialTimeOut,
		Endpoints:   config.Etcd.Connection,
	})
	defer cli.Close()

	kv := clientv3.NewKV(cli)
	//wt := clientv3.NewWatcher(cli)

	switch howToConfig {
	case "etcd":

		etcdKeyConfig, err = kv.Get(ctx, config.Etcd.ConfigKey)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = config.UnmarshalByte([]byte(etcdKeyConfig.Kvs[0].Value))
		if err != nil {
			log.Fatal(err.Error())
		}

	}*/

	for projectName, projectConfig := range config.AffluentRiver {

		go func(name string, config capib.Project) {
			var err error

			//mutexTerminateServer[ name ].Add(1)

			server := http.NewServeMux()

			//defer mutexTerminateServer[ name ].Done()

			for _, staticPath := range config.Static {

				if _, err = os.Stat(staticPath.FilePath); os.IsNotExist(err) {
					log.Fatalf("static dir error: %v\n", err.Error())
				}

				server.Handle("/"+staticPath.ServerPath+"/", http.StripPrefix("/"+staticPath.ServerPath+"/", http.FileServer(http.Dir(staticPath.FilePath))))
			}

			server.HandleFunc("/", config.HandleFunc)

			if config.Sll.Enabled == true {

				var certificatesList = tls.Certificate{}

				if config.Sll.X509.Certificate != "" && config.Sll.X509.CertificateKey != "" {

					if _, err = os.Stat(config.Sll.X509.Certificate); os.IsNotExist(err) {
						log.Fatalf("sll x509 certificate error: %v\n", err.Error())
					}

					if _, err = os.Stat(config.Sll.X509.CertificateKey); os.IsNotExist(err) {
						log.Fatalf("sll x509 certificate key error: %v\n", err.Error())
					}

					certificatesList, err = tls.LoadX509KeyPair(config.Sll.X509.Certificate, config.Sll.X509.CertificateKey)
					if err != nil {
						log.Fatalf("sll x509 certificate load pair error: %v\n", err.Error())
					}

				}

				var tlsMinVersion uint16 = 0
				if config.Sll.Version.Min == 10 {
					tlsMinVersion = tls.VersionTLS10
				} else if config.Sll.Version.Min == 11 {
					tlsMinVersion = tls.VersionTLS11
				} else if config.Sll.Version.Min == 12 {
					tlsMinVersion = tls.VersionTLS12
				} else if config.Sll.Version.Min == 30 {
					tlsMinVersion = tls.VersionSSL30
				}

				var tlsMaxVersion uint16 = 0
				if config.Sll.Version.Max == 10 {
					tlsMaxVersion = tls.VersionTLS10
				} else if config.Sll.Version.Max == 11 {
					tlsMaxVersion = tls.VersionTLS11
				} else if config.Sll.Version.Max == 12 {
					tlsMaxVersion = tls.VersionTLS12
				} else if config.Sll.Version.Max == 30 {
					tlsMaxVersion = tls.VersionSSL30
				}

				var curveIdList = make([]tls.CurveID, len(config.Sll.CurvePreferences.([]string)))
				for k, v := range config.Sll.CurvePreferences.([]string) {
					if v == "P256" {
						curveIdList[k] = tls.CurveP256
					} else if v == "P384" {
						curveIdList[k] = tls.CurveP384
					} else if v == "P521" {
						curveIdList[k] = tls.CurveP521
					} else if v == "X25519" {
						curveIdList[k] = tls.X25519
					}
				}

				var cipherSuitesList = make([]uint16, len(config.Sll.CipherSuites.([]string)))
				for k, v := range config.Sll.CurvePreferences.([]string) {
					if v == "TLS_RSA_WITH_RC4_128_SHA" {
						cipherSuitesList[k] = tls.TLS_RSA_WITH_RC4_128_SHA
					} else if v == "TLS_RSA_WITH_3DES_EDE_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA
					} else if v == "TLS_RSA_WITH_AES_128_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_128_CBC_SHA
					} else if v == "TLS_RSA_WITH_AES_256_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_256_CBC_SHA
					} else if v == "TLS_RSA_WITH_AES_128_CBC_SHA256" {
						cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_128_CBC_SHA256
					} else if v == "TLS_RSA_WITH_AES_128_GCM_SHA256" {
						cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_128_GCM_SHA256
					} else if v == "TLS_RSA_WITH_AES_256_GCM_SHA384" {
						cipherSuitesList[k] = tls.TLS_RSA_WITH_AES_256_GCM_SHA384
					} else if v == "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA" {
						cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA
					} else if v == "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
					} else if v == "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
					} else if v == "TLS_ECDHE_RSA_WITH_RC4_128_SHA" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA
					} else if v == "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
					} else if v == "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
					} else if v == "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
					} else if v == "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256" {
						cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
					} else if v == "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
					} else if v == "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
					} else if v == "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256" {
						cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
					} else if v == "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
					} else if v == "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384" {
						cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
					} else if v == "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305" {
						cipherSuitesList[k] = tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305
					} else if v == "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305" {
						cipherSuitesList[k] = tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305
					} else if v == "TLS_FALLBACK_SCSV" {
						cipherSuitesList[k] = tls.TLS_FALLBACK_SCSV
					}
				}

				newServer := &http.Server{
					TLSConfig: &tls.Config{
						MinVersion:               tlsMinVersion,
						MaxVersion:               tlsMaxVersion,
						CurvePreferences:         curveIdList,
						PreferServerCipherSuites: config.Sll.PreferServerCipherSuites,
						CipherSuites:             cipherSuitesList,
						Certificates:             []tls.Certificate{certificatesList},
					},
					//TLSNextProto:               make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
					Addr:    config.Listen,
					Handler: server,
				}

				if config.DebugServerEnable == true {
					newServer.ErrorLog = log.New(capib.DebugLogger{}, "", 0)
				}

				log.Fatal(newServer.ListenAndServeTLS(config.Sll.Certificate, config.Sll.CertificateKey))

			} else {

				newServer := &http.Server{
					Addr:    config.Listen,
					Handler: server,
				}

				if config.DebugServerEnable == true {
					newServer.ErrorLog = log.New(capib.DebugLogger{}, "", 0)
				}

				log.Fatal(newServer.ListenAndServe())

			}

		}(projectName, projectConfig)

	}
}
