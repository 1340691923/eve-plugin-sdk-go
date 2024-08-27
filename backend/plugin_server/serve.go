package plugin_server

import (
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/1340691923/eve-plugin-sdk-go/check_health"
)

type ServeOpts struct {
	PluginJson *build.PluginJsonData
	backend.CallResourceHandler

	backend.PluginInfoHandler

	GRPCSettings backend.GRPCSettings

	EvRpcPort string

	EvRpcKey string

	Migration *build.Gormigrate

	ExitCallback func()
}

func Serve(opts ServeOpts) {
	backend.Serve(backend.ServeOpts{
		PluginJson:          opts.PluginJson,
		CallResourceHandler: opts.CallResourceHandler,
		PluginInfoHandler:   opts.PluginInfoHandler,
		CheckHealthHandler:  check_health.NewCheckHealthSvr(opts.PluginJson, opts.Migration),
		GRPCSettings:        opts.GRPCSettings,
		EvRpcPort:           opts.EvRpcPort,
		EvRpcKey:            opts.EvRpcKey,
		ExitCallback:        opts.ExitCallback,
	})

}
