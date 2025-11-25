package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/MichalGul/http_server_go/internal/headers"
)

type StatusCode int

const (
	OkStatusCode StatusCode = 200
	BadRequestStatusCode StatusCode = 400
	InternalServerErrorStatusCode StatusCode = 500
)

type WriteState int

const (
	Initialize = iota
	StatusLineWrote
	HeadersWrote
	BodyWrote

)

type Writer struct {
	Connection io.Writer
	WriteState WriteState
}

func NewWritter(conn io.Writer) *Writer{
	return &Writer{
		Connection: conn,
		WriteState: Initialize,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.WriteState != Initialize {
		return fmt.Errorf("error: atempt to write to response in incorrect state")
	}

	byteStatusLine := getStatusLine(statusCode)
	_, err := w.Connection.Write(byteStatusLine)
	if err != nil{
		return err
	}
	w.WriteState = StatusLineWrote
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.WriteState != StatusLineWrote{
		return fmt.Errorf("error: atempt to write headers in incorrect state")
	}

	err := WriteHeaders(w.Connection, headers)
	if err != nil {
		return err
	}

	w.WriteState = HeadersWrote
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.WriteState != HeadersWrote {
		return 0, fmt.Errorf("error: atempt to write body in incorrect state")
	}

	w.WriteState = BodyWrote
	return w.Connection.Write(p)
}


func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case OkStatusCode:
		reasonPhrase = "OK"
	case BadRequestStatusCode:
		reasonPhrase = "Bad Request"
	case InternalServerErrorStatusCode:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}


func GetDefaultHeaders(contentLen int) headers.Headers {

	headers := headers.NewHeaders()
	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	
	for name, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n",name, value)))
		if err != nil {			
			return err
		}
	}
	w.Write([]byte("\r\n"))

	return nil
}