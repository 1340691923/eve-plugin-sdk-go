package grpcplugin

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type TaskServer interface {
	pluginv2.TaskServer
}

type TaskClient interface {
	pluginv2.TaskClient
}

type TaskGRPCPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	plugin.GRPCPlugin
	TaskServer TaskServer
}

func (p *TaskGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pluginv2.RegisterTaskServer(s, &TaskGRPCServer{
		server: p.TaskServer,
	})
	return nil
}

func (p *TaskGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &TaskGRPCClient{client: pluginv2.NewTaskClient(c)}, nil
}

type TaskGRPCServer struct {
	server TaskServer
}

func (s *TaskGRPCServer) TaskExec(ctx context.Context, req *pluginv2.TaskRequest) (*pluginv2.TaskResponse, error) {
	return s.server.TaskExec(ctx, req)
}

type TaskGRPCClient struct {
	client pluginv2.TaskClient
}

func (m *TaskGRPCClient) TaskExec(ctx context.Context, in *pluginv2.TaskRequest, opts ...grpc.CallOption) (*pluginv2.TaskResponse, error) {
	return m.client.TaskExec(ctx, in, opts...)
}

var _ TaskServer = &TaskGRPCServer{}
var _ TaskClient = &TaskGRPCClient{}
