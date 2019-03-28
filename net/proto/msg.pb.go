// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.proto

package proto

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

//info
type ServerInfo struct {
	ServerId             int64    `protobuf:"varint,1,opt,name=serverId,proto3" json:"serverId,omitempty"`
	ServerType           int32    `protobuf:"varint,2,opt,name=serverType,proto3" json:"serverType,omitempty"`
	ServerVersion        int32    `protobuf:"varint,3,opt,name=serverVersion,proto3" json:"serverVersion,omitempty"`
	State                int64    `protobuf:"varint,4,opt,name=state,proto3" json:"state,omitempty"`
	ServerName           string   `protobuf:"bytes,5,opt,name=serverName,proto3" json:"serverName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerInfo) Reset()         { *m = ServerInfo{} }
func (m *ServerInfo) String() string { return proto.CompactTextString(m) }
func (*ServerInfo) ProtoMessage()    {}
func (*ServerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

func (m *ServerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerInfo.Unmarshal(m, b)
}
func (m *ServerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerInfo.Marshal(b, m, deterministic)
}
func (m *ServerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerInfo.Merge(m, src)
}
func (m *ServerInfo) XXX_Size() int {
	return xxx_messageInfo_ServerInfo.Size(m)
}
func (m *ServerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ServerInfo proto.InternalMessageInfo

func (m *ServerInfo) GetServerId() int64 {
	if m != nil {
		return m.ServerId
	}
	return 0
}

func (m *ServerInfo) GetServerType() int32 {
	if m != nil {
		return m.ServerType
	}
	return 0
}

func (m *ServerInfo) GetServerVersion() int32 {
	if m != nil {
		return m.ServerVersion
	}
	return 0
}

func (m *ServerInfo) GetState() int64 {
	if m != nil {
		return m.State
	}
	return 0
}

func (m *ServerInfo) GetServerName() string {
	if m != nil {
		return m.ServerName
	}
	return ""
}

type UserInfo struct {
	UId                  int64    `protobuf:"varint,1,opt,name=UId,proto3" json:"UId,omitempty"`
	UsersId              []int64  `protobuf:"varint,2,rep,packed,name=UsersId,proto3" json:"UsersId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserInfo) Reset()         { *m = UserInfo{} }
func (m *UserInfo) String() string { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()    {}
func (*UserInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
}

func (m *UserInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserInfo.Unmarshal(m, b)
}
func (m *UserInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserInfo.Marshal(b, m, deterministic)
}
func (m *UserInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserInfo.Merge(m, src)
}
func (m *UserInfo) XXX_Size() int {
	return xxx_messageInfo_UserInfo.Size(m)
}
func (m *UserInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_UserInfo.DiscardUnknown(m)
}

var xxx_messageInfo_UserInfo proto.InternalMessageInfo

func (m *UserInfo) GetUId() int64 {
	if m != nil {
		return m.UId
	}
	return 0
}

func (m *UserInfo) GetUsersId() []int64 {
	if m != nil {
		return m.UsersId
	}
	return nil
}

<<<<<<< HEAD
type NewConnect struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
=======
//msg
type ConnectRS struct {
	Result               int32    `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

<<<<<<< HEAD
func (m *NewConnect) Reset()         { *m = NewConnect{} }
func (m *NewConnect) String() string { return proto.CompactTextString(m) }
func (*NewConnect) ProtoMessage()    {}
func (*NewConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{2}
}

func (m *NewConnect) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewConnect.Unmarshal(m, b)
}
func (m *NewConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewConnect.Marshal(b, m, deterministic)
}
func (m *NewConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewConnect.Merge(m, src)
}
func (m *NewConnect) XXX_Size() int {
	return xxx_messageInfo_NewConnect.Size(m)
}
func (m *NewConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_NewConnect.DiscardUnknown(m)
}

var xxx_messageInfo_NewConnect proto.InternalMessageInfo

func (m *NewConnect) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type DelConnect struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DelConnect) Reset()         { *m = DelConnect{} }
func (m *DelConnect) String() string { return proto.CompactTextString(m) }
func (*DelConnect) ProtoMessage()    {}
func (*DelConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{3}
}

func (m *DelConnect) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DelConnect.Unmarshal(m, b)
}
func (m *DelConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DelConnect.Marshal(b, m, deterministic)
}
func (m *DelConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DelConnect.Merge(m, src)
}
func (m *DelConnect) XXX_Size() int {
	return xxx_messageInfo_DelConnect.Size(m)
}
func (m *DelConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_DelConnect.DiscardUnknown(m)
}

var xxx_messageInfo_DelConnect proto.InternalMessageInfo

func (m *DelConnect) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type ServerDelConnect struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerDelConnect) Reset()         { *m = ServerDelConnect{} }
func (m *ServerDelConnect) String() string { return proto.CompactTextString(m) }
func (*ServerDelConnect) ProtoMessage()    {}
func (*ServerDelConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{4}
}

func (m *ServerDelConnect) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerDelConnect.Unmarshal(m, b)
}
func (m *ServerDelConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerDelConnect.Marshal(b, m, deterministic)
}
func (m *ServerDelConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerDelConnect.Merge(m, src)
}
func (m *ServerDelConnect) XXX_Size() int {
	return xxx_messageInfo_ServerDelConnect.Size(m)
}
func (m *ServerDelConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerDelConnect.DiscardUnknown(m)
}

var xxx_messageInfo_ServerDelConnect proto.InternalMessageInfo

func (m *ServerDelConnect) GetId() int64 {
	if m != nil {
		return m.Id
=======
func (m *ConnectRS) Reset()         { *m = ConnectRS{} }
func (m *ConnectRS) String() string { return proto.CompactTextString(m) }
func (*ConnectRS) ProtoMessage()    {}
func (*ConnectRS) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{2}
}

func (m *ConnectRS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectRS.Unmarshal(m, b)
}
func (m *ConnectRS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectRS.Marshal(b, m, deterministic)
}
func (m *ConnectRS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectRS.Merge(m, src)
}
func (m *ConnectRS) XXX_Size() int {
	return xxx_messageInfo_ConnectRS.Size(m)
}
func (m *ConnectRS) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectRS.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectRS proto.InternalMessageInfo

func (m *ConnectRS) GetResult() int32 {
	if m != nil {
		return m.Result
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
	}
	return 0
}

type ServerTick struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	State                int64    `protobuf:"varint,2,opt,name=state,proto3" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerTick) Reset()         { *m = ServerTick{} }
func (m *ServerTick) String() string { return proto.CompactTextString(m) }
func (*ServerTick) ProtoMessage()    {}
func (*ServerTick) Descriptor() ([]byte, []int) {
<<<<<<< HEAD
	return fileDescriptor_c06e4cca6c2cc899, []int{5}
=======
	return fileDescriptor_c06e4cca6c2cc899, []int{3}
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (m *ServerTick) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerTick.Unmarshal(m, b)
}
func (m *ServerTick) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerTick.Marshal(b, m, deterministic)
}
func (m *ServerTick) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerTick.Merge(m, src)
}
func (m *ServerTick) XXX_Size() int {
	return xxx_messageInfo_ServerTick.Size(m)
}
func (m *ServerTick) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerTick.DiscardUnknown(m)
}

