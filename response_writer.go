package accesslog

import (
	"bytes"
	"net/http"
	"time"
)

// ResponseWriter is an interface that extends http.ResponseWriter with
// additional methods that allow you to inspect the response that is being
// written to the client
type ResponseWriter interface {
	http.ResponseWriter
	BytesWritten() int64
	Status() int
	StartTime() time.Time
	// End is called when the downsteam handler has finished processing.
	End()
}

// ResponseWriterBuilder is an interface that allows you to wrap an existing
// http.ResponseWriter and return a ResponseWriter object
//
// By default `DefaultResponseWriterBuilder` is used.
type ResponseWriterBuilder interface {
	Wrap(http.ResponseWriter, *http.Request, bool) ResponseWriter
}

func DefaultResponseWriterBuilder() ResponseWriterBuilder {
	return defaultResponseWriterBuilder{}
}

type defaultResponseWriterBuilder struct{}

func (defaultResponseWriterBuilder) Wrap(w http.ResponseWriter, _ *http.Request, recordResponse bool) ResponseWriter {
	rw := &responseWriter{
		ResponseWriter: w,
		start:          time.Now(),
	}

	if recordResponse {
		rw.responseBody = &bytes.Buffer{}
	}
	return rw
}

type responseWriter struct {
	http.ResponseWriter
	bytesWritten int64
	code         int
	start        time.Time
	end          time.Time
	responseBody *bytes.Buffer
}

func (rw *responseWriter) End() {
	rw.end = time.Now()
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)

}

func (rw *responseWriter) Write(b []byte) (int, error) {
	written, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(written)
	if rb := rw.responseBody; rb != nil {
		rb.Write(b)
	}
	return written, err
}

func (rw *responseWriter) BytesWritten() int64 {
	return rw.bytesWritten
}

func (rw *responseWriter) Status() int {
	return rw.code
}

func (rw *responseWriter) StartTime() time.Time {
	return rw.start
}

func (rw *responseWriter) EndTime() time.Time {
	return rw.end
}

func (rw *responseWriter) Body() []byte {
	if rb := rw.responseBody; rb != nil {
		return rb.Bytes()
	}
	return nil
}
