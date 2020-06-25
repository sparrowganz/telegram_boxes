// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: services/admin/protobuf/admin.proto

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

type SendErrorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Status   string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Error    string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SendErrorRequest) Reset() {
	*x = SendErrorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_admin_protobuf_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendErrorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendErrorRequest) ProtoMessage() {}

func (x *SendErrorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_admin_protobuf_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendErrorRequest.ProtoReflect.Descriptor instead.
func (*SendErrorRequest) Descriptor() ([]byte, []int) {
	return file_services_admin_protobuf_admin_proto_rawDescGZIP(), []int{0}
}

func (x *SendErrorRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SendErrorRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *SendErrorRequest) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type SendErrorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendErrorResponse) Reset() {
	*x = SendErrorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_admin_protobuf_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendErrorResponse) ProtoMessage() {}

func (x *SendErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_admin_protobuf_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendErrorResponse.ProtoReflect.Descriptor instead.
func (*SendErrorResponse) Descriptor() ([]byte, []int) {
	return file_services_admin_protobuf_admin_proto_rawDescGZIP(), []int{1}
}

var File_services_admin_protobuf_admin_proto protoreflect.FileDescriptor

var file_services_admin_protobuf_admin_proto_rawDesc = []byte{
	0x0a, 0x23, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22,
	0x5c, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x13, 0x0a,
	0x11, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x32, 0x4f, 0x0a, 0x05, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x46, 0x0a, 0x09, 0x53,
	0x65, 0x6e, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x19, 0x5a, 0x17, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_admin_protobuf_admin_proto_rawDescOnce sync.Once
	file_services_admin_protobuf_admin_proto_rawDescData = file_services_admin_protobuf_admin_proto_rawDesc
)

func file_services_admin_protobuf_admin_proto_rawDescGZIP() []byte {
	file_services_admin_protobuf_admin_proto_rawDescOnce.Do(func() {
		file_services_admin_protobuf_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_admin_protobuf_admin_proto_rawDescData)
	})
	return file_services_admin_protobuf_admin_proto_rawDescData
}

var file_services_admin_protobuf_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_services_admin_protobuf_admin_proto_goTypes = []interface{}{
	(*SendErrorRequest)(nil),  // 0: protobuf.SendErrorRequest
	(*SendErrorResponse)(nil), // 1: protobuf.SendErrorResponse
}
var file_services_admin_protobuf_admin_proto_depIdxs = []int32{
	0, // 0: protobuf.Admin.SendError:input_type -> protobuf.SendErrorRequest
	1, // 1: protobuf.Admin.SendError:output_type -> protobuf.SendErrorResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_services_admin_protobuf_admin_proto_init() }
func file_services_admin_protobuf_admin_proto_init() {
	if File_services_admin_protobuf_admin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_admin_protobuf_admin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendErrorRequest); i {
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
		file_services_admin_protobuf_admin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendErrorResponse); i {
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
			RawDescriptor: file_services_admin_protobuf_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_admin_protobuf_admin_proto_goTypes,
		DependencyIndexes: file_services_admin_protobuf_admin_proto_depIdxs,
		MessageInfos:      file_services_admin_protobuf_admin_proto_msgTypes,
	}.Build()
	File_services_admin_protobuf_admin_proto = out.File
	file_services_admin_protobuf_admin_proto_rawDesc = nil
	file_services_admin_protobuf_admin_proto_goTypes = nil
	file_services_admin_protobuf_admin_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AdminClient is the client API for Admin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdminClient interface {
	SendError(ctx context.Context, in *SendErrorRequest, opts ...grpc.CallOption) (*SendErrorResponse, error)
}

type adminClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminClient(cc grpc.ClientConnInterface) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) SendError(ctx context.Context, in *SendErrorRequest, opts ...grpc.CallOption) (*SendErrorResponse, error) {
	out := new(SendErrorResponse)
	err := c.cc.Invoke(ctx, "/protobuf.Admin/SendError", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServer is the server API for Admin service.
type AdminServer interface {
	SendError(context.Context, *SendErrorRequest) (*SendErrorResponse, error)
}

// UnimplementedAdminServer can be embedded to have forward compatible implementations.
type UnimplementedAdminServer struct {
}

func (*UnimplementedAdminServer) SendError(context.Context, *SendErrorRequest) (*SendErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendError not implemented")
}

func RegisterAdminServer(s *grpc.Server, srv AdminServer) {
	s.RegisterService(&_Admin_serviceDesc, srv)
}

func _Admin_SendError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendErrorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).SendError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Admin/SendError",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).SendError(ctx, req.(*SendErrorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Admin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendError",
			Handler:    _Admin_SendError_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/admin/protobuf/admin.proto",
}
