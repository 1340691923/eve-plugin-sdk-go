// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 上下文包
	"context"
	// Protobuf生成的包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
)

// taskSDKAdapter 任务SDK适配器结构
type taskSDKAdapter struct {
	// 任务处理器
	taskHandler TaskHandler
}

// newTaskSDKAdapter 创建新的任务SDK适配器
// 参数：
//   - handler: 任务处理器
//
// 返回：
//   - *taskSDKAdapter: 任务SDK适配器实例
func newTaskSDKAdapter(handler TaskHandler) *taskSDKAdapter {
	return &taskSDKAdapter{
		taskHandler: handler,
	}
}

// taskResponseSenderFunc 任务响应发送函数类型
type taskResponseSenderFunc func(resp *TaskResponse) error

// Send 实现任务响应发送接口
// 参数：
//   - resp: 任务响应
//
// 返回：
//   - error: 错误信息
func (fn taskResponseSenderFunc) Send(resp *TaskResponse) error {
	return fn(resp)
}

// Exec 实现Protobuf任务执行接口
// 参数：
//   - ctx: 上下文
//   - protoReq: Protobuf任务请求
//
// 返回：
//   - *pluginv2.TaskResponse: Protobuf任务响应
//   - error: 错误信息
func (a *taskSDKAdapter) TaskExec(ctx context.Context, protoReq *pluginv2.TaskRequest) (*pluginv2.TaskResponse, error) {
	if a.taskHandler != nil {
		parsedReq := FromProto().TaskRequest(protoReq)

		res, err := a.taskHandler.TaskExec(ctx, parsedReq)
		if err != nil {
			return nil, err
		}
		return ToProto().TaskResponse(res), nil
	}

	return &pluginv2.TaskResponse{
		Status: pluginv2.TaskResponse_SUCCESS,
	}, nil
}
