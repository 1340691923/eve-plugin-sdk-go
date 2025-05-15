// util包用于提供各种实用工具函数
package util

// 导入所需的包
import (
	// 导入高性能JSON处理库
	"github.com/goccy/go-json"
	// 导入错误处理库
	"github.com/pkg/errors"
	// 导入文件操作库
	"io/ioutil"
)

// LoadJsonAndParse 从文件加载JSON并解析到指定结构
func LoadJsonAndParse(f string, res interface{}) error {
	// 读取文件内容
	b, err := ioutil.ReadFile(f)
	// 如果读取失败，返回包装的错误
	if err != nil {
		return errors.Wrap(err, "加载插件配置文件失败")
	}
	// 将JSON数据解析到结构体
	err = json.Unmarshal(b, &res)
	// 如果解析失败，返回包装的错误
	if err != nil {
		return errors.Wrap(err, "解析插件配置文件失败")
	}
	// 操作成功，返回nil
	return nil
}
