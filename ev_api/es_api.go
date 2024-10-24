package ev_api

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// 用于更方便的操作es
type EvApiAdapter struct {
	EsConnId int
	UserId   int
}

func NewEvWrapApi(esConnId int, userId int) *EvApiAdapter {
	return &EvApiAdapter{EsConnId: esConnId, UserId: userId}
}

// 执行sql
func (this *EvApiAdapter) StoreExec(ctx context.Context, sql string, args ...interface{}) (rowsAffected int64, err error) {
	return GetEvApi().StoreExec(ctx, sql, args...)
}

// 查询索引 dist参数必须是一个切片
func (this *EvApiAdapter) StoreSelect(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	return GetEvApi().StoreSelect(ctx, dest, sql, args...)
}

func (this *EvApiAdapter) StoreFirst(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	return GetEvApi().StoreFirst(ctx, dest, sql, args...)
}

func (this *EvApiAdapter) LoadDebugPlugin(ctx context.Context, req *dto.LoadDebugPlugin) (err error) {
	return GetEvApi().LoadDebugPlugin(ctx, req)
}

func (this *EvApiAdapter) StopDebugPlugin(ctx context.Context, req *dto.StopDebugPlugin) (err error) {
	return GetEvApi().StopDebugPlugin(ctx, req)
}

func (this *EvApiAdapter) EsRunDsl(ctx context.Context, req *dto.PluginRunDsl2) (res *proto.Response, err error) {

	return GetEvApi().EsRunDsl(ctx, &dto.PluginRunDsl{
		EsConnectData: &dto.EsConnectData{
			UserID:    this.UserId,
			EsConnect: this.EsConnId,
		},
		Params:     req.Params,
		HttpMethod: req.HttpMethod,
		Path:       req.Path,
		Dsl:        req.Dsl,
	})
}

func (this *EvApiAdapter) EsVersion() (version int, err error) {
	verson, err := GetEvApi().EsVersion(context.Background(), this.buildEsConnectData())
	if err != nil {
		logger.DefaultLogger.Error("get es version err", err)
		return 0, err
	}
	return verson, nil
}

func (this *EvApiAdapter) MysqlExecSql(ctx context.Context, dbName, sql string, args ...interface{}) (rowsAffected int64, err error) {
	return GetEvApi().MysqlExecSql(ctx, &dto.MysqlExecReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Sql:           sql,
		Args:          args,
	})
}

func (this *EvApiAdapter) MysqlSelectSql(ctx context.Context, dbName, sql string, args ...interface{}) (res []map[string]interface{}, err error) {
	return GetEvApi().MysqlSelectSql(ctx, &dto.MysqlSelectReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Sql:           sql,
		Args:          args,
	})
}

func (this *EvApiAdapter) MysqlFirstSql(ctx context.Context, dbName, sql string, args ...interface{}) (res map[string]interface{}, err error) {
	return GetEvApi().MysqlFirstSql(ctx, &dto.MysqlSelectReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,
		Sql:           sql,
		Args:          args,
	})
}

func (this *EvApiAdapter) RedisExecCommand(ctx context.Context, dbName int, args ...interface{}) (data interface{}, err error) {
	return GetEvApi().RedisExecCommand(ctx, &dto.RedisExecReq{
		EsConnectData: this.buildEsConnectData(),
		DbName:        dbName,

		Args: args,
	})
}

func (this *EvApiAdapter) EsCatNodes(ctx context.Context, h []string) (res *proto.Response, err error) {
	return GetEvApi().EsCatNodes(ctx, dto.CatNodesReq{
		EsConnectData:  this.buildEsConnectData(),
		CatNodeReqData: dto.CatNodeReqData{H: h},
	})
}

func (this *EvApiAdapter) EsClusterStats(ctx context.Context, human bool) (res *proto.Response, err error) {
	return GetEvApi().EsClusterStats(ctx, dto.ClusterStatsReq{
		EsConnectData:       this.buildEsConnectData(),
		ClusterStatsReqData: dto.ClusterStatsReqData{Human: human},
	})
}

