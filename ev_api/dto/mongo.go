// dto包提供数据传输对象定义
package dto

// 导入所需的包
import (
	// MongoDB BSON包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	// 时间处理包
	"time"
)

// MongoExecReq MongoDB执行请求结构
type MongoExecReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
	// 数据库名称
	DbName string `json:"db_name"`
	// MongoDB命令
	Command bson.D `json:"command"`
	// 超时时间
	Timeout time.Duration `json:"timeout"`
}

// ShowMongoDbsReq 显示MongoDB数据库请求结构
type ShowMongoDbsReq struct {
	// ES连接数据
	EsConnectData EsConnectData `json:"es_connect_data"`
}
