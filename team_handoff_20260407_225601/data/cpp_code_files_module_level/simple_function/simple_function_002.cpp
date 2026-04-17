#ifndef MLIR_HLO_UTILS_PLACEMENT_UTILS_H
#define MLIR_HLO_UTILS_PLACEMENT_UTILS_H
#include "llvm/ADT/StringRef.h"
namespace mlir {
namespace mhlo {
namespace placement_utils {
constexpr llvm::StringRef cCpu = "cpu";
constexpr llvm::StringRef cGpu = "gpu";
}  
}  
}  
#endif  
#include "tensorflow/core/common_runtime/eager/placement_utils.h"
#include <variant>
#include "absl/status/status.h"
#include "tensorflow/c/eager/immediate_execution_tensor_handle.h"
#include "tensorflow/core/common_runtime/eager/attr_builder.h"
#include "tensorflow/core/common_runtime/eager/custom_device.h"
#include "tensorflow/core/common_runtime/eager/eager_operation.h"
#include "tensorflow/core/common_runtime/input_colocation_exemption_registry.h"
#include "tensorflow/core/framework/op_def.pb.h"
#include "tensorflow/core/framework/types.pb.h"
#include "tensorflow/core/platform/errors.h"
namespace tensorflow {
namespace eager {
static bool IsPinnableOp(StringPiece op_name) {
  static const gtl::FlatSet<string>* unpinnable_ops = new gtl::FlatSet<string>({
      "RandomUniform",
      "RandomUniformInt",
      "RandomStandardNormal",
      "StatelessRandomUniform",
      "StatelessRandomUniformInt",
      "StatelessRandomUniformFullInt",
      "StatelessRandomNormal",
  });
  return unpinnable_ops->find(string(op_name)) == unpinnable_ops->end() &&
         !absl::StartsWith(op_name, "XRT");
}
static Status ValidateTensorHandleRemoteDevice(EagerContext* ctx,
                                               int64_t device_incarnation) {
  if (ctx->remote_device_mgr()->ContainsDevice(device_incarnation)) {
    return absl::OkStatus();
  }
  return errors::InvalidArgument(
      "Resource input tensor contains an invalid device. This might happen "
      "when the client has connected to a different cluster, or some remote "
      "workers have been restarted.");
}
bool IsColocationExempt(StringPiece op_name) {
  const auto& exempt_ops = InputColocationExemptionRegistry::Global()->Get();
  return exempt_ops.find(string(op_name)) != exempt_ops.end();
}
bool IsFunction(StringPiece op_name) {
  const OpDef* op_def = nullptr;
  Status s = OpDefForOp(string(op_name), &op_def);
  if (!s.ok()) {
    if (!absl::IsNotFound(s)) {
      LOG(WARNING) << "Looking up OpDef failed with error: " << s;
    }
    return true;
  }
  return false;
}
Status MaybePinSmallOpsToCpu(
    bool* result, StringPiece op_name,
    absl::Span<ImmediateExecutionTensorHandle* const> args,
    StringPiece cpu_device_name) {
  if (IsFunction(op_name) || IsColocationExempt(op_name) ||
      !IsPinnableOp(op_name)) {
    *result = false;
    return absl::OkStatus();
  }
  if (args.empty()) {
    *result = false;
    return absl::OkStatus();
  }
  int i = 0;
  for (auto* arg : args) {
    Status s;
    const char* device_name = arg->DeviceName(&s);
    DataType dtype = arg->DataType();
    TF_RETURN_IF_ERROR(s);
    DVLOG(2) << "for op " << op_name << " input " << i << " "
             << DataTypeString(dtype) << " input device = " << device_name;
    if (device_name != cpu_device_name) {
      *result = false;
      return absl::OkStatus();
    }
    if (dtype != DataType::DT_INT32 && dtype != DataType::DT_INT64) {
      *result = false;
      return absl::OkStatus();
    }
    int64_t num_elements;
    TF_RETURN_IF_ERROR(arg->NumElements(&num_elements));
    if (num_elements > 64) {
      *result = false;
      return absl::OkStatus();
    }
    i++;
  }
  DVLOG(1) << "Forcing op " << op_name
           << " to be on the CPU since all input tensors have an "
              "int32/int64 dtype, and are small (less than 64 elements).";
  *result = true;
  return absl::OkStatus();
}
Status MaybePinToResourceDevice(Device** device, const EagerOperation& op) {
  if (op.colocation_exempt()) {
    return absl::OkStatus();
  }
  EagerContext& ctx = op.EagerContext();
  const absl::InlinedVector<TensorHandle*, 4>* inputs;
  TF_RETURN_IF_ERROR(op.TensorHandleInputs(&inputs));
  Device* op_device = op.Device() == kVariantDeviceNull
                          ? ctx.HostCPU()
                          : std::get<Device*>(op.Device());
  for (int i = 0; i < inputs->size(); ++i) {
    TensorHandle* tensor_handle = (*inputs)[i];
    if (tensor_handle->dtype == DT_RESOURCE) {
      if (tensor_handle->resource_remote_device_incarnation() != 0) {
        TF_RETURN_IF_ERROR(ValidateTensorHandleRemoteDevice(
            &ctx, tensor_handle->resource_remote_device_incarnation()));
      }
      Device* resource_device = tensor_handle->resource_device();
      DVLOG(2) << "for op " << op.Name() << " input " << i << " "
               << DataTypeString(tensor_handle->dtype)
               << " input device = " << resource_device->name()
               << ", op device = " << op_device->name();
      if (resource_device != op_device || op.Device() == kVariantDeviceNull) {
        DVLOG(1) << (resource_device != op_device ? "Changing " : "Setting ")
                 << "device of operation " << op.Name() << " to "
                 << resource_device->name() << " because input #" << i
                 << " is a resource in this device.";
        *device = resource_device;
        return absl::OkStatus();
      }
    }
  }
  return absl::OkStatus();
}
}  
}  