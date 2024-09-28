package check_health

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/goccy/go-json"
)

type CheckHealthSvr struct {
	pluginJson *build.PluginJsonData
	Migration  *build.Gormigrate
}

func NewCheckHealthSvr(pluginJson *build.PluginJsonData, migration *build.Gormigrate) *CheckHealthSvr {
	return &CheckHealthSvr{pluginJson: pluginJson, Migration: migration}
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
	js, _ := json.Marshal(pluginInitRespData)
	return js
}
