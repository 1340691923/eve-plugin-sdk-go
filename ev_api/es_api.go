// ev_api包提供EVE API的接口和实现
package ev_api

// 导入所需的包
import (
	// 上下文包
	"context"
	// 日志包
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	// MongoDB BSON包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	// 数据传输对象包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	// Protobuf协议包
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	// 错误处理包
	"github.com/pkg/errors"
	// IO操作包
	"io/ioutil"
	// 日志包
	"log"
	// HTTP包
	"net/http"
	// 时间处理包
	"time"
)

// EvApiAdapter ES API适配器结构体
type EvApiAdapter struct {
	// 连接ID
	ConnId int
	// 用户ID
	UserId int
}

// MongoAggregateDocuments 执行MongoDB聚合查询
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - collectionName: 集合名称
//   - pipeline: 聚合管道
//
// 返回：
//   - []bson.M: 聚合查询结果
//   - error: 错误信息
func (this *EvApiAdapter) MongoAggregateDocuments(ctx context.Context, dbName, collectionName string, pipeline bson.Pipeline) ([]bson.M, error) {
	return GetEvApi().AggregateMongoDocuments(ctx, &dto.AggregateMongoDocumentsReq{
		EsConnectData:  this.buildEsConnectData(),
		DbName:         dbName,
		CollectionName: collectionName,
		Pipeline:       pipeline,
	})
}

// MongoCountDocuments 统计MongoDB文档数量
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - collectionName: 集合名称
//   - filter: 查询过滤条件
//
// 返回：
//   - int64: 文档数量
//   - error: 错误信息
func (this *EvApiAdapter) MongoCountDocuments(ctx context.Context, dbName, collectionName string, filter bson.M) (int64, error) {
	return GetEvApi().CountMongoDocuments(ctx, &dto.CountMongoDocumentsReq{
		EsConnectData:  this.buildEsConnectData(),
		DbName:         dbName,
		CollectionName: collectionName,
		Filter:         filter,
	})
}

// MongoInsertManyDocuments 批量插入MongoDB文档
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - collectionName: 集合名称
//   - docs: 要插入的文档列表
//
// 返回：
//   - []interface{}: 插入后的文档ID列表
//   - error: 错误信息
func (this *EvApiAdapter) MongoInsertManyDocuments(ctx context.Context, dbName, collectionName string, docs []bson.M) (insertIds []string, err error) {
	return GetEvApi().InsertManyMongoDocuments(ctx, &dto.InsertManyMongoDocumentsReq{
		EsConnectData:  this.buildEsConnectData(),
		DbName:         dbName,
		CollectionName: collectionName,
		Docs:           docs,
	})
}

// MongoInsertDocument 插入单个MongoDB文档
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - collectionName: 集合名称
//   - doc: 要插入的文档
//
// 返回：
//   - interface{}: 插入后的文档ID
//   - error: 错误信息
func (this *EvApiAdapter) MongoInsertDocument(ctx context.Context, dbName, collectionName string, doc bson.M) (insertId string, err error) {
	result, err := GetEvApi().InsertMongoDocument(ctx, &dto.InsertMongoDocumentReq{
		EsConnectData:  this.buildEsConnectData(),
		DbName:         dbName,
		CollectionName: collectionName,
		Doc:            doc,
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

// MongoDeleteDocument 删除单个MongoDB文档
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - collectionName: 集合名称
//   - docId: 要删除的文档ID
//
// 返回：
//   - int64: 删除的文档数量
//   - error: 错误信息
func (this *EvApiAdapter) MongoDeleteDocument(ctx context.Context, dbName, collectionName string, docId interface{}, filter bson.M) (deleteCnt int64, err error) {
	return GetEvApi().DeleteMongoDocument(ctx, &dto.DeleteMongoDocumentReq{
		EsConnectData:  this.buildEsConnectData(),
		DbName:         dbName,
		CollectionName: collectionName,
		DocId:          docId,
		Filter:         filter,
	})
}

// MongoUpdateDocument 更新MongoDB文档
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - collectionName: 集合名称
//   - docId: 要更新的文档ID
//   - update: 更新操作
//
// 返回：
//   - matchedCount: 匹配的文档数量
//   - modifiedCount: 修改的文档数量
//   - upsertedCount: 插入的文档数量
//   - upsertedID: 插入的文档ID
//   - error: 错误信息
func (this *EvApiAdapter) MongoUpdateDocument(ctx context.Context, dbName, collectionName string, docId interface{}, filter bson.M, update bson.M) (matchedCount int64, modifiedCount int64, upsertedCount int64, upsertedID interface{}, err error) {
	result, err := GetEvApi().UpdateMongoDocument(ctx, &dto.UpdateMongoDocumentReq{
		EsConnectData:  this.buildEsConnectData(),
		DbName:         dbName,
		CollectionName: collectionName,
		DocId:          docId,
		Update:         update,
		Filter:         filter,
	})
	if err != nil {
		return 0, 0, 0, nil, err
	}
	return result.MatchedCount, result.ModifiedCount, result.UpsertedCount, result.UpsertedID, nil
}

// MongoFindDocuments 查找MongoDB文档
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - collectionName: 集合名称
//   - filter: 查询过滤条件
//   - sort: 排序条件
//   - skip: 跳过的文档数量
//   - limit: 返回的文档数量限制
//
// 返回：
//   - []bson.M: 查询结果文档列表
//   - error: 错误信息
func (this *EvApiAdapter) MongoFindDocuments(ctx context.Context, dbName, collectionName string, projection bson.M, filter bson.M, sort bson.D, skip int64, limit int64) ([]bson.M, error) {
	return GetEvApi().FindMongoDocuments(ctx, &dto.FindMongoDocumentsReq{
		EsConnectData:  this.buildEsConnectData(),
		DbName:         dbName,
		CollectionName: collectionName,
		Filter:         filter,
		Projection:     projection,
		Sort:           sort,
		Skip:           skip,
		Limit:          limit,
	})
}

// MongoGetCollections 获取MongoDB数据库中的集合列表
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//
// 返回：
//   - []string: 集合名称列表
//   - error: 错误信息
func (this *EvApiAdapter) MongoGetCollections(ctx context.Context, dbName string) ([]string, error) {
	return GetEvApi().GetMongoCollections(ctx, &dto.GetMongoCollectionsReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
	})
}
// ShowMongoDbs 显示MongoDB数据库列表
// 参数：
//   - ctx: 上下文
//
// 返回：
//   - []string: 数据库名称列表
//   - error: 错误信息
func (this *EvApiAdapter) ShowMongoDbs(ctx context.Context) ([]string, error) {
	return GetEvApi().ShowMongoDbs(ctx, &dto.ShowMongoDbsReq{EsConnectData: this.buildEsConnectData()})
}
// DsType 获取数据源类型
// 返回：
//   - string: 数据源类型字符串
func (this *EvApiAdapter) DsType() string {
	dsType, err := GetEvApi().DsType(context.Background(), &dto.DsTypeReq{
		EsConnectData: this.buildEsConnectData(),
	})
	if err != nil {
		logger.DefaultLogger.Error("get ds type err", err)
		return ""
	}
	return dsType
}

// NewEvWrapApi 创建一个新的ES API适配器
// 参数：
//   - connId: 连接ID
//   - userId: 用户ID
//
// 返回：
//   - *EvApiAdapter: ES API适配器实例
func NewEvWrapApi(connId int, userId int) *EvApiAdapter {
	return &EvApiAdapter{ConnId: connId, UserId: userId}
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
func (this *EvApiAdapter) StoreExec(ctx context.Context, sql string, args ...interface{}) (rowsAffected int64, err error) {
	return GetEvApi().StoreExec(ctx, sql, args...)
}

// StoreMoreExec 批量执行SQL语句，带事务
// 参数：
//   - ctx: 上下文
//   - sqls: SQL语句列表
//
// 返回：
//   - err: 错误信息
func (this *EvApiAdapter) StoreMoreExec(ctx context.Context, sqls []dto.ExecSql) (err error) {
	return GetEvApi().StoreMoreExec(ctx, sqls)
}

// LiveBroadcast 通过长连接广播消息（给每个订阅该频道的用户）
// 参数：
//   - ctx: 上下文
//   - channel: 频道名称
//   - data: 广播数据
//
// 返回：
//   - isNoSub: 是否没有订阅者
//   - err: 错误信息
func (this *EvApiAdapter) LiveBroadcast(ctx context.Context, channel string, data interface{}) (isNoSub bool, err error) {
	return GetEvApi().LiveBroadcast(ctx, channel, data)
}

// BatchLiveBroadcast 批量广播消息
// 参数：
//   - ctx: 上下文
//   - channel: 频道名称
//   - datas: 广播数据列表
//
// 返回：
//   - noSub: 是否没有订阅者
//   - err: 错误信息
func (this *EvApiAdapter) BatchLiveBroadcast(ctx context.Context, channel string, datas ...interface{}) (noSub bool, err error) {
	return GetEvApi().BatchLiveBroadcast(ctx, channel, datas)
}

// StoreSelect 执行查询SQL，结果存入dest切片
// 参数：
//   - ctx: 上下文
//   - dest: 结果接收对象（必须是切片）
//   - sql: SQL语句
//   - args: SQL参数
//
// 返回：
//   - err: 错误信息
func (this *EvApiAdapter) StoreSelect(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	return GetEvApi().StoreSelect(ctx, dest, sql, args...)
}

// GetRoles4UserID 获取用户角色ID列表
// 参数：
//   - ctx: 上下文
//   - userId: 用户ID
//
// 返回：
//   - roleIds: 角色ID列表
//   - err: 错误信息
func (this *EvApiAdapter) GetRoles4UserID(ctx context.Context, userId int) (roleIds []int, err error) {
	return GetEvApi().GetRoles4UserID(ctx, userId)
}

// StoreFirst 执行查询SQL，获取第一条结果
// 参数：
//   - ctx: 上下文
//   - dest: 结果接收对象
//   - sql: SQL语句
//   - args: SQL参数
//
// 返回：
//   - err: 错误信息
func (this *EvApiAdapter) StoreFirst(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	return GetEvApi().StoreFirst(ctx, dest, sql, args...)
}

// LoadDebugPlugin 加载调试插件
// 参数：
//   - ctx: 上下文
//   - req: 加载调试插件请求
//
// 返回：
//   - err: 错误信息
func (this *EvApiAdapter) LoadDebugPlugin(ctx context.Context, req *dto.LoadDebugPlugin) (err error) {
	return GetEvApi().LoadDebugPlugin(ctx, req)
}

// StopDebugPlugin 停止调试插件
// 参数：
//   - ctx: 上下文
//   - req: 停止调试插件请求
//
// 返回：
//   - err: 错误信息
func (this *EvApiAdapter) StopDebugPlugin(ctx context.Context, req *dto.StopDebugPlugin) (err error) {
	return GetEvApi().StopDebugPlugin(ctx, req)
}

// EsRunDsl 执行ES DSL查询
// 参数：
//   - ctx: 上下文
//   - req: 插件运行DSL请求
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsRunDsl(ctx context.Context, req *dto.PluginRunDsl2) (res *proto.Response, err error) {

	return GetEvApi().EsRunDsl(ctx, &dto.PluginRunDsl{
		EsConnectData: &dto.EsConnectData{
			UserID:    this.UserId,
			EsConnect: this.ConnId,
		},
		Params:     req.Params,
		HttpMethod: req.HttpMethod,
		Path:       req.Path,
		Dsl:        req.Dsl,
	})
}

// EsVersion 获取ES版本
// 返回：
//   - version: ES版本号
//   - err: 错误信息
func (this *EvApiAdapter) EsVersion() (version int, err error) {
	verson, err := GetEvApi().EsVersion(context.Background(), this.buildEsConnectData())
	if err != nil {
		logger.DefaultLogger.Error("get es version err", err)
		return 0, err
	}
	return verson, nil
}

// MysqlExecSql 执行MySQL SQL语句（INSERT、UPDATE、DELETE）
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - sql: SQL语句
//   - args: SQL参数
//
// 返回：
//   - rowsAffected: 影响的行数
//   - err: 错误信息
func (this *EvApiAdapter) MysqlExecSql(ctx context.Context, dbName, sql string, args ...interface{}) (rowsAffected int64, err error) {
	return GetEvApi().MysqlExecSql(ctx, &dto.MysqlExecReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Sql:           sql,
		Args:          args,
	})
}

// MysqlSelectSql 执行MySQL查询语句
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - sql: SQL语句
//   - args: SQL参数
//
// 返回：
//   - columns: 列名列表
//   - res: 查询结果
//   - err: 错误信息
func (this *EvApiAdapter) MysqlSelectSql(ctx context.Context, dbName, sql string, args ...interface{}) (columns []string, res []map[string]interface{}, err error) {
	return GetEvApi().MysqlSelectSql(ctx, &dto.MysqlSelectReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Sql:           sql,
		Args:          args,
	})
}

// MysqlFirstSql 执行MySQL查询并获取第一条记录
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - sql: SQL语句
//   - args: SQL参数
//
// 返回：
//   - res: 查询结果
//   - err: 错误信息
func (this *EvApiAdapter) MysqlFirstSql(ctx context.Context, dbName, sql string, args ...interface{}) (res map[string]interface{}, err error) {
	return GetEvApi().MysqlFirstSql(ctx, &dto.MysqlSelectReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Sql:           sql,
		Args:          args,
	})
}

// MysqlDbs 获取MySQL数据库列表
// 参数：
//   - ctx: 上下文
//
// 返回：
//   - dbs: 数据库名称列表
//   - err: 错误信息
func (this *EvApiAdapter) MysqlDbs(ctx context.Context) (dbs []string, err error) {
	return GetEvApi().MysqlDbs(ctx, &dto.MysqlDbsReq{
		EsConnectData: this.buildEsConnectData(),
	})
}

// MysqlTables 获取MySQL数据库中指定数据库的表列表
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//
// 返回：
//   - tables: 表名称列表
//   - err: 错误信息
func (this *EvApiAdapter) MysqlTables(ctx context.Context, dbName string) (tables []string, err error) {
	return GetEvApi().MysqlTables(ctx, &dto.MysqlTablesReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
	})
}

// RedisExecCommand 执行Redis命令
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库编号
//   - args: Redis命令参数
//
// 返回：
//   - data: 执行结果
//   - err: 错误信息
func (this *EvApiAdapter) RedisExecCommand(ctx context.Context, dbName int, args ...interface{}) (data interface{}, err error) {
	return GetEvApi().RedisExecCommand(ctx, &dto.RedisExecReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Args:          args,
	})
}

// ExecMongoCommand 执行MongoDB命令
// 参数：
//   - ctx: 上下文
//   - dbName: 数据库名称
//   - command: MongoDB命令
//   - timeout: 超时时间
//
// 返回：
//   - res: 执行结果
//   - err: 错误信息
func (this *EvApiAdapter) ExecMongoCommand(ctx context.Context, dbName string, command bson.D, timeout time.Duration) (res bson.M, err error) {
	return GetEvApi().ExecMongoCommand(ctx, &dto.MongoExecReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Command:       command,
		Timeout:       timeout,
	})
}



// EsCatNodes 获取ES节点信息
// 参数：
//   - ctx: 上下文
//   - h: 显示的字段列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCatNodes(ctx context.Context, h []string) (res *proto.Response, err error) {
	return GetEvApi().EsCatNodes(ctx, dto.CatNodesReq{
		EsConnectData:  this.buildEsConnectData(),
		CatNodeReqData: dto.CatNodeReqData{H: h},
	})
}

// EsClusterStats 获取ES集群统计信息
// 参数：
//   - ctx: 上下文
//   - human: 是否使用人类可读格式
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsClusterStats(ctx context.Context, human bool) (res *proto.Response, err error) {
	return GetEvApi().EsClusterStats(ctx, dto.ClusterStatsReq{
		EsConnectData:       this.buildEsConnectData(),
		ClusterStatsReqData: dto.ClusterStatsReqData{Human: human},
	})
}

// EsIndicesSegmentsRequest 获取ES索引段信息
// 参数：
//   - ctx: 上下文
//   - human: 是否使用人类可读格式
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsIndicesSegmentsRequest(ctx context.Context, human bool) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesSegmentsRequest(ctx, dto.IndicesSegmentsRequest{
		EsConnectData:              this.buildEsConnectData(),
		IndicesSegmentsRequestData: dto.IndicesSegmentsRequestData{Human: human},
	})
}

// EsPerformRequest 执行自定义ES请求
// 参数：
//   - ctx: 上下文
//   - req: HTTP请求对象
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsPerformRequest(ctx context.Context, req *http.Request) (res *proto.Response, err error) {

	request := &dto.Request{
		Method:        req.Method,
		URL:           req.URL,
		Header:        req.Header,
		Form:          req.Form,
		PostForm:      req.PostForm,
		MultipartForm: req.MultipartForm,
	}
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		request.JsonBody = string(body)
	}

	return GetEvApi().EsPerformRequest(ctx, dto.PerformRequest{
		EsConnectData: this.buildEsConnectData(),
		Request:       request,
	})
}

