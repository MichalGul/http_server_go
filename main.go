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


func getLinesChannel(f io.ReadCloser) <- chan string {

	lines := make(chan string)

	go func () {
		var currentLine string = ""
		for {
			dataChunk := make([]byte, 8)
			numOfBytes, readError := f.Read(dataChunk)
			if readError != nil {
				if errors.Is(readError, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", readError.Error())
				break
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
		close(lines)
		f.Close()
	}()

	return lines
}



func main() {

	// Loading file
	file, err := os.Open(messagesFilename)
	if err != nil {
		log.Fatalf("File %s failed to open: %v \n", messagesFilename, err)
	}
	defer file.Close()

	linesChannel := getLinesChannel(file)

	for elem := range linesChannel {
		fmt.Printf("read: %s\n", elem)

	}


}
