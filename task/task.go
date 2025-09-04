package task

import (
	"context"
	"errors"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/goccy/go-json"
)

type TaskHandle interface {
	TaskExec(ctx context.Context, taskName string, userID int, jsonData json.RawMessage) (message string, err error)
}

// Live 实现实时数据处理功能
type Task struct {
	// 实时处理器
	taskHandle TaskHandle
}

// NewLive 创建一个新的实时数据处理实例
func NewTask(taskHandle TaskHandle) *Task {
	// 返回初始化的Live结构体
	return &Task{taskHandle: taskHandle}
}

// Pub2Channel 实现向频道发布消息的功能
func (this *Task) TaskExec(ctx context.Context, req *backend.TaskRequest) (*backend.TaskResponse, error) {
	// 检查实时处理器是否存在
	if this.taskHandle == nil {
		// 如果不存在，返回错误状态和错误信息
		return &backend.TaskResponse{Status: backend.TaskStatusFailed}, errors.New("该插件没有实现计划任务接口")
	}

	// 调用实时处理器发布消息
	message, err := this.taskHandle.TaskExec(ctx, req.TaskName, req.UserID, req.JSONData)
	// 检查是否有错误
	if err != nil {
		// 如果有错误，返回错误状态和错误信息
		return &backend.TaskResponse{Status: backend.TaskStatusFailed}, err
	}

	// 返回成功状态和JSON详情
	return &backend.TaskResponse{Status: backend.TaskStatusSuccess, Message: message}, nil
}
