// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Packet.proto

package Packet

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Header_Type int32

const (
	Header_NONE             Header_Type = 0
	Header_SC_Error         Header_Type = 1
	Header_CS_Enter         Header_Type = 10000
	Header_CS_Enter_Ack     Header_Type = 10001
	Header_CS_Broadcast     Header_Type = 10002
	Header_CS_Broadcast_Ack Header_Type = 10003
	Header_SC_Broadcast     Header_Type = 20000
)

var Header_Type_name = map[int32]string{
	0:     "NONE",
	1:     "SC_Error",
	10000: "CS_Enter",
	10001: "CS_Enter_Ack",
	10002: "CS_Broadcast",
	10003: "CS_Broadcast_Ack",
	20000: "SC_Broadcast",
}

var Header_Type_value = map[string]int32{
	"NONE":             0,
	"SC_Error":         1,
	"CS_Enter":         10000,
	"CS_Enter_Ack":     10001,
	"CS_Broadcast":     10002,
	"CS_Broadcast_Ack": 10003,
	"SC_Broadcast":     20000,
}

func (x Header_Type) String() string {
	return proto.EnumName(Header_Type_name, int32(x))
}

func (Header_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{6, 0}
}

//서버 -> 클라이언트
type SC_Error struct {
	ErrorCode            int32    `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,2,opt,name=ErrorMessage,proto3" json:"ErrorMessage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SC_Error) Reset()         { *m = SC_Error{} }
func (m *SC_Error) String() string { return proto.CompactTextString(m) }
func (*SC_Error) ProtoMessage()    {}
func (*SC_Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{0}
}

func (m *SC_Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SC_Error.Unmarshal(m, b)
}
func (m *SC_Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SC_Error.Marshal(b, m, deterministic)
}
func (m *SC_Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SC_Error.Merge(m, src)
}
func (m *SC_Error) XXX_Size() int {
	return xxx_messageInfo_SC_Error.Size(m)
}
func (m *SC_Error) XXX_DiscardUnknown() {
	xxx_messageInfo_SC_Error.DiscardUnknown(m)
}

var xxx_messageInfo_SC_Error proto.InternalMessageInfo

func (m *SC_Error) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

func (m *SC_Error) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

//클라이언트 -> 서버
type CS_Enter struct {
	GameUID              uint64   `protobuf:"varint,1,opt,name=GameUID,proto3" json:"GameUID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CS_Enter) Reset()         { *m = CS_Enter{} }
func (m *CS_Enter) String() string { return proto.CompactTextString(m) }
func (*CS_Enter) ProtoMessage()    {}
func (*CS_Enter) Descriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{1}
}

func (m *CS_Enter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CS_Enter.Unmarshal(m, b)
}
func (m *CS_Enter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CS_Enter.Marshal(b, m, deterministic)
}
func (m *CS_Enter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CS_Enter.Merge(m, src)
}
func (m *CS_Enter) XXX_Size() int {
	return xxx_messageInfo_CS_Enter.Size(m)
}
func (m *CS_Enter) XXX_DiscardUnknown() {
	xxx_messageInfo_CS_Enter.DiscardUnknown(m)
}

var xxx_messageInfo_CS_Enter proto.InternalMessageInfo

func (m *CS_Enter) GetGameUID() uint64 {
	if m != nil {
		return m.GameUID
	}
	return 0
}

type CS_Enter_Ack struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CS_Enter_Ack) Reset()         { *m = CS_Enter_Ack{} }
func (m *CS_Enter_Ack) String() string { return proto.CompactTextString(m) }
func (*CS_Enter_Ack) ProtoMessage()    {}
func (*CS_Enter_Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{2}
}

func (m *CS_Enter_Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CS_Enter_Ack.Unmarshal(m, b)
}
func (m *CS_Enter_Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CS_Enter_Ack.Marshal(b, m, deterministic)
}
func (m *CS_Enter_Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CS_Enter_Ack.Merge(m, src)
}
func (m *CS_Enter_Ack) XXX_Size() int {
	return xxx_messageInfo_CS_Enter_Ack.Size(m)
}
func (m *CS_Enter_Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_CS_Enter_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_CS_Enter_Ack proto.InternalMessageInfo

type CS_Broadcast struct {
	Message              string   `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CS_Broadcast) Reset()         { *m = CS_Broadcast{} }
func (m *CS_Broadcast) String() string { return proto.CompactTextString(m) }
func (*CS_Broadcast) ProtoMessage()    {}
func (*CS_Broadcast) Descriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{3}
}

func (m *CS_Broadcast) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CS_Broadcast.Unmarshal(m, b)
}
func (m *CS_Broadcast) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CS_Broadcast.Marshal(b, m, deterministic)
}
func (m *CS_Broadcast) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CS_Broadcast.Merge(m, src)
}
func (m *CS_Broadcast) XXX_Size() int {
	return xxx_messageInfo_CS_Broadcast.Size(m)
}
func (m *CS_Broadcast) XXX_DiscardUnknown() {
	xxx_messageInfo_CS_Broadcast.DiscardUnknown(m)
}

var xxx_messageInfo_CS_Broadcast proto.InternalMessageInfo

func (m *CS_Broadcast) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type CS_Broadcast_Ack struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CS_Broadcast_Ack) Reset()         { *m = CS_Broadcast_Ack{} }
func (m *CS_Broadcast_Ack) String() string { return proto.CompactTextString(m) }
func (*CS_Broadcast_Ack) ProtoMessage()    {}
func (*CS_Broadcast_Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{4}
}

func (m *CS_Broadcast_Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CS_Broadcast_Ack.Unmarshal(m, b)
}
func (m *CS_Broadcast_Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CS_Broadcast_Ack.Marshal(b, m, deterministic)
}
func (m *CS_Broadcast_Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CS_Broadcast_Ack.Merge(m, src)
}
func (m *CS_Broadcast_Ack) XXX_Size() int {
	return xxx_messageInfo_CS_Broadcast_Ack.Size(m)
}
func (m *CS_Broadcast_Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_CS_Broadcast_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_CS_Broadcast_Ack proto.InternalMessageInfo

type SC_Broadcast struct {
	Message              string   `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SC_Broadcast) Reset()         { *m = SC_Broadcast{} }
func (m *SC_Broadcast) String() string { return proto.CompactTextString(m) }
func (*SC_Broadcast) ProtoMessage()    {}
func (*SC_Broadcast) Descriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{5}
}

