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
//	OrderLimit ->
//	repeated
//	   Order ->
//	   Chunk ->
//	PieceHash signed by uplink ->
//	   <- PieceHash signed by storage node
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
//	{OrderLimit, Chunk} ->
//	go repeated
//	   Order -> (async)
//	go repeated
//	   <- PieceDownloadResponse.Chunk
type PieceDownloadRequest struct {
	// first message to show that we are allowed to upload
	Limit *OrderLimit `protobuf:"bytes,1,opt,name=limit,proto3" json:"limit,omitempty"`
	// order for downloading
	Order *Order `protobuf:"bytes,2,opt,name=order,proto3" json:"order,omitempty"`
	// request for the chunk
	Chunk *PieceDownloadRequest_Chunk `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
	// maximum_chunk_size is an advisory request from the uplink to
	// the storage node on how big the chunks should be. smaller
	// maximum sizes may reduce time to first byte and peak memory usage.
	MaximumChunkSize     int32    `protobuf:"varint,4,opt,name=maximum_chunk_size,json=maximumChunkSize,proto3" json:"maximum_chunk_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
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

func (m *PieceDownloadRequest) GetMaximumChunkSize() int32 {
	if m != nil {
		return m.MaximumChunkSize
	}
	return 0
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
	// 907 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0xdd, 0x6e, 0x1b, 0x45,
	0x14, 0xce, 0xfa, 0xaf, 0xc9, 0x89, 0xd7, 0x4d, 0x27, 0x49, 0x65, 0x56, 0x80, 0xcd, 0x42, 0xa8,
	0x85, 0xca, 0xa6, 0xb8, 0x57, 0x20, 0xda, 0x2a, 0x3f, 0x54, 0x44, 0x6a, 0x69, 0x99, 0xa6, 0xbd,
	0xe0, 0x66, 0x35, 0xd9, 0x1d, 0xdb, 0x03, 0xde, 0x1d, 0xb3, 0x33, 0x86, 0xaa, 0xf7, 0x48, 0x5c,
	0xc2, 0x7b, 0xf0, 0x0c, 0x5c, 0xf3, 0x0c, 0x08, 0x95, 0x57, 0x41, 0xf3, 0xe7, 0x78, 0x63, 0x27,
	0x21, 0xb9, 0xb2, 0xe7, 0x9c, 0xf3, 0x9d, 0x73, 0xe6, 0x9b, 0xef, 0x9c, 0x85, 0x5b, 0x13, 0x46,
	0x13, 0x2a, 0x24, 0x2f, 0x68, 0x3f, 0x9a, 0x14, 0x5c, 0x72, 0x04, 0xa7, 0xa6, 0x00, 0x86, 0x7c,
	0xc8, 0x8d, 0x3d, 0xe8, 0x0c, 0x39, 0x1f, 0x8e, 0xe9, 0xae, 0x3e, 0x9d, 0x4c, 0x07, 0xbb, 0x92,
	0x65, 0x54, 0x48, 0x92, 0x4d, 0x6c, 0x40, 0x93, 0x17, 0x29, 0x2d, 0x84, 0x39, 0x85, 0x7f, 0x56,
	0x00, 0x3d, 0x57, 0x99, 0x5e, 0x4e, 0xc6, 0x9c, 0xa4, 0x98, 0xfe, 0x38, 0xa5, 0x42, 0xa2, 0x1e,
	0xd4, 0xc7, 0x2c, 0x63, 0xb2, 0xed, 0x75, 0xbd, 0xde, 0x7a, 0x1f, 0x45, 0x16, 0xf4, 0x4c, 0xfd,
	0x3c, 0x51, 0x1e, 0x6c, 0x02, 0xd0, 0x1e, 0xb4, 0x46, 0x44, 0x8c, 0x62, 0x32, 0x1e, 0xf2, 0x82,
	0xc9, 0x51, 0xd6, 0xae, 0x77, 0xbd, 0x5e, 0xab, 0x1f, 0x38, 0x88, 0xce, 0xfe, 0x35, 0x11, 0xa3,
	0x3d, 0x17, 0x81, 0xfd, 0xd1, 0xfc, 0x11, 0x7d, 0x08, 0x75, 0x1d, 0xdb, 0xae, 0xe8, 0x62, 0x7e,
	0xa9, 0x18, 0x36, 0x3e, 0xf4, 0x05, 0xd4, 0x93, 0xd1, 0x34, 0xff, 0xa1, 0x5d, 0xd5, 0x41, 0x1f,
	0x45, 0xa7, 0xf7, 0x8f, 0x16, 0x2f, 0x10, 0x1d, 0xa8, 0x58, 0x6c, 0x20, 0x68, 0x07, 0x6a, 0x29,
	0xcf, 0x69, 0xbb, 0xa6, 0xa1, 0xb7, 0x16, 0x3a, 0xc3, 0xda, 0x1d, 0xdc, 0x87, 0xba, 0x86, 0xa1,
	0xdb, 0xd0, 0xe0, 0x83, 0x81, 0xa0, 0xe6, 0xfa, 0x55, 0x6c, 0x4f, 0x08, 0x41, 0x2d, 0x25, 0x92,
	0xe8, 0x3e, 0x9b, 0x58, 0xff, 0x0f, 0x13, 0xd8, 0x2c, 0x95, 0x17, 0x13, 0x9e, 0x0b, 0x3a, 0x2b,
	0xe9, 0x5d, 0x58, 0x12, 0xed, 0x40, 0x2b, 0xe7, 0x29, 0x8d, 0x13, 0x5a, 0xc8, 0x64, 0x44, 0x58,
	0x6e, 0x73, 0xfb, 0xca, 0x7a, 0xe0, 0x8c, 0xe1, 0xef, 0x15, 0xd8, 0xd2, 0xd0, 0x43, 0xfe, 0x73,
	0x7e, 0xbd, 0x77, 0xfa, 0x5f, 0x24, 0x7f, 0x59, 0x26, 0xf9, 0xe3, 0x05, 0x92, 0xcf, 0xd4, 0x2f,
	0xd3, 0x7c, 0x17, 0x50, 0x46, 0x5e, 0xb3, 0x6c, 0x9a, 0xc5, 0xda, 0x10, 0x0b, 0xf6, 0xc6, 0x90,
	0x5e, 0xc7, 0x1b, 0xd6, 0xa3, 0x01, 0x2f, 0xd8, 0x1b, 0x1a, 0x3c, 0xbc, 0x8c, 0xed, 0xf7, 0x00,
	0xe6, 0xd2, 0x54, 0xb4, 0x6f, 0x2d, 0x71, 0xf8, 0xf0, 0x1f, 0x0f, 0xb6, 0xcf, 0xf4, 0x64, 0xb9,
	0x7f, 0xe0, 0x6e, 0x61, 0x48, 0xb9, 0x73, 0xc1, 0x2d, 0x0c, 0x62, 0x41, 0x2d, 0x4a, 0x9f, 0x96,
	0xa8, 0x65, 0x4f, 0xa7, 0xdc, 0xa7, 0xd4, 0x57, 0x2f, 0xa1, 0xfe, 0x7a, 0xba, 0x7a, 0x68, 0xe7,
	0xf2, 0x90, 0x8e, 0xa9, 0xa4, 0x57, 0x7e, 0xef, 0x70, 0xdb, 0xea, 0xd2, 0xe1, 0xcd, 0x4d, 0xc3,
	0x03, 0xd8, 0x34, 0x16, 0xed, 0x14, 0x2e, 0xef, 0x5d, 0x58, 0xd3, 0x24, 0xc5, 0x2c, 0x15, 0x6d,
	0xaf, 0x5b, 0xed, 0x35, 0xf7, 0x6f, 0xfe, 0xf5, 0xb6, 0xb3, 0xf2, 0xf7, 0xdb, 0xce, 0x0d, 0x1d,
	0x79, 0x74, 0x88, 0x57, 0x75, 0xc4, 0x51, 0x2a, 0xc2, 0x47, 0xb0, 0x55, 0x4e, 0x62, 0x89, 0xbf,
	0x03, 0x37, 0xa7, 0xf9, 0x88, 0xe4, 0xe9, 0x98, 0xa6, 0x71, 0xc2, 0xa7, 0xb9, 0xbb, 0x68, 0x6b,
	0x66, 0x3e, 0x50, 0xd6, 0xb0, 0x00, 0x1f, 0x53, 0x49, 0x58, 0xee, 0xea, 0x1f, 0x81, 0x9f, 0x14,
	0x94, 0x48, 0xc6, 0xf3, 0x38, 0x25, 0xd2, 0xcd, 0x4d, 0x10, 0x99, 0x6d, 0x16, 0xb9, 0x6d, 0x16,
	0x1d, 0xbb, 0x6d, 0xb6, 0xbf, 0xaa, 0xfa, 0xfb, 0xed, 0xdf, 0x8e, 0x87, 0x9b, 0x0e, 0x7a, 0x48,
	0x24, 0x55, 0x24, 0x0f, 0xd8, 0x58, 0x5a, 0xa5, 0x37, 0xb1, 0x3d, 0x85, 0x1b, 0xd0, 0x72, 0x35,
	0x2d, 0x17, 0xdb, 0xb0, 0x89, 0x8d, 0x2c, 0x8e, 0x0b, 0xf5, 0xae, 0xa6, 0x97, 0xf0, 0x36, 0x6c,
	0x95, 0xcd, 0x36, 0xfc, 0x97, 0x2a, 0xac, 0x1b, 0x11, 0x50, 0xa2, 0x86, 0xe5, 0x09, 0xb4, 0x06,
	0xbc, 0xc8, 0x88, 0x8c, 0x7f, 0xa2, 0x85, 0x60, 0x3c, 0xd7, 0x4d, 0xb7, 0xfa, 0x3b, 0x0b, 0x7a,
	0x33, 0x80, 0xe8, 0xb1, 0x8e, 0x7e, 0x65, 0x82, 0xb1, 0x3f, 0x98, 0x3f, 0x2a, 0x0d, 0xcc, 0x54,
	0xd7, 0xb4, 0x12, 0x5b, 0xdc, 0xad, 0x8d, 0xab, 0xee, 0xd6, 0x79, 0x62, 0xd5, 0x97, 0xc0, 0xaa,
	0xf5, 0x8a, 0xc4, 0x2a, 0x27, 0x7a, 0x17, 0xd6, 0x04, 0x1b, 0xe6, 0x44, 0x4e, 0x0b, 0x33, 0xd5,
	0x4d, 0x7c, 0x6a, 0x40, 0x9f, 0xc3, 0xba, 0x6e, 0x2a, 0x36, 0xfa, 0xac, 0x9f, 0xa7, 0xcf, 0xfd,
	0x9a, 0x4a, 0x8f, 0x81, 0xcf, 0x2c, 0xe1, 0xa7, 0xe0, 0x97, 0xa8, 0x41, 0x3e, 0xac, 0x3d, 0x7e,
	0x86, 0x9f, 0xee, 0x1d, 0xc7, 0xaf, 0xee, 0x6d, 0xac, 0xcc, 0x1f, 0x3f, 0xdb, 0xf0, 0xc2, 0x07,
	0xe0, 0x7f, 0xf5, 0x9a, 0x09, 0x79, 0x4d, 0xf1, 0x7e, 0x02, 0x2d, 0x07, 0xb7, 0xb2, 0x6d, 0xc3,
	0x8d, 0x8c, 0x09, 0xc1, 0xf2, 0xa1, 0x46, 0xfb, 0xd8, 0x1d, 0xfb, 0x7f, 0xd4, 0x00, 0x9e, 0xcf,
	0x1e, 0x13, 0x3d, 0x85, 0x86, 0x59, 0xf3, 0xe8, 0xfd, 0x8b, 0x3f, 0x3f, 0x41, 0xe7, 0x5c, 0xbf,
	0x15, 0xd3, 0x4a, 0xcf, 0x43, 0x2f, 0x61, 0xd5, 0x6d, 0x22, 0xd4, 0xbd, 0x6c, 0xd5, 0x06, 0x1f,
	0x5c, 0xba, 0xc6, 0x54, 0xd2, 0x7b, 0x1e, 0xfa, 0x06, 0x1a, 0x66, 0x3a, 0x97, 0x74, 0x59, 0xda,
	0x26, 0x4b, 0xba, 0x3c, 0xb3, 0x2d, 0xaa, 0xbf, 0x56, 0x3c, 0xf4, 0x2d, 0x34, 0xe7, 0xa7, 0x1d,
	0x95, 0x50, 0x4b, 0x96, 0x49, 0xd0, 0x3d, 0x3f, 0xc0, 0x32, 0xfe, 0x08, 0x1a, 0x66, 0x16, 0xd1,
	0x3b, 0xf3, 0xb1, 0xa5, 0x9d, 0x10, 0x04, 0xcb, 0x5c, 0x36, 0xc1, 0x0b, 0x68, 0xce, 0xcf, 0x68,
	0xb9, 0xa7, 0x25, 0x43, 0x5d, 0xee, 0x69, 0xe9, 0x78, 0xaf, 0xa8, 0xae, 0x8c, 0x32, 0xca, 0x5d,
	0x95, 0xc4, 0x56, 0xee, 0xaa, 0x2c, 0xa4, 0xfd, 0xad, 0xef, 0x90, 0xb2, 0x7f, 0x1f, 0x31, 0xbe,
	0x9b, 0xf0, 0x2c, 0xe3, 0xf9, 0xee, 0xe4, 0xe4, 0xa4, 0xa1, 0x67, 0xec, 0xfe, 0x7f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x4a, 0x05, 0x24, 0xe5, 0xc5, 0x09, 0x00, 0x00,
}
