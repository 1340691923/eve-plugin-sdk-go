// live包提供实时数据处理功能
package live

// 导入错误处理库
import "github.com/pkg/errors"

// NoSubscriberErr 定义无订阅者错误
var NoSubscriberErr = errors.New("此频道订阅数为0")
