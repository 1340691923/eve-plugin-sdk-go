// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 导入JSON处理库
	"github.com/goccy/go-json"
	// 导入时间处理库
	"time"
)

// User 定义用户信息结构
type User struct {
	// 用户登录名
	Login string
	// 用户名称
	Name string
	// 用户邮箱
	Email string
	// 用户角色
	Role string
}

// AppInstanceSettings 定义应用实例设置结构
type AppInstanceSettings struct {
	// 应用的JSON配置数据
	JSONData json.RawMessage

	// 已解密的安全JSON数据
	DecryptedSecureJSONData map[string]string

	// 更新时间
	Updated time.Time
}

// DataSourceInstanceSettings 定义数据源实例设置结构
type DataSourceInstanceSettings struct {
	// 数据源ID
	ID int64

	// 数据源唯一标识符
	UID string

	// 数据源类型
	Type string

	// 数据源名称
	Name string

	// 数据源URL
	URL string

	// 数据源用户
	User string

	// 数据源数据库名
	Database string

	// 是否启用基本身份验证
	BasicAuthEnabled bool

	// 基本身份验证用户
	BasicAuthUser string

	// 数据源的JSON配置数据
	JSONData json.RawMessage

	// 已解密的安全JSON数据
	DecryptedSecureJSONData map[string]string

	// 更新时间
	Updated time.Time
}

// PluginContext 定义插件上下文结构
type PluginContext struct {
	// 组织ID
	OrgID int64

	// 插件ID
	PluginID string

	// 用户信息
	User *User

	// 应用实例设置
	AppInstanceSettings *AppInstanceSettings

	// 数据源实例设置
	DataSourceInstanceSettings *DataSourceInstanceSettings
}
