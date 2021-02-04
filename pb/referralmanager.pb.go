// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: referralmanager.proto

package pb

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"

	drpc "storj.io/drpc"
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

type GetTokensRequest struct {
	OwnerUserId          []byte   `protobuf:"bytes,1,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	OwnerSatelliteId     NodeID   `protobuf:"bytes,2,opt,name=owner_satellite_id,json=ownerSatelliteId,proto3,customtype=NodeID" json:"owner_satellite_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTokensRequest) Reset()         { *m = GetTokensRequest{} }
func (m *GetTokensRequest) String() string { return proto.CompactTextString(m) }
func (*GetTokensRequest) ProtoMessage()    {}
func (*GetTokensRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d96ad24f1e021c, []int{0}
}
func (m *GetTokensRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTokensRequest.Unmarshal(m, b)
}
func (m *GetTokensRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTokensRequest.Marshal(b, m, deterministic)
}
func (m *GetTokensRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTokensRequest.Merge(m, src)
}
func (m *GetTokensRequest) XXX_Size() int {
	return xxx_messageInfo_GetTokensRequest.Size(m)
}
func (m *GetTokensRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTokensRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTokensRequest proto.InternalMessageInfo

func (m *GetTokensRequest) GetOwnerUserId() []byte {
	if m != nil {
		return m.OwnerUserId
	}
	return nil
}

type GetTokensResponse struct {
	TokenSecrets         [][]byte `protobuf:"bytes,1,rep,name=token_secrets,json=tokenSecrets,proto3" json:"token_secrets,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTokensResponse) Reset()         { *m = GetTokensResponse{} }
func (m *GetTokensResponse) String() string { return proto.CompactTextString(m) }
func (*GetTokensResponse) ProtoMessage()    {}
func (*GetTokensResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d96ad24f1e021c, []int{1}
}
func (m *GetTokensResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTokensResponse.Unmarshal(m, b)
}
func (m *GetTokensResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTokensResponse.Marshal(b, m, deterministic)
}
func (m *GetTokensResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTokensResponse.Merge(m, src)
}
func (m *GetTokensResponse) XXX_Size() int {
	return xxx_messageInfo_GetTokensResponse.Size(m)
}
func (m *GetTokensResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTokensResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTokensResponse proto.InternalMessageInfo

func (m *GetTokensResponse) GetTokenSecrets() [][]byte {
	if m != nil {
		return m.TokenSecrets
	}
	return nil
}

type RedeemTokenRequest struct {
	Token                []byte   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	RedeemUserId         []byte   `protobuf:"bytes,2,opt,name=redeem_user_id,json=redeemUserId,proto3" json:"redeem_user_id,omitempty"`
	RedeemSatelliteId    NodeID   `protobuf:"bytes,3,opt,name=redeem_satellite_id,json=redeemSatelliteId,proto3,customtype=NodeID" json:"redeem_satellite_id"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RedeemTokenRequest) Reset()         { *m = RedeemTokenRequest{} }
func (m *RedeemTokenRequest) String() string { return proto.CompactTextString(m) }
func (*RedeemTokenRequest) ProtoMessage()    {}
func (*RedeemTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d96ad24f1e021c, []int{2}
}
func (m *RedeemTokenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RedeemTokenRequest.Unmarshal(m, b)
}
func (m *RedeemTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RedeemTokenRequest.Marshal(b, m, deterministic)
}
func (m *RedeemTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RedeemTokenRequest.Merge(m, src)
}
func (m *RedeemTokenRequest) XXX_Size() int {
	return xxx_messageInfo_RedeemTokenRequest.Size(m)
}
func (m *RedeemTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RedeemTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RedeemTokenRequest proto.InternalMessageInfo

func (m *RedeemTokenRequest) GetToken() []byte {
	if m != nil {
		return m.Token
	}
	return nil
}

func (m *RedeemTokenRequest) GetRedeemUserId() []byte {
	if m != nil {
		return m.RedeemUserId
	}
	return nil
}

type RedeemTokenResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RedeemTokenResponse) Reset()         { *m = RedeemTokenResponse{} }
func (m *RedeemTokenResponse) String() string { return proto.CompactTextString(m) }
func (*RedeemTokenResponse) ProtoMessage()    {}
func (*RedeemTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_45d96ad24f1e021c, []int{3}
}
func (m *RedeemTokenResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RedeemTokenResponse.Unmarshal(m, b)
}
func (m *RedeemTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RedeemTokenResponse.Marshal(b, m, deterministic)
}
func (m *RedeemTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RedeemTokenResponse.Merge(m, src)
}
func (m *RedeemTokenResponse) XXX_Size() int {
	return xxx_messageInfo_RedeemTokenResponse.Size(m)
}
func (m *RedeemTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RedeemTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RedeemTokenResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GetTokensRequest)(nil), "referralmanager.GetTokensRequest")
	proto.RegisterType((*GetTokensResponse)(nil), "referralmanager.GetTokensResponse")
	proto.RegisterType((*RedeemTokenRequest)(nil), "referralmanager.RedeemTokenRequest")
	proto.RegisterType((*RedeemTokenResponse)(nil), "referralmanager.RedeemTokenResponse")
}