func (m *SC_Broadcast) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SC_Broadcast.Unmarshal(m, b)
}
func (m *SC_Broadcast) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SC_Broadcast.Marshal(b, m, deterministic)
}
func (m *SC_Broadcast) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SC_Broadcast.Merge(m, src)
}
func (m *SC_Broadcast) XXX_Size() int {
	return xxx_messageInfo_SC_Broadcast.Size(m)
}
func (m *SC_Broadcast) XXX_DiscardUnknown() {
	xxx_messageInfo_SC_Broadcast.DiscardUnknown(m)
}

var xxx_messageInfo_SC_Broadcast proto.InternalMessageInfo

func (m *SC_Broadcast) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Header struct {
	Type                 Header_Type `protobuf:"varint,1,opt,name=type,proto3,enum=Header_Type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_365103b9ca9fe667, []int{6}
}

func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (m *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(m, src)
}
func (m *Header) XXX_Size() int {
	return xxx_messageInfo_Header.Size(m)
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetType() Header_Type {
	if m != nil {
		return m.Type
	}
	return Header_NONE
}

func init() {
	proto.RegisterEnum("Header_Type", Header_Type_name, Header_Type_value)
	proto.RegisterType((*SC_Error)(nil), "SC_Error")
	proto.RegisterType((*CS_Enter)(nil), "CS_Enter")
	proto.RegisterType((*CS_Enter_Ack)(nil), "CS_Enter_Ack")
	proto.RegisterType((*CS_Broadcast)(nil), "CS_Broadcast")
	proto.RegisterType((*CS_Broadcast_Ack)(nil), "CS_Broadcast_Ack")
	proto.RegisterType((*SC_Broadcast)(nil), "SC_Broadcast")
	proto.RegisterType((*Header)(nil), "Header")
}

func init() { proto.RegisterFile("Packet.proto", fileDescriptor_365103b9ca9fe667) }

var fileDescriptor_365103b9ca9fe667 = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x09, 0x48, 0x4c, 0xce,
	0x4e, 0x2d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xf2, 0xe1, 0xe2, 0x08, 0x76, 0x8e, 0x77,
	0x2d, 0x2a, 0xca, 0x2f, 0x12, 0x92, 0xe1, 0xe2, 0x04, 0x33, 0x9c, 0xf3, 0x53, 0x52, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x58, 0x83, 0x10, 0x02, 0x42, 0x4a, 0x5c, 0x3c, 0x60, 0x8e, 0x6f, 0x6a, 0x71,
	0x71, 0x62, 0x7a, 0xaa, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8a, 0x98, 0x92, 0x0a, 0x17,
	0x87, 0x73, 0x70, 0xbc, 0x6b, 0x5e, 0x49, 0x6a, 0x91, 0x90, 0x04, 0x17, 0xbb, 0x7b, 0x62, 0x6e,
	0x6a, 0xa8, 0xa7, 0x0b, 0xd8, 0x2c, 0x96, 0x20, 0x18, 0x57, 0x89, 0x8f, 0x8b, 0x07, 0xa6, 0x2a,
	0xde, 0x31, 0x39, 0x5b, 0x49, 0x03, 0xcc, 0x77, 0x2a, 0xca, 0x4f, 0x4c, 0x49, 0x4e, 0x2c, 0x2e,
	0x01, 0xe9, 0x84, 0x59, 0xc2, 0x08, 0xb6, 0x04, 0xc6, 0x55, 0x12, 0xe2, 0x12, 0x40, 0x56, 0x09,
	0xd3, 0x1d, 0xec, 0x4c, 0x94, 0xee, 0x55, 0x8c, 0x5c, 0x6c, 0x1e, 0xa9, 0x89, 0x29, 0xa9, 0x45,
	0x42, 0x0a, 0x5c, 0x2c, 0x25, 0x95, 0x05, 0x10, 0x15, 0x7c, 0x46, 0x3c, 0x7a, 0x10, 0x61, 0xbd,
	0x90, 0xca, 0x82, 0xd4, 0x20, 0xb0, 0x8c, 0x52, 0x1d, 0x17, 0x0b, 0x88, 0x27, 0xc4, 0xc1, 0xc5,
	0xe2, 0xe7, 0xef, 0xe7, 0x2a, 0xc0, 0x20, 0xc4, 0x83, 0x08, 0x2a, 0x01, 0x46, 0x21, 0x5e, 0x84,
	0x57, 0x05, 0x26, 0xf8, 0x09, 0x09, 0xa2, 0xfa, 0x49, 0x60, 0x22, 0x4c, 0x08, 0xee, 0x30, 0x81,
	0x49, 0x7e, 0x42, 0xa2, 0x98, 0xee, 0x17, 0x98, 0xec, 0x27, 0x24, 0x84, 0xea, 0x05, 0x81, 0x05,
	0x73, 0x18, 0x93, 0xd8, 0xc0, 0xf1, 0x63, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xfe, 0xd7, 0xee,
	0xb8, 0xaf, 0x01, 0x00, 0x00,
}
