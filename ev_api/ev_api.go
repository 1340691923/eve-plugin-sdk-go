// ev_api包提供EVE API的接口和实现
package ev_api

// 导入所需的包
import (
	"bytes"

	// 上下文包
	"context"
	"io"
	"net/http"

	// JSON编码包
	"encoding/json"
	// 格式化包
	"fmt"

	// 枚举包
	"github.com/1340691923/eve-plugin-sdk-go/enum"
	// MongoDB BSON包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	// 数据传输对象包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	// Protobuf协议包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	// 视图对象包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/vo"
	// 生成的Protobuf包
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	// 高性能JSON包
	json2 "github.com/goccy/go-json"
	// 错误处理包
	"github.com/pkg/errors"
	// 类型转换包
	"github.com/spf13/cast"
	// 高性能HTTP客户端
	// Protobuf编码包
	protobuf "google.golang.org/protobuf/proto"
	// URL处理包
	"net/url"
	// 路径处理包
	"path"
	// 并发控制包
	"sync"
	// 时间处理包
	"time"
)

// evApi 实现EVE API的结构体
type evApi struct {
	// RPC端口
	rpcPort string
	// 调试模式标志
	debug bool
	// 插件ID
	pluginId string
	client   *http.Client
}

// 全局变量
var (
	// 单例模式控制
	once *sync.Once
	// 全局API对象
	evApiObj *evApi
)

// init 初始化函数
func init() {
	// 初始化once
	once = new(sync.Once)
}

// SetEvApi 设置EVE API实例
// 参数：
//   - rpcPort: RPC服务端口
//   - pluginId: 插件ID
//   - debug: 调试模式标志
//
// 返回：
//   - *evApi: EVE API实例
func SetEvApi(rpcPort, pluginId string, debug bool) *evApi {
	// 创建HTTP客户端，设置超时

	// 使用sync.Once确保单例模式
	once.Do(func() {
		evApiObj = &evApi{
			rpcPort:  rpcPort,
			pluginId: pluginId,
			debug:    debug,
			client: &http.Client{
				Timeout: 300 * time.Second,
				Transport: &http.Transport{
					MaxIdleConns:        10000,
					MaxIdleConnsPerHost: 10000,
					IdleConnTimeout:     300 * time.Second,
				},
			},
		}
	})

	return evApiObj
}

// GetEvApi 获取全局EVE API实例
// 返回：
//   - *evApi: EVE API实例
func GetEvApi() *evApi {
	return evApiObj
}

// EsVersion 获取Elasticsearch版本
// 参数：
//   - ctx: 上下文
//   - req: ES连接数据
//
// 返回：
//   - version: ES版本号
//   - err: 错误信息
func (this *evApi) EsVersion(ctx context.Context, req dto.EsConnectData) (version int, err error) {
	// 定义返回结构
	res := vo.ApiCommonRes{}
	// 发送请求
	err = this.request(ctx, "api/plugin_util/EsVersion", req, &res)
	if err != nil {
		return 0, err
	}
	// 转换并返回版本号
	return cast.ToInt(res.Data), nil
}

// EsCatNodes 获取ES节点信息
// 参数：
//   - ctx: 上下文
//   - req: Cat节点请求
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCatNodes(ctx context.Context, req dto.CatNodesReq) (res *proto.Response, err error) {
	// 发送Protobuf请求
	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatNodes", req)
	if err != nil {
		return
	}
	return
}

// EsClusterStats 获取ES集群统计信息
// 参数：
//   - ctx: 上下文
//   - req: 集群统计请求
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsClusterStats(ctx context.Context, req dto.ClusterStatsReq) (res *proto.Response, err error) {
	// 发送Protobuf请求
	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsClusterStats", req)
	if err != nil {
		return
	}
	return
}

// EsPerformRequest 执行ES请求
// 参数：
//   - ctx: 上下文
//   - req: 执行请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsPerformRequest(ctx context.Context, req dto.PerformRequest) (res *proto.Response, err error) {
	// 发送Protobuf请求
	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsPerformRequest", req)
	if err != nil {
		return
	}
	return
}

// EsIndicesSegmentsRequest 获取ES索引分段信息
// 参数：
//   - ctx: 上下文
//   - req: 索引分段请求
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsIndicesSegmentsRequest(ctx context.Context, req dto.IndicesSegmentsRequest) (res *proto.Response, err error) {
	// 发送Protobuf请求
	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesSegmentsRequest", req)
	if err != nil {
		return
	}
	return
}

// Ping 测试ES连接
// 参数：
//   - ctx: 上下文
//   - req: Ping请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) Ping(ctx context.Context, req dto.PingReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/Ping", req)
	if err != nil {
		return
	}
	return
}

// EsRefresh 刷新ES索引
// 参数：
//   - ctx: 上下文
//   - req: 刷新请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsRefresh(ctx context.Context, req dto.RefreshReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsRefresh", req)
	if err != nil {
		return
	}
	return
}

// EsOpen 打开ES索引
// 参数：
//   - ctx: 上下文
//   - req: 打开索引请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsOpen(ctx context.Context, req dto.OpenReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsOpen", req)
	if err != nil {
		return
	}
	return
}

// EsFlush 执行ES索引刷新操作
// 参数：
//   - ctx: 上下文
//   - req: 刷新操作请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsFlush(ctx context.Context, req dto.FlushReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsFlush", req)
	if err != nil {
		return
	}
	return
}

// EsIndicesClearCache 清除ES索引缓存
// 参数：
//   - ctx: 上下文
//   - req: 清除缓存请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsIndicesClearCache(ctx context.Context, req dto.IndicesClearCacheReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesClearCache", req)
	if err != nil {
		return
	}
	return
}

// EsIndicesClose 关闭ES索引
// 参数：
//   - ctx: 上下文
//   - req: 关闭索引请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsIndicesClose(ctx context.Context, req dto.IndicesCloseReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesClose", req)
	if err != nil {
		return
	}
	return
}

// EsIndicesForcemerge 强制合并ES索引
// 参数：
//   - ctx: 上下文
//   - req: 强制合并请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsIndicesForcemerge(ctx context.Context, req dto.IndicesForcemergeReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesForcemerge", req)
	if err != nil {
		return
	}
	return
}

// EsDeleteByQuery 按查询条件删除ES文档
// 参数：
//   - ctx: 上下文
//   - req: 按查询删除请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsDeleteByQuery(ctx context.Context, req dto.DeleteByQueryReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsDeleteByQuery", req)
	if err != nil {
		return
	}
	return
}

// EsSnapshotCreate 创建ES快照
// 参数：
//   - ctx: 上下文
//   - req: 创建快照请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsSnapshotCreate(ctx context.Context, req dto.SnapshotCreateReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotCreate", req)
	if err != nil {
		return
	}
	return
}

// EsSnapshotDelete 删除ES快照
// 参数：
//   - ctx: 上下文
//   - req: 删除快照请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsSnapshotDelete(ctx context.Context, req dto.SnapshotDeleteReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotDelete", req)
	if err != nil {
		return
	}
	return
}

// EsRestoreSnapshot 恢复ES快照
// 参数：
//   - ctx: 上下文
//   - req: 恢复快照请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsRestoreSnapshot(ctx context.Context, req dto.RestoreSnapshotReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsRestoreSnapshot", req)
	if err != nil {
		return
	}
	return
}

// EsSnapshotStatus 获取ES快照状态
// 参数：
//   - ctx: 上下文
//   - req: 快照状态请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsSnapshotStatus(ctx context.Context, req dto.SnapshotStatusReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotStatus", req)
	if err != nil {
		return
	}
	return
}

// EsSnapshotGetRepository 获取ES快照仓库
// 参数：
//   - ctx: 上下文
//   - req: 获取快照仓库请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsSnapshotGetRepository(ctx context.Context, req dto.SnapshotGetRepositoryReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotGetRepository", req)
	if err != nil {
		return
	}
	return
}

// EsSnapshotCreateRepository 创建ES快照仓库
// 参数：
//   - ctx: 上下文
//   - req: 创建快照仓库请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsSnapshotCreateRepository(ctx context.Context, req dto.SnapshotCreateRepositoryReq) (res *proto.Response, err error) {
	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotCreateRepository", req)
	if err != nil {
		return
	}
	return
}

// EsSnapshotDeleteRepository 删除ES快照仓库
// 参数：
//   - ctx: 上下文
//   - req: 删除快照仓库请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsSnapshotDeleteRepository(ctx context.Context, req dto.SnapshotDeleteRepositoryReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotDeleteRepository", req)
	if err != nil {
		return
	}
	return
}

// EsGetIndices 获取ES索引列表
// 参数：
//   - ctx: 上下文
//   - req: 获取索引列表请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsGetIndices(ctx context.Context, req dto.GetIndicesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsGetIndices", req)
	if err != nil {
		return
	}
	return
}

// EsCatHealth 获取ES集群健康状态
// 参数：
//   - ctx: 上下文
//   - req: 集群健康状态请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCatHealth(ctx context.Context, req dto.CatHealthReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatHealth", req)
	if err != nil {
		return
	}
	return
}

// EsCatShards 获取ES分片信息
// 参数：
//   - ctx: 上下文
//   - req: 分片信息请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCatShards(ctx context.Context, req dto.CatShardsReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatShards", req)
	if err != nil {
		return
	}
	return
}

// EsCatCount 获取ES索引文档数量
// 参数：
//   - ctx: 上下文
//   - req: 文档数量请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCatCount(ctx context.Context, req dto.CatCountReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatCount", req)
	if err != nil {
		return
	}
	return
}

// EsCatAllocationRequest 获取ES分片分配信息
// 参数：
//   - ctx: 上下文
//   - req: 分片分配请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCatAllocationRequest(ctx context.Context, req dto.CatAllocationRequest) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatAllocationRequest", req)
	if err != nil {
		return
	}
	return
}

// EsCatAliases 获取ES别名信息
// 参数：
//   - ctx: 上下文
//   - req: 别名信息请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCatAliases(ctx context.Context, req dto.CatAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatAliases", req)
	if err != nil {
		return
	}
	return
}

// EsDelete 删除ES文档
// 参数：
//   - ctx: 上下文
//   - req: 删除文档请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsDelete(ctx context.Context, req dto.DeleteReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsDelete", req)
	if err != nil {
		return
	}
	return
}

// EsUpdate 更新ES文档
// 参数：
//   - ctx: 上下文
//   - req: 更新文档请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsUpdate(ctx context.Context, req dto.UpdateReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsUpdate", req)
	if err != nil {
		return
	}
	return
}

// EsCreate 创建ES文档
// 参数：
//   - ctx: 上下文
//   - req: 创建文档请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCreate(ctx context.Context, req dto.CreateReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCreate", req)
	if err != nil {
		return
	}
	return
}

// EsSearch 搜索ES文档
// 参数：
//   - ctx: 上下文
//   - req: 搜索请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsSearch(ctx context.Context, req dto.SearchReq) (res *proto.Response, err error) {

	return this.requestProtobuf(ctx, "api/plugin_util/EsSearch", req)
}

// EsIndicesPutSettingsRequest 设置ES索引配置
// 参数：
//   - ctx: 上下文
//   - req: 索引配置请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsIndicesPutSettingsRequest(ctx context.Context, req dto.IndicesPutSettingsRequest) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesPutSettingsRequest", req)
	if err != nil {
		return
	}
	return
}

// EsCreateIndex 创建ES索引
// 参数：
//   - ctx: 上下文
//   - req: 创建索引请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCreateIndex(ctx context.Context, req dto.CreateIndexReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCreateIndex", req)
	if err != nil {
		return
	}
	return
}

// EsDeleteIndex 删除ES索引
// 参数：
//   - ctx: 上下文
//   - req: 删除索引请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsDeleteIndex(ctx context.Context, req dto.DeleteIndexReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsDeleteIndex", req)
	if err != nil {
		return
	}
	return
}

// EsReindex 重新索引ES数据
// 参数：
//   - ctx: 上下文
//   - req: 重新索引请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsReindex(ctx context.Context, req dto.ReindexReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsReindex", req)
	if err != nil {
		return
	}
	return
}

// EsIndicesGetSettingsRequest 获取ES索引配置
// 参数：
//   - ctx: 上下文
//   - req: 获取索引配置请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsIndicesGetSettingsRequest(ctx context.Context, req dto.IndicesGetSettingsRequestReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesGetSettingsRequest", req)
	if err != nil {
		return
	}
	return
}

// EsPutMapping 设置ES索引映射
// 参数：
//   - ctx: 上下文
//   - req: 设置映射请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsPutMapping(ctx context.Context, req dto.PutMappingReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsPutMapping", req)
	if err != nil {
		return
	}
	return
}

// EsGetMapping 获取ES索引映射
// 参数：
//   - ctx: 上下文
//   - req: 获取映射请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsGetMapping(ctx context.Context, req dto.GetMappingReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsGetMapping", req)
	if err != nil {
		return
	}
	return
}

