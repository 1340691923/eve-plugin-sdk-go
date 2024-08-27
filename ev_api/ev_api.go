package ev_api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/vo"
	"github.com/imroc/req/v2"
	"github.com/spf13/cast"
	"sync"
)

type evApi struct {
	rpcKey   string
	rpcPort  string
	debug    bool
	client   *req.Client
	pluginId string
}

var (
	once     *sync.Once
	evApiObj *evApi
)

func init() {
	once = new(sync.Once)
}

func SetEvApi(rpcKey, rpcPort, pluginId string, debug bool) *evApi {

	once.Do(func() {
		var client *req.Client

		if debug {
			client = req.C().DevMode()
		} else {
			client = req.C()
		}
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

func getEvApi() *evApi {
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

func (this *evApi) CatNodes(ctx context.Context, req dto.CatNodesReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/CatNodes", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) ClusterStats(ctx context.Context, req dto.ClusterStatsReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/ClusterStats", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) PerformRequest(ctx context.Context, req dto.PerformRequest) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/PerformRequest", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) IndicesSegmentsRequest(ctx context.Context, req dto.IndicesSegmentsRequest) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/IndicesSegmentsRequest", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Ping(ctx context.Context, req dto.PingReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Ping", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Refresh(ctx context.Context, req dto.RefreshReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Refresh", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Open(ctx context.Context, req dto.OpenReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Open", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Flush(ctx context.Context, req dto.FlushReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Flush", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) IndicesClearCache(ctx context.Context, req dto.IndicesClearCacheReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/IndicesClearCache", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) IndicesClose(ctx context.Context, req dto.IndicesCloseReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/IndicesClose", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) IndicesForcemerge(ctx context.Context, req dto.IndicesForcemergeReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/IndicesForcemerge", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) DeleteByQuery(ctx context.Context, req dto.DeleteByQueryReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/DeleteByQuery", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) SnapshotCreate(ctx context.Context, req dto.SnapshotCreateReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/SnapshotCreate", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) SnapshotDelete(ctx context.Context, req dto.SnapshotDeleteReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/SnapshotDelete", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) RestoreSnapshot(ctx context.Context, req dto.RestoreSnapshotReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/RestoreSnapshot", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) SnapshotStatus(ctx context.Context, req dto.SnapshotStatusReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/SnapshotStatus", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) SnapshotGetRepository(ctx context.Context, req dto.SnapshotGetRepositoryReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/SnapshotGetRepository", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) SnapshotCreateRepository(ctx context.Context, req dto.SnapshotCreateRepositoryReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/SnapshotCreateRepository", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) SnapshotDeleteRepository(ctx context.Context, req dto.SnapshotDeleteRepositoryReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/SnapshotDeleteRepository", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) GetIndices(ctx context.Context, req dto.GetIndicesReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/GetIndices", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) CatHealth(ctx context.Context, req dto.CatHealthReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/CatHealth", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) CatShards(ctx context.Context, req dto.CatShardsReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/CatShards", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) CatCount(ctx context.Context, req dto.CatCountReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/CatCount", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) CatAllocationRequest(ctx context.Context, req dto.CatAllocationRequest) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/CatAllocationRequest", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) CatAliases(ctx context.Context, req dto.CatAliasesReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/CatAliases", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Delete(ctx context.Context, req dto.DeleteReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Delete", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Update(ctx context.Context, req dto.UpdateReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Update", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Create(ctx context.Context, req dto.CreateReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Create", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Search(ctx context.Context, req dto.SearchReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Search", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) IndicesPutSettingsRequest(ctx context.Context, req dto.IndicesPutSettingsRequest) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/IndicesPutSettingsRequest", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) CreateIndex(ctx context.Context, req dto.CreateIndexReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/CreateIndex", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) DeleteIndex(ctx context.Context, req dto.DeleteIndexReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/DeleteIndex", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) Reindex(ctx context.Context, req dto.ReindexReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/Reindex", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) IndicesGetSettingsRequest(ctx context.Context, req dto.IndicesGetSettingsRequestReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/IndicesGetSettingsRequest", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) PutMapping(ctx context.Context, req dto.PutMappingReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/PutMapping", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) GetMapping(ctx context.Context, req dto.GetMappingReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/GetMapping", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) GetAliases(ctx context.Context, req dto.GetAliasesReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/GetAliases", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) AddAliases(ctx context.Context, req dto.AddAliasesReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/AddAliases", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) RemoveAliases(ctx context.Context, req dto.RemoveAliasesReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/RemoveAliases", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) MoveToAnotherIndexAliases(ctx context.Context, req dto.MoveToAnotherIndexAliasesReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/MoveToAnotherIndexAliases", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) TaskList(ctx context.Context, req dto.TaskListReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/TaskList", req, &commonRes)
	if err != nil {
		return
	}
	return
}

func (this *evApi) TasksCancel(ctx context.Context, req dto.TasksCancelReq) (res *proto.Response, err error) {
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	err = this.request(ctx, "api/plugin_util/TasksCancel", req, &commonRes)
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
	res = new(proto.Response)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res
	if req.Params != nil {
		req.Path = fmt.Sprintf("%s?%s", req.Path, req.Params.Encode())
	}
	err = this.request(ctx, "api/plugin_util/EsRunDsl", req, &commonRes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *evApi) DslHistoryListAction(ctx context.Context, req *dto.DslHistoryListReq) (res *vo.DisHistoryListRes, err error) {
	res = new(vo.DisHistoryListRes)
	commonRes := vo.ApiCommonRes{}
	commonRes.Data = res

	err = this.request(ctx, "api/plugin_util/DslHistoryListAction", req, &commonRes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *evApi) CleanDslHistoryAction(ctx context.Context, req *dto.EsConnectData) (err error) {
	err = this.request(ctx, "api/plugin_util/CleanDslHistoryAction", req, &vo.ApiCommonRes{})
	if err != nil {
		return err
	}
	return nil
}

func (this *evApi) request(ctx context.Context, api API, requestData interface{}, result interface{}) error {

	requestDataJSON, _ := json.Marshal(requestData)

	// 计算 HMAC-SHA256 签名
	mac := hmac.New(sha256.New, []byte(this.rpcKey))
	mac.Write(requestDataJSON)
	signatureBytes := mac.Sum(nil)

	// 将签名转换为 Base64 编码的字符串
	signature := base64.StdEncoding.EncodeToString(signatureBytes)

	_, err := this.client.R().
		SetHeaders(map[string]string{
			"X-Plugin-ID":        this.pluginId,
			"X-Plugin-Signature": signature,
		}).
		SetContext(ctx).
		SetBody(string(requestDataJSON)).
		SetResult(&result).
		Post(fmt.Sprintf("http://127.0.0.1:%s/%s", this.rpcPort, api))
	if err != nil {
		return err
	}
	switch result.(type) {
	case *vo.ApiCommonRes:
		return result.(*vo.ApiCommonRes).Error()
	}

	return nil
}
