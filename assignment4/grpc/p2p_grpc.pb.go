// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: grpc/p2p.proto

package ping

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RingClient is the client API for Ring service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RingClient interface {
	HandNext(ctx context.Context, in *MSG, opts ...grpc.CallOption) (*Reply, error)
}

type ringClient struct {
	cc grpc.ClientConnInterface
}

func NewRingClient(cc grpc.ClientConnInterface) RingClient {
	return &ringClient{cc}
}

func (c *ringClient) HandNext(ctx context.Context, in *MSG, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/p2p.Ring/HandNext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RingServer is the server API for Ring service.
// All implementations must embed UnimplementedRingServer
// for forward compatibility
type RingServer interface {
	HandNext(context.Context, *MSG) (*Reply, error)
	mustEmbedUnimplementedRingServer()
}

// UnimplementedRingServer must be embedded to have forward compatible implementations.
type UnimplementedRingServer struct {
}

func (UnimplementedRingServer) HandNext(context.Context, *MSG) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandNext not implemented")
}
func (UnimplementedRingServer) mustEmbedUnimplementedRingServer() {}

// UnsafeRingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RingServer will
// result in compilation errors.
type UnsafeRingServer interface {
	mustEmbedUnimplementedRingServer()
}

func RegisterRingServer(s grpc.ServiceRegistrar, srv RingServer) {
	s.RegisterService(&Ring_ServiceDesc, srv)
}

func _Ring_HandNext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MSG)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingServer).HandNext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2p.Ring/HandNext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingServer).HandNext(ctx, req.(*MSG))
	}
	return interceptor(ctx, in, info, handler)
}

// Ring_ServiceDesc is the grpc.ServiceDesc for Ring service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ring_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "p2p.Ring",
	HandlerType: (*RingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandNext",
			Handler:    _Ring_HandNext_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/p2p.proto",
}
