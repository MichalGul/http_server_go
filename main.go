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

const messagesFilename string = "messages.txt"
const port = ":42069"


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

	// Loading file
	// file, err := os.Open(messagesFilename)
	// if err != nil {
	// 	log.Fatalf("File %s failed to open: %v \n", messagesFilename, err)
	// }
	// defer file.Close()

	// Reading from tcp connection

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

		linesChannel := getLinesChannel(connection)
		for elem := range linesChannel {
			fmt.Printf("%s\n", elem)
		}

		fmt.Println("Connection to", connection.RemoteAddr(), "closed")
	

	}

}
