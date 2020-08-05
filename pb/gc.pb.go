// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: gc.proto

package pb

import (
	context "context"
	fmt "fmt"
	math "math"
	time "time"

	proto "github.com/gogo/protobuf/proto"

	drpc "storj.io/drpc"
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

type PingRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{0}
}
func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (m *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(m, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

type PingResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResponse) Reset()         { *m = PingResponse{} }
func (m *PingResponse) String() string { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()    {}
func (*PingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{1}
}
func (m *PingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResponse.Unmarshal(m, b)
}
func (m *PingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResponse.Marshal(b, m, deterministic)
}
func (m *PingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResponse.Merge(m, src)
}
func (m *PingResponse) XXX_Size() int {
	return xxx_messageInfo_PingResponse.Size(m)
}
func (m *PingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PingResponse proto.InternalMessageInfo

type StartSessionRequest struct {
	SessionId            time.Time `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3,stdtime" json:"session_id"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *StartSessionRequest) Reset()         { *m = StartSessionRequest{} }
func (m *StartSessionRequest) String() string { return proto.CompactTextString(m) }
func (*StartSessionRequest) ProtoMessage()    {}
func (*StartSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{2}
}
func (m *StartSessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartSessionRequest.Unmarshal(m, b)
}
func (m *StartSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartSessionRequest.Marshal(b, m, deterministic)
}
func (m *StartSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartSessionRequest.Merge(m, src)
}
func (m *StartSessionRequest) XXX_Size() int {
	return xxx_messageInfo_StartSessionRequest.Size(m)
}
func (m *StartSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartSessionRequest proto.InternalMessageInfo

func (m *StartSessionRequest) GetSessionId() time.Time {
	if m != nil {
		return m.SessionId
	}
	return time.Time{}
}

type StartSessionResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartSessionResponse) Reset()         { *m = StartSessionResponse{} }
func (m *StartSessionResponse) String() string { return proto.CompactTextString(m) }
func (*StartSessionResponse) ProtoMessage()    {}
func (*StartSessionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{3}
}
func (m *StartSessionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartSessionResponse.Unmarshal(m, b)
}
func (m *StartSessionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartSessionResponse.Marshal(b, m, deterministic)
}
func (m *StartSessionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartSessionResponse.Merge(m, src)
}
func (m *StartSessionResponse) XXX_Size() int {
	return xxx_messageInfo_StartSessionResponse.Size(m)
}
func (m *StartSessionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StartSessionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StartSessionResponse proto.InternalMessageInfo

type AddPieceRequest struct {
	// session_id indicates which GC session this piece ID belongs to
	SessionId time.Time `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3,stdtime" json:"session_id"`
	// piece is the piece_id that should be added to the bloom filter
	PieceId PieceID `protobuf:"bytes,2,opt,name=piece_id,json=pieceId,proto3,customtype=PieceID" json:"piece_id"`
	// sequence_number is the ordered number assigned to this request so the worker can confirm they received all the correct pieces
	SequenceNumber int64 `protobuf:"varint,3,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"`
	// storage_node_id is the id of the storage node this piece is stored on
	StorageNodeId        NodeID   `protobuf:"bytes,4,opt,name=storage_node_id,json=storageNodeId,proto3,customtype=NodeID" json:"storage_node_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddPieceRequest) Reset()         { *m = AddPieceRequest{} }
func (m *AddPieceRequest) String() string { return proto.CompactTextString(m) }
func (*AddPieceRequest) ProtoMessage()    {}
func (*AddPieceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{4}
}
func (m *AddPieceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddPieceRequest.Unmarshal(m, b)
}
func (m *AddPieceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddPieceRequest.Marshal(b, m, deterministic)
}
func (m *AddPieceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddPieceRequest.Merge(m, src)
}
func (m *AddPieceRequest) XXX_Size() int {
	return xxx_messageInfo_AddPieceRequest.Size(m)
}
func (m *AddPieceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddPieceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddPieceRequest proto.InternalMessageInfo

func (m *AddPieceRequest) GetSessionId() time.Time {
	if m != nil {
		return m.SessionId
	}
	return time.Time{}
}

func (m *AddPieceRequest) GetSequenceNumber() int64 {
	if m != nil {
		return m.SequenceNumber
	}
	return 0
}

type AddPieceResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddPieceResponse) Reset()         { *m = AddPieceResponse{} }
func (m *AddPieceResponse) String() string { return proto.CompactTextString(m) }
func (*AddPieceResponse) ProtoMessage()    {}
func (*AddPieceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{5}
}
func (m *AddPieceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddPieceResponse.Unmarshal(m, b)
}
func (m *AddPieceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddPieceResponse.Marshal(b, m, deterministic)
}
func (m *AddPieceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddPieceResponse.Merge(m, src)
}
func (m *AddPieceResponse) XXX_Size() int {
	return xxx_messageInfo_AddPieceResponse.Size(m)
}
func (m *AddPieceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddPieceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddPieceResponse proto.InternalMessageInfo

type EndSessionRequest struct {
	// session_id indicates which session to end
	SessionId time.Time `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3,stdtime" json:"session_id"`
	// node_ending_sequence is a map of storage node ID to its corresponding ending sequence number for how many pieces it should have processed
	NodeEndingSequence   map[string]int64 `protobuf:"bytes,2,rep,name=node_ending_sequence,json=nodeEndingSequence,proto3" json:"node_ending_sequence,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *EndSessionRequest) Reset()         { *m = EndSessionRequest{} }
func (m *EndSessionRequest) String() string { return proto.CompactTextString(m) }
func (*EndSessionRequest) ProtoMessage()    {}
func (*EndSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{6}
}
func (m *EndSessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndSessionRequest.Unmarshal(m, b)
}
func (m *EndSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndSessionRequest.Marshal(b, m, deterministic)
}
func (m *EndSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndSessionRequest.Merge(m, src)
}
func (m *EndSessionRequest) XXX_Size() int {
	return xxx_messageInfo_EndSessionRequest.Size(m)
}
func (m *EndSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EndSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EndSessionRequest proto.InternalMessageInfo

func (m *EndSessionRequest) GetSessionId() time.Time {
	if m != nil {
		return m.SessionId
	}
	return time.Time{}
}

func (m *EndSessionRequest) GetNodeEndingSequence() map[string]int64 {
	if m != nil {
		return m.NodeEndingSequence
	}
	return nil
}

type EndSessionResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EndSessionResponse) Reset()         { *m = EndSessionResponse{} }
func (m *EndSessionResponse) String() string { return proto.CompactTextString(m) }
func (*EndSessionResponse) ProtoMessage()    {}
func (*EndSessionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5502b0b1493f7734, []int{7}
}
func (m *EndSessionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EndSessionResponse.Unmarshal(m, b)
}
func (m *EndSessionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EndSessionResponse.Marshal(b, m, deterministic)
}
func (m *EndSessionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EndSessionResponse.Merge(m, src)
}
func (m *EndSessionResponse) XXX_Size() int {
	return xxx_messageInfo_EndSessionResponse.Size(m)
}
func (m *EndSessionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EndSessionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EndSessionResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*PingRequest)(nil), "gc.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "gc.PingResponse")
	proto.RegisterType((*StartSessionRequest)(nil), "gc.StartSessionRequest")
	proto.RegisterType((*StartSessionResponse)(nil), "gc.StartSessionResponse")
	proto.RegisterType((*AddPieceRequest)(nil), "gc.AddPieceRequest")
	proto.RegisterType((*AddPieceResponse)(nil), "gc.AddPieceResponse")
	proto.RegisterType((*EndSessionRequest)(nil), "gc.EndSessionRequest")
	proto.RegisterMapType((map[string]int64)(nil), "gc.EndSessionRequest.NodeEndingSequenceEntry")
	proto.RegisterType((*EndSessionResponse)(nil), "gc.EndSessionResponse")
}

func init() { proto.RegisterFile("gc.proto", fileDescriptor_5502b0b1493f7734) }

var fileDescriptor_5502b0b1493f7734 = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xad, 0xed, 0xd2, 0xa6, 0xd3, 0x34, 0x4e, 0xa7, 0xa6, 0x8d, 0x7c, 0x49, 0xe4, 0x0b, 0x11,
	0x08, 0x47, 0x0a, 0x12, 0x20, 0x24, 0x0e, 0x24, 0x44, 0x28, 0x97, 0xaa, 0x72, 0x38, 0xf5, 0x62,
	0xd9, 0xde, 0x65, 0x65, 0x88, 0x77, 0x8d, 0xbd, 0x41, 0xea, 0x8d, 0x4f, 0xe0, 0xb3, 0xf8, 0x06,
	0x0e, 0xe5, 0xc0, 0x77, 0x20, 0xa1, 0x5d, 0xdb, 0x4d, 0xda, 0x94, 0x1b, 0xdc, 0x3c, 0xcf, 0x6f,
	0xde, 0x9b, 0x99, 0x7d, 0xd0, 0x62, 0x89, 0x9f, 0x17, 0x42, 0x0a, 0x34, 0x59, 0xe2, 0x02, 0x13,
	0x4c, 0x54, 0xb5, 0xdb, 0x67, 0x42, 0xb0, 0x25, 0x1d, 0xe9, 0x2a, 0x5e, 0x7d, 0x18, 0xc9, 0x34,
	0xa3, 0xa5, 0x8c, 0xb2, 0xbc, 0x22, 0x78, 0x47, 0x70, 0x78, 0x91, 0x72, 0x16, 0xd0, 0xcf, 0x2b,
	0x5a, 0x4a, 0xaf, 0x03, 0xed, 0xaa, 0x2c, 0x73, 0xc1, 0x4b, 0xea, 0x5d, 0xc2, 0xc9, 0x42, 0x46,
	0x85, 0x5c, 0xd0, 0xb2, 0x4c, 0x05, 0xaf, 0x69, 0x38, 0x05, 0x28, 0x2b, 0x24, 0x4c, 0x49, 0xcf,
	0x18, 0x18, 0xc3, 0xc3, 0xb1, 0xeb, 0x57, 0x5e, 0x7e, 0xe3, 0xe5, 0xbf, 0x6f, 0xbc, 0x26, 0xad,
	0xef, 0xd7, 0xfd, 0x9d, 0x6f, 0x3f, 0xfb, 0x46, 0x70, 0x50, 0xf7, 0xcd, 0x89, 0x77, 0x0a, 0xce,
	0x6d, 0xed, 0xda, 0xf3, 0x97, 0x01, 0xf6, 0x1b, 0x42, 0x2e, 0x52, 0x9a, 0xd0, 0x7f, 0x69, 0x88,
	0x8f, 0xa1, 0x95, 0x2b, 0x51, 0x25, 0x61, 0x0e, 0x8c, 0x61, 0x7b, 0x62, 0x2b, 0xda, 0x8f, 0xeb,
	0xfe, 0xbe, 0x36, 0x9b, 0xbf, 0x0d, 0xf6, 0x35, 0x61, 0x4e, 0xf0, 0x11, 0xd8, 0xa5, 0xf2, 0xe6,
	0x09, 0x0d, 0xf9, 0x2a, 0x8b, 0x69, 0xd1, 0xb3, 0x06, 0xc6, 0xd0, 0x0a, 0x3a, 0x0d, 0x7c, 0xae,
	0x51, 0x7c, 0x0e, 0x76, 0x29, 0x45, 0x11, 0x31, 0x1a, 0x72, 0x41, 0xb4, 0xf6, 0xae, 0xd6, 0xee,
	0xd4, 0xda, 0x7b, 0xe7, 0x82, 0x28, 0xe9, 0xa3, 0x9a, 0xa6, 0x4b, 0xe2, 0x21, 0x74, 0xd7, 0x4b,
	0xd6, 0x9b, 0x7f, 0x35, 0xe1, 0x78, 0xc6, 0xc9, 0x7f, 0x38, 0x36, 0x86, 0xe0, 0xe8, 0xf1, 0x28,
	0x27, 0x29, 0x67, 0x61, 0xb3, 0x44, 0xcf, 0x1c, 0x58, 0xc3, 0xc3, 0xf1, 0x53, 0x9f, 0x25, 0xfe,
	0x96, 0xb3, 0xaf, 0x46, 0x9d, 0xe9, 0x86, 0x45, 0xcd, 0x9f, 0x71, 0x59, 0x5c, 0x05, 0xc8, 0xb7,
	0x7e, 0xb8, 0x33, 0x38, 0xfb, 0x0b, 0x1d, 0xbb, 0x60, 0x7d, 0xa2, 0x57, 0x7a, 0xf2, 0x83, 0x40,
	0x7d, 0xa2, 0x03, 0x0f, 0xbe, 0x44, 0xcb, 0x15, 0xd5, 0xcf, 0x60, 0x05, 0x55, 0xf1, 0xca, 0x7c,
	0x69, 0x78, 0x0e, 0xe0, 0xe6, 0x1c, 0xd5, 0x61, 0xc6, 0xbf, 0x0d, 0x38, 0x7e, 0x17, 0x15, 0x71,
	0xc4, 0xe8, 0x54, 0x2c, 0x97, 0x34, 0x91, 0xa9, 0xe0, 0xf8, 0x04, 0x76, 0x55, 0x58, 0xd1, 0x56,
	0xd3, 0x6f, 0xa4, 0xd8, 0xed, 0xae, 0x81, 0xfa, 0xb2, 0x3b, 0x38, 0x85, 0xf6, 0x66, 0xda, 0xf0,
	0x4c, 0x71, 0xee, 0xc9, 0xb6, 0xdb, 0xdb, 0xfe, 0x71, 0x23, 0xf2, 0x02, 0x5a, 0xcd, 0xa3, 0xe1,
	0x89, 0xe2, 0xdd, 0xc9, 0xa9, 0xeb, 0xdc, 0x06, 0x6f, 0x1a, 0x5f, 0x03, 0xac, 0xd7, 0xc2, 0x87,
	0xf7, 0x9e, 0xdb, 0x3d, 0xbd, 0x0b, 0x37, 0xed, 0x13, 0xe7, 0x12, 0x55, 0x7a, 0x3e, 0xfa, 0xa9,
	0x18, 0x25, 0x22, 0xcb, 0x04, 0x1f, 0xe5, 0x71, 0xbc, 0xa7, 0x1f, 0xff, 0xd9, 0x9f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x8a, 0x2d, 0x49, 0xec, 0xff, 0x03, 0x00, 0x00,
}

// --- DRPC BEGIN ---

type DRPCGarbageCollectionClient interface {
	DRPCConn() drpc.Conn

	Ping(ctx context.Context, in *PingRequest) (*PingResponse, error)
	// StartSession begins a new garbage collection session
	StartSession(ctx context.Context, in *StartSessionRequest) (*StartSessionResponse, error)
	// AddPiece adds a piece ID for a specifc storage node along with a
	// sequence number that helps ensuring no pieces were missed
	AddPiece(ctx context.Context, in *AddPieceRequest) (*AddPieceResponse, error)
	// EndSession ends the garbage collection session indicating how many pieces
	// should have been processed for each storage node
	EndSession(ctx context.Context, in *EndSessionRequest) (*EndSessionResponse, error)
}

type drpcGarbageCollectionClient struct {
	cc drpc.Conn
}

func NewDRPCGarbageCollectionClient(cc drpc.Conn) DRPCGarbageCollectionClient {
	return &drpcGarbageCollectionClient{cc}
}

func (c *drpcGarbageCollectionClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcGarbageCollectionClient) Ping(ctx context.Context, in *PingRequest) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/gc.GarbageCollection/Ping", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcGarbageCollectionClient) StartSession(ctx context.Context, in *StartSessionRequest) (*StartSessionResponse, error) {
	out := new(StartSessionResponse)
	err := c.cc.Invoke(ctx, "/gc.GarbageCollection/StartSession", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcGarbageCollectionClient) AddPiece(ctx context.Context, in *AddPieceRequest) (*AddPieceResponse, error) {
	out := new(AddPieceResponse)
	err := c.cc.Invoke(ctx, "/gc.GarbageCollection/AddPiece", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcGarbageCollectionClient) EndSession(ctx context.Context, in *EndSessionRequest) (*EndSessionResponse, error) {
	out := new(EndSessionResponse)
	err := c.cc.Invoke(ctx, "/gc.GarbageCollection/EndSession", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCGarbageCollectionServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	// StartSession begins a new garbage collection session
	StartSession(context.Context, *StartSessionRequest) (*StartSessionResponse, error)
	// AddPiece adds a piece ID for a specifc storage node along with a
	// sequence number that helps ensuring no pieces were missed
	AddPiece(context.Context, *AddPieceRequest) (*AddPieceResponse, error)
	// EndSession ends the garbage collection session indicating how many pieces
	// should have been processed for each storage node
	EndSession(context.Context, *EndSessionRequest) (*EndSessionResponse, error)
}

type DRPCGarbageCollectionDescription struct{}

func (DRPCGarbageCollectionDescription) NumMethods() int { return 4 }

func (DRPCGarbageCollectionDescription) Method(n int) (string, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/gc.GarbageCollection/Ping",
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCGarbageCollectionServer).
					Ping(
						ctx,
						in1.(*PingRequest),
					)
			}, DRPCGarbageCollectionServer.Ping, true
	case 1:
		return "/gc.GarbageCollection/StartSession",
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCGarbageCollectionServer).
					StartSession(
						ctx,
						in1.(*StartSessionRequest),
					)
			}, DRPCGarbageCollectionServer.StartSession, true
	case 2:
		return "/gc.GarbageCollection/AddPiece",
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCGarbageCollectionServer).
					AddPiece(
						ctx,
						in1.(*AddPieceRequest),
					)
			}, DRPCGarbageCollectionServer.AddPiece, true
	case 3:
		return "/gc.GarbageCollection/EndSession",
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCGarbageCollectionServer).
					EndSession(
						ctx,
						in1.(*EndSessionRequest),
					)
			}, DRPCGarbageCollectionServer.EndSession, true
	default:
		return "", nil, nil, false
	}
}

func DRPCRegisterGarbageCollection(mux drpc.Mux, impl DRPCGarbageCollectionServer) error {
	return mux.Register(impl, DRPCGarbageCollectionDescription{})
}

type DRPCGarbageCollection_PingStream interface {
	drpc.Stream
	SendAndClose(*PingResponse) error
}

type drpcGarbageCollectionPingStream struct {
	drpc.Stream
}

func (x *drpcGarbageCollectionPingStream) SendAndClose(m *PingResponse) error {
	if err := x.MsgSend(m); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCGarbageCollection_StartSessionStream interface {
	drpc.Stream
	SendAndClose(*StartSessionResponse) error
}

type drpcGarbageCollectionStartSessionStream struct {
	drpc.Stream
}

func (x *drpcGarbageCollectionStartSessionStream) SendAndClose(m *StartSessionResponse) error {
	if err := x.MsgSend(m); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCGarbageCollection_AddPieceStream interface {
	drpc.Stream
	SendAndClose(*AddPieceResponse) error
}

type drpcGarbageCollectionAddPieceStream struct {
	drpc.Stream
}

func (x *drpcGarbageCollectionAddPieceStream) SendAndClose(m *AddPieceResponse) error {
	if err := x.MsgSend(m); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCGarbageCollection_EndSessionStream interface {
	drpc.Stream
	SendAndClose(*EndSessionResponse) error
}

type drpcGarbageCollectionEndSessionStream struct {
	drpc.Stream
}

func (x *drpcGarbageCollectionEndSessionStream) SendAndClose(m *EndSessionResponse) error {
	if err := x.MsgSend(m); err != nil {
		return err
	}
	return x.CloseSend()
}

// --- DRPC END ---