// EsGetAliases 获取ES别名
// 参数：
//   - ctx: 上下文
//   - req: 获取别名请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsGetAliases(ctx context.Context, req dto.GetAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsGetAliases", req)
	if err != nil {
		return
	}
	return
}

// EsAddAliases 添加ES别名
// 参数：
//   - ctx: 上下文
//   - req: 添加别名请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsAddAliases(ctx context.Context, req dto.AddAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsAddAliases", req)
	if err != nil {
		return
	}
	return
}

// EsRemoveAliases 移除ES别名
// 参数：
//   - ctx: 上下文
//   - req: 移除别名请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsRemoveAliases(ctx context.Context, req dto.RemoveAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsRemoveAliases", req)
	if err != nil {
		return
	}
	return
}

// EsMoveToAnotherIndexAliases 将别名移动到另一个索引
// 参数：
//   - ctx: 上下文
//   - req: 移动别名请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsMoveToAnotherIndexAliases(ctx context.Context, req dto.MoveToAnotherIndexAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsMoveToAnotherIndexAliases", req)
	if err != nil {
		return
	}
	return
}

// EsTaskList 获取ES任务列表
// 参数：
//   - ctx: 上下文
//   - req: 任务列表请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsTaskList(ctx context.Context, req dto.TaskListReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsTaskList", req)
	if err != nil {
		return
	}
	return
}

// EsTasksCancel 取消ES任务
// 参数：
//   - ctx: 上下文
//   - req: 取消任务请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsTasksCancel(ctx context.Context, req dto.TasksCancelReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsTasksCancel", req)
	if err != nil {
		return
	}
	return
}

