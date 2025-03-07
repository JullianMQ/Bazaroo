package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

var PORT = ":3000"

func ServeHttp() {
	mux := http.NewServeMux()
	corsOpts := cors.Options{
		AllowedOrigins: []string{ "*" },
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}

	handler := cors.New(corsOpts).Handler(mux)
	fmt.Printf("Listening at port %s\n", PORT)
	mux.HandleFunc("GET /", getRoot)
	if err := http.ListenAndServe(PORT, handler); err != nil {
		log.Fatal(err)
	}
}
