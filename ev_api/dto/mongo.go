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

type GetMongoCollectionsReq struct {
	EsConnectData EsConnectData `json:"es_connect_data"`
	DbName        string        `json:"db_name"`
}

type FindMongoDocumentsReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
	Filter         bson.M        `json:"filter"`
	Projection     bson.M        `json:"projection"`
	Sort           bson.D        `json:"sort"`
	Skip           int64         `json:"skip"`
	Limit          int64         `json:"limit"`
}

type UpdateMongoDocumentReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
	DocId          interface{}   `json:"doc_id"`
	Update         bson.M        `json:"update"`
	Filter         bson.M        `json:"filter"`
}

type DeleteMongoDocumentReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
	DocId          interface{}   `json:"doc_id"`
	Filter         bson.M        `json:"filter"`
}

type InsertMongoDocumentReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
	Doc            bson.M        `json:"doc"`
}

type InsertManyMongoDocumentsReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
	Docs           []bson.M      `json:"docs"`
}

type DeleteManyMongoDocumentsReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DocIds         []interface{} `json:"doc_ids"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
}

type CountMongoDocumentsReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
	Filter         bson.M        `json:"filter"`
}

type AggregateMongoDocumentsReq struct {
	EsConnectData  EsConnectData `json:"es_connect_data"`
	DbName         string        `json:"db_name"`
	CollectionName string        `json:"collection_name"`
	Pipeline       bson.Pipeline `json:"pipeline"`
}