// Ping 测试ES连接
// 参数：
//   - ctx: 上下文
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) Ping(ctx context.Context) (res *proto.Response, err error) {
	return GetEvApi().Ping(ctx, dto.PingReq{EsConnectData: this.buildEsConnectData()})
}

// EsRefresh 刷新ES索引
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsRefresh(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsRefresh(ctx, dto.RefreshReq{
		EsConnectData:  this.buildEsConnectData(),
		RefreshReqData: dto.RefreshReqData{IndexNames: indexNames},
	})
}

// EsOpen 打开ES索引
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsOpen(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsOpen(ctx, dto.OpenReq{
		EsConnectData: this.buildEsConnectData(),
		OpenReqData:   dto.OpenReqData{IndexNames: indexNames},
	})
}

// EsFlush 强制刷新ES索引
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsFlush(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsFlush(ctx, dto.FlushReq{
		EsConnectData: this.buildEsConnectData(),
		FlushReqData:  dto.FlushReqData{IndexNames: indexNames},
	})
}

// EsIndicesClearCache 清除ES索引缓存
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsIndicesClearCache(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesClearCache(ctx, dto.IndicesClearCacheReq{
		EsConnectData:            this.buildEsConnectData(),
		IndicesClearCacheReqData: dto.IndicesClearCacheReqData{IndexNames: indexNames},
	})
}

// EsIndicesClose 关闭ES索引
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsIndicesClose(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesClose(ctx, dto.IndicesCloseReq{
		EsConnectData:       this.buildEsConnectData(),
		IndicesCloseReqData: dto.IndicesCloseReqData{IndexNames: indexNames},
	})
}

// EsIndicesForcemerge 强制合并ES索引段
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//   - maxNumSegments: 最大段数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsIndicesForcemerge(ctx context.Context, indexNames []string, maxNumSegments *int) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesForcemerge(ctx, dto.IndicesForcemergeReq{
		EsConnectData:            this.buildEsConnectData(),
		IndicesForcemergeReqData: dto.IndicesForcemergeReqData{IndexNames: indexNames, MaxNumSegments: maxNumSegments},
	})
}

// EsDeleteByQuery 按查询条件删除ES文档
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//   - documents: 文档类型列表
//   - body: 查询条件
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsDeleteByQuery(ctx context.Context, indexNames []string, documents []string, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsDeleteByQuery(ctx, dto.DeleteByQueryReq{
		EsConnectData:        this.buildEsConnectData(),
		DeleteByQueryReqData: dto.DeleteByQueryReqData{IndexNames: indexNames, Documents: documents, Body: body},
	})
}

// EsSnapshotCreate 创建ES快照
// 参数：
//   - ctx: 上下文
//   - repository: 仓库名称
//   - snapshot: 快照名称
//   - waitForCompletion: 是否等待完成
//   - reqJson: 请求JSON
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsSnapshotCreate(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotCreate(ctx, dto.SnapshotCreateReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotCreateReqData: dto.SnapshotCreateReqData{Repository: repository, Snapshot: snapshot, WaitForCompletion: waitForCompletion, ReqJson: reqJson},
	})
}

// EsSnapshotDelete 删除ES快照
// 参数：
//   - ctx: 上下文
//   - repository: 仓库名称
//   - snapshot: 快照名称
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsSnapshotDelete(ctx context.Context, repository string, snapshot string) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotDelete(ctx, dto.SnapshotDeleteReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotDeleteReqData: dto.SnapshotDeleteReqData{Repository: repository, Snapshot: snapshot},
	})
}

// EsRestoreSnapshot 恢复ES快照
// 参数：
//   - ctx: 上下文
//   - repository: 仓库名称
//   - snapshot: 快照名称
//   - waitForCompletion: 是否等待完成
//   - reqJson: 请求JSON
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsRestoreSnapshot(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	return GetEvApi().EsRestoreSnapshot(ctx, dto.RestoreSnapshotReq{
		EsConnectData:          this.buildEsConnectData(),
		RestoreSnapshotReqData: dto.RestoreSnapshotReqData{Repository: repository, Snapshot: snapshot, WaitForCompletion: waitForCompletion, ReqJson: reqJson},
	})
}

