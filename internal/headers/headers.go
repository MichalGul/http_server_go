package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

var crlf = []byte("\r\n")

func NewHeaders() Headers {
	return Headers{}
}

func IsValidHeaderName(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		// Only ASCII allowed
		if r > 127 {
			return false
		}
		// No controls or whitespace
		if r <= 31 || r == 127 || unicode.IsSpace(r) {
			return false
		}
		switch r {
		// Disallowed separators per HTTP token rule
		case '(', ')', '<', '>', '@', // angles/at
			',', ';', ':', // comma/semicolon/colon
			'\\', '"', '/', // backslash/quote/slash
			'[', ']', '?', '=', // brackets/question/equals
			'{', '}', // braces
			' ':
			return false
		}
		// Otherwise allowed tchar (ALPHA/DIGIT and listed symbols)
	}
	return true
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2) // case for : in field value
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}

	if strings.HasSuffix(string(parts[0]), " ") {
		return "", "", fmt.Errorf("whitespace between field name and colon detected. Malformed header")
	}

	headerNameBytes := bytes.TrimSpace(parts[0])
	headerValue := bytes.TrimSpace(parts[1])

	headerName := string(headerNameBytes)

	if !IsValidHeaderName(headerName) {
		return "", "", fmt.Errorf("not allowed character in header key")
	}

	return strings.ToLower(headerName), string(headerValue), nil

}

// Parse raw string headers to Headers map
// Headers structure: ```field-line   = field-name ":" OWS field-value OWS``` OWS whitespaces zero or more
// gets header data in bytes parses it according to headers structure and check if last crlf was found meainng end of headers.
// Caller will handle multiple calls to parse multiple headers
func (h Headers) Parse(data []byte) (int, bool, error) {

	dataRead := 0
	done := false

	index := bytes.Index(data, crlf)
	if index == -1 { // did not find clrf not enough data for parsing
		return 0, done, nil
	}

	// Empty Header, so we have read everything we needed to read, consume CRLF
	if index == 0 {
		done = true
		return len(crlf), done, nil
	}

	headerBytes := data[:index]
	name, value, err := parseHeader(headerBytes)
	if err != nil {
		return 0, false, err
	}

	// Set header
	h[name] = value
	dataRead += len(headerBytes)

	return index + len(crlf), false, nil //index plus len(crls)

}
