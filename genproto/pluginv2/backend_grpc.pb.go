package pluginv2

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion7

const (
	Resource_CallResource_FullMethodName = "/pluginv2.Resource/CallResource"
)

type ResourceClient interface {
	CallResource(ctx context.Context, in *CallResourceRequest, opts ...grpc.CallOption) (Resource_CallResourceClient, error)
}

type resourceClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceClient(cc grpc.ClientConnInterface) ResourceClient {
	return &resourceClient{cc}
}

func (c *resourceClient) CallResource(ctx context.Context, in *CallResourceRequest, opts ...grpc.CallOption) (Resource_CallResourceClient, error) {
	stream, err := c.cc.NewStream(ctx, &Resource_ServiceDesc.Streams[0], Resource_CallResource_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &resourceCallResourceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Resource_CallResourceClient interface {
	Recv() (*CallResourceResponse, error)
	grpc.ClientStream
}

type resourceCallResourceClient struct {
	grpc.ClientStream
}

func (x *resourceCallResourceClient) Recv() (*CallResourceResponse, error) {
	m := new(CallResourceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type ResourceServer interface {
	CallResource(*CallResourceRequest, Resource_CallResourceServer) error
}

type UnimplementedResourceServer struct {
}

func (UnimplementedResourceServer) CallResource(*CallResourceRequest, Resource_CallResourceServer) error {
	return status.Errorf(codes.Unimplemented, "method CallResource not implemented")
}

type UnsafeResourceServer interface {
	mustEmbedUnimplementedResourceServer()
}

func RegisterResourceServer(s grpc.ServiceRegistrar, srv ResourceServer) {
	s.RegisterService(&Resource_ServiceDesc, srv)
}

func _Resource_CallResource_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CallResourceRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ResourceServer).CallResource(m, &resourceCallResourceServer{stream})
}

type Resource_CallResourceServer interface {
	Send(*CallResourceResponse) error
	grpc.ServerStream
}

type resourceCallResourceServer struct {
	grpc.ServerStream
}

func (x *resourceCallResourceServer) Send(m *CallResourceResponse) error {
	return x.ServerStream.SendMsg(m)
}

var Resource_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pluginv2.Resource",
	HandlerType: (*ResourceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CallResource",
			Handler:       _Resource_CallResource_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "backend.proto",
}

const (
	PluginInfo_Get_FullMethodName         = "/pluginv2.PluginInfo/Get"
	PluginInfo_CheckHealth_FullMethodName = "/pluginv2.PluginInfo/CheckHealth"
)

type PluginInfoClient interface {
	Get(ctx context.Context, in *PluginInfoGetReq, opts ...grpc.CallOption) (*PluginInfoGetRes, error)
	CheckHealth(ctx context.Context, in *CheckHealthRequest, opts ...grpc.CallOption) (*CheckHealthResponse, error)
}

type pluginInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginInfoClient(cc grpc.ClientConnInterface) PluginInfoClient {
	return &pluginInfoClient{cc}
}

func (c *pluginInfoClient) Get(ctx context.Context, in *PluginInfoGetReq, opts ...grpc.CallOption) (*PluginInfoGetRes, error) {
	out := new(PluginInfoGetRes)
	err := c.cc.Invoke(ctx, PluginInfo_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginInfoClient) CheckHealth(ctx context.Context, in *CheckHealthRequest, opts ...grpc.CallOption) (*CheckHealthResponse, error) {
	out := new(CheckHealthResponse)
	err := c.cc.Invoke(ctx, PluginInfo_CheckHealth_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type PluginInfoServer interface {
	Get(context.Context, *PluginInfoGetReq) (*PluginInfoGetRes, error)
	CheckHealth(context.Context, *CheckHealthRequest) (*CheckHealthResponse, error)
}

type UnimplementedPluginInfoServer struct {
}

func (UnimplementedPluginInfoServer) Get(context.Context, *PluginInfoGetReq) (*PluginInfoGetRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedPluginInfoServer) CheckHealth(context.Context, *CheckHealthRequest) (*CheckHealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHealth not implemented")
}

type UnsafePluginInfoServer interface {
	mustEmbedUnimplementedPluginInfoServer()
}

func RegisterPluginInfoServer(s grpc.ServiceRegistrar, srv PluginInfoServer) {
	s.RegisterService(&PluginInfo_ServiceDesc, srv)
}

func _PluginInfo_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PluginInfoGetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInfoServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginInfo_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInfoServer).Get(ctx, req.(*PluginInfoGetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginInfo_CheckHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckHealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginInfoServer).CheckHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PluginInfo_CheckHealth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginInfoServer).CheckHealth(ctx, req.(*CheckHealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var PluginInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pluginv2.PluginInfo",
	HandlerType: (*PluginInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _PluginInfo_Get_Handler,
		},
		{
			MethodName: "CheckHealth",
			Handler:    _PluginInfo_CheckHealth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backend.proto",
}
