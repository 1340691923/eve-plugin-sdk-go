// build包提供插件构建和打包相关功能
package build

// 导入格式化包
import "fmt"

// PluginJsonData 定义插件配置数据结构
type PluginJsonData struct {
	// 开发者信息
	Developer string `json:"developer"`
	// 插件版本
	Version string `json:"version"`
	// 主Go文件路径
	MainGoFile string `json:"main_go_file"`
	// 插件名称
	PluginName string `json:"plugin_name"`
	// 后端调试模式开关
	BackendDebug bool `json:"backend_debug"`
	// 前端调试模式开关
	FrontendDebug bool `json:"frontend_debug"`
	// 插件别名
	PluginAlias string `json:"plugin_alias"`
	// 前端路由配置
	FrontendRoutes []*Route `json:"frontend_routes"`
	// 前端开发服务端口
	FrontendDevPort int `json:"frontend_dev_port"`
	// 前端2c模式开关
	Frontend2c bool `json:"frontend_2c"`
	// 后端路由配置
	BackendRoutes []*BackendRoute `json:"backend_routes"`
}

// BackendRoute 定义后端路由结构
type BackendRoute struct {
	// 路由路径
	Path string `json:"path"`
	// 路由说明
	Remark string `json:"remark"`
	// 是否需要认证
	NeedAuth bool `json:"needAuth"`
}

// Route 定义前端路由结构
type Route struct {
	// 路由路径
	Path string `json:"path"`
	// 路由名称
	Name string `json:"name"`
	// 路由元数据
	Meta *RouteMeta `json:"meta"`
	// 子路由
	Children []*Route `json:"children"`
}

// RouteMeta 定义路由元数据结构
type RouteMeta struct {
	// 标题
	Title string `json:"title"`
	// 图标
	Icon string `json:"icon"`
	// 是否隐藏
	Hidden bool `json:"hidden"`
	// 是否为服务
	Service bool `json:"service"`
}

// String 返回插件数据的字符串表示
func (this *PluginJsonData) String() string {
	// 格式化插件基本信息
	return fmt.Sprintf("插件名：%s，开发者：%s,main文件：%s,插件别名：%s,版本号：%s", this.PluginName, this.Developer, this.MainGoFile, this.PluginAlias, this.Version)
}

// Gormigrate 定义数据库迁移结构
type Gormigrate struct {
	// 迁移项列表
	Migrations []*Migration `json:"migrations"`
}

// ExecSql 定义SQL执行结构
type ExecSql struct {
	// SQL语句
	Sql string `json:"sql"`
	// SQL参数
	Args []interface{} `json:"args"`
}

// Migration 定义数据库迁移项
type Migration struct {
	// 版本ID
	ID string `json:"id"`
	// SQLite迁移SQL语句列表
	SqliteMigrateSqls []*ExecSql `json:"migrate_sqls"`
	// MySQL迁移SQL语句列表
	MysqlMigrateSqls []*ExecSql `json:"mysql_migrate_sqls"`
	// SQLite回滚SQL语句列表
	SqliteRollback []*ExecSql `json:"rollback"`
	// MySQL回滚SQL语句列表
	MysqlRollback []*ExecSql `json:"mysql_rollback"`
}

// PluginInitRespData 定义插件初始化响应数据
type PluginInitRespData struct {
	// 插件JSON数据
	PluginJsonData *PluginJsonData `json:"plugin_json_data"`
	// 数据库迁移配置
	Gormigrate *Gormigrate `json:"gormigrate"`
}
