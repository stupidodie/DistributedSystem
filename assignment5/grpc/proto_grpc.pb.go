// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: grpc/proto.proto

package as5

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TradeClient is the client API for Trade service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TradeClient interface {
	Bid(ctx context.Context, in *Price, opts ...grpc.CallOption) (*Ack, error)
	Result(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Outcome, error)
}

type tradeClient struct {
	cc grpc.ClientConnInterface
}

func NewTradeClient(cc grpc.ClientConnInterface) TradeClient {
	return &tradeClient{cc}
}

func (c *tradeClient) Bid(ctx context.Context, in *Price, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/as5.Trade/Bid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tradeClient) Result(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Outcome, error) {
	out := new(Outcome)
	err := c.cc.Invoke(ctx, "/as5.Trade/Result", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TradeServer is the server API for Trade service.
// All implementations must embed UnimplementedTradeServer
// for forward compatibility
type TradeServer interface {
	Bid(context.Context, *Price) (*Ack, error)
	Result(context.Context, *emptypb.Empty) (*Outcome, error)
	mustEmbedUnimplementedTradeServer()
}

// UnimplementedTradeServer must be embedded to have forward compatible implementations.
type UnimplementedTradeServer struct {
}

func (UnimplementedTradeServer) Bid(context.Context, *Price) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedTradeServer) Result(context.Context, *emptypb.Empty) (*Outcome, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedTradeServer) mustEmbedUnimplementedTradeServer() {}

// UnsafeTradeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TradeServer will
// result in compilation errors.
type UnsafeTradeServer interface {
	mustEmbedUnimplementedTradeServer()
}

func RegisterTradeServer(s grpc.ServiceRegistrar, srv TradeServer) {
	s.RegisterService(&Trade_ServiceDesc, srv)
}

func _Trade_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Price)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TradeServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/as5.Trade/Bid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TradeServer).Bid(ctx, req.(*Price))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trade_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TradeServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/as5.Trade/Result",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TradeServer).Result(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Trade_ServiceDesc is the grpc.ServiceDesc for Trade service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Trade_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "as5.Trade",
	HandlerType: (*TradeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bid",
			Handler:    _Trade_Bid_Handler,
		},
		{
			MethodName: "Result",
			Handler:    _Trade_Result_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/proto.proto",
}
