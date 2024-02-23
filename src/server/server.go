package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	port     string
	ip       string
	stop     chan bool
	handlers []Handler
}

type Handler struct {
	path    string
	handler func(conn net.Conn)
}

func NewServer(ip, port string) *Server {
	return &Server{
		port: port,
		ip:   ip,
		stop: make(chan bool),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.ip+":"+s.port)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println("Server started on", s.ip+":"+s.port)

	for {
		select {
		case <-s.stop:
			break
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				log.Println(err)
				continue
			}

			request := string(buffer[:n])

			lines := strings.Split(request, "\n")

			method := strings.Split(lines[0], " ")[0]
			path := strings.Split(lines[0], " ")[1]

			if method != "GET" {
				conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n"))
				conn.Write([]byte("Allow: GET\r\n"))
				conn.Write([]byte("\r\n"))
				conn.Close()
				continue
			}

			for _, handler := range s.handlers {
				if handler.path == path {
					go handler.handler(conn)
					break
				}
			}
		}
	}
}

func (s *Server) Get(path string, handler func(conn net.Conn)) {
	s.handlers = append(s.handlers, Handler{path: path, handler: handler})
}

func (s *Server) Stop() {
	s.stop <- true
}