func init() { proto.RegisterFile("referralmanager.proto", fileDescriptor_45d96ad24f1e021c) }

var fileDescriptor_45d96ad24f1e021c = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcf, 0x4e, 0xf2, 0x40,
	0x14, 0xc5, 0xbf, 0x7e, 0x44, 0x12, 0x2f, 0xe5, 0xdf, 0x00, 0x09, 0x61, 0x03, 0x16, 0x16, 0xac,
	0x20, 0xd1, 0x8d, 0x0b, 0xe3, 0x82, 0x98, 0x18, 0x16, 0xba, 0x18, 0x34, 0x31, 0x6e, 0x48, 0xa1,
	0x57, 0x52, 0x6d, 0x7b, 0xeb, 0xcc, 0x10, 0x5f, 0xc3, 0x37, 0x72, 0xeb, 0x33, 0xb8, 0xe0, 0x59,
	0x4c, 0x67, 0x0a, 0x29, 0xc5, 0xb0, 0xec, 0xe9, 0xef, 0xf6, 0x9e, 0x73, 0x6e, 0xa1, 0x25, 0xf0,
	0x05, 0x85, 0x70, 0x83, 0xd0, 0x8d, 0xdc, 0x15, 0x8a, 0x51, 0x2c, 0x48, 0x11, 0xab, 0xe6, 0xe4,
	0x0e, 0xac, 0x68, 0x45, 0xe6, 0xa5, 0xa3, 0xa0, 0x76, 0x8b, 0xea, 0x81, 0xde, 0x30, 0x92, 0x1c,
	0xdf, 0xd7, 0x28, 0x15, 0x73, 0xa0, 0x4c, 0x1f, 0x11, 0x8a, 0xf9, 0x5a, 0xa2, 0x98, 0xfb, 0x5e,
	0xdb, 0xea, 0x59, 0x43, 0x9b, 0x97, 0xb4, 0xf8, 0x28, 0x51, 0x4c, 0x3d, 0x76, 0x05, 0xcc, 0x30,
	0xd2, 0x55, 0x18, 0x04, 0xbe, 0xc2, 0x04, 0xfc, 0x9f, 0x80, 0x93, 0xca, 0xf7, 0xa6, 0xfb, 0xef,
	0x67, 0xd3, 0x2d, 0xde, 0x93, 0x87, 0xd3, 0x1b, 0x5e, 0xd3, 0xe4, 0x6c, 0x0b, 0x4e, 0x3d, 0xe7,
	0x12, 0xea, 0x99, 0xad, 0x32, 0xa6, 0x48, 0x22, 0xeb, 0x43, 0x59, 0x25, 0xca, 0x5c, 0xe2, 0x52,
	0xa0, 0x92, 0x6d, 0xab, 0x57, 0x18, 0xda, 0xdc, 0xd6, 0xe2, 0xcc, 0x68, 0xce, 0xa7, 0x05, 0x8c,
	0xa3, 0x87, 0x18, 0xea, 0xe9, 0xad, 0xe5, 0x26, 0x9c, 0x68, 0x2c, 0xb5, 0x6a, 0x1e, 0xd8, 0x00,
	0x2a, 0x42, 0xb3, 0xbb, 0x24, 0xda, 0x20, 0xb7, 0x8d, 0x9a, 0x46, 0xb9, 0x86, 0x46, 0x4a, 0xed,
	0x65, 0x29, 0xfc, 0x99, 0xa5, 0x6e, 0xd0, 0x6c, 0x98, 0x16, 0x34, 0xf6, 0x1c, 0x99, 0x38, 0xe7,
	0x5f, 0x16, 0x54, 0x79, 0xda, 0xfc, 0x9d, 0x69, 0x9e, 0x71, 0x38, 0xdd, 0xe5, 0x66, 0x67, 0xa3,
	0xfc, 0xbd, 0xf2, 0x97, 0xe8, 0x38, 0xc7, 0x90, 0xb4, 0xb6, 0x27, 0x28, 0x65, 0xd6, 0xb3, 0xfe,
	0xc1, 0xc8, 0x61, 0x5d, 0x9d, 0xc1, 0x71, 0xc8, 0x7c, 0x79, 0xd2, 0x7c, 0x66, 0x52, 0x91, 0x78,
	0x1d, 0xf9, 0x34, 0x5e, 0x52, 0x18, 0x52, 0x34, 0x8e, 0x17, 0x8b, 0xa2, 0xfe, 0x71, 0x2e, 0x7e,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x85, 0x79, 0x33, 0xba, 0x6e, 0x02, 0x00, 0x00,
}

