package main

import (
	"log"
	"net"
	"bufio"
	"os"
	"fmt"
)

func main() {
	serverAddr := "localhost:42069"

	destinationConnection, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("error resolving for UDP traffic: %s\n", err.Error())
	}

	udpConnection, connectionError := net.DialUDP("udp", nil, destinationConnection)
	if connectionError != nil {
		log.Fatalf("error connecting to UDP traffic: %s\n", connectionError.Error())

	}
	defer udpConnection.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddr)
	stdInReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		data, err := stdInReader.ReadString('\n') // reads data until new line form standard input
		if err!= nil {
			fmt.Printf("Reading from standard in error: %s\n", err.Error())
			os.Exit(1)
		}

		_, err = udpConnection.Write([]byte(data))
		if err != nil {
			fmt.Printf("Writting to connection : %s\n", err.Error())
		}

		fmt.Printf("Message sent: %s", data)
	}

}
