package capibaribe

import (
	"log"
	"net/http"
	"os"
)

func ConfigStatic(static *[]static, server *http.ServeMux) {
	var err error

	for _, staticPath := range *static {

		if _, err = os.Stat(staticPath.FilePath); os.IsNotExist(err) {
			log.Fatalf("static dir error: %v\n", err.Error())
		}

		server.Handle("/"+staticPath.ServerPath+"/", http.StripPrefix("/"+staticPath.ServerPath+"/", http.FileServer(http.Dir(staticPath.FilePath))))
	}

}
