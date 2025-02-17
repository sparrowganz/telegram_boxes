// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: services/box/protobuf/box.proto

package protobuf

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type StartBroadcastRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     string    `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	FileLink string    `protobuf:"bytes,2,opt,name=fileLink,proto3" json:"fileLink,omitempty"`
	Buttons  []*Button `protobuf:"bytes,3,rep,name=buttons,proto3" json:"buttons,omitempty"`
	Text     string    `protobuf:"bytes,4,opt,name=text,proto3" json:"text,omitempty"`
	ChatID   int64     `protobuf:"varint,5,opt,name=chatID,proto3" json:"chatID,omitempty"`
}

func (x *StartBroadcastRequest) Reset() {
	*x = StartBroadcastRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartBroadcastRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartBroadcastRequest) ProtoMessage() {}

func (x *StartBroadcastRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartBroadcastRequest.ProtoReflect.Descriptor instead.
func (*StartBroadcastRequest) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{0}
}

func (x *StartBroadcastRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *StartBroadcastRequest) GetFileLink() string {
	if x != nil {
		return x.FileLink
	}
	return ""
}

func (x *StartBroadcastRequest) GetButtons() []*Button {
	if x != nil {
		return x.Buttons
	}
	return nil
}

func (x *StartBroadcastRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *StartBroadcastRequest) GetChatID() int64 {
	if x != nil {
		return x.ChatID
	}
	return 0
}

type Button struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Url  string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Button) Reset() {
	*x = Button{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Button) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Button) ProtoMessage() {}

func (x *Button) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Button.ProtoReflect.Descriptor instead.
func (*Button) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{1}
}

func (x *Button) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Button) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success int64 `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Fail    int64 `protobuf:"varint,2,opt,name=fail,proto3" json:"fail,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{2}
}

func (x *Stats) GetSuccess() int64 {
	if x != nil {
		return x.Success
	}
	return 0
}

func (x *Stats) GetFail() int64 {
	if x != nil {
		return x.Fail
	}
	return 0
}

//----------------------------------------------------------------------------------------------------------------------
//  RemoveCheckTask
//----------------------------------------------------------------------------------------------------------------------
type RemoveCheckTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskID string `protobuf:"bytes,1,opt,name=taskID,proto3" json:"taskID,omitempty"`
}

func (x *RemoveCheckTaskRequest) Reset() {
	*x = RemoveCheckTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveCheckTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveCheckTaskRequest) ProtoMessage() {}

func (x *RemoveCheckTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveCheckTaskRequest.ProtoReflect.Descriptor instead.
func (*RemoveCheckTaskRequest) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{3}
}

func (x *RemoveCheckTaskRequest) GetTaskID() string {
	if x != nil {
		return x.TaskID
	}
	return ""
}

type RemoveCheckTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveCheckTaskResponse) Reset() {
	*x = RemoveCheckTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveCheckTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveCheckTaskResponse) ProtoMessage() {}

func (x *RemoveCheckTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveCheckTaskResponse.ProtoReflect.Descriptor instead.
func (*RemoveCheckTaskResponse) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{4}
}

//----------------------------------------------------------------------------------------------------------------------
//  Check
//----------------------------------------------------------------------------------------------------------------------
type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChatID int64 `protobuf:"varint,1,opt,name=chatID,proto3" json:"chatID,omitempty"`
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{5}
}

func (x *CheckRequest) GetChatID() int64 {
	if x != nil {
		return x.ChatID
	}
	return 0
}

type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{6}
}

type GetStatisticsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetStatisticsRequest) Reset() {
	*x = GetStatisticsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStatisticsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatisticsRequest) ProtoMessage() {}

func (x *GetStatisticsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatisticsRequest.ProtoReflect.Descriptor instead.
func (*GetStatisticsRequest) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{7}
}

type Statistic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	All     int64 `protobuf:"varint,1,opt,name=all,proto3" json:"all,omitempty"`
	Blocked int64 `protobuf:"varint,2,opt,name=blocked,proto3" json:"blocked,omitempty"`
	Current int64 `protobuf:"varint,3,opt,name=current,proto3" json:"current,omitempty"`
}

func (x *Statistic) Reset() {
	*x = Statistic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_box_protobuf_box_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Statistic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Statistic) ProtoMessage() {}

func (x *Statistic) ProtoReflect() protoreflect.Message {
	mi := &file_services_box_protobuf_box_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Statistic.ProtoReflect.Descriptor instead.
func (*Statistic) Descriptor() ([]byte, []int) {
	return file_services_box_protobuf_box_proto_rawDescGZIP(), []int{8}
}

func (x *Statistic) GetAll() int64 {
	if x != nil {
		return x.All
	}
	return 0
}

func (x *Statistic) GetBlocked() int64 {
	if x != nil {
		return x.Blocked
	}
	return 0
}

func (x *Statistic) GetCurrent() int64 {
	if x != nil {
		return x.Current
	}
	return 0
}

var File_services_box_protobuf_box_proto protoreflect.FileDescriptor

var file_services_box_protobuf_box_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x62, 0x6f, 0x78, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x62, 0x6f, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x9f, 0x01, 0x0a, 0x15,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x2a, 0x0a, 0x07, 0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x52, 0x07, 0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x44, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x44, 0x22, 0x2e, 0x0a,
	0x06, 0x42, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x35, 0x0a,
	0x05, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04,
	0x66, 0x61, 0x69, 0x6c, 0x22, 0x30, 0x0a, 0x16, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x22, 0x19, 0x0a, 0x17, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x26, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x44, 0x22, 0x0f, 0x0a, 0x0d, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x51, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x12,
	0x10, 0x0a, 0x03, 0x61, 0x6c, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x61, 0x6c,
	0x6c, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x32, 0xab, 0x02, 0x0a, 0x03, 0x42, 0x6f, 0x78, 0x12, 0x3a, 0x0a,
	0x05, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0d, 0x47, 0x65, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74,
	0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x22,
	0x00, 0x12, 0x58, 0x0a, 0x0f, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x54, 0x61, 0x73, 0x6b, 0x12, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x61, 0x73,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0e, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x1f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x42, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x22,
	0x00, 0x30, 0x01, 0x42, 0x17, 0x5a, 0x15, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f,
	0x62, 0x6f, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_box_protobuf_box_proto_rawDescOnce sync.Once
	file_services_box_protobuf_box_proto_rawDescData = file_services_box_protobuf_box_proto_rawDesc
)

func file_services_box_protobuf_box_proto_rawDescGZIP() []byte {
	file_services_box_protobuf_box_proto_rawDescOnce.Do(func() {
		file_services_box_protobuf_box_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_box_protobuf_box_proto_rawDescData)
	})
	return file_services_box_protobuf_box_proto_rawDescData
}

var file_services_box_protobuf_box_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_services_box_protobuf_box_proto_goTypes = []interface{}{
	(*StartBroadcastRequest)(nil),   // 0: protobuf.StartBroadcastRequest
	(*Button)(nil),                  // 1: protobuf.Button
	(*Stats)(nil),                   // 2: protobuf.Stats
	(*RemoveCheckTaskRequest)(nil),  // 3: protobuf.RemoveCheckTaskRequest
	(*RemoveCheckTaskResponse)(nil), // 4: protobuf.RemoveCheckTaskResponse
	(*CheckRequest)(nil),            // 5: protobuf.CheckRequest
	(*CheckResponse)(nil),           // 6: protobuf.CheckResponse
	(*GetStatisticsRequest)(nil),    // 7: protobuf.GetStatisticsRequest
	(*Statistic)(nil),               // 8: protobuf.Statistic
}
var file_services_box_protobuf_box_proto_depIdxs = []int32{
	1, // 0: protobuf.StartBroadcastRequest.buttons:type_name -> protobuf.Button
	5, // 1: protobuf.Box.Check:input_type -> protobuf.CheckRequest
	7, // 2: protobuf.Box.GetStatistics:input_type -> protobuf.GetStatisticsRequest
	3, // 3: protobuf.Box.RemoveCheckTask:input_type -> protobuf.RemoveCheckTaskRequest
	0, // 4: protobuf.Box.StartBroadcast:input_type -> protobuf.StartBroadcastRequest
	6, // 5: protobuf.Box.Check:output_type -> protobuf.CheckResponse
	8, // 6: protobuf.Box.GetStatistics:output_type -> protobuf.Statistic
	4, // 7: protobuf.Box.RemoveCheckTask:output_type -> protobuf.RemoveCheckTaskResponse
	2, // 8: protobuf.Box.StartBroadcast:output_type -> protobuf.Stats
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_services_box_protobuf_box_proto_init() }
func file_services_box_protobuf_box_proto_init() {
	if File_services_box_protobuf_box_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_box_protobuf_box_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartBroadcastRequest); i {
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
		file_services_box_protobuf_box_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Button); i {
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
		file_services_box_protobuf_box_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stats); i {
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
		file_services_box_protobuf_box_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveCheckTaskRequest); i {
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
		file_services_box_protobuf_box_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveCheckTaskResponse); i {
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
		file_services_box_protobuf_box_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
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
		file_services_box_protobuf_box_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
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
		file_services_box_protobuf_box_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStatisticsRequest); i {
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
		file_services_box_protobuf_box_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Statistic); i {
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
			RawDescriptor: file_services_box_protobuf_box_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_box_protobuf_box_proto_goTypes,
		DependencyIndexes: file_services_box_protobuf_box_proto_depIdxs,
		MessageInfos:      file_services_box_protobuf_box_proto_msgTypes,
	}.Build()
	File_services_box_protobuf_box_proto = out.File
	file_services_box_protobuf_box_proto_rawDesc = nil
	file_services_box_protobuf_box_proto_goTypes = nil
	file_services_box_protobuf_box_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BoxClient is the client API for Box service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BoxClient interface {
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	GetStatistics(ctx context.Context, in *GetStatisticsRequest, opts ...grpc.CallOption) (*Statistic, error)
	RemoveCheckTask(ctx context.Context, in *RemoveCheckTaskRequest, opts ...grpc.CallOption) (*RemoveCheckTaskResponse, error)
	StartBroadcast(ctx context.Context, in *StartBroadcastRequest, opts ...grpc.CallOption) (Box_StartBroadcastClient, error)
}

type boxClient struct {
	cc grpc.ClientConnInterface
}

func NewBoxClient(cc grpc.ClientConnInterface) BoxClient {
	return &boxClient{cc}
}

func (c *boxClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/protobuf.Box/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boxClient) GetStatistics(ctx context.Context, in *GetStatisticsRequest, opts ...grpc.CallOption) (*Statistic, error) {
	out := new(Statistic)
	err := c.cc.Invoke(ctx, "/protobuf.Box/GetStatistics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boxClient) RemoveCheckTask(ctx context.Context, in *RemoveCheckTaskRequest, opts ...grpc.CallOption) (*RemoveCheckTaskResponse, error) {
	out := new(RemoveCheckTaskResponse)
	err := c.cc.Invoke(ctx, "/protobuf.Box/RemoveCheckTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boxClient) StartBroadcast(ctx context.Context, in *StartBroadcastRequest, opts ...grpc.CallOption) (Box_StartBroadcastClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Box_serviceDesc.Streams[0], "/protobuf.Box/StartBroadcast", opts...)
	if err != nil {
		return nil, err
	}
	x := &boxStartBroadcastClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Box_StartBroadcastClient interface {
	Recv() (*Stats, error)
	grpc.ClientStream
}

type boxStartBroadcastClient struct {
	grpc.ClientStream
}

func (x *boxStartBroadcastClient) Recv() (*Stats, error) {
	m := new(Stats)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BoxServer is the server API for Box service.
type BoxServer interface {
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	GetStatistics(context.Context, *GetStatisticsRequest) (*Statistic, error)
	RemoveCheckTask(context.Context, *RemoveCheckTaskRequest) (*RemoveCheckTaskResponse, error)
	StartBroadcast(*StartBroadcastRequest, Box_StartBroadcastServer) error
}

// UnimplementedBoxServer can be embedded to have forward compatible implementations.
type UnimplementedBoxServer struct {
}

func (*UnimplementedBoxServer) Check(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (*UnimplementedBoxServer) GetStatistics(context.Context, *GetStatisticsRequest) (*Statistic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatistics not implemented")
}
func (*UnimplementedBoxServer) RemoveCheckTask(context.Context, *RemoveCheckTaskRequest) (*RemoveCheckTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveCheckTask not implemented")
}
func (*UnimplementedBoxServer) StartBroadcast(*StartBroadcastRequest, Box_StartBroadcastServer) error {
	return status.Errorf(codes.Unimplemented, "method StartBroadcast not implemented")
}

func RegisterBoxServer(s *grpc.Server, srv BoxServer) {
	s.RegisterService(&_Box_serviceDesc, srv)
}

func _Box_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoxServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Box/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoxServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Box_GetStatistics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatisticsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoxServer).GetStatistics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Box/GetStatistics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoxServer).GetStatistics(ctx, req.(*GetStatisticsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Box_RemoveCheckTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveCheckTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoxServer).RemoveCheckTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Box/RemoveCheckTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoxServer).RemoveCheckTask(ctx, req.(*RemoveCheckTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Box_StartBroadcast_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StartBroadcastRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BoxServer).StartBroadcast(m, &boxStartBroadcastServer{stream})
}

type Box_StartBroadcastServer interface {
	Send(*Stats) error
	grpc.ServerStream
}

type boxStartBroadcastServer struct {
	grpc.ServerStream
}

func (x *boxStartBroadcastServer) Send(m *Stats) error {
	return x.ServerStream.SendMsg(m)
}

var _Box_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Box",
	HandlerType: (*BoxServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _Box_Check_Handler,
		},
		{
			MethodName: "GetStatistics",
			Handler:    _Box_GetStatistics_Handler,
		},
		{
			MethodName: "RemoveCheckTask",
			Handler:    _Box_RemoveCheckTask_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartBroadcast",
			Handler:       _Box_StartBroadcast_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "services/box/protobuf/box.proto",
}
