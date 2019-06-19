// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gc.proto

package pb

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	math "math"
	time "time"
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

type RetainRequest struct {
	CreationDate         time.Time `protobuf:"bytes,1,opt,name=creation_date,json=creationDate,proto3,stdtime" json:"creation_date"`
	Filter               []byte    `protobuf:"bytes,2,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *RetainRequest) Reset()         { *m = RetainRequest{} }
func (m *RetainRequest) String() string { return proto.CompactTextString(m) }
func (*RetainRequest) ProtoMessage()    {}
func (*RetainRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{0}
}
func (m *RetainRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RetainRequest.Unmarshal(m, b)
}
func (m *RetainRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RetainRequest.Marshal(b, m, deterministic)
}
func (m *RetainRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RetainRequest.Merge(m, src)
}
func (m *RetainRequest) XXX_Size() int {
	return xxx_messageInfo_RetainRequest.Size(m)
}
func (m *RetainRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RetainRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RetainRequest proto.InternalMessageInfo

func (m *RetainRequest) GetCreationDate() time.Time {
	if m != nil {
		return m.CreationDate
	}
	return time.Time{}
}

func (m *RetainRequest) GetFilter() []byte {
	if m != nil {
		return m.Filter
	}
	return nil
}

type RetainResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RetainResponse) Reset()         { *m = RetainResponse{} }
func (m *RetainResponse) String() string { return proto.CompactTextString(m) }
func (*RetainResponse) ProtoMessage()    {}
func (*RetainResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{1}
}
func (m *RetainResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RetainResponse.Unmarshal(m, b)
}
func (m *RetainResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RetainResponse.Marshal(b, m, deterministic)
}
func (m *RetainResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RetainResponse.Merge(m, src)
}
func (m *RetainResponse) XXX_Size() int {
	return xxx_messageInfo_RetainResponse.Size(m)
}
func (m *RetainResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RetainResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RetainResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RetainRequest)(nil), "gc.RetainRequest")
	proto.RegisterType((*RetainResponse)(nil), "gc.RetainResponse")
}

func init() { proto.RegisterFile("gc.proto", fileDescriptor_5502b0b1493f7734) }

var fileDescriptor_5502b0b1493f7734 = []byte{
	// 222 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8e, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x4d, 0x90, 0x50, 0xc6, 0x56, 0xec, 0x1e, 0xa4, 0xe4, 0x92, 0xd2, 0x53, 0x4f, 0x1b,
	0xa8, 0x6f, 0x50, 0x0b, 0xe2, 0x35, 0x78, 0xf2, 0x22, 0x9b, 0x75, 0x3a, 0x04, 0xb6, 0x99, 0x75,
	0x77, 0xfa, 0x1e, 0x3e, 0x96, 0x4f, 0xa1, 0xaf, 0x22, 0x66, 0x5d, 0xd0, 0xe3, 0x37, 0xfc, 0xf3,
	0xfd, 0x3f, 0xcc, 0xc8, 0x6a, 0x1f, 0x58, 0x58, 0x95, 0x64, 0x6b, 0x20, 0x26, 0x4e, 0x5c, 0x37,
	0xc4, 0x4c, 0x0e, 0xdb, 0x89, 0xfa, 0xf3, 0xb1, 0x95, 0xe1, 0x84, 0x51, 0xcc, 0xc9, 0xa7, 0xc0,
	0x26, 0xc0, 0xa2, 0x43, 0x31, 0xc3, 0xd8, 0xe1, 0xdb, 0x19, 0xa3, 0xa8, 0x47, 0x58, 0xd8, 0x80,
	0x46, 0x06, 0x1e, 0x5f, 0x5e, 0x8d, 0xe0, 0xaa, 0x58, 0x17, 0xdb, 0xab, 0x5d, 0xad, 0x93, 0x49,
	0x67, 0x93, 0x7e, 0xca, 0xa6, 0xfd, 0xec, 0xe3, 0xb3, 0xb9, 0x78, 0xff, 0x6a, 0x8a, 0x6e, 0x9e,
	0x5f, 0x0f, 0x46, 0x50, 0xdd, 0x42, 0x75, 0x1c, 0x9c, 0x60, 0x58, 0x95, 0xeb, 0x62, 0x3b, 0xef,
	0x7e, 0x69, 0x73, 0x03, 0xd7, 0xb9, 0x33, 0x7a, 0x1e, 0x23, 0xee, 0x0e, 0xb0, 0x7c, 0x30, 0xa1,
	0x37, 0x84, 0xf7, 0xec, 0x1c, 0xda, 0x1f, 0x85, 0x6a, 0xa1, 0x4a, 0x31, 0xb5, 0xd4, 0x64, 0xf5,
	0xbf, 0x99, 0xb5, 0xfa, 0x7b, 0x4a, 0x96, 0xfd, 0xe5, 0x73, 0xe9, 0xfb, 0xbe, 0x9a, 0x16, 0xde,
	0x7d, 0x07, 0x00, 0x00, 0xff, 0xff, 0xe4, 0x72, 0x5d, 0xf5, 0x15, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GarbageCollectionClient is the client API for GarbageCollection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GarbageCollectionClient interface {
	Retain(ctx context.Context, in *RetainRequest, opts ...grpc.CallOption) (*RetainResponse, error)
}

type garbageCollectionClient struct {
	cc *grpc.ClientConn
}

func NewGarbageCollectionClient(cc *grpc.ClientConn) GarbageCollectionClient {
	return &garbageCollectionClient{cc}
}

func (c *garbageCollectionClient) Retain(ctx context.Context, in *RetainRequest, opts ...grpc.CallOption) (*RetainResponse, error) {
	out := new(RetainResponse)
	err := c.cc.Invoke(ctx, "/gc.GarbageCollection/Retain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GarbageCollectionServer is the server API for GarbageCollection service.
type GarbageCollectionServer interface {
	Retain(context.Context, *RetainRequest) (*RetainResponse, error)
}

func RegisterGarbageCollectionServer(s *grpc.Server, srv GarbageCollectionServer) {
	s.RegisterService(&_GarbageCollection_serviceDesc, srv)
}

func _GarbageCollection_Retain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GarbageCollectionServer).Retain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gc.GarbageCollection/Retain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GarbageCollectionServer).Retain(ctx, req.(*RetainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GarbageCollection_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gc.GarbageCollection",
	HandlerType: (*GarbageCollectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Retain",
			Handler:    _GarbageCollection_Retain_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gc.proto",
}