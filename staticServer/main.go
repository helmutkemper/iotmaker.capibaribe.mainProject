package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	serverAndPortFlag := flag.String("server", "127.0.0.1:3000", "This option set server and port. Ex.: :3000")
	staticDirPathFlag := flag.String("staticPath", "/docker/static", "This option set the file path to the static files server. Ex.: /docker/static")

	flag.Parse()

	fs := http.FileServer(http.Dir(*staticDirPathFlag))
	http.Handle("/", fs)

	fmt.Println("Start config info:")
	fmt.Printf("- Web server and port: %v\n", *serverAndPortFlag)
	fmt.Printf("- Static file path: %v\n", *staticDirPathFlag)
	fmt.Println()
	fmt.Println("> Starting http server...")
	fmt.Println()
	fmt.Println()

	if err := http.ListenAndServe(*serverAndPortFlag, nil); err != nil {
		log.Fatal(err.Error())
	}
}
