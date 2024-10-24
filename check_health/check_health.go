package check_health

import (
	"context"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/goccy/go-json"
)

type CheckHealthSvr struct {
	pluginJson *build.PluginJsonData
	Migration  *build.Gormigrate
	webEngine  *web_engine.WebEngine
}

func NewCheckHealthSvr(pluginJson *build.PluginJsonData, migration *build.Gormigrate, webEngine *web_engine.WebEngine) *CheckHealthSvr {
	return &CheckHealthSvr{pluginJson: pluginJson, Migration: migration, webEngine: webEngine}
}

func (c *CheckHealthSvr) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:      backend.HealthStatusOk,
		Message:     "pong",
		JSONDetails: c.initJSONDetails(),
	}, nil
}

func (c *CheckHealthSvr) initJSONDetails() []byte {
	var pluginInitRespData build.PluginInitRespData
	pluginInitRespData.PluginJsonData = c.pluginJson
	pluginInitRespData.Gormigrate = c.Migration
	//pluginInitRespData.PluginJsonData.BackendRoutes = []*build.BackendRoute{}
	group := c.webEngine.GetRouterConfigGroups()

	pluginJsonBackendRoutes := map[string]struct{}{}

	for _, v := range pluginInitRespData.PluginJsonData.BackendRoutes {
		pluginJsonBackendRoutes[v.Path] = struct{}{}
	}

	for _, v := range group {
		for _, routerConfig := range v.RouterConfigs {
			if _, ok := pluginJsonBackendRoutes[routerConfig.Url]; ok {
				continue
			}
			pluginInitRespData.PluginJsonData.BackendRoutes =
				append(pluginInitRespData.PluginJsonData.BackendRoutes, &build.BackendRoute{
					Path:     routerConfig.Url,
					Remark:   fmt.Sprintf("[%s]%s", v.GroupRemark, routerConfig.Remark),
					NeedAuth: routerConfig.NeedAuth,
				})
		}
	}
	js, _ := json.Marshal(pluginInitRespData)
	return js
}
