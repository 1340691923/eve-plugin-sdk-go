package dto

type RedisExecReq struct {
	EsConnectData EsConnectData `json:"es_connect_data"`
	Args          []interface{} `json:"args"`
	DbName        int           `json:"dbName"`
}
