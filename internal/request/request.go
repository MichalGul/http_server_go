package request

import (
	"bytes"
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

const crlf = "\r\n"

func parseRequestLine(data []byte) (*RequestLine, error) {
	
	// Find endline /r/n so everything until first CR on http request
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}
	return requestLine, nil
}

func requestLineFromString(requestLineString string) (*RequestLine, error) {

	fmt.Printf("requestLineRaw: %s \n", requestLineString)

	requestLineParts := strings.Split(requestLineString, " ")
	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", requestLineParts)
	}

	method := requestLineParts[0]
	requestTarget := requestLineParts[1]
	httpVersion := requestLineParts[2]

	if method != strings.ToUpper(method) {
		return nil, fmt.Errorf("invalid http method: %s", method)
	}

	httpVersionsPart := strings.Split(httpVersion, "/")

	if len(httpVersionsPart) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", requestLineString)
	}

	httpPart := httpVersionsPart[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := httpVersionsPart[1]
	if version != "1.1" {
		return nil, fmt.Errorf("bad http version")
	}

	fmt.Printf("Method: %s \n requestTarget: %s \n httpVersion: %s \n", method, requestTarget, httpVersion)

	parsedRequestLine := RequestLine{
		HttpVersion:   version,
		Method:        method,
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
	requestLine, err := parseRequestLine(requestString)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}

	request := Request{
		RequestLine: *requestLine,
	}

	return &request, nil
}
