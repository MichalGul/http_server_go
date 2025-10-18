package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const messagesFilename string = "messages.txt"


func main() {

	// Loading file
	file, err := os.Open(messagesFilename)
	if err != nil {
		log.Fatalf("File %s failed to open: %v \n", messagesFilename, err)
	}
	defer file.Close()

	var currentLine string = ""

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

		dataString := string(dataChunk[:numOfBytes])
		currentLine += dataString

		parts := strings.Split(currentLine, "\n")

		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s\n", parts[i])
		}

		currentLine = parts[len(parts)-1]

	}

	if currentLine != "" {
		fmt.Printf("read: %s\n", currentLine)
	}
}
