package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct{}

const PORT = "80"

func main() {
	app := Config{}

	log.Println("Starting mail service on port", PORT)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
