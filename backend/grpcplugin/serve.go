package grpcplugin

import (
	"context"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"log"
	"os"
)

type ServeOpts struct {
	PluginID       string
	PluginJson     *build.PluginJsonData
	ResourceServer ResourceServer
	LiveServer     LiveServer

	PluginInfoServer PluginInfoServer

	GRPCServer func(options []grpc.ServerOption) *grpc.Server

	Debug bool

	ReadyCallback func()

	ExitCallback func()
}

func Serve(opts ServeOpts) {
	versionedPlugins := make(map[int]plugin.PluginSet)
	pSet := make(plugin.PluginSet)

	if opts.ResourceServer != nil {
		pSet["resource"] = &ResourceGRPCPlugin{
			ResourceServer: opts.ResourceServer,
		}
	}

	if opts.LiveServer != nil {
		pSet["live"] = &LiveGRPCPlugin{
			LiveServer: opts.LiveServer,
		}
	}

	if opts.PluginInfoServer != nil {
		pSet["basic"] = &PluginInfoGRPCPlugin{
			PluginInfoServer: opts.PluginInfoServer,
		}
	}

	versionedPlugins[ProtocolVersion] = pSet

	if opts.GRPCServer == nil {
		opts.GRPCServer = plugin.DefaultGRPCServer
	}

	plugKeys := make([]string, 0, len(pSet))
	for k := range pSet {
		plugKeys = append(plugKeys, k)
	}

	reattachConfig := make(chan *plugin.ReattachConfig)
	closeCh := make(chan struct{})
	exitCtx, cancelFn := context.WithCancel(context.Background())
	if opts.Debug {
		os.Setenv(MagicCookieKey, MagicCookieValue)
		go plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig:  handshake,
			GRPCServer:       opts.GRPCServer,
			VersionedPlugins: versionedPlugins,
			Plugins:          pSet,
			Test: &plugin.ServeTestConfig{
				Context:          exitCtx,
				ReattachConfigCh: reattachConfig,
				CloseCh:          closeCh,
				SyncStdio:        false,
			},
		})
		select {
		case pluginReattachConfig := <-reattachConfig:
			err := ev_api.GetEvApi().LoadDebugPlugin(context.Background(), &dto.LoadDebugPlugin{
				ID:   opts.PluginID,
				Addr: pluginReattachConfig.Addr.String(),
				Pid:  pluginReattachConfig.Pid,
			})
			if err != nil {
				panic(fmt.Sprintf("链接ev基座异常:%s", err.Error()))
			} else {
				if opts.ReadyCallback != nil {
					go func() {
						if r := recover(); r != nil {
							log.Println("ReadyCallback 发生 panic:", r)
						}
						opts.ReadyCallback()
					}()
				}
				log.Println(fmt.Sprintf("正常链接ev基座"))
			}
		}
		util.WaitQuit(func() {
			cancelFn()
		})
	} else {
		go func() {
			if opts.ReadyCallback != nil {
				go func() {
					if r := recover(); r != nil {
						log.Println("ReadyCallback 发生 panic:", r)
					}
					opts.ReadyCallback()
				}()
			}
			plugin.Serve(&plugin.ServeConfig{
				HandshakeConfig:  handshake,
				VersionedPlugins: versionedPlugins,
				GRPCServer:       opts.GRPCServer,
			})
			closeCh <- struct{}{}
		}()
	}

	<-closeCh

	if opts.Debug {
		err := ev_api.GetEvApi().StopDebugPlugin(context.Background(), &dto.StopDebugPlugin{
			ID: opts.PluginID,
		})
		if err != nil {
			logger.DefaultLogger.Debug(fmt.Sprintf("停止调试插件进程异常:%s", err.Error()))
		}
	}
	if opts.ExitCallback != nil {
		opts.ExitCallback()
	}
	logger.DefaultLogger.Debug("Plugin server exited")

	return
}
