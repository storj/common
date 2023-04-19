// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: piecestore2.proto

package pb

import (
	fmt "fmt"
	math "math"
	time "time"

	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type PieceHeader_FormatVersion int32

const (
	PieceHeader_FORMAT_V0 PieceHeader_FormatVersion = 0
	PieceHeader_FORMAT_V1 PieceHeader_FormatVersion = 1
)

var PieceHeader_FormatVersion_name = map[int32]string{
	0: "FORMAT_V0",
	1: "FORMAT_V1",
}

var PieceHeader_FormatVersion_value = map[string]int32{
	"FORMAT_V0": 0,
	"FORMAT_V1": 1,
}

func (x PieceHeader_FormatVersion) String() string {
	return proto.EnumName(PieceHeader_FormatVersion_name, int32(x))
}

func (PieceHeader_FormatVersion) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{12, 0}
}

// Expected order of messages from uplink:
//
//   OrderLimit ->
//   repeated
//      Order ->
//      Chunk ->
//   PieceHash signed by uplink ->
//      <- PieceHash signed by storage node
type PieceUploadRequest struct {
	// first message to show that we are allowed to upload
	Limit *OrderLimit `protobuf:"bytes,1,opt,name=limit,proto3" json:"limit,omitempty"`
	// first message must have it if (!) not the default sha256 is used, as it
	// should be initialized by the storagenode before upload.
	// should match with the algorithm in the done field of the last message
	HashAlgorithm PieceHashAlgorithm `protobuf:"varint,5,opt,name=hash_algorithm,json=hashAlgorithm,proto3,enum=orders.PieceHashAlgorithm" json:"hash_algorithm,omitempty"`
	// order for uploading
	Order *Order                    `protobuf:"bytes,2,opt,name=order,proto3" json:"order,omitempty"`
	Chunk *PieceUploadRequest_Chunk `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
	// final message
	Done                 *PieceHash `protobuf:"bytes,4,opt,name=done,proto3" json:"done,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *PieceUploadRequest) Reset()         { *m = PieceUploadRequest{} }
func (m *PieceUploadRequest) String() string { return proto.CompactTextString(m) }
func (*PieceUploadRequest) ProtoMessage()    {}
func (*PieceUploadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{0}
}
func (m *PieceUploadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceUploadRequest.Unmarshal(m, b)
}
func (m *PieceUploadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceUploadRequest.Marshal(b, m, deterministic)
}
func (m *PieceUploadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceUploadRequest.Merge(m, src)
}
func (m *PieceUploadRequest) XXX_Size() int {
	return xxx_messageInfo_PieceUploadRequest.Size(m)
}
func (m *PieceUploadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceUploadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PieceUploadRequest proto.InternalMessageInfo

func (m *PieceUploadRequest) GetLimit() *OrderLimit {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *PieceUploadRequest) GetHashAlgorithm() PieceHashAlgorithm {
	if m != nil {
		return m.HashAlgorithm
	}
	return PieceHashAlgorithm_SHA256
}

func (m *PieceUploadRequest) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *PieceUploadRequest) GetChunk() *PieceUploadRequest_Chunk {
	if m != nil {
		return m.Chunk
	}
	return nil
}

func (m *PieceUploadRequest) GetDone() *PieceHash {
	if m != nil {
		return m.Done
	}
	return nil
}

// data message
type PieceUploadRequest_Chunk struct {
	Offset               int64    `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PieceUploadRequest_Chunk) Reset()         { *m = PieceUploadRequest_Chunk{} }
func (m *PieceUploadRequest_Chunk) String() string { return proto.CompactTextString(m) }
func (*PieceUploadRequest_Chunk) ProtoMessage()    {}
func (*PieceUploadRequest_Chunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{0, 0}
}
func (m *PieceUploadRequest_Chunk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceUploadRequest_Chunk.Unmarshal(m, b)
}
func (m *PieceUploadRequest_Chunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceUploadRequest_Chunk.Marshal(b, m, deterministic)
}
func (m *PieceUploadRequest_Chunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceUploadRequest_Chunk.Merge(m, src)
}
func (m *PieceUploadRequest_Chunk) XXX_Size() int {
	return xxx_messageInfo_PieceUploadRequest_Chunk.Size(m)
}
func (m *PieceUploadRequest_Chunk) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceUploadRequest_Chunk.DiscardUnknown(m)
}

var xxx_messageInfo_PieceUploadRequest_Chunk proto.InternalMessageInfo

func (m *PieceUploadRequest_Chunk) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *PieceUploadRequest_Chunk) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type PieceUploadResponse struct {
	Done *PieceHash `protobuf:"bytes,1,opt,name=done,proto3" json:"done,omitempty"`
	// this is for validating the PieceHash signature if the cert chain is
	// unable to be pulled off the connection (Noise instead of TLS).
	NodeCertchain        []byte   `protobuf:"bytes,2,opt,name=node_certchain,json=nodeCertchain,proto3" json:"node_certchain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PieceUploadResponse) Reset()         { *m = PieceUploadResponse{} }
func (m *PieceUploadResponse) String() string { return proto.CompactTextString(m) }
func (*PieceUploadResponse) ProtoMessage()    {}
func (*PieceUploadResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{1}
}
func (m *PieceUploadResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceUploadResponse.Unmarshal(m, b)
}
func (m *PieceUploadResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceUploadResponse.Marshal(b, m, deterministic)
}
func (m *PieceUploadResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceUploadResponse.Merge(m, src)
}
func (m *PieceUploadResponse) XXX_Size() int {
	return xxx_messageInfo_PieceUploadResponse.Size(m)
}
func (m *PieceUploadResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceUploadResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PieceUploadResponse proto.InternalMessageInfo

func (m *PieceUploadResponse) GetDone() *PieceHash {
	if m != nil {
		return m.Done
	}
	return nil
}

func (m *PieceUploadResponse) GetNodeCertchain() []byte {
	if m != nil {
		return m.NodeCertchain
	}
	return nil
}

// Expected order of messages from uplink:
//
//   {OrderLimit, Chunk} ->
//   go repeated
//      Order -> (async)
//   go repeated
//      <- PieceDownloadResponse.Chunk
type PieceDownloadRequest struct {
	// first message to show that we are allowed to upload
	Limit *OrderLimit `protobuf:"bytes,1,opt,name=limit,proto3" json:"limit,omitempty"`
	// order for downloading
	Order *Order `protobuf:"bytes,2,opt,name=order,proto3" json:"order,omitempty"`
	// request for the chunk
	Chunk                *PieceDownloadRequest_Chunk `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *PieceDownloadRequest) Reset()         { *m = PieceDownloadRequest{} }
func (m *PieceDownloadRequest) String() string { return proto.CompactTextString(m) }
func (*PieceDownloadRequest) ProtoMessage()    {}
func (*PieceDownloadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{2}
}
func (m *PieceDownloadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceDownloadRequest.Unmarshal(m, b)
}
func (m *PieceDownloadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceDownloadRequest.Marshal(b, m, deterministic)
}
func (m *PieceDownloadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceDownloadRequest.Merge(m, src)
}
func (m *PieceDownloadRequest) XXX_Size() int {
	return xxx_messageInfo_PieceDownloadRequest.Size(m)
}
func (m *PieceDownloadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceDownloadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PieceDownloadRequest proto.InternalMessageInfo

func (m *PieceDownloadRequest) GetLimit() *OrderLimit {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *PieceDownloadRequest) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *PieceDownloadRequest) GetChunk() *PieceDownloadRequest_Chunk {
	if m != nil {
		return m.Chunk
	}
	return nil
}

// Chunk that we wish to download
type PieceDownloadRequest_Chunk struct {
	Offset               int64    `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	ChunkSize            int64    `protobuf:"varint,2,opt,name=chunk_size,json=chunkSize,proto3" json:"chunk_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PieceDownloadRequest_Chunk) Reset()         { *m = PieceDownloadRequest_Chunk{} }
func (m *PieceDownloadRequest_Chunk) String() string { return proto.CompactTextString(m) }
func (*PieceDownloadRequest_Chunk) ProtoMessage()    {}
func (*PieceDownloadRequest_Chunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{2, 0}
}
func (m *PieceDownloadRequest_Chunk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceDownloadRequest_Chunk.Unmarshal(m, b)
}
func (m *PieceDownloadRequest_Chunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceDownloadRequest_Chunk.Marshal(b, m, deterministic)
}
func (m *PieceDownloadRequest_Chunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceDownloadRequest_Chunk.Merge(m, src)
}
func (m *PieceDownloadRequest_Chunk) XXX_Size() int {
	return xxx_messageInfo_PieceDownloadRequest_Chunk.Size(m)
}
func (m *PieceDownloadRequest_Chunk) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceDownloadRequest_Chunk.DiscardUnknown(m)
}

var xxx_messageInfo_PieceDownloadRequest_Chunk proto.InternalMessageInfo

func (m *PieceDownloadRequest_Chunk) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *PieceDownloadRequest_Chunk) GetChunkSize() int64 {
	if m != nil {
		return m.ChunkSize
	}
	return 0
}

type PieceDownloadResponse struct {
	Chunk                *PieceDownloadResponse_Chunk `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
	Hash                 *PieceHash                   `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Limit                *OrderLimit                  `protobuf:"bytes,3,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *PieceDownloadResponse) Reset()         { *m = PieceDownloadResponse{} }
func (m *PieceDownloadResponse) String() string { return proto.CompactTextString(m) }
func (*PieceDownloadResponse) ProtoMessage()    {}
func (*PieceDownloadResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{3}
}
func (m *PieceDownloadResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceDownloadResponse.Unmarshal(m, b)
}
func (m *PieceDownloadResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceDownloadResponse.Marshal(b, m, deterministic)
}
func (m *PieceDownloadResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceDownloadResponse.Merge(m, src)
}
func (m *PieceDownloadResponse) XXX_Size() int {
	return xxx_messageInfo_PieceDownloadResponse.Size(m)
}
func (m *PieceDownloadResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceDownloadResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PieceDownloadResponse proto.InternalMessageInfo

func (m *PieceDownloadResponse) GetChunk() *PieceDownloadResponse_Chunk {
	if m != nil {
		return m.Chunk
	}
	return nil
}

func (m *PieceDownloadResponse) GetHash() *PieceHash {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *PieceDownloadResponse) GetLimit() *OrderLimit {
	if m != nil {
		return m.Limit
	}
	return nil
}

// Chunk response for download request
type PieceDownloadResponse_Chunk struct {
	Offset               int64    `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PieceDownloadResponse_Chunk) Reset()         { *m = PieceDownloadResponse_Chunk{} }
func (m *PieceDownloadResponse_Chunk) String() string { return proto.CompactTextString(m) }
func (*PieceDownloadResponse_Chunk) ProtoMessage()    {}
func (*PieceDownloadResponse_Chunk) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{3, 0}
}
func (m *PieceDownloadResponse_Chunk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceDownloadResponse_Chunk.Unmarshal(m, b)
}
func (m *PieceDownloadResponse_Chunk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceDownloadResponse_Chunk.Marshal(b, m, deterministic)
}
func (m *PieceDownloadResponse_Chunk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceDownloadResponse_Chunk.Merge(m, src)
}
func (m *PieceDownloadResponse_Chunk) XXX_Size() int {
	return xxx_messageInfo_PieceDownloadResponse_Chunk.Size(m)
}
func (m *PieceDownloadResponse_Chunk) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceDownloadResponse_Chunk.DiscardUnknown(m)
}

var xxx_messageInfo_PieceDownloadResponse_Chunk proto.InternalMessageInfo

func (m *PieceDownloadResponse_Chunk) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *PieceDownloadResponse_Chunk) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type PieceDeleteRequest struct {
	Limit                *OrderLimit `protobuf:"bytes,1,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PieceDeleteRequest) Reset()         { *m = PieceDeleteRequest{} }
func (m *PieceDeleteRequest) String() string { return proto.CompactTextString(m) }
func (*PieceDeleteRequest) ProtoMessage()    {}
func (*PieceDeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{4}
}
func (m *PieceDeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceDeleteRequest.Unmarshal(m, b)
}
func (m *PieceDeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceDeleteRequest.Marshal(b, m, deterministic)
}
func (m *PieceDeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceDeleteRequest.Merge(m, src)
}
func (m *PieceDeleteRequest) XXX_Size() int {
	return xxx_messageInfo_PieceDeleteRequest.Size(m)
}
func (m *PieceDeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceDeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PieceDeleteRequest proto.InternalMessageInfo

func (m *PieceDeleteRequest) GetLimit() *OrderLimit {
	if m != nil {
		return m.Limit
	}
	return nil
}

type PieceDeleteResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PieceDeleteResponse) Reset()         { *m = PieceDeleteResponse{} }
func (m *PieceDeleteResponse) String() string { return proto.CompactTextString(m) }
func (*PieceDeleteResponse) ProtoMessage()    {}
func (*PieceDeleteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{5}
}
func (m *PieceDeleteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceDeleteResponse.Unmarshal(m, b)
}
func (m *PieceDeleteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceDeleteResponse.Marshal(b, m, deterministic)
}
func (m *PieceDeleteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceDeleteResponse.Merge(m, src)
}
func (m *PieceDeleteResponse) XXX_Size() int {
	return xxx_messageInfo_PieceDeleteResponse.Size(m)
}
func (m *PieceDeleteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceDeleteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PieceDeleteResponse proto.InternalMessageInfo

type DeletePiecesRequest struct {
	PieceIds             []PieceID `protobuf:"bytes,1,rep,name=piece_ids,json=pieceIds,proto3,customtype=PieceID" json:"piece_ids"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *DeletePiecesRequest) Reset()         { *m = DeletePiecesRequest{} }
func (m *DeletePiecesRequest) String() string { return proto.CompactTextString(m) }
func (*DeletePiecesRequest) ProtoMessage()    {}
func (*DeletePiecesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{6}
}
func (m *DeletePiecesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletePiecesRequest.Unmarshal(m, b)
}
func (m *DeletePiecesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletePiecesRequest.Marshal(b, m, deterministic)
}
func (m *DeletePiecesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletePiecesRequest.Merge(m, src)
}
func (m *DeletePiecesRequest) XXX_Size() int {
	return xxx_messageInfo_DeletePiecesRequest.Size(m)
}
func (m *DeletePiecesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletePiecesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeletePiecesRequest proto.InternalMessageInfo

type DeletePiecesResponse struct {
	UnhandledCount       int64    `protobuf:"varint,1,opt,name=unhandled_count,json=unhandledCount,proto3" json:"unhandled_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeletePiecesResponse) Reset()         { *m = DeletePiecesResponse{} }
func (m *DeletePiecesResponse) String() string { return proto.CompactTextString(m) }
func (*DeletePiecesResponse) ProtoMessage()    {}
func (*DeletePiecesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{7}
}
func (m *DeletePiecesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletePiecesResponse.Unmarshal(m, b)
}
func (m *DeletePiecesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletePiecesResponse.Marshal(b, m, deterministic)
}
func (m *DeletePiecesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletePiecesResponse.Merge(m, src)
}
func (m *DeletePiecesResponse) XXX_Size() int {
	return xxx_messageInfo_DeletePiecesResponse.Size(m)
}
func (m *DeletePiecesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletePiecesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeletePiecesResponse proto.InternalMessageInfo

func (m *DeletePiecesResponse) GetUnhandledCount() int64 {
	if m != nil {
		return m.UnhandledCount
	}
	return 0
}

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
	return fileDescriptor_23ff32dd550c2439, []int{8}
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
	return fileDescriptor_23ff32dd550c2439, []int{9}
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

type RestoreTrashRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RestoreTrashRequest) Reset()         { *m = RestoreTrashRequest{} }
func (m *RestoreTrashRequest) String() string { return proto.CompactTextString(m) }
func (*RestoreTrashRequest) ProtoMessage()    {}
func (*RestoreTrashRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{10}
}
func (m *RestoreTrashRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RestoreTrashRequest.Unmarshal(m, b)
}
func (m *RestoreTrashRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RestoreTrashRequest.Marshal(b, m, deterministic)
}
func (m *RestoreTrashRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RestoreTrashRequest.Merge(m, src)
}
func (m *RestoreTrashRequest) XXX_Size() int {
	return xxx_messageInfo_RestoreTrashRequest.Size(m)
}
func (m *RestoreTrashRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RestoreTrashRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RestoreTrashRequest proto.InternalMessageInfo

type RestoreTrashResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RestoreTrashResponse) Reset()         { *m = RestoreTrashResponse{} }
func (m *RestoreTrashResponse) String() string { return proto.CompactTextString(m) }
func (*RestoreTrashResponse) ProtoMessage()    {}
func (*RestoreTrashResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{11}
}
func (m *RestoreTrashResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RestoreTrashResponse.Unmarshal(m, b)
}
func (m *RestoreTrashResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RestoreTrashResponse.Marshal(b, m, deterministic)
}
func (m *RestoreTrashResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RestoreTrashResponse.Merge(m, src)
}
func (m *RestoreTrashResponse) XXX_Size() int {
	return xxx_messageInfo_RestoreTrashResponse.Size(m)
}
func (m *RestoreTrashResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RestoreTrashResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RestoreTrashResponse proto.InternalMessageInfo

// PieceHeader is used in piece storage to keep track of piece attributes.
type PieceHeader struct {
	// the storage format version being used for this piece. The piece filename should agree with this.
	// The inclusion of this field is intended to aid repairability when filenames are damaged.
	FormatVersion PieceHeader_FormatVersion `protobuf:"varint,1,opt,name=format_version,json=formatVersion,proto3,enum=piecestore.PieceHeader_FormatVersion" json:"format_version,omitempty"`
	// content hash of the piece
	Hash []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	// the algorithm of the hash
	HashAlgorithm PieceHashAlgorithm `protobuf:"varint,6,opt,name=hash_algorithm,json=hashAlgorithm,proto3,enum=orders.PieceHashAlgorithm" json:"hash_algorithm,omitempty"`
	// timestamp when upload occurred, as given by the "timestamp" field in the original orders.PieceHash
	CreationTime time.Time `protobuf:"bytes,3,opt,name=creation_time,json=creationTime,proto3,stdtime" json:"creation_time"`
	// signature from uplink over the original orders.PieceHash (the corresponding PieceHashSigning
	// is reconstructable using the piece id from the piecestore, the piece size from the
	// filesystem (minus the piece header size), and these (hash, upload_time, signature) fields).
	Signature []byte `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	// the OrderLimit authorizing storage of this piece, as signed by the satellite and sent by
	// the uplink
	OrderLimit           OrderLimit `protobuf:"bytes,5,opt,name=order_limit,json=orderLimit,proto3" json:"order_limit"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *PieceHeader) Reset()         { *m = PieceHeader{} }
func (m *PieceHeader) String() string { return proto.CompactTextString(m) }
func (*PieceHeader) ProtoMessage()    {}
func (*PieceHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{12}
}
func (m *PieceHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceHeader.Unmarshal(m, b)
}
func (m *PieceHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceHeader.Marshal(b, m, deterministic)
}
func (m *PieceHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceHeader.Merge(m, src)
}
func (m *PieceHeader) XXX_Size() int {
	return xxx_messageInfo_PieceHeader.Size(m)
}
func (m *PieceHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceHeader.DiscardUnknown(m)
}

var xxx_messageInfo_PieceHeader proto.InternalMessageInfo

func (m *PieceHeader) GetFormatVersion() PieceHeader_FormatVersion {
	if m != nil {
		return m.FormatVersion
	}
	return PieceHeader_FORMAT_V0
}

func (m *PieceHeader) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *PieceHeader) GetHashAlgorithm() PieceHashAlgorithm {
	if m != nil {
		return m.HashAlgorithm
	}
	return PieceHashAlgorithm_SHA256
}

func (m *PieceHeader) GetCreationTime() time.Time {
	if m != nil {
		return m.CreationTime
	}
	return time.Time{}
}

func (m *PieceHeader) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *PieceHeader) GetOrderLimit() OrderLimit {
	if m != nil {
		return m.OrderLimit
	}
	return OrderLimit{}
}

type ExistsRequest struct {
	PieceIds             []PieceID `protobuf:"bytes,1,rep,name=piece_ids,json=pieceIds,proto3,customtype=PieceID" json:"piece_ids"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ExistsRequest) Reset()         { *m = ExistsRequest{} }
func (m *ExistsRequest) String() string { return proto.CompactTextString(m) }
func (*ExistsRequest) ProtoMessage()    {}
func (*ExistsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{13}
}
func (m *ExistsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExistsRequest.Unmarshal(m, b)
}
func (m *ExistsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExistsRequest.Marshal(b, m, deterministic)
}
func (m *ExistsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExistsRequest.Merge(m, src)
}
func (m *ExistsRequest) XXX_Size() int {
	return xxx_messageInfo_ExistsRequest.Size(m)
}
func (m *ExistsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExistsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExistsRequest proto.InternalMessageInfo

type ExistsResponse struct {
	// input piece ids indices of the missing pieces
	Missing              []uint32 `protobuf:"varint,1,rep,packed,name=missing,proto3" json:"missing,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExistsResponse) Reset()         { *m = ExistsResponse{} }
func (m *ExistsResponse) String() string { return proto.CompactTextString(m) }
func (*ExistsResponse) ProtoMessage()    {}
func (*ExistsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_23ff32dd550c2439, []int{14}
}
func (m *ExistsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExistsResponse.Unmarshal(m, b)
}
func (m *ExistsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExistsResponse.Marshal(b, m, deterministic)
}
func (m *ExistsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExistsResponse.Merge(m, src)
}
func (m *ExistsResponse) XXX_Size() int {
	return xxx_messageInfo_ExistsResponse.Size(m)
}
func (m *ExistsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExistsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExistsResponse proto.InternalMessageInfo

func (m *ExistsResponse) GetMissing() []uint32 {
	if m != nil {
		return m.Missing
	}
	return nil
}

func init() {
	proto.RegisterEnum("piecestore.PieceHeader_FormatVersion", PieceHeader_FormatVersion_name, PieceHeader_FormatVersion_value)
	proto.RegisterType((*PieceUploadRequest)(nil), "piecestore.PieceUploadRequest")
	proto.RegisterType((*PieceUploadRequest_Chunk)(nil), "piecestore.PieceUploadRequest.Chunk")
	proto.RegisterType((*PieceUploadResponse)(nil), "piecestore.PieceUploadResponse")
	proto.RegisterType((*PieceDownloadRequest)(nil), "piecestore.PieceDownloadRequest")
	proto.RegisterType((*PieceDownloadRequest_Chunk)(nil), "piecestore.PieceDownloadRequest.Chunk")
	proto.RegisterType((*PieceDownloadResponse)(nil), "piecestore.PieceDownloadResponse")
	proto.RegisterType((*PieceDownloadResponse_Chunk)(nil), "piecestore.PieceDownloadResponse.Chunk")
	proto.RegisterType((*PieceDeleteRequest)(nil), "piecestore.PieceDeleteRequest")
	proto.RegisterType((*PieceDeleteResponse)(nil), "piecestore.PieceDeleteResponse")
	proto.RegisterType((*DeletePiecesRequest)(nil), "piecestore.DeletePiecesRequest")
	proto.RegisterType((*DeletePiecesResponse)(nil), "piecestore.DeletePiecesResponse")
	proto.RegisterType((*RetainRequest)(nil), "piecestore.RetainRequest")
	proto.RegisterType((*RetainResponse)(nil), "piecestore.RetainResponse")
	proto.RegisterType((*RestoreTrashRequest)(nil), "piecestore.RestoreTrashRequest")
	proto.RegisterType((*RestoreTrashResponse)(nil), "piecestore.RestoreTrashResponse")
	proto.RegisterType((*PieceHeader)(nil), "piecestore.PieceHeader")
	proto.RegisterType((*ExistsRequest)(nil), "piecestore.ExistsRequest")
	proto.RegisterType((*ExistsResponse)(nil), "piecestore.ExistsResponse")
}

func init() { proto.RegisterFile("piecestore2.proto", fileDescriptor_23ff32dd550c2439) }

var fileDescriptor_23ff32dd550c2439 = []byte{
	// 888 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x4d, 0x6f, 0x1b, 0x45,
	0x18, 0xce, 0xc6, 0x1f, 0x4d, 0xde, 0xec, 0xba, 0xe9, 0x24, 0xa9, 0xcc, 0x0a, 0xb0, 0x59, 0x08,
	0xb5, 0x10, 0x6c, 0x8a, 0x7b, 0x02, 0xd1, 0x56, 0xf9, 0xa0, 0x22, 0x52, 0x4b, 0xcb, 0x34, 0xed,
	0x81, 0xcb, 0x6a, 0xe3, 0x1d, 0xdb, 0x03, 0xde, 0x1d, 0xb3, 0x33, 0x06, 0xd4, 0x3b, 0x12, 0x47,
	0x7e, 0x08, 0xbf, 0x81, 0x33, 0xbf, 0x01, 0xa1, 0x70, 0xe0, 0x8f, 0xa0, 0xf9, 0x72, 0x3c, 0xb1,
	0x13, 0x93, 0x9c, 0xec, 0x79, 0x3f, 0x9f, 0x79, 0xe6, 0x79, 0xdf, 0x85, 0x3b, 0x63, 0x4a, 0x7a,
	0x84, 0x0b, 0x56, 0x92, 0x6e, 0x3c, 0x2e, 0x99, 0x60, 0x08, 0xce, 0x4d, 0x21, 0x0c, 0xd8, 0x80,
	0x69, 0x7b, 0xd8, 0x1a, 0x30, 0x36, 0x18, 0x91, 0x3d, 0x75, 0x3a, 0x9d, 0xf4, 0xf7, 0x04, 0xcd,
	0x09, 0x17, 0x69, 0x3e, 0x36, 0x01, 0x3e, 0x2b, 0x33, 0x52, 0x72, 0x7d, 0x8a, 0xfe, 0x58, 0x05,
	0xf4, 0x42, 0x56, 0x7a, 0x35, 0x1e, 0xb1, 0x34, 0xc3, 0xe4, 0x87, 0x09, 0xe1, 0x02, 0x75, 0xa0,
	0x36, 0xa2, 0x39, 0x15, 0x4d, 0xaf, 0xed, 0x75, 0x36, 0xba, 0x28, 0x36, 0x49, 0xcf, 0xe5, 0xcf,
	0x53, 0xe9, 0xc1, 0x3a, 0x00, 0xed, 0x43, 0x63, 0x98, 0xf2, 0x61, 0x92, 0x8e, 0x06, 0xac, 0xa4,
	0x62, 0x98, 0x37, 0x6b, 0x6d, 0xaf, 0xd3, 0xe8, 0x86, 0x36, 0x45, 0x55, 0xff, 0x2a, 0xe5, 0xc3,
	0x7d, 0x1b, 0x81, 0x83, 0xe1, 0xec, 0x11, 0xbd, 0x0f, 0x35, 0x15, 0xdb, 0x5c, 0x55, 0xcd, 0x02,
	0xa7, 0x19, 0xd6, 0x3e, 0xf4, 0x39, 0xd4, 0x7a, 0xc3, 0x49, 0xf1, 0x7d, 0xb3, 0xa2, 0x82, 0x3e,
	0x88, 0xcf, 0xef, 0x1f, 0xcf, 0x5f, 0x20, 0x3e, 0x94, 0xb1, 0x58, 0xa7, 0xa0, 0x5d, 0xa8, 0x66,
	0xac, 0x20, 0xcd, 0xaa, 0x4a, 0xbd, 0x33, 0x87, 0x0c, 0x2b, 0x77, 0xf8, 0x00, 0x6a, 0x2a, 0x0d,
	0xdd, 0x85, 0x3a, 0xeb, 0xf7, 0x39, 0xd1, 0xd7, 0xaf, 0x60, 0x73, 0x42, 0x08, 0xaa, 0x59, 0x2a,
	0x52, 0x85, 0xd3, 0xc7, 0xea, 0x7f, 0xd4, 0x83, 0x2d, 0xa7, 0x3d, 0x1f, 0xb3, 0x82, 0x93, 0x69,
	0x4b, 0xef, 0xca, 0x96, 0x68, 0x17, 0x1a, 0x05, 0xcb, 0x48, 0xd2, 0x23, 0xa5, 0xe8, 0x0d, 0x53,
	0x5a, 0x98, 0xda, 0x81, 0xb4, 0x1e, 0x5a, 0x63, 0xf4, 0xaf, 0x07, 0xdb, 0x2a, 0xf5, 0x88, 0xfd,
	0x54, 0xdc, 0xec, 0x9d, 0xfe, 0x17, 0xc9, 0x5f, 0xb8, 0x24, 0x7f, 0x38, 0x47, 0xf2, 0x85, 0xfe,
	0x0e, 0xcd, 0xe1, 0xa3, 0x65, 0xfc, 0xbd, 0x03, 0xa0, 0x22, 0x13, 0x4e, 0xdf, 0x10, 0x05, 0xa4,
	0x82, 0xd7, 0x95, 0xe5, 0x25, 0x7d, 0x43, 0xa2, 0xbf, 0x3d, 0xd8, 0xb9, 0xd0, 0xc5, 0xb0, 0xf9,
	0xd0, 0xe2, 0xd2, 0xd7, 0xbc, 0x77, 0x05, 0x2e, 0x9d, 0x31, 0xf7, 0xfe, 0x52, 0x71, 0xe6, 0xea,
	0x8b, 0x1e, 0x43, 0xba, 0xcf, 0xc9, 0xac, 0x2c, 0x21, 0xf3, 0x66, 0x4a, 0x79, 0x64, 0x26, 0xed,
	0x88, 0x8c, 0x88, 0x20, 0xd7, 0x7e, 0xc1, 0x68, 0xc7, 0x28, 0xcd, 0xe6, 0xeb, 0x9b, 0x46, 0x87,
	0xb0, 0xa5, 0x2d, 0xca, 0xc9, 0x6d, 0xdd, 0x8f, 0x61, 0x5d, 0x91, 0x94, 0xd0, 0x8c, 0x37, 0xbd,
	0x76, 0xa5, 0xe3, 0x1f, 0xdc, 0xfe, 0xf3, 0xac, 0xb5, 0xf2, 0xd7, 0x59, 0xeb, 0x96, 0x8a, 0x3c,
	0x3e, 0xc2, 0x6b, 0x2a, 0xe2, 0x38, 0xe3, 0xd1, 0x63, 0xd8, 0x76, 0x8b, 0x18, 0xe2, 0xef, 0xc1,
	0xed, 0x49, 0x31, 0x4c, 0x8b, 0x6c, 0x44, 0xb2, 0xa4, 0xc7, 0x26, 0x85, 0xbd, 0x68, 0x63, 0x6a,
	0x3e, 0x94, 0xd6, 0xa8, 0x84, 0x00, 0x13, 0x91, 0xd2, 0xc2, 0xf6, 0x3f, 0x86, 0xa0, 0x57, 0x92,
	0x54, 0x50, 0x56, 0x24, 0x59, 0x2a, 0xec, 0x24, 0x84, 0xb1, 0xde, 0x4f, 0xb1, 0xdd, 0x4f, 0xf1,
	0x89, 0xdd, 0x4f, 0x07, 0x6b, 0x12, 0xdf, 0x6f, 0xff, 0xb4, 0x3c, 0xec, 0xdb, 0xd4, 0xa3, 0x54,
	0x10, 0x49, 0x72, 0x9f, 0x8e, 0x84, 0xd1, 0xae, 0x8f, 0xcd, 0x29, 0xda, 0x84, 0x86, 0xed, 0x69,
	0xb8, 0xd8, 0x81, 0x2d, 0xac, 0x65, 0x71, 0x52, 0xca, 0x77, 0xd5, 0x58, 0xa2, 0xbb, 0xb0, 0xed,
	0x9a, 0x4d, 0xf8, 0x2f, 0x15, 0xd8, 0xd0, 0x22, 0x20, 0xa9, 0x94, 0xff, 0x53, 0x68, 0xf4, 0x59,
	0x99, 0xa7, 0x22, 0xf9, 0x91, 0x94, 0x9c, 0xb2, 0x42, 0x81, 0x6e, 0x74, 0x77, 0xe7, 0xf4, 0xa6,
	0x13, 0xe2, 0x27, 0x2a, 0xfa, 0xb5, 0x0e, 0xc6, 0x41, 0x7f, 0xf6, 0x28, 0x35, 0x30, 0x55, 0x9d,
	0x6f, 0x24, 0x36, 0xbf, 0x2d, 0xeb, 0xd7, 0xdd, 0x96, 0xb3, 0xc4, 0xca, 0xdd, 0x6e, 0xd4, 0x7a,
	0x4d, 0x62, 0xa5, 0x13, 0xbd, 0x0d, 0xeb, 0x9c, 0x0e, 0x8a, 0x54, 0x4c, 0x4a, 0xbd, 0x1c, 0x7d,
	0x7c, 0x6e, 0x40, 0x9f, 0xc1, 0x86, 0x02, 0x95, 0x68, 0x7d, 0xd6, 0x2e, 0xd3, 0xe7, 0x41, 0x55,
	0x96, 0xc7, 0xc0, 0xa6, 0x96, 0xe8, 0x13, 0x08, 0x1c, 0x6a, 0x50, 0x00, 0xeb, 0x4f, 0x9e, 0xe3,
	0x67, 0xfb, 0x27, 0xc9, 0xeb, 0xfb, 0x9b, 0x2b, 0xb3, 0xc7, 0x4f, 0x37, 0xbd, 0xe8, 0x21, 0x04,
	0x5f, 0xfe, 0x4c, 0xb9, 0xb8, 0xa1, 0x78, 0x3f, 0x82, 0x86, 0x4d, 0x37, 0xb2, 0x6d, 0xc2, 0xad,
	0x9c, 0x72, 0x4e, 0x8b, 0x81, 0xca, 0x0e, 0xb0, 0x3d, 0x76, 0x7f, 0xaf, 0x02, 0xbc, 0x98, 0x3e,
	0x26, 0x7a, 0x06, 0x75, 0xbd, 0xb8, 0xd1, 0xbb, 0x57, 0x7f, 0x50, 0xc2, 0xd6, 0xa5, 0x7e, 0x23,
	0xa6, 0x95, 0x8e, 0x87, 0x5e, 0xc1, 0x9a, 0xdd, 0x44, 0xa8, 0xbd, 0x6c, 0x79, 0x86, 0xef, 0x2d,
	0x5d, 0x63, 0xb2, 0xe8, 0x7d, 0x0f, 0x7d, 0x0d, 0x75, 0x3d, 0x9d, 0x0b, 0x50, 0x3a, 0xdb, 0x64,
	0x01, 0xca, 0x0b, 0xdb, 0xa2, 0xf2, 0xeb, 0xaa, 0x87, 0xbe, 0x01, 0x7f, 0x76, 0xda, 0x91, 0x93,
	0xb5, 0x60, 0x99, 0x84, 0xed, 0xcb, 0x03, 0x0c, 0xe3, 0x8f, 0xa1, 0xae, 0x67, 0x11, 0xbd, 0x35,
	0x1b, 0xeb, 0xec, 0x84, 0x30, 0x5c, 0xe4, 0x32, 0x05, 0x5e, 0x82, 0x3f, 0x3b, 0xa3, 0x2e, 0xa6,
	0x05, 0x43, 0xed, 0x62, 0x5a, 0x38, 0xde, 0x2b, 0x12, 0x95, 0x56, 0x86, 0x8b, 0xca, 0x11, 0x9b,
	0x8b, 0xca, 0x15, 0xd2, 0xc1, 0xf6, 0xb7, 0x48, 0xda, 0xbf, 0x8b, 0x29, 0xdb, 0xeb, 0xb1, 0x3c,
	0x67, 0xc5, 0xde, 0xf8, 0xf4, 0xb4, 0xae, 0x66, 0xec, 0xc1, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x30, 0x0c, 0x7c, 0x59, 0x97, 0x09, 0x00, 0x00,
}
