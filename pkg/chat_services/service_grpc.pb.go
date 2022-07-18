// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: service.proto

package chat_services

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

// ChatServicesClient is the client API for ChatServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServicesClient interface {
	// Get all active rooms.
	GetActiveRooms(ctx context.Context, in *ActiveRoomsRequest, opts ...grpc.CallOption) (*RoomList, error)
	// Get room information based on it.
	GetRoom(ctx context.Context, in *RoomId, opts ...grpc.CallOption) (*Room, error)
	// Request to create a new chatroom.
	CreateRoom(ctx context.Context, in *RoomRequest, opts ...grpc.CallOption) (*Room, error)
	// Bidirectional rpc stream for sending and receiving messages.
	RouteChat(ctx context.Context, opts ...grpc.CallOption) (ChatServices_RouteChatClient, error)
}

type chatServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServicesClient(cc grpc.ClientConnInterface) ChatServicesClient {
	return &chatServicesClient{cc}
}

func (c *chatServicesClient) GetActiveRooms(ctx context.Context, in *ActiveRoomsRequest, opts ...grpc.CallOption) (*RoomList, error) {
	out := new(RoomList)
	err := c.cc.Invoke(ctx, "/main.ChatServices/GetActiveRooms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServicesClient) GetRoom(ctx context.Context, in *RoomId, opts ...grpc.CallOption) (*Room, error) {
	out := new(Room)
	err := c.cc.Invoke(ctx, "/main.ChatServices/GetRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServicesClient) CreateRoom(ctx context.Context, in *RoomRequest, opts ...grpc.CallOption) (*Room, error) {
	out := new(Room)
	err := c.cc.Invoke(ctx, "/main.ChatServices/CreateRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServicesClient) RouteChat(ctx context.Context, opts ...grpc.CallOption) (ChatServices_RouteChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatServices_ServiceDesc.Streams[0], "/main.ChatServices/RouteChat", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServicesRouteChatClient{stream}
	return x, nil
}

type ChatServices_RouteChatClient interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chatServicesRouteChatClient struct {
	grpc.ClientStream
}

func (x *chatServicesRouteChatClient) Send(m *ChatMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServicesRouteChatClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServicesServer is the server API for ChatServices service.
// All implementations must embed UnimplementedChatServicesServer
// for forward compatibility
type ChatServicesServer interface {
	// Get all active rooms.
	GetActiveRooms(context.Context, *ActiveRoomsRequest) (*RoomList, error)
	// Get room information based on it.
	GetRoom(context.Context, *RoomId) (*Room, error)
	// Request to create a new chatroom.
	CreateRoom(context.Context, *RoomRequest) (*Room, error)
	// Bidirectional rpc stream for sending and receiving messages.
	RouteChat(ChatServices_RouteChatServer) error
	mustEmbedUnimplementedChatServicesServer()
}

// UnimplementedChatServicesServer must be embedded to have forward compatible implementations.
type UnimplementedChatServicesServer struct {
}

func (UnimplementedChatServicesServer) GetActiveRooms(context.Context, *ActiveRoomsRequest) (*RoomList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActiveRooms not implemented")
}
func (UnimplementedChatServicesServer) GetRoom(context.Context, *RoomId) (*Room, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoom not implemented")
}
func (UnimplementedChatServicesServer) CreateRoom(context.Context, *RoomRequest) (*Room, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedChatServicesServer) RouteChat(ChatServices_RouteChatServer) error {
	return status.Errorf(codes.Unimplemented, "method RouteChat not implemented")
}
func (UnimplementedChatServicesServer) mustEmbedUnimplementedChatServicesServer() {}

// UnsafeChatServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServicesServer will
// result in compilation errors.
type UnsafeChatServicesServer interface {
	mustEmbedUnimplementedChatServicesServer()
}

func RegisterChatServicesServer(s grpc.ServiceRegistrar, srv ChatServicesServer) {
	s.RegisterService(&ChatServices_ServiceDesc, srv)
}

func _ChatServices_GetActiveRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActiveRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServicesServer).GetActiveRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.ChatServices/GetActiveRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServicesServer).GetActiveRooms(ctx, req.(*ActiveRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatServices_GetRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServicesServer).GetRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.ChatServices/GetRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServicesServer).GetRoom(ctx, req.(*RoomId))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatServices_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServicesServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.ChatServices/CreateRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServicesServer).CreateRoom(ctx, req.(*RoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatServices_RouteChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServicesServer).RouteChat(&chatServicesRouteChatServer{stream})
}

type ChatServices_RouteChatServer interface {
	Send(*ChatMessage) error
	Recv() (*ChatMessage, error)
	grpc.ServerStream
}

type chatServicesRouteChatServer struct {
	grpc.ServerStream
}

func (x *chatServicesRouteChatServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServicesRouteChatServer) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServices_ServiceDesc is the grpc.ServiceDesc for ChatServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.ChatServices",
	HandlerType: (*ChatServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetActiveRooms",
			Handler:    _ChatServices_GetActiveRooms_Handler,
		},
		{
			MethodName: "GetRoom",
			Handler:    _ChatServices_GetRoom_Handler,
		},
		{
			MethodName: "CreateRoom",
			Handler:    _ChatServices_CreateRoom_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RouteChat",
			Handler:       _ChatServices_RouteChat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}
