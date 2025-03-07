package main

import (
	"fmt"

	// "github.com/joho/godotenv"
	"github.com/JullianMQ/Bazaroo/server"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Real")
	server.ConnDB()
	server.ServeHttp()
}
