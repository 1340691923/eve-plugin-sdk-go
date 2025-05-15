// util包用于提供各种实用工具函数
package util

// 导入所需的包
import (
	// 导入eve-plugin-sdk-go的枚举包
	"github.com/1340691923/eve-plugin-sdk-go/enum"
	// 导入gin框架
	"github.com/gin-gonic/gin"
	// 导入字符串转换包
	"strconv"
)

// GetEvUserID 从gin上下文中获取EVE用户ID
func GetEvUserID(ctx *gin.Context) int {
	// 从请求头获取用户ID并转换为整型
	userId, err := strconv.Atoi(ctx.GetHeader(enum.EvUserID))
	// 如果转换过程中出现错误，返回默认值0
	if err != nil {
		return 0
	}
	// 返回获取到的用户ID
	return userId
}

func GetPluginID(ctx *gin.Context) string {
	return ctx.GetHeader(enum.EvFromPluginID)
}
