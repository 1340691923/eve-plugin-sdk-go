package ev_api

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/backend/logger"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/vo"
	"net/http"
)

// 用于更方便的操作es
type EvApiAdapter struct {
	EsConnId int
	RoleId   int
	UserId   int
}

func NewEvWrapApi(esConnId int, roleId int, userId int) *EvApiAdapter {
	return &EvApiAdapter{EsConnId: esConnId, RoleId: roleId, UserId: userId}
}

// 执行sql
func (this *EvApiAdapter) StoreExec(ctx context.Context, sql string, args ...interface{}) (rowsAffected int64, err error) {
	return getEvApi().StoreExec(ctx, sql, args...)
}

// 查询索引 dist参数必须是一个切片
func (this *EvApiAdapter) StoreSelect(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	return getEvApi().StoreSelect(ctx, dest, sql, args...)
}

func (this *EvApiAdapter) StoreFirst(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error) {
	return getEvApi().StoreFirst(ctx, dest, sql, args...)
}

func (this *EvApiAdapter) LoadDebugPlugin(ctx context.Context, req *dto.LoadDebugPlugin) (err error) {
	return getEvApi().LoadDebugPlugin(ctx, req)
}

func (this *EvApiAdapter) StopDebugPlugin(ctx context.Context, req *dto.StopDebugPlugin) (err error) {
	return getEvApi().StopDebugPlugin(ctx, req)
}

func (this *EvApiAdapter) EsRunDsl(ctx context.Context, req *dto.PluginRunDsl2) (res *proto.Response, err error) {

	return getEvApi().EsRunDsl(ctx, &dto.PluginRunDsl{
		EsConnectData: &dto.EsConnectData{
			RoleID:    this.RoleId,
			UserID:    this.UserId,
			EsConnect: this.EsConnId,
		},
		Params:     req.Params,
		HttpMethod: req.HttpMethod,
		Path:       req.Path,
		Dsl:        req.Dsl,
	})
}

func (this *EvApiAdapter) DslHistoryList(ctx context.Context, req *dto.DslHistoryListReq) (res *vo.DisHistoryListRes, err error) {
	req.EsConnectData = this.buildEsConnectData()
	return getEvApi().DslHistoryListAction(ctx, req)
}

func (this *EvApiAdapter) CleanDslHistory(ctx context.Context) (err error) {
	req := this.buildEsConnectData()
	return getEvApi().CleanDslHistoryAction(ctx, &req)
}

func (this *EvApiAdapter) Version() int {
	verson, err := getEvApi().EsVersion(context.Background(), this.buildEsConnectData())
	if err != nil {
		logger.DefaultLogger.Error("get es version err", err)
		return 0
	}
	return verson
}

func (this *EvApiAdapter) CatNodes(ctx context.Context, h []string) (res *proto.Response, err error) {
	return getEvApi().CatNodes(ctx, dto.CatNodesReq{
		EsConnectData:  this.buildEsConnectData(),
		CatNodeReqData: dto.CatNodeReqData{H: h},
	})
}

func (this *EvApiAdapter) ClusterStats(ctx context.Context, human bool) (res *proto.Response, err error) {
	return getEvApi().ClusterStats(ctx, dto.ClusterStatsReq{
		EsConnectData:       this.buildEsConnectData(),
		ClusterStatsReqData: dto.ClusterStatsReqData{Human: human},
	})
}

func (this *EvApiAdapter) IndicesSegmentsRequest(ctx context.Context, human bool) (res *proto.Response, err error) {
	return getEvApi().IndicesSegmentsRequest(ctx, dto.IndicesSegmentsRequest{
		EsConnectData:              this.buildEsConnectData(),
		IndicesSegmentsRequestData: dto.IndicesSegmentsRequestData{Human: human},
	})
}

func (this *EvApiAdapter) PerformRequest(ctx context.Context, req *http.Request) (res *proto.Response, err error) {

	return getEvApi().PerformRequest(ctx, dto.PerformRequest{
		EsConnectData: this.buildEsConnectData(),
		Request:       req,
	})
}

func (this *EvApiAdapter) Ping(ctx context.Context) (res *proto.Response, err error) {
	return getEvApi().Ping(ctx, dto.PingReq{EsConnectData: this.buildEsConnectData()})
}

func (this *EvApiAdapter) Refresh(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return getEvApi().Refresh(ctx, dto.RefreshReq{
		EsConnectData:  this.buildEsConnectData(),
		RefreshReqData: dto.RefreshReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) Open(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return getEvApi().Open(ctx, dto.OpenReq{
		EsConnectData: this.buildEsConnectData(),
		OpenReqData:   dto.OpenReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) Flush(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return getEvApi().Flush(ctx, dto.FlushReq{
		EsConnectData: this.buildEsConnectData(),
		FlushReqData:  dto.FlushReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) IndicesClearCache(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return getEvApi().IndicesClearCache(ctx, dto.IndicesClearCacheReq{
		EsConnectData:            this.buildEsConnectData(),
		IndicesClearCacheReqData: dto.IndicesClearCacheReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) IndicesClose(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return getEvApi().IndicesClose(ctx, dto.IndicesCloseReq{
		EsConnectData:       this.buildEsConnectData(),
		IndicesCloseReqData: dto.IndicesCloseReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) IndicesForcemerge(ctx context.Context, indexNames []string, maxNumSegments *int) (res *proto.Response, err error) {
	return getEvApi().IndicesForcemerge(ctx, dto.IndicesForcemergeReq{
		EsConnectData:            this.buildEsConnectData(),
		IndicesForcemergeReqData: dto.IndicesForcemergeReqData{IndexNames: indexNames, MaxNumSegments: maxNumSegments},
	})
}

func (this *EvApiAdapter) DeleteByQuery(ctx context.Context, indexNames []string, documents []string, body interface{}) (res *proto.Response, err error) {
	return getEvApi().DeleteByQuery(ctx, dto.DeleteByQueryReq{
		EsConnectData:        this.buildEsConnectData(),
		DeleteByQueryReqData: dto.DeleteByQueryReqData{IndexNames: indexNames, Documents: documents, Body: body},
	})
}

func (this *EvApiAdapter) SnapshotCreate(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	return getEvApi().SnapshotCreate(ctx, dto.SnapshotCreateReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotCreateReqData: dto.SnapshotCreateReqData{Repository: repository, Snapshot: snapshot, WaitForCompletion: waitForCompletion, ReqJson: reqJson},
	})
}

func (this *EvApiAdapter) SnapshotDelete(ctx context.Context, repository string, snapshot string) (res *proto.Response, err error) {
	return getEvApi().SnapshotDelete(ctx, dto.SnapshotDeleteReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotDeleteReqData: dto.SnapshotDeleteReqData{Repository: repository, Snapshot: snapshot},
	})
}

func (this *EvApiAdapter) RestoreSnapshot(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	return getEvApi().RestoreSnapshot(ctx, dto.RestoreSnapshotReq{
		EsConnectData:          this.buildEsConnectData(),
		RestoreSnapshotReqData: dto.RestoreSnapshotReqData{Repository: repository, Snapshot: snapshot, WaitForCompletion: waitForCompletion, ReqJson: reqJson},
	})
}

func (this *EvApiAdapter) SnapshotStatus(ctx context.Context, repository string, snapshot []string, ignoreUnavailable *bool) (res *proto.Response, err error) {
	return getEvApi().SnapshotStatus(ctx, dto.SnapshotStatusReq{
		EsConnectData:         this.buildEsConnectData(),
		SnapshotStatusReqData: dto.SnapshotStatusReqData{Repository: repository, Snapshot: snapshot, IgnoreUnavailable: ignoreUnavailable},
	})
}

func (this *EvApiAdapter) SnapshotGetRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	return getEvApi().SnapshotGetRepository(ctx, dto.SnapshotGetRepositoryReq{
		EsConnectData:                this.buildEsConnectData(),
		SnapshotGetRepositoryReqData: dto.SnapshotGetRepositoryReqData{Repository: repository},
	})
}

func (this *EvApiAdapter) SnapshotCreateRepository(ctx context.Context, repository string, reqJson proto.Json) (res *proto.Response, err error) {
	return getEvApi().SnapshotCreateRepository(ctx, dto.SnapshotCreateRepositoryReq{
		EsConnectData:                   this.buildEsConnectData(),
		SnapshotCreateRepositoryReqData: dto.SnapshotCreateRepositoryReqData{Repository: repository, ReqJson: reqJson},
	})
}

func (this *EvApiAdapter) SnapshotDeleteRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	return getEvApi().SnapshotDeleteRepository(ctx, dto.SnapshotDeleteRepositoryReq{
		EsConnectData:                   this.buildEsConnectData(),
		SnapshotDeleteRepositoryReqData: dto.SnapshotDeleteRepositoryReqData{Repository: repository},
	})
}

func (this *EvApiAdapter) GetIndices(ctx context.Context, catIndicesRequest proto.CatIndicesRequest) (res *proto.Response, err error) {
	return getEvApi().GetIndices(ctx, dto.GetIndicesReq{
		EsConnectData:     this.buildEsConnectData(),
		GetIndicesReqData: dto.GetIndicesReqData{CatIndicesRequest: catIndicesRequest},
	})
}

func (this *EvApiAdapter) CatHealth(ctx context.Context, catRequest proto.CatHealthRequest) (res *proto.Response, err error) {
	return getEvApi().CatHealth(ctx, dto.CatHealthReq{
		EsConnectData:    this.buildEsConnectData(),
		CatHealthReqData: dto.CatHealthReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) CatShards(ctx context.Context, catRequest proto.CatShardsRequest) (res *proto.Response, err error) {
	return getEvApi().CatShards(ctx, dto.CatShardsReq{
		EsConnectData:    this.buildEsConnectData(),
		CatShardsReqData: dto.CatShardsReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) CatCount(ctx context.Context, catRequest proto.CatCountRequest) (res *proto.Response, err error) {
	return getEvApi().CatCount(ctx, dto.CatCountReq{
		EsConnectData:   this.buildEsConnectData(),
		CatCountReqData: dto.CatCountReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) CatAllocationRequest(ctx context.Context, catRequest proto.CatAllocationRequest) (res *proto.Response, err error) {
	return getEvApi().CatAllocationRequest(ctx, dto.CatAllocationRequest{
		EsConnectData:            this.buildEsConnectData(),
		CatAllocationRequestData: dto.CatAllocationRequestData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) CatAliases(ctx context.Context, catRequest proto.CatAliasesRequest) (res *proto.Response, err error) {
	return getEvApi().CatAliases(ctx, dto.CatAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		CatAliasesReqData: dto.CatAliasesReqData{CatRequest: catRequest},
	})
}

func (this *EvApiAdapter) Delete(ctx context.Context, deleteRequest proto.DeleteRequest) (res *proto.Response, err error) {
	return getEvApi().Delete(ctx, dto.DeleteReq{
		EsConnectData: this.buildEsConnectData(),
		DeleteReqData: dto.DeleteReqData{DeleteRequest: deleteRequest},
	})
}

func (this *EvApiAdapter) Update(ctx context.Context, updateRequest proto.UpdateRequest, body interface{}) (res *proto.Response, err error) {
	return getEvApi().Update(ctx, dto.UpdateReq{
		EsConnectData: this.buildEsConnectData(),
		UpdateReqData: dto.UpdateReqData{UpdateRequest: updateRequest, Body: body},
	})
}

func (this *EvApiAdapter) Create(ctx context.Context, createRequest proto.CreateRequest, body interface{}) (res *proto.Response, err error) {
	return getEvApi().Create(ctx, dto.CreateReq{
		EsConnectData: this.buildEsConnectData(),
		CreateReqData: dto.CreateReqData{CreateRequest: createRequest, Body: body},
	})
}

func (this *EvApiAdapter) Search(ctx context.Context, searchRequest proto.SearchRequest, query interface{}) (res *proto.Response, err error) {
	return getEvApi().Search(ctx, dto.SearchReq{
		EsConnectData: this.buildEsConnectData(),
		SearchReqData: dto.SearchReqData{SearchRequest: searchRequest, Query: query},
	})
}

func (this *EvApiAdapter) IndicesPutSettingsRequest(ctx context.Context, indexSettingsRequest proto.IndicesPutSettingsRequest, body interface{}) (res *proto.Response, err error) {
	return getEvApi().IndicesPutSettingsRequest(ctx, dto.IndicesPutSettingsRequest{
		EsConnectData:                 this.buildEsConnectData(),
		IndicesPutSettingsRequestData: dto.IndicesPutSettingsRequestData{IndexSettingsRequest: indexSettingsRequest, Body: body},
	})
}

func (this *EvApiAdapter) CreateIndex(ctx context.Context, indexCreateRequest proto.IndicesCreateRequest, body interface{}) (res *proto.Response, err error) {
	return getEvApi().CreateIndex(ctx, dto.CreateIndexReq{
		EsConnectData:      this.buildEsConnectData(),
		CreateIndexReqData: dto.CreateIndexReqData{IndexCreateRequest: indexCreateRequest, Body: body},
	})
}

func (this *EvApiAdapter) DeleteIndex(ctx context.Context, indicesDeleteRequest proto.IndicesDeleteRequest) (res *proto.Response, err error) {
	return getEvApi().DeleteIndex(ctx, dto.DeleteIndexReq{
		EsConnectData:      this.buildEsConnectData(),
		DeleteIndexReqData: dto.DeleteIndexReqData{IndicesDeleteRequest: indicesDeleteRequest},
	})
}

func (this *EvApiAdapter) Reindex(ctx context.Context, reindexRequest proto.ReindexRequest, body interface{}) (res *proto.Response, err error) {
	return getEvApi().Reindex(ctx, dto.ReindexReq{
		EsConnectData:  this.buildEsConnectData(),
		ReindexReqData: dto.ReindexReqData{ReindexRequest: reindexRequest, Body: body},
	})
}

func (this *EvApiAdapter) IndicesGetSettingsRequest(ctx context.Context, indicesGetSettingsRequest proto.IndicesGetSettingsRequest) (res *proto.Response, err error) {
	return getEvApi().IndicesGetSettingsRequest(ctx, dto.IndicesGetSettingsRequestReq{
		EsConnectData:                    this.buildEsConnectData(),
		IndicesGetSettingsRequestReqData: dto.IndicesGetSettingsRequestReqData{IndicesGetSettingsRequest: indicesGetSettingsRequest},
	})
}

func (this *EvApiAdapter) PutMapping(ctx context.Context, indicesPutMappingRequest proto.IndicesPutMappingRequest, body interface{}) (res *proto.Response, err error) {
	return getEvApi().PutMapping(ctx, dto.PutMappingReq{
		EsConnectData:     this.buildEsConnectData(),
		PutMappingReqData: dto.PutMappingReqData{IndicesPutMappingRequest: indicesPutMappingRequest, Body: body},
	})
}

func (this *EvApiAdapter) GetMapping(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return getEvApi().GetMapping(ctx, dto.GetMappingReq{
		EsConnectData:     this.buildEsConnectData(),
		GetMappingReqData: dto.GetMappingReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) GetAliases(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	return getEvApi().GetAliases(ctx, dto.GetAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		GetAliasesReqData: dto.GetAliasesReqData{IndexNames: indexNames},
	})
}

func (this *EvApiAdapter) AddAliases(ctx context.Context, indexName []string, aliasName string) (res *proto.Response, err error) {
	return getEvApi().AddAliases(ctx, dto.AddAliasesReq{
		EsConnectData:     this.buildEsConnectData(),
		AddAliasesReqData: dto.AddAliasesReqData{IndexName: indexName, AliasName: aliasName},
	})
}

func (this *EvApiAdapter) RemoveAliases(ctx context.Context, indexName []string, aliasName []string) (res *proto.Response, err error) {
	return getEvApi().RemoveAliases(ctx, dto.RemoveAliasesReq{
		EsConnectData:        this.buildEsConnectData(),
		RemoveAliasesReqData: dto.RemoveAliasesReqData{IndexName: indexName, AliasName: aliasName},
	})
}

func (this *EvApiAdapter) MoveToAnotherIndexAliases(ctx context.Context, body proto.AliasAction) (res *proto.Response, err error) {
	return getEvApi().MoveToAnotherIndexAliases(ctx, dto.MoveToAnotherIndexAliasesReq{
		EsConnectData:                    this.buildEsConnectData(),
		MoveToAnotherIndexAliasesReqData: dto.MoveToAnotherIndexAliasesReqData{Body: body},
	})
}

func (this *EvApiAdapter) TaskList(ctx context.Context) (res *proto.Response, err error) {
	return getEvApi().TaskList(ctx, dto.TaskListReq{
		EsConnectData: this.buildEsConnectData(),
	})
}

func (this *EvApiAdapter) TasksCancel(ctx context.Context, taskId string) (res *proto.Response, err error) {
	return getEvApi().TasksCancel(ctx, dto.TasksCancelReq{
		EsConnectData:      this.buildEsConnectData(),
		TasksCancelReqData: dto.TasksCancelReqData{TaskId: taskId},
	})
}

func (this *EvApiAdapter) buildEsConnectData() dto.EsConnectData {
	return dto.EsConnectData{
		RoleID:    this.RoleId,
		UserID:    this.UserId,
		EsConnect: this.EsConnId,
	}
}
