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
	URL
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

type URL struct {
	URL string `protobuf:"bytes,1,opt,name=URL" json:"URL,omitempty"`
}

func (m *URL) Reset()                    { *m = URL{} }
func (m *URL) String() string            { return proto.CompactTextString(m) }
func (*URL) ProtoMessage()               {}
func (*URL) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *URL) GetURL() string {
	if m != nil {
		return m.URL
	}
	return ""
}

func init() {
	proto.RegisterType((*ImageStream)(nil), "webcammixer.ImageStream")
	proto.RegisterType((*VideoDeviceDef)(nil), "webcammixer.VideoDeviceDef")
	proto.RegisterType((*FrameStream)(nil), "webcammixer.FrameStream")
	proto.RegisterType((*LoopbackInfo)(nil), "webcammixer.LoopbackInfo")
	proto.RegisterType((*IdleTextRequest)(nil), "webcammixer.IdleTextRequest")
	proto.RegisterType((*URL)(nil), "webcammixer.URL")
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
	// switch to images-client
	SwitchToLiveImages(ctx context.Context, in *URL, opts ...grpc.CallOption) (*common.Void, error)
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

func (c *webCamMixerClient) SwitchToLiveImages(ctx context.Context, in *URL, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/SwitchToLiveImages", in, out, c.cc, opts...)
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
	// switch to images-client
	SwitchToLiveImages(context.Context, *URL) (*common.Void, error)
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

func _WebCamMixer_SwitchToLiveImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(URL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).SwitchToLiveImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/SwitchToLiveImages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).SwitchToLiveImages(ctx, req.(*URL))
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
		{
			MethodName: "SwitchToLiveImages",
			Handler:    _WebCamMixer_SwitchToLiveImages_Handler,
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
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x53, 0x5d, 0x8b, 0xd3, 0x40,
	0x14, 0xa5, 0x76, 0x5d, 0xdc, 0x9b, 0x6a, 0x97, 0x8b, 0x68, 0x5c, 0x17, 0xad, 0x01, 0x21, 0xf8,
	0x90, 0x82, 0x8a, 0xe0, 0x22, 0x2c, 0x68, 0x59, 0x2d, 0xc4, 0x45, 0x92, 0xfd, 0x78, 0x9e, 0x24,
	0x77, 0xd3, 0xc1, 0x4e, 0xa6, 0x26, 0xe3, 0x36, 0x4f, 0xfe, 0x0a, 0x7f, 0xb0, 0xcc, 0x34, 0xc5,
	0x99, 0x14, 0x7c, 0xf0, 0x29, 0x77, 0xce, 0xdc, 0x73, 0xcf, 0xe4, 0x1c, 0x2e, 0x9c, 0xae, 0x6a,
	0xa9, 0x64, 0x33, 0x2d, 0xe5, 0x92, 0x55, 0x65, 0x94, 0xcb, 0xaa, 0x66, 0xc5, 0x5a, 0xca, 0x22,
	0xaa, 0x48, 0x4d, 0xd9, 0x8a, 0x37, 0xd3, 0x35, 0x65, 0x39, 0x13, 0x82, 0xb7, 0x54, 0xdb, 0x75,
	0x64, 0x98, 0xe8, 0x59, 0xd0, 0x51, 0xf4, 0x8f, 0x31, 0xb9, 0x14, 0x42, 0x56, 0xdd, 0x67, 0x43,
	0x0e, 0x4e, 0xc1, 0x9b, 0x0b, 0x56, 0x52, 0xaa, 0x6a, 0x62, 0x02, 0x8f, 0xe1, 0xe0, 0x9c, 0x5a,
	0x65, 0x20, 0x7f, 0x30, 0x19, 0x84, 0xf7, 0x92, 0xbf, 0x00, 0x22, 0xec, 0xcd, 0x98, 0x62, 0xfe,
	0x9d, 0xc9, 0x20, 0x1c, 0x25, 0xa6, 0x0e, 0x4e, 0xe0, 0xc1, 0x15, 0x2f, 0x48, 0xce, 0xe8, 0x96,
	0xe7, 0x34, 0xa3, 0x1b, 0x0c, 0x61, 0x6c, 0x21, 0xe7, 0x4c, 0x6c, 0x26, 0x1d, 0x24, 0x7d, 0x58,
	0x8b, 0x9f, 0xd5, 0x4c, 0xfc, 0xbf, 0xf8, 0x2f, 0x18, 0xc5, 0x52, 0xae, 0x32, 0x96, 0x7f, 0x9f,
	0x57, 0x37, 0x12, 0x9f, 0x01, 0xec, 0xa8, 0x5a, 0x08, 0x3e, 0x84, 0xbb, 0xd7, 0xbc, 0x50, 0x0b,
	0x33, 0xe4, 0x7e, 0xb2, 0x39, 0xe0, 0x23, 0xd8, 0xff, 0x42, 0xbc, 0x5c, 0x28, 0x7f, 0x68, 0xe0,
	0xee, 0x84, 0x13, 0xf0, 0xbe, 0xf1, 0x96, 0x96, 0x67, 0xb2, 0x16, 0x4c, 0xf9, 0x7b, 0xe6, 0xd2,
	0x86, 0x82, 0x97, 0x30, 0x9e, 0x17, 0x4b, 0xba, 0xa0, 0x56, 0x25, 0xf4, 0xe3, 0x27, 0x35, 0x4a,
	0x3f, 0x53, 0x1f, 0x3b, 0x71, 0x53, 0x07, 0x8f, 0x61, 0x78, 0x99, 0xc4, 0x78, 0x68, 0x3e, 0xdd,
	0x8d, 0x2e, 0x5f, 0xff, 0x1e, 0x82, 0x77, 0x4d, 0xd9, 0x27, 0x26, 0xbe, 0xea, 0xf4, 0xf0, 0x1d,
	0x40, 0x4a, 0x55, 0x61, 0x7e, 0xb8, 0x41, 0x3f, 0xb2, 0xc3, 0xb6, 0x62, 0x3a, 0x1a, 0x45, 0x5d,
	0x88, 0x57, 0x92, 0x17, 0xe1, 0x60, 0xcb, 0x33, 0x66, 0xf6, 0x79, 0x96, 0xc3, 0x3b, 0xbc, 0x0f,
	0x30, 0xd6, 0x3c, 0x2b, 0x17, 0x7c, 0xea, 0x90, 0xdd, 0x68, 0x5d, 0x3e, 0xbe, 0x82, 0x51, 0xba,
	0xe6, 0x2a, 0x5f, 0x5c, 0x48, 0xed, 0x02, 0x3a, 0xb7, 0xbd, 0xde, 0x13, 0x18, 0x7f, 0x26, 0xe5,
	0x84, 0xe5, 0xb6, 0x3f, 0x71, 0x74, 0x9d, 0xc6, 0xf7, 0xe0, 0xa5, 0xa4, 0xb6, 0x46, 0xe3, 0xb1,
	0x6b, 0x8b, 0xeb, 0x7f, 0x4f, 0xf6, 0x2d, 0xe0, 0xf6, 0x89, 0x31, 0xbf, 0xa5, 0xce, 0xd8, 0x43,
	0x67, 0xc2, 0x65, 0x12, 0xbb, 0xac, 0x8f, 0x2f, 0xe0, 0x79, 0x45, 0xca, 0xde, 0x21, 0xbd, 0x3f,
	0x36, 0x29, 0xdb, 0x37, 0xeb, 0xf3, 0xe6, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6d, 0x1f, 0xe2,
	0x55, 0xbe, 0x03, 0x00, 0x00,
}
