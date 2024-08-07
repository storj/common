// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: debug.proto

package pb

import (
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
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

type CollectRuntimeTracesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CollectRuntimeTracesRequest) Reset()         { *m = CollectRuntimeTracesRequest{} }
func (m *CollectRuntimeTracesRequest) String() string { return proto.CompactTextString(m) }
func (*CollectRuntimeTracesRequest) ProtoMessage()    {}
func (*CollectRuntimeTracesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d9d361be58531fb, []int{0}
}
func (m *CollectRuntimeTracesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollectRuntimeTracesRequest.Unmarshal(m, b)
}
func (m *CollectRuntimeTracesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollectRuntimeTracesRequest.Marshal(b, m, deterministic)
}
func (m *CollectRuntimeTracesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollectRuntimeTracesRequest.Merge(m, src)
}
func (m *CollectRuntimeTracesRequest) XXX_Size() int {
	return xxx_messageInfo_CollectRuntimeTracesRequest.Size(m)
}
func (m *CollectRuntimeTracesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CollectRuntimeTracesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CollectRuntimeTracesRequest proto.InternalMessageInfo

type CollectRuntimeTracesResponse struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CollectRuntimeTracesResponse) Reset()         { *m = CollectRuntimeTracesResponse{} }
func (m *CollectRuntimeTracesResponse) String() string { return proto.CompactTextString(m) }
func (*CollectRuntimeTracesResponse) ProtoMessage()    {}
func (*CollectRuntimeTracesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8d9d361be58531fb, []int{1}
}
func (m *CollectRuntimeTracesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollectRuntimeTracesResponse.Unmarshal(m, b)
}
func (m *CollectRuntimeTracesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollectRuntimeTracesResponse.Marshal(b, m, deterministic)
}
func (m *CollectRuntimeTracesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollectRuntimeTracesResponse.Merge(m, src)
}
func (m *CollectRuntimeTracesResponse) XXX_Size() int {
	return xxx_messageInfo_CollectRuntimeTracesResponse.Size(m)
}
func (m *CollectRuntimeTracesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CollectRuntimeTracesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CollectRuntimeTracesResponse proto.InternalMessageInfo

func (m *CollectRuntimeTracesResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*CollectRuntimeTracesRequest)(nil), "debug.CollectRuntimeTracesRequest")
	proto.RegisterType((*CollectRuntimeTracesResponse)(nil), "debug.CollectRuntimeTracesResponse")
}

func init() { proto.RegisterFile("debug.proto", fileDescriptor_8d9d361be58531fb) }

var fileDescriptor_8d9d361be58531fb = []byte{
	// 157 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0x49, 0x4d, 0x2a,
	0x4d, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x64, 0xb9, 0xa4, 0x9d,
	0xf3, 0x73, 0x72, 0x52, 0x93, 0x4b, 0x82, 0x4a, 0xf3, 0x4a, 0x32, 0x73, 0x53, 0x43, 0x8a, 0x12,
	0x93, 0x53, 0x8b, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x94, 0x8c, 0xb8, 0x64, 0xb0, 0x4b,
	0x17, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x09, 0x71, 0xb1, 0xa4, 0x24, 0x96, 0x24, 0x4a, 0x30,
	0x2a, 0x30, 0x6a, 0xf0, 0x04, 0x81, 0xd9, 0x46, 0x59, 0x5c, 0xac, 0x2e, 0x20, 0xb3, 0x85, 0x12,
	0xb9, 0x44, 0xb0, 0x69, 0x16, 0x52, 0xd2, 0x83, 0x38, 0x04, 0x8f, 0xc5, 0x52, 0xca, 0x78, 0xd5,
	0x40, 0x6c, 0x37, 0x60, 0x74, 0x12, 0x89, 0x12, 0x2a, 0x2e, 0xc9, 0x2f, 0xca, 0xd2, 0xcb, 0xcc,
	0xd7, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xd3, 0x2f, 0x48, 0x4a, 0x62, 0x03, 0x7b, 0xd1, 0x18,
	0x10, 0x00, 0x00, 0xff, 0xff, 0x32, 0x73, 0x6f, 0x88, 0xf1, 0x00, 0x00, 0x00,
}