var xxx_messageInfo_ServerTick proto.InternalMessageInfo

func (m *ServerTick) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *ServerTick) GetState() int64 {
	if m != nil {
		return m.State
	}
	return 0
}

type ServerLoginRQ struct {
	Serverinfo           *ServerInfo `protobuf:"bytes,1,opt,name=serverinfo,proto3" json:"serverinfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ServerLoginRQ) Reset()         { *m = ServerLoginRQ{} }
func (m *ServerLoginRQ) String() string { return proto.CompactTextString(m) }
func (*ServerLoginRQ) ProtoMessage()    {}
func (*ServerLoginRQ) Descriptor() ([]byte, []int) {
<<<<<<< HEAD
	return fileDescriptor_c06e4cca6c2cc899, []int{6}
=======
	return fileDescriptor_c06e4cca6c2cc899, []int{4}
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (m *ServerLoginRQ) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerLoginRQ.Unmarshal(m, b)
}
func (m *ServerLoginRQ) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerLoginRQ.Marshal(b, m, deterministic)
}
func (m *ServerLoginRQ) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerLoginRQ.Merge(m, src)
}
func (m *ServerLoginRQ) XXX_Size() int {
	return xxx_messageInfo_ServerLoginRQ.Size(m)
}
func (m *ServerLoginRQ) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerLoginRQ.DiscardUnknown(m)
}

var xxx_messageInfo_ServerLoginRQ proto.InternalMessageInfo

func (m *ServerLoginRQ) GetServerinfo() *ServerInfo {
	if m != nil {
		return m.Serverinfo
	}
	return nil
}

type ServerLoginRS struct {
	Result               int32       `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	Serverinfo           *ServerInfo `protobuf:"bytes,2,opt,name=serverinfo,proto3" json:"serverinfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ServerLoginRS) Reset()         { *m = ServerLoginRS{} }
func (m *ServerLoginRS) String() string { return proto.CompactTextString(m) }
func (*ServerLoginRS) ProtoMessage()    {}
func (*ServerLoginRS) Descriptor() ([]byte, []int) {
<<<<<<< HEAD
	return fileDescriptor_c06e4cca6c2cc899, []int{7}
=======
	return fileDescriptor_c06e4cca6c2cc899, []int{5}
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (m *ServerLoginRS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerLoginRS.Unmarshal(m, b)
}
func (m *ServerLoginRS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerLoginRS.Marshal(b, m, deterministic)
}
func (m *ServerLoginRS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerLoginRS.Merge(m, src)
}
func (m *ServerLoginRS) XXX_Size() int {
	return xxx_messageInfo_ServerLoginRS.Size(m)
}
func (m *ServerLoginRS) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerLoginRS.DiscardUnknown(m)
}

var xxx_messageInfo_ServerLoginRS proto.InternalMessageInfo

func (m *ServerLoginRS) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

func (m *ServerLoginRS) GetServerinfo() *ServerInfo {
	if m != nil {
		return m.Serverinfo
	}
	return nil
}

type ServerRegister struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Fuid                 []uint32 `protobuf:"varint,2,rep,packed,name=fuid,proto3" json:"fuid,omitempty"`
	Type                 int32    `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerRegister) Reset()         { *m = ServerRegister{} }
