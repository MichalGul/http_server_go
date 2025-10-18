package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const messagesFilename string = "messages.txt" 

func main() {

	// Loading file
	file, err := os.Open(messagesFilename)
	if err != nil {
		log.Fatalf("File %s failed to open: %v \n", messagesFilename, err)
	}
	defer file.Close()


	for {
		dataChunk := make([]byte, 8)
		numOfBytes, readError := file.Read(dataChunk)

		if readError != nil {
			if errors.Is(readError, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		fmt.Printf("read: %s\n", dataChunk[:numOfBytes])
		
	}

}
