package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MichalGul/http_server_go/internal/request"
	"github.com/MichalGul/http_server_go/internal/response"
	"github.com/MichalGul/http_server_go/internal/server"
)

const port = 42069

const BAD_REQUEST =`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const SERVER_ERROR =`<html>
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

func basicHandler(w *response.Writer, req *request.Request) {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		w.WriteStatusLine(response.BadRequestStatusCode)
		headers := response.GetDefaultHeaders(len(BAD_REQUEST))
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody([]byte(BAD_REQUEST))

	} else if req.RequestLine.RequestTarget == "/myproblem" {
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
