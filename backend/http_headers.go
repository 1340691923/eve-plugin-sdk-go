package backend

import (
	"fmt"
	"net/http"
	"net/textproto"
	"strings"
)

const (
	OAuthIdentityTokenHeaderName	= "Authorization"

	OAuthIdentityIDTokenHeaderName	= "X-Id-Token"

	CookiesHeaderName	= "Cookie"

	httpHeaderPrefix	= "http_"
)

type ForwardHTTPHeaders interface {
	SetHTTPHeader(key, value string)

	DeleteHTTPHeader(key string)

	GetHTTPHeader(key string) string

	GetHTTPHeaders() http.Header
}

func setHTTPHeaderInStringMap(headers map[string]string, key string, value string) {
	if headers == nil {
		headers = map[string]string{}
	}

	headers[fmt.Sprintf("%s%s", httpHeaderPrefix, key)] = value
}

func getHTTPHeadersFromStringMap(headers map[string]string) http.Header {
	httpHeaders := http.Header{}

	for k, v := range headers {
		if textproto.CanonicalMIMEHeaderKey(k) == OAuthIdentityTokenHeaderName {
			httpHeaders.Set(k, v)
		}

		if textproto.CanonicalMIMEHeaderKey(k) == OAuthIdentityIDTokenHeaderName {
			httpHeaders.Set(k, v)
		}

		if textproto.CanonicalMIMEHeaderKey(k) == CookiesHeaderName {
			httpHeaders.Set(k, v)
		}

		if strings.HasPrefix(k, httpHeaderPrefix) {
			hKey := strings.TrimPrefix(k, httpHeaderPrefix)
			httpHeaders.Set(hKey, v)
		}
	}

	return httpHeaders
}

func deleteHTTPHeaderInStringMap(headers map[string]string, key string) {
	for k := range headers {
		if textproto.CanonicalMIMEHeaderKey(k) == textproto.CanonicalMIMEHeaderKey(key) ||
			textproto.CanonicalMIMEHeaderKey(k) == textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", httpHeaderPrefix, key)) {
			delete(headers, k)
			break
		}
	}
}
