// dto包提供数据传输对象定义
package dto

// DslHistoryListReq2 DSL历史列表请求结构
type DslHistoryListReq2 struct {
	// 索引名
	IndexName string `json:"indexName"`
	// 开始时间与结束时间（格式："年-月-日 时:分:秒"）
	Date []string `json:"date"`
	// 当前页码
	Page int `json:"page"`
	// 每页条数
	Limit int `json:"limit"`
}
