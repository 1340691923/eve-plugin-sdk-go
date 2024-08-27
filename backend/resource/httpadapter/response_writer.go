package httpadapter

import (
	"bytes"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	"net/http"
)

type callResourceResponseWriter struct {
	stream backend.CallResourceResponseSender

	Code int

	HeaderMap http.Header

	Body *bytes.Buffer

	Flushed bool

	wroteHeader     bool
	sentFirstStream bool
}

func newResponseWriter(stream backend.CallResourceResponseSender) *callResourceResponseWriter {
	return &callResourceResponseWriter{
		stream:    stream,
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

func (rw *callResourceResponseWriter) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

func (rw *callResourceResponseWriter) writeHeader(b []byte, str string) {
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

func (rw *callResourceResponseWriter) Write(buf []byte) (int, error) {
	rw.writeHeader(buf, "")
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

func (rw *callResourceResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.Code = code
	rw.wroteHeader = true
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(http.Header)
	}
}

func (rw *callResourceResponseWriter) Flush() {
	if !rw.wroteHeader {
		rw.WriteHeader(200)
	}

	resp := rw.getResponse()
	if resp != nil {
		if err := rw.stream.Send(resp); err != nil {
			logger.DefaultLogger.Error("Failed to send resource response", "error", err)
		}
	}

	rw.Body = new(bytes.Buffer)
}

func (rw *callResourceResponseWriter) getResponse() *backend.CallResourceResponse {
	if !rw.sentFirstStream {
		res := &backend.CallResourceResponse{
			Status:  rw.Code,
			Headers: rw.Header().Clone(),
		}

		if rw.Body != nil {
			res.Body = rw.Body.Bytes()
		}

		rw.sentFirstStream = true
		return res
	}

	if rw.Body != nil && rw.Body.Len() > 0 {
		return &backend.CallResourceResponse{
			Body: rw.Body.Bytes(),
		}
	}

	return nil
}

func (rw *callResourceResponseWriter) close() {
	rw.Flush()
}
