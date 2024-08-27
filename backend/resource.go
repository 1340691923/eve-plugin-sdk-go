package backend

import (
	"context"
)

type CallResourceRequest struct {
	PluginContext	PluginContext
	Path	string
	Method	string
	URL	string
	Headers	map[string][]string
	Body	[]byte
}

type CallResourceResponse struct {
	Status	int
	Headers	map[string][]string
	Body	[]byte
}

type CallResourceResponseSender interface {
	Send(*CallResourceResponse) error
}

type CallResourceHandler interface {
	CallResource(ctx context.Context, req *CallResourceRequest, sender CallResourceResponseSender) error
}

type CallResourceHandlerFunc func(ctx context.Context, req *CallResourceRequest, sender CallResourceResponseSender) error

func (fn CallResourceHandlerFunc) CallResource(ctx context.Context, req *CallResourceRequest, sender CallResourceResponseSender) error {
	return fn(ctx, req, sender)
}
