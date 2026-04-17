#ifndef TENSORFLOW_COMPILER_MLIR_TENSORFLOW_TRANSFORMS_HOST_RUNTIME_TPU_METADATA_UTILS_H_
#define TENSORFLOW_COMPILER_MLIR_TENSORFLOW_TRANSFORMS_HOST_RUNTIME_TPU_METADATA_UTILS_H_
#include <optional>
#include "mlir/IR/Diagnostics.h"  
#include "mlir/Support/LogicalResult.h"  
#include "tensorflow/compiler/mlir/tensorflow/ir/tf_device.h"
#include "xla/xla.pb.h"
#include "xla/xla_data.pb.h"
#include "tensorflow/core/framework/tensor_shape.pb.h"
#include "tensorflow/core/framework/types.pb.h"
#include "tensorflow/core/protobuf/tpu/compile_metadata.pb.h"
namespace mlir {
namespace TFTPU {
LogicalResult SetMetadataProtoFromClusterFuncOp(
    tf_device::ClusterFuncOp op, int num_replicas, int num_cores_per_replica,
    std::optional<xla::DeviceAssignmentProto>&& xla_device_assignment,
    tensorflow::tpu::TPUCompileMetadataProto* metadata);
}  
}  
#endif  
#include "tensorflow/compiler/mlir/tensorflow/transforms/host_runtime/tpu_metadata_utils.h"
#include <optional>
#include <string>
#include <utility>
#include "llvm/ADT/ArrayRef.h"
#include "llvm/ADT/STLExtras.h"
#include "llvm/ADT/SmallSet.h"
#include "llvm/ADT/StringRef.h"
#include "llvm/Support/FormatVariadic.h"
#include "mlir/Dialect/Func/IR/FuncOps.h"  
#include "mlir/IR/Attributes.h"  
#include "mlir/IR/Builders.h"  
#include "mlir/IR/BuiltinAttributes.h"  
#include "mlir/IR/BuiltinTypes.h"  
#include "mlir/IR/Diagnostics.h"  
#include "mlir/IR/Operation.h"  
#include "mlir/IR/Types.h"  
#include "mlir/Support/LLVM.h"  
#include "mlir/Support/LogicalResult.h"  
#include "tensorflow/compiler/mlir/tensorflow/ir/tf_device.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/attribute_utils.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/convert_tensor.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/convert_type.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/xla_sharding_util.h"
#include "xla/xla.pb.h"
#include "xla/xla_data.pb.h"
#include "tensorflow/core/framework/tensor_shape.h"
#include "tensorflow/core/framework/tensor_shape.pb.h"
#include "tensorflow/core/framework/types.pb.h"
#include "tensorflow/core/lib/core/status.h"
#include "tensorflow/core/protobuf/tpu/compile_metadata.pb.h"
namespace mlir {
namespace TFTPU {
namespace {
constexpr char kStepMarkerLocationAttr[] = "step_marker_location";
constexpr char kUseXlaSpmdAttr[] = "use_spmd_for_xla_partitioning";
constexpr char kBadStringArrayElementMsg[] =
    "bad '{0}' attribute at index {1}, not a string";
constexpr char kBadArrayElementMsg[] =
    "bad '{0}' attribute at index {1} with value '{2}': failed to parse to {3}";
constexpr char kBadArrayAttrLengthMsg[] =
    "bad '{0}' attribute, expected array attribute of size {1}, got size {2}";
std::string CreateMissingAttributeMsg(llvm::StringRef attribute) {
  return llvm::formatv("requires attribute '{0}'", attribute).str();
}
LogicalResult SetMetadataProtoStepMarkerLocation(
    tf_device::ClusterFuncOp op,
    tensorflow::tpu::TPUCompileMetadataProto* metadata) {
  auto step_marker_location =
      op->getAttrOfType<StringAttr>(kStepMarkerLocationAttr);
  if (!step_marker_location)
    return op.emitOpError(CreateMissingAttributeMsg(kStepMarkerLocationAttr));
  xla::DebugOptions::StepMarkerLocation location =
      xla::DebugOptions::STEP_MARK_AT_ENTRY;
  if (!step_marker_location.getValue().empty() &&
      !xla::DebugOptions::StepMarkerLocation_Parse(
          std::string(step_marker_location.getValue()), &location))
    return op.emitOpError(llvm::formatv("bad '{0}' attribute with value '{1}'",
                                        kStepMarkerLocationAttr,
                                        step_marker_location.getValue()));
  metadata->set_step_marker_location(location);
  return success();
}
LogicalResult SetOpSharding(Operation* op, Attribute attr, llvm::StringRef name,
                            int index, xla::OpSharding* sharding_ptr) {
  auto sharding_attr = mlir::dyn_cast<StringAttr>(attr);
  if (!sharding_attr)
    return op->emitOpError(
        llvm::formatv(kBadStringArrayElementMsg, name, index));
  if (tensorflow::DecodeShardingAttribute(sharding_attr, *sharding_ptr)
          .failed()) {
    return op->emitOpError(llvm::formatv(kBadArrayElementMsg, name, index,
                                         sharding_attr.getValue(),
                                         "xla::OpSharding"));
  }
  return success();
}
LogicalResult SetMetadataProtoArgs(
    tf_device::ClusterFuncOp op,
    tensorflow::tpu::TPUCompileMetadataProto* metadata) {
  auto input_shardings =
      op->getAttrOfType<ArrayAttr>(tensorflow::kInputShardingAttr);
  if (!input_shardings)
    return op.emitOpError(
        CreateMissingAttributeMsg(tensorflow::kInputShardingAttr));
  if (input_shardings.size() != op.getNumOperands())
    return op.emitOpError(
        llvm::formatv(kBadArrayAttrLengthMsg, tensorflow::kInputShardingAttr,
                      op.getNumOperands(), input_shardings.size()));
  mlir::StringAttr replication_attr_name = mlir::StringAttr::get(
      op.getContext(), "mhlo.is_same_data_across_replicas");
  auto dynamic_arg_idx = op->getAttrOfType<ArrayAttr>(TF::kDynamicArgIndexAttr);
  llvm::SmallSet<int, 4> dynamic_arg_idx_set;
  if (dynamic_arg_idx) {
    for (auto idx : dynamic_arg_idx.getValue()) {
      dynamic_arg_idx_set.insert(mlir::dyn_cast<IntegerAttr>(idx).getInt());
    }
  }
  for (auto operand_type_and_idx : llvm::enumerate(op.getOperandTypes())) {
    Type operand_type = operand_type_and_idx.value();
    int index = operand_type_and_idx.index();
    tensorflow::tpu::TPUCompileMetadataProto::Arg* arg = metadata->add_args();
    tensorflow::DataType dtype;
    tensorflow::Status status =
        tensorflow::ConvertToDataType(operand_type, &dtype);
    if (!status.ok())
      return op.emitOpError(
          llvm::formatv("failed to determine operand type at index {0}: {1}",
                        index, status.message()));
    arg->set_dtype(dtype);
    if (dtype == tensorflow::DT_RESOURCE)
      arg->set_kind(tensorflow::tpu::TPUCompileMetadataProto::Arg::VARIABLE);
    else
      arg->set_kind(tensorflow::tpu::TPUCompileMetadataProto::Arg::PARAMETER);
    *arg->mutable_shape() = tensorflow::TensorShapeProto();
    if (auto ranked_tensor_type =
            mlir::dyn_cast<RankedTensorType>(operand_type)) {
      tensorflow::TensorShapeProto shape_proto;
      ConvertToTensorShapeProto(ranked_tensor_type.getShape(), &shape_proto);
      *arg->mutable_shape() = std::move(shape_proto);
    } else {
      arg->mutable_shape()->set_unknown_rank(true);
    }
    if (failed(SetOpSharding(op, input_shardings.getValue()[index],
                             tensorflow::kInputShardingAttr, index,
                             arg->mutable_sharding())))
      return failure();
    auto attr = op.getFuncOp().getArgAttrOfType<mlir::BoolAttr>(
        index, replication_attr_name);
    arg->set_is_same_data_across_replicas(attr != nullptr && attr.getValue());
    arg->mutable_is_bounded_dynamic_dim()->Add(
        dynamic_arg_idx_set.contains(index));
  }
  return success();
}
LogicalResult SetMetadataProtoRetvals(
    tf_device::ClusterFuncOp op,
    tensorflow::tpu::TPUCompileMetadataProto* metadata) {
  auto output_shardings =
      op->getAttrOfType<ArrayAttr>(tensorflow::kOutputShardingAttr);
  if (!output_shardings)
    return op.emitOpError(
        CreateMissingAttributeMsg(tensorflow::kOutputShardingAttr));
  if (output_shardings.size() != op.getNumResults())
    return op.emitOpError(
        llvm::formatv(kBadArrayAttrLengthMsg, tensorflow::kOutputShardingAttr,
                      op.getNumResults(), output_shardings.size()));
  for (auto output_sharding_and_idx : llvm::enumerate(output_shardings))
    if (failed(SetOpSharding(op, output_sharding_and_idx.value(),
                             tensorflow::kOutputShardingAttr,
                             output_sharding_and_idx.index(),
                             metadata->add_retvals()->mutable_sharding())))
      return failure();
  return success();
}
}  
LogicalResult SetMetadataProtoFromClusterFuncOp(
    tf_device::ClusterFuncOp op, int num_replicas, int num_cores_per_replica,
    std::optional<xla::DeviceAssignmentProto>&& xla_device_assignment,
    tensorflow::tpu::TPUCompileMetadataProto* metadata) {
  if (auto options_attr =
          op->getAttrOfType<StringAttr>("tpu_compile_options_proto")) {
    if (!metadata->mutable_compile_options()->ParseFromArray(
            options_attr.data(), options_attr.size())) {
      return failure();
    }
  }
  metadata->set_num_replicas(num_replicas);
  metadata->set_num_cores_per_replica(num_cores_per_replica);
  if (failed(SetMetadataProtoStepMarkerLocation(op, metadata)))
    return failure();
  if (xla_device_assignment.has_value())
    *metadata->mutable_device_assignment() =
        std::move(xla_device_assignment.value());
  auto use_spmd_attr = op->getAttrOfType<BoolAttr>(kUseXlaSpmdAttr);
  if (!use_spmd_attr)
    return op.emitOpError(CreateMissingAttributeMsg(kUseXlaSpmdAttr));
  metadata->set_use_spmd_for_xla_partitioning(use_spmd_attr.getValue());
  if (failed(SetMetadataProtoArgs(op, metadata))) return failure();
  return SetMetadataProtoRetvals(op, metadata);
}
}  
}  