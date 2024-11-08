// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.28.3
// source: timeline.proto

package timeline

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

// TimelineServiceClient is the client API for TimelineService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimelineServiceClient interface {
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
	TimelineUpdates(ctx context.Context, in *TimelineUpdateRequest, opts ...grpc.CallOption) (TimelineService_TimelineUpdatesClient, error)
	LikePost(ctx context.Context, in *LikePostRequest, opts ...grpc.CallOption) (*Empty, error)
}

type timelineServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTimelineServiceClient(cc grpc.ClientConnInterface) TimelineServiceClient {
	return &timelineServiceClient{cc}
}

func (c *timelineServiceClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, "/timeline.TimelineService/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *timelineServiceClient) TimelineUpdates(ctx context.Context, in *TimelineUpdateRequest, opts ...grpc.CallOption) (TimelineService_TimelineUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &TimelineService_ServiceDesc.Streams[0], "/timeline.TimelineService/TimelineUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &timelineServiceTimelineUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TimelineService_TimelineUpdatesClient interface {
	Recv() (*TimelineUpdateResponse, error)
	grpc.ClientStream
}

type timelineServiceTimelineUpdatesClient struct {
	grpc.ClientStream
}

func (x *timelineServiceTimelineUpdatesClient) Recv() (*TimelineUpdateResponse, error) {
	m := new(TimelineUpdateResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *timelineServiceClient) LikePost(ctx context.Context, in *LikePostRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/timeline.TimelineService/LikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TimelineServiceServer is the server API for TimelineService service.
// All implementations must embed UnimplementedTimelineServiceServer
// for forward compatibility
type TimelineServiceServer interface {
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	TimelineUpdates(*TimelineUpdateRequest, TimelineService_TimelineUpdatesServer) error
	LikePost(context.Context, *LikePostRequest) (*Empty, error)
	mustEmbedUnimplementedTimelineServiceServer()
}

// UnimplementedTimelineServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTimelineServiceServer struct {
}

func (UnimplementedTimelineServiceServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedTimelineServiceServer) TimelineUpdates(*TimelineUpdateRequest, TimelineService_TimelineUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method TimelineUpdates not implemented")
}
func (UnimplementedTimelineServiceServer) LikePost(context.Context, *LikePostRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikePost not implemented")
}
func (UnimplementedTimelineServiceServer) mustEmbedUnimplementedTimelineServiceServer() {}

// UnsafeTimelineServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimelineServiceServer will
// result in compilation errors.
type UnsafeTimelineServiceServer interface {
	mustEmbedUnimplementedTimelineServiceServer()
}

func RegisterTimelineServiceServer(s grpc.ServiceRegistrar, srv TimelineServiceServer) {
	s.RegisterService(&TimelineService_ServiceDesc, srv)
}

func _TimelineService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimelineServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/timeline.TimelineService/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimelineServiceServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TimelineService_TimelineUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TimelineUpdateRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TimelineServiceServer).TimelineUpdates(m, &timelineServiceTimelineUpdatesServer{stream})
}

type TimelineService_TimelineUpdatesServer interface {
	Send(*TimelineUpdateResponse) error
	grpc.ServerStream
}

type timelineServiceTimelineUpdatesServer struct {
	grpc.ServerStream
}

func (x *timelineServiceTimelineUpdatesServer) Send(m *TimelineUpdateResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _TimelineService_LikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimelineServiceServer).LikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/timeline.TimelineService/LikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimelineServiceServer).LikePost(ctx, req.(*LikePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TimelineService_ServiceDesc is the grpc.ServiceDesc for TimelineService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TimelineService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "timeline.TimelineService",
	HandlerType: (*TimelineServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _TimelineService_Search_Handler,
		},
		{
			MethodName: "LikePost",
			Handler:    _TimelineService_LikePost_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TimelineUpdates",
			Handler:       _TimelineService_TimelineUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "timeline.proto",
}
