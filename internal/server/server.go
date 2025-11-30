package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/MichalGul/http_server_go/internal/request"
	"github.com/MichalGul/http_server_go/internal/response"
)

type Server struct {
	isClosed           atomic.Bool
	connectionListener net.Listener
	handler 		   Handler
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

// func (he HandlerError) Write(w io.Writer) {

// 	response.WriteStatusLine(w, he.StatusCode)
// 	messageBytes := []byte(he.Message)
// 	headers := response.GetDefaultHeaders(len(messageBytes))
// 	response.WriteHeaders(w, headers)

// 	// Error message
// 	w.Write(messageBytes)

// }

func Serve(port int, handler Handler) (*Server, error) {

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
		handler: handler,
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

	responseWritter := response.NewWritter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		responseWritter.WriteStatusLine(response.BadRequestStatusCode)
		defaultHeaders := response.GetDefaultHeaders(len(err.Error()))
		responseWritter.WriteHeaders(defaultHeaders)
		responseWritter.WriteBody([]byte(err.Error()))
		return
	}

	s.handler(responseWritter, req)

	return
}
