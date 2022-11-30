// Code generated by protoc-gen-go.
// source: protos/golang.conradwood.net/apis/webcammixer/webcammixer.proto
// DO NOT EDIT!

/*
Package webcammixer is a generated protocol buffer package.

It is generated from these files:
	protos/golang.conradwood.net/apis/webcammixer/webcammixer.proto

It has these top-level messages:
	ImageStream
	VideoDeviceDef
	FrameStream
	LoopbackInfo
	IdleTextRequest
*/
package webcammixer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "golang.conradwood.net/apis/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// break down one or more images into transportable chunks
type ImageStream struct {
	NextImage bool   `protobuf:"varint,1,opt,name=NextImage" json:"NextImage,omitempty"`
	Data      []byte `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (m *ImageStream) Reset()                    { *m = ImageStream{} }
func (m *ImageStream) String() string            { return proto.CompactTextString(m) }
func (*ImageStream) ProtoMessage()               {}
func (*ImageStream) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ImageStream) GetNextImage() bool {
	if m != nil {
		return m.NextImage
	}
	return false
}

func (m *ImageStream) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type VideoDeviceDef struct {
	VideoDeviceName string `protobuf:"bytes,1,opt,name=VideoDeviceName" json:"VideoDeviceName,omitempty"`
}

func (m *VideoDeviceDef) Reset()                    { *m = VideoDeviceDef{} }
func (m *VideoDeviceDef) String() string            { return proto.CompactTextString(m) }
func (*VideoDeviceDef) ProtoMessage()               {}
func (*VideoDeviceDef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *VideoDeviceDef) GetVideoDeviceName() string {
	if m != nil {
		return m.VideoDeviceName
	}
	return ""
}

// break down one or more frames into transportable chunks
type FrameStream struct {
	NextImage bool   `protobuf:"varint,1,opt,name=NextImage" json:"NextImage,omitempty"`
	Data      []byte `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (m *FrameStream) Reset()                    { *m = FrameStream{} }
func (m *FrameStream) String() string            { return proto.CompactTextString(m) }
func (*FrameStream) ProtoMessage()               {}
func (*FrameStream) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *FrameStream) GetNextImage() bool {
	if m != nil {
		return m.NextImage
	}
	return false
}

func (m *FrameStream) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type LoopbackInfo struct {
	DeviceName  string `protobuf:"bytes,1,opt,name=DeviceName" json:"DeviceName,omitempty"`
	Width       uint32 `protobuf:"varint,2,opt,name=Width" json:"Width,omitempty"`
	Height      uint32 `protobuf:"varint,3,opt,name=Height" json:"Height,omitempty"`
	PixelFormat uint32 `protobuf:"varint,4,opt,name=PixelFormat" json:"PixelFormat,omitempty"`
}

func (m *LoopbackInfo) Reset()                    { *m = LoopbackInfo{} }
func (m *LoopbackInfo) String() string            { return proto.CompactTextString(m) }
func (*LoopbackInfo) ProtoMessage()               {}
func (*LoopbackInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LoopbackInfo) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

func (m *LoopbackInfo) GetWidth() uint32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *LoopbackInfo) GetHeight() uint32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *LoopbackInfo) GetPixelFormat() uint32 {
	if m != nil {
		return m.PixelFormat
	}
	return 0
}

type IdleTextRequest struct {
	Text string `protobuf:"bytes,1,opt,name=Text" json:"Text,omitempty"`
}