// --- DRPC BEGIN ---

type DRPCReferralManagerClient interface {
	DRPCConn() drpc.Conn

	// GetTokens retrieves a list of unredeemed tokens for a user
	GetTokens(ctx context.Context, in *GetTokensRequest) (*GetTokensResponse, error)
	// RedeemToken saves newly created user info in referral manager
	RedeemToken(ctx context.Context, in *RedeemTokenRequest) (*RedeemTokenResponse, error)
}

type drpcReferralManagerClient struct {
	cc drpc.Conn
}

func NewDRPCReferralManagerClient(cc drpc.Conn) DRPCReferralManagerClient {
	return &drpcReferralManagerClient{cc}
}

func (c *drpcReferralManagerClient) DRPCConn() drpc.Conn { return c.cc }

func (c *drpcReferralManagerClient) GetTokens(ctx context.Context, in *GetTokensRequest) (*GetTokensResponse, error) {
	out := new(GetTokensResponse)
	err := c.cc.Invoke(ctx, "/referralmanager.ReferralManager/GetTokens", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *drpcReferralManagerClient) RedeemToken(ctx context.Context, in *RedeemTokenRequest) (*RedeemTokenResponse, error) {
	out := new(RedeemTokenResponse)
	err := c.cc.Invoke(ctx, "/referralmanager.ReferralManager/RedeemToken", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type DRPCReferralManagerServer interface {
	// GetTokens retrieves a list of unredeemed tokens for a user
	GetTokens(context.Context, *GetTokensRequest) (*GetTokensResponse, error)
	// RedeemToken saves newly created user info in referral manager
	RedeemToken(context.Context, *RedeemTokenRequest) (*RedeemTokenResponse, error)
}

type DRPCReferralManagerDescription struct{}

func (DRPCReferralManagerDescription) NumMethods() int { return 2 }

func (DRPCReferralManagerDescription) Method(n int) (string, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		return "/referralmanager.ReferralManager/GetTokens",
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCReferralManagerServer).
					GetTokens(
						ctx,
						in1.(*GetTokensRequest),
					)
			}, DRPCReferralManagerServer.GetTokens, true
	case 1:
		return "/referralmanager.ReferralManager/RedeemToken",
			func(srv interface{}, ctx context.Context, in1, in2 interface{}) (drpc.Message, error) {
				return srv.(DRPCReferralManagerServer).
					RedeemToken(
						ctx,
						in1.(*RedeemTokenRequest),
					)
			}, DRPCReferralManagerServer.RedeemToken, true
	default:
		return "", nil, nil, false
	}
}

func DRPCRegisterReferralManager(mux drpc.Mux, impl DRPCReferralManagerServer) error {
	return mux.Register(impl, DRPCReferralManagerDescription{})
}

type DRPCReferralManager_GetTokensStream interface {
	drpc.Stream
	SendAndClose(*GetTokensResponse) error
}

type drpcReferralManagerGetTokensStream struct {
	drpc.Stream
}

func (x *drpcReferralManagerGetTokensStream) SendAndClose(m *GetTokensResponse) error {
	if err := x.MsgSend(m); err != nil {
		return err
	}
	return x.CloseSend()
}

type DRPCReferralManager_RedeemTokenStream interface {
	drpc.Stream
	SendAndClose(*RedeemTokenResponse) error
}

type drpcReferralManagerRedeemTokenStream struct {
	drpc.Stream
}

func (x *drpcReferralManagerRedeemTokenStream) SendAndClose(m *RedeemTokenResponse) error {
	if err := x.MsgSend(m); err != nil {
		return err
	}
	return x.CloseSend()
}

// --- DRPC END ---
