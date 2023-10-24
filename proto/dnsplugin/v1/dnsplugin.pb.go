// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: dnsplugin/v1/dnsplugin.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConfigureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProviderName         string            `protobuf:"bytes,1,opt,name=provider_name,json=providerName,proto3" json:"provider_name,omitempty"`
	Config               map[string]string `protobuf:"bytes,2,rep,name=config,proto3" json:"config,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RecursiveNameservers []string          `protobuf:"bytes,3,rep,name=recursive_nameservers,json=recursiveNameservers,proto3" json:"recursive_nameservers,omitempty"`
}

func (x *ConfigureRequest) Reset() {
	*x = ConfigureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureRequest) ProtoMessage() {}

func (x *ConfigureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigureRequest.ProtoReflect.Descriptor instead.
func (*ConfigureRequest) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigureRequest) GetProviderName() string {
	if x != nil {
		return x.ProviderName
	}
	return ""
}

func (x *ConfigureRequest) GetConfig() map[string]string {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *ConfigureRequest) GetRecursiveNameservers() []string {
	if x != nil {
		return x.RecursiveNameservers
	}
	return nil
}

type ConfigureResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConfigureResponse) Reset() {
	*x = ConfigureResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigureResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigureResponse) ProtoMessage() {}

func (x *ConfigureResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigureResponse.ProtoReflect.Descriptor instead.
func (*ConfigureResponse) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{1}
}

type PresentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Domain  string `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`
	Token   string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	KeyAuth string `protobuf:"bytes,3,opt,name=key_auth,json=keyAuth,proto3" json:"key_auth,omitempty"`
}

func (x *PresentRequest) Reset() {
	*x = PresentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PresentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresentRequest) ProtoMessage() {}

func (x *PresentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresentRequest.ProtoReflect.Descriptor instead.
func (*PresentRequest) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{2}
}

func (x *PresentRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *PresentRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *PresentRequest) GetKeyAuth() string {
	if x != nil {
		return x.KeyAuth
	}
	return ""
}

type PresentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PresentResponse) Reset() {
	*x = PresentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PresentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PresentResponse) ProtoMessage() {}

func (x *PresentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PresentResponse.ProtoReflect.Descriptor instead.
func (*PresentResponse) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{3}
}

type CleanUpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Domain  string `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`
	Token   string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	KeyAuth string `protobuf:"bytes,3,opt,name=key_auth,json=keyAuth,proto3" json:"key_auth,omitempty"`
}

func (x *CleanUpRequest) Reset() {
	*x = CleanUpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CleanUpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CleanUpRequest) ProtoMessage() {}

func (x *CleanUpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CleanUpRequest.ProtoReflect.Descriptor instead.
func (*CleanUpRequest) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{4}
}

func (x *CleanUpRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *CleanUpRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *CleanUpRequest) GetKeyAuth() string {
	if x != nil {
		return x.KeyAuth
	}
	return ""
}

type CleanUpResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CleanUpResponse) Reset() {
	*x = CleanUpResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CleanUpResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CleanUpResponse) ProtoMessage() {}

func (x *CleanUpResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CleanUpResponse.ProtoReflect.Descriptor instead.
func (*CleanUpResponse) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{5}
}

type TimeoutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TimeoutRequest) Reset() {
	*x = TimeoutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeoutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeoutRequest) ProtoMessage() {}

func (x *TimeoutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeoutRequest.ProtoReflect.Descriptor instead.
func (*TimeoutRequest) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{6}
}

type TimeoutResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timeout  *durationpb.Duration `protobuf:"bytes,1,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Interval *durationpb.Duration `protobuf:"bytes,2,opt,name=interval,proto3" json:"interval,omitempty"`
}

func (x *TimeoutResponse) Reset() {
	*x = TimeoutResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeoutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeoutResponse) ProtoMessage() {}

func (x *TimeoutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dnsplugin_v1_dnsplugin_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeoutResponse.ProtoReflect.Descriptor instead.
func (*TimeoutResponse) Descriptor() ([]byte, []int) {
	return file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP(), []int{7}
}

func (x *TimeoutResponse) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *TimeoutResponse) GetInterval() *durationpb.Duration {
	if x != nil {
		return x.Interval
	}
	return nil
}

var File_dnsplugin_v1_dnsplugin_proto protoreflect.FileDescriptor

