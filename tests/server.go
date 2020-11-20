package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory")
)

func main() {
	flag.Parse()
	log.Printf("Open your browser in http://localhost%v...", *listen)
	log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir))))
}
