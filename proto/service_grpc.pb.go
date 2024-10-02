// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.28.2
// source: service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ColorService_GetColor_FullMethodName = "/ColorService/GetColor"
)

// ColorServiceClient is the client API for ColorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ColorServiceClient interface {
	GetColor(ctx context.Context, in *ColorRequest, opts ...grpc.CallOption) (*ColorResponse, error)
}

type colorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewColorServiceClient(cc grpc.ClientConnInterface) ColorServiceClient {
	return &colorServiceClient{cc}
}

func (c *colorServiceClient) GetColor(ctx context.Context, in *ColorRequest, opts ...grpc.CallOption) (*ColorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ColorResponse)
	err := c.cc.Invoke(ctx, ColorService_GetColor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ColorServiceServer is the server API for ColorService service.
// All implementations must embed UnimplementedColorServiceServer
// for forward compatibility
type ColorServiceServer interface {
	GetColor(context.Context, *ColorRequest) (*ColorResponse, error)
	mustEmbedUnimplementedColorServiceServer()
}

// UnimplementedColorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedColorServiceServer struct {
}

func (UnimplementedColorServiceServer) GetColor(context.Context, *ColorRequest) (*ColorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetColor not implemented")
}
func (UnimplementedColorServiceServer) mustEmbedUnimplementedColorServiceServer() {}

// UnsafeColorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ColorServiceServer will
// result in compilation errors.
type UnsafeColorServiceServer interface {
	mustEmbedUnimplementedColorServiceServer()
}

func RegisterColorServiceServer(s grpc.ServiceRegistrar, srv ColorServiceServer) {
	s.RegisterService(&ColorService_ServiceDesc, srv)
}

func _ColorService_GetColor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ColorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ColorServiceServer).GetColor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ColorService_GetColor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ColorServiceServer).GetColor(ctx, req.(*ColorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ColorService_ServiceDesc is the grpc.ServiceDesc for ColorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ColorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ColorService",
	HandlerType: (*ColorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetColor",
			Handler:    _ColorService_GetColor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}