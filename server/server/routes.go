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
	mux.HandleFunc("GET /", GetRoot)

	mux.HandleFunc("GET /v1/addr", GetAddr)
	mux.HandleFunc("POST /v1/addr", PostAddr)

	mux.HandleFunc("GET /v1/offices", GetOffices)
	mux.HandleFunc("POST /v1/offices", PostOffice)

	mux.HandleFunc("GET /v1/emps", GetEmps)
	mux.HandleFunc("GET /v1/emps/", GetEmpId)
	mux.HandleFunc("POST /v1/emps", PostEmp)

	mux.HandleFunc("GET /v1/vendors", GetVendors)
	mux.HandleFunc("GET /v1/vendors/", GetVendorId)
	mux.HandleFunc("POST /v1/vendors", PostVendor)

	mux.HandleFunc("GET /v1/prodlines", GetProductLine)
	mux.HandleFunc("POST /v1/prodlines", PostProductLine)

	if err := http.ListenAndServe(PORT, handler); err != nil {
		log.Fatal(err)
	}
}
