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

package v2

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

type StatusRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRequest) Reset()         { *m = StatusRequest{} }
func (m *StatusRequest) String() string { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()    {}
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}
func (m *StatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRequest.Unmarshal(m, b)
}
func (m *StatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRequest.Marshal(b, m, deterministic)
}
func (m *StatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRequest.Merge(m, src)
}
func (m *StatusRequest) XXX_Size() int {
	return xxx_messageInfo_StatusRequest.Size(m)
}
func (m *StatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRequest proto.InternalMessageInfo

type StatusResponse struct {
	// Version of the KMS gRPC plugin API. Must equal v2 to v2beta1 (v2 is recommended, but both are equivalent).
	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	// Any value other than "ok" is failing healthz.  On failure, the associated API server healthz endpoint will contain this value as part of the error message.
	Healthz string `protobuf:"bytes,2,opt,name=healthz,proto3" json:"healthz,omitempty"`
	// the current write key, used to determine staleness of data updated via value.Transformer.TransformFromStorage.
	// keyID must satisfy the following constraints:
	// 1. The keyID is not empty.
	// 2. The size of keyID is less than 1 kB.
	KeyId                string   `protobuf:"bytes,3,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}
func (m *StatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResponse.Unmarshal(m, b)
}
func (m *StatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResponse.Marshal(b, m, deterministic)
}
func (m *StatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResponse.Merge(m, src)
}
func (m *StatusResponse) XXX_Size() int {
	return xxx_messageInfo_StatusResponse.Size(m)
}
func (m *StatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *StatusResponse) GetHealthz() string {
	if m != nil {
		return m.Healthz
	}
	return ""
}

func (m *StatusResponse) GetKeyId() string {
	if m != nil {
		return m.KeyId
	}
	return ""
}

type DecryptRequest struct {
	// The data to be decrypted.
	Ciphertext []byte `protobuf:"bytes,1,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	// UID is a unique identifier for the request.
	Uid string `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	// The keyID that was provided to the apiserver during encryption.
	// This represents the KMS KEK that was used to encrypt the data.
	KeyId string `protobuf:"bytes,3,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	// Additional metadata that was sent by the KMS plugin during encryption.
	Annotations          map[string][]byte `protobuf:"bytes,4,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DecryptRequest) Reset()         { *m = DecryptRequest{} }
func (m *DecryptRequest) String() string { return proto.CompactTextString(m) }
func (*DecryptRequest) ProtoMessage()    {}
func (*DecryptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}
func (m *DecryptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DecryptRequest.Unmarshal(m, b)
}
func (m *DecryptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DecryptRequest.Marshal(b, m, deterministic)
}
func (m *DecryptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DecryptRequest.Merge(m, src)
}
func (m *DecryptRequest) XXX_Size() int {
	return xxx_messageInfo_DecryptRequest.Size(m)
}
func (m *DecryptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DecryptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DecryptRequest proto.InternalMessageInfo

func (m *DecryptRequest) GetCiphertext() []byte {
	if m != nil {
		return m.Ciphertext
	}
	return nil
}

func (m *DecryptRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *DecryptRequest) GetKeyId() string {
	if m != nil {
		return m.KeyId
	}
	return ""
}

func (m *DecryptRequest) GetAnnotations() map[string][]byte {
	if m != nil {
		return m.Annotations
	}
	return nil
}

type DecryptResponse struct {
	// The decrypted data.
	Plaintext            []byte   `protobuf:"bytes,1,opt,name=plaintext,proto3" json:"plaintext,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DecryptResponse) Reset()         { *m = DecryptResponse{} }
func (m *DecryptResponse) String() string { return proto.CompactTextString(m) }
func (*DecryptResponse) ProtoMessage()    {}
func (*DecryptResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}
func (m *DecryptResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DecryptResponse.Unmarshal(m, b)
}
func (m *DecryptResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DecryptResponse.Marshal(b, m, deterministic)
}
func (m *DecryptResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DecryptResponse.Merge(m, src)
}
func (m *DecryptResponse) XXX_Size() int {
	return xxx_messageInfo_DecryptResponse.Size(m)
}
func (m *DecryptResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DecryptResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DecryptResponse proto.InternalMessageInfo

func (m *DecryptResponse) GetPlaintext() []byte {
	if m != nil {
		return m.Plaintext
	}
	return nil
}

type EncryptRequest struct {
	// The data to be encrypted.
	Plaintext []byte `protobuf:"bytes,1,opt,name=plaintext,proto3" json:"plaintext,omitempty"`
	// UID is a unique identifier for the request.
	Uid                  string   `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EncryptRequest) Reset()         { *m = EncryptRequest{} }
func (m *EncryptRequest) String() string { return proto.CompactTextString(m) }
func (*EncryptRequest) ProtoMessage()    {}
func (*EncryptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}
func (m *EncryptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptRequest.Unmarshal(m, b)
}
func (m *EncryptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptRequest.Marshal(b, m, deterministic)
}
func (m *EncryptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptRequest.Merge(m, src)
}
func (m *EncryptRequest) XXX_Size() int {
	return xxx_messageInfo_EncryptRequest.Size(m)
}
func (m *EncryptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptRequest proto.InternalMessageInfo

func (m *EncryptRequest) GetPlaintext() []byte {
	if m != nil {
		return m.Plaintext
	}
	return nil
}

func (m *EncryptRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

type EncryptResponse struct {
	// The encrypted data.
	// ciphertext must satisfy the following constraints:
	// 1. The ciphertext is not empty.
	// 2. The ciphertext is less than 1 kB.
	Ciphertext []byte `protobuf:"bytes,1,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	// The KMS key ID used to encrypt the data. This must always refer to the KMS KEK and not any local KEKs that may be in use.
	// This can be used to inform staleness of data updated via value.Transformer.TransformFromStorage.
	// keyID must satisfy the following constraints:
	// 1. The keyID is not empty.
	// 2. The size of keyID is less than 1 kB.
	KeyId string `protobuf:"bytes,2,opt,name=key_id,json=keyId,proto3" json:"key_id,omitempty"`
	// Additional metadata to be stored with the encrypted data.
	// This data is stored in plaintext in etcd. KMS plugin implementations are responsible for pre-encrypting any sensitive data.
	// Annotations must satisfy the following constraints:
	//  1. Annotation key must be a fully qualified domain name that conforms to the definition in DNS (RFC 1123).
	//  2. The size of annotations keys + values is less than 32 kB.
	Annotations          map[string][]byte `protobuf:"bytes,3,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *EncryptResponse) Reset()         { *m = EncryptResponse{} }
func (m *EncryptResponse) String() string { return proto.CompactTextString(m) }
func (*EncryptResponse) ProtoMessage()    {}
func (*EncryptResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}
func (m *EncryptResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptResponse.Unmarshal(m, b)
}
func (m *EncryptResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptResponse.Marshal(b, m, deterministic)
}
func (m *EncryptResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptResponse.Merge(m, src)
}
func (m *EncryptResponse) XXX_Size() int {
	return xxx_messageInfo_EncryptResponse.Size(m)
}
func (m *EncryptResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptResponse proto.InternalMessageInfo

func (m *EncryptResponse) GetCiphertext() []byte {
	if m != nil {
		return m.Ciphertext
	}
	return nil
}

func (m *EncryptResponse) GetKeyId() string {
	if m != nil {
		return m.KeyId
	}
	return ""
}

func (m *EncryptResponse) GetAnnotations() map[string][]byte {
	if m != nil {
		return m.Annotations
	}
	return nil
}

func init() {
	proto.RegisterType((*StatusRequest)(nil), "v2.StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "v2.StatusResponse")
	proto.RegisterType((*DecryptRequest)(nil), "v2.DecryptRequest")
	proto.RegisterMapType((map[string][]byte)(nil), "v2.DecryptRequest.AnnotationsEntry")
	proto.RegisterType((*DecryptResponse)(nil), "v2.DecryptResponse")
	proto.RegisterType((*EncryptRequest)(nil), "v2.EncryptRequest")
	proto.RegisterType((*EncryptResponse)(nil), "v2.EncryptResponse")
	proto.RegisterMapType((map[string][]byte)(nil), "v2.EncryptResponse.AnnotationsEntry")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0xcd, 0x6e, 0xda, 0x40,
	0x10, 0xc7, 0xb1, 0x5d, 0x40, 0x0c, 0x14, 0xe8, 0x96, 0x4a, 0x16, 0xaa, 0x2a, 0xb4, 0xed, 0x81,
	0x93, 0xad, 0xba, 0x3d, 0xa0, 0x1e, 0xaa, 0xb6, 0x2a, 0x95, 0xaa, 0xaa, 0x17, 0x73, 0x6b, 0x0f,
	0xd1, 0x06, 0x46, 0x61, 0x65, 0x58, 0x3b, 0xde, 0xb5, 0x15, 0xe7, 0xbd, 0xf2, 0x1e, 0x79, 0x84,
	0x3c, 0x4a, 0x64, 0x7b, 0x01, 0x1b, 0x94, 0xe4, 0x94, 0x9b, 0xe7, 0xf3, 0x3f, 0xf3, 0xdb, 0x31,
	0x74, 0x58, 0xc4, 0x9d, 0x28, 0x0e, 0x55, 0x48, 0xcc, 0xd4, 0xa3, 0x03, 0x78, 0xb9, 0x50, 0x4c,
	0x25, 0xd2, 0xc7, 0xcb, 0x04, 0xa5, 0xa2, 0xff, 0xa1, 0xbf, 0x73, 0xc8, 0x28, 0x14, 0x12, 0x89,
	0x0d, 0xed, 0x14, 0x63, 0xc9, 0x43, 0x61, 0x1b, 0x13, 0x63, 0xda, 0xf1, 0x77, 0x66, 0x1e, 0x59,
	0x23, 0xdb, 0xa8, 0xf5, 0xb5, 0x6d, 0x96, 0x11, 0x6d, 0x92, 0x37, 0xd0, 0x0a, 0x30, 0x3b, 0xe3,
	0x2b, 0xdb, 0x2a, 0x02, 0xcd, 0x00, 0xb3, 0xdf, 0x2b, 0x7a, 0x67, 0x40, 0xff, 0x27, 0x2e, 0xe3,
	0x2c, 0x52, 0x5a, 0x8f, 0xbc, 0x03, 0x58, 0xf2, 0x68, 0x8d, 0xb1, 0xc2, 0x2b, 0x55, 0x08, 0xf4,
	0xfc, 0x8a, 0x87, 0x0c, 0xc1, 0x4a, 0xf8, 0x4a, 0xf7, 0xcf, 0x3f, 0x1f, 0xe8, 0x4d, 0xe6, 0xd0,
	0x65, 0x42, 0x84, 0x8a, 0x29, 0x1e, 0x0a, 0x69, 0xbf, 0x98, 0x58, 0xd3, 0xae, 0xf7, 0xde, 0x49,
	0x3d, 0xa7, 0xae, 0xe8, 0x7c, 0x3f, 0x64, 0xcd, 0x85, 0x8a, 0x33, 0xbf, 0x5a, 0x37, 0xfe, 0x0a,
	0xc3, 0xe3, 0x84, 0x7c, 0x86, 0x00, 0x33, 0xbd, 0x7d, 0xfe, 0x49, 0x46, 0xd0, 0x4c, 0xd9, 0x26,
	0xc1, 0x62, 0xae, 0x9e, 0x5f, 0x1a, 0x5f, 0xcc, 0x99, 0x41, 0x5d, 0x18, 0xec, 0xf5, 0x34, 0xc0,
	0xb7, 0xd0, 0x89, 0x36, 0x8c, 0x8b, 0xca, 0x86, 0x07, 0x07, 0xfd, 0x06, 0xfd, 0xb9, 0xa8, 0x21,
	0x79, 0x34, 0xff, 0x14, 0x08, 0xbd, 0x35, 0x60, 0xb0, 0x6f, 0xa1, 0x35, 0x9f, 0xc2, 0x7a, 0x80,
	0x68, 0x56, 0x21, 0xfe, 0xaa, 0x43, 0xb4, 0x0a, 0x88, 0x1f, 0x72, 0x88, 0x47, 0x02, 0xcf, 0x4b,
	0xd1, 0xbb, 0x31, 0x60, 0xf4, 0x07, 0xb3, 0xbf, 0x4c, 0xb0, 0x0b, 0xdc, 0xa2, 0x50, 0x0b, 0x8c,
	0x53, 0xbe, 0x44, 0xf2, 0x11, 0x5a, 0xe5, 0x79, 0x92, 0x57, 0xf9, 0x54, 0xb5, 0xdb, 0x1d, 0x93,
	0xaa, 0xab, 0x9c, 0x93, 0x36, 0xc8, 0x67, 0x68, 0xeb, 0x17, 0x21, 0xe4, 0xf4, 0x1c, 0xc6, 0xaf,
	0x6b, 0xbe, 0x6a, 0x95, 0x5e, 0xb9, 0xac, 0xaa, 0xbf, 0x51, 0x59, 0x75, 0xc4, 0x84, 0x36, 0x7e,
	0x8c, 0xfe, 0x91, 0x60, 0x26, 0x1d, 0x1e, 0xba, 0xc1, 0x56, 0xba, 0x2c, 0xe2, 0xd2, 0x4d, 0xbd,
	0xf3, 0x56, 0xf1, 0xbf, 0x7d, 0xba, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x5f, 0xf8, 0x49, 0x17, 0x7c,
	0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// KeyManagementServiceClient is the client API for KeyManagementService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type KeyManagementServiceClient interface {
	// this API is meant to be polled
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	// Execute decryption operation in KMS provider.
	Decrypt(ctx context.Context, in *DecryptRequest, opts ...grpc.CallOption) (*DecryptResponse, error)
	// Execute encryption operation in KMS provider.
	Encrypt(ctx context.Context, in *EncryptRequest, opts ...grpc.CallOption) (*EncryptResponse, error)
}

type keyManagementServiceClient struct {
	cc *grpc.ClientConn
}

func NewKeyManagementServiceClient(cc *grpc.ClientConn) KeyManagementServiceClient {
	return &keyManagementServiceClient{cc}
}

func (c *keyManagementServiceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/v2.KeyManagementService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyManagementServiceClient) Decrypt(ctx context.Context, in *DecryptRequest, opts ...grpc.CallOption) (*DecryptResponse, error) {
	out := new(DecryptResponse)
	err := c.cc.Invoke(ctx, "/v2.KeyManagementService/Decrypt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyManagementServiceClient) Encrypt(ctx context.Context, in *EncryptRequest, opts ...grpc.CallOption) (*EncryptResponse, error) {
	out := new(EncryptResponse)
	err := c.cc.Invoke(ctx, "/v2.KeyManagementService/Encrypt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyManagementServiceServer is the server API for KeyManagementService service.
type KeyManagementServiceServer interface {
	// this API is meant to be polled
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	// Execute decryption operation in KMS provider.
	Decrypt(context.Context, *DecryptRequest) (*DecryptResponse, error)
	// Execute encryption operation in KMS provider.
	Encrypt(context.Context, *EncryptRequest) (*EncryptResponse, error)
}

// UnimplementedKeyManagementServiceServer can be embedded to have forward compatible implementations.
type UnimplementedKeyManagementServiceServer struct {
}

func (*UnimplementedKeyManagementServiceServer) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (*UnimplementedKeyManagementServiceServer) Decrypt(ctx context.Context, req *DecryptRequest) (*DecryptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Decrypt not implemented")
}
func (*UnimplementedKeyManagementServiceServer) Encrypt(ctx context.Context, req *EncryptRequest) (*EncryptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Encrypt not implemented")
}

func RegisterKeyManagementServiceServer(s *grpc.Server, srv KeyManagementServiceServer) {
	s.RegisterService(&_KeyManagementService_serviceDesc, srv)
}

func _KeyManagementService_Status_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyManagementServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2.KeyManagementService/Status",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KeyManagementServiceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyManagementService_Decrypt_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(DecryptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyManagementServiceServer).Decrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2.KeyManagementService/Decrypt",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KeyManagementServiceServer).Decrypt(ctx, req.(*DecryptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyManagementService_Encrypt_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(EncryptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyManagementServiceServer).Encrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2.KeyManagementService/Encrypt",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(KeyManagementServiceServer).Encrypt(ctx, req.(*EncryptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _KeyManagementService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v2.KeyManagementService",
	HandlerType: (*KeyManagementServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _KeyManagementService_Status_Handler,
		},
		{
			MethodName: "Decrypt",
			Handler:    _KeyManagementService_Decrypt_Handler,
		},
		{
			MethodName: "Encrypt",
			Handler:    _KeyManagementService_Encrypt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