func (m *IdleTextRequest) Reset()                    { *m = IdleTextRequest{} }
func (m *IdleTextRequest) String() string            { return proto.CompactTextString(m) }
func (*IdleTextRequest) ProtoMessage()               {}
func (*IdleTextRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *IdleTextRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func init() {
	proto.RegisterType((*ImageStream)(nil), "webcammixer.ImageStream")
	proto.RegisterType((*VideoDeviceDef)(nil), "webcammixer.VideoDeviceDef")
	proto.RegisterType((*FrameStream)(nil), "webcammixer.FrameStream")
	proto.RegisterType((*LoopbackInfo)(nil), "webcammixer.LoopbackInfo")
	proto.RegisterType((*IdleTextRequest)(nil), "webcammixer.IdleTextRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for WebCamMixer service

type WebCamMixerClient interface {
	// send one or more images to v4l
	SendImages(ctx context.Context, opts ...grpc.CallOption) (WebCamMixer_SendImagesClient, error)
	//
	// send one or more frames (that is images converted to v4l-out format) to v4l
	// no conversion will be done on this frame.
	// mostly useful for testing converters
	SendFrames(ctx context.Context, opts ...grpc.CallOption) (WebCamMixer_SendFramesClient, error)
	// connect to a video device
	SendVideoDevice(ctx context.Context, in *VideoDeviceDef, opts ...grpc.CallOption) (*common.Void, error)
	// switch to "idle"
	SwitchToIdle(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error)
	// get loopback config info
	GetLoopbackInfo(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*LoopbackInfo, error)
	// set idle text
	SetIdleText(ctx context.Context, in *IdleTextRequest, opts ...grpc.CallOption) (*common.Void, error)
}

type webCamMixerClient struct {
	cc *grpc.ClientConn
}

func NewWebCamMixerClient(cc *grpc.ClientConn) WebCamMixerClient {
	return &webCamMixerClient{cc}
}

func (c *webCamMixerClient) SendImages(ctx context.Context, opts ...grpc.CallOption) (WebCamMixer_SendImagesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WebCamMixer_serviceDesc.Streams[0], c.cc, "/webcammixer.WebCamMixer/SendImages", opts...)
	if err != nil {
		return nil, err
	}
	x := &webCamMixerSendImagesClient{stream}
	return x, nil
}

type WebCamMixer_SendImagesClient interface {
	Send(*ImageStream) error
	CloseAndRecv() (*common.Void, error)
	grpc.ClientStream
}

type webCamMixerSendImagesClient struct {
	grpc.ClientStream
}

func (x *webCamMixerSendImagesClient) Send(m *ImageStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *webCamMixerSendImagesClient) CloseAndRecv() (*common.Void, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(common.Void)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *webCamMixerClient) SendFrames(ctx context.Context, opts ...grpc.CallOption) (WebCamMixer_SendFramesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_WebCamMixer_serviceDesc.Streams[1], c.cc, "/webcammixer.WebCamMixer/SendFrames", opts...)
	if err != nil {
		return nil, err
	}
	x := &webCamMixerSendFramesClient{stream}
	return x, nil
}

type WebCamMixer_SendFramesClient interface {
	Send(*FrameStream) error
	CloseAndRecv() (*common.Void, error)
	grpc.ClientStream
}

type webCamMixerSendFramesClient struct {
	grpc.ClientStream
}

func (x *webCamMixerSendFramesClient) Send(m *FrameStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *webCamMixerSendFramesClient) CloseAndRecv() (*common.Void, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(common.Void)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *webCamMixerClient) SendVideoDevice(ctx context.Context, in *VideoDeviceDef, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/SendVideoDevice", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webCamMixerClient) SwitchToIdle(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/SwitchToIdle", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webCamMixerClient) GetLoopbackInfo(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*LoopbackInfo, error) {
	out := new(LoopbackInfo)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/GetLoopbackInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webCamMixerClient) SetIdleText(ctx context.Context, in *IdleTextRequest, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/SetIdleText", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for WebCamMixer service

type WebCamMixerServer interface {
	// send one or more images to v4l
	SendImages(WebCamMixer_SendImagesServer) error
	//
	// send one or more frames (that is images converted to v4l-out format) to v4l
	// no conversion will be done on this frame.
	// mostly useful for testing converters
	SendFrames(WebCamMixer_SendFramesServer) error
	// connect to a video device
	SendVideoDevice(context.Context, *VideoDeviceDef) (*common.Void, error)
	// switch to "idle"
	SwitchToIdle(context.Context, *common.Void) (*common.Void, error)
	// get loopback config info
	GetLoopbackInfo(context.Context, *common.Void) (*LoopbackInfo, error)
	// set idle text
	SetIdleText(context.Context, *IdleTextRequest) (*common.Void, error)
}

func RegisterWebCamMixerServer(s *grpc.Server, srv WebCamMixerServer) {
	s.RegisterService(&_WebCamMixer_serviceDesc, srv)
}

func _WebCamMixer_SendImages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WebCamMixerServer).SendImages(&webCamMixerSendImagesServer{stream})
}

