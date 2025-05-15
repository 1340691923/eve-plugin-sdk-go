// ev_api包提供EVE API的接口和实现
package ev_api

// 导入所需的包
import (
	// 上下文包
	"context"
	// JSON编码包
	"encoding/json"
	// 格式化包
	"fmt"
	// 日志包
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
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
	"github.com/valyala/fasthttp"
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
	// HTTP客户端
	client *fasthttp.Client
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
	client := &fasthttp.Client{
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	// 使用sync.Once确保单例模式
	once.Do(func() {
		evApiObj = &evApi{
			rpcPort:  rpcPort,
			pluginId: pluginId,
			debug:    debug,
			client:   client,
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
//   - res: Protobuf响应
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
//   - res: Protobuf响应
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
//   - res: Protobuf响应
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
//   - res: Protobuf响应
//   - err: 错误信息
func (this *evApi) EsIndicesSegmentsRequest(ctx context.Context, req dto.IndicesSegmentsRequest) (res *proto.Response, err error) {
	// 发送Protobuf请求
	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesSegmentsRequest", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Ping(ctx context.Context, req dto.PingReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/Ping", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsRefresh(ctx context.Context, req dto.RefreshReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsRefresh", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsOpen(ctx context.Context, req dto.OpenReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsOpen", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsFlush(ctx context.Context, req dto.FlushReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsFlush", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsIndicesClearCache(ctx context.Context, req dto.IndicesClearCacheReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesClearCache", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsIndicesClose(ctx context.Context, req dto.IndicesCloseReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesClose", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsIndicesForcemerge(ctx context.Context, req dto.IndicesForcemergeReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesForcemerge", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsDeleteByQuery(ctx context.Context, req dto.DeleteByQueryReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsDeleteByQuery", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsSnapshotCreate(ctx context.Context, req dto.SnapshotCreateReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotCreate", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsSnapshotDelete(ctx context.Context, req dto.SnapshotDeleteReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotDelete", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsRestoreSnapshot(ctx context.Context, req dto.RestoreSnapshotReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsRestoreSnapshot", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsSnapshotStatus(ctx context.Context, req dto.SnapshotStatusReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotStatus", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsSnapshotGetRepository(ctx context.Context, req dto.SnapshotGetRepositoryReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotGetRepository", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsSnapshotCreateRepository(ctx context.Context, req dto.SnapshotCreateRepositoryReq) (res *proto.Response, err error) {
	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotCreateRepository", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsSnapshotDeleteRepository(ctx context.Context, req dto.SnapshotDeleteRepositoryReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsSnapshotDeleteRepository", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsGetIndices(ctx context.Context, req dto.GetIndicesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsGetIndices", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsCatHealth(ctx context.Context, req dto.CatHealthReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatHealth", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsCatShards(ctx context.Context, req dto.CatShardsReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatShards", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsCatCount(ctx context.Context, req dto.CatCountReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatCount", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsCatAllocationRequest(ctx context.Context, req dto.CatAllocationRequest) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatAllocationRequest", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsCatAliases(ctx context.Context, req dto.CatAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatAliases", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsDelete(ctx context.Context, req dto.DeleteReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsDelete", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsUpdate(ctx context.Context, req dto.UpdateReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsUpdate", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsCreate(ctx context.Context, req dto.CreateReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCreate", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsSearch(ctx context.Context, req dto.SearchReq) (res *proto.Response, err error) {

	return this.requestProtobuf(ctx, "api/plugin_util/EsSearch", req)
}

func (this *evApi) EsIndicesPutSettingsRequest(ctx context.Context, req dto.IndicesPutSettingsRequest) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesPutSettingsRequest", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsCreateIndex(ctx context.Context, req dto.CreateIndexReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCreateIndex", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsDeleteIndex(ctx context.Context, req dto.DeleteIndexReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsDeleteIndex", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsReindex(ctx context.Context, req dto.ReindexReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsReindex", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsIndicesGetSettingsRequest(ctx context.Context, req dto.IndicesGetSettingsRequestReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsIndicesGetSettingsRequest", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsPutMapping(ctx context.Context, req dto.PutMappingReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsPutMapping", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsGetMapping(ctx context.Context, req dto.GetMappingReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsGetMapping", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsGetAliases(ctx context.Context, req dto.GetAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsGetAliases", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsAddAliases(ctx context.Context, req dto.AddAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsAddAliases", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsRemoveAliases(ctx context.Context, req dto.RemoveAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsRemoveAliases", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsMoveToAnotherIndexAliases(ctx context.Context, req dto.MoveToAnotherIndexAliasesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsMoveToAnotherIndexAliases", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsTaskList(ctx context.Context, req dto.TaskListReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsTaskList", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsTasksCancel(ctx context.Context, req dto.TasksCancelReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsTasksCancel", req)
	if err != nil {
		return
	}
	return
}

// 执行sql
func (this *evApi) StoreExec(ctx context.Context, sql string, args ...interface{}) (rowsAffected int64, err error) {
	data := &vo.ExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/ExecSql", &dto.ExecSqlReq{PluginId: this.pluginId, Sql: sql, Args: args}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

// 执行sql
func (this *evApi) StoreMoreExec(ctx context.Context, sqls []dto.ExecSql) (err error) {
	err = this.request(ctx, "api/plugin_util/ExecMoreSql", &dto.ExecMoreReq{PluginId: this.pluginId, Sqls: sqls}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

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

func (this *evApi) StoreSave(ctx context.Context, table string, data interface{}) (err error) {
	err = this.request(ctx, "api/plugin_util/SaveDb", &dto.SaveDb{PluginId: this.pluginId, TableName: table, Data: data}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *evApi) StoreUpdate(ctx context.Context, table string, updateData map[string]interface{}, whereSql string, whereArgs ...interface{}) (rowsAffected int64, err error) {
	data := &vo.ExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/UpdateDb",
		&dto.UpdateDb{PluginId: this.pluginId, TableName: table, Data: updateData, UpdateArgs: whereArgs, UpdateSql: whereSql}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

func (this *evApi) StoreDelete(ctx context.Context, tableName, whereSql string, whereArgs ...interface{}) (rowsAffected int64, err error) {
	data := &vo.ExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/DeleteDb",
		&dto.DeleteDb{PluginId: this.pluginId, TableName: tableName, WhereArgs: whereArgs, WhereSql: whereSql}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

func (this *evApi) StoreInsertOrUpdate(ctx context.Context, table string, upsertData map[string]interface{}, uniqueKeys ...string) (err error) {
	err = this.request(ctx, "api/plugin_util/InsertOrUpdateDb",
		&dto.InsertOrUpdateDb{PluginId: this.pluginId, TableName: table, UpsertData: upsertData, UniqueKeys: uniqueKeys}, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

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

func (this *evApi) GetEveToken(ctx context.Context) (token string, err error) {

	res := &vo.ApiCommonRes{Data: ""}
	err = this.request(ctx, "api/plugin_util/GetEveToken", map[string]interface{}{}, res)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return cast.ToString(res.Data), nil
}

// 查询索引 dist参数必须是一个切片
func (this *evApi) StoreSelect(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	data := &vo.SelectRes{}
	data.Result = &dest
	err = this.request(ctx, "api/plugin_util/SelectSql", &dto.SelectReq{Sql: sql, PluginId: this.pluginId, Args: args}, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// 查询索引 dist参数必须是一个切片
func (this *evApi) GetRoles4UserID(ctx context.Context, userId int) (roleIds []int, err error) {
	data := &vo.GetRoles4UserIdRes{}
	err = this.request(ctx, "api/plugin_util/GetRoles4UserID",
		&dto.GetRoles4UserIdReq{UserId: userId}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return data.RoleIds, err
	}
	return data.RoleIds, nil
}

func (this *evApi) StoreFirst(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	data := &vo.SelectRes{}
	data.Result = &dest
	err = this.request(ctx, "api/plugin_util/FirstSql", &dto.SelectReq{Sql: sql, PluginId: this.pluginId, Args: args}, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *evApi) LoadDebugPlugin(ctx context.Context, req *dto.LoadDebugPlugin) (err error) {
	err = this.request(ctx, "api/plugin_util/LoadDebugPlugin", req, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *evApi) StopDebugPlugin(ctx context.Context, req *dto.StopDebugPlugin) (err error) {
	err = this.request(ctx, "api/plugin_util/LoadDebugPlugin", req, &vo.ApiCommonRes{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *evApi) EsRunDsl(ctx context.Context, req *dto.PluginRunDsl) (res *proto.Response, err error) {

	if req.Params != nil {
		req.Path = fmt.Sprintf("%s?%s", req.Path, req.Params.Encode())
	}

	return this.requestProtobuf(ctx, "api/plugin_util/EsRunDsl", req)
}

// MysqlExecSql MysqlSelectSql  MysqlFirstSql
// 执行sql
func (this *evApi) MysqlExecSql(ctx context.Context, req *dto.MysqlExecReq) (rowsAffected int64, err error) {
	data := &vo.MysqlExecSqlRes{}
	err = this.request(ctx, "api/plugin_util/MysqlExecSql", req, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return 0, err
	}
	return data.RowsAffected, nil
}

// 查询索引 dist参数必须是一个切片
func (this *evApi) MysqlSelectSql(ctx context.Context, req *dto.MysqlSelectReq) (columns []string, result []map[string]interface{}, err error) {
	data := &vo.MysqlSelectSqlRes{}
	err = this.request(ctx, "api/plugin_util/MysqlSelectSql", req, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return data.Columns, data.Result, nil
}

func (this *evApi) MysqlFirstSql(ctx context.Context, req *dto.MysqlSelectReq) (result map[string]interface{}, err error) {
	data := &vo.MysqlFirstSqlRes{}
	err = this.request(ctx, "api/plugin_util/MysqlFirstSql", req, &vo.ApiCommonRes{Data: data}, true)
	if err != nil {
		return data.Result, err
	}
	return data.Result, nil
}

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

func (this *evApi) ShowMongoDbs(ctx context.Context, req *dto.ShowMongoDbsReq) (dbList []string, err error) {
	res := &vo.ApiCommonRes{Data: dbList}
	err = this.request(ctx, "api/plugin_util/ShowMongoDbs", req, res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cast.ToStringSlice(res.Data), nil
}

func (this *evApi) request(ctx context.Context, api API, requestData interface{}, result interface{}, nativeParse ...bool) error {
	var requestDataJSON = []byte(`{}`)
	if requestData != nil {
		requestDataJSON, _ = json2.Marshal(requestData)
	}

	t1 := time.Now()
	res, err := this.SendRequest(ctx, api, fasthttp.MethodPost, requestDataJSON)
	if err != nil {
		return errors.WithStack(err)
	}
	if this.debug {
		logger.DefaultLogger.Info("debug network",
			"api", api,
			"reqBody", string(requestDataJSON),
			"lose time", api, time.Now().Sub(t1).String())
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

func (this *evApi) requestProtobuf(ctx context.Context, api API, requestData interface{}) (result *proto.Response, err error) {

	requestDataJSON, err := json2.Marshal(requestData)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	t1 := time.Now()
	res, err := this.SendRequest(ctx, api, fasthttp.MethodPost, requestDataJSON)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if this.debug {
		logger.DefaultLogger.Info("debug network",
			"api", api,
			"reqBody", string(requestDataJSON),
			"lose time", api, time.Now().Sub(t1).String())
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
	// 构建请求对象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 释放请求对象，防止内存泄漏

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 释放响应对象，防止内存泄漏

	// 设置请求 URL
	url := fmt.Sprintf("http://127.0.0.1:%s/%s", this.rpcPort, api)
	req.SetRequestURI(url)

	// 设置请求方法为 POST
	req.Header.SetMethod(method)
	req.Header.Set("Content-Type", "application/json")
	// 设置自定义头
	req.Header.Set(enum.EvFromPluginID, this.pluginId)

	// 设置请求体
	req.SetBody(requestDataJSON)

	// 发起请求

	errCh := make(chan error, 1)

	// 启动异步请求
	go func() {
		errCh <- this.client.Do(req, resp)
	}()

	select {
	case <-ctx.Done():
		// 防止 goroutine 泄漏
		go func() { <-errCh }()
		return nil, errors.WithStack(ctx.Err())

	case err := <-errCh:
		if err != nil {
			return nil, errors.WithStack(fmt.Errorf("request failed: %w", err))
		}
	}

	// 返回响应体
	return resp.Body(), nil
}

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

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 使用 path.Join 处理路径拼接，避免 // 或缺失 /
	fullPath := path.Join(pluginAlias, api)

	// 构建 URL + query 参数
	url := fmt.Sprintf("http://127.0.0.1:%s/api/plugin_util/CallPlugin/%s", this.rpcPort, fullPath)
	if len(opts.QueryParams) > 0 {
		url += "?" + opts.QueryParams.Encode()
	}
	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	// 固定头
	req.Header.Set(enum.EvFromPluginID, this.pluginId)

	if opts.Headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if opts.UserId > 0 {
		req.Header.Set(enum.EvUserID, cast.ToString(opts.UserId))
	}

	// 用户自定义头
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	// 设置 Body（不自动设置 Content-Type）
	if len(body) > 0 {
		req.SetBody(body)
	}

	// 判断是否设置超时时间
	timeout := opts.Timeout

	// 使用 channel 管理 goroutine 错误
	errCh := make(chan error, 1)

	go func() {
		var err error
		if timeout > 0 {
			err = this.client.DoTimeout(req, resp, timeout)
		} else {
			err = this.client.Do(req, resp)
		}
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		// 防止 goroutine 泄漏
		go func() { <-errCh }()
		return nil, errors.WithStack(ctx.Err())

	case err := <-errCh:
		if err != nil {
			return nil, errors.WithStack(fmt.Errorf("request failed: %w", err))
		}
	}

	return resp.Body(), nil
}