// EsSnapshotStatus 获取ES快照状态
// 参数：
//   - ctx: 上下文
//   - repository: 仓库名称
//   - snapshot: 快照名称列表
//   - ignoreUnavailable: 是否忽略不可用的快照
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsSnapshotStatus(ctx context.Context, repository string, snapshot []string, ignoreUnavailable *bool) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotStatus(ctx, dto.SnapshotStatusReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotStatusReqData: dto.SnapshotStatusReqData{Repository: repository, Snapshot: snapshot, IgnoreUnavailable: ignoreUnavailable},
	})
}

// EsSnapshotGetRepository 获取ES快照仓库
// 参数：
//   - ctx: 上下文
//   - repository: 仓库名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsSnapshotGetRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotGetRepository(ctx, dto.SnapshotGetRepositoryReq{
		EsConnectData:                this.buildEsConnectData(),
		SnapshotGetRepositoryReqData: dto.SnapshotGetRepositoryReqData{Repository: repository},
	})
}

// EsSnapshotCreateRepository 创建ES快照仓库
// 参数：
//   - ctx: 上下文
//   - repository: 仓库名称
//   - reqJson: 请求JSON
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsSnapshotCreateRepository(ctx context.Context, repository string, reqJson proto.Json) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotCreateRepository(ctx, dto.SnapshotCreateRepositoryReq{
		EsConnectData:                   this.buildEsConnectData(),
		SnapshotCreateRepositoryReqData: dto.SnapshotCreateRepositoryReqData{Repository: repository, ReqJson: reqJson},
	})
}

// EsSnapshotDeleteRepository 删除ES快照仓库
// 参数：
//   - ctx: 上下文
//   - repository: 仓库名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsSnapshotDeleteRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotDeleteRepository(ctx, dto.SnapshotDeleteRepositoryReq{
		EsConnectData:                   this.buildEsConnectData(),
		SnapshotDeleteRepositoryReqData: dto.SnapshotDeleteRepositoryReqData{Repository: repository},
	})
}

// EsGetIndices 获取ES索引列表
// 参数：
//   - ctx: 上下文
//   - catIndicesRequest: 索引列表请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsGetIndices(ctx context.Context, catIndicesRequest proto.CatIndicesRequest) (res *proto.Response, err error) {
	return GetEvApi().EsGetIndices(ctx, dto.GetIndicesReq{
		EsConnectData:     this.buildEsConnectData(),
		GetIndicesReqData: dto.GetIndicesReqData{CatIndicesRequest: catIndicesRequest},
	})
}

// EsCatHealth 获取ES集群健康状态
// 参数：
//   - ctx: 上下文
//   - catRequest: 健康状态请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCatHealth(ctx context.Context, catRequest proto.CatHealthRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatHealth(ctx, dto.CatHealthReq{
		EsConnectData:    this.buildEsConnectData(),
		CatHealthReqData: dto.CatHealthReqData{CatRequest: catRequest},
	})
}

// EsCatShards 获取ES分片信息
// 参数：
//   - ctx: 上下文
//   - catRequest: 分片信息请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCatShards(ctx context.Context, catRequest proto.CatShardsRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatShards(ctx, dto.CatShardsReq{
		EsConnectData:    this.buildEsConnectData(),
		CatShardsReqData: dto.CatShardsReqData{CatRequest: catRequest},
	})
}

// EsCatCount 获取ES索引文档计数
// 参数：
//   - ctx: 上下文
//   - catRequest: 计数请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCatCount(ctx context.Context, catRequest proto.CatCountRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatCount(ctx, dto.CatCountReq{
		EsConnectData:   this.buildEsConnectData(),
		CatCountReqData: dto.CatCountReqData{CatRequest: catRequest},
	})
}

// EsCatAllocationRequest 获取ES分片分配信息
// 参数：
//   - ctx: 上下文
//   - catRequest: 分配信息请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCatAllocationRequest(ctx context.Context, catRequest proto.CatAllocationRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatAllocationRequest(ctx, dto.CatAllocationRequest{
		EsConnectData:            this.buildEsConnectData(),
		CatAllocationRequestData: dto.CatAllocationRequestData{CatRequest: catRequest},
	})
}

// EsCatAliases 获取ES别名信息
// 参数：
//   - ctx: 上下文
//   - catRequest: 别名信息请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCatAliases(ctx context.Context, catRequest proto.CatAliasesRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatAliases(ctx, dto.CatAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		CatAliasesReqData: dto.CatAliasesReqData{CatRequest: catRequest},
	})
}

