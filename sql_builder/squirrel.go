package sql_builder

import (
	"fmt"
	"github.com/Masterminds/squirrel"
)

var SqlBuilder = squirrel.StatementBuilder

type (
	Eq            = squirrel.Eq
	Or            = squirrel.Or
	And           = squirrel.And
	NotEq         = squirrel.NotEq
	Gt            = squirrel.Gt
	Lt            = squirrel.Lt
	GtOrEq        = squirrel.GtOrEq
	LtOrEq        = squirrel.LtOrEq
	Like          = squirrel.Like
	Gte           = squirrel.GtOrEq
	Lte           = squirrel.LtOrEq
	SelectBuilder = squirrel.SelectBuilder
	InsertBuilder = squirrel.InsertBuilder
	UpdateBuilder = squirrel.UpdateBuilder
)

// 创建分页查询
func CreatePage(page, limit int) uint64 {
	tmp := (page - 1) * limit
	return uint64(tmp)
}

// 创建模糊查询
func CreateLike(column string) string {
	return fmt.Sprint("%", column, "%")
}
