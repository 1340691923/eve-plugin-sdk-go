// ev_api包提供EVE API的接口和实现
package ev_api

// 导入错误处理库
import 	"github.com/pkg/errors"

// NoSubscriberErr 定义无订阅者错误
var NoSubscriberErr  = errors.New("此频道订阅数为0")
