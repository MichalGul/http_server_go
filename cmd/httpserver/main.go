package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/MichalGul/http_server_go/internal/headers"
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

	if path == "/video" {
		videoHandler(w, req)
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

func videoHandler(w *response.Writer, req *request.Request) {

	w.WriteStatusLine(response.OkStatusCode)
	file, err := os.ReadFile("/home/michal/workspace/httpfromtcp/http_server_go/assets/vim.mp4")

	if err != nil {
		handler500(w, req)
	}

	h := response.GetDefaultHeaders(len(file))
	h.Set("Content-Type", "video/mp4")

	w.WriteHeaders(h)
	w.WriteBody(file)
	return

}

func proxyHandler(w *response.Writer, req *request.Request) {

	path := req.RequestLine.RequestTarget

	trimmedPath := strings.TrimPrefix(path, "/httpbin/")
	proxyPath := fmt.Sprintf("%s/%s", PROXY_TARGET, trimmedPath)

	w.WriteStatusLine(response.OkStatusCode)
	h := response.GetDefaultHeaders(0)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	delete(h, "Content-Length")

	trailers := headers.NewHeaders()
	delete(trailers, "Content-Length")

	// Write headers
	err := w.WriteHeaders(h)
	if err != nil {
		fmt.Printf("error: error calling proxy server %s", proxyPath)
		return
	}

	// Request to proxy
	proxyResponse, err := http.Get(proxyPath)
	if err != nil {
		fmt.Printf("error: error calling proxy server %s", proxyPath)
		return
	}
	defer proxyResponse.Body.Close()

	//Handle response from proxy
	buf := make([]byte, 256)

	// Store whole message
	var fullMessage []byte

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

		fullMessage = append(fullMessage, buf[:n]...)
		_, err = w.WriteChunkedBody(buf[:n]) // Write read :n bytes to response
		if err != nil {
			fmt.Printf("error: writting data error %v", err)
		}
	}
	// Mark finishing writting body
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Printf("error: writting chunked done data error %v", err)
	}

	// Calculate hex32 of whole message
	messageSha := sha256.Sum256(fullMessage)
	hashHex := hex.EncodeToString(messageSha[:])

	msgLength := len(fullMessage)

	// Add trailers after body
	trailers.Set("X-Content-SHA256", hashHex)
	trailers.Set("X-Content-Length", strconv.Itoa(msgLength))

	// Write trailers to client
	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Printf("error: writting trailers %v", err)
	}

}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.InternalServerErrorStatusCode)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
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
