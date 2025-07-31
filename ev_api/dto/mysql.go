// dto包提供数据传输对象定义
package dto

// MysqlExecReq MySQL执行请求结构
type MysqlExecReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
	// SQL语句
	Sql string `json:"sql"`
	// SQL参数
	Args []interface{} `json:"args"`
	// 数据库名称
	DbName string `json:"dbName"`
}

// MysqlSelectReq MySQL查询请求结构
type MysqlSelectReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
	// SQL语句
	Sql string `json:"sql"`
	// SQL参数
	Args []interface{} `json:"args"`
	// 数据库名称
	DbName string `json:"dbName"`
}

type MysqlDbsReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
}

type MysqlTablesReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
	// 数据库名称
	DbName string `json:"dbName"`
}

type BatchInsertDataReq struct {
	EsConnectData EsConnectData `json:"es_connect_data"`
	DbName string `json:"db_name"`
	TableName string `json:"table_name"`
	Cols []string `json:"cols"`
	Data [][]interface{} `json:"data"`
}



type DsTypeReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
}
