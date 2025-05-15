// sql_builder包提供SQL构建工具和辅助功能
package sql_builder

// 导入所需的包
import (
	// 导入格式化包
	"fmt"
	// 导入squirrel SQL构建库
	"github.com/Masterminds/squirrel"
)

// SqlBuilder 是squirrel的SQL语句构建器实例
var SqlBuilder = squirrel.StatementBuilder

// 类型别名定义，便于使用squirrel库的各种条件和构建器
type (
	// Eq 等于条件别名
	Eq = squirrel.Eq
	// Or 或条件别名
	Or = squirrel.Or
	// And 与条件别名
	And = squirrel.And
	// NotEq 不等于条件别名
	NotEq = squirrel.NotEq
	// Gt 大于条件别名
	Gt = squirrel.Gt
	// Lt 小于条件别名
	Lt = squirrel.Lt
	// GtOrEq 大于等于条件别名
	GtOrEq = squirrel.GtOrEq
	// LtOrEq 小于等于条件别名
	LtOrEq = squirrel.LtOrEq
	// Like 模糊匹配条件别名
	Like = squirrel.Like
	// Gte 大于等于条件别名（与GtOrEq相同）
	Gte = squirrel.GtOrEq
	// Lte 小于等于条件别名（与LtOrEq相同）
	Lte = squirrel.LtOrEq
	// SelectBuilder 查询构建器别名
	SelectBuilder = squirrel.SelectBuilder
	// InsertBuilder 插入构建器别名
	InsertBuilder = squirrel.InsertBuilder
	// UpdateBuilder 更新构建器别名
	UpdateBuilder = squirrel.UpdateBuilder
)

// CreatePage 创建分页查询的偏移量
func CreatePage(page, limit int) uint64 {
	// 计算偏移量：(页码-1) * 每页限制数
	tmp := (page - 1) * limit
	// 转换为无符号整数并返回
	return uint64(tmp)
}

// CreateLike 创建模糊查询的匹配字符串
func CreateLike(column string) string {
	// 在字符串两端添加%符号，用于SQL的LIKE查询
	return fmt.Sprint("%", column, "%")
}
