package grpcplugin

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type PluginInfoServer interface {
	pluginv2.PluginInfoServer
}

type PluginInfoClient interface {
	pluginv2.PluginInfoClient
}

type PluginInfoGRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	plugin.GRPCPlugin
	PluginInfoServer PluginInfoServer
}

func (p *PluginInfoGRPCPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	pluginv2.RegisterPluginInfoServer(s, &pluginInfoGRPCServer{
		server: p.PluginInfoServer,
	})
	return nil
}

func (p *PluginInfoGRPCPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &pluginInfoGRPCClient{client: pluginv2.NewPluginInfoClient(c)}, nil
}

type pluginInfoGRPCServer struct {
	server PluginInfoServer
}

func (s *pluginInfoGRPCServer) Get(ctx context.Context, req *pluginv2.PluginInfoGetReq) (*pluginv2.PluginInfoGetRes, error) {
	return s.server.Get(ctx, req)
}

func (s *pluginInfoGRPCServer) CheckHealth(ctx context.Context, req *pluginv2.CheckHealthRequest) (*pluginv2.CheckHealthResponse, error) {
	return s.server.CheckHealth(ctx, req)
}

type pluginInfoGRPCClient struct {
	client pluginv2.PluginInfoClient
}

func (s *pluginInfoGRPCClient) Get(ctx context.Context, req *pluginv2.PluginInfoGetReq, opts ...grpc.CallOption) (*pluginv2.PluginInfoGetRes, error) {
	return s.client.Get(ctx, req, opts...)
}

func (s *pluginInfoGRPCClient) CheckHealth(ctx context.Context, req *pluginv2.CheckHealthRequest, opts ...grpc.CallOption) (*pluginv2.CheckHealthResponse, error) {
	return s.client.CheckHealth(ctx, req, opts...)
}

var _ PluginInfoServer = &pluginInfoGRPCServer{}
var _ PluginInfoClient = &pluginInfoGRPCClient{}
