// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: gpu.proto

package gpu

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

type CheckpointRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Directory string `protobuf:"bytes,1,opt,name=directory,proto3" json:"directory,omitempty"`
}

func (x *CheckpointRequest) Reset() {
	*x = CheckpointRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckpointRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckpointRequest) ProtoMessage() {}

func (x *CheckpointRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckpointRequest.ProtoReflect.Descriptor instead.
func (*CheckpointRequest) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{0}
}

func (x *CheckpointRequest) GetDirectory() string {
	if x != nil {
		return x.Directory
	}
	return ""
}

type CheckpointResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success  bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	MemPath  string `protobuf:"bytes,2,opt,name=memPath,proto3" json:"memPath,omitempty"`
	CkptPath string `protobuf:"bytes,3,opt,name=ckptPath,proto3" json:"ckptPath,omitempty"`
}

func (x *CheckpointResponse) Reset() {
	*x = CheckpointResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckpointResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckpointResponse) ProtoMessage() {}

func (x *CheckpointResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckpointResponse.ProtoReflect.Descriptor instead.
func (*CheckpointResponse) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{1}
}

func (x *CheckpointResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CheckpointResponse) GetMemPath() string {
	if x != nil {
		return x.MemPath
	}
	return ""
}

func (x *CheckpointResponse) GetCkptPath() string {
	if x != nil {
		return x.CkptPath
	}
	return ""
}

type RestoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Directory string `protobuf:"bytes,1,opt,name=directory,proto3" json:"directory,omitempty"`
}

func (x *RestoreRequest) Reset() {
	*x = RestoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreRequest) ProtoMessage() {}

func (x *RestoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreRequest.ProtoReflect.Descriptor instead.
func (*RestoreRequest) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{2}
}

func (x *RestoreRequest) GetDirectory() string {
	if x != nil {
		return x.Directory
	}
	return ""
}

type RestoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success         bool             `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	GpuRestoreStats *GPURestoreStats `protobuf:"bytes,2,opt,name=gpuRestoreStats,proto3" json:"gpuRestoreStats,omitempty"`
}

func (x *RestoreResponse) Reset() {
	*x = RestoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreResponse) ProtoMessage() {}

func (x *RestoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreResponse.ProtoReflect.Descriptor instead.
func (*RestoreResponse) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{3}
}

func (x *RestoreResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *RestoreResponse) GetGpuRestoreStats() *GPURestoreStats {
	if x != nil {
		return x.GpuRestoreStats
	}
	return nil
}

type StartupPollRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StartupPollRequest) Reset() {
	*x = StartupPollRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartupPollRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartupPollRequest) ProtoMessage() {}

func (x *StartupPollRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartupPollRequest.ProtoReflect.Descriptor instead.
func (*StartupPollRequest) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{4}
}

type StartupPollResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *StartupPollResponse) Reset() {
	*x = StartupPollResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartupPollResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartupPollResponse) ProtoMessage() {}

func (x *StartupPollResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartupPollResponse.ProtoReflect.Descriptor instead.
func (*StartupPollResponse) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{5}
}

func (x *StartupPollResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type HealthCheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthCheckRequest) Reset() {
	*x = HealthCheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckRequest) ProtoMessage() {}

func (x *HealthCheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckRequest.ProtoReflect.Descriptor instead.
func (*HealthCheckRequest) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{6}
}

type HealthCheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success           bool               `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Version           string             `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	DeviceCount       int32              `protobuf:"varint,3,opt,name=deviceCount,proto3" json:"deviceCount,omitempty"`
	AvailableCUDAAPIs *AvailableCUDAAPIs `protobuf:"bytes,4,opt,name=availableCUDAAPIs,proto3" json:"availableCUDAAPIs,omitempty"`
}

func (x *HealthCheckResponse) Reset() {
	*x = HealthCheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckResponse) ProtoMessage() {}

func (x *HealthCheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckResponse.ProtoReflect.Descriptor instead.
func (*HealthCheckResponse) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{7}
}

func (x *HealthCheckResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *HealthCheckResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *HealthCheckResponse) GetDeviceCount() int32 {
	if x != nil {
		return x.DeviceCount
	}
	return 0
}

func (x *HealthCheckResponse) GetAvailableCUDAAPIs() *AvailableCUDAAPIs {
	if x != nil {
		return x.AvailableCUDAAPIs
	}
	return nil
}

type AvailableCUDAAPIs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CuDNNVersion   int32 `protobuf:"varint,1,opt,name=cuDNNVersion,proto3" json:"cuDNNVersion,omitempty"`
	CuBLASVersion  int32 `protobuf:"varint,2,opt,name=cuBLASVersion,proto3" json:"cuBLASVersion,omitempty"`
	NcclVersion    int32 `protobuf:"varint,3,opt,name=ncclVersion,proto3" json:"ncclVersion,omitempty"`
	DriverVersion  int32 `protobuf:"varint,4,opt,name=driverVersion,proto3" json:"driverVersion,omitempty"`
	RuntimeVersion int32 `protobuf:"varint,5,opt,name=runtimeVersion,proto3" json:"runtimeVersion,omitempty"`
}

func (x *AvailableCUDAAPIs) Reset() {
	*x = AvailableCUDAAPIs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AvailableCUDAAPIs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AvailableCUDAAPIs) ProtoMessage() {}

func (x *AvailableCUDAAPIs) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AvailableCUDAAPIs.ProtoReflect.Descriptor instead.
func (*AvailableCUDAAPIs) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{8}
}

func (x *AvailableCUDAAPIs) GetCuDNNVersion() int32 {
	if x != nil {
		return x.CuDNNVersion
	}
	return 0
}

func (x *AvailableCUDAAPIs) GetCuBLASVersion() int32 {
	if x != nil {
		return x.CuBLASVersion
	}
	return 0
}

func (x *AvailableCUDAAPIs) GetNcclVersion() int32 {
	if x != nil {
		return x.NcclVersion
	}
	return 0
}

func (x *AvailableCUDAAPIs) GetDriverVersion() int32 {
	if x != nil {
		return x.DriverVersion
	}
	return 0
}

func (x *AvailableCUDAAPIs) GetRuntimeVersion() int32 {
	if x != nil {
		return x.RuntimeVersion
	}
	return 0
}

type GPURestoreStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CopyMemTime     int64 `protobuf:"varint,1,opt,name=copyMemTime,proto3" json:"copyMemTime,omitempty"`
	ReplayCallsTime int64 `protobuf:"varint,2,opt,name=replayCallsTime,proto3" json:"replayCallsTime,omitempty"`
}

func (x *GPURestoreStats) Reset() {
	*x = GPURestoreStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gpu_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GPURestoreStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GPURestoreStats) ProtoMessage() {}

func (x *GPURestoreStats) ProtoReflect() protoreflect.Message {
	mi := &file_gpu_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GPURestoreStats.ProtoReflect.Descriptor instead.
func (*GPURestoreStats) Descriptor() ([]byte, []int) {
	return file_gpu_proto_rawDescGZIP(), []int{9}
}

func (x *GPURestoreStats) GetCopyMemTime() int64 {
	if x != nil {
		return x.CopyMemTime
	}
	return 0
}

func (x *GPURestoreStats) GetReplayCallsTime() int64 {
	if x != nil {
		return x.ReplayCallsTime
	}
	return 0
}

var File_gpu_proto protoreflect.FileDescriptor

var file_gpu_proto_rawDesc = []byte{
	0x0a, 0x09, 0x67, 0x70, 0x75, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x65, 0x64,
	0x61, 0x6e, 0x61, 0x67, 0x70, 0x75, 0x22, 0x31, 0x0a, 0x11, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x22, 0x64, 0x0a, 0x12, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x6d,
	0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x50,
	0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6b, 0x70, 0x74, 0x50, 0x61, 0x74, 0x68, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6b, 0x70, 0x74, 0x50, 0x61, 0x74, 0x68, 0x22,
	0x2e, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x22,
	0x71, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x44, 0x0a, 0x0f,
	0x67, 0x70, 0x75, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67, 0x70,
	0x75, 0x2e, 0x47, 0x50, 0x55, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x0f, 0x67, 0x70, 0x75, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x22, 0x14, 0x0a, 0x12, 0x53, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x50, 0x6f, 0x6c,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x2f, 0x0a, 0x13, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x75, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x14, 0x0a, 0x12, 0x48, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0xb7, 0x01, 0x0a, 0x13, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x4a, 0x0a,
	0x11, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x55, 0x44, 0x41, 0x41, 0x50,
	0x49, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e,
	0x61, 0x67, 0x70, 0x75, 0x2e, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x55,
	0x44, 0x41, 0x41, 0x50, 0x49, 0x73, 0x52, 0x11, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x43, 0x55, 0x44, 0x41, 0x41, 0x50, 0x49, 0x73, 0x22, 0xcd, 0x01, 0x0a, 0x11, 0x41, 0x76,
	0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x55, 0x44, 0x41, 0x41, 0x50, 0x49, 0x73, 0x12,
	0x22, 0x0a, 0x0c, 0x63, 0x75, 0x44, 0x4e, 0x4e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x63, 0x75, 0x44, 0x4e, 0x4e, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x75, 0x42, 0x4c, 0x41, 0x53, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x63, 0x75, 0x42, 0x4c,
	0x41, 0x53, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x63, 0x63,
	0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b,
	0x6e, 0x63, 0x63, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x64,
	0x72, 0x69, 0x76, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0d, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x26, 0x0a, 0x0e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x72, 0x75, 0x6e, 0x74, 0x69,
	0x6d, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x5d, 0x0a, 0x0f, 0x47, 0x50, 0x55,
	0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x20, 0x0a, 0x0b,
	0x63, 0x6f, 0x70, 0x79, 0x4d, 0x65, 0x6d, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0b, 0x63, 0x6f, 0x70, 0x79, 0x4d, 0x65, 0x6d, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x28,
	0x0a, 0x0f, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x43, 0x61, 0x6c, 0x6c, 0x73, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x43,
	0x61, 0x6c, 0x6c, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x32, 0xbc, 0x02, 0x0a, 0x09, 0x43, 0x65, 0x64,
	0x61, 0x6e, 0x61, 0x47, 0x50, 0x55, 0x12, 0x4b, 0x0a, 0x0a, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67, 0x70, 0x75,
	0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67, 0x70, 0x75, 0x2e, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x19,
	0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67, 0x70, 0x75, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x63, 0x65, 0x64, 0x61,
	0x6e, 0x61, 0x67, 0x70, 0x75, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x75, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x12, 0x1d, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67,
	0x70, 0x75, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67, 0x70,
	0x75, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x50, 0x6f, 0x6c, 0x6c, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4e, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x1d, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67,
	0x70, 0x75, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x67, 0x70,
	0x75, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x65, 0x64, 0x61, 0x6e, 0x61, 0x2f, 0x63, 0x65, 0x64,
	0x61, 0x6e, 0x61, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2f, 0x67, 0x70, 0x75, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gpu_proto_rawDescOnce sync.Once
	file_gpu_proto_rawDescData = file_gpu_proto_rawDesc
)

func file_gpu_proto_rawDescGZIP() []byte {
	file_gpu_proto_rawDescOnce.Do(func() {
		file_gpu_proto_rawDescData = protoimpl.X.CompressGZIP(file_gpu_proto_rawDescData)
	})
	return file_gpu_proto_rawDescData
}

var file_gpu_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_gpu_proto_goTypes = []interface{}{
	(*CheckpointRequest)(nil),   // 0: cedanagpu.CheckpointRequest
	(*CheckpointResponse)(nil),  // 1: cedanagpu.CheckpointResponse
	(*RestoreRequest)(nil),      // 2: cedanagpu.RestoreRequest
	(*RestoreResponse)(nil),     // 3: cedanagpu.RestoreResponse
	(*StartupPollRequest)(nil),  // 4: cedanagpu.StartupPollRequest
	(*StartupPollResponse)(nil), // 5: cedanagpu.StartupPollResponse
	(*HealthCheckRequest)(nil),  // 6: cedanagpu.HealthCheckRequest
	(*HealthCheckResponse)(nil), // 7: cedanagpu.HealthCheckResponse
	(*AvailableCUDAAPIs)(nil),   // 8: cedanagpu.AvailableCUDAAPIs
	(*GPURestoreStats)(nil),     // 9: cedanagpu.GPURestoreStats
}
var file_gpu_proto_depIdxs = []int32{
	9, // 0: cedanagpu.RestoreResponse.gpuRestoreStats:type_name -> cedanagpu.GPURestoreStats
	8, // 1: cedanagpu.HealthCheckResponse.availableCUDAAPIs:type_name -> cedanagpu.AvailableCUDAAPIs
	0, // 2: cedanagpu.CedanaGPU.Checkpoint:input_type -> cedanagpu.CheckpointRequest
	2, // 3: cedanagpu.CedanaGPU.Restore:input_type -> cedanagpu.RestoreRequest
	4, // 4: cedanagpu.CedanaGPU.StartupPoll:input_type -> cedanagpu.StartupPollRequest
	6, // 5: cedanagpu.CedanaGPU.HealthCheck:input_type -> cedanagpu.HealthCheckRequest
	1, // 6: cedanagpu.CedanaGPU.Checkpoint:output_type -> cedanagpu.CheckpointResponse
	3, // 7: cedanagpu.CedanaGPU.Restore:output_type -> cedanagpu.RestoreResponse
	5, // 8: cedanagpu.CedanaGPU.StartupPoll:output_type -> cedanagpu.StartupPollResponse
	7, // 9: cedanagpu.CedanaGPU.HealthCheck:output_type -> cedanagpu.HealthCheckResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_gpu_proto_init() }
func file_gpu_proto_init() {
	if File_gpu_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gpu_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckpointRequest); i {
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
		file_gpu_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckpointResponse); i {
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
		file_gpu_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreRequest); i {
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
		file_gpu_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreResponse); i {
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
		file_gpu_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartupPollRequest); i {
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
		file_gpu_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartupPollResponse); i {
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
		file_gpu_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheckRequest); i {
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
		file_gpu_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheckResponse); i {
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
		file_gpu_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AvailableCUDAAPIs); i {
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
		file_gpu_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GPURestoreStats); i {
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
			RawDescriptor: file_gpu_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gpu_proto_goTypes,
		DependencyIndexes: file_gpu_proto_depIdxs,
		MessageInfos:      file_gpu_proto_msgTypes,
	}.Build()
	File_gpu_proto = out.File
	file_gpu_proto_rawDesc = nil
	file_gpu_proto_goTypes = nil
	file_gpu_proto_depIdxs = nil
}
