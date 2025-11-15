package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/MichalGul/http_server_go/internal/response"
)

type Server struct {
	isClosed           atomic.Bool
	connectionListener net.Listener
}

func Serve(port int) (*Server, error) {

	portStr := strconv.Itoa(port)
	listener, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		return nil, fmt.Errorf("error creating listener %v", err)
	}

	var isClosed atomic.Bool
	isClosed.Store(false)

	server := &Server{
		isClosed:           isClosed,
		connectionListener: listener,
	}

	// Accept listen for connections in gorutine
	go server.listen()

	return server, nil

}

func (s *Server) Close() error {

	s.isClosed.Store(true)
	if s.connectionListener != nil {
		return s.connectionListener.Close()
	}
	return nil
}

func (s *Server) listen() {

	for {
		connection, err := s.connectionListener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return // if server is closed ignore errors
			}
			fmt.Printf("accept error: %v", err)
			continue
		}

		go s.handle(connection)

	}

}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	response.WriteStatusLine(conn, response.OkStatusCode)
	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, headers)


	// staticResponse := "HTTP/1.1 200 OK\r\n" + // Status line
	// 	"Content-Type: text/plain\r\n" + // Example header
	// 	"\r\n" + // Blank line to separate headers from the body
	// 	"Hello World!\n" // Body

	// conn.Write([]byte(staticResponse))
	return
}
