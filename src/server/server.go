package server

import (
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
	// Initialize the server with the given IP and port
	return &Server{
		port: port,
		ip:   ip,
		stop: make(chan bool),
	}
}

func (s *Server) Start() {
	// Create a listener on the given IP and port
	listener, err := net.Listen("tcp", s.ip+":"+s.port)
	if err != nil {
		log.Fatal(err)
	}

	// Close the listener when the server is stopped
	defer listener.Close()

	go func() {
		// Accept incoming connections
		for {
			// Accept a connection
			conn, err := listener.Accept()
			if err != nil {
				// If error is due to closed network connection, return
				if opErr, ok := err.(*net.OpError); ok &&
					opErr.Err.Error() == "use of closed network connection" {
					return
				}
				// Otherwise, log the error and continue
				log.Println(err)
				continue
			}

			// Create a buffer to read the request
			buffer := make([]byte, 1024)

			n, err := conn.Read(buffer)
			if err != nil {
				log.Println(err)
				continue
			}

			// Convert the buffer to a string
			request := string(buffer[:n])

			// Split the request into lines
			lines := strings.Split(request, "\n")

			// Parse the method and path from the request
			method := strings.Split(lines[0], " ")[0]
			path := strings.Split(lines[0], " ")[1]

			// If the method is not GET, return a 405 Method Not Allowed
			if method != "GET" {
				conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n"))
				conn.Write([]byte("Allow: GET\r\n"))
				conn.Write([]byte("\r\n"))
				conn.Close()
				continue
			}

			// Find the handler for the given path
			found := false
			for _, handler := range s.handlers {
				if handler.path == path {
					found = true
					go handler.handler(conn)
					break
				}
			}

			// If no handler is found, return a 404 Not Found
			if found == false {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
				conn.Write([]byte("\r\n"))
				conn.Close()
			}
		}
	}()

	// Wait for the server to be stopped
	<-s.stop
}

func (s *Server) Get(path string, handler func(conn net.Conn)) {
	// Add a handler for the given path
	s.handlers = append(s.handlers, Handler{path: path, handler: handler})
}

func (s *Server) Stop() {
	// Send a signal to stop the server
	s.stop <- true
}
