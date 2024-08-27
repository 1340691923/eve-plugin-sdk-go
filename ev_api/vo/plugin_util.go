package vo

type SelectRes struct {
	Result interface{} `json:"result"`
}

type ExecSqlRes struct {
	RowsAffected int64 `json:"rows_affected"`
}
