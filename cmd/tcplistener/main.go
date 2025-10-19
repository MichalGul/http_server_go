package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"net"
	"log"
)

const port = ":42069"

// Returns a chanel to read data from
// Runs goorutine that reads from io.Read in chunks of 8 bytes
// Goorutine searchs for new line and then pust while line into
// channel. Closes channel when there is no more data in io.ReadCloser to read
// note: io can be anything from file to tcp connection
func getLinesChannel(f io.ReadCloser) <- chan string {

	lines := make(chan string)

	// Goroutine to write file contents to channel line by line
	go func () {
		defer close(lines)
		defer f.Close()
		var currentLine string = ""
		for {
			dataChunk := make([]byte, 8)
			numOfBytes, readError := f.Read(dataChunk)
			if readError != nil {
				if errors.Is(readError, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", readError.Error())
				return
			}

			dataString := string(dataChunk[:numOfBytes])
			currentLine += dataString

			// Read until new line character, then place whole line
			// into clannel and accumulate next line.
			parts := strings.Split(currentLine, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- parts[i]
			}

			currentLine = parts[len(parts)-1]

		}
		if currentLine != "" {
			lines <- currentLine
		}
	}()

	return lines
}



func main() {

	// Opening TCP connection
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	for {
		// Wait for a connection.
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Connecting error: %s\n", err.Error())
			os.Exit(0)
		}
		fmt.Println("Accepted connection from", connection.RemoteAddr())

		// Get channel and read data from it
		linesChannel := getLinesChannel(connection)
		for elem := range linesChannel {
			fmt.Printf("%s\n", elem)
		}

		fmt.Println("Connection to", connection.RemoteAddr(), "closed")
	

	}

}
