package dto

type ExecSqlReq struct {
	PluginId string        `json:"plugin_id"`
	Sql      string        `json:"sql"`
	Args     []interface{} `json:"args"`
}

type SelectReq struct {
	PluginId string        `json:"plugin_id"`
	Sql      string        `json:"sql"`
	Args     []interface{} `json:"args"`
}
