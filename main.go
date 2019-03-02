//todo: tem que ter um log avisando quando foi a ultima requisição
//todo: microsservico deveria ter um id inicial que permita seguir uma requisição em caso de erro
package main

import (
	capib "./capibaribe"

	"encoding/json"
	"flag"
	//"fmt"
	//"github.com/etcd-io/etcd/clientv3"
	"log"
	"net/http"
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

		if projectConfig.ListenAndServer.InAddress != "" {

			go func(name string, config capib.Project) {

				server := http.NewServeMux()

				capib.ConfigStatic(&config.Static, server)

				server.HandleFunc("/", config.HandleFunc)

				newServer := &http.Server{
					//TLSNextProto:               make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
					Addr:    config.ListenAndServer.InAddress,
					Handler: server,
				}

				if config.DebugServerEnable == true {
					newServer.ErrorLog = log.New(capib.DebugLogger{}, "", 0)
				}

				capib.ConfigCertificates(config.Sll, newServer)

				if config.Sll.Certificate != "" && config.Sll.CertificateKey != "" {

					log.Fatal(newServer.ListenAndServeTLS(config.Sll.Certificate, config.Sll.CertificateKey))

				} else {

					log.Fatal(newServer.ListenAndServe())

				}

			}(projectName, projectConfig)

		} else if projectConfig.Listen.InAddress != "" {

			projectConfig.Listen.Pygocentrus = projectConfig.Pygocentrus
			log.Fatal(projectConfig.Listen.Listen())

		}
	}
}
