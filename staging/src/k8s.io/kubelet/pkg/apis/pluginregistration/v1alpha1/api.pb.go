/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api.proto

package v1alpha1

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// PluginInfo is the message sent from a plugin to the Kubelet pluginwatcher for plugin registration
type PluginInfo struct {
	// Type of the Plugin. CSIPlugin or DevicePlugin
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Plugin name that uniquely identifies the plugin for the given plugin type.
	// For DevicePlugin, this is the resource name that the plugin manages and
	// should follow the extended resource name convention.
	// For CSI, this is the CSI driver registrar name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Optional endpoint location. If found set by Kubelet component,
	// Kubelet component will use this endpoint for specific requests.
	// This allows the plugin to register using one endpoint and possibly use
	// a different socket for control operations. CSI uses this model to delegate
	// its registration external from the plugin.
	Endpoint string `protobuf:"bytes,3,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// Plugin service API versions the plugin supports.
	// For DevicePlugin, this maps to the deviceplugin API versions the
	// plugin supports at the given socket.
	// The Kubelet component communicating with the plugin should be able
	// to choose any preferred version from this list, or returns an error
	// if none of the listed versions is supported.
	SupportedVersions    []string `protobuf:"bytes,4,rep,name=supported_versions,json=supportedVersions,proto3" json:"supported_versions,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PluginInfo) Reset()      { *m = PluginInfo{} }
func (*PluginInfo) ProtoMessage() {}
func (*PluginInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}
func (m *PluginInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PluginInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PluginInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PluginInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PluginInfo.Merge(m, src)
}
func (m *PluginInfo) XXX_Size() int {
	return m.Size()
}
func (m *PluginInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PluginInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PluginInfo proto.InternalMessageInfo

func (m *PluginInfo) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *PluginInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PluginInfo) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *PluginInfo) GetSupportedVersions() []string {
	if m != nil {
		return m.SupportedVersions
	}
	return nil
}

// RegistrationStatus is the message sent from Kubelet pluginwatcher to the plugin for notification on registration status
type RegistrationStatus struct {
	// True if plugin gets registered successfully at Kubelet
	PluginRegistered bool `protobuf:"varint,1,opt,name=plugin_registered,json=pluginRegistered,proto3" json:"plugin_registered,omitempty"`
	// Error message in case plugin fails to register, empty string otherwise
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegistrationStatus) Reset()      { *m = RegistrationStatus{} }
func (*RegistrationStatus) ProtoMessage() {}
func (*RegistrationStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}
func (m *RegistrationStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RegistrationStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RegistrationStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RegistrationStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegistrationStatus.Merge(m, src)
}
func (m *RegistrationStatus) XXX_Size() int {
	return m.Size()
}
func (m *RegistrationStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_RegistrationStatus.DiscardUnknown(m)
}

var xxx_messageInfo_RegistrationStatus proto.InternalMessageInfo

func (m *RegistrationStatus) GetPluginRegistered() bool {
	if m != nil {
		return m.PluginRegistered
	}
	return false
}

func (m *RegistrationStatus) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

