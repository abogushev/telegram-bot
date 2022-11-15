// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: report.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReportResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64              `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Start  string             `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	End    string             `protobuf:"bytes,3,opt,name=end,proto3" json:"end,omitempty"`
	Data   map[string]float64 `protobuf:"bytes,4,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
}

func (x *ReportResult) Reset() {
	*x = ReportResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_report_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReportResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReportResult) ProtoMessage() {}

func (x *ReportResult) ProtoReflect() protoreflect.Message {
	mi := &file_report_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReportResult.ProtoReflect.Descriptor instead.
func (*ReportResult) Descriptor() ([]byte, []int) {
	return file_report_proto_rawDescGZIP(), []int{0}
}

func (x *ReportResult) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReportResult) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *ReportResult) GetEnd() string {
	if x != nil {
		return x.End
	}
	return ""
}

func (x *ReportResult) GetData() map[string]float64 {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_report_proto protoreflect.FileDescriptor

var file_report_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xbb, 0x01, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x65, 0x6e, 0x64, 0x12, 0x32, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x32, 0x40, 0x0a, 0x06, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x36, 0x0a, 0x04, 0x53,
	0x65, 0x6e, 0x64, 0x12, 0x14, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x7a,
	0x6f, 0x6e, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x61, 0x6c, 0x65, 0x78, 0x2e, 0x62, 0x6f, 0x67, 0x75,
	0x73, 0x68, 0x65, 0x76, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d, 0x2d, 0x62, 0x6f,
	0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_report_proto_rawDescOnce sync.Once
	file_report_proto_rawDescData = file_report_proto_rawDesc
)

func file_report_proto_rawDescGZIP() []byte {
	file_report_proto_rawDescOnce.Do(func() {
		file_report_proto_rawDescData = protoimpl.X.CompressGZIP(file_report_proto_rawDescData)
	})
	return file_report_proto_rawDescData
}

var file_report_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_report_proto_goTypes = []interface{}{
	(*ReportResult)(nil),  // 0: report.ReportResult
	nil,                   // 1: report.ReportResult.DataEntry
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_report_proto_depIdxs = []int32{
	1, // 0: report.ReportResult.data:type_name -> report.ReportResult.DataEntry
	0, // 1: report.Report.Send:input_type -> report.ReportResult
	2, // 2: report.Report.Send:output_type -> google.protobuf.Empty
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_report_proto_init() }
func file_report_proto_init() {
	if File_report_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_report_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReportResult); i {
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
			RawDescriptor: file_report_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_report_proto_goTypes,
		DependencyIndexes: file_report_proto_depIdxs,
		MessageInfos:      file_report_proto_msgTypes,
	}.Build()
	File_report_proto = out.File
	file_report_proto_rawDesc = nil
	file_report_proto_goTypes = nil
	file_report_proto_depIdxs = nil
}
