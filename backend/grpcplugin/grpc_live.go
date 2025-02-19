package grpcplugin

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type LiveServer interface {
	pluginv2.LiveServer
}

type LiveClient interface {
	pluginv2.LiveClient
}

type LiveGRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	plugin.GRPCPlugin
	LiveServer LiveServer
}

func (p *LiveGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pluginv2.RegisterLiveServer(s, &LiveGRPCServer{
		server: p.LiveServer,
	})
	return nil
}

func (p *LiveGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &LiveGRPCClient{client: pluginv2.NewLiveClient(c)}, nil
}

type LiveGRPCServer struct {
	server LiveServer
}

func (s *LiveGRPCServer) Pub2Channel(ctx context.Context, req *pluginv2.Pub2ChannelRequest) (*pluginv2.Pub2ChannelResponse, error) {
	return s.server.Pub2Channel(ctx, req)
}

type LiveGRPCClient struct {
	client pluginv2.LiveClient
}

func (m *LiveGRPCClient) Pub2Channel(ctx context.Context, in *pluginv2.Pub2ChannelRequest, opts ...grpc.CallOption) (*pluginv2.Pub2ChannelResponse, error) {
	return m.client.Pub2Channel(ctx, in, opts...)
}

var _ LiveServer = &LiveGRPCServer{}
var _ LiveClient = &LiveGRPCClient{}
