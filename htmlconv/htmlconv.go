package htmlconv

import (
	"bytes"
	"io"
	"net/http"
)

func StringToResponse(s string) *http.Response {
	t := http.Response{
		Body: io.NopCloser(bytes.NewBufferString(s)),
	}
	return &t
}

func ByteToResponse(b []byte) *http.Response {
	t := http.Response{
		Body: io.NopCloser(bytes.NewBuffer(b)),
	}
	return &t
}
