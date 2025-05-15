// util包用于提供各种实用工具函数
package util

// 导入版本比较的依赖包
import "github.com/hashicorp/go-version"

// LessThan 比较两个版本字符串，判断v1是否小于v2
func LessThan(v1, v2 string) bool {
	// 将第一个版本字符串解析为版本对象
	ver1, err1 := version.NewVersion(v1)
	// 如果解析出错，返回false
	if err1 != nil {
		return false
	}
	// 将第二个版本字符串解析为版本对象
	ver2, err2 := version.NewVersion(v2)
	// 如果解析出错，返回false
	if err2 != nil {
		return false
	}

	// 返回版本比较结果，判断ver1是否小于ver2
	return ver1.LessThan(ver2)
}