// RegistrationStatusResponse is sent by plugin to kubelet in response to RegistrationStatus RPC
type RegistrationStatusResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegistrationStatusResponse) Reset()      { *m = RegistrationStatusResponse{} }
func (*RegistrationStatusResponse) ProtoMessage() {}
func (*RegistrationStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}
func (m *RegistrationStatusResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RegistrationStatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RegistrationStatusResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RegistrationStatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegistrationStatusResponse.Merge(m, src)
}
func (m *RegistrationStatusResponse) XXX_Size() int {
	return m.Size()
}
func (m *RegistrationStatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegistrationStatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegistrationStatusResponse proto.InternalMessageInfo

// InfoRequest is the empty request message from Kubelet
type InfoRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoRequest) Reset()      { *m = InfoRequest{} }
func (*InfoRequest) ProtoMessage() {}
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}
func (m *InfoRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InfoRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoRequest.Merge(m, src)
}
func (m *InfoRequest) XXX_Size() int {
	return m.Size()
}
func (m *InfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InfoRequest proto.InternalMessageInfo

func init() {
	proto.RegisterType((*PluginInfo)(nil), "pluginregistration.PluginInfo")
	proto.RegisterType((*RegistrationStatus)(nil), "pluginregistration.RegistrationStatus")
	proto.RegisterType((*RegistrationStatusResponse)(nil), "pluginregistration.RegistrationStatusResponse")
	proto.RegisterType((*InfoRequest)(nil), "pluginregistration.InfoRequest")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 373 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x41, 0x4b, 0xe3, 0x40,
	0x14, 0xce, 0x6c, 0xbb, 0xbb, 0xed, 0xec, 0x2e, 0x6c, 0x87, 0x3d, 0x84, 0xb0, 0x8c, 0x25, 0x07,
	0x29, 0x48, 0x13, 0x6a, 0x2f, 0x9e, 0xbd, 0x88, 0x20, 0xa2, 0x11, 0x14, 0xbc, 0x94, 0xc4, 0xbe,
	0xa6, 0x43, 0xdb, 0x99, 0x71, 0x66, 0x52, 0xe8, 0x49, 0x7f, 0x82, 0x3f, 0xab, 0x47, 0xf1, 0xe4,
	0xd1, 0xc6, 0x3f, 0x22, 0x9d, 0x94, 0x58, 0x48, 0x0f, 0xde, 0xde, 0xf7, 0xbd, 0xef, 0xbd, 0x79,
	0xdf, 0xc7, 0xe0, 0x66, 0x2c, 0x59, 0x20, 0x95, 0x30, 0x82, 0x10, 0x39, 0xcd, 0x52, 0xc6, 0x15,
	0xa4, 0x4c, 0x1b, 0x15, 0x1b, 0x26, 0xb8, 0xd7, 0x4d, 0x99, 0x19, 0x67, 0x49, 0x70, 0x27, 0x66,
	0x61, 0x2a, 0x52, 0x11, 0x5a, 0x69, 0x92, 0x8d, 0x2c, 0xb2, 0xc0, 0x56, 0xc5, 0x0a, 0xff, 0x01,
	0xe3, 0x0b, 0xbb, 0xe4, 0x94, 0x8f, 0x04, 0x21, 0xb8, 0x6e, 0x16, 0x12, 0x5c, 0xd4, 0x46, 0x9d,
	0x66, 0x64, 0xeb, 0x35, 0xc7, 0xe3, 0x19, 0xb8, 0xdf, 0x0a, 0x6e, 0x5d, 0x13, 0x0f, 0x37, 0x80,
	0x0f, 0xa5, 0x60, 0xdc, 0xb8, 0x35, 0xcb, 0x97, 0x98, 0x74, 0x31, 0xd1, 0x99, 0x94, 0x42, 0x19,
	0x18, 0x0e, 0xe6, 0xa0, 0x34, 0x13, 0x5c, 0xbb, 0xf5, 0x76, 0xad, 0xd3, 0x8c, 0x5a, 0x65, 0xe7,
	0x7a, 0xd3, 0xf0, 0x6f, 0x30, 0x89, 0xb6, 0xee, 0xbf, 0x32, 0xb1, 0xc9, 0x34, 0x39, 0xc0, 0xad,
	0xc2, 0xdb, 0xa0, 0x30, 0x07, 0x0a, 0x86, 0xf6, 0xaa, 0x46, 0xf4, 0xb7, 0x68, 0x44, 0x25, 0x4f,
	0xfe, 0xe1, 0xef, 0xa0, 0x94, 0x50, 0x9b, 0x13, 0x0b, 0xe0, 0xff, 0xc7, 0x5e, 0x75, 0x71, 0x04,
	0x5a, 0x0a, 0xae, 0xc1, 0xff, 0x83, 0x7f, 0xad, 0x1d, 0x47, 0x70, 0x9f, 0x81, 0x36, 0x87, 0x2f,
	0x08, 0xff, 0xde, 0x56, 0x93, 0x33, 0xfc, 0xf3, 0x04, 0x8c, 0x0d, 0x65, 0x2f, 0xa8, 0xc6, 0x1c,
	0x6c, 0x0d, 0x7b, 0x74, 0x97, 0xe0, 0x33, 0x55, 0xdf, 0x21, 0x06, 0xbb, 0xe7, 0xc2, 0xb0, 0xd1,
	0x62, 0x87, 0xd5, 0xfd, 0x5d, 0xd3, 0x55, 0x9d, 0x17, 0x7c, 0x4d, 0x57, 0x3a, 0x74, 0x8e, 0x2f,
	0x97, 0x2b, 0x8a, 0x5e, 0x57, 0xd4, 0x79, 0xcc, 0x29, 0x5a, 0xe6, 0x14, 0x3d, 0xe7, 0x14, 0xbd,
	0xe5, 0x14, 0x3d, 0xbd, 0x53, 0xe7, 0xb6, 0x3f, 0x39, 0xd2, 0x01, 0x13, 0xe1, 0x24, 0x4b, 0x60,
	0x0a, 0x26, 0x94, 0x93, 0x34, 0x8c, 0x25, 0xd3, 0x61, 0xf5, 0x99, 0x70, 0xde, 0x8b, 0xa7, 0x72,
	0x1c, 0xf7, 0x92, 0x1f, 0xf6, 0xd7, 0xf4, 0x3f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x60, 0x59, 0x70,
	0xa3, 0x85, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RegistrationClient is the client API for Registration service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RegistrationClient interface {
	GetInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*PluginInfo, error)
	NotifyRegistrationStatus(ctx context.Context, in *RegistrationStatus, opts ...grpc.CallOption) (*RegistrationStatusResponse, error)
}

type registrationClient struct {
	cc *grpc.ClientConn
}

func NewRegistrationClient(cc *grpc.ClientConn) RegistrationClient {
	return &registrationClient{cc}
}

func (c *registrationClient) GetInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*PluginInfo, error) {
	out := new(PluginInfo)
	err := c.cc.Invoke(ctx, "/pluginregistration.Registration/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrationClient) NotifyRegistrationStatus(ctx context.Context, in *RegistrationStatus, opts ...grpc.CallOption) (*RegistrationStatusResponse, error) {
	out := new(RegistrationStatusResponse)
	err := c.cc.Invoke(ctx, "/pluginregistration.Registration/NotifyRegistrationStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistrationServer is the server API for Registration service.
type RegistrationServer interface {
	GetInfo(context.Context, *InfoRequest) (*PluginInfo, error)
	NotifyRegistrationStatus(context.Context, *RegistrationStatus) (*RegistrationStatusResponse, error)
}

// UnimplementedRegistrationServer can be embedded to have forward compatible implementations.
type UnimplementedRegistrationServer struct {
}

func (*UnimplementedRegistrationServer) GetInfo(ctx context.Context, req *InfoRequest) (*PluginInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (*UnimplementedRegistrationServer) NotifyRegistrationStatus(ctx context.Context, req *RegistrationStatus) (*RegistrationStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyRegistrationStatus not implemented")
}

func RegisterRegistrationServer(s *grpc.Server, srv RegistrationServer) {
	s.RegisterService(&_Registration_serviceDesc, srv)
}

func _Registration_GetInfo_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pluginregistration.Registration/GetInfo",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RegistrationServer).GetInfo(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registration_NotifyRegistrationStatus_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(RegistrationStatus)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrationServer).NotifyRegistrationStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pluginregistration.Registration/NotifyRegistrationStatus",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RegistrationServer).NotifyRegistrationStatus(ctx, req.(*RegistrationStatus))
	}
	return interceptor(ctx, in, info, handler)
}

var _Registration_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pluginregistration.Registration",
	HandlerType: (*RegistrationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _Registration_GetInfo_Handler,
		},
		{
			MethodName: "NotifyRegistrationStatus",
			Handler:    _Registration_NotifyRegistrationStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func (m *PluginInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PluginInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PluginInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SupportedVersions) > 0 {
		for iNdEx := len(m.SupportedVersions) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.SupportedVersions[iNdEx])
			copy(dAtA[i:], m.SupportedVersions[iNdEx])
			i = encodeVarintApi(dAtA, i, uint64(len(m.SupportedVersions[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Endpoint) > 0 {
		i -= len(m.Endpoint)
		copy(dAtA[i:], m.Endpoint)
		i = encodeVarintApi(dAtA, i, uint64(len(m.Endpoint)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintApi(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintApi(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RegistrationStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RegistrationStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RegistrationStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Error) > 0 {
		i -= len(m.Error)
		copy(dAtA[i:], m.Error)
		i = encodeVarintApi(dAtA, i, uint64(len(m.Error)))
		i--
		dAtA[i] = 0x12
	}
	if m.PluginRegistered {
		i--
		if m.PluginRegistered {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *RegistrationStatusResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RegistrationStatusResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RegistrationStatusResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *InfoRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InfoRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InfoRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintApi(dAtA []byte, offset int, v uint64) int {
	offset -= sovApi(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PluginInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	l = len(m.Endpoint)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	if len(m.SupportedVersions) > 0 {
		for _, s := range m.SupportedVersions {
			l = len(s)
			n += 1 + l + sovApi(uint64(l))
		}
	}
	return n
}

func (m *RegistrationStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PluginRegistered {
		n += 2
	}
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovApi(uint64(l))
	}
	return n
}

func (m *RegistrationStatusResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *InfoRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovApi(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozApi(x uint64) (n int) {
	return sovApi(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *PluginInfo) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PluginInfo{`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Endpoint:` + fmt.Sprintf("%v", this.Endpoint) + `,`,
		`SupportedVersions:` + fmt.Sprintf("%v", this.SupportedVersions) + `,`,
		`}`,
	}, "")
	return s
}
func (this *RegistrationStatus) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RegistrationStatus{`,
		`PluginRegistered:` + fmt.Sprintf("%v", this.PluginRegistered) + `,`,
		`Error:` + fmt.Sprintf("%v", this.Error) + `,`,
		`}`,
	}, "")
	return s
}
func (this *RegistrationStatusResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RegistrationStatusResponse{`,
		`}`,
	}, "")
	return s
}
func (this *InfoRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&InfoRequest{`,
		`}`,
	}, "")
	return s
}
func valueToStringApi(v any) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *PluginInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PluginInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PluginInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthApi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthApi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Endpoint", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthApi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Endpoint = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SupportedVersions", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthApi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SupportedVersions = append(m.SupportedVersions, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthApi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RegistrationStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RegistrationStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RegistrationStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PluginRegistered", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.PluginRegistered = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthApi
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthApi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthApi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RegistrationStatusResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RegistrationStatusResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RegistrationStatusResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipApi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthApi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *InfoRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: InfoRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InfoRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipApi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthApi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipApi(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowApi
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowApi
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowApi
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthApi
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupApi
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthApi
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthApi        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowApi          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupApi = fmt.Errorf("proto: unexpected end of group")
)
