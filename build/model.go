package build

import "fmt"

type PluginJsonData struct {
	Version         string   `json:"version"`
	MainGoFile      string   `json:"main_go_file"`
	PluginName      string   `json:"plugin_name"`
	BackendDebug    bool     `json:"backend_debug"`
	FrontendDebug   bool     `json:"frontend_debug"`
	PluginAlias     string   `json:"plugin_alias"`
	FrontendRoutes  []*Route `json:"frontend_routes"`
	FrontendDevPort int      `json:"frontend_dev_port"`

	BackendRoutes []*BackendRoute `json:"backend_routes"`
}

type BackendRoute struct {
	Path     string `json:"path"`
	Remark   string `json:"remark"`
	NeedAuth bool   `json:"needAuth"`
}

type Route struct {
	Path     string     `json:"path"`
	Name     string     `json:"name"`
	Meta     *RouteMeta `json:"meta"`
	Children []*Route   `json:"children"`
}

type RouteMeta struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

func (this *PluginJsonData) String() string {
	return fmt.Sprintf("插件名：%s，main文件：%s,插件别名：%s,版本号：%s", this.PluginName, this.MainGoFile, this.PluginAlias, this.Version)
}

type Gormigrate struct {
	Migrations []*Migration `json:"migrations"`
}

type ExecSql struct {
	Sql  string        `json:"sql"`
	Args []interface{} `json:"args"`
}

type Migration struct {
	// 版本ID
	ID string `json:"id"`
	// 迁移到当前版本需要的执行Sql
	MigrateSqls []*ExecSql `json:"migrate_sqls"`
	// 回退到当前版本需要的执行Sql
	Rollback []*ExecSql `json:"rollback"`
}

type PluginInitRespData struct {
	PluginJsonData *PluginJsonData `json:"plugin_json_data"`
	Gormigrate     *Gormigrate     `json:"gormigrate"`
}
