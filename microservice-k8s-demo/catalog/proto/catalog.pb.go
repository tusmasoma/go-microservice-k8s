// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: proto/catalog.proto

package proto

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

type ListCatalogItemsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListCatalogItemsRequest) Reset() {
	*x = ListCatalogItemsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCatalogItemsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCatalogItemsRequest) ProtoMessage() {}

func (x *ListCatalogItemsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCatalogItemsRequest.ProtoReflect.Descriptor instead.
func (*ListCatalogItemsRequest) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{0}
}

type ListCatalogItemsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*CatalogItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ListCatalogItemsResponse) Reset() {
	*x = ListCatalogItemsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCatalogItemsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCatalogItemsResponse) ProtoMessage() {}

func (x *ListCatalogItemsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCatalogItemsResponse.ProtoReflect.Descriptor instead.
func (*ListCatalogItemsResponse) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{1}
}

func (x *ListCatalogItemsResponse) GetItems() []*CatalogItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type ListCatalogItemsByNameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ListCatalogItemsByNameRequest) Reset() {
	*x = ListCatalogItemsByNameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCatalogItemsByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCatalogItemsByNameRequest) ProtoMessage() {}

func (x *ListCatalogItemsByNameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCatalogItemsByNameRequest.ProtoReflect.Descriptor instead.
func (*ListCatalogItemsByNameRequest) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{2}
}

func (x *ListCatalogItemsByNameRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ListCatalogItemsByNameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*CatalogItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ListCatalogItemsByNameResponse) Reset() {
	*x = ListCatalogItemsByNameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCatalogItemsByNameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCatalogItemsByNameResponse) ProtoMessage() {}

func (x *ListCatalogItemsByNameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCatalogItemsByNameResponse.ProtoReflect.Descriptor instead.
func (*ListCatalogItemsByNameResponse) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{3}
}

func (x *ListCatalogItemsByNameResponse) GetItems() []*CatalogItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type CatalogItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Price float32 `protobuf:"fixed32,3,opt,name=price,proto3" json:"price,omitempty"`
}

func (x *CatalogItem) Reset() {
	*x = CatalogItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CatalogItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatalogItem) ProtoMessage() {}

func (x *CatalogItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatalogItem.ProtoReflect.Descriptor instead.
func (*CatalogItem) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{4}
}

func (x *CatalogItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CatalogItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CatalogItem) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

type CreateCatalogItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Price float32 `protobuf:"fixed32,2,opt,name=price,proto3" json:"price,omitempty"`
}

func (x *CreateCatalogItemRequest) Reset() {
	*x = CreateCatalogItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCatalogItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCatalogItemRequest) ProtoMessage() {}

func (x *CreateCatalogItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCatalogItemRequest.ProtoReflect.Descriptor instead.
func (*CreateCatalogItemRequest) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{5}
}

func (x *CreateCatalogItemRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateCatalogItemRequest) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

type CreateCatalogItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateCatalogItemResponse) Reset() {
	*x = CreateCatalogItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCatalogItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCatalogItemResponse) ProtoMessage() {}

func (x *CreateCatalogItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCatalogItemResponse.ProtoReflect.Descriptor instead.
func (*CreateCatalogItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{6}
}

type UpdateCatalogItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name  string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Price float32 `protobuf:"fixed32,3,opt,name=price,proto3" json:"price,omitempty"`
}

func (x *UpdateCatalogItemRequest) Reset() {
	*x = UpdateCatalogItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCatalogItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCatalogItemRequest) ProtoMessage() {}

func (x *UpdateCatalogItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCatalogItemRequest.ProtoReflect.Descriptor instead.
func (*UpdateCatalogItemRequest) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateCatalogItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateCatalogItemRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateCatalogItemRequest) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

type UpdateCatalogItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateCatalogItemResponse) Reset() {
	*x = UpdateCatalogItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCatalogItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCatalogItemResponse) ProtoMessage() {}

func (x *UpdateCatalogItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCatalogItemResponse.ProtoReflect.Descriptor instead.
func (*UpdateCatalogItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{8}
}

type DeleteCatalogItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteCatalogItemRequest) Reset() {
	*x = DeleteCatalogItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCatalogItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCatalogItemRequest) ProtoMessage() {}

func (x *DeleteCatalogItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCatalogItemRequest.ProtoReflect.Descriptor instead.
func (*DeleteCatalogItemRequest) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteCatalogItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteCatalogItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteCatalogItemResponse) Reset() {
	*x = DeleteCatalogItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_catalog_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCatalogItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCatalogItemResponse) ProtoMessage() {}

func (x *DeleteCatalogItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_catalog_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCatalogItemResponse.ProtoReflect.Descriptor instead.
func (*DeleteCatalogItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_catalog_proto_rawDescGZIP(), []int{10}
}

var File_proto_catalog_proto protoreflect.FileDescriptor

var file_proto_catalog_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x22, 0x19,
	0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x46, 0x0a, 0x18, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x43,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x33, 0x0a, 0x1d, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x73, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4c, 0x0a, 0x1e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f,
	0x67, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69,
	0x74, 0x65, 0x6d, 0x73, 0x22, 0x47, 0x0a, 0x0b, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0x44, 0x0a,
	0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x22, 0x1b, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x54, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f,
	0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0x1b, 0x0a, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x2a, 0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x1b, 0x0a, 0x19, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xe8, 0x03, 0x0a,
	0x0e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x57, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74,
	0x65, 0x6d, 0x73, 0x12, 0x20, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x69, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x42, 0x79, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x26, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x42, 0x79, 0x4e,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x63, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x73, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x21, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c,
	0x6f, 0x67, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61,
	0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x5a, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x12, 0x21, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f,
	0x67, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49,
	0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x11, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d,
	0x12, 0x21, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x08, 0x5a, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_catalog_proto_rawDescOnce sync.Once
	file_proto_catalog_proto_rawDescData = file_proto_catalog_proto_rawDesc
)

func file_proto_catalog_proto_rawDescGZIP() []byte {
	file_proto_catalog_proto_rawDescOnce.Do(func() {
		file_proto_catalog_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_catalog_proto_rawDescData)
	})
	return file_proto_catalog_proto_rawDescData
}

var file_proto_catalog_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_catalog_proto_goTypes = []interface{}{
	(*ListCatalogItemsRequest)(nil),        // 0: catalog.ListCatalogItemsRequest
	(*ListCatalogItemsResponse)(nil),       // 1: catalog.ListCatalogItemsResponse
	(*ListCatalogItemsByNameRequest)(nil),  // 2: catalog.ListCatalogItemsByNameRequest
	(*ListCatalogItemsByNameResponse)(nil), // 3: catalog.ListCatalogItemsByNameResponse
	(*CatalogItem)(nil),                    // 4: catalog.CatalogItem
	(*CreateCatalogItemRequest)(nil),       // 5: catalog.CreateCatalogItemRequest
	(*CreateCatalogItemResponse)(nil),      // 6: catalog.CreateCatalogItemResponse
	(*UpdateCatalogItemRequest)(nil),       // 7: catalog.UpdateCatalogItemRequest
	(*UpdateCatalogItemResponse)(nil),      // 8: catalog.UpdateCatalogItemResponse
	(*DeleteCatalogItemRequest)(nil),       // 9: catalog.DeleteCatalogItemRequest
	(*DeleteCatalogItemResponse)(nil),      // 10: catalog.DeleteCatalogItemResponse
}
var file_proto_catalog_proto_depIdxs = []int32{
	4,  // 0: catalog.ListCatalogItemsResponse.items:type_name -> catalog.CatalogItem
	4,  // 1: catalog.ListCatalogItemsByNameResponse.items:type_name -> catalog.CatalogItem
	0,  // 2: catalog.CatalogService.ListCatalogItems:input_type -> catalog.ListCatalogItemsRequest
	2,  // 3: catalog.CatalogService.ListCatalogItemsByName:input_type -> catalog.ListCatalogItemsByNameRequest
	5,  // 4: catalog.CatalogService.CreateCatalogItem:input_type -> catalog.CreateCatalogItemRequest
	7,  // 5: catalog.CatalogService.UpdateCatalogItem:input_type -> catalog.UpdateCatalogItemRequest
	9,  // 6: catalog.CatalogService.DeleteCatalogItem:input_type -> catalog.DeleteCatalogItemRequest
	1,  // 7: catalog.CatalogService.ListCatalogItems:output_type -> catalog.ListCatalogItemsResponse
	3,  // 8: catalog.CatalogService.ListCatalogItemsByName:output_type -> catalog.ListCatalogItemsByNameResponse
	6,  // 9: catalog.CatalogService.CreateCatalogItem:output_type -> catalog.CreateCatalogItemResponse
	8,  // 10: catalog.CatalogService.UpdateCatalogItem:output_type -> catalog.UpdateCatalogItemResponse
	10, // 11: catalog.CatalogService.DeleteCatalogItem:output_type -> catalog.DeleteCatalogItemResponse
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_proto_catalog_proto_init() }
func file_proto_catalog_proto_init() {
	if File_proto_catalog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_catalog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCatalogItemsRequest); i {
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
		file_proto_catalog_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCatalogItemsResponse); i {
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
		file_proto_catalog_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCatalogItemsByNameRequest); i {
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
		file_proto_catalog_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListCatalogItemsByNameResponse); i {
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
		file_proto_catalog_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CatalogItem); i {
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
		file_proto_catalog_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCatalogItemRequest); i {
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
		file_proto_catalog_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCatalogItemResponse); i {
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
		file_proto_catalog_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateCatalogItemRequest); i {
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
		file_proto_catalog_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateCatalogItemResponse); i {
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
		file_proto_catalog_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCatalogItemRequest); i {
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
		file_proto_catalog_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCatalogItemResponse); i {
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
			RawDescriptor: file_proto_catalog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_catalog_proto_goTypes,
		DependencyIndexes: file_proto_catalog_proto_depIdxs,
		MessageInfos:      file_proto_catalog_proto_msgTypes,
	}.Build()
	File_proto_catalog_proto = out.File
	file_proto_catalog_proto_rawDesc = nil
	file_proto_catalog_proto_goTypes = nil
	file_proto_catalog_proto_depIdxs = nil
}
