// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: link.proto

package ova_link_api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      uint64   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Url         string   `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Tags        []string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *CreateLinkRequest) Reset() {
	*x = CreateLinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_link_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLinkRequest) ProtoMessage() {}

func (x *CreateLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_link_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLinkRequest.ProtoReflect.Descriptor instead.
func (*CreateLinkRequest) Descriptor() ([]byte, []int) {
	return file_link_proto_rawDescGZIP(), []int{0}
}

func (x *CreateLinkRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateLinkRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CreateLinkRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateLinkRequest) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type DeleteLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteLinkRequest) Reset() {
	*x = DeleteLinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_link_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteLinkRequest) ProtoMessage() {}

func (x *DeleteLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_link_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteLinkRequest.ProtoReflect.Descriptor instead.
func (*DeleteLinkRequest) Descriptor() ([]byte, []int) {
	return file_link_proto_rawDescGZIP(), []int{1}
}

func (x *DeleteLinkRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DescribeLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DescribeLinkRequest) Reset() {
	*x = DescribeLinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_link_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeLinkRequest) ProtoMessage() {}

func (x *DescribeLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_link_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeLinkRequest.ProtoReflect.Descriptor instead.
func (*DescribeLinkRequest) Descriptor() ([]byte, []int) {
	return file_link_proto_rawDescGZIP(), []int{2}
}

func (x *DescribeLinkRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type DescribeLinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId      uint64                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Url         string                 `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Description string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Tags        []string               `protobuf:"bytes,5,rep,name=tags,proto3" json:"tags,omitempty"`
	DateCreated *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=date_created,json=dateCreated,proto3" json:"date_created,omitempty"`
}

func (x *DescribeLinkResponse) Reset() {
	*x = DescribeLinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_link_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeLinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeLinkResponse) ProtoMessage() {}

func (x *DescribeLinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_link_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeLinkResponse.ProtoReflect.Descriptor instead.
func (*DescribeLinkResponse) Descriptor() ([]byte, []int) {
	return file_link_proto_rawDescGZIP(), []int{3}
}

func (x *DescribeLinkResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DescribeLinkResponse) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *DescribeLinkResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *DescribeLinkResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *DescribeLinkResponse) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *DescribeLinkResponse) GetDateCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.DateCreated
	}
	return nil
}

type ListLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit  *uint64 `protobuf:"varint,1,opt,name=limit,proto3,oneof" json:"limit,omitempty"`
	Offset *uint64 `protobuf:"varint,2,opt,name=offset,proto3,oneof" json:"offset,omitempty"`
}

func (x *ListLinkRequest) Reset() {
	*x = ListLinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_link_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLinkRequest) ProtoMessage() {}

func (x *ListLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_link_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLinkRequest.ProtoReflect.Descriptor instead.
func (*ListLinkRequest) Descriptor() ([]byte, []int) {
	return file_link_proto_rawDescGZIP(), []int{4}
}

func (x *ListLinkRequest) GetLimit() uint64 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *ListLinkRequest) GetOffset() uint64 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

type ListLinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*DescribeLinkResponse `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ListLinkResponse) Reset() {
	*x = ListLinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_link_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLinkResponse) ProtoMessage() {}

func (x *ListLinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_link_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLinkResponse.ProtoReflect.Descriptor instead.
func (*ListLinkResponse) Descriptor() ([]byte, []int) {
	return file_link_proto_rawDescGZIP(), []int{5}
}

func (x *ListLinkResponse) GetItems() []*DescribeLinkResponse {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_link_proto protoreflect.FileDescriptor

var file_link_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6f, 0x76,
	0x61, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x74, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61,
	0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x23,
	0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x25, 0x0a, 0x13, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4c,
	0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22, 0xc6, 0x01, 0x0a, 0x14, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x12, 0x3d, 0x0a, 0x0c, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x64, 0x61, 0x74, 0x65, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x22, 0x5e, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x1b, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x48, 0x01, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x42, 0x08,
	0x0a, 0x06, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x22, 0x4c, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x6c, 0x69, 0x6e,
	0x6b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4c, 0x69,
	0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x32, 0xc1, 0x02, 0x0a, 0x07, 0x4c, 0x69, 0x6e, 0x6b, 0x41, 0x50, 0x49, 0x12, 0x47, 0x0a,
	0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x1f, 0x2e, 0x6f, 0x76,
	0x61, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x0c, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x21, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x6c, 0x69, 0x6e,
	0x6b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x4c, 0x69,
	0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6f, 0x76, 0x61, 0x2e,
	0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4b, 0x0a, 0x08, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x1d, 0x2e, 0x6f, 0x76,
	0x61, 0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c,
	0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6f, 0x76, 0x61,
	0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69,
	0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x0a,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x1f, 0x2e, 0x6f, 0x76, 0x61,
	0x2e, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x1c, 0x5a, 0x1a, 0x2f, 0x6f, 0x76, 0x61, 0x2d, 0x6c, 0x69,
	0x6e, 0x6b, 0x2d, 0x61, 0x70, 0x69, 0x3b, 0x6f, 0x76, 0x61, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x5f,
	0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_link_proto_rawDescOnce sync.Once
	file_link_proto_rawDescData = file_link_proto_rawDesc
)

func file_link_proto_rawDescGZIP() []byte {
	file_link_proto_rawDescOnce.Do(func() {
		file_link_proto_rawDescData = protoimpl.X.CompressGZIP(file_link_proto_rawDescData)
	})
	return file_link_proto_rawDescData
}

var file_link_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_link_proto_goTypes = []interface{}{
	(*CreateLinkRequest)(nil),     // 0: ova.link.api.CreateLinkRequest
	(*DeleteLinkRequest)(nil),     // 1: ova.link.api.DeleteLinkRequest
	(*DescribeLinkRequest)(nil),   // 2: ova.link.api.DescribeLinkRequest
	(*DescribeLinkResponse)(nil),  // 3: ova.link.api.DescribeLinkResponse
	(*ListLinkRequest)(nil),       // 4: ova.link.api.ListLinkRequest
	(*ListLinkResponse)(nil),      // 5: ova.link.api.ListLinkResponse
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 7: google.protobuf.Empty
}
var file_link_proto_depIdxs = []int32{
	6, // 0: ova.link.api.DescribeLinkResponse.date_created:type_name -> google.protobuf.Timestamp
	3, // 1: ova.link.api.ListLinkResponse.items:type_name -> ova.link.api.DescribeLinkResponse
	0, // 2: ova.link.api.LinkAPI.CreateLink:input_type -> ova.link.api.CreateLinkRequest
	2, // 3: ova.link.api.LinkAPI.DescribeLink:input_type -> ova.link.api.DescribeLinkRequest
	4, // 4: ova.link.api.LinkAPI.ListLink:input_type -> ova.link.api.ListLinkRequest
	1, // 5: ova.link.api.LinkAPI.DeleteLink:input_type -> ova.link.api.DeleteLinkRequest
	7, // 6: ova.link.api.LinkAPI.CreateLink:output_type -> google.protobuf.Empty
	3, // 7: ova.link.api.LinkAPI.DescribeLink:output_type -> ova.link.api.DescribeLinkResponse
	5, // 8: ova.link.api.LinkAPI.ListLink:output_type -> ova.link.api.ListLinkResponse
	7, // 9: ova.link.api.LinkAPI.DeleteLink:output_type -> google.protobuf.Empty
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_link_proto_init() }
func file_link_proto_init() {
	if File_link_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_link_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLinkRequest); i {
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
		file_link_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteLinkRequest); i {
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
		file_link_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeLinkRequest); i {
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
		file_link_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeLinkResponse); i {
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
		file_link_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLinkRequest); i {
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
		file_link_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListLinkResponse); i {
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
	file_link_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_link_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_link_proto_goTypes,
		DependencyIndexes: file_link_proto_depIdxs,
		MessageInfos:      file_link_proto_msgTypes,
	}.Build()
	File_link_proto = out.File
	file_link_proto_rawDesc = nil
	file_link_proto_goTypes = nil
	file_link_proto_depIdxs = nil
}
