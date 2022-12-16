// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.5
// source: msg.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 交易方向
type TradingType int32

const (
	TradingType_BUYIN   TradingType = 0
	TradingType_SELLOUT TradingType = 1
)

// Enum value maps for TradingType.
var (
	TradingType_name = map[int32]string{
		0: "BUYIN",
		1: "SELLOUT",
	}
	TradingType_value = map[string]int32{
		"BUYIN":   0,
		"SELLOUT": 1,
	}
)

func (x TradingType) Enum() *TradingType {
	p := new(TradingType)
	*p = x
	return p
}

func (x TradingType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TradingType) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_proto_enumTypes[0].Descriptor()
}

func (TradingType) Type() protoreflect.EnumType {
	return &file_msg_proto_enumTypes[0]
}

func (x TradingType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TradingType.Descriptor instead.
func (TradingType) EnumDescriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

// 没有这只股票（0）：210，新增记录（1）：211，增加值（2）：212，减少值（3）213，
// 删除记录（4）: 214，没有持股可以卖（5）:215，没有足量的持股（6）:216
type ResultType int32

const (
	ResultType_NOSTOCK  ResultType = 0
	ResultType_ADDSTOCK ResultType = 1
	ResultType_ADDVALUE ResultType = 2
	ResultType_DECVALUE ResultType = 3
	ResultType_DELSTOCK ResultType = 4
	ResultType_NORECORD ResultType = 5
	ResultType_NOEHOUGH ResultType = 6
)

// Enum value maps for ResultType.
var (
	ResultType_name = map[int32]string{
		0: "NOSTOCK",
		1: "ADDSTOCK",
		2: "ADDVALUE",
		3: "DECVALUE",
		4: "DELSTOCK",
		5: "NORECORD",
		6: "NOEHOUGH",
	}
	ResultType_value = map[string]int32{
		"NOSTOCK":  0,
		"ADDSTOCK": 1,
		"ADDVALUE": 2,
		"DECVALUE": 3,
		"DELSTOCK": 4,
		"NORECORD": 5,
		"NOEHOUGH": 6,
	}
)

func (x ResultType) Enum() *ResultType {
	p := new(ResultType)
	*p = x
	return p
}

func (x ResultType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResultType) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_proto_enumTypes[1].Descriptor()
}

func (ResultType) Type() protoreflect.EnumType {
	return &file_msg_proto_enumTypes[1]
}

func (x ResultType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResultType.Descriptor instead.
func (ResultType) EnumDescriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

// 同步用户的Uid
type SyncUid struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid int32 `protobuf:"varint,1,opt,name=Uid,proto3" json:"Uid,omitempty"`
}

func (x *SyncUid) Reset() {
	*x = SyncUid{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncUid) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncUid) ProtoMessage() {}

func (x *SyncUid) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncUid.ProtoReflect.Descriptor instead.
func (*SyncUid) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

func (x *SyncUid) GetUid() int32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

// 拉取客户的持仓列表
type PullHoldingInfos struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cid string `protobuf:"bytes,1,opt,name=Cid,proto3" json:"Cid,omitempty"`
}

func (x *PullHoldingInfos) Reset() {
	*x = PullHoldingInfos{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullHoldingInfos) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullHoldingInfos) ProtoMessage() {}

func (x *PullHoldingInfos) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullHoldingInfos.ProtoReflect.Descriptor instead.
func (*PullHoldingInfos) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

func (x *PullHoldingInfos) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

// 响应客户的持仓列表
type AckHoldingInfos struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HInfo []*AckHoldingInfo `protobuf:"bytes,1,rep,name=HInfo,proto3" json:"HInfo,omitempty"`
}

func (x *AckHoldingInfos) Reset() {
	*x = AckHoldingInfos{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckHoldingInfos) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckHoldingInfos) ProtoMessage() {}

func (x *AckHoldingInfos) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckHoldingInfos.ProtoReflect.Descriptor instead.
func (*AckHoldingInfos) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{2}
}

func (x *AckHoldingInfos) GetHInfo() []*AckHoldingInfo {
	if x != nil {
		return x.HInfo
	}
	return nil
}

// 客户持仓
type AckHoldingInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string  `protobuf:"bytes,1,opt,name=ClientId,proto3" json:"ClientId,omitempty"`
	ExchangeType string  `protobuf:"bytes,2,opt,name=exchangeType,proto3" json:"exchangeType,omitempty"` //市场，1表示沪，2表示深
	StockCode    string  `protobuf:"bytes,3,opt,name=stockCode,proto3" json:"stockCode,omitempty"`       //代码	沪市：600000-600099，深市：000001-000100
	LastPrice    float64 `protobuf:"fixed64,4,opt,name=lastPrice,proto3" json:"lastPrice,omitempty"`     //最新价格 10.00-1000.00之间的随机数
	HoldAmount   int32   `protobuf:"varint,5,opt,name=holdAmount,proto3" json:"holdAmount,omitempty"`    //持仓数量
	MarketValue  float64 `protobuf:"fixed64,6,opt,name=marketValue,proto3" json:"marketValue,omitempty"` //市值 holdAmount * stock.GetLastPrice()
}

func (x *AckHoldingInfo) Reset() {
	*x = AckHoldingInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckHoldingInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckHoldingInfo) ProtoMessage() {}

func (x *AckHoldingInfo) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckHoldingInfo.ProtoReflect.Descriptor instead.
func (*AckHoldingInfo) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{3}
}

func (x *AckHoldingInfo) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *AckHoldingInfo) GetExchangeType() string {
	if x != nil {
		return x.ExchangeType
	}
	return ""
}

func (x *AckHoldingInfo) GetStockCode() string {
	if x != nil {
		return x.StockCode
	}
	return ""
}

func (x *AckHoldingInfo) GetLastPrice() float64 {
	if x != nil {
		return x.LastPrice
	}
	return 0
}

func (x *AckHoldingInfo) GetHoldAmount() int32 {
	if x != nil {
		return x.HoldAmount
	}
	return 0
}

func (x *AckHoldingInfo) GetMarketValue() float64 {
	if x != nil {
		return x.MarketValue
	}
	return 0
}

// 用户执行交易命令
type TradingInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string      `protobuf:"bytes,1,opt,name=ClientId,proto3" json:"ClientId,omitempty"`
	ExchangeType string      `protobuf:"bytes,2,opt,name=exchangeType,proto3" json:"exchangeType,omitempty"`
	StockCode    string      `protobuf:"bytes,3,opt,name=stockCode,proto3" json:"stockCode,omitempty"`
	TraType      TradingType `protobuf:"varint,4,opt,name=traType,proto3,enum=pb.TradingType" json:"traType,omitempty"`
	Number       int32       `protobuf:"varint,5,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *TradingInfo) Reset() {
	*x = TradingInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradingInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradingInfo) ProtoMessage() {}

func (x *TradingInfo) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradingInfo.ProtoReflect.Descriptor instead.
func (*TradingInfo) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{4}
}

func (x *TradingInfo) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *TradingInfo) GetExchangeType() string {
	if x != nil {
		return x.ExchangeType
	}
	return ""
}

func (x *TradingInfo) GetStockCode() string {
	if x != nil {
		return x.StockCode
	}
	return ""
}

func (x *TradingInfo) GetTraType() TradingType {
	if x != nil {
		return x.TraType
	}
	return TradingType_BUYIN
}

func (x *TradingInfo) GetNumber() int32 {
	if x != nil {
		return x.Number
	}
	return 0
}

// 用户执行命令之后的结果，返回的是一个股票，因为数量和用户的情况在客户端是相同的
type ChangeStock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RstType      ResultType `protobuf:"varint,1,opt,name=rstType,proto3,enum=pb.ResultType" json:"rstType,omitempty"`
	ExchangeType string     `protobuf:"bytes,2,opt,name=exchangeType,proto3" json:"exchangeType,omitempty"`
	StockCode    string     `protobuf:"bytes,3,opt,name=stockCode,proto3" json:"stockCode,omitempty"`
	LastPrice    float64    `protobuf:"fixed64,4,opt,name=lastPrice,proto3" json:"lastPrice,omitempty"`
	Number       int32      `protobuf:"varint,5,opt,name=number,proto3" json:"number,omitempty"`
}

func (x *ChangeStock) Reset() {
	*x = ChangeStock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeStock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeStock) ProtoMessage() {}

func (x *ChangeStock) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeStock.ProtoReflect.Descriptor instead.
func (*ChangeStock) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{5}
}

func (x *ChangeStock) GetRstType() ResultType {
	if x != nil {
		return x.RstType
	}
	return ResultType_NOSTOCK
}

func (x *ChangeStock) GetExchangeType() string {
	if x != nil {
		return x.ExchangeType
	}
	return ""
}

func (x *ChangeStock) GetStockCode() string {
	if x != nil {
		return x.StockCode
	}
	return ""
}

func (x *ChangeStock) GetLastPrice() float64 {
	if x != nil {
		return x.LastPrice
	}
	return 0
}

func (x *ChangeStock) GetNumber() int32 {
	if x != nil {
		return x.Number
	}
	return 0
}

// 拉取历史记录
type PullHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cid string `protobuf:"bytes,1,opt,name=Cid,proto3" json:"Cid,omitempty"`
}

func (x *PullHistory) Reset() {
	*x = PullHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PullHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PullHistory) ProtoMessage() {}

func (x *PullHistory) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PullHistory.ProtoReflect.Descriptor instead.
func (*PullHistory) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{6}
}

func (x *PullHistory) GetCid() string {
	if x != nil {
		return x.Cid
	}
	return ""
}

// 历史交易记录
type HistoryRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HisRcd string `protobuf:"bytes,1,opt,name=hisRcd,proto3" json:"hisRcd,omitempty"`
}

func (x *HistoryRecord) Reset() {
	*x = HistoryRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistoryRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryRecord) ProtoMessage() {}

func (x *HistoryRecord) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryRecord.ProtoReflect.Descriptor instead.
func (*HistoryRecord) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{7}
}

func (x *HistoryRecord) GetHisRcd() string {
	if x != nil {
		return x.HisRcd
	}
	return ""
}

var File_msg_proto protoreflect.FileDescriptor

var file_msg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22,
	0x1b, 0x0a, 0x07, 0x53, 0x79, 0x6e, 0x63, 0x55, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x55, 0x69, 0x64, 0x22, 0x24, 0x0a, 0x10,
	0x50, 0x75, 0x6c, 0x6c, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x43, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x43,
	0x69, 0x64, 0x22, 0x3b, 0x0a, 0x0f, 0x41, 0x63, 0x6b, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67,
	0x49, 0x6e, 0x66, 0x6f, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x48, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x41, 0x63, 0x6b, 0x48, 0x6f, 0x6c,
	0x64, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x48, 0x49, 0x6e, 0x66, 0x6f, 0x22,
	0xce, 0x01, 0x0a, 0x0e, 0x41, 0x63, 0x6b, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x22,
	0x0a, 0x0c, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x68, 0x6f, 0x6c, 0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x68, 0x6f, 0x6c, 0x64, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x20,
	0x0a, 0x0b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x0b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0xae, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c,
	0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x29,
	0x0a, 0x07, 0x74, 0x72, 0x61, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x07, 0x74, 0x72, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x22, 0xaf, 0x01, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x74, 0x6f, 0x63,
	0x6b, 0x12, 0x28, 0x0a, 0x07, 0x72, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x07, 0x72, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x65,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x22, 0x1f, 0x0a, 0x0b, 0x50, 0x75, 0x6c, 0x6c, 0x48, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x43, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x43, 0x69, 0x64, 0x22, 0x27, 0x0a, 0x0d, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x69, 0x73, 0x52, 0x63, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x68, 0x69, 0x73, 0x52, 0x63, 0x64, 0x2a, 0x25, 0x0a,
	0x0b, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05,
	0x42, 0x55, 0x59, 0x49, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x45, 0x4c, 0x4c, 0x4f,
	0x55, 0x54, 0x10, 0x01, 0x2a, 0x6d, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x4f, 0x53, 0x54, 0x4f, 0x43, 0x4b, 0x10, 0x00, 0x12,
	0x0c, 0x0a, 0x08, 0x41, 0x44, 0x44, 0x53, 0x54, 0x4f, 0x43, 0x4b, 0x10, 0x01, 0x12, 0x0c, 0x0a,
	0x08, 0x41, 0x44, 0x44, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x44,
	0x45, 0x43, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x4c,
	0x53, 0x54, 0x4f, 0x43, 0x4b, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x4f, 0x52, 0x45, 0x43,
	0x4f, 0x52, 0x44, 0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x4f, 0x45, 0x48, 0x4f, 0x55, 0x47,
	0x48, 0x10, 0x06, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_proto_rawDescOnce sync.Once
	file_msg_proto_rawDescData = file_msg_proto_rawDesc
)

func file_msg_proto_rawDescGZIP() []byte {
	file_msg_proto_rawDescOnce.Do(func() {
		file_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_rawDescData)
	})
	return file_msg_proto_rawDescData
}

var file_msg_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_msg_proto_goTypes = []interface{}{
	(TradingType)(0),         // 0: pb.TradingType
	(ResultType)(0),          // 1: pb.ResultType
	(*SyncUid)(nil),          // 2: pb.SyncUid
	(*PullHoldingInfos)(nil), // 3: pb.PullHoldingInfos
	(*AckHoldingInfos)(nil),  // 4: pb.AckHoldingInfos
	(*AckHoldingInfo)(nil),   // 5: pb.AckHoldingInfo
	(*TradingInfo)(nil),      // 6: pb.TradingInfo
	(*ChangeStock)(nil),      // 7: pb.ChangeStock
	(*PullHistory)(nil),      // 8: pb.PullHistory
	(*HistoryRecord)(nil),    // 9: pb.HistoryRecord
}
var file_msg_proto_depIdxs = []int32{
	5, // 0: pb.AckHoldingInfos.HInfo:type_name -> pb.AckHoldingInfo
	0, // 1: pb.TradingInfo.traType:type_name -> pb.TradingType
	1, // 2: pb.ChangeStock.rstType:type_name -> pb.ResultType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_msg_proto_init() }
func file_msg_proto_init() {
	if File_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncUid); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullHoldingInfos); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckHoldingInfos); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckHoldingInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TradingInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeStock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PullHistory); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msg_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistoryRecord); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msg_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_proto_goTypes,
		DependencyIndexes: file_msg_proto_depIdxs,
		EnumInfos:         file_msg_proto_enumTypes,
		MessageInfos:      file_msg_proto_msgTypes,
	}.Build()
	File_msg_proto = out.File
	file_msg_proto_rawDesc = nil
	file_msg_proto_goTypes = nil
	file_msg_proto_depIdxs = nil
}