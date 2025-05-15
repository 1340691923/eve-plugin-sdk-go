# ElasticView Golang Plugin SDK

## 这是用于开发ElasticView后端插件的SDK

### 简介
本SDK提供了开发ElasticView后端插件所需的工具和接口，帮助开发者快速构建与ElasticView系统集成的插件。

### 主要功能
- 提供了插件与ElasticView系统通信的标准接口
- 支持实时数据处理
- 提供资源调用功能
- 支持健康检查机制
- 内置SQL构建工具
- 提供便捷的Web服务支持

### 包结构
- `backend`: 提供插件后端的核心功能
- `call_resource`: 提供资源调用和HTTP服务功能
- `check_health`: 提供插件健康检查功能
- `enum`: 定义系统中使用的各种枚举常量
- `live`: 提供实时数据处理功能
- `sql_builder`: 提供SQL构建工具和辅助功能
- `util`: 提供各种实用工具函数

### 使用方法
请参考SDK中的示例代码和各包的文档，了解如何使用本SDK开发ElasticView插件。
