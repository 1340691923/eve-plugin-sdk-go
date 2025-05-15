// proto包包含Protocol Buffers定义和相关代码
package proto

// protoc命令用于生成Go语言的gRPC代码
// 生成的代码会输出到../genproto/pluginv2目录
// 指定不要求实现所有服务接口
//protoc -I ./  --go_out=../genproto/pluginv2  --go-grpc_out=../genproto/pluginv2 --go-grpc_opt=require_unimplemented_servers=false backend.proto
