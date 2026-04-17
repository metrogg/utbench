#ifndef TENSORFLOW_CORE_TFRT_IFRT_TF_HOST_CALLBACK_H_
#define TENSORFLOW_CORE_TFRT_IFRT_TF_HOST_CALLBACK_H_
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/types/span.h"
#include "tensorflow/compiler/mlir/tfrt/transforms/ifrt/ifrt_types.h"
#include "tensorflow/core/common_runtime/device_mgr.h"
#include "tensorflow/core/common_runtime/eager/context.h"
#include "tensorflow/core/protobuf/config.pb.h"
namespace tensorflow {
namespace ifrt_serving {
class TfHostCallback {
 public:
  static absl::StatusOr<std::unique_ptr<TfHostCallback>> Create(
      absl::Span<const tensorflow::FunctionDef> functions,
      absl::string_view entry_function_name,
      absl::Span<const DtypeAndShape> operand_type_and_shapes,
      absl::Span<const DtypeAndShape> result_type_and_shapes,
      tensorflow::DeviceMgr* device_mgr);
  absl::Status Call(void** inputs, void** outputs);
 private:
  TfHostCallback(absl::string_view entry_function_name,
                 absl::Span<const DtypeAndShape> operand_type_and_shapes,
                 absl::Span<const DtypeAndShape> result_type_and_shape,
                 tensorflow::EagerContextPtr ctx)
      : ctx_(std::move(ctx)),
        entry_function_name_(entry_function_name),
        operand_type_and_shapes_(operand_type_and_shapes.begin(),
                                 operand_type_and_shapes.end()),
        result_type_and_shapes_(result_type_and_shape.begin(),
                                result_type_and_shape.end()) {}
  tensorflow::EagerContextPtr ctx_;
  std::string entry_function_name_;
  std::vector<DtypeAndShape> operand_type_and_shapes_;
  std::vector<DtypeAndShape> result_type_and_shapes_;
};
absl::StatusOr<std::unique_ptr<tensorflow::StaticDeviceMgr>>
CreateTfStaticDeviceMgr();
}  
}  
#endif  
#include "tensorflow/core/tfrt/ifrt/tf_host_callback.h"
#include <cstddef>
#include <cstring>
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/cleanup/cleanup.h"
#include "absl/container/fixed_array.h"
#include "absl/log/check.h"
#include "absl/memory/memory.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/types/span.h"
#include "tensorflow/c/eager/abstract_tensor_handle.h"
#include "tensorflow/c/eager/immediate_execution_context.h"
#include "tensorflow/c/eager/immediate_execution_operation.h"
#include "tensorflow/compiler/mlir/tfrt/transforms/ifrt/ifrt_types.h"
#include "tensorflow/core/common_runtime/device_mgr.h"
#include "tensorflow/core/common_runtime/eager/context.h"
#include "tensorflow/core/common_runtime/eager/tensor_handle.h"
#include "tensorflow/core/framework/device.h"
#include "tensorflow/core/framework/device_factory.h"
#include "tensorflow/core/framework/tensor.h"
#include "tensorflow/core/framework/types.h"
#include "tensorflow/core/protobuf/config.pb.h"
#include "tsl/platform/casts.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/refcount.h"
#include "tsl/profiler/lib/traceme.h"
namespace tensorflow {
namespace ifrt_serving {
namespace {
using RefCountHandle = ::tsl::core::RefCountPtr<tensorflow::TensorHandle>;
size_t GetSizeInBytes(const tensorflow::Tensor& tensor) {
  return tensor.shape().num_elements() * DataTypeSize(tensor.dtype());
}
tensorflow::Tensor GetTensor(const DtypeAndShape& dtype_and_shape, void* src) {
  DCHECK(DataTypeCanUseMemcpy(dtype_and_shape.dtype));
  tensorflow::Tensor t(dtype_and_shape.dtype, dtype_and_shape.shape);
  std::memcpy(t.data(), src, GetSizeInBytes(t));
  return t;
}
void CopyToBuffer(void* dst, const tensorflow::Tensor& tensor) {
  DCHECK(DataTypeCanUseMemcpy(tensor.dtype()));
  std::memcpy(dst, tensor.data(), GetSizeInBytes(tensor));
}
}  
absl::Status TfHostCallback::Call(void** inputs, void** outputs) {
  tsl::profiler::TraceMe trace_me("TfHostCallback::Call");
  tensorflow::ImmediateOpPtr op(ctx_->CreateOperation());
  TF_RETURN_IF_ERROR(
      op->Reset(entry_function_name_.c_str(), nullptr));
  ctx_->StartStep();
  absl::Cleanup cleanup_step = [this]() { ctx_->EndStep(); };
  for (int i = 0; i < operand_type_and_shapes_.size(); ++i) {
    tensorflow::Tensor t = GetTensor(operand_type_and_shapes_[i], inputs[i]);
    RefCountHandle handle(tensorflow::down_cast<tensorflow::TensorHandle*>(
        ctx_->CreateLocalHandleFromTFTensor(t, nullptr)));
    TF_RETURN_IF_ERROR(op->AddInput(handle.get()));
  }
  int num_outputs = result_type_and_shapes_.size();
  absl::FixedArray<tensorflow::AbstractTensorHandle*> output_raw_handles(
      num_outputs);
  TF_RETURN_IF_ERROR(
      op->Execute(absl::MakeSpan(output_raw_handles), &num_outputs));
  std::vector<RefCountHandle> output_handles;
  output_handles.reserve(num_outputs);
  for (auto* output_raw_handle : output_raw_handles) {
    output_handles.emplace_back(
        tensorflow::down_cast<tensorflow::TensorHandle*>(output_raw_handle));
  }
  if (result_type_and_shapes_.size() != num_outputs) {
    return absl::InternalError(absl::StrCat(
        "TF host callback invocation expected ", result_type_and_shapes_.size(),
        " results, instead got ", num_outputs));
  }
  for (int i = 0; i < num_outputs; ++i) {
    const tensorflow::Tensor* tensor;
    TF_RETURN_IF_ERROR(output_handles[i]->Tensor(&tensor));
    CopyToBuffer(outputs[i], *tensor);
  }
  return absl::OkStatus();
}
absl::StatusOr<std::unique_ptr<TfHostCallback>> TfHostCallback::Create(
    absl::Span<const tensorflow::FunctionDef> functions,
    absl::string_view entry_function_name,
    absl::Span<const DtypeAndShape> operand_type_and_shapes,
    absl::Span<const DtypeAndShape> result_type_and_shapes,
    tensorflow::DeviceMgr* device_mgr) {
  tensorflow::SessionOptions options;
  options.config.add_device_filters("/device:CPU:*");
  DCHECK(device_mgr != nullptr);
  tensorflow::EagerContextPtr ctx(new tensorflow::EagerContext(
      options,
      tensorflow::ContextDevicePlacementPolicy::DEVICE_PLACEMENT_SILENT,
      false, device_mgr,
      false,
      nullptr,
      nullptr,
      nullptr,
      true));
  for (const tensorflow::FunctionDef& function : functions) {
    TF_RETURN_IF_ERROR(ctx->AddFunctionDef(function));
  }
  return absl::WrapUnique(
      new TfHostCallback(entry_function_name, operand_type_and_shapes,
                         result_type_and_shapes, std::move(ctx)));
}
absl::StatusOr<std::unique_ptr<tensorflow::StaticDeviceMgr>>
CreateTfStaticDeviceMgr() {
  std::vector<std::unique_ptr<tensorflow::Device>> devices;
  TF_RETURN_IF_ERROR(tensorflow::DeviceFactory::AddCpuDevices(
      tensorflow::SessionOptions(), "/job:localhost/replica:0/task:0",
      &devices));
  return std::make_unique<tensorflow::StaticDeviceMgr>(std::move(devices));
}
}  
}  