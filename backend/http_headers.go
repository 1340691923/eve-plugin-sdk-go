// backend包提供插件后端的核心功能
package backend

// 导入所需的包
import (
	// 导入格式化包
	"fmt"
	// 导入HTTP处理包
	"net/http"
	// 导入文本协议处理包
	"net/textproto"
	// 导入字符串处理包
	"strings"
)

// 定义HTTP头常量
const (
	// OAuth身份令牌头名称
	OAuthIdentityTokenHeaderName = "Authorization"

	// OAuth身份ID令牌头名称
	OAuthIdentityIDTokenHeaderName = "X-Id-Token"

	// Cookies头名称
	CookiesHeaderName = "Cookie"

	// HTTP头前缀
	httpHeaderPrefix = "http_"
)

// ForwardHTTPHeaders 定义HTTP头转发接口
type ForwardHTTPHeaders interface {
	// 设置HTTP头
	SetHTTPHeader(key, value string)

	// 删除HTTP头
	DeleteHTTPHeader(key string)

	// 获取指定HTTP头
	GetHTTPHeader(key string) string

	// 获取所有HTTP头
	GetHTTPHeaders() http.Header
}

// setHTTPHeaderInStringMap 在字符串映射中设置HTTP头
func setHTTPHeaderInStringMap(headers map[string]string, key string, value string) {
	// 如果headers为空，则初始化一个新的map
	if headers == nil {
		headers = map[string]string{}
	}

	// 设置带有前缀的HTTP头
	headers[fmt.Sprintf("%s%s", httpHeaderPrefix, key)] = value
}

// getHTTPHeadersFromStringMap 从字符串映射中获取HTTP头
func getHTTPHeadersFromStringMap(headers map[string]string) http.Header {
	// 创建一个新的HTTP头对象
	httpHeaders := http.Header{}

	// 遍历所有头信息
	for k, v := range headers {
		// 检查是否为OAuth身份令牌头
		if textproto.CanonicalMIMEHeaderKey(k) == OAuthIdentityTokenHeaderName {
			httpHeaders.Set(k, v)
		}

		// 检查是否为OAuth身份ID令牌头
		if textproto.CanonicalMIMEHeaderKey(k) == OAuthIdentityIDTokenHeaderName {
			httpHeaders.Set(k, v)
		}

		// 检查是否为Cookies头
		if textproto.CanonicalMIMEHeaderKey(k) == CookiesHeaderName {
			httpHeaders.Set(k, v)
		}

		// 检查是否为带前缀的HTTP头
		if strings.HasPrefix(k, httpHeaderPrefix) {
			// 去除前缀
			hKey := strings.TrimPrefix(k, httpHeaderPrefix)
			// 设置HTTP头
			httpHeaders.Set(hKey, v)
		}
	}

	// 返回HTTP头对象
	return httpHeaders
}

// deleteHTTPHeaderInStringMap 从字符串映射中删除HTTP头
func deleteHTTPHeaderInStringMap(headers map[string]string, key string) {
	// 遍历所有头信息
	for k := range headers {
		// 检查是否匹配要删除的键（包括原始键和带前缀的键）
		if textproto.CanonicalMIMEHeaderKey(k) == textproto.CanonicalMIMEHeaderKey(key) ||
			textproto.CanonicalMIMEHeaderKey(k) == textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", httpHeaderPrefix, key)) {
			// 从映射中删除匹配的键
			delete(headers, k)
			// 找到并删除后跳出循环
			break
		}
	}
}