// StoreExec 执行SQL语句
// 参数：
//   - ctx: 上下文
//   - sql: SQL语句
//   - args: SQL参数
//
// 返回：
//   - rowsAffected: 影响的行数
//   - err: 错误信息
func (this *evApi) StoreExec(ctx context.Context, sql string, args ...interface{}) (rowsAffected int64, err error) {
	data := &vo.ExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/ExecSql", &dto.ExecSqlReq{PluginId: this.pluginId, Sql: sql, Args: args}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

// StoreMoreExec 批量执行多条SQL语句
// 参数：
//   - ctx: 上下文
//   - sqls: SQL语句列表
//
// 返回：
//   - err: 错误信息
func (this *evApi) StoreMoreExec(ctx context.Context, sqls []dto.ExecSql) (err error) {
	err = this.request(ctx, "api/plugin_util/ExecMoreSql", &dto.ExecMoreReq{PluginId: this.pluginId, Sqls: sqls}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LiveBroadcastEvMsg2All 向所有用户广播EVE消息
// 参数：
//   - ctx: 上下文
//   - notice: 通知数据
//
// 返回：
//   - err: 错误信息
func (this *evApi) LiveBroadcastEvMsg2All(ctx context.Context, notice *dto.NoticeData) (err error) {
	err = notice.Validate()
	if err != nil {
		return errors.WithStack(err)
	}

	notice.PluginAlias = this.pluginId

	err = this.request(ctx, "api/plugin_util/LiveBroadcastEvMsg2All", &dto.LiveBroadcastEvMsg2AllReq{NoticeData: notice}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LiveBroadcastEvMsg2Roles 向指定角色广播EVE消息
// 参数：
//   - ctx: 上下文
//   - notice: 通知数据
//   - roleIds: 角色ID列表
//
// 返回：
//   - err: 错误信息
func (this *evApi) LiveBroadcastEvMsg2Roles(ctx context.Context, notice *dto.NoticeData, roleIds []int) (err error) {
	err = notice.Validate()
	if err != nil {
		return errors.WithStack(err)
	}

	notice.PluginAlias = this.pluginId

	err = this.request(ctx, "api/plugin_util/LiveBroadcastEvMsg2Roles", &dto.LiveBroadcastEvMsg2RolesReq{NoticeData: notice, RoleIds: roleIds}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LiveBroadcastEvMsg2Users 向指定用户广播EVE消息
// 参数：
//   - ctx: 上下文
//   - notice: 通知数据
//   - userIds: 用户ID列表
//
// 返回：
//   - err: 错误信息
func (this *evApi) LiveBroadcastEvMsg2Users(ctx context.Context, notice *dto.NoticeData, userIds []int) (err error) {
	err = notice.Validate()
	if err != nil {
		return errors.WithStack(err)
	}
	notice.PluginAlias = this.pluginId
	err = this.request(ctx, "api/plugin_util/LiveBroadcastEvMsg2Users", &dto.LiveBroadcastEvMsg2UsersReq{NoticeData: notice, UserIds: userIds}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// StoreSave 保存数据到数据库
// 参数：
//   - ctx: 上下文
//   - table: 表名
//   - data: 要保存的数据
//
// 返回：
//   - err: 错误信息
func (this *evApi) StoreSave(ctx context.Context, table string, data interface{}) (err error) {
	err = this.request(ctx, "api/plugin_util/SaveDb", &dto.SaveDb{PluginId: this.pluginId, TableName: table, Data: data}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// StoreUpdate 更新数据库记录
// 参数：
//   - ctx: 上下文
//   - table: 表名
//   - updateData: 更新的数据
//   - whereSql: WHERE条件SQL
//   - whereArgs: WHERE条件参数
//
// 返回：
//   - rowsAffected: 影响的行数
//   - err: 错误信息
func (this *evApi) StoreUpdate(ctx context.Context, table string, updateData map[string]interface{}, whereSql string, whereArgs ...interface{}) (rowsAffected int64, err error) {
	data := &vo.ExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/UpdateDb",
		&dto.UpdateDb{PluginId: this.pluginId, TableName: table, Data: updateData, UpdateArgs: whereArgs, UpdateSql: whereSql}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

// StoreDelete 删除数据库记录
// 参数：
//   - ctx: 上下文
//   - tableName: 表名
//   - whereSql: WHERE条件SQL
//   - whereArgs: WHERE条件参数
//
// 返回：
//   - rowsAffected: 影响的行数
//   - err: 错误信息
func (this *evApi) StoreDelete(ctx context.Context, tableName, whereSql string, whereArgs ...interface{}) (rowsAffected int64, err error) {
	data := &vo.ExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/DeleteDb",
		&dto.DeleteDb{PluginId: this.pluginId, TableName: tableName, WhereArgs: whereArgs, WhereSql: whereSql}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

// StoreInsertOrUpdate 插入或更新数据库记录
// 参数：
//   - ctx: 上下文
//   - table: 表名
//   - upsertData: 插入或更新的数据
//   - uniqueKeys: 唯一键字段
//
// 返回：
//   - err: 错误信息
func (this *evApi) StoreInsertOrUpdate(ctx context.Context, table string, upsertData map[string]interface{}, uniqueKeys ...string) (err error) {
	err = this.request(ctx, "api/plugin_util/InsertOrUpdateDb",
		&dto.InsertOrUpdateDb{PluginId: this.pluginId, TableName: table, UpsertData: upsertData, UniqueKeys: uniqueKeys}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LiveBroadcast 实时广播消息到指定频道
// 参数：
//   - ctx: 上下文
//   - channel: 频道名
//   - data: 广播数据
//
// 返回：
//   - noSub: 是否无订阅者
//   - err: 错误信息
func (this *evApi) LiveBroadcast(ctx context.Context, channel string, data interface{}) (noSub bool, err error) {
	err = this.request(ctx, "api/plugin_util/LiveBroadcast", map[string]interface{}{
		"channel": this.pluginId + "$v$" + channel,
		"data":    data,
	}, &vo.ApiCommonRes{})
	if err != nil {
		if err.Error() == NoSubscriberErr.Error() {
			return true, nil
		}
		return false, errors.WithStack(err)
	}
	return false, nil
}

// BatchLiveBroadcast 批量实时广播消息到指定频道
// 参数：
//   - ctx: 上下文
//   - channel: 频道名
//   - datas: 广播数据列表
//
// 返回：
//   - noSub: 是否无订阅者
//   - err: 错误信息
func (this *evApi) BatchLiveBroadcast(ctx context.Context, channel string, datas ...interface{}) (noSub bool, err error) {

	channel = this.pluginId + "$v$" + channel
	list := []interface{}{}

	for _, data := range datas {
		list = append(list, data)
	}

	err = this.request(ctx, "api/plugin_util/BatchLiveBroadcast", &dto.BatchLiveBroadcast{List: list}, &vo.ApiCommonRes{})
	if err != nil {
		if err.Error() == NoSubscriberErr.Error() {
			return true, nil
		}
		return false, errors.WithStack(err)
	}
	return false, nil
}

// GetEveToken 获取EVE令牌
// 参数：
//   - ctx: 上下文
//
// 返回：
//   - token: EVE令牌
//   - err: 错误信息
func (this *evApi) GetEveToken(ctx context.Context) (token string, err error) {

	res := &vo.ApiCommonRes{Data: ""}
	err = this.request(ctx, "api/plugin_util/GetEveToken", map[string]interface{}{}, res)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return cast.ToString(res.Data), nil
}

// StoreSelect 查询数据库记录，dest参数必须是一个切片
// 参数：
//   - ctx: 上下文
//   - dest: 查询结果目标对象
//   - sql: SQL查询语句
//   - args: SQL查询参数
//
// 返回：
//   - err: 错误信息
func (this *evApi) StoreSelect(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	data := &vo.SelectRes{}
	data.Result = &dest
	err = this.request(ctx, "api/plugin_util/SelectSql", &dto.SelectReq{Sql: sql, PluginId: this.pluginId, Args: args}, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// GetRoles4UserID 根据用户ID获取角色列表
// 参数：
//   - ctx: 上下文
//   - userId: 用户ID
//
// 返回：
//   - roleIds: 角色ID列表
//   - err: 错误信息
func (this *evApi) GetRoles4UserID(ctx context.Context, userId int) (roleIds []int, err error) {
	data := &vo.GetRoles4UserIdRes{}
	err = this.request(ctx, "api/plugin_util/GetRoles4UserID",
		&dto.GetRoles4UserIdReq{UserId: userId}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return data.RoleIds, err
	}
	return data.RoleIds, nil
}

// StoreFirst 查询数据库第一条记录
// 参数：
//   - ctx: 上下文
//   - dest: 查询结果目标对象
//   - sql: SQL查询语句
//   - args: SQL查询参数
//
// 返回：
//   - err: 错误信息
func (this *evApi) StoreFirst(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	data := &vo.SelectRes{}
	data.Result = &dest
	err = this.request(ctx, "api/plugin_util/FirstSql", &dto.SelectReq{Sql: sql, PluginId: this.pluginId, Args: args}, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// LoadDebugPlugin 加载调试插件
// 参数：
//   - ctx: 上下文
//   - req: 加载调试插件请求
//
// 返回：
//   - err: 错误信息
func (this *evApi) LoadDebugPlugin(ctx context.Context, req *dto.LoadDebugPlugin) (err error) {
	err = this.request(ctx, "api/plugin_util/LoadDebugPlugin", req, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// StopDebugPlugin 停止调试插件
// 参数：
//   - ctx: 上下文
//   - req: 停止调试插件请求
//
// 返回：
//   - err: 错误信息
func (this *evApi) StopDebugPlugin(ctx context.Context, req *dto.StopDebugPlugin) (err error) {
	err = this.request(ctx, "api/plugin_util/LoadDebugPlugin", req, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// EsRunDsl 执行ES DSL查询
// 参数：
//   - ctx: 上下文
//   - req: DSL查询请求
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsRunDsl(ctx context.Context, req *dto.PluginRunDsl) (res *proto.Response, err error) {

	if req.Params != nil {
		req.Path = fmt.Sprintf("%s?%s", req.Path, req.Params.Encode())
	}

	return this.requestProtobuf(ctx, "api/plugin_util/EsRunDsl", req)
}

// MysqlExecSql 执行MySQL SQL语句
// 参数：
//   - ctx: 上下文
//   - req: MySQL执行请求
//
// 返回：
//   - rowsAffected: 影响的行数
//   - err: 错误信息
func (this *evApi) MysqlExecSql(ctx context.Context, req *dto.MysqlExecReq) (rowsAffected int64, err error) {
	data := &vo.MysqlExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/MysqlExecSql", req, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

// MysqlSelectSql 查询MySQL数据，返回列名和结果
// 参数：
//   - ctx: 上下文
//   - req: MySQL查询请求
//
// 返回：
//   - columns: 列名列表
//   - result: 查询结果
//   - err: 错误信息
func (this *evApi) MysqlSelectSql(ctx context.Context, req *dto.MysqlSelectReq) (columns []string, result []map[string]interface{}, err error) {
	data := &vo.MysqlSelectSqlRes{}
	req.DbName = fmt.Sprintf("`%s`", req.DbName)
	err = this.request(ctx, "api/plugin_util/MysqlSelectSql", req, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return data.Columns, data.Result, nil
}

// MysqlFirstSql 查询MySQL第一条记录
// 参数：
//   - ctx: 上下文
//   - req: MySQL查询请求
//
// 返回：
//   - result: 查询结果
//   - err: 错误信息
func (this *evApi) MysqlFirstSql(ctx context.Context, req *dto.MysqlSelectReq) (result map[string]interface{}, err error) {
	data := &vo.MysqlFirstSqlRes{}
	req.DbName = fmt.Sprintf("`%s`", req.DbName)
	err = this.request(ctx, "api/plugin_util/MysqlFirstSql", req, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return data.Result, err
	}
	return data.Result, nil
}

func (this *evApi) MysqlDbs(ctx context.Context, req *dto.MysqlDbsReq) (dbs []string, err error) {
	data := &vo.MysqlDbsRes{}
	err = this.request(ctx, "api/plugin_util/MysqlFirstSql", req, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return data.Dbs, err
	}
	return data.Dbs, nil
}

func (this *evApi) DsType(ctx context.Context, req *dto.DsTypeReq) (dsType string, err error) {
	data := &vo.DsTypeRes{}
	err = this.request(ctx, "api/plugin_util/DsType", req, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return data.DsType, err
	}
	return data.DsType, nil
}

func (this *evApi) MysqlTables(ctx context.Context, req *dto.MysqlTablesReq) (tables []string, err error) {
	data := &vo.MysqlTablesRes{}
	req.DbName = fmt.Sprintf("`%s`", req.DbName)
	err = this.request(ctx, "api/plugin_util/MysqlFirstSql", req, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return data.Tables, err
	}
	return data.Tables, nil
}

// RedisExecCommand 执行Redis命令
// 参数：
//   - ctx: 上下文
//   - req: Redis执行请求
//
// 返回：
//   - data: 执行结果
//   - err: 错误信息
func (this *evApi) RedisExecCommand(ctx context.Context, req *dto.RedisExecReq) (data interface{}, err error) {

	result, err := this.requestProtobuf(ctx, "api/plugin_util/RedisExecCommand", req)
	if err != nil {
		return data, err
	}

	if result.StatusErr() != nil {
		return data, result.StatusErr()
	}

	res := map[string]interface{}{}

	err = json2.Unmarshal(result.ResByte(), &res)

	if err != nil {
		return data, err
	}

	return res["data"], nil
}

// ExecMongoCommand 执行MongoDB命令
// 参数：
//   - ctx: 上下文
//   - req: MongoDB执行请求
//
// 返回：
//   - data: 执行结果
//   - err: 错误信息
func (this *evApi) ExecMongoCommand(ctx context.Context, req *dto.MongoExecReq) (data bson.M, err error) {

	res, err := this.requestProtobuf(ctx, "api/plugin_util/MongoExecCommand", req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if res.StatusErr() != nil {
		return nil, res.StatusErr()
	}

	data = map[string]interface{}{}

	err = json2.Unmarshal(res.ResByte(), &data)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

// ShowMongoDbs 显示MongoDB数据库列表
// 参数：
//   - ctx: 上下文
//   - req: 显示数据库请求
//
// 返回：
//   - dbList: 数据库名称列表
//   - err: 错误信息
func (this *evApi) ShowMongoDbs(ctx context.Context, req *dto.ShowMongoDbsReq) (dbList []string, err error) {
	res := &vo.ApiCommonRes{Data: dbList}
	err = this.request(ctx, "api/plugin_util/ShowMongoDbs", req, res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cast.ToStringSlice(res.Data), nil
}

// request 发送HTTP请求的内部方法
// 参数：
//   - ctx: 上下文
//   - api: API路径
//   - requestData: 请求数据
//   - result: 响应结果
//   - nativeParse: 是否使用原生JSON解析
//
// 返回：
//   - error: 错误信息
func (this *evApi) request(ctx context.Context, api API, requestData interface{}, result interface{}, nativeParse ...bool) error {
	var requestDataJSON = []byte(`{}`)
	if requestData != nil {
		requestDataJSON, _ = json2.Marshal(requestData)
	}

	t1 := time.Now()
	res, err := this.SendRequest(ctx, api, "POST", requestDataJSON)
	if err != nil {
		return errors.WithStack(err)
	}
	if this.debug {
		_ = t1
		/*logger.DefaultLogger.Info("debug network",
		"api", api,
		"reqBody", string(requestDataJSON),
		"lose time", api, time.Now().Sub(t1).String())*/
	}
	if len(nativeParse) > 0 {
		err = json.Unmarshal(res, result)
	} else {
		err = json2.Unmarshal(res, result)
	}

	if err != nil {
		return errors.WithStack(err)
	}

	switch result.(type) {
	case *vo.ApiCommonRes:
		return result.(*vo.ApiCommonRes).Error()
	}

	return nil
}

// requestProtobuf 发送Protobuf格式HTTP请求的内部方法
// 参数：
//   - ctx: 上下文
//   - api: API路径
//   - requestData: 请求数据
//
// 返回：
//   - result: *proto.Response
//   - err: 错误信息
func (this *evApi) requestProtobuf(ctx context.Context, api API, requestData interface{}) (result *proto.Response, err error) {

	requestDataJSON, err := json2.Marshal(requestData)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	t1 := time.Now()
	res, err := this.SendRequest(ctx, api, "POST", requestDataJSON)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if this.debug {
		_ = t1
		/*logger.DefaultLogger.Info("debug network",
		"api", api,
		"reqBody", string(requestDataJSON),
		"lose time", api, time.Now().Sub(t1).String())*/
	}

	p := &pluginv2.CallResourceResponse{}

	err = protobuf.Unmarshal(res, p)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	headers := map[string][]string{}
	for k, values := range p.Headers {
		headers[k] = values.Values
	}

	result = proto.NewResponseWithProto(int(p.Code), headers, p.Body)

	//202为ev自定义code报错
	if p.Code == 202 {

		if _, ok := headers["EV-MSG"]; ok {
			if len(headers["EV-MSG"]) > 0 {
				err = errors.New(headers["EV-MSG"][0])
			}
		}

	}

	return result, err
}

func (this *evApi) SendRequest(ctx context.Context, api API, method string, requestDataJSON []byte) ([]byte, error) {
	// 创建请求对象
	url := fmt.Sprintf("http://127.0.0.1:%s/%s", this.rpcPort, api)
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(requestDataJSON))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(enum.EvFromPluginID, this.pluginId)
	// 发送请求
	resp, err := this.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("request failed: %w", err))
	}
	defer resp.Body.Close() // 确保关闭响应体

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to read response: %w", err))
	}

	// 检查HTTP状态码
	if resp.StatusCode >= 400 {
		return body, errors.WithStack(fmt.Errorf("HTTP error: %s", resp.Status))
	}

	return body, nil
}

// PluginRequestOptions 插件请求选项配置
type PluginRequestOptions struct {
	QueryParams url.Values        // URL 查询参数
	Headers     map[string]string // 自定义 Headers（包含 Content-Type）
	Timeout     time.Duration     // 请求超时时间（优先于 ctx 的超时）
	UserId      int               //当前操作者id
}

func (this *evApi) CallPlugin(
	ctx context.Context,
	pluginAlias string,
	api string,
	method string,
	body []byte,
	opts *PluginRequestOptions,
) ([]byte, error) {
	// 保护 opts 为非 nil
	if opts == nil {
		opts = &PluginRequestOptions{}
	}

	// 使用 path.Join 处理路径拼接
	fullPath := path.Join(pluginAlias, api)
	url := fmt.Sprintf("http://127.0.0.1:%s/api/plugin_util/CallPlugin/%s", this.rpcPort, fullPath)

	// 创建请求对象
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req = req.WithContext(ctx) // 关联上下文

	// 设置固定头
	req.Header.Set(enum.EvFromPluginID, this.pluginId)

	// 设置默认 Content-Type
	if _, ok := opts.Headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置用户ID
	if opts.UserId > 0 {
		req.Header.Set(enum.EvUserID, cast.ToString(opts.UserId))
	}

	// 设置自定义头
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	// 设置查询参数
	if len(opts.QueryParams) > 0 {
		q := req.URL.Query()
		for k, values := range opts.QueryParams {
			for _, v := range values {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	// 处理超时逻辑
	var resp *http.Response
	if opts.Timeout > 0 {
		// 创建带超时的子上下文
		timeoutCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
		defer cancel()

		resp, err = this.client.Do(req.WithContext(timeoutCtx))
	} else {
		resp, err = this.client.Do(req)
	}

	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("request failed: %w", err))
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to read response body: %w", err))
	}

	// 检查HTTP状态码
	if resp.StatusCode >= 400 {
		return respBody, errors.WithStack(fmt.Errorf("HTTP error: %s", resp.Status))
	}

	return respBody, nil
}