type WebCamMixer_SendImagesServer interface {
	SendAndClose(*common.Void) error
	Recv() (*ImageStream, error)
	grpc.ServerStream
}

type webCamMixerSendImagesServer struct {
	grpc.ServerStream
}

func (x *webCamMixerSendImagesServer) SendAndClose(m *common.Void) error {
	return x.ServerStream.SendMsg(m)
}

func (x *webCamMixerSendImagesServer) Recv() (*ImageStream, error) {
	m := new(ImageStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _WebCamMixer_SendFrames_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WebCamMixerServer).SendFrames(&webCamMixerSendFramesServer{stream})
}

type WebCamMixer_SendFramesServer interface {
	SendAndClose(*common.Void) error
	Recv() (*FrameStream, error)
	grpc.ServerStream
}

type webCamMixerSendFramesServer struct {
	grpc.ServerStream
}

func (x *webCamMixerSendFramesServer) SendAndClose(m *common.Void) error {
	return x.ServerStream.SendMsg(m)
}

func (x *webCamMixerSendFramesServer) Recv() (*FrameStream, error) {
	m := new(FrameStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _WebCamMixer_SendVideoDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VideoDeviceDef)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).SendVideoDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/SendVideoDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).SendVideoDevice(ctx, req.(*VideoDeviceDef))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebCamMixer_SwitchToIdle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).SwitchToIdle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/SwitchToIdle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).SwitchToIdle(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebCamMixer_GetLoopbackInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).GetLoopbackInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/GetLoopbackInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).GetLoopbackInfo(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebCamMixer_SetIdleText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdleTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).SetIdleText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/SetIdleText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).SetIdleText(ctx, req.(*IdleTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WebCamMixer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "webcammixer.WebCamMixer",
	HandlerType: (*WebCamMixerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendVideoDevice",
			Handler:    _WebCamMixer_SendVideoDevice_Handler,
		},
		{
			MethodName: "SwitchToIdle",
			Handler:    _WebCamMixer_SwitchToIdle_Handler,
		},
		{
			MethodName: "GetLoopbackInfo",
			Handler:    _WebCamMixer_GetLoopbackInfo_Handler,
		},
		{
			MethodName: "SetIdleText",
			Handler:    _WebCamMixer_SetIdleText_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendImages",
			Handler:       _WebCamMixer_SendImages_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SendFrames",
			Handler:       _WebCamMixer_SendFrames_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "protos/golang.conradwood.net/apis/webcammixer/webcammixer.proto",
}

func init() {
	proto.RegisterFile("protos/golang.conradwood.net/apis/webcammixer/webcammixer.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 406 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x52, 0xd1, 0x8a, 0xd3, 0x40,
	0x14, 0x25, 0xb5, 0x16, 0x7b, 0x13, 0x0d, 0x0c, 0x22, 0xb1, 0x16, 0xad, 0x05, 0x21, 0xf8, 0x90,
	0x82, 0x82, 0x60, 0x11, 0x0a, 0x5a, 0xaa, 0x05, 0x2d, 0x92, 0x94, 0xf6, 0x79, 0x92, 0xdc, 0xa6,
	0x83, 0x9d, 0x4c, 0x4d, 0x66, 0xb7, 0x79, 0xda, 0x4f, 0xdd, 0x6f, 0x59, 0x32, 0x4d, 0xd9, 0x99,
	0x16, 0xf6, 0x61, 0x9f, 0x32, 0xf7, 0xcc, 0x3d, 0xf7, 0x64, 0xce, 0xb9, 0x30, 0xd9, 0x17, 0x42,
	0x8a, 0x72, 0x94, 0x89, 0x1d, 0xcd, 0xb3, 0x20, 0x11, 0x79, 0x41, 0xd3, 0x83, 0x10, 0x69, 0x90,
	0xa3, 0x1c, 0xd1, 0x3d, 0x2b, 0x47, 0x07, 0x8c, 0x13, 0xca, 0x39, 0xab, 0xb0, 0xd0, 0xcf, 0x81,
	0x62, 0x12, 0x5b, 0x83, 0x7a, 0xc1, 0x03, 0x63, 0x12, 0xc1, 0xb9, 0xc8, 0x9b, 0xcf, 0x91, 0x3c,
	0x9c, 0x80, 0x3d, 0xe7, 0x34, 0xc3, 0x48, 0x16, 0x48, 0x39, 0xe9, 0x43, 0x77, 0x81, 0x95, 0x54,
	0x90, 0x67, 0x0d, 0x2c, 0xff, 0x59, 0x78, 0x0f, 0x10, 0x02, 0xed, 0x29, 0x95, 0xd4, 0x6b, 0x0d,
	0x2c, 0xdf, 0x09, 0xd5, 0x79, 0x38, 0x86, 0x17, 0x2b, 0x96, 0xa2, 0x98, 0xe2, 0x35, 0x4b, 0x70,
	0x8a, 0x1b, 0xe2, 0x83, 0xab, 0x21, 0x0b, 0xca, 0x8f, 0x93, 0xba, 0xe1, 0x39, 0x5c, 0x8b, 0xcf,
	0x0a, 0xca, 0x1f, 0x2f, 0x7e, 0x03, 0xce, 0x6f, 0x21, 0xf6, 0x31, 0x4d, 0xfe, 0xcd, 0xf3, 0x8d,
	0x20, 0x6f, 0x01, 0x2e, 0x54, 0x35, 0x84, 0xbc, 0x84, 0xa7, 0x6b, 0x96, 0xca, 0xad, 0x1a, 0xf2,
	0x3c, 0x3c, 0x16, 0xe4, 0x15, 0x74, 0x7e, 0x21, 0xcb, 0xb6, 0xd2, 0x7b, 0xa2, 0xe0, 0xa6, 0x22,
	0x03, 0xb0, 0xff, 0xb2, 0x0a, 0x77, 0x33, 0x51, 0x70, 0x2a, 0xbd, 0xb6, 0xba, 0xd4, 0xa1, 0xe1,
	0x07, 0x70, 0xe7, 0xe9, 0x0e, 0x97, 0x58, 0xc9, 0x10, 0xff, 0x5f, 0x61, 0x29, 0xeb, 0xdf, 0xac,
	0xcb, 0x46, 0x5c, 0x9d, 0x3f, 0xdd, 0xb6, 0xc0, 0x5e, 0x63, 0xfc, 0x83, 0xf2, 0x3f, 0x75, 0x48,
	0xe4, 0x0b, 0x40, 0x84, 0x79, 0xaa, 0xde, 0x55, 0x12, 0x2f, 0xd0, 0x33, 0xd5, 0xd2, 0xe8, 0x39,
	0x41, 0x93, 0xd5, 0x4a, 0xb0, 0xd4, 0xb7, 0x4e, 0x3c, 0xe5, 0xd9, 0x39, 0x4f, 0x33, 0xf2, 0x82,
	0xf7, 0x0d, 0xdc, 0x9a, 0xa7, 0xd9, 0x4f, 0xde, 0x18, 0x64, 0x33, 0x41, 0x93, 0x4f, 0x3e, 0x82,
	0x13, 0x1d, 0x98, 0x4c, 0xb6, 0x4b, 0x51, 0x3f, 0x96, 0x18, 0xb7, 0x67, 0xbd, 0x63, 0x70, 0x7f,
	0xa2, 0x34, 0x32, 0x31, 0xdb, 0x5f, 0x1b, 0xba, 0x46, 0xe3, 0x57, 0xb0, 0x23, 0x94, 0x27, 0x3f,
	0x49, 0xdf, 0xb4, 0xc5, 0xb4, 0xd9, 0x94, 0xfd, 0xfe, 0x1e, 0xde, 0xe5, 0x28, 0xf5, 0xa5, 0xaf,
	0x17, 0x5e, 0x1f, 0x10, 0x77, 0xd4, 0xbe, 0x7f, 0xbe, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xb8, 0x8a,
	0x0a, 0x2e, 0x6f, 0x03, 0x00, 0x00,
}
