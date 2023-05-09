package main

import (
	"log"
	"net/http"
)

const webPort = ":80"

func main() {
	srv := &http.Server{
		Handler: routes(),
		Addr:    webPort,
	}

	log.Println("Starting medApiService at port:", webPort)

	err := srv.ListenAndServe()

	log.Fatal(err)
}
