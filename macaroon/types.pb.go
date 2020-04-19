// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types.proto

package macaroon

import (
	fmt "fmt"
	math "math"
	time "time"

	proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Caveat struct {
	// if any of these three are set, disallow that type of access
	DisallowReads   bool           `protobuf:"varint,1,opt,name=disallow_reads,json=disallowReads,proto3" json:"disallow_reads,omitempty"`
	DisallowWrites  bool           `protobuf:"varint,2,opt,name=disallow_writes,json=disallowWrites,proto3" json:"disallow_writes,omitempty"`
	DisallowLists   bool           `protobuf:"varint,3,opt,name=disallow_lists,json=disallowLists,proto3" json:"disallow_lists,omitempty"`
	DisallowDeletes bool           `protobuf:"varint,4,opt,name=disallow_deletes,json=disallowDeletes,proto3" json:"disallow_deletes,omitempty"`
	AllowedPaths    []*Caveat_Path `protobuf:"bytes,10,rep,name=allowed_paths,json=allowedPaths,proto3" json:"allowed_paths,omitempty"`
	// if set, the validity time window
	NotAfter  *time.Time `protobuf:"bytes,20,opt,name=not_after,json=notAfter,proto3,stdtime" json:"not_after,omitempty"`
	NotBefore *time.Time `protobuf:"bytes,21,opt,name=not_before,json=notBefore,proto3,stdtime" json:"not_before,omitempty"`
	// nonce is set to some random bytes so that you can make arbitrarily
	// many restricted macaroons with the same (or no) restrictions.
	Nonce                []byte   `protobuf:"bytes,30,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Caveat) Reset()         { *m = Caveat{} }
func (m *Caveat) String() string { return proto.CompactTextString(m) }
func (*Caveat) ProtoMessage()    {}
func (*Caveat) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}
func (m *Caveat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Caveat.Unmarshal(m, b)
}
func (m *Caveat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Caveat.Marshal(b, m, deterministic)
}
func (m *Caveat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Caveat.Merge(m, src)
}
func (m *Caveat) XXX_Size() int {
	return xxx_messageInfo_Caveat.Size(m)
}
func (m *Caveat) XXX_DiscardUnknown() {
	xxx_messageInfo_Caveat.DiscardUnknown(m)
}

var xxx_messageInfo_Caveat proto.InternalMessageInfo

func (m *Caveat) GetDisallowReads() bool {
	if m != nil {
		return m.DisallowReads
	}
	return false
}

func (m *Caveat) GetDisallowWrites() bool {
	if m != nil {
		return m.DisallowWrites
	}
	return false
}

func (m *Caveat) GetDisallowLists() bool {
	if m != nil {
		return m.DisallowLists
	}
	return false
}

func (m *Caveat) GetDisallowDeletes() bool {
	if m != nil {
		return m.DisallowDeletes
	}
	return false
}

func (m *Caveat) GetAllowedPaths() []*Caveat_Path {
	if m != nil {
		return m.AllowedPaths
	}
	return nil
}

func (m *Caveat) GetNotAfter() *time.Time {
	if m != nil {
		return m.NotAfter
	}
	return nil
}

func (m *Caveat) GetNotBefore() *time.Time {
	if m != nil {
		return m.NotBefore
	}
	return nil
}

func (m *Caveat) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

// If any entries exist, require all access to happen in at least
// one of them.
type Caveat_Path struct {
	Bucket               []byte   `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	EncryptedPathPrefix  []byte   `protobuf:"bytes,2,opt,name=encrypted_path_prefix,json=encryptedPathPrefix,proto3" json:"encrypted_path_prefix,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Caveat_Path) Reset()         { *m = Caveat_Path{} }
func (m *Caveat_Path) String() string { return proto.CompactTextString(m) }
func (*Caveat_Path) ProtoMessage()    {}
func (*Caveat_Path) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0, 0}
}
func (m *Caveat_Path) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Caveat_Path.Unmarshal(m, b)
}
func (m *Caveat_Path) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Caveat_Path.Marshal(b, m, deterministic)
}
func (m *Caveat_Path) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Caveat_Path.Merge(m, src)
}
func (m *Caveat_Path) XXX_Size() int {
	return xxx_messageInfo_Caveat_Path.Size(m)
}
func (m *Caveat_Path) XXX_DiscardUnknown() {
	xxx_messageInfo_Caveat_Path.DiscardUnknown(m)
}

var xxx_messageInfo_Caveat_Path proto.InternalMessageInfo

func (m *Caveat_Path) GetBucket() []byte {
	if m != nil {
		return m.Bucket
	}
	return nil
}

func (m *Caveat_Path) GetEncryptedPathPrefix() []byte {
	if m != nil {
		return m.EncryptedPathPrefix
	}
	return nil
}

func init() {
	proto.RegisterType((*Caveat)(nil), "macaroon.Caveat")
	proto.RegisterType((*Caveat_Path)(nil), "macaroon.Caveat.Path")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 343 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x50, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x25, 0xaf, 0x7d, 0xa5, 0xef, 0x36, 0x7d, 0xef, 0x31, 0xb6, 0x12, 0xb2, 0xb0, 0x41, 0x10,
	0xe3, 0x66, 0x0a, 0x75, 0x27, 0x88, 0x58, 0x5d, 0xba, 0x28, 0x83, 0xe0, 0x32, 0x4c, 0x92, 0x9b,
	0x34, 0x98, 0x66, 0xc2, 0xcc, 0xd4, 0xda, 0xbf, 0xf0, 0xd3, 0xfc, 0x03, 0x7f, 0x45, 0x66, 0xd2,
	0x04, 0xba, 0x73, 0x79, 0xce, 0x3d, 0xe7, 0xdc, 0x7b, 0x0f, 0x8c, 0xf4, 0xbe, 0x46, 0x45, 0x6b,
	0x29, 0xb4, 0x20, 0xc3, 0x0d, 0x4f, 0xb8, 0x14, 0xa2, 0xf2, 0x21, 0x17, 0xb9, 0x68, 0x58, 0x7f,
	0x96, 0x0b, 0x91, 0x97, 0x38, 0xb7, 0x28, 0xde, 0x66, 0x73, 0x5d, 0x6c, 0x50, 0x69, 0xbe, 0xa9,
	0x1b, 0xc1, 0xf9, 0x67, 0x0f, 0x06, 0x0f, 0xfc, 0x0d, 0xb9, 0x26, 0x17, 0xf0, 0x37, 0x2d, 0x14,
	0x2f, 0x4b, 0xb1, 0x8b, 0x24, 0xf2, 0x54, 0x79, 0x4e, 0xe0, 0x84, 0x43, 0x36, 0x6e, 0x59, 0x66,
	0x48, 0x72, 0x09, 0xff, 0x3a, 0xd9, 0x4e, 0x16, 0x1a, 0x95, 0xf7, 0xcb, 0xea, 0x3a, 0xf7, 0x8b,
	0x65, 0x8f, 0xf2, 0xca, 0x42, 0x69, 0xe5, 0xf5, 0x8e, 0xf3, 0x9e, 0x0c, 0x49, 0xae, 0xe0, 0x7f,
	0x27, 0x4b, 0xb1, 0x44, 0x13, 0xd8, 0xb7, 0xc2, 0x6e, 0xcf, 0x63, 0x43, 0x93, 0x1b, 0x18, 0x5b,
	0x8c, 0x69, 0x54, 0x73, 0xbd, 0x56, 0x1e, 0x04, 0xbd, 0x70, 0xb4, 0x98, 0xd2, 0xf6, 0x77, 0xda,
	0xbc, 0x42, 0x57, 0x5c, 0xaf, 0x99, 0x7b, 0xd0, 0x1a, 0xa0, 0xc8, 0x2d, 0xfc, 0xa9, 0x84, 0x8e,
	0x78, 0xa6, 0x51, 0x7a, 0x93, 0xc0, 0x09, 0x47, 0x0b, 0x9f, 0x36, 0xed, 0xd0, 0xb6, 0x1d, 0xfa,
	0xdc, 0xb6, 0xb3, 0xec, 0x7f, 0x7c, 0xcd, 0x1c, 0x36, 0xac, 0x84, 0xbe, 0x37, 0x0e, 0x72, 0x07,
	0x60, 0xec, 0x31, 0x66, 0x42, 0xa2, 0x37, 0xfd, 0xa1, 0xdf, 0xac, 0x5c, 0x5a, 0x0b, 0x99, 0xc0,
	0xef, 0x4a, 0x54, 0x09, 0x7a, 0x67, 0x81, 0x13, 0xba, 0xac, 0x01, 0x3e, 0x83, 0xbe, 0x39, 0x8f,
	0x9c, 0xc2, 0x20, 0xde, 0x26, 0xaf, 0xa8, 0x6d, 0xe7, 0x2e, 0x3b, 0x20, 0xb2, 0x80, 0x29, 0x56,
	0x89, 0xdc, 0xd7, 0xfa, 0xf0, 0x73, 0x54, 0x4b, 0xcc, 0x8a, 0x77, 0x5b, 0xb9, 0xcb, 0x4e, 0xba,
	0xa1, 0x49, 0x59, 0xd9, 0x51, 0x3c, 0xb0, 0xe7, 0x5c, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0xca,
	0x7b, 0x7d, 0xfc, 0x1f, 0x02, 0x00, 0x00,
}