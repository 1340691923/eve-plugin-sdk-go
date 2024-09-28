package vo

type MysqlExecSqlRes struct {
	RowsAffected int64 `json:"rows_affected"`
}

type MysqlSelectSqlRes struct {
	Result []map[string]interface{} `json:"result"`
}

type MysqlFirstSqlRes struct {
	Result map[string]interface{} `json:"result"`
}
