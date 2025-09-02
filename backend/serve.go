// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
	// GRPC插件包
	"github.com/1340691923/eve-plugin-sdk-go/backend/grpcplugin"
	// 日志包
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	// 构建包
	"github.com/1340691923/eve-plugin-sdk-go/build"
	// EVE API包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api"
	// Protobuf生成的包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	// 网络包
	"net"

	// GRPC中间件包
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	// GRPC Prometheus集成包
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	// HashiCorp插件系统包
	"github.com/hashicorp/go-plugin"
	// GRPC包
	"google.golang.org/grpc"
)

// 默认服务器最大接收消息大小（16MB）
const defaultServerMaxReceiveMessageSize = 1024 * 1024 * 16

// GRPCSettings GRPC服务设置结构
type GRPCSettings struct {
	// 最大接收消息大小
	MaxReceiveMsgSize int
	// 最大发送消息大小
	MaxSendMsgSize int
}

// ServeOpts 服务选项结构
type ServeOpts struct {
	// 插件ID
	pluginID string
	// 插件版本
	pluginVersion string
	// 插件JSON数据
	PluginJson *build.PluginJsonData
	// 健康检查处理器
	CheckHealthHandler CheckHealthHandler
	// 插件信息处理器
	PluginInfoHandler PluginInfoHandler
	// 资源调用处理器
	CallResourceHandler CallResourceHandler
	// 实时处理器
	LiveHandler LiveHandler
	// 任务处理器
	TaskHandler TaskHandler
	// GRPC设置
	GRPCSettings GRPCSettings
	// 调试模式标志
	debug bool
	// EVE RPC端口
	EvRpcPort string
	// 退出回调函数
	ExitCallback func()
	// 就绪回调函数
	ReadyCallback func(ctx context.Context)
}

// DefaultPluginInfoHandler 默认插件信息处理器结构
type DefaultPluginInfoHandler struct {
	// 插件ID
	PluginID string
	// 插件版本
	PluginVersion string
}

// NewDefaultPluginInfoHandler 创建新的默认插件信息处理器
// 参数：
//   - pluginID: 插件ID
//   - pluginVersion: 插件版本
//
// 返回：
//   - *DefaultPluginInfoHandler: 默认插件信息处理器实例
func NewDefaultPluginInfoHandler(pluginID string, pluginVersion string) *DefaultPluginInfoHandler {
	return &DefaultPluginInfoHandler{PluginID: pluginID, PluginVersion: pluginVersion}
}

// PluginInfo 获取插件信息
// 参数：
//   - ctx: 上下文
//   - req: 插件信息获取请求
//
// 返回：
//   - *PluginInfoGetRes: 插件信息响应
//   - error: 错误信息
func (d DefaultPluginInfoHandler) PluginInfo(ctx context.Context, req *PluginInfoGetReq) (*PluginInfoGetRes, error) {
	return &PluginInfoGetRes{
		PluginID:      d.PluginID,
		PluginVersion: d.PluginVersion,
	}, nil
}

// asGRPCServeOpts 将ServeOpts转换为GRPC服务选项
// 参数：
//   - opts: 服务选项
//
// 返回：
//   - grpcplugin.ServeOpts: GRPC插件服务选项
func asGRPCServeOpts(opts ServeOpts) grpcplugin.ServeOpts {

	if opts.PluginInfoHandler == nil {
		opts.PluginInfoHandler = NewDefaultPluginInfoHandler(opts.pluginID, opts.pluginID)
	}

	pluginOpts := grpcplugin.ServeOpts{
		PluginInfoServer: newPluginInfoSDKAdapter(opts.PluginInfoHandler, opts.CheckHealthHandler),
	}

	if opts.LiveHandler != nil {
		pluginOpts.LiveServer = newLiveSDKAdapter(opts.LiveHandler)
	}

	if opts.TaskHandler != nil {
		pluginOpts.TaskServer = newTaskSDKAdapter(opts.TaskHandler)
	}

	if opts.CallResourceHandler != nil {
		pluginOpts.ResourceServer = newResourceSDKAdapter(opts.CallResourceHandler)
	}

	return pluginOpts
}

// Serve 启动插件服务
// 参数：
//   - opts: 服务选项
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
	pluginOpts.ReadyCallback = opts.ReadyCallback

	pluginOpts.PluginID = opts.pluginID
	pluginOpts.PluginJson = pluginJson
	grpcplugin.Serve(pluginOpts)
}

// StandaloneServe 独立启动插件服务
// 参数：
//   - dsopts: 服务选项
//   - address: 服务地址
//
// 返回：
//   - error: 错误信息
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

	if opts.TaskServer != nil {
		pluginv2.RegisterTaskServer(server, opts.TaskServer)
		plugKeys = append(plugKeys, "task")
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