// EsDelete 删除ES文档
// 参数：
//   - ctx: 上下文
//   - deleteRequest: 删除请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsDelete(ctx context.Context, deleteRequest proto.DeleteRequest) (res *proto.Response, err error) {
	return GetEvApi().EsDelete(ctx, dto.DeleteReq{
		EsConnectData: this.buildEsConnectData(),
		DeleteReqData: dto.DeleteReqData{DeleteRequest: deleteRequest},
	})
}

// EsUpdate 更新ES文档
// 参数：
//   - ctx: 上下文
//   - updateRequest: 更新请求参数
//   - body: 更新数据
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsUpdate(ctx context.Context, updateRequest proto.UpdateRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsUpdate(ctx, dto.UpdateReq{
		EsConnectData: this.buildEsConnectData(),
		UpdateReqData: dto.UpdateReqData{UpdateRequest: updateRequest, Body: body},
	})
}

// EsCreate 创建ES文档
// 参数：
//   - ctx: 上下文
//   - createRequest: 创建请求参数
//   - body: 创建数据
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCreate(ctx context.Context, createRequest proto.CreateRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsCreate(ctx, dto.CreateReq{
		EsConnectData: this.buildEsConnectData(),
		CreateReqData: dto.CreateReqData{CreateRequest: createRequest, Body: body},
	})
}

// EsSearch 搜索ES文档
// 参数：
//   - ctx: 上下文
//   - searchRequest: 搜索请求参数
//   - query: 搜索查询
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsSearch(ctx context.Context, searchRequest proto.SearchRequest, query interface{}) (res *proto.Response, err error) {
	t := time.Now()
	defer func() {
		log.Println("lose time", time.Now().Sub(t).String())
	}()

	return GetEvApi().EsSearch(ctx, dto.SearchReq{
		EsConnectData: this.buildEsConnectData(),
		SearchReqData: dto.SearchReqData{SearchRequest: searchRequest, Query: query},
	})
}

// EsIndicesPutSettingsRequest 设置ES索引配置
// 参数：
//   - ctx: 上下文
//   - indexSettingsRequest: 索引配置请求参数
//   - body: 配置数据
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsIndicesPutSettingsRequest(ctx context.Context, indexSettingsRequest proto.IndicesPutSettingsRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesPutSettingsRequest(ctx, dto.IndicesPutSettingsRequest{
		EsConnectData:                 this.buildEsConnectData(),
		IndicesPutSettingsRequestData: dto.IndicesPutSettingsRequestData{IndexSettingsRequest: indexSettingsRequest, Body: body},
	})
}

// EsCreateIndex 创建ES索引
// 参数：
//   - ctx: 上下文
//   - indexCreateRequest: 创建索引请求参数
//   - body: 索引配置数据
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsCreateIndex(ctx context.Context, indexCreateRequest proto.IndicesCreateRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsCreateIndex(ctx, dto.CreateIndexReq{
		EsConnectData:      this.buildEsConnectData(),
		CreateIndexReqData: dto.CreateIndexReqData{IndexCreateRequest: indexCreateRequest, Body: body},
	})
}

// EsDeleteIndex 删除ES索引
// 参数：
//   - ctx: 上下文
//   - indicesDeleteRequest: 删除索引请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsDeleteIndex(ctx context.Context, indicesDeleteRequest proto.IndicesDeleteRequest) (res *proto.Response, err error) {
	return GetEvApi().EsDeleteIndex(ctx, dto.DeleteIndexReq{
		EsConnectData:      this.buildEsConnectData(),
		DeleteIndexReqData: dto.DeleteIndexReqData{IndicesDeleteRequest: indicesDeleteRequest},
	})
}

// EsReindex 重新索引ES数据
// 参数：
//   - ctx: 上下文
//   - reindexRequest: 重新索引请求参数
//   - body: 重新索引配置
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsReindex(ctx context.Context, reindexRequest proto.ReindexRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsReindex(ctx, dto.ReindexReq{
		EsConnectData:  this.buildEsConnectData(),
		ReindexReqData: dto.ReindexReqData{ReindexRequest: reindexRequest, Body: body},
	})
}

