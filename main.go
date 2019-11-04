//todo: tem que ter um log avisando quando foi a ultima requisição
//todo: microsservico deveria ter um id inicial que permita seguir uma requisição em caso de erro
package main

import (
	capib "./capibaribe"

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
	flag.Parse()

	loadConf(*filePath)

	//etcdTest()

	wg.Wait()
}

func loadConf(filePath string) {
	var err error
	err = config.Unmarshal(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for projectName, projectConfig := range config.AffluentRiver {

		// fixme: isto tem lógica se ficar em branco?
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

				// fixme: melhorar isto
				// enabled == true e sem certificados é um erro
				if config.Sll.Enabled == true && config.Sll.Certificate != "" && config.Sll.CertificateKey != "" {

					log.Fatal(newServer.ListenAndServeTLS(config.Sll.Certificate, config.Sll.CertificateKey))

				} else {

					log.Fatal(newServer.ListenAndServe())

				}

			}(projectName, projectConfig)

			// todo: pygocentrus não deveria está aqui
		} else if projectConfig.Listen.InAddress != "" {

			//projectConfig.Listen.Pygocentrus = projectConfig.Pygocentrus
			log.Fatal(projectConfig.Listen.Listen())

		}
	}
}
