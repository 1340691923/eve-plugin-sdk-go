// dto包提供数据传输对象定义
package dto

// RedisExecReq Redis执行请求结构
type RedisExecReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
	// 命令参数
	Args []interface{} `json:"args"`
	// 数据库索引
	DbName int `json:"dbName"`
}

type RedisBatchExecReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
	// 命令参数
	Args [][]interface{} `json:"args"`
	// 数据库索引
	DbName int `json:"dbName"`
}
