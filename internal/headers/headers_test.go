package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsingHeaders(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)

	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Parsing 2 headers"
	headers = NewHeaders()
	data = []byte("Host1: localhost:11111\r\nHost2: localhost:99999\r\n\r\n")
	n, done, err = headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)

	assert.Equal(t, "localhost:11111", headers["host1"])
	assert.Equal(t, 24, n)
	assert.False(t, done)

	data = data[n:]
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	assert.Equal(t, "localhost:99999", headers["host2"])
	assert.Equal(t, 24, n)
	assert.False(t, done)

	// Test: Valid done after reading 2 headers
	data = data[n:]
	n, done, err = headers.Parse(data)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Valid 2 headers with existing headers
	headers = map[string]string{"host": "localhost:42069"}
	data = []byte("User-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Invalid spacing header"
	headers = NewHeaders()
	data = []byte("Host : localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.ErrorContains(t, err, "whitespace between field name and colon detected")

	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid key caracter header"
	headers = NewHeaders()
	data = []byte("H©)t: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.ErrorContains(t, err, "not allowed character in header key")

	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Append to presend header
	headers = map[string]string{"set-example-header": "example-header-value1, example-header-value2"}
	data = []byte("Set-Example-Header: example-header-value3\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	assert.Equal(t, "example-header-value1, example-header-value2, example-header-value3", headers["set-example-header"])
	assert.Equal(t, 43, n)
	assert.False(t, done)

}