var file_dnsplugin_v1_dnsplugin_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x64,
	0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xeb, 0x01, 0x0a,
	0x10, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x42, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x33, 0x0a, 0x15, 0x72, 0x65,
	0x63, 0x75, 0x72, 0x73, 0x69, 0x76, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x14, 0x72, 0x65, 0x63, 0x75, 0x72,
	0x73, 0x69, 0x76, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x1a,
	0x39, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x13, 0x0a, 0x11, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x59, 0x0a, 0x0e, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x19, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68, 0x22, 0x11, 0x0a, 0x0f, 0x50, 0x72,
	0x65, 0x73, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x59, 0x0a,
	0x0e, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a,
	0x08, 0x6b, 0x65, 0x79, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68, 0x22, 0x11, 0x0a, 0x0f, 0x43, 0x6c, 0x65, 0x61,
	0x6e, 0x55, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x54,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x7d, 0x0a,
	0x0f, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x33, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x32, 0xc2, 0x02, 0x0a,
	0x12, 0x44, 0x4e, 0x53, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65,
	0x12, 0x1e, 0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1f, 0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x12, 0x1c,
	0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72,
	0x65, 0x73, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x64,
	0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x65, 0x73,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a,
	0x07, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x12, 0x1c, 0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x07, 0x54, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x1c, 0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6d, 0x79, 0x6b, 0x6c, 0x73, 0x74, 0x2f, 0x74, 0x65, 0x72, 0x72, 0x61, 0x66, 0x6f, 0x72, 0x6d,
	0x2d, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2d, 0x61, 0x63, 0x6d, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x6e, 0x73, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dnsplugin_v1_dnsplugin_proto_rawDescOnce sync.Once
	file_dnsplugin_v1_dnsplugin_proto_rawDescData = file_dnsplugin_v1_dnsplugin_proto_rawDesc
)

func file_dnsplugin_v1_dnsplugin_proto_rawDescGZIP() []byte {
	file_dnsplugin_v1_dnsplugin_proto_rawDescOnce.Do(func() {
		file_dnsplugin_v1_dnsplugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_dnsplugin_v1_dnsplugin_proto_rawDescData)
	})
	return file_dnsplugin_v1_dnsplugin_proto_rawDescData
}

var file_dnsplugin_v1_dnsplugin_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_dnsplugin_v1_dnsplugin_proto_goTypes = []interface{}{
	(*ConfigureRequest)(nil),    // 0: dnsplugin.v1.ConfigureRequest
	(*ConfigureResponse)(nil),   // 1: dnsplugin.v1.ConfigureResponse
	(*PresentRequest)(nil),      // 2: dnsplugin.v1.PresentRequest
	(*PresentResponse)(nil),     // 3: dnsplugin.v1.PresentResponse
	(*CleanUpRequest)(nil),      // 4: dnsplugin.v1.CleanUpRequest
	(*CleanUpResponse)(nil),     // 5: dnsplugin.v1.CleanUpResponse
	(*TimeoutRequest)(nil),      // 6: dnsplugin.v1.TimeoutRequest
	(*TimeoutResponse)(nil),     // 7: dnsplugin.v1.TimeoutResponse
	nil,                         // 8: dnsplugin.v1.ConfigureRequest.ConfigEntry
	(*durationpb.Duration)(nil), // 9: google.protobuf.Duration
}
var file_dnsplugin_v1_dnsplugin_proto_depIdxs = []int32{
	8, // 0: dnsplugin.v1.ConfigureRequest.config:type_name -> dnsplugin.v1.ConfigureRequest.ConfigEntry
	9, // 1: dnsplugin.v1.TimeoutResponse.timeout:type_name -> google.protobuf.Duration
	9, // 2: dnsplugin.v1.TimeoutResponse.interval:type_name -> google.protobuf.Duration
	0, // 3: dnsplugin.v1.DNSProviderService.Configure:input_type -> dnsplugin.v1.ConfigureRequest
	2, // 4: dnsplugin.v1.DNSProviderService.Present:input_type -> dnsplugin.v1.PresentRequest
	4, // 5: dnsplugin.v1.DNSProviderService.CleanUp:input_type -> dnsplugin.v1.CleanUpRequest
	6, // 6: dnsplugin.v1.DNSProviderService.Timeout:input_type -> dnsplugin.v1.TimeoutRequest
	1, // 7: dnsplugin.v1.DNSProviderService.Configure:output_type -> dnsplugin.v1.ConfigureResponse
	3, // 8: dnsplugin.v1.DNSProviderService.Present:output_type -> dnsplugin.v1.PresentResponse
	5, // 9: dnsplugin.v1.DNSProviderService.CleanUp:output_type -> dnsplugin.v1.CleanUpResponse
	7, // 10: dnsplugin.v1.DNSProviderService.Timeout:output_type -> dnsplugin.v1.TimeoutResponse
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_dnsplugin_v1_dnsplugin_proto_init() }
func file_dnsplugin_v1_dnsplugin_proto_init() {
	if File_dnsplugin_v1_dnsplugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigureRequest); i {
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
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigureResponse); i {
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
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PresentRequest); i {
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
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PresentResponse); i {
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
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CleanUpRequest); i {
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
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CleanUpResponse); i {
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
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeoutRequest); i {
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
		file_dnsplugin_v1_dnsplugin_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeoutResponse); i {
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
			RawDescriptor: file_dnsplugin_v1_dnsplugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dnsplugin_v1_dnsplugin_proto_goTypes,
		DependencyIndexes: file_dnsplugin_v1_dnsplugin_proto_depIdxs,
		MessageInfos:      file_dnsplugin_v1_dnsplugin_proto_msgTypes,
	}.Build()
	File_dnsplugin_v1_dnsplugin_proto = out.File
	file_dnsplugin_v1_dnsplugin_proto_rawDesc = nil
	file_dnsplugin_v1_dnsplugin_proto_goTypes = nil
	file_dnsplugin_v1_dnsplugin_proto_depIdxs = nil
}
