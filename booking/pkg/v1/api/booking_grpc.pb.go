// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: booking.proto

package api

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

// BookingServiceClient is the client API for BookingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookingServiceClient interface {
	CreateBooking(ctx context.Context, in *CreateBookingRequest, opts ...grpc.CallOption) (*CreateBookingResponse, error)
	GetBookings(ctx context.Context, in *GetOrderBookingsRequest, opts ...grpc.CallOption) (*GetOrderBookingsResponse, error)
	ExpireBookings(ctx context.Context, in *ExpireBookingsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteOrderBookings(ctx context.Context, in *DeleteOrderBookingsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateStock(ctx context.Context, in *CreateStockRequest, opts ...grpc.CallOption) (*CreateStockResponse, error)
	GetStocks(ctx context.Context, in *GetStocksRequest, opts ...grpc.CallOption) (*GetStocksResponse, error)
	GetStock(ctx context.Context, in *GetStockRequest, opts ...grpc.CallOption) (*GetStockResponse, error)
	DeleteStock(ctx context.Context, in *DeleteStockRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type bookingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBookingServiceClient(cc grpc.ClientConnInterface) BookingServiceClient {
	return &bookingServiceClient{cc}
}

func (c *bookingServiceClient) CreateBooking(ctx context.Context, in *CreateBookingRequest, opts ...grpc.CallOption) (*CreateBookingResponse, error) {
	out := new(CreateBookingResponse)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/CreateBooking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) GetBookings(ctx context.Context, in *GetOrderBookingsRequest, opts ...grpc.CallOption) (*GetOrderBookingsResponse, error) {
	out := new(GetOrderBookingsResponse)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/GetBookings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) ExpireBookings(ctx context.Context, in *ExpireBookingsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/ExpireBookings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) DeleteOrderBookings(ctx context.Context, in *DeleteOrderBookingsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/DeleteOrderBookings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) CreateStock(ctx context.Context, in *CreateStockRequest, opts ...grpc.CallOption) (*CreateStockResponse, error) {
	out := new(CreateStockResponse)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/CreateStock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) GetStocks(ctx context.Context, in *GetStocksRequest, opts ...grpc.CallOption) (*GetStocksResponse, error) {
	out := new(GetStocksResponse)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/GetStocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) GetStock(ctx context.Context, in *GetStockRequest, opts ...grpc.CallOption) (*GetStockResponse, error) {
	out := new(GetStockResponse)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/GetStock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookingServiceClient) DeleteStock(ctx context.Context, in *DeleteStockRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/booking.api.v1.BookingService/DeleteStock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookingServiceServer is the server API for BookingService service.
// All implementations must embed UnimplementedBookingServiceServer
// for forward compatibility
type BookingServiceServer interface {
	CreateBooking(context.Context, *CreateBookingRequest) (*CreateBookingResponse, error)
	GetBookings(context.Context, *GetOrderBookingsRequest) (*GetOrderBookingsResponse, error)
	ExpireBookings(context.Context, *ExpireBookingsRequest) (*emptypb.Empty, error)
	DeleteOrderBookings(context.Context, *DeleteOrderBookingsRequest) (*emptypb.Empty, error)
	CreateStock(context.Context, *CreateStockRequest) (*CreateStockResponse, error)
	GetStocks(context.Context, *GetStocksRequest) (*GetStocksResponse, error)
	GetStock(context.Context, *GetStockRequest) (*GetStockResponse, error)
	DeleteStock(context.Context, *DeleteStockRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedBookingServiceServer()
}

// UnimplementedBookingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBookingServiceServer struct {
}

func (UnimplementedBookingServiceServer) CreateBooking(context.Context, *CreateBookingRequest) (*CreateBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBooking not implemented")
}
func (UnimplementedBookingServiceServer) GetBookings(context.Context, *GetOrderBookingsRequest) (*GetOrderBookingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBookings not implemented")
}
func (UnimplementedBookingServiceServer) ExpireBookings(context.Context, *ExpireBookingsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExpireBookings not implemented")
}
func (UnimplementedBookingServiceServer) DeleteOrderBookings(context.Context, *DeleteOrderBookingsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteOrderBookings not implemented")
}
func (UnimplementedBookingServiceServer) CreateStock(context.Context, *CreateStockRequest) (*CreateStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStock not implemented")
}
func (UnimplementedBookingServiceServer) GetStocks(context.Context, *GetStocksRequest) (*GetStocksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStocks not implemented")
}
func (UnimplementedBookingServiceServer) GetStock(context.Context, *GetStockRequest) (*GetStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStock not implemented")
}
func (UnimplementedBookingServiceServer) DeleteStock(context.Context, *DeleteStockRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStock not implemented")
}
func (UnimplementedBookingServiceServer) mustEmbedUnimplementedBookingServiceServer() {}

// UnsafeBookingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookingServiceServer will
// result in compilation errors.
type UnsafeBookingServiceServer interface {
	mustEmbedUnimplementedBookingServiceServer()
}

func RegisterBookingServiceServer(s grpc.ServiceRegistrar, srv BookingServiceServer) {
	s.RegisterService(&BookingService_ServiceDesc, srv)
}

func _BookingService_CreateBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).CreateBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/CreateBooking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).CreateBooking(ctx, req.(*CreateBookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_GetBookings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrderBookingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).GetBookings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/GetBookings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).GetBookings(ctx, req.(*GetOrderBookingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_ExpireBookings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExpireBookingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).ExpireBookings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/ExpireBookings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).ExpireBookings(ctx, req.(*ExpireBookingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_DeleteOrderBookings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteOrderBookingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).DeleteOrderBookings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/DeleteOrderBookings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).DeleteOrderBookings(ctx, req.(*DeleteOrderBookingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_CreateStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).CreateStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/CreateStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).CreateStock(ctx, req.(*CreateStockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_GetStocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStocksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).GetStocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/GetStocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).GetStocks(ctx, req.(*GetStocksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_GetStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).GetStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/GetStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).GetStock(ctx, req.(*GetStockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BookingService_DeleteStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookingServiceServer).DeleteStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/booking.api.v1.BookingService/DeleteStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookingServiceServer).DeleteStock(ctx, req.(*DeleteStockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BookingService_ServiceDesc is the grpc.ServiceDesc for BookingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BookingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "booking.api.v1.BookingService",
	HandlerType: (*BookingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBooking",
			Handler:    _BookingService_CreateBooking_Handler,
		},
		{
			MethodName: "GetBookings",
			Handler:    _BookingService_GetBookings_Handler,
		},
		{
			MethodName: "ExpireBookings",
			Handler:    _BookingService_ExpireBookings_Handler,
		},
		{
			MethodName: "DeleteOrderBookings",
			Handler:    _BookingService_DeleteOrderBookings_Handler,
		},
		{
			MethodName: "CreateStock",
			Handler:    _BookingService_CreateStock_Handler,
		},
		{
			MethodName: "GetStocks",
			Handler:    _BookingService_GetStocks_Handler,
		},
		{
			MethodName: "GetStock",
			Handler:    _BookingService_GetStock_Handler,
		},
		{
			MethodName: "DeleteStock",
			Handler:    _BookingService_DeleteStock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "booking.proto",
}
