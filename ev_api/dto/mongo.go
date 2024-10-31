package dto

import (
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	"time"
)

type MongoExecReq struct {
	EsConnectData EsConnectData `json:"es_connect_data"`
	DbName        string        `json:"db_name"`
	Command       bson.D        `json:"command"`
	Timeout       time.Duration `json:"timeout"`
}

type ShowMongoDbsReq struct {
	EsConnectData EsConnectData `json:"es_connect_data"`
}
