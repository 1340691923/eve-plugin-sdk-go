package dto

type MysqlExecReq struct {
	EsConnectData EsConnectData `json:"es_connect_data"`
	Sql           string        `json:"sql"`
	Args          []interface{} `json:"args"`
	DbName        string        `json:"dbName"`
}

type MysqlSelectReq struct {
	EsConnectData EsConnectData `json:"es_connect_data"`
	Sql           string        `json:"sql"`
	Args          []interface{} `json:"args"`
	DbName        string        `json:"dbName"`
}
