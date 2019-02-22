package apiplugin

import (
	"bytes"
	"net/http"
)

// WrapResponseWriter is wrap of ResponseWriter, and implements io.Reader
type WrapResponseWriter struct {
	body        *bytes.Buffer
	code        int
	header      http.Header
	wroteHeader bool
}

// NewWrapResponseWriter return an instance of WrapResponseWriter
func NewWrapResponseWriter() *WrapResponseWriter {
	return &WrapResponseWriter{
		body:   new(bytes.Buffer),
		header: make(http.Header),
	}
}

func (rw *WrapResponseWriter) writeHeader(b []byte, str string) {
	if rw.wroteHeader {
		return
	}
	if len(str) > 512 {
		str = str[:512]
	}

	m := rw.Header()

	_, hasType := m["Content-Type"]
	hasTE := m.Get("Transfer-Encoding") != ""
	if !hasType && !hasTE {
		if b == nil {
			b = []byte(str)
		}
		m.Set("Content-Type", http.DetectContentType(b))
	}

	rw.WriteHeader(200)
}

// Header implements ResponseWriter.Header
func (rw *WrapResponseWriter) Header() http.Header {
	m := rw.header
	if m == nil {
		m = make(http.Header)
		rw.header = m
	}
	return m
}

// Write implements ResponseWriter.Write
func (rw *WrapResponseWriter) Write(buf []byte) (int, error) {
	rw.writeHeader(buf, "")
	if rw.body != nil {
		rw.body.Write(buf)
	}
	return len(buf), nil
}

// WriteHeader implements ResponseWriter.WriteHeader
func (rw *WrapResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.code = code
	rw.wroteHeader = true
	if rw.header == nil {
		rw.header = make(http.Header)
	}
}

// Read implements io.Reader.Read
func (rw *WrapResponseWriter) Read(p []byte) (int, error) {
	return rw.body.Read(p)
}

// Code returns status code
func (rw *WrapResponseWriter) Code() int {
	return rw.code
}
