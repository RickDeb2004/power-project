package main

import (
	"net/http"
)

// ResponseWriter is a custom implementation of http.ResponseWriter
type ResponseWriter struct {
	http.ResponseWriter
	Status int
	BodyContent    []byte
	// Add any additional fields or methods you need
}

// NewResponseWriter creates a new instance of ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		Status:     http.StatusOK,
		BodyContent:           nil,
		// Initialize any additional fields if needed
	}
}
func (rw *ResponseWriter) Write(content []byte) (int, error) {
	rw.BodyContent = content
	return len(content), nil
}

func (w *ResponseWriter) GetStatus() int {
	return w.Status
}
func(w*ResponseWriter)SetStatus(status int){
	w.Status=status
}

func (rw *ResponseWriter) Body() []byte {
	return rw.BodyContent
}
