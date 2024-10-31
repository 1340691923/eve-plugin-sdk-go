package ev_api

import (
	"context"
	"errors"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/vo"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	"github.com/goccy/go-json"
	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
	protobuf "google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

type evApi struct {
	rpcKey   string
	rpcPort  string
	debug    bool
	pluginId string
	client   *fasthttp.Client
}

var (
	once     *sync.Once
	evApiObj *evApi
)

func init() {
	once = new(sync.Once)
}

func SetEvApi(rpcKey, rpcPort, pluginId string, debug bool) *evApi {
	client := &fasthttp.Client{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	once.Do(func() {
		evApiObj = &evApi{
			rpcKey:   rpcKey,
			rpcPort:  rpcPort,
			pluginId: pluginId,
			debug:    debug,
			client:   client,
		}
	})

	return evApiObj
}

func GetEvApi() *evApi {
	return evApiObj
}

func (this *evApi) EsVersion(ctx context.Context, req dto.EsConnectData) (version int, err error) {
	res := vo.ApiCommonRes{}
	err = this.request(ctx, "api/plugin_util/EsVersion", req, &res)
	if err != nil {
		return 0, err
	}
	return cast.ToInt(res.Data), nil
}

func (this *evApi) EsCatNodes(ctx context.Context, req dto.CatNodesReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsCatNodes", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsClusterStats(ctx context.Context, req dto.ClusterStatsReq) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsClusterStats", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsPerformRequest(ctx context.Context, req dto.PerformRequest) (res *proto.Response, err error) {

	res, err = this.requestProtobuf(ctx, "api/plugin_util/EsPerformRequest", req)
	if err != nil {
		return
	}
	return
}

func (this *evApi) EsIndicesSegmentsRequest(ctx context.Context, req dto.IndicesSegmentsRequest) (res *proto.Response, err error) {

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

// 查询索引 dist参数必须是一个切片
func (this *evApi) StoreSelect(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	data := &vo.SelectRes{}
	data.Result = &dest
	err = this.request(ctx, "api/plugin_util/SelectSql", &dto.SelectReq{Sql: sql, PluginId: this.pluginId, Args: args}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return err
	}
	return nil
}

func (this *evApi) StoreFirst(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	data := &vo.SelectRes{}
	data.Result = &dest
	err = this.request(ctx, "api/plugin_util/FirstSql", &dto.SelectReq{Sql: sql, PluginId: this.pluginId, Args: args}, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return err
	}
	return nil
}

func (this *evApi) LoadDebugPlugin(ctx context.Context, req *dto.LoadDebugPlugin) (err error) {
	err = this.request(ctx, "api/plugin_util/LoadDebugPlugin", req, &vo.ApiCommonRes{})
	if err != nil {
		return err
	}
	return nil
}

func (this *evApi) StopDebugPlugin(ctx context.Context, req *dto.StopDebugPlugin) (err error) {
	err = this.request(ctx, "api/plugin_util/LoadDebugPlugin", req, &vo.ApiCommonRes{})
	if err != nil {
		return err
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
func (this *evApi) MysqlSelectSql(ctx context.Context, req *dto.MysqlSelectReq) (result []map[string]interface{}, err error) {
	data := &vo.MysqlSelectSqlRes{}
	err = this.request(ctx, "api/plugin_util/MysqlSelectSql", req, &vo.ApiCommonRes{Data: data})
	if err != nil {
		return nil, err
	}
	return data.Result, nil
}

func (this *evApi) MysqlFirstSql(ctx context.Context, req *dto.MysqlSelectReq) (result map[string]interface{}, err error) {
	data := &vo.MysqlFirstSqlRes{}
	err = this.request(ctx, "api/plugin_util/MysqlFirstSql", req, &vo.ApiCommonRes{Data: data})
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

	err = json.Unmarshal(result.ResByte(), &res)

	if err != nil {
		return data, err
	}

	return res["data"], nil
}

func (this *evApi) ExecMongoCommand(ctx context.Context, req *dto.MongoExecReq) (data bson.M, err error) {

	res, err := this.requestProtobuf(ctx, "api/plugin_util/MongoExecCommand", req)
	if err != nil {
		return nil, err
	}
	if res.StatusErr() != nil {
		return nil, res.StatusErr()
	}

	data = map[string]interface{}{}

	err = json.Unmarshal(res.ResByte(), &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *evApi) ShowMongoDbs(ctx context.Context, req *dto.ShowMongoDbsReq) (dbList []string, err error) {
	res := &vo.ApiCommonRes{Data: dbList}
	err = this.request(ctx, "api/plugin_util/ShowMongoDbs", req, res)
	if err != nil {
		return nil, err
	}

	return cast.ToStringSlice(res.Data), nil
}

func (this *evApi) request(ctx context.Context, api API, requestData interface{}, result interface{}) error {
	var requestDataJSON = []byte(`{}`)
	if requestData != nil {
		requestDataJSON, _ = json.Marshal(requestData)
	}

	t1 := time.Now()
	res, err := this.SendRequest(ctx, api, requestDataJSON)
	if err != nil {
		return err
	}

	if this.debug {
		logger.DefaultLogger.Info("debug network",
			"api", api,
			"reqBody", string(requestDataJSON),
			"lose time", api, time.Now().Sub(t1).String())
	}

	err = json.Unmarshal(res, result)

	if err != nil {
		return err
	}

	switch result.(type) {
	case *vo.ApiCommonRes:
		return result.(*vo.ApiCommonRes).Error()
	}

	return nil
}

func (this *evApi) requestProtobuf(ctx context.Context, api API, requestData interface{}) (result *proto.Response, err error) {

	requestDataJSON, err := json.Marshal(requestData)

	if err != nil {
		log.Println("err", err)
		return nil, err
	}

	t1 := time.Now()
	res, err := this.SendRequest(ctx, api, requestDataJSON)
	if err != nil {
		return nil, err
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
		return nil, err
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

func (this *evApi) SendRequest(ctx context.Context, api API, requestDataJSON []byte) ([]byte, error) {
	// 构建请求对象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 释放请求对象，防止内存泄漏

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 释放响应对象，防止内存泄漏

	// 设置请求 URL
	url := fmt.Sprintf("http://127.0.0.1:%s/%s", this.rpcPort, api)
	req.SetRequestURI(url)

	// 设置请求方法为 POST
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	// 设置自定义头
	req.Header.Set("X-Plugin-ID", this.pluginId)
	//req.Header.Set("X-Plugin-Signature", signature)

	// 设置请求体
	req.SetBody(requestDataJSON)

	// 发起请求

	errCh := make(chan error, 1)

	// 启动异步请求
	go func() {
		t := time.Now()
		errCh <- this.client.Do(req, resp)
		log.Println("client.do", time.Now().Sub(t).String())
	}()

	select {
	case <-ctx.Done(): // 如果 context 超时或取消
		return nil, ctx.Err()
	case err := <-errCh: // 请求完成
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
	}

	// 返回响应体
	return resp.Body(), nil
}
