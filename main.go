package main

import (
	capib "./capibaribe"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	var err error
	var wg sync.WaitGroup

	filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
	flag.Parse()

	config := capib.MainConfig{}
	err = config.Unmarshal(*filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	for projectName, projectConfig := range config.AffluentRiver {

		_ = projectName

		wg.Add(1)
		go func(config capib.Project) {
			var err error
			server := http.NewServeMux()

			defer wg.Done()

			for _, staticPath := range config.Static {

				if _, err = os.Stat(staticPath.FilePath); os.IsNotExist(err) {
					log.Fatalf("static dir error: %v\n", err.Error())
				}

				server.Handle("/"+staticPath.ServerPath+"/", http.StripPrefix("/"+staticPath.ServerPath+"/", http.FileServer(http.Dir(staticPath.FilePath))))
			}

			server.HandleFunc("/", config.HandleFunc)

			newServer := &http.Server{
				Addr:     config.Listen,
				Handler:  server,
				ErrorLog: log.New(capib.DebugLogger{}, "", 0),
			}
			log.Fatal(newServer.ListenAndServe())

		}(projectConfig)

	}

	wg.Wait()
}
