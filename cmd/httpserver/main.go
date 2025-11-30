package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/MichalGul/http_server_go/internal/request"
	"github.com/MichalGul/http_server_go/internal/response"
	"github.com/MichalGul/http_server_go/internal/server"
)

const port = 42069

const BAD_REQUEST = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const SERVER_ERROR = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const OK = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

const PROXY_TARGET = "https://httpbin.org"

func basicHandler(w *response.Writer, req *request.Request) {

	path := req.RequestLine.RequestTarget

	if strings.HasPrefix(path, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	if path == "/yourproblem" {
		w.WriteStatusLine(response.BadRequestStatusCode)
		headers := response.GetDefaultHeaders(len(BAD_REQUEST))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody([]byte(BAD_REQUEST))

	} else if path == "/myproblem" {
		w.WriteStatusLine(response.InternalServerErrorStatusCode)
		headers := response.GetDefaultHeaders(len(SERVER_ERROR))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody([]byte(SERVER_ERROR))
	} else {
		w.WriteStatusLine(response.OkStatusCode)
		headers := response.GetDefaultHeaders(len(OK))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody([]byte(OK))
	}

}

func proxyHandler(w *response.Writer, req *request.Request) {

	path := req.RequestLine.RequestTarget

	trimmedPath := strings.TrimPrefix(path, "/httpbin/")
	proxyPath := fmt.Sprintf("%s/%s", PROXY_TARGET, trimmedPath)

	w.WriteStatusLine(response.OkStatusCode)
	headers := response.GetDefaultHeaders(0)

	delete(headers, "Content-Length")
	headers.Set("Transfer-Encoding", "chunked")

	// Write headers
	w.WriteHeaders(headers)

	// Request to proxy
	proxyResponse, err := http.Get(proxyPath)
	if err != nil {
		fmt.Printf("error: error calling proxy server %s", proxyPath)
		return
	}
	defer proxyResponse.Body.Close()

	//Handle response from proxy
	buf := make([]byte, 256)

	// handling reading streaming data
	for {
		// Read data from proxy info buffor
		n, err := proxyResponse.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Write done of the chunked body
				break
			}
			fmt.Printf("error: proxing data error %v", err)
			panic(err)
		}

		fmt.Printf("chunk (%d bytes): %s\n", n, buf[:n])
		w.WriteChunkedBody(buf[:n]) // Write read :n bytes to response
	}

	w.WriteChunkedBodyDone()


}

func main() {

	serv, err := server.Serve(port, basicHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer serv.Close()
	log.Println("Server started on port", port)

	// Gracefully shut down the server
	// Because server.Serve returns immediately (it handles requests in the background in goroutines)
	// if we exit main immediately, the server will just stop. We want to wait for a signal (like CTRL+C) before
	// we stop the server.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
