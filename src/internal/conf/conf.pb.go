// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: src/internal/conf/conf.proto

package conf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Bootstrap struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Server        *Server                `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	Mongodb       *MongoDbConnection     `protobuf:"bytes,2,opt,name=mongodb,proto3" json:"mongodb,omitempty"`
	GrpcServer    *GRPCServer            `protobuf:"bytes,3,opt,name=grpc_server,json=grpcServer,proto3" json:"grpc_server,omitempty"`
	Consul        *Consul                `protobuf:"bytes,4,opt,name=consul,proto3" json:"consul,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Bootstrap) Reset() {
	*x = Bootstrap{}
	mi := &file_src_internal_conf_conf_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Bootstrap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bootstrap) ProtoMessage() {}

func (x *Bootstrap) ProtoReflect() protoreflect.Message {
	mi := &file_src_internal_conf_conf_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bootstrap.ProtoReflect.Descriptor instead.
func (*Bootstrap) Descriptor() ([]byte, []int) {
	return file_src_internal_conf_conf_proto_rawDescGZIP(), []int{0}
}

func (x *Bootstrap) GetServer() *Server {
	if x != nil {
		return x.Server
	}
	return nil
}

func (x *Bootstrap) GetMongodb() *MongoDbConnection {
	if x != nil {
		return x.Mongodb
	}
	return nil
}

func (x *Bootstrap) GetGrpcServer() *GRPCServer {
	if x != nil {
		return x.GrpcServer
	}
	return nil
}

func (x *Bootstrap) GetConsul() *Consul {
	if x != nil {
		return x.Consul
	}
	return nil
}

type Server struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Http          *Server_HTTP           `protobuf:"bytes,1,opt,name=http,proto3" json:"http,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Server) Reset() {
	*x = Server{}
	mi := &file_src_internal_conf_conf_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_src_internal_conf_conf_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_src_internal_conf_conf_proto_rawDescGZIP(), []int{1}
}

func (x *Server) GetHttp() *Server_HTTP {
	if x != nil {
		return x.Http
	}
	return nil
}

type MongoDbConnection struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Uri           string                 `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	Database      string                 `protobuf:"bytes,2,opt,name=database,proto3" json:"database,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MongoDbConnection) Reset() {
	*x = MongoDbConnection{}
	mi := &file_src_internal_conf_conf_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MongoDbConnection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MongoDbConnection) ProtoMessage() {}

func (x *MongoDbConnection) ProtoReflect() protoreflect.Message {
	mi := &file_src_internal_conf_conf_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MongoDbConnection.ProtoReflect.Descriptor instead.
func (*MongoDbConnection) Descriptor() ([]byte, []int) {
	return file_src_internal_conf_conf_proto_rawDescGZIP(), []int{2}
}

func (x *MongoDbConnection) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

func (x *MongoDbConnection) GetDatabase() string {
	if x != nil {
		return x.Database
	}
	return ""
}

type GRPCServer struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Port          int32                  `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Host          string                 `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GRPCServer) Reset() {
	*x = GRPCServer{}
	mi := &file_src_internal_conf_conf_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GRPCServer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GRPCServer) ProtoMessage() {}

func (x *GRPCServer) ProtoReflect() protoreflect.Message {
	mi := &file_src_internal_conf_conf_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GRPCServer.ProtoReflect.Descriptor instead.
func (*GRPCServer) Descriptor() ([]byte, []int) {
	return file_src_internal_conf_conf_proto_rawDescGZIP(), []int{3}
}

func (x *GRPCServer) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *GRPCServer) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

type Consul struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Address       string                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Consul) Reset() {
	*x = Consul{}
	mi := &file_src_internal_conf_conf_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Consul) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Consul) ProtoMessage() {}

func (x *Consul) ProtoReflect() protoreflect.Message {
	mi := &file_src_internal_conf_conf_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Consul.ProtoReflect.Descriptor instead.
func (*Consul) Descriptor() ([]byte, []int) {
	return file_src_internal_conf_conf_proto_rawDescGZIP(), []int{4}
}

func (x *Consul) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type Server_HTTP struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Network       string                 `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	Addr          string                 `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Timeout       *durationpb.Duration   `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Server_HTTP) Reset() {
	*x = Server_HTTP{}
	mi := &file_src_internal_conf_conf_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Server_HTTP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server_HTTP) ProtoMessage() {}

