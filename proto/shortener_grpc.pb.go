// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: proto/shortener.proto

package proto

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

const (
	Shortener_IndexPage_FullMethodName  = "/shortener.Shortener/IndexPage"
	Shortener_RedirectTo_FullMethodName = "/shortener.Shortener/RedirectTo"
	Shortener_ListURLS_FullMethodName   = "/shortener.Shortener/ListURLS"
	Shortener_DeleteURL_FullMethodName  = "/shortener.Shortener/DeleteURL"
	Shortener_Stats_FullMethodName      = "/shortener.Shortener/Stats"
)

// ShortenerClient is the client API for Shortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortenerClient interface {
	IndexPage(ctx context.Context, in *IndexPageRequest, opts ...grpc.CallOption) (*IndexPageResponse, error)
	RedirectTo(ctx context.Context, in *RedirectToRequest, opts ...grpc.CallOption) (*RedirectToResponse, error)
	ListURLS(ctx context.Context, in *ListURLsRequest, opts ...grpc.CallOption) (*ListURLsResponse, error)
	DeleteURL(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatsResponse, error)
}

type shortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewShortenerClient(cc grpc.ClientConnInterface) ShortenerClient {
	return &shortenerClient{cc}
}

func (c *shortenerClient) IndexPage(ctx context.Context, in *IndexPageRequest, opts ...grpc.CallOption) (*IndexPageResponse, error) {
	out := new(IndexPageResponse)
	err := c.cc.Invoke(ctx, Shortener_IndexPage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) RedirectTo(ctx context.Context, in *RedirectToRequest, opts ...grpc.CallOption) (*RedirectToResponse, error) {
	out := new(RedirectToResponse)
	err := c.cc.Invoke(ctx, Shortener_RedirectTo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) ListURLS(ctx context.Context, in *ListURLsRequest, opts ...grpc.CallOption) (*ListURLsResponse, error) {
	out := new(ListURLsResponse)
	err := c.cc.Invoke(ctx, Shortener_ListURLS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) DeleteURL(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, Shortener_DeleteURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortenerClient) Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, Shortener_Stats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortenerServer is the server API for Shortener service.
// All implementations must embed UnimplementedShortenerServer
// for forward compatibility
type ShortenerServer interface {
	IndexPage(context.Context, *IndexPageRequest) (*IndexPageResponse, error)
	RedirectTo(context.Context, *RedirectToRequest) (*RedirectToResponse, error)
	ListURLS(context.Context, *ListURLsRequest) (*ListURLsResponse, error)
	DeleteURL(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Stats(context.Context, *emptypb.Empty) (*StatsResponse, error)
	mustEmbedUnimplementedShortenerServer()
}

// UnimplementedShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedShortenerServer struct {
}

func (UnimplementedShortenerServer) IndexPage(context.Context, *IndexPageRequest) (*IndexPageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexPage not implemented")
}
func (UnimplementedShortenerServer) RedirectTo(context.Context, *RedirectToRequest) (*RedirectToResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RedirectTo not implemented")
}
func (UnimplementedShortenerServer) ListURLS(context.Context, *ListURLsRequest) (*ListURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListURLS not implemented")
}
func (UnimplementedShortenerServer) DeleteURL(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURL not implemented")
}
func (UnimplementedShortenerServer) Stats(context.Context, *emptypb.Empty) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedShortenerServer) mustEmbedUnimplementedShortenerServer() {}

// UnsafeShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortenerServer will
// result in compilation errors.
type UnsafeShortenerServer interface {
	mustEmbedUnimplementedShortenerServer()
}

func RegisterShortenerServer(s grpc.ServiceRegistrar, srv ShortenerServer) {
	s.RegisterService(&Shortener_ServiceDesc, srv)
}

func _Shortener_IndexPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexPageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).IndexPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_IndexPage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).IndexPage(ctx, req.(*IndexPageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_RedirectTo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RedirectToRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).RedirectTo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_RedirectTo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).RedirectTo(ctx, req.(*RedirectToRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_ListURLS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).ListURLS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_ListURLS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).ListURLS(ctx, req.(*ListURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_DeleteURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).DeleteURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_DeleteURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).DeleteURL(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortener_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortenerServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Shortener_Stats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortenerServer).Stats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortener_ServiceDesc is the grpc.ServiceDesc for Shortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.Shortener",
	HandlerType: (*ShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IndexPage",
			Handler:    _Shortener_IndexPage_Handler,
		},
		{
			MethodName: "RedirectTo",
			Handler:    _Shortener_RedirectTo_Handler,
		},
		{
			MethodName: "ListURLS",
			Handler:    _Shortener_ListURLS_Handler,
		},
		{
			MethodName: "DeleteURL",
			Handler:    _Shortener_DeleteURL_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _Shortener_Stats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shortener.proto",
}
