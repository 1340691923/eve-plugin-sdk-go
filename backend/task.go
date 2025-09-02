// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
	// JSON包
	"encoding/json"
	// 数值转换包
	"strconv"
)

// TaskRequest 任务请求结构
type TaskRequest struct {
	// 请求数据
	JSONData json.RawMessage
	// 任务名称
	TaskName string
	// 用户ID
	UserID int
}

// TaskStatus 任务执行状态枚举
type TaskStatus int

// 定义任务状态常量
const (
	// 未知状态
	TaskStatusUnknown TaskStatus = iota
	// 成功状态
	TaskStatusSuccess
	// 失败状态
	TaskStatusFailed
	// 执行中状态
	TaskStatusRunning
)

// 任务状态名称映射
var taskStatusNames = map[int]string{
	0: "UNKNOWN",
	1: "SUCCESS",
	2: "FAILED",
	3: "RUNNING",
}

// TaskResponse 任务响应结构
type TaskResponse struct {
	// 执行状态
	Status TaskStatus
	// 任务输出
	Message string
}

// String 将任务状态转换为字符串
// 返回：
//   - string: 状态字符串表示
func (ts TaskStatus) String() string {
	s, exists := taskStatusNames[int(ts)]
	if exists {
		return s
	}
	return strconv.Itoa(int(ts))
}

// TaskHandler 任务处理器接口
type TaskHandler interface {
	// Exec 执行任务
	// 参数：
	//   - ctx: 上下文
	//   - req: 任务请求
	//
	// 返回：
	//   - *TaskResponse: 任务响应
	//   - error: 错误信息
	TaskExec(ctx context.Context, req *TaskRequest) (*TaskResponse, error)
}

// TaskHandlerFunc 任务处理函数类型
type TaskHandlerFunc func(ctx context.Context, req *TaskRequest) (*TaskResponse, error)

// Exec 实现TaskHandler接口
func (fn TaskHandlerFunc) TaskExec(ctx context.Context, req *TaskRequest) (*TaskResponse, error) {
	return fn(ctx, req)
}
