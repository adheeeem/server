package main

import (
	"log"
	"net"
	"os"
	"server/pkg/server"
	"strconv"
)

func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", func(conn net.Conn) {
		body := "Welcome to our web-site"
		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: 	text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				body,
		))
		if err != nil {
			log.Print(err)
		}
	})
	srv.Register("/about", func(conn net.Conn) {
		body := "About Golang Academy"
		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: 	text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				body,
		))
		if err != nil {
			log.Print(err)
		}
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
