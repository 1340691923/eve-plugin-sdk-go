package backend

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/backend/grpcplugin"
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const defaultServerMaxReceiveMessageSize = 1024 * 1024 * 16

type GRPCSettings struct {
	MaxReceiveMsgSize int

	MaxSendMsgSize int
}

type ServeOpts struct {
	pluginID string

	pluginVersion string

	PluginJson *build.PluginJsonData

	CheckHealthHandler CheckHealthHandler

	PluginInfoHandler PluginInfoHandler

	CallResourceHandler CallResourceHandler

	GRPCSettings GRPCSettings

	debug bool

	EvRpcPort string

	ExitCallback func()
}

type DefaultPluginInfoHandler struct {
	PluginID      string
	PluginVersion string
}

func NewDefaultPluginInfoHandler(pluginID string, pluginVersion string) *DefaultPluginInfoHandler {
	return &DefaultPluginInfoHandler{PluginID: pluginID, PluginVersion: pluginVersion}
}

func (d DefaultPluginInfoHandler) PluginInfo(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error) {
	return &PluginInfoGetRes{
		PluginID:      d.PluginID,
		PluginVersion: d.PluginVersion,
	}, nil
}

func asGRPCServeOpts(opts ServeOpts) grpcplugin.ServeOpts {

	if opts.PluginInfoHandler == nil {
		opts.PluginInfoHandler = NewDefaultPluginInfoHandler(opts.pluginID, opts.pluginID)
	}

	pluginOpts := grpcplugin.ServeOpts{
		PluginInfoServer: newPluginInfoSDKAdapter(opts.PluginInfoHandler, opts.CheckHealthHandler),
	}

	if opts.CallResourceHandler != nil {
		pluginOpts.ResourceServer = newResourceSDKAdapter(opts.CallResourceHandler)
	}

	return pluginOpts
}

func Serve(opts ServeOpts) {
	var pluginJson = opts.PluginJson
	ev_api.SetEvApi(opts.EvRpcPort, pluginJson.PluginAlias, pluginJson.BackendDebug)
	opts.debug = pluginJson.BackendDebug
	opts.pluginID = pluginJson.PluginAlias
	opts.pluginVersion = pluginJson.Version
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpcMiddlewares := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
		)),
	}

	if opts.GRPCSettings.MaxReceiveMsgSize <= 0 {
		opts.GRPCSettings.MaxReceiveMsgSize = defaultServerMaxReceiveMessageSize
	}

	grpcMiddlewares = append([]grpc.ServerOption{grpc.MaxRecvMsgSize(opts.GRPCSettings.MaxReceiveMsgSize)}, grpcMiddlewares...)

	if opts.GRPCSettings.MaxSendMsgSize > 0 {
		grpcMiddlewares = append([]grpc.ServerOption{grpc.MaxSendMsgSize(opts.GRPCSettings.MaxSendMsgSize)}, grpcMiddlewares...)
	}

	pluginOpts := asGRPCServeOpts(opts)
	pluginOpts.Debug = opts.debug
	pluginOpts.GRPCServer = func(opts []grpc.ServerOption) *grpc.Server {
		opts = append(opts, grpcMiddlewares...)
		return grpc.NewServer(opts...)
	}
	pluginOpts.ExitCallback = opts.ExitCallback
	pluginOpts.PluginID = opts.pluginID
	pluginOpts.PluginJson = pluginJson
	grpcplugin.Serve(pluginOpts)
}

func StandaloneServe(dsopts ServeOpts, address string) error {
	opts := asGRPCServeOpts(dsopts)

	if opts.GRPCServer == nil {
		opts.GRPCServer = plugin.DefaultGRPCServer
	}

	server := opts.GRPCServer(nil)

	plugKeys := []string{}

	if opts.ResourceServer != nil {
		pluginv2.RegisterResourceServer(server, opts.ResourceServer)
		plugKeys = append(plugKeys, "resources")
	}

	logger.DefaultLogger.Debug("Standalone plugin server", "capabilities", plugKeys)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	err = server.Serve(listener)
	if err != nil {
		return err
	}
	logger.DefaultLogger.Debug("Plugin server exited")

	return nil
}
