package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(requestString string) (*RequestLine, error) {

	requestParts := strings.Split(requestString, "\r\n")
	requestLineRaw := requestParts[0]
	fmt.Printf("requestLineRaw: %s \n", requestLineRaw)

	requestLineParts := strings.Split(requestLineRaw, " ")

	fmt.Printf("requestLineParts: %s \n", requestLineParts)

	if len(requestLineParts) != 3 {
		return nil, errors.New("http request line missing all request line parts")
	}

	method := requestLineParts[0]
	requestTarget := requestLineParts[1]
	httpVersion := requestLineParts[2]

	if method != strings.ToUpper(method) {
		return nil, errors.New("http method not all in caps")
	}

	httpVersionsPart := strings.Split(httpVersion, "/")
	
	if len(httpVersion) <2 && string(httpVersionsPart[1]) != "1.1" {
		return nil, errors.New("bad http version")
	}
	httpVersion = string(httpVersionsPart[1])

	fmt.Printf("Method: %s \n requestTarget: %s \n httpVersion: %s \n", method, requestTarget, httpVersion)

	parsedRequestLine := RequestLine{
		HttpVersion: httpVersion,
		Method: method,
		RequestTarget: requestTarget,
	}

	return &parsedRequestLine, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	requestString, err := io.ReadAll(reader)

	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}

	fmt.Printf("Reader data: \n %s \n", string(requestString))
	requestLine, err := parseRequestLine(string(requestString))
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}

	request := Request{
		RequestLine: *requestLine,
	}

	return &request, nil
}