// EsIndicesGetSettingsRequest 获取ES索引设置
// 参数：
//   - ctx: 上下文
//   - indicesGetSettingsRequest: 获取索引设置请求参数
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsIndicesGetSettingsRequest(ctx context.Context, indicesGetSettingsRequest proto.IndicesGetSettingsRequest) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesGetSettingsRequest(ctx, dto.IndicesGetSettingsRequestReq{
		EsConnectData:                    this.buildEsConnectData(),
		IndicesGetSettingsRequestReqData: dto.IndicesGetSettingsRequestReqData{IndicesGetSettingsRequest: indicesGetSettingsRequest},
	})
}

// EsPutMapping 设置ES索引映射
// 参数：
//   - ctx: 上下文
//   - indicesPutMappingRequest: 设置映射请求参数
//   - body: 映射配置
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsPutMapping(ctx context.Context, indicesPutMappingRequest proto.IndicesPutMappingRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsPutMapping(ctx, dto.PutMappingReq{
		EsConnectData:     this.buildEsConnectData(),
		PutMappingReqData: dto.PutMappingReqData{IndicesPutMappingRequest: indicesPutMappingRequest, Body: body},
	})
}

// EsGetMapping 获取ES索引映射
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsGetMapping(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsGetMapping(ctx, dto.GetMappingReq{
		EsConnectData:     this.buildEsConnectData(),
		GetMappingReqData: dto.GetMappingReqData{IndexNames: indexNames},
	})
}

// EsGetAliases 获取ES索引别名
// 参数：
//   - ctx: 上下文
//   - indexNames: 索引名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsGetAliases(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsGetAliases(ctx, dto.GetAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		GetAliasesReqData: dto.GetAliasesReqData{IndexNames: indexNames},
	})
}

// EsAddAliases 添加ES索引别名
// 参数：
//   - ctx: 上下文
//   - indexName: 索引名称列表
//   - aliasName: 别名名称
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsAddAliases(ctx context.Context, indexName []string, aliasName string) (res *proto.Response, err error) {
	return GetEvApi().EsAddAliases(ctx, dto.AddAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		AddAliasesReqData: dto.AddAliasesReqData{IndexName: indexName, AliasName: aliasName},
	})
}

// EsRemoveAliases 移除ES索引别名
// 参数：
//   - ctx: 上下文
//   - indexName: 索引名称列表
//   - aliasName: 别名名称列表
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsRemoveAliases(ctx context.Context, indexName []string, aliasName []string) (res *proto.Response, err error) {
	return GetEvApi().EsRemoveAliases(ctx, dto.RemoveAliasesReq{
		EsConnectData:        this.buildEsConnectData(),
		RemoveAliasesReqData: dto.RemoveAliasesReqData{IndexName: indexName, AliasName: aliasName},
	})
}

// EsMoveToAnotherIndexAliases 将别名移动到另一个索引
// 参数：
//   - ctx: 上下文
//   - body: 别名操作配置
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsMoveToAnotherIndexAliases(ctx context.Context, body proto.AliasAction) (res *proto.Response, err error) {
	return GetEvApi().EsMoveToAnotherIndexAliases(ctx, dto.MoveToAnotherIndexAliasesReq{
		EsConnectData:                    this.buildEsConnectData(),
		MoveToAnotherIndexAliasesReqData: dto.MoveToAnotherIndexAliasesReqData{Body: body},
	})
}

// EsTaskList 获取ES任务列表
// 参数：
//   - ctx: 上下文
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsTaskList(ctx context.Context) (res *proto.Response, err error) {
	return GetEvApi().EsTaskList(ctx, dto.TaskListReq{
		EsConnectData: this.buildEsConnectData(),
	})
}

// EsTasksCancel 取消ES任务
// 参数：
//   - ctx: 上下文
//   - taskId: 任务ID
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *EvApiAdapter) EsTasksCancel(ctx context.Context, taskId string) (res *proto.Response, err error) {
	return GetEvApi().EsTasksCancel(ctx, dto.TasksCancelReq{
		EsConnectData:      this.buildEsConnectData(),
		TasksCancelReqData: dto.TasksCancelReqData{TaskId: taskId},
	})
}

// buildEsConnectData 构建ES连接数据
// 返回：
//   - dto.EsConnectData: ES连接数据对象
func (this *EvApiAdapter) buildEsConnectData() dto.EsConnectData {
	return dto.EsConnectData{
		UserID:    this.UserId,
		EsConnect: this.ConnId,
	}
}
