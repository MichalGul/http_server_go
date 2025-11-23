package server

import (
	"bytes"
	"fmt"
	"io"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) Write(w io.Writer) {

	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)

	// Error message
	w.Write(messageBytes)

}

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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		handleError := &HandlerError{
			StatusCode: response.BadRequestStatusCode,
			Message:    err.Error(),
		}
		handleError.Write(conn)
	}

	dataBuff := bytes.NewBuffer(nil)
	handlerError := s.handler(dataBuff, req)

	if handlerError != nil {
		handlerError.Write(conn)
		return
	}

	b := dataBuff.Bytes()
	response.WriteStatusLine(conn, response.OkStatusCode)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)

	// staticResponse := "HTTP/1.1 200 OK\r\n" + // Status line
	// 	"Content-Type: text/plain\r\n" + // Example header
	// 	"\r\n" + // Blank line to separate headers from the body
	// 	"Hello World!\n" // Body

	// conn.Write([]byte(staticResponse))

	return
}
