package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"errors"
)

type RequestParsingState int

const (
	Initialized RequestParsingState = iota
	Done
)

type Request struct {
	RequestLine  RequestLine
	ParsingState RequestParsingState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

//
//
//
func (r *Request) parse (data []byte) (int, error) {

	switch r.ParsingState {
		case Initialized:
			numOfBytes, requestLine, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}

			if numOfBytes == 0 && err == nil {
				// needs more data from the stream
				return 0, nil
			}

			// Succesfuly parsed bytes
			r.RequestLine = *requestLine
			r.ParsingState = Done

			return numOfBytes, nil


		case Done:
			return 0, fmt.Errorf("error: trying to read data in a done state")

		default:
			return 0, fmt.Errorf("error: unknown request parsting state")

	}

}



const crlf = "\r\n"
const streamBufferSize = 8

func parseRequestLine(data []byte) (int, *RequestLine, error) {

	// Find endline /r/n so everything until first CR on http request
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		// return nil, fmt.Errorf("could not find CRLF in request-line")\
		return 0, nil, nil // needs more byte to read
	}
	requestLineText := string(data[:idx])

	// Rest of message after first request line, meaning headers and body separaterd wit /r/n
	// bytesRead := idx + len(crlf)
	// resetOfMessage := string(data[idx + len(crlf)])

	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return 0, nil, err
	}
	return idx +2, requestLine, nil
}

// Parse http request as string to RequestLine object
// Perform structure checs for HTTP request standard
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

// Main method to parse incoming data through tcp connection
// reads from io.Reader that is tcp connection or file
// It uses []byte as buffor for data with set streamBufferSize
// Create Request object with Init state, Check for done state in loop,
// Read bytes from io.Reader, to buffer, acknowledge number of bytes read
// atempt to parse bytes to RequestLine, and move buffor 
// readingRequest.parse determines if whole Request line was read and changes state to Done
func RequestFromReader(reader io.Reader) (*Request, error) {

	// Buffer chunk size to read data from stream (by streamBufferSize bytes at the time untile streaming data is finished)
	databuffor := make([]byte, streamBufferSize, streamBufferSize)
	readToIndex := 0 // track how much data read from io.Reader into the buffer
	readingRequest := &Request{
		ParsingState: Initialized,
	}


	for readingRequest.ParsingState != Done {

		if readToIndex >= len(databuffor){ // when buffor is full
			// make new slice with capacity x2 and copy data
			newBuffor := make([]byte, 2*len(databuffor))
			copy(newBuffor, databuffor)
			databuffor = newBuffor
		}

		numOfBytesRead, readError := reader.Read(databuffor[readToIndex:])
		if readError != nil {
			if errors.Is(readError, io.EOF) { // Read all
				readingRequest.ParsingState = Done
				break
			}
		}

		endReadIndex := readToIndex + numOfBytesRead
		readToIndex = endReadIndex

		// numOfParsedBytes will be 0 untile whole request line present
		numOfParsedBytes, parseError := readingRequest.parse(databuffor[:endReadIndex])
		if parseError != nil {
			return nil, parseError
		}

		// Move past by read data, w don't need them in buffor.
		copy(databuffor, databuffor[numOfParsedBytes:])
		readToIndex -= numOfParsedBytes

	}

	return readingRequest, nil
}
