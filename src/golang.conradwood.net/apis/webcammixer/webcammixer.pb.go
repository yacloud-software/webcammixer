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
	CaptureDeviceList
	CaptureDevice
	CountdownRequest
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

type CaptureDeviceList struct {
	Devices []*CaptureDevice `protobuf:"bytes,1,rep,name=Devices" json:"Devices,omitempty"`
}

func (m *CaptureDeviceList) Reset()                    { *m = CaptureDeviceList{} }
func (m *CaptureDeviceList) String() string            { return proto.CompactTextString(m) }
func (*CaptureDeviceList) ProtoMessage()               {}
func (*CaptureDeviceList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CaptureDeviceList) GetDevices() []*CaptureDevice {
	if m != nil {
		return m.Devices
	}
	return nil
}

type CaptureDevice struct {
	Device string `protobuf:"bytes,1,opt,name=Device" json:"Device,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	Type   uint32 `protobuf:"varint,3,opt,name=Type" json:"Type,omitempty"`
}

func (m *CaptureDevice) Reset()                    { *m = CaptureDevice{} }
func (m *CaptureDevice) String() string            { return proto.CompactTextString(m) }
func (*CaptureDevice) ProtoMessage()               {}
func (*CaptureDevice) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CaptureDevice) GetDevice() string {
	if m != nil {
		return m.Device
	}
	return ""
}

func (m *CaptureDevice) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CaptureDevice) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

type CountdownRequest struct {
	Text    string `protobuf:"bytes,1,opt,name=Text" json:"Text,omitempty"`
	Seconds uint32 `protobuf:"varint,2,opt,name=Seconds" json:"Seconds,omitempty"`
}

func (m *CountdownRequest) Reset()                    { *m = CountdownRequest{} }
func (m *CountdownRequest) String() string            { return proto.CompactTextString(m) }
func (*CountdownRequest) ProtoMessage()               {}
func (*CountdownRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *CountdownRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *CountdownRequest) GetSeconds() uint32 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func init() {
	proto.RegisterType((*ImageStream)(nil), "webcammixer.ImageStream")
	proto.RegisterType((*VideoDeviceDef)(nil), "webcammixer.VideoDeviceDef")
	proto.RegisterType((*FrameStream)(nil), "webcammixer.FrameStream")
	proto.RegisterType((*LoopbackInfo)(nil), "webcammixer.LoopbackInfo")
	proto.RegisterType((*IdleTextRequest)(nil), "webcammixer.IdleTextRequest")
	proto.RegisterType((*URL)(nil), "webcammixer.URL")
	proto.RegisterType((*CaptureDeviceList)(nil), "webcammixer.CaptureDeviceList")
	proto.RegisterType((*CaptureDevice)(nil), "webcammixer.CaptureDevice")
	proto.RegisterType((*CountdownRequest)(nil), "webcammixer.CountdownRequest")
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
	// connect to a LOCAL video device
	SendVideoDevice(ctx context.Context, in *VideoDeviceDef, opts ...grpc.CallOption) (*common.Void, error)
	// connect to any capture source
	SendFromCaptureDevice(ctx context.Context, in *CaptureDevice, opts ...grpc.CallOption) (*common.Void, error)
	// switch to "idle"
	SwitchToIdle(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error)
	// get loopback config info
	GetLoopbackInfo(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*LoopbackInfo, error)
	// set idle text
	SetIdleText(ctx context.Context, in *IdleTextRequest, opts ...grpc.CallOption) (*common.Void, error)
	// switch to images-client
	SwitchToLiveImages(ctx context.Context, in *URL, opts ...grpc.CallOption) (*common.Void, error)
	// get list of video devices
	GetCaptureDevices(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*CaptureDeviceList, error)
	// rpc display a message and countdown
	SetCountdown(ctx context.Context, in *CountdownRequest, opts ...grpc.CallOption) (*common.Void, error)
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

func (c *webCamMixerClient) SendFromCaptureDevice(ctx context.Context, in *CaptureDevice, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/SendFromCaptureDevice", in, out, c.cc, opts...)
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

func (c *webCamMixerClient) GetCaptureDevices(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*CaptureDeviceList, error) {
	out := new(CaptureDeviceList)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/GetCaptureDevices", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webCamMixerClient) SetCountdown(ctx context.Context, in *CountdownRequest, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/webcammixer.WebCamMixer/SetCountdown", in, out, c.cc, opts...)
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
	// connect to a LOCAL video device
	SendVideoDevice(context.Context, *VideoDeviceDef) (*common.Void, error)
	// connect to any capture source
	SendFromCaptureDevice(context.Context, *CaptureDevice) (*common.Void, error)
	// switch to "idle"
	SwitchToIdle(context.Context, *common.Void) (*common.Void, error)
	// get loopback config info
	GetLoopbackInfo(context.Context, *common.Void) (*LoopbackInfo, error)
	// set idle text
	SetIdleText(context.Context, *IdleTextRequest) (*common.Void, error)
	// switch to images-client
	SwitchToLiveImages(context.Context, *URL) (*common.Void, error)
	// get list of video devices
	GetCaptureDevices(context.Context, *common.Void) (*CaptureDeviceList, error)
	// rpc display a message and countdown
	SetCountdown(context.Context, *CountdownRequest) (*common.Void, error)
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

func _WebCamMixer_SendFromCaptureDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CaptureDevice)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).SendFromCaptureDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/SendFromCaptureDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).SendFromCaptureDevice(ctx, req.(*CaptureDevice))
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

func _WebCamMixer_GetCaptureDevices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).GetCaptureDevices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/GetCaptureDevices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).GetCaptureDevices(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebCamMixer_SetCountdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountdownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebCamMixerServer).SetCountdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webcammixer.WebCamMixer/SetCountdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebCamMixerServer).SetCountdown(ctx, req.(*CountdownRequest))
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
			MethodName: "SendFromCaptureDevice",
			Handler:    _WebCamMixer_SendFromCaptureDevice_Handler,
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
		{
			MethodName: "GetCaptureDevices",
			Handler:    _WebCamMixer_GetCaptureDevices_Handler,
		},
		{
			MethodName: "SetCountdown",
			Handler:    _WebCamMixer_SetCountdown_Handler,
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
	// 576 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x54, 0xdf, 0x6f, 0xd3, 0x30,
	0x10, 0x56, 0xb7, 0xb1, 0xb1, 0x4b, 0x47, 0x3b, 0x8b, 0x1f, 0xa1, 0x8c, 0x51, 0x45, 0x02, 0x55,
	0x3c, 0xa4, 0xd2, 0x98, 0x90, 0x18, 0x48, 0x03, 0x5a, 0x6d, 0x54, 0x2a, 0x03, 0x25, 0xfb, 0x21,
	0xf1, 0xe6, 0x26, 0xb7, 0xd6, 0xa2, 0x8e, 0x4b, 0xe2, 0xae, 0xe5, 0x85, 0x7f, 0x93, 0x7f, 0x07,
	0xd9, 0x49, 0xc1, 0x4e, 0x45, 0x85, 0x78, 0xf2, 0xdd, 0xf9, 0xbe, 0xfb, 0xee, 0xce, 0x9f, 0x0c,
	0xc7, 0x93, 0x54, 0x48, 0x91, 0xb5, 0x87, 0x62, 0x4c, 0x93, 0xa1, 0x1f, 0x89, 0x24, 0xa5, 0xf1,
	0x4c, 0x88, 0xd8, 0x4f, 0x50, 0xb6, 0xe9, 0x84, 0x65, 0xed, 0x19, 0x0e, 0x22, 0xca, 0x39, 0x9b,
	0x63, 0x6a, 0xda, 0xbe, 0x46, 0x12, 0xc7, 0x08, 0x35, 0xfc, 0x15, 0x65, 0x22, 0xc1, 0xb9, 0x48,
	0x8a, 0x23, 0x07, 0x7b, 0xc7, 0xe0, 0xf4, 0x38, 0x1d, 0x62, 0x28, 0x53, 0xa4, 0x9c, 0xec, 0xc1,
	0xf6, 0x19, 0xce, 0xa5, 0x0e, 0xb9, 0x95, 0x66, 0xa5, 0x75, 0x3b, 0xf8, 0x13, 0x20, 0x04, 0x36,
	0xba, 0x54, 0x52, 0x77, 0xad, 0x59, 0x69, 0x55, 0x03, 0x6d, 0x7b, 0x47, 0x70, 0xe7, 0x92, 0xc5,
	0x28, 0xba, 0x78, 0xc3, 0x22, 0xec, 0xe2, 0x35, 0x69, 0x41, 0xcd, 0x88, 0x9c, 0x51, 0x9e, 0x57,
	0xda, 0x0e, 0xca, 0x61, 0x45, 0x7e, 0x92, 0x52, 0xfe, 0xff, 0xe4, 0x3f, 0xa0, 0xda, 0x17, 0x62,
	0x32, 0xa0, 0xd1, 0xd7, 0x5e, 0x72, 0x2d, 0xc8, 0x3e, 0xc0, 0x12, 0xab, 0x11, 0x21, 0x77, 0xe1,
	0xd6, 0x15, 0x8b, 0xe5, 0x48, 0x17, 0xd9, 0x09, 0x72, 0x87, 0xdc, 0x87, 0xcd, 0x0f, 0xc8, 0x86,
	0x23, 0xe9, 0xae, 0xeb, 0x70, 0xe1, 0x91, 0x26, 0x38, 0x9f, 0xd9, 0x1c, 0xc7, 0x27, 0x22, 0xe5,
	0x54, 0xba, 0x1b, 0xfa, 0xd2, 0x0c, 0x79, 0x4f, 0xa1, 0xd6, 0x8b, 0xc7, 0x78, 0x8e, 0x73, 0x19,
	0xe0, 0xb7, 0x29, 0x66, 0x52, 0xb5, 0xa9, 0xdc, 0x82, 0x5c, 0xdb, 0xde, 0x03, 0x58, 0xbf, 0x08,
	0xfa, 0xa4, 0xae, 0x8f, 0xe2, 0x46, 0x99, 0x5e, 0x0f, 0x76, 0x3b, 0x74, 0x22, 0xa7, 0x29, 0xe6,
	0x4d, 0xf6, 0x59, 0x26, 0xc9, 0x21, 0x6c, 0xe5, 0x5e, 0xe6, 0x56, 0x9a, 0xeb, 0x2d, 0xe7, 0xa0,
	0xe1, 0x9b, 0x8f, 0x6e, 0x01, 0x82, 0x45, 0xaa, 0xf7, 0x09, 0x76, 0xac, 0x1b, 0x35, 0x55, 0x6e,
	0x15, 0x84, 0x85, 0xa7, 0x1a, 0xd4, 0xdb, 0x59, 0xcb, 0x1b, 0xd4, 0x7b, 0x51, 0x4d, 0x7f, 0x9f,
	0x60, 0x31, 0xbf, 0xb6, 0xbd, 0xb7, 0x50, 0xef, 0x88, 0x69, 0x22, 0x63, 0x31, 0x4b, 0x56, 0x0c,
	0x47, 0x5c, 0xd8, 0x0a, 0x31, 0x12, 0x49, 0x9c, 0x15, 0x5b, 0x5d, 0xb8, 0x07, 0x3f, 0x37, 0xc0,
	0xb9, 0xc2, 0x41, 0x87, 0xf2, 0x8f, 0xaa, 0x73, 0xf2, 0x12, 0x20, 0xc4, 0x24, 0xd6, 0xcf, 0x99,
	0x11, 0xd7, 0x9a, 0xca, 0x10, 0x61, 0xa3, 0xea, 0x17, 0x12, 0xbd, 0x14, 0x2c, 0x6e, 0x55, 0x16,
	0x38, 0x2d, 0x95, 0x32, 0xce, 0xd0, 0xcf, 0x12, 0xee, 0x0d, 0xd4, 0x14, 0xce, 0x50, 0x1d, 0x79,
	0x64, 0x81, 0x6d, 0xe1, 0xda, 0x78, 0xf2, 0x0e, 0xee, 0xe5, 0xac, 0x82, 0xdb, 0x8b, 0x5d, 0xf1,
	0x1c, 0xa5, 0x12, 0xcf, 0xa1, 0x1a, 0xce, 0x98, 0x8c, 0x46, 0xe7, 0x42, 0xc9, 0x84, 0x58, 0xb7,
	0xa5, 0xdc, 0x23, 0xa8, 0x9d, 0xa2, 0xb4, 0xd4, 0x6c, 0xa7, 0x3f, 0xb4, 0x68, 0xad, 0xc4, 0x57,
	0xe0, 0x84, 0x28, 0x17, 0x4a, 0x24, 0x7b, 0xf6, 0x66, 0x6d, 0x81, 0x96, 0x68, 0x0f, 0x81, 0x2c,
	0x5a, 0xec, 0xb3, 0x1b, 0x2c, 0xde, 0xa6, 0x6e, 0x55, 0xb8, 0x08, 0xfa, 0x4b, 0xbb, 0xd9, 0x3d,
	0x45, 0x69, 0x8d, 0x9e, 0x95, 0xda, 0xdd, 0xff, 0xfb, 0x96, 0xb4, 0xca, 0x5f, 0x43, 0x35, 0x44,
	0xf9, 0x5b, 0x61, 0xe4, 0xb1, 0x9d, 0x5f, 0x52, 0x9e, 0xcd, 0xff, 0xbe, 0x07, 0x4f, 0x12, 0x94,
	0xe6, 0x27, 0xa7, 0x3e, 0x38, 0xb3, 0xc2, 0x97, 0x67, 0xff, 0xf6, 0x9f, 0x0e, 0x36, 0xf5, 0x3f,
	0xf8, 0xe2, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe9, 0x22, 0x74, 0x87, 0x87, 0x05, 0x00, 0x00,
}
