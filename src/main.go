package main

import (
	"net"

	"adventune/the-bad-webserver/server"
)

func main() {
	server := server.NewServer("127.0.0.1", "8080")

	server.Get("/", func(conn net.Conn) {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/html\r\n"))
		conn.Write([]byte("\r\n"))
		conn.Write([]byte("<h1>Hello, World!</h1>"))

		conn.Close()
	})

	server.Start()
}