func (m *ServerRegister) String() string { return proto.CompactTextString(m) }
func (*ServerRegister) ProtoMessage()    {}
func (*ServerRegister) Descriptor() ([]byte, []int) {
<<<<<<< HEAD
	return fileDescriptor_c06e4cca6c2cc899, []int{8}
=======
	return fileDescriptor_c06e4cca6c2cc899, []int{6}
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (m *ServerRegister) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerRegister.Unmarshal(m, b)
}
func (m *ServerRegister) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerRegister.Marshal(b, m, deterministic)
}
func (m *ServerRegister) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerRegister.Merge(m, src)
}
func (m *ServerRegister) XXX_Size() int {
	return xxx_messageInfo_ServerRegister.Size(m)
}
func (m *ServerRegister) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerRegister.DiscardUnknown(m)
}

var xxx_messageInfo_ServerRegister proto.InternalMessageInfo

func (m *ServerRegister) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *ServerRegister) GetFuid() []uint32 {
	if m != nil {
		return m.Fuid
	}
	return nil
}

func (m *ServerRegister) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

type ServerDelFunc struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Fuid                 []uint32 `protobuf:"varint,2,rep,packed,name=fuid,proto3" json:"fuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServerDelFunc) Reset()         { *m = ServerDelFunc{} }
func (m *ServerDelFunc) String() string { return proto.CompactTextString(m) }
func (*ServerDelFunc) ProtoMessage()    {}
func (*ServerDelFunc) Descriptor() ([]byte, []int) {
<<<<<<< HEAD
	return fileDescriptor_c06e4cca6c2cc899, []int{9}
=======
	return fileDescriptor_c06e4cca6c2cc899, []int{7}
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (m *ServerDelFunc) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerDelFunc.Unmarshal(m, b)
}
func (m *ServerDelFunc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerDelFunc.Marshal(b, m, deterministic)
}
func (m *ServerDelFunc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerDelFunc.Merge(m, src)
}
func (m *ServerDelFunc) XXX_Size() int {
	return xxx_messageInfo_ServerDelFunc.Size(m)
}
func (m *ServerDelFunc) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerDelFunc.DiscardUnknown(m)
}

var xxx_messageInfo_ServerDelFunc proto.InternalMessageInfo

func (m *ServerDelFunc) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *ServerDelFunc) GetFuid() []uint32 {
	if m != nil {
		return m.Fuid
	}
	return nil
}

type ServerCall struct {
	User                 *UserInfo `protobuf:"bytes,1,opt,name=User,proto3" json:"User,omitempty"`
	FromMId              int64     `protobuf:"varint,2,opt,name=fromMId,proto3" json:"fromMId,omitempty"`
	ToMId                int64     `protobuf:"varint,3,opt,name=toMId,proto3" json:"toMId,omitempty"`
	FUId                 uint32    `protobuf:"varint,4,opt,name=fUId,proto3" json:"fUId,omitempty"`
	CUId                 string    `protobuf:"bytes,5,opt,name=cUId,proto3" json:"cUId,omitempty"`
	Msginfo              []byte    `protobuf:"bytes,6,opt,name=msginfo,proto3" json:"msginfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ServerCall) Reset()         { *m = ServerCall{} }
func (m *ServerCall) String() string { return proto.CompactTextString(m) }
func (*ServerCall) ProtoMessage()    {}
func (*ServerCall) Descriptor() ([]byte, []int) {
<<<<<<< HEAD
	return fileDescriptor_c06e4cca6c2cc899, []int{10}
=======
	return fileDescriptor_c06e4cca6c2cc899, []int{8}
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (m *ServerCall) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerCall.Unmarshal(m, b)
}
func (m *ServerCall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerCall.Marshal(b, m, deterministic)
}
func (m *ServerCall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerCall.Merge(m, src)
}
func (m *ServerCall) XXX_Size() int {
	return xxx_messageInfo_ServerCall.Size(m)
}
func (m *ServerCall) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerCall.DiscardUnknown(m)
}

var xxx_messageInfo_ServerCall proto.InternalMessageInfo

func (m *ServerCall) GetUser() *UserInfo {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *ServerCall) GetFromMId() int64 {
	if m != nil {
		return m.FromMId
	}
	return 0
}

func (m *ServerCall) GetToMId() int64 {
	if m != nil {
		return m.ToMId
	}
	return 0
}

func (m *ServerCall) GetFUId() uint32 {
	if m != nil {
		return m.FUId
	}
	return 0
}

func (m *ServerCall) GetCUId() string {
	if m != nil {
		return m.CUId
	}
	return ""
}

func (m *ServerCall) GetMsginfo() []byte {
	if m != nil {
		return m.Msginfo
	}
	return nil
}

type ServerCallBack struct {
	User                 *UserInfo `protobuf:"bytes,1,opt,name=User,proto3" json:"User,omitempty"`
	FromMId              int64     `protobuf:"varint,2,opt,name=fromMId,proto3" json:"fromMId,omitempty"`
	ToMId                int64     `protobuf:"varint,3,opt,name=toMId,proto3" json:"toMId,omitempty"`
	CUId                 string    `protobuf:"bytes,4,opt,name=cUId,proto3" json:"cUId,omitempty"`
	Msginfo              []byte    `protobuf:"bytes,5,opt,name=msginfo,proto3" json:"msginfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ServerCallBack) Reset()         { *m = ServerCallBack{} }
func (m *ServerCallBack) String() string { return proto.CompactTextString(m) }
func (*ServerCallBack) ProtoMessage()    {}
func (*ServerCallBack) Descriptor() ([]byte, []int) {
<<<<<<< HEAD
	return fileDescriptor_c06e4cca6c2cc899, []int{11}
=======
	return fileDescriptor_c06e4cca6c2cc899, []int{9}
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}

func (m *ServerCallBack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerCallBack.Unmarshal(m, b)
}
func (m *ServerCallBack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerCallBack.Marshal(b, m, deterministic)
}
func (m *ServerCallBack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerCallBack.Merge(m, src)
}
func (m *ServerCallBack) XXX_Size() int {
	return xxx_messageInfo_ServerCallBack.Size(m)
}
func (m *ServerCallBack) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerCallBack.DiscardUnknown(m)
}

var xxx_messageInfo_ServerCallBack proto.InternalMessageInfo

func (m *ServerCallBack) GetUser() *UserInfo {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *ServerCallBack) GetFromMId() int64 {
	if m != nil {
		return m.FromMId
	}
	return 0
}

func (m *ServerCallBack) GetToMId() int64 {
	if m != nil {
		return m.ToMId
	}
	return 0
}

func (m *ServerCallBack) GetCUId() string {
	if m != nil {
		return m.CUId
	}
	return ""
}

func (m *ServerCallBack) GetMsginfo() []byte {
	if m != nil {
		return m.Msginfo
	}
	return nil
}

func init() {
	proto.RegisterType((*ServerInfo)(nil), "proto.ServerInfo")
	proto.RegisterType((*UserInfo)(nil), "proto.UserInfo")
<<<<<<< HEAD
	proto.RegisterType((*NewConnect)(nil), "proto.NewConnect")
	proto.RegisterType((*DelConnect)(nil), "proto.DelConnect")
	proto.RegisterType((*ServerDelConnect)(nil), "proto.ServerDelConnect")
=======
	proto.RegisterType((*ConnectRS)(nil), "proto.ConnectRS")
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
	proto.RegisterType((*ServerTick)(nil), "proto.ServerTick")
	proto.RegisterType((*ServerLoginRQ)(nil), "proto.ServerLoginRQ")
	proto.RegisterType((*ServerLoginRS)(nil), "proto.ServerLoginRS")
	proto.RegisterType((*ServerRegister)(nil), "proto.ServerRegister")
	proto.RegisterType((*ServerDelFunc)(nil), "proto.ServerDelFunc")
	proto.RegisterType((*ServerCall)(nil), "proto.ServerCall")
	proto.RegisterType((*ServerCallBack)(nil), "proto.ServerCallBack")
}

func init() { proto.RegisterFile("msg.proto", fileDescriptor_c06e4cca6c2cc899) }

var fileDescriptor_c06e4cca6c2cc899 = []byte{
<<<<<<< HEAD
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0xc1, 0x8e, 0xd3, 0x30,
	0x10, 0x55, 0x92, 0xa6, 0xec, 0x0e, 0x74, 0x59, 0x2c, 0x84, 0x22, 0x84, 0x50, 0x64, 0x38, 0xe4,
	0xb4, 0x12, 0x20, 0xf6, 0x03, 0x76, 0x57, 0x48, 0x8b, 0xa0, 0x12, 0x6e, 0xcb, 0x81, 0x5b, 0x48,
	0x9d, 0xc8, 0x6a, 0x12, 0x57, 0xb1, 0x0b, 0xe2, 0x3b, 0xf8, 0x04, 0xf8, 0x50, 0x34, 0xe3, 0x38,
	0x4d, 0x11, 0xa0, 0x5e, 0x38, 0x65, 0xde, 0x4c, 0xfc, 0xfc, 0xfc, 0xe6, 0xc1, 0x69, 0x63, 0xaa,
	0x8b, 0x6d, 0xa7, 0xad, 0x66, 0x31, 0x7d, 0xf8, 0xcf, 0x00, 0x60, 0x21, 0xbb, 0x2f, 0xb2, 0xbb,
	0x6d, 0x4b, 0xcd, 0x1e, 0xc3, 0x89, 0x71, 0x68, 0x9d, 0x04, 0x69, 0x90, 0x45, 0x62, 0xc0, 0xec,
	0x29, 0x80, 0xab, 0x97, 0xdf, 0xb6, 0x32, 0x09, 0xd3, 0x20, 0x8b, 0xc5, 0xa8, 0xc3, 0x9e, 0xc3,
	0xcc, 0xa1, 0x8f, 0xb2, 0x33, 0x4a, 0xb7, 0x49, 0x44, 0xbf, 0x1c, 0x36, 0xd9, 0x43, 0x88, 0x8d,
	0xcd, 0xad, 0x4c, 0x26, 0x44, 0xef, 0xc0, 0x9e, 0x7b, 0x9e, 0x37, 0x32, 0x89, 0xd3, 0x20, 0x3b,
	0x15, 0xa3, 0x0e, 0xbf, 0x84, 0x93, 0x95, 0xe9, 0x35, 0x9e, 0x43, 0xb4, 0x1a, 0xe4, 0x61, 0xc9,
	0x12, 0xb8, 0x83, 0x53, 0x73, 0xbb, 0x4e, 0xc2, 0x34, 0xca, 0x22, 0xe1, 0x21, 0x7f, 0x02, 0x30,
	0x97, 0x5f, 0xaf, 0x75, 0xdb, 0xca, 0xc2, 0xb2, 0x33, 0x08, 0x95, 0x3f, 0x18, 0x2a, 0x9a, 0xde,
	0xc8, 0xfa, 0x6f, 0x53, 0x0e, 0xe7, 0xce, 0x99, 0x7f, 0xfc, 0x73, 0xe9, 0xdd, 0x5b, 0xaa, 0x62,
	0xc3, 0x18, 0x4c, 0xac, 0x6a, 0x64, 0x3f, 0xa7, 0x7a, 0xff, 0xde, 0x70, 0xf4, 0x5e, 0x7e, 0x05,
	0x33, 0x77, 0xee, 0x9d, 0xae, 0x54, 0x2b, 0x3e, 0xb0, 0x17, 0xde, 0x00, 0xd5, 0x96, 0x9a, 0x08,
	0xee, 0xbe, 0x7c, 0xe0, 0x56, 0x75, 0xb1, 0xdf, 0x8f, 0x18, 0xfd, 0xc4, 0x3f, 0x1d, 0x72, 0x2c,
	0xd8, 0x23, 0x98, 0x76, 0xd2, 0xec, 0x6a, 0x4b, 0xe7, 0x63, 0xd1, 0xa3, 0xdf, 0xb8, 0xc3, 0x63,
	0xb8, 0xdf, 0xc2, 0x99, 0x9b, 0x08, 0x59, 0x29, 0x63, 0x65, 0x87, 0xae, 0xef, 0x86, 0xa7, 0x63,
	0x89, 0xaf, 0x2d, 0xb1, 0x85, 0x96, 0xcf, 0x04, 0xd5, 0xe4, 0x00, 0xa6, 0xc3, 0xad, 0x9e, 0x6a,
	0xfe, 0xda, 0xeb, 0xbc, 0x91, 0xf5, 0x9b, 0x5d, 0x5b, 0x1c, 0x47, 0xc5, 0x7f, 0x0c, 0xc9, 0xbc,
	0xce, 0xeb, 0x9a, 0x3d, 0x83, 0x09, 0x2e, 0xb5, 0xb7, 0xe6, 0x7e, 0x2f, 0xdf, 0x87, 0x42, 0xd0,
	0x10, 0x83, 0x50, 0x76, 0xba, 0x79, 0x4f, 0x41, 0x40, 0x76, 0x0f, 0x71, 0x0d, 0x56, 0x63, 0x3f,
	0x72, 0x6b, 0x20, 0x40, 0xf7, 0x62, 0x96, 0x30, 0x8b, 0x78, 0xef, 0xca, 0xf5, 0x0a, 0xec, 0xb9,
	0x10, 0x52, 0x8d, 0xbc, 0x8d, 0xa9, 0xc8, 0xbe, 0x69, 0x1a, 0x64, 0xf7, 0x84, 0x87, 0xfc, 0x7b,
	0xe0, 0x9d, 0x42, 0x95, 0x57, 0x79, 0xb1, 0xf9, 0x6f, 0x4a, 0x0b, 0xaf, 0xf4, 0x0f, 0xaa, 0xe2,
	0x03, 0x55, 0x9f, 0xa7, 0x74, 0xe7, 0xab, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x43, 0x75, 0x8d,
	0x41, 0xf0, 0x03, 0x00, 0x00,
=======
	// 419 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0x4d, 0xab, 0xd3, 0x40,
	0x14, 0x25, 0x49, 0x53, 0x5f, 0xaf, 0xd6, 0x8f, 0x41, 0x24, 0xb8, 0x90, 0x30, 0xcf, 0x45, 0x56,
	0x0f, 0x54, 0x7c, 0x3f, 0xe0, 0x3d, 0x11, 0x2a, 0x2a, 0x38, 0x6d, 0x5d, 0xb8, 0x8b, 0xe9, 0x24,
	0x0c, 0x4d, 0x32, 0x25, 0x33, 0x11, 0xfc, 0x1d, 0xfe, 0x04, 0xfd, 0xa1, 0x72, 0xef, 0x64, 0xd2,
	0x54, 0x14, 0xba, 0x71, 0x95, 0x7b, 0xce, 0x24, 0x67, 0x4e, 0xce, 0xb9, 0xb0, 0x68, 0x4c, 0x75,
	0x75, 0xe8, 0xb4, 0xd5, 0x2c, 0xa6, 0x07, 0xff, 0x15, 0x00, 0xac, 0x65, 0xf7, 0x4d, 0x76, 0xab,
	0xb6, 0xd4, 0xec, 0x29, 0x5c, 0x18, 0x87, 0x76, 0x49, 0x90, 0x06, 0x59, 0x24, 0x46, 0xcc, 0x9e,
	0x01, 0xb8, 0x79, 0xf3, 0xfd, 0x20, 0x93, 0x30, 0x0d, 0xb2, 0x58, 0x4c, 0x18, 0xf6, 0x1c, 0x96,
	0x0e, 0x7d, 0x96, 0x9d, 0x51, 0xba, 0x4d, 0x22, 0x7a, 0xe5, 0x94, 0x64, 0x8f, 0x21, 0x36, 0x36,
	0xb7, 0x32, 0x99, 0x91, 0xbc, 0x03, 0x47, 0xed, 0x8f, 0x79, 0x23, 0x93, 0x38, 0x0d, 0xb2, 0x85,
	0x98, 0x30, 0xfc, 0x1a, 0x2e, 0xb6, 0x66, 0xf0, 0xf8, 0x10, 0xa2, 0xed, 0x68, 0x0f, 0x47, 0x96,
	0xc0, 0x1d, 0x3c, 0x35, 0xab, 0x5d, 0x12, 0xa6, 0x51, 0x16, 0x09, 0x0f, 0xf9, 0x25, 0x2c, 0x6e,
	0x75, 0xdb, 0xca, 0xc2, 0x8a, 0x35, 0x7b, 0x02, 0xf3, 0x4e, 0x9a, 0xbe, 0xb6, 0xf4, 0x6d, 0x2c,
	0x06, 0xc4, 0xaf, 0x7d, 0x04, 0x1b, 0x55, 0xec, 0x19, 0x83, 0x99, 0x55, 0x8d, 0x1c, 0xf4, 0x69,
	0x3e, 0x9a, 0x0e, 0x27, 0xa6, 0xf9, 0x0d, 0x2c, 0xdd, 0x77, 0xef, 0x75, 0xa5, 0x5a, 0xf1, 0x89,
	0xbd, 0xf0, 0x7f, 0xa1, 0xda, 0x52, 0x93, 0xc0, 0xdd, 0x97, 0x8f, 0x5c, 0xde, 0x57, 0xc7, 0x90,
	0xc5, 0xe4, 0x25, 0xfe, 0xe5, 0x54, 0xe3, 0x9f, 0x26, 0xff, 0xd0, 0x0e, 0xcf, 0xd1, 0x7e, 0x07,
	0xf7, 0xdd, 0x89, 0x90, 0x95, 0x32, 0x56, 0x76, 0x18, 0x5d, 0xaf, 0xc6, 0xe8, 0x7a, 0xb5, 0xc3,
	0xbf, 0x2d, 0x91, 0xc2, 0xdc, 0x96, 0x82, 0x66, 0x4a, 0x00, 0x2b, 0x76, 0xfd, 0xd1, 0xcc, 0x5f,
	0x7b, 0x9f, 0x6f, 0x64, 0xfd, 0xb6, 0x6f, 0x8b, 0xf3, 0xa4, 0xf8, 0xcf, 0x71, 0xbd, 0x6e, 0xf3,
	0xba, 0x66, 0x97, 0x30, 0xc3, 0x66, 0x86, 0x68, 0x1e, 0x0c, 0xf6, 0x7d, 0xb3, 0x82, 0x0e, 0xb1,
	0xcd, 0xb2, 0xd3, 0xcd, 0x07, 0x6a, 0x13, 0xd5, 0x3d, 0xc4, 0x1a, 0xac, 0x46, 0x3e, 0x72, 0x35,
	0x10, 0xa0, 0x7b, 0x71, 0x21, 0x70, 0xa1, 0xf0, 0xde, 0xad, 0xe3, 0x0a, 0xe4, 0xdc, 0x26, 0xd1,
	0x8c, 0xba, 0x8d, 0xa9, 0x28, 0xbe, 0x79, 0x1a, 0x64, 0xf7, 0x84, 0x87, 0xfc, 0x47, 0xe0, 0x93,
	0x42, 0x97, 0x37, 0x79, 0xb1, 0xff, 0x6f, 0x4e, 0x0b, 0xef, 0xf4, 0x2f, 0xae, 0xe2, 0x13, 0x57,
	0x5f, 0xe7, 0x74, 0xe7, 0xab, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6e, 0x54, 0x3b, 0xcf, 0xb5,
	0x03, 0x00, 0x00,
>>>>>>> c9145dbb314b8f9dafdcbd3aa6c14cf2edd80729
}
