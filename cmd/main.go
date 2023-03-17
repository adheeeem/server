package main

import (
	"log"
	"net"
	"os"
	"server/pkg/server"
)

func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/payments/{id}", func(req *server.Request) {
		id := req.PathParams["id"]
		foo := req.QueryParams["id"]
		log.Print(id, foo, req.Headers)
	})

	return srv.Start()
}

func main() {
	host := "localhost"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}
