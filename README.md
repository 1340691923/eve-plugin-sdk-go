# ElasticView Golang Plugin SDK

## 这是用于开发ElasticView后端插件的SDK

### 简介
本SDK提供了开发ElasticView后端插件所需的工具和接口，帮助开发者快速构建与ElasticView系统集成的插件。

### SDK介绍
- 提供了插件与ElasticView系统通信的标准接口
- 支持长连接数据推送
- 提供http资源/接口调用功能
- 内置SQL构建工具
- 提供便捷的Web路由鉴权


### 快速开始

#### 1. 初始化插件(以ev工具箱插件为例)

```go
package main

import (
	"context"
	"embed"
	_ "embed"
	"ev-plugin/backend/migrate"
	"ev-plugin/backend/router"
	"ev-plugin/frontend"
	"flag"
	"github.com/1340691923/eve-plugin-sdk-go/backend/plugin_server"
	"github.com/1340691923/eve-plugin-sdk-go/build"
)

//go:embed plugin.json
var pluginJsonBytes []byte

//go:embed logo.png
var logoPng embed.FS

func main() {
	//必须先解析参数
	flag.Parse()
	//启动插件服务
	plugin_server.Serve(plugin_server.ServeOpts{
		Assets: &plugin_server.Assets{
			PluginJsonBytes: pluginJsonBytes,//plugin.json 插件配置，必填
			FrontendFiles:   frontend.StatisFs,//前端工程编译后产物，必填
			Icon: logoPng,//logo图标，必填
		},
		ReadyCallBack: func(ctx context.Context) {
			// 插件就绪后的回调，可以在这里进行初始化操作，可选
		},
		Migration: &build.Gormigrate{Migrations: []*build.Migration{
			migrate.V0_0_1(),
		}}, //插件存储的版本迁移配置，可选
		RegisterRoutes: router.NewRouter, // 注册插件http接口路由，必填
	})
}


```

#### 2. 路由注册

```go
package router

import (
	"github.com/1340691923/eve-plugin-sdk-go/backend/web_engine"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewRouter 注册插件路由
func NewRouter(engine *web_engine.WebEngine) {
	// 创建API路由组
	userGroup := engine.Group("用戶管理", "/api/user")
	//参数从左到右依次为是否鉴权，接口备注，路由path,handler
	userGroup.GET(true, "获取用户列表", "/list", getUserList)
}

// handle处理函数示例
func getUserList(c *gin.Context) {

	// 获取EV API实例
	api := ev_api.GetEvApi()
	
	// 查询用户列表
	var users []User
	err := api.StoreSelect(c.Request.Context(), &users, 
		"SELECT * FROM users WHERE status = ? ORDER BY created_at DESC", 1)
	
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
			"code":500,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": users,
		"msg":"success"
	})
}
```

#### 4. 获取基座通讯API实例(内部已处理成单例)

```go
import "github.com/1340691923/eve-plugin-sdk-go/ev_api"

// 获取全局EV API实例
api := ev_api.GetEvApi()
```

### 插件数据存储相关（若无数据存储需求则跳过不看）

#### 3. 插件存储版本迁移配置

```go
package migrate

import "github.com/1340691923/eve-plugin-sdk-go/build"

// V0_0_1 创建版本迁移  分别为 sqlite 和 mysql
func V0_0_1() *build.Migration {
	return &build.Migration{
		ID: "v0.0.1", // 版本号
		SqliteMigrateSqls: []*build.ExecSql{
			{
				Sql: `CREATE TABLE IF NOT EXISTS user_settings (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER NOT NULL,
					setting_key TEXT NOT NULL,
					setting_value TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				)`,
				Args: []interface{}{},
			},
			{
				Sql: "CREATE INDEX idx_user_settings_user_id ON user_settings(user_id)",
				Args: []interface{}{},
			},
		},
		MysqlMigrateSqls: []*build.ExecSql{
			{
				Sql: `CREATE TABLE IF NOT EXISTS user_settings (
					id INT AUTO_INCREMENT PRIMARY KEY,
					user_id INT NOT NULL,
					setting_key VARCHAR(255) NOT NULL,
					setting_value TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
				) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
				Args: []interface{}{},
			},
			{
				Sql: "CREATE INDEX idx_user_settings_user_id ON user_settings(user_id)",
				Args: []interface{}{},
			},
		},
		SqliteRollback: []*build.ExecSql{
			{Sql: "DROP TABLE IF EXISTS user_settings", Args: []interface{}{}},
		},
		MysqlRollback: []*build.ExecSql{
			{Sql: "DROP TABLE IF EXISTS user_settings", Args: []interface{}{}},
		},
	}
}
```

#### 6. 插件存储操作

##### 查询操作
```go
// StoreSelect 执行查询操作
// 参数：
//   - ctx: 上下文
//   - dest: 结果接收器
//   - sql: SQL语句
//   - args: 参数
func (this *evApi) StoreSelect(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error)

// 使用示例
var users []User
err := api.StoreSelect(ctx, &users, "SELECT * FROM users WHERE age > ?", 18)
```

##### 单条查询
```go
// StoreFirst 查询单条记录
func (this *evApi) StoreFirst(ctx context.Context, dest interface{}, sql string, args ...interface{}) (err error)

// 使用示例
var user User
err := api.StoreFirst(ctx, &user, "SELECT * FROM users WHERE id = ?", userId)
```

##### 执行操作
```go
// StoreExec 执行SQL语句（INSERT、UPDATE、DELETE）
func (this *evApi) StoreExec(ctx context.Context, sql string, args ...interface{}) (rowsAffected int64, err error)

// 使用示例
rowsAffected, err := api.StoreExec(ctx, "UPDATE users SET name = ? WHERE id = ?", "新名称", userId)
```

##### 批量执行（内部有事务处理，单条sql执行不成功则都执行不成功）
```go
// StoreMoreExec 批量执行SQL语句
func (this *evApi) StoreMoreExec(ctx context.Context, sqls []dto.ExecSql) (err error)

// 使用示例
sqls := []dto.ExecSql{
    {Sql: "INSERT INTO users (name) VALUES (?)", Args: []interface{}{"用户1"}},
    {Sql: "INSERT INTO users (name) VALUES (?)", Args: []interface{}{"用户2"}},
}
err := api.StoreMoreExec(ctx, sqls)
```

##### 便捷数据操作
```go
// 保存数据
err := api.StoreSave(ctx, "users", map[string]interface{}{
    "name": "张三",
    "email": "zhangsan@example.com",
    "age": 25,
})

// 更新数据
rowsAffected, err := api.StoreUpdate(ctx, "users", 
    map[string]interface{}{"name": "李四"}, 
    "id = ?", userId)

// 删除数据
rowsAffected, err := api.StoreDelete(ctx, "users", "id = ?", userId)

// 插入或更新（Upsert）
err := api.StoreInsertOrUpdate(ctx, "user_settings", 
    map[string]interface{}{
        "user_id": userId,
        "setting_key": "theme",
        "setting_value": "dark",
    }, "user_id", "setting_key")
```

#### 7. SQL构建工具

```go
import "github.com/1340691923/eve-plugin-sdk-go/sql_builder"

// 使用SQL构建器
query := sql_builder.SqlBuilder.
    Select("id", "name", "email").
    From("users").
    Where(sql_builder.And{
        sql_builder.Eq{"status": 1},
        sql_builder.Gt{"age": 18},
    }).
    OrderBy("created_at DESC").
    Limit(10).
    Offset(sql_builder.CreatePage(2, 10)) // 第2页，每页10条

sql, args, err := query.ToSql()

```

### 更快捷的操作数据源

#### 9. 第三方Elasticsearch操作（操作后台设置的数据源，默认兼容Elasticsearc 6，7，8版本）

##### 获取ES版本
```go
// EsVersion 获取Elasticsearch版本
// 参数：
//   - ctx: 上下文
//   - req: ES连接数据
//
// 返回：
//   - version: ES版本号
//   - err: 错误信息
func (this *evApi) EsVersion(ctx context.Context, req dto.EsConnectData) (version int, err error)

// 使用示例
version, err := api.EsVersion(ctx, dto.EsConnectData{
    UserID:    userID,    // 用户ID
    EsConnect: connectID, // ES连接配置ID
})
```

##### 获取ES节点信息
```go
// EsCatNodes 获取ES节点信息
// 参数：
//   - ctx: 上下文
//   - req: Cat节点请求
//
// 返回：
//   - res: *proto.Response
//   - err: 错误信息
func (this *evApi) EsCatNodes(ctx context.Context, req dto.CatNodesReq) (res *proto.Response, err error)

// 使用示例
res, err := api.EsCatNodes(ctx, dto.CatNodesReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    CatNodeReqData: dto.CatNodeReqData{
        H: []string{"name", "heap.percent", "ram.percent", "cpu", "load_1m"},
    },
})
```

##### ES集群统计
```go
// EsClusterStats 获取ES集群统计信息
func (this *evApi) EsClusterStats(ctx context.Context, req dto.ClusterStatsReq) (res *proto.Response, err error)

// 使用示例
res, err := api.EsClusterStats(ctx, dto.ClusterStatsReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    ClusterStatsReqData: dto.ClusterStatsReqData{
        Human: true, // 返回人类可读的格式
    },
})
```

##### ES连接测试
```go
// Ping 测试ES连接
func (this *evApi) Ping(ctx context.Context, req dto.PingReq) (res *proto.Response, err error)

// 使用示例
res, err := api.Ping(ctx, dto.PingReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
})
```

##### ES搜索操作
```go
// EsSearch 执行ES搜索
func (this *evApi) EsSearch(ctx context.Context, req dto.SearchReq) (res *proto.Response, err error)

// 使用示例
searchRes, err := api.EsSearch(ctx, dto.SearchReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    SearchReqData: dto.SearchReqData{
        SearchRequest: proto.SearchRequest{
            Index: []string{"my-index"},
            Size:  proto.IntPtr(100),
            From:  proto.IntPtr(0),
        },
        Query: map[string]interface{}{
            "match_all": map[string]interface{}{},
        },
    },
})

// 复杂查询示例
complexSearchRes, err := api.EsSearch(ctx, dto.SearchReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    SearchReqData: dto.SearchReqData{
        SearchRequest: proto.SearchRequest{
            Index: []string{"logs-*"},
            Size:  proto.IntPtr(50),
        },
        Query: map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "range": map[string]interface{}{
                            "@timestamp": map[string]interface{}{
                                "gte": "now-1h",
                            },
                        },
                    },
                    map[string]interface{}{
                        "term": map[string]interface{}{
                            "level": "ERROR",
                        },
                    },
                },
            },
        },
    },
})
```

##### ES索引管理
```go
// 创建索引
res, err := api.EsCreateIndex(ctx, dto.CreateIndexReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    CreateIndexReqData: dto.CreateIndexReqData{
        IndexCreateRequest: proto.IndicesCreateRequest{
            Index: "new-index",
        },
        Body: map[string]interface{}{
            "settings": map[string]interface{}{
                "number_of_shards":   1,
                "number_of_replicas": 0,
            },
            "mappings": map[string]interface{}{
                "properties": map[string]interface{}{
                    "title": map[string]interface{}{
                        "type": "text",
                        "analyzer": "standard",
                    },
                    "content": map[string]interface{}{
                        "type": "text",
                    },
                    "created_at": map[string]interface{}{
                        "type": "date",
                    },
                },
            },
        },
    },
})

// 删除索引
res, err := api.EsDeleteIndex(ctx, dto.DeleteIndexReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DeleteIndexReqData: dto.DeleteIndexReqData{
        IndicesDeleteRequest: proto.IndicesDeleteRequest{
            Index: []string{"old-index"},
        },
    },
})

// 获取索引列表
res, err := api.EsGetIndices(ctx, dto.GetIndicesReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    GetIndicesReqData: dto.GetIndicesReqData{
        CatIndicesRequest: proto.CatIndicesRequest{
            Index:  []string{"*"},
            Format: proto.StringPtr("json"),
            H:      []string{"index", "docs.count", "store.size"},
        },
    },
})
```

##### ES文档操作
```go
// 创建文档
res, err := api.EsCreate(ctx, dto.CreateReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    CreateReqData: dto.CreateReqData{
        CreateRequest: proto.CreateRequest{
            Index:      "my-index",
            DocumentID: "1",
        },
        Body: map[string]interface{}{
            "title":      "示例文档",
            "content":    "这是一个示例文档的内容",
            "created_at": time.Now().Format(time.RFC3339),
        },
    },
})

// 更新文档
res, err := api.EsUpdate(ctx, dto.UpdateReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    UpdateReqData: dto.UpdateReqData{
        UpdateRequest: proto.UpdateRequest{
            Index:      "my-index",
            DocumentID: "1",
        },
        Body: map[string]interface{}{
            "doc": map[string]interface{}{
                "title":      "更新后的标题",
                "updated_at": time.Now().Format(time.RFC3339),
            },
        },
    },
})

// 删除文档
res, err := api.EsDelete(ctx, dto.DeleteReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DeleteReqData: dto.DeleteReqData{
        DeleteRequest: proto.DeleteRequest{
            Index:      "my-index",
            DocumentID: "1",
        },
    },
})
```

##### ES索引操作
```go
// 刷新索引
res, err := api.EsRefresh(ctx, dto.RefreshReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    RefreshReqData: dto.RefreshReqData{
        IndexNames: []string{"my-index"},
    },
})

// 打开索引
res, err := api.EsOpen(ctx, dto.OpenReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    OpenReqData: dto.OpenReqData{
        IndexNames: []string{"my-index"},
    },
})

// 关闭索引
res, err := api.EsIndicesClose(ctx, dto.IndicesCloseReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    IndicesCloseReqData: dto.IndicesCloseReqData{
        IndexNames: []string{"my-index"},
    },
})

// 强制合并
maxSegments := 1
res, err := api.EsIndicesForcemerge(ctx, dto.IndicesForcemergeReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    IndicesForcemergeReqData: dto.IndicesForcemergeReqData{
        IndexNames:     []string{"my-index"},
        MaxNumSegments: &maxSegments,
    },
})
```

##### ES别名管理
```go
// 获取别名
res, err := api.EsGetAliases(ctx, dto.GetAliasesReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    GetAliasesReqData: dto.GetAliasesReqData{
        IndexNames: []string{"my-index"},
    },
})

// 添加别名
res, err := api.EsAddAliases(ctx, dto.AddAliasesReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    AddAliasesReqData: dto.AddAliasesReqData{
        IndexName: []string{"my-index"},
        AliasName: "my-alias",
    },
})

// 移除别名
res, err := api.EsRemoveAliases(ctx, dto.RemoveAliasesReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    RemoveAliasesReqData: dto.RemoveAliasesReqData{
        IndexName: []string{"my-index"},
        AliasName: []string{"my-alias"},
    },
})
```

##### ES快照管理
```go
// 创建快照
res, err := api.EsSnapshotCreate(ctx, dto.SnapshotCreateReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    SnapshotCreateReqData: dto.SnapshotCreateReqData{
        Repository: "my-repo",
        Snapshot:   "snapshot-" + time.Now().Format("20060102-150405"),
        ReqJson: proto.Json{
            "indices":               "my-index",
            "ignore_unavailable":    true,
            "include_global_state":  false,
        },
    },
})

// 恢复快照
res, err := api.EsRestoreSnapshot(ctx, dto.RestoreSnapshotReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    RestoreSnapshotReqData: dto.RestoreSnapshotReqData{
        Repository: "my-repo",
        Snapshot:   "my-snapshot",
        ReqJson: proto.Json{
            "indices":               "my-index",
            "ignore_unavailable":    true,
            "include_global_state":  false,
        },
    },
})

// 删除快照
res, err := api.EsSnapshotDelete(ctx, dto.SnapshotDeleteReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    SnapshotDeleteReqData: dto.SnapshotDeleteReqData{
        Repository: "my-repo",
        Snapshot:   "my-snapshot",
    },
})
```

##### ES任务管理
```go
// 获取任务列表
res, err := api.EsTaskList(ctx, dto.TaskListReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
})

// 取消任务
res, err := api.EsTasksCancel(ctx, dto.TasksCancelReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    TasksCancelReqData: dto.TasksCancelReqData{
        TaskId: "node1:12345",
    },
})
```

#### 10. 第三方数据库操作（操作后台设置的数据源，可支持mysql,clickhouse,oracle,sqlserver,postgressql,MariaDb,达梦）
```go
// MySQL查询
columns, result, err := api.MysqlSelectSql(ctx, &dto.MysqlSelectReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: "my_database",
    Sql:    "SELECT * FROM products WHERE price > ?",
    Args:   []interface{}{100},
})

// MySQL单条查询
result, err := api.MysqlFirstSql(ctx, &dto.MysqlSelectReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: "my_database", 
    Sql:    "SELECT * FROM users WHERE id = ?",
    Args:   []interface{}{userId},
})

// MySQL执行
rowsAffected, err := api.MysqlExecSql(ctx, &dto.MysqlExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: "my_database",
    Sql:    "UPDATE products SET price = price * 1.1 WHERE category = ?",
    Args:   []interface{}{"electronics"},
})
```

#### 11.  第三方Redis操作（操作后台设置的数据源）
```go
// Redis命令执行
data, err := api.RedisExecCommand(ctx, &dto.RedisExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: 0,
    Args:   []interface{}{"GET", "user:123"},
})

// Redis SET操作
data, err := api.RedisExecCommand(ctx, &dto.RedisExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: 0,
    Args:   []interface{}{"SET", "user:123", "张三", "EX", 3600},
})

// Redis Hash操作
data, err := api.RedisExecCommand(ctx, &dto.RedisExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: 0,
    Args:   []interface{}{"HGETALL", "user:profile:123"},
})

// Redis列表操作
data, err := api.RedisExecCommand(ctx, &dto.RedisExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: 0,
    Args:   []interface{}{"LPUSH", "notifications", "新消息"},
})
```

#### 12.  第三方mongo操作（操作后台设置的数据源）
```go
// 执行MongoDB查询命令
result, err := api.ExecMongoCommand(ctx, &dto.MongoExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: "mydb",
    Command: bson.D{
        {"find", "users"},
        {"filter", bson.D{
            {"age", bson.D{{"$gt", 18}}},
            {"status", "active"},
        }},
        {"limit", 100},
    },
    Timeout: 30 * time.Second,
})

// MongoDB聚合查询
result, err := api.ExecMongoCommand(ctx, &dto.MongoExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: "mydb",
    Command: bson.D{
        {"aggregate", "orders"},
        {"pipeline", bson.A{
            bson.D{{"$match", bson.D{{"status", "completed"}}}},
            bson.D{{"$group", bson.D{
                {"_id", "$category"},
                {"total", bson.D{{"$sum", "$amount"}}},
                {"count", bson.D{{"$sum", 1}}},
            }}},
            bson.D{{"$sort", bson.D{{"total", -1}}}},
        }},
        {"cursor", bson.D{}},
    },
    Timeout: 60 * time.Second,
})

// MongoDB插入文档
result, err := api.ExecMongoCommand(ctx, &dto.MongoExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: "mydb",
    Command: bson.D{
        {"insert", "users"},
        {"documents", bson.A{
            bson.D{
                {"name", "张三"},
                {"email", "zhangsan@example.com"},
                {"age", 25},
                {"created_at", time.Now()},
            },
        }},
    },
    Timeout: 10 * time.Second,
})

// MongoDB更新文档
result, err := api.ExecMongoCommand(ctx, &dto.MongoExecReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
    DbName: "mydb",
    Command: bson.D{
        {"update", "users"},
        {"updates", bson.A{
            bson.D{
                {"q", bson.D{{"_id", userObjectID}}},
                {"u", bson.D{{"$set", bson.D{
                    {"last_login", time.Now()},
                    {"status", "online"},
                }}}},
            },
        }},
    },
    Timeout: 10 * time.Second,
})

// 显示数据库列表
dbList, err := api.ShowMongoDbs(ctx, &dto.ShowMongoDbsReq{
    EsConnectData: dto.EsConnectData{
        UserID:    userID,
        EsConnect: connectID,
    },
})
```

