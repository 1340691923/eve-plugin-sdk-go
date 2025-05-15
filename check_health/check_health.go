// check_health包提供插件健康检查功能
package check_health

// 导入所需的包
import (
	// 导入上下文包
	"context"
	// 导入格式化包
	"fmt"
	// 导入backend核心功能包
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	// 导入Web引擎包
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	// 导入构建工具包
	"github.com/1340691923/eve-plugin-sdk-go/build"
	// 导入JSON处理库
	"github.com/goccy/go-json"
)

// CheckHealthSvr 定义健康检查服务结构
type CheckHealthSvr struct {
	// 插件JSON数据
	pluginJson *build.PluginJsonData
	// 数据库迁移工具
	Migration *build.Gormigrate
	// Web引擎实例
	webEngine *web_engine.WebEngine
}

// NewCheckHealthSvr 创建一个新的健康检查服务实例
func NewCheckHealthSvr(pluginJson *build.PluginJsonData, migration *build.Gormigrate, webEngine *web_engine.WebEngine) *CheckHealthSvr {
	// 返回初始化的CheckHealthSvr结构体
	return &CheckHealthSvr{pluginJson: pluginJson, Migration: migration, webEngine: webEngine}
}

// CheckHealth 实现健康检查功能
func (c *CheckHealthSvr) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// 返回健康检查结果
	return &backend.CheckHealthResult{
		// 设置健康状态为正常
		Status: backend.HealthStatusOk,
		// 设置响应消息
		Message: "pong",
		// 添加详细的JSON信息
		JSONDetails: c.initJSONDetails(),
	}, nil
}

// initJSONDetails 初始化JSON详细信息
func (c *CheckHealthSvr) initJSONDetails() []byte {
	// 创建插件初始化响应数据
	var pluginInitRespData build.PluginInitRespData
	// 设置插件JSON数据
	pluginInitRespData.PluginJsonData = c.pluginJson
	// 设置数据库迁移工具
	pluginInitRespData.Gormigrate = c.Migration
	// 获取路由配置组
	group := c.webEngine.GetRouterConfigGroups()

	// 创建插件JSON后端路由映射，用于快速查找
	pluginJsonBackendRoutes := map[string]struct{}{}

	// 将现有后端路由添加到映射中
	for _, v := range pluginInitRespData.PluginJsonData.BackendRoutes {
		pluginJsonBackendRoutes[v.Path] = struct{}{}
	}

	// 处理路由配置组
	for _, v := range group {
		// 遍历每个路由配置
		for _, routerConfig := range v.RouterConfigs {
			// 如果路由已存在于插件JSON中，则跳过
			if _, ok := pluginJsonBackendRoutes[routerConfig.Url]; ok {
				continue
			}
			// 将新路由添加到插件JSON的后端路由列表中
			pluginInitRespData.PluginJsonData.BackendRoutes =
				append(pluginInitRespData.PluginJsonData.BackendRoutes, &build.BackendRoute{
					Path:     routerConfig.Url,
					Remark:   fmt.Sprintf("[%s]%s", v.GroupRemark, routerConfig.Remark),
					NeedAuth: routerConfig.NeedAuth,
				})
		}
	}
	// 将插件初始化响应数据序列化为JSON
	js, _ := json.Marshal(pluginInitRespData)
	// 返回JSON字节数组
	return js
}
