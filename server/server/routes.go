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
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
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
	mux.HandleFunc("POST /v1/emps/signup", PostEmployeeSignUp)
	mux.HandleFunc("POST /v1/emps/login", PostEmpLogin)

	mux.HandleFunc("GET /v1/vendors", GetVendors)
	mux.HandleFunc("GET /v1/vendors/", GetVendorId)
	mux.HandleFunc("POST /v1/vendors", PostVendor)

	mux.HandleFunc("GET /v1/prodlines", GetProductLine)
	mux.HandleFunc("POST /v1/prodlines", PostProductLine)

	mux.HandleFunc("GET /v1/products", GetProducts)
	mux.HandleFunc("POST /v1/products", PostProduct)

	mux.HandleFunc("GET /v1/customers", GetCustomers)
	mux.HandleFunc("POST /v1/customers", PostCustomer)
	mux.HandleFunc("POST /v1/customers/signup", PostCustomerSignUp)
	mux.HandleFunc("POST /v1/customers/login", PostCustomerLogIn)

	mux.HandleFunc("GET /v1/orders", GetOrders)
	mux.HandleFunc("GET /v1/orders/", GetOrderByCustId)
	mux.HandleFunc("POST /v1/orders", PostOrder)

	mux.HandleFunc("GET /v1/payments/", GetPaymentsByCustId)
	mux.HandleFunc("POST /v1/payments/", PostPayment)

	mux.HandleFunc("GET /v1/orderdetails/", GetOrderDetailsByOrderId)
	mux.HandleFunc("POST /v1/orderdetails", PostOrderDetail)

	if err := http.ListenAndServe(PORT, handler); err != nil {
		log.Fatal(err)
	}
}
