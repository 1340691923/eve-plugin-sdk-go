package pkg

import (
	"context"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"net/http"
	"time"
)

type ClientInterface interface {
	Ping(ctx context.Context) (res *proto.Response, err error)
	EsVersion() (version int, err error)
	EsCatNodes(ctx context.Context, h []string) (res *proto.Response, err error)
	EsClusterStats(ctx context.Context, human bool) (res *proto.Response, err error)
	EsIndicesSegmentsRequest(ctx context.Context, human bool) (res *proto.Response, err error)
	EsPerformRequest(ctx context.Context, req *http.Request) (res *proto.Response, err error)
	EsRefresh(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	EsOpen(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	EsFlush(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	EsIndicesClearCache(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	EsIndicesClose(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	EsIndicesForcemerge(ctx context.Context, indexNames []string, maxNumSegments *int) (res *proto.Response, err error)
	EsDeleteByQuery(ctx context.Context, indexNames []string, documents []string, body interface{}) (res *proto.Response, err error)
	EsSnapshotCreate(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error)
	EsSnapshotDelete(ctx context.Context, repository string, snapshot string) (res *proto.Response, err error)
	EsRestoreSnapshot(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error)
	EsSnapshotStatus(ctx context.Context, repository string, snapshot []string, ignoreUnavailable *bool) (res *proto.Response, err error)

	EsSnapshotGetRepository(ctx context.Context, repository []string) (res *proto.Response, err error)
	EsSnapshotCreateRepository(ctx context.Context, repository string, reqJson proto.Json) (res *proto.Response, err error)
	EsSnapshotDeleteRepository(ctx context.Context, repository []string) (res *proto.Response, err error)

	EsGetIndices(ctx context.Context, catIndicesRequest proto.CatIndicesRequest) (res *proto.Response, err error)
	EsCatHealth(ctx context.Context, catRequest proto.CatHealthRequest) (res *proto.Response, err error)
	EsCatShards(ctx context.Context, catRequest proto.CatShardsRequest) (res *proto.Response, err error)
	EsCatCount(ctx context.Context, catRequest proto.CatCountRequest) (res *proto.Response, err error)
	EsCatAllocationRequest(ctx context.Context, catRequest proto.CatAllocationRequest) (res *proto.Response, err error)
	EsCatAliases(ctx context.Context, catRequest proto.CatAliasesRequest) (res *proto.Response, err error)

	EsDelete(ctx context.Context, deleteRequest proto.DeleteRequest) (res *proto.Response, err error)
	EsUpdate(ctx context.Context, updateRequest proto.UpdateRequest, body interface{}) (res *proto.Response, err error)
	EsCreate(ctx context.Context, createRequest proto.CreateRequest, body interface{}) (res *proto.Response, err error)
	EsSearch(ctx context.Context, searchRequest proto.SearchRequest, query interface{}) (res *proto.Response, err error)

	EsIndicesPutSettingsRequest(ctx context.Context, indexSettingsRequest proto.IndicesPutSettingsRequest, body interface{}) (res *proto.Response, err error)
	EsCreateIndex(ctx context.Context, indexCreateRequest proto.IndicesCreateRequest, body interface{}) (res *proto.Response, err error)
	EsDeleteIndex(ctx context.Context, indicesDeleteRequest proto.IndicesDeleteRequest) (res *proto.Response, err error)
	EsReindex(ctx context.Context, reindexRequest proto.ReindexRequest, body interface{}) (res *proto.Response, err error)
	EsIndicesGetSettingsRequest(ctx context.Context, indicesGetSettingsRequest proto.IndicesGetSettingsRequest) (res *proto.Response, err error)

	EsPutMapping(ctx context.Context, indicesPutMappingRequest proto.IndicesPutMappingRequest, body interface{}) (res *proto.Response, err error)
	EsGetMapping(ctx context.Context, indexNames []string) (res *proto.Response, err error)

	EsGetAliases(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	EsAddAliases(ctx context.Context, indexName []string, aliasName string) (res *proto.Response, err error)
	EsRemoveAliases(ctx context.Context, indexName []string, aliasName []string) (res *proto.Response, err error)
	EsMoveToAnotherIndexAliases(ctx context.Context, body proto.AliasAction) (res *proto.Response, err error)

	EsTaskList(ctx context.Context) (res *proto.Response, err error)
	EsTasksCancel(ctx context.Context, taskId string) (res *proto.Response, err error)

	//mysql数据源接口

	MysqlExecSql(ctx context.Context, dbName, sql string, args ...interface{}) (rowsAffected int64, err error)
	MysqlSelectSql(ctx context.Context, dbName, sql string, args ...interface{}) (list []map[string]interface{}, err error)
	MysqlFirstSql(ctx context.Context, dbName, sql string, args ...interface{}) (data map[string]interface{}, err error)

	//redis数据源接口
	RedisExecCommand(ctx context.Context, dbName int, args ...interface{}) (data interface{}, err error)

	//mongo数据源接口

	ExecMongoCommand(ctx context.Context, dbName string, command bson.D, timeout time.Duration) (res bson.M, err error)

	ShowMongoDbs(ctx context.Context) ([]string, error)
}
