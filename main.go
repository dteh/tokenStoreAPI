package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dteh/tokenStoreAPI/store"
)

func main() {
	DB, err := store.CreateOrGetDB("tokens.db")
	if err != nil {
		log.Fatal("Unable to access database", err)
		os.Exit(1)
	}

	// set up api routes
	h := store.NewAPIHandler(DB)
	http.HandleFunc("/tokens", h.TokenEndpoint)

	// start the http server
	go http.ListenAndServe("0.0.0.0:1234", nil)
}
