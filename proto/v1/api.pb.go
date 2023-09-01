// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: proto/v1/api.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Route int32

const (
	Route_ROUTE_UNSPECIFIED Route = 0
	Route_ROUTE_KITCHEN     Route = 1
)

// Enum value maps for Route.
var (
	Route_name = map[int32]string{
		0: "ROUTE_UNSPECIFIED",
		1: "ROUTE_KITCHEN",
	}
	Route_value = map[string]int32{
		"ROUTE_UNSPECIFIED": 0,
		"ROUTE_KITCHEN":     1,
	}
)

func (x Route) Enum() *Route {
	p := new(Route)
	*p = x
	return p
}

func (x Route) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Route) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_v1_api_proto_enumTypes[0].Descriptor()
}

func (Route) Type() protoreflect.EnumType {
	return &file_proto_v1_api_proto_enumTypes[0]
}

func (x Route) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Route.Descriptor instead.
func (Route) EnumDescriptor() ([]byte, []int) {
	return file_proto_v1_api_proto_rawDescGZIP(), []int{0}
}

type StartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Route Route `protobuf:"varint,1,opt,name=route,proto3,enum=proto.v1.Route" json:"route,omitempty"`
}

func (x *StartRequest) Reset() {
	*x = StartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartRequest) ProtoMessage() {}

func (x *StartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartRequest.ProtoReflect.Descriptor instead.
func (*StartRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_api_proto_rawDescGZIP(), []int{0}
}

func (x *StartRequest) GetRoute() Route {
	if x != nil {
		return x.Route
	}
	return Route_ROUTE_UNSPECIFIED
}

type StartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StartResponse) Reset() {
	*x = StartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartResponse) ProtoMessage() {}

func (x *StartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartResponse.ProtoReflect.Descriptor instead.
func (*StartResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_api_proto_rawDescGZIP(), []int{1}
}

var File_proto_v1_api_proto protoreflect.FileDescriptor

var file_proto_v1_api_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x0c,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x05,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x05, 0x72, 0x6f,
	0x75, 0x74, 0x65, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x31, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x15, 0x0a,
	0x11, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x5f, 0x4b, 0x49,
	0x54, 0x43, 0x48, 0x45, 0x4e, 0x10, 0x01, 0x32, 0x67, 0x0a, 0x0d, 0x41, 0x69, 0x72, 0x62, 0x6f,
	0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x22, 0x14, 0x2f, 0x61, 0x69, 0x72,
	0x62, 0x6f, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_api_proto_rawDescOnce sync.Once
	file_proto_v1_api_proto_rawDescData = file_proto_v1_api_proto_rawDesc
)

func file_proto_v1_api_proto_rawDescGZIP() []byte {
	file_proto_v1_api_proto_rawDescOnce.Do(func() {
		file_proto_v1_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_api_proto_rawDescData)
	})
	return file_proto_v1_api_proto_rawDescData
}

var file_proto_v1_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_v1_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_v1_api_proto_goTypes = []interface{}{
	(Route)(0),            // 0: proto.v1.Route
	(*StartRequest)(nil),  // 1: proto.v1.StartRequest
	(*StartResponse)(nil), // 2: proto.v1.StartResponse
}
var file_proto_v1_api_proto_depIdxs = []int32{
	0, // 0: proto.v1.StartRequest.route:type_name -> proto.v1.Route
	1, // 1: proto.v1.AirbotService.Start:input_type -> proto.v1.StartRequest
	2, // 2: proto.v1.AirbotService.Start:output_type -> proto.v1.StartResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_v1_api_proto_init() }
func file_proto_v1_api_proto_init() {
	if File_proto_v1_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartRequest); i {
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
		file_proto_v1_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartResponse); i {
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
			RawDescriptor: file_proto_v1_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_api_proto_goTypes,
		DependencyIndexes: file_proto_v1_api_proto_depIdxs,
		EnumInfos:         file_proto_v1_api_proto_enumTypes,
		MessageInfos:      file_proto_v1_api_proto_msgTypes,
	}.Build()
	File_proto_v1_api_proto = out.File
	file_proto_v1_api_proto_rawDesc = nil
	file_proto_v1_api_proto_goTypes = nil
	file_proto_v1_api_proto_depIdxs = nil
}