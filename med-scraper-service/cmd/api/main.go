package main

import (
	"log"
	"net/http"
)

const (
	webPort = ":80"
)

func main() {
	srv := &http.Server{
		Addr:    webPort,
		Handler: routes(),
	}

	log.Println("Starting medScraperService at port", webPort)

	err := srv.ListenAndServe()

	log.Fatal(err)
}