func (x *Server_HTTP) ProtoReflect() protoreflect.Message {
	mi := &file_src_internal_conf_conf_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server_HTTP.ProtoReflect.Descriptor instead.
func (*Server_HTTP) Descriptor() ([]byte, []int) {
	return file_src_internal_conf_conf_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Server_HTTP) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Server_HTTP) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Server_HTTP) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

var File_src_internal_conf_conf_proto protoreflect.FileDescriptor

var file_src_internal_conf_conf_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x73, 0x72, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63,
	0x6f, 0x6e, 0x66, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd5, 0x01, 0x0a, 0x09, 0x42,
	0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x12, 0x2a, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f,
	0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x06, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x07, 0x6d, 0x6f, 0x6e, 0x67, 0x6f, 0x64, 0x62, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x4d, 0x6f, 0x6e, 0x67, 0x6f, 0x44, 0x62, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6d, 0x6f, 0x6e, 0x67, 0x6f, 0x64, 0x62, 0x12, 0x37, 0x0a,
	0x0b, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x47, 0x52, 0x50, 0x43, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x0a, 0x67, 0x72, 0x70, 0x63,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6c, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x73,
	0x75, 0x6c, 0x22, 0xa0, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x2b, 0x0a,
	0x04, 0x68, 0x74, 0x74, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6b, 0x72,
	0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e,
	0x48, 0x54, 0x54, 0x50, 0x52, 0x04, 0x68, 0x74, 0x74, 0x70, 0x1a, 0x69, 0x0a, 0x04, 0x48, 0x54,
	0x54, 0x50, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x12, 0x0a, 0x04,
	0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x12, 0x33, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x41, 0x0a, 0x11, 0x4d, 0x6f, 0x6e, 0x67, 0x6f, 0x44, 0x62,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x12, 0x1a, 0x0a, 0x08,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x22, 0x34, 0x0a, 0x0a, 0x47, 0x52, 0x50, 0x43,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f,
	0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0x22,
	0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x42, 0x18, 0x5a, 0x16, 0x73, 0x72, 0x63, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_src_internal_conf_conf_proto_rawDescOnce sync.Once
	file_src_internal_conf_conf_proto_rawDescData []byte
)

func file_src_internal_conf_conf_proto_rawDescGZIP() []byte {
	file_src_internal_conf_conf_proto_rawDescOnce.Do(func() {
		file_src_internal_conf_conf_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_src_internal_conf_conf_proto_rawDesc), len(file_src_internal_conf_conf_proto_rawDesc)))
	})
	return file_src_internal_conf_conf_proto_rawDescData
}

var file_src_internal_conf_conf_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_src_internal_conf_conf_proto_goTypes = []any{
	(*Bootstrap)(nil),           // 0: kratos.api.Bootstrap
	(*Server)(nil),              // 1: kratos.api.Server
	(*MongoDbConnection)(nil),   // 2: kratos.api.MongoDbConnection
	(*GRPCServer)(nil),          // 3: kratos.api.GRPCServer
	(*Consul)(nil),              // 4: kratos.api.Consul
	(*Server_HTTP)(nil),         // 5: kratos.api.Server.HTTP
	(*durationpb.Duration)(nil), // 6: google.protobuf.Duration
}
var file_src_internal_conf_conf_proto_depIdxs = []int32{
	1, // 0: kratos.api.Bootstrap.server:type_name -> kratos.api.Server
	2, // 1: kratos.api.Bootstrap.mongodb:type_name -> kratos.api.MongoDbConnection
	3, // 2: kratos.api.Bootstrap.grpc_server:type_name -> kratos.api.GRPCServer
	4, // 3: kratos.api.Bootstrap.consul:type_name -> kratos.api.Consul
	5, // 4: kratos.api.Server.http:type_name -> kratos.api.Server.HTTP
	6, // 5: kratos.api.Server.HTTP.timeout:type_name -> google.protobuf.Duration
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_src_internal_conf_conf_proto_init() }
func file_src_internal_conf_conf_proto_init() {
	if File_src_internal_conf_conf_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_src_internal_conf_conf_proto_rawDesc), len(file_src_internal_conf_conf_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_src_internal_conf_conf_proto_goTypes,
		DependencyIndexes: file_src_internal_conf_conf_proto_depIdxs,
		MessageInfos:      file_src_internal_conf_conf_proto_msgTypes,
	}.Build()
	File_src_internal_conf_conf_proto = out.File
	file_src_internal_conf_conf_proto_goTypes = nil
	file_src_internal_conf_conf_proto_depIdxs = nil
}
