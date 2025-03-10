package main

import (
	"sync"

	"github.com/JullianMQ/Bazaroo/server"
	_ "github.com/lib/pq"
)

var wg = sync.WaitGroup{}

func main() {
	server.ConnDB()
	server.ServeHttp()
}
