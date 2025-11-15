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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
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