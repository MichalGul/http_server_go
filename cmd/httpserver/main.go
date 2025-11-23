package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MichalGul/http_server_go/internal/request"
	"github.com/MichalGul/http_server_go/internal/response"
	"github.com/MichalGul/http_server_go/internal/server"
)

const port = 42069

func basicHandler(w io.Writer, req *request.Request) *server.HandlerError {

	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.BadRequestStatusCode,
			Message: "Your problem is not my problem\n",
		}

	} else if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.InternalServerErrorStatusCode,
			Message: "Woopsie, my bad\n",
		}
	} else {
		w.Write([]byte("All good, frfr\n"))
		return nil
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
