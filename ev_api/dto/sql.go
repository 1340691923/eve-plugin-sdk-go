// dto包提供数据传输对象定义
package dto

// ExecSqlReq SQL执行请求结构
type ExecSqlReq struct {
	// 插件ID
	PluginId string `json:"plugin_id"`
	// SQL语句
	Sql string `json:"sql"`
	// SQL参数
	Args []interface{} `json:"args"`
}

// SelectReq SQL查询请求结构
type SelectReq struct {
	// 插件ID
	PluginId string `json:"plugin_id"`
	// SQL语句
	Sql string `json:"sql"`
	// SQL参数
	Args []interface{} `json:"args"`
}

// ExecMoreReq 批量SQL执行请求结构
type ExecMoreReq struct {
	// 插件ID
	PluginId string `json:"plugin_id"`
	// SQL语句列表
	Sqls []ExecSql `json:"sqls"`
}

// ExecSql SQL执行结构
type ExecSql struct {
	// SQL语句
	Sql string `json:"sql"`
	// SQL参数
	Args []interface{} `json:"args"`
}

// SaveDb 保存数据到数据库请求结构
type SaveDb struct {
	// 插件ID
	PluginId string `json:"plugin_id"`
	// 目标表名
	TableName string `json:"table"` // 目标表名
	// 要插入或更新的数据
	Data interface{} `json:"data"` // 要插入或更新的数据
}

// UpdateDb 更新数据库请求结构
type UpdateDb struct {
	// 插件ID
	PluginId string `json:"plugin_id"`
	// 目标表名
	TableName string `json:"table"` // 目标表名
	// 更新SQL语句
	UpdateSql string `json:"update_sql"`
	// 更新SQL参数
	UpdateArgs []interface{} `json:"update_args"`
	// 要插入或更新的数据
	Data map[string]interface{} `json:"data"` // 要插入或更新的数据
}

// InsertOrUpdateDb 插入或更新数据库请求结构
type InsertOrUpdateDb struct {
	// 插件ID
	PluginId string `json:"plugin_id"`
	// 目标表名
	TableName string `json:"table"` // 目标表名
	// 没有则新增，有则更新的数据
	UpsertData map[string]interface{} // 没有则新增，有则更新
	// 冲突检查的唯一键
	UniqueKeys []string // 冲突检查的唯一键
}

// DeleteDb 删除数据库记录请求结构
type DeleteDb struct {
	// 插件ID
	PluginId string `json:"plugin_id"`
	// 目标表名
	TableName string `json:"table"` // 目标表名
	// WHERE条件SQL
	WhereSql string `json:"where_sql"`
	// WHERE条件参数
	WhereArgs []interface{} `json:"where_args"`
}
