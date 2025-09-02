package httpadapter

import (
	"bytes"
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func New(handler http.Handler) backend.CallResourceHandler {
	return &httpResourceHandler{
		handler: handler,
	}
}

type httpResourceHandler struct {
	handler http.Handler
}

func (h *httpResourceHandler) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	var reqBodyReader io.Reader
	if len(req.Body) > 0 {
		reqBodyReader = bytes.NewReader(req.Body)
	}

	reqURL, err := url.Parse(req.URL)
	if err != nil {
		return err
	}

	resourceURL := req.Path
	if reqURL.RawQuery != "" {
		resourceURL += "?" + reqURL.RawQuery
	}

	if !strings.HasPrefix(resourceURL, "/") {
		resourceURL = "/" + resourceURL
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, resourceURL, reqBodyReader)
	if err != nil {
		return err
	}

	for key, values := range req.Headers {
		httpReq.Header[key] = values
	}

	writer := newResponseWriter(sender)
	h.handler.ServeHTTP(writer, httpReq)
	writer.close()

	return nil
}
