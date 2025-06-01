package plugin_server

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/1340691923/eve-plugin-sdk-go/call_resource"
	"github.com/1340691923/eve-plugin-sdk-go/check_health"
	"github.com/1340691923/eve-plugin-sdk-go/enum"
	"log"
	"os"
)

var PluginJson *build.PluginJsonData

type Assets struct {
	PluginJsonBytes []byte
	FrontendFiles   embed.FS
	Icon embed.FS
}

type ServeOpts struct {
	RegisterRoutes func(engine *web_engine.WebEngine)

	Assets *Assets

	LiveHandler backend.LiveHandler

	GRPCSettings backend.GRPCSettings

	Migration *build.Gormigrate

	ReadyCallBack func(ctx context.Context)

	ExitCallback func()
}

var (
	EvRpcPort        string
	Debug            bool
	TmpFileStorePath string
	PluginAlias      string
	DbType           string
)

func init() {
	flag.StringVar(&TmpFileStorePath, "tmpFileStorePath", "store_file_dir", "临时文件存放目录")
	flag.StringVar(&EvRpcPort, "evRpcPort", "8091", "ev基座内网访问端口")
	flag.BoolVar(&Debug, "debug", false, "是否开启调试")
	flag.StringVar(&DbType, "dbType", enum.SqliteDbTyp, "存储类型")
}

func Serve(opts ServeOpts) {

	evRpcPort := EvRpcPort
	pluginJson := new(build.PluginJsonData)
	if opts.Assets == nil {
		panic("静态资源不能为空")
	}

	if opts.Assets.PluginJsonBytes == nil {
		panic("插件配置资源不能为空")
	}

	if len(opts.Assets.PluginJsonBytes) > 0 {
		err := json.Unmarshal(opts.Assets.PluginJsonBytes, &pluginJson)
		if err != nil {
			log.Println("plugin.json解析失败")
			panic(err)
		}
		PluginJson = pluginJson
	}
	if pluginJson != nil {
		pluginJson.BackendDebug = Debug
	}
	if opts.Migration == nil {
		opts.Migration = new(build.Gormigrate)
	}

	PluginAlias = pluginJson.PluginAlias

	webEngine := web_engine.NewWebEngine()

	if opts.RegisterRoutes != nil {
		opts.RegisterRoutes(webEngine)
	}

	backend.Serve(backend.ServeOpts{
		PluginJson:          pluginJson,
		CallResourceHandler: call_resource.NewResourceHandler(webEngine, opts.Assets.FrontendFiles,opts.Assets.Icon),
		CheckHealthHandler:  check_health.NewCheckHealthSvr(pluginJson, opts.Migration, webEngine),
		GRPCSettings:        opts.GRPCSettings,
		LiveHandler:         opts.LiveHandler,
		EvRpcPort:           evRpcPort,
		ExitCallback:        opts.ExitCallback,
		ReadyCallback:       opts.ReadyCallBack,
	})

}

func GetTmpFileStorePath() string {
	if _, err := os.Stat(TmpFileStorePath); os.IsNotExist(err) {
		os.MkdirAll(TmpFileStorePath, os.ModePerm)
	}
	return TmpFileStorePath
}
