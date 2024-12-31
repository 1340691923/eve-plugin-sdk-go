package build

import "fmt"

type PluginJsonData struct {
	Developer string `json:"developer"`
	Version         string          `json:"version"`
	MainGoFile      string          `json:"main_go_file"`
	PluginName      string          `json:"plugin_name"`
	BackendDebug    bool            `json:"backend_debug"`
	FrontendDebug   bool            `json:"frontend_debug"`
	PluginAlias     string          `json:"plugin_alias"`
	FrontendRoutes  []*Route        `json:"frontend_routes"`
	FrontendDevPort int             `json:"frontend_dev_port"`
	Frontend2c      bool            `json:"frontend_2c"`
	BackendRoutes   []*BackendRoute `json:"backend_routes"`
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
	Title   string `json:"title"`
	Icon    string `json:"icon"`
	Hidden  bool   `json:"hidden"`
	Service bool   `json:"service"`
}

func (this *PluginJsonData) String() string {
	return fmt.Sprintf("插件名：%s，开发者：%s,main文件：%s,插件别名：%s,版本号：%s", this.PluginName,this.Developer, this.MainGoFile, this.PluginAlias, this.Version)
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
	SqliteMigrateSqls []*ExecSql `json:"migrate_sqls"`
	MysqlMigrateSqls []*ExecSql `json:"mysql_migrate_sqls"`
	// 回退到当前版本需要的执行Sql
	SqliteRollback []*ExecSql `json:"rollback"`
	MysqlRollback []*ExecSql `json:"mysql_rollback"`
}

type PluginInitRespData struct {
	PluginJsonData *PluginJsonData `json:"plugin_json_data"`
	Gormigrate     *Gormigrate     `json:"gormigrate"`
}