func (this *EvApiAdapter) EsIndicesSegmentsRequest(ctx context.Context, human bool) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesSegmentsRequest(ctx, dto.IndicesSegmentsRequest{
		EsConnectData:              this.buildEsConnectData(),
		IndicesSegmentsRequestData: dto.IndicesSegmentsRequestData{Human: human},
	})
}

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

func (this *EvApiAdapter) Ping(ctx context.Context) (res *proto.Response, err error) {
	return GetEvApi().Ping(ctx, dto.PingReq{EsConnectData: this.buildEsConnectData()})
}

func (this *EvApiAdapter) EsRefresh(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsRefresh(ctx, dto.RefreshReq{
		EsConnectData:  this.buildEsConnectData(),
		RefreshReqData: dto.RefreshReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) EsOpen(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsOpen(ctx, dto.OpenReq{
		EsConnectData: this.buildEsConnectData(),
		OpenReqData:   dto.OpenReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) EsFlush(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsFlush(ctx, dto.FlushReq{
		EsConnectData: this.buildEsConnectData(),
		FlushReqData:  dto.FlushReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) EsIndicesClearCache(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesClearCache(ctx, dto.IndicesClearCacheReq{
		EsConnectData:            this.buildEsConnectData(),
		IndicesClearCacheReqData: dto.IndicesClearCacheReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) EsIndicesClose(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesClose(ctx, dto.IndicesCloseReq{
		EsConnectData:       this.buildEsConnectData(),
		IndicesCloseReqData: dto.IndicesCloseReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) EsIndicesForcemerge(ctx context.Context, indexNames []string, maxNumSegments *int) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesForcemerge(ctx, dto.IndicesForcemergeReq{
		EsConnectData:            this.buildEsConnectData(),
		IndicesForcemergeReqData: dto.IndicesForcemergeReqData{IndexNames: indexNames, MaxNumSegments: maxNumSegments},
	})
}

func (this *EvApiAdapter) EsDeleteByQuery(ctx context.Context, indexNames []string, documents []string, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsDeleteByQuery(ctx, dto.DeleteByQueryReq{
		EsConnectData:        this.buildEsConnectData(),
		DeleteByQueryReqData: dto.DeleteByQueryReqData{IndexNames: indexNames, Documents: documents, Body: body},
	})
}

func (this *EvApiAdapter) EsSnapshotCreate(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotCreate(ctx, dto.SnapshotCreateReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotCreateReqData: dto.SnapshotCreateReqData{Repository: repository, Snapshot: snapshot, WaitForCompletion: waitForCompletion, ReqJson: reqJson},
	})
}

func (this *EvApiAdapter) EsSnapshotDelete(ctx context.Context, repository string, snapshot string) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotDelete(ctx, dto.SnapshotDeleteReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotDeleteReqData: dto.SnapshotDeleteReqData{Repository: repository, Snapshot: snapshot},
	})
}

func (this *EvApiAdapter) EsRestoreSnapshot(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	return GetEvApi().EsRestoreSnapshot(ctx, dto.RestoreSnapshotReq{
		EsConnectData:          this.buildEsConnectData(),
		RestoreSnapshotReqData: dto.RestoreSnapshotReqData{Repository: repository, Snapshot: snapshot, WaitForCompletion: waitForCompletion, ReqJson: reqJson},
	})
}

func (this *EvApiAdapter) EsSnapshotStatus(ctx context.Context, repository string, snapshot []string, ignoreUnavailable *bool) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotStatus(ctx, dto.SnapshotStatusReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotStatusReqData: dto.SnapshotStatusReqData{Repository: repository, Snapshot: snapshot, IgnoreUnavailable: ignoreUnavailable},
	})
}

func (this *EvApiAdapter) EsSnapshotGetRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotGetRepository(ctx, dto.SnapshotGetRepositoryReq{
		EsConnectData:                this.buildEsConnectData(),
		SnapshotGetRepositoryReqData: dto.SnapshotGetRepositoryReqData{Repository: repository},
	})
}

func (this *EvApiAdapter) EsSnapshotCreateRepository(ctx context.Context, repository string, reqJson proto.Json) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotCreateRepository(ctx, dto.SnapshotCreateRepositoryReq{
		EsConnectData:                   this.buildEsConnectData(),
		SnapshotCreateRepositoryReqData: dto.SnapshotCreateRepositoryReqData{Repository: repository, ReqJson: reqJson},
	})
}

func (this *EvApiAdapter) EsSnapshotDeleteRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	return GetEvApi().EsSnapshotDeleteRepository(ctx, dto.SnapshotDeleteRepositoryReq{
		EsConnectData:                   this.buildEsConnectData(),
		SnapshotDeleteRepositoryReqData: dto.SnapshotDeleteRepositoryReqData{Repository: repository},
	})
}

func (this *EvApiAdapter) EsGetIndices(ctx context.Context, catIndicesRequest proto.CatIndicesRequest) (res *proto.Response, err error) {
	return GetEvApi().EsGetIndices(ctx, dto.GetIndicesReq{
		EsConnectData:     this.buildEsConnectData(),
		GetIndicesReqData: dto.GetIndicesReqData{CatIndicesRequest: catIndicesRequest},
	})
}

func (this *EvApiAdapter) EsCatHealth(ctx context.Context, catRequest proto.CatHealthRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatHealth(ctx, dto.CatHealthReq{
		EsConnectData:    this.buildEsConnectData(),
		CatHealthReqData: dto.CatHealthReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) EsCatShards(ctx context.Context, catRequest proto.CatShardsRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatShards(ctx, dto.CatShardsReq{
		EsConnectData:    this.buildEsConnectData(),
		CatShardsReqData: dto.CatShardsReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) EsCatCount(ctx context.Context, catRequest proto.CatCountRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatCount(ctx, dto.CatCountReq{
		EsConnectData:   this.buildEsConnectData(),
		CatCountReqData: dto.CatCountReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) EsCatAllocationRequest(ctx context.Context, catRequest proto.CatAllocationRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatAllocationRequest(ctx, dto.CatAllocationRequest{
		EsConnectData:            this.buildEsConnectData(),
		CatAllocationRequestData: dto.CatAllocationRequestData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) EsCatAliases(ctx context.Context, catRequest proto.CatAliasesRequest) (res *proto.Response, err error) {
	return GetEvApi().EsCatAliases(ctx, dto.CatAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		CatAliasesReqData: dto.CatAliasesReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) EsDelete(ctx context.Context, deleteRequest proto.DeleteRequest) (res *proto.Response, err error) {
	return GetEvApi().EsDelete(ctx, dto.DeleteReq{
		EsConnectData: this.buildEsConnectData(),
		DeleteReqData: dto.DeleteReqData{DeleteRequest: deleteRequest},
	})
}

func (this *EvApiAdapter) EsUpdate(ctx context.Context, updateRequest proto.UpdateRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsUpdate(ctx, dto.UpdateReq{
		EsConnectData: this.buildEsConnectData(),
		UpdateReqData: dto.UpdateReqData{UpdateRequest: updateRequest, Body: body},
	})
}

func (this *EvApiAdapter) EsCreate(ctx context.Context, createRequest proto.CreateRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsCreate(ctx, dto.CreateReq{
		EsConnectData: this.buildEsConnectData(),
		CreateReqData: dto.CreateReqData{CreateRequest: createRequest, Body: body},
	})
}

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

func (this *EvApiAdapter) EsIndicesPutSettingsRequest(ctx context.Context, indexSettingsRequest proto.IndicesPutSettingsRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesPutSettingsRequest(ctx, dto.IndicesPutSettingsRequest{
		EsConnectData:                 this.buildEsConnectData(),
		IndicesPutSettingsRequestData: dto.IndicesPutSettingsRequestData{IndexSettingsRequest: indexSettingsRequest, Body: body},
	})
}

func (this *EvApiAdapter) EsCreateIndex(ctx context.Context, indexCreateRequest proto.IndicesCreateRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsCreateIndex(ctx, dto.CreateIndexReq{
		EsConnectData:      this.buildEsConnectData(),
		CreateIndexReqData: dto.CreateIndexReqData{IndexCreateRequest: indexCreateRequest, Body: body},
	})
}

func (this *EvApiAdapter) EsDeleteIndex(ctx context.Context, indicesDeleteRequest proto.IndicesDeleteRequest) (res *proto.Response, err error) {
	return GetEvApi().EsDeleteIndex(ctx, dto.DeleteIndexReq{
		EsConnectData:      this.buildEsConnectData(),
		DeleteIndexReqData: dto.DeleteIndexReqData{IndicesDeleteRequest: indicesDeleteRequest},
	})
}

func (this *EvApiAdapter) EsReindex(ctx context.Context, reindexRequest proto.ReindexRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsReindex(ctx, dto.ReindexReq{
		EsConnectData:  this.buildEsConnectData(),
		ReindexReqData: dto.ReindexReqData{ReindexRequest: reindexRequest, Body: body},
	})
}

func (this *EvApiAdapter) EsIndicesGetSettingsRequest(ctx context.Context, indicesGetSettingsRequest proto.IndicesGetSettingsRequest) (res *proto.Response, err error) {
	return GetEvApi().EsIndicesGetSettingsRequest(ctx, dto.IndicesGetSettingsRequestReq{
		EsConnectData:                    this.buildEsConnectData(),
		IndicesGetSettingsRequestReqData: dto.IndicesGetSettingsRequestReqData{IndicesGetSettingsRequest: indicesGetSettingsRequest},
	})
}

func (this *EvApiAdapter) EsPutMapping(ctx context.Context, indicesPutMappingRequest proto.IndicesPutMappingRequest, body interface{}) (res *proto.Response, err error) {
	return GetEvApi().EsPutMapping(ctx, dto.PutMappingReq{
		EsConnectData:     this.buildEsConnectData(),
		PutMappingReqData: dto.PutMappingReqData{IndicesPutMappingRequest: indicesPutMappingRequest, Body: body},
	})
}

func (this *EvApiAdapter) EsGetMapping(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsGetMapping(ctx, dto.GetMappingReq{
		EsConnectData:     this.buildEsConnectData(),
		GetMappingReqData: dto.GetMappingReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) EsGetAliases(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return GetEvApi().EsGetAliases(ctx, dto.GetAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		GetAliasesReqData: dto.GetAliasesReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) EsAddAliases(ctx context.Context, indexName []string, aliasName string) (res *proto.Response, err error) {
	return GetEvApi().EsAddAliases(ctx, dto.AddAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		AddAliasesReqData: dto.AddAliasesReqData{IndexName: indexName, AliasName: aliasName},
	})
}

func (this *EvApiAdapter) EsRemoveAliases(ctx context.Context, indexName []string, aliasName []string) (res *proto.Response, err error) {
	return GetEvApi().EsRemoveAliases(ctx, dto.RemoveAliasesReq{
		EsConnectData:        this.buildEsConnectData(),
		RemoveAliasesReqData: dto.RemoveAliasesReqData{IndexName: indexName, AliasName: aliasName},
	})
}

func (this *EvApiAdapter) EsMoveToAnotherIndexAliases(ctx context.Context, body proto.AliasAction) (res *proto.Response, err error) {
	return GetEvApi().EsMoveToAnotherIndexAliases(ctx, dto.MoveToAnotherIndexAliasesReq{
		EsConnectData:                    this.buildEsConnectData(),
		MoveToAnotherIndexAliasesReqData: dto.MoveToAnotherIndexAliasesReqData{Body: body},
	})
}

func (this *EvApiAdapter) EsTaskList(ctx context.Context) (res *proto.Response, err error) {
	return GetEvApi().EsTaskList(ctx, dto.TaskListReq{
		EsConnectData: this.buildEsConnectData(),
	})
}

func (this *EvApiAdapter) EsTasksCancel(ctx context.Context, taskId string) (res *proto.Response, err error) {
	return GetEvApi().EsTasksCancel(ctx, dto.TasksCancelReq{
		EsConnectData:      this.buildEsConnectData(),
		TasksCancelReqData: dto.TasksCancelReqData{TaskId: taskId},
	})
}

func (this *EvApiAdapter) buildEsConnectData() dto.EsConnectData {
	return dto.EsConnectData{

		UserID:    this.UserId,
		EsConnect: this.EsConnId,
	}
}
