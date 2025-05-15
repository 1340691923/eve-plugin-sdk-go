// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on gopkg.in/mgo.v2/bson by Gustavo Niemeyer
// See THIRD-PARTY-NOTICES for original license terms.

// bson包提供了与MongoDB BSON格式交互的类型和函数
// 本包是对MongoDB官方驱动中bson包的简化封装，提供了基本的BSON类型
package bson // import "go.mongodb.org/mongo-driver/bson"

// 导入原始类型包
import "github.com/1340691923/eve-plugin-sdk-go/ev_api/primitive"

// Zeroer 允许自定义结构类型实现零值状态报告
// 所有未实现Zeroer接口或IsZero返回false的结构类型被认为是非零值
type Zeroer interface {
	// IsZero 返回该类型是否处于零值状态
	IsZero() bool
}

// D 是BSON文档的有序表示
// 当元素顺序很重要时应使用此类型，例如MongoDB命令文档
// 如果元素顺序不重要，应该使用M类型
//
// # D不应构造具有重复键名的文档，因为这可能导致未定义的服务器行为
//
// 示例用法:
//
//	bson.D{{"foo", "bar"}, {"hello", "world"}, {"pi", 3.14159}}
type D = primitive.D

// E 表示D中的BSON元素
// 通常在D内部使用
type E = primitive.E

// M 是BSON文档的无序表示
// 当元素顺序不重要时应使用此类型
// 编码和解码时此类型被视为普通的map[string]interface{}
// 元素将按未定义的随机顺序序列化
// 如果元素顺序很重要，应该使用D类型
//
// 示例用法:
//
//	bson.M{"foo": "bar", "hello": "world", "pi": 3.14159}
type M = primitive.M

// A 是BSON数组的有序表示
//
// 示例用法:
//
//	bson.A{"bar", "world", 3.14159, bson.D{{"qux", 12345}}}
type A = primitive.A
