package plugin_server

import (
	"embed"
	"encoding/json"
	"flag"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/1340691923/eve-plugin-sdk-go/call_resource"
	"github.com/1340691923/eve-plugin-sdk-go/check_health"
	"log"
	"os"
)

type ServeOpts struct {
	PluginJsonBytes []byte

	pluginJson *build.PluginJsonData

	WebEngine *web_engine.WebEngine

	FrontendFiles embed.FS

	LiveHandler backend.LiveHandler

	GRPCSettings backend.GRPCSettings

	evRpcPort string

	Migration *build.Gormigrate

	ExitCallback func()
}

var (
	EvRpcPort        string
	Debug            bool
	TmpFileStorePath string
	PluginAlias      string
	DbType           string
)

const (
	MysqlDbTyp  = "mysql"
	SqliteDbTyp = "sqlite3"
)

func init() {
	flag.StringVar(&TmpFileStorePath, "tmpFileStorePath", "store_file_dir", "临时文件存放目录")
	flag.StringVar(&EvRpcPort, "evRpcPort", "8091", "ev基座内网访问端口")
	flag.BoolVar(&Debug, "debug", false, "是否开启调试")
	flag.StringVar(&DbType, "dbType", SqliteDbTyp, "存储类型")
}

func Serve(opts ServeOpts) {

	opts.evRpcPort = EvRpcPort

	if len(opts.PluginJsonBytes) > 0 && opts.pluginJson == nil {
		opts.pluginJson = new(build.PluginJsonData)
		err := json.Unmarshal(opts.PluginJsonBytes, &opts.pluginJson)
		if err != nil {
			log.Println("plugin.json解析失败")
			panic(err)
		}
	}
	if opts.pluginJson != nil {
		opts.pluginJson.BackendDebug = Debug
	}
	if opts.Migration == nil {
		opts.Migration = new(build.Gormigrate)
	}

	PluginAlias = opts.pluginJson.PluginAlias

	backend.Serve(backend.ServeOpts{
		PluginJson:          opts.pluginJson,
		CallResourceHandler: call_resource.NewResourceHandler(opts.WebEngine, opts.FrontendFiles),
		CheckHealthHandler:  check_health.NewCheckHealthSvr(opts.pluginJson, opts.Migration, opts.WebEngine),
		GRPCSettings:        opts.GRPCSettings,
		LiveHandler:         opts.LiveHandler,
		EvRpcPort:           opts.evRpcPort,
		ExitCallback:        opts.ExitCallback,
	})

}

func GetTmpFileStorePath() string {
	if _, err := os.Stat(TmpFileStorePath); os.IsNotExist(err) {
		os.MkdirAll(TmpFileStorePath, os.ModePerm)
	}
	return TmpFileStorePath
}
