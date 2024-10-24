package plugin_server

import (
	"embed"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/1340691923/eve-plugin-sdk-go/call_resource"
	"github.com/1340691923/eve-plugin-sdk-go/check_health"
)

type ServeOpts struct {
	PluginJson *build.PluginJsonData

	WebEngine *web_engine.WebEngine

	FrontendFiles embed.FS

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
		CallResourceHandler: call_resource.NewResourceHandler(opts.WebEngine, opts.FrontendFiles),
		PluginInfoHandler:   opts.PluginInfoHandler,
		CheckHealthHandler:  check_health.NewCheckHealthSvr(opts.PluginJson, opts.Migration, opts.WebEngine),
		GRPCSettings:        opts.GRPCSettings,
		EvRpcPort:           opts.EvRpcPort,
		EvRpcKey:            opts.EvRpcKey,
		ExitCallback:        opts.ExitCallback,
	})

}
