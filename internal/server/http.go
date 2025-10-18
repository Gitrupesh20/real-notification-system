package server

import (
	"log"
	"net/http"
)

func StartServer(addr string, handler http.Handler) {

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
