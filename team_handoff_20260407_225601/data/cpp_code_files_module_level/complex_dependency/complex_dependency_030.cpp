#ifndef TENSORFLOW_COMPILER_MLIR_TENSORFLOW_TRANSFORMS_HOST_RUNTIME_LOWER_CLUSTER_TO_RUNTIME_OPS_H_
#define TENSORFLOW_COMPILER_MLIR_TENSORFLOW_TRANSFORMS_HOST_RUNTIME_LOWER_CLUSTER_TO_RUNTIME_OPS_H_
#include "absl/base/attributes.h"
#include "llvm/ADT/StringRef.h"
#include "mlir/IR/BuiltinOps.h"  
#include "mlir/Pass/PassManager.h"  
#include "xla/tsl/framework/device_type.h"
#include "tensorflow/core/lib/core/status.h"
namespace tensorflow {
namespace tfrt_compiler {
tensorflow::Status RunLowerClusterToRuntimeOpsPassPipeline(
    mlir::ModuleOp module, tsl::DeviceType xla_device_type,
    llvm::StringRef module_name = llvm::StringRef());
void RegisterTPULowerClusterToRuntimeOpsPassPipeline();
void RegisterNonTPULowerClusterToRuntimeOpsPassPipeline();
}  
}  
#endif  
#include <memory>
#include <string>
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "llvm/ADT/StringRef.h"
#include "mlir/Dialect/Func/IR/FuncOps.h"  
#include "mlir/IR/BuiltinOps.h"  
#include "mlir/Pass/PassManager.h"  
#include "mlir/Pass/PassRegistry.h"  
#include "mlir/Support/LogicalResult.h"  
#include "mlir/Transforms/Passes.h"  
#include "tensorflow/compiler/jit/flags.h"
#include "tensorflow/compiler/mlir/tensorflow/transforms/host_runtime/runtime_passes.h"
#include "tensorflow/compiler/mlir/tensorflow/transforms/passes.h"
#include "tensorflow/compiler/mlir/tensorflow/transforms/sparsecore/sparsecore_passes.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/attribute_utils.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/data_dumper_logger_config.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/dump_mlir_util.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/error_util.h"
#include "xla/tsl/framework/device_type.h"
#include "tensorflow/core/framework/metrics.h"
#include "tensorflow/core/framework/types.h"
#include "tensorflow/core/platform/error_payloads.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/tpu/tpu_defs.h"
#include "tensorflow/core/util/debug_data_dumper.h"
#include "tsl/platform/error_logging.h"
#include "tsl/platform/errors.h"
namespace tensorflow {
namespace tfrt_compiler {
namespace {
using mlir::LogicalResult;
using mlir::OpPassManager;
using mlir::PassManager;
using mlir::func::FuncOp;
using mlir::TF::StandardPipelineOptions;
void EnablePassIRPrinting(PassManager& pm, const std::string& dump_group_name,
                          llvm::StringRef module_name) {
  pm.getContext()->disableMultithreading();
  pm.enableIRPrinting(std::make_unique<::tensorflow::DataDumperLoggerConfig>(
      [module_name, dump_group_name](const std::string& pass_tag_name,
                                     mlir::Operation* op) {
        return DEBUG_DATA_DUMPER()->GetDumpFilename(
            module_name.str(), dump_group_name, pass_tag_name);
      },
      "",
      true));
  pm.enableTiming();
}
}  
void AddTPULowerClusterToRuntimeOpsPassPipeline(OpPassManager& pm,
                                                llvm::StringRef module_name) {
  pm.addPass(mlir::TFTPU::CreateTPURewritePass(module_name));
  pm.addPass(mlir::createSymbolDCEPass());
  pm.addNestedPass<FuncOp>(
      mlir::TFDevice::CreateReplicateInvariantOpHoistingPass());
  pm.addNestedPass<FuncOp>(mlir::TFDevice::CreateEmbeddingProgramKeyPass());
  pm.addPass(mlir::TFTPU::CreateTPUMergeVariablesWithExecutePass());
  pm.addNestedPass<FuncOp>(
      mlir::TFTPU::CreateExtractTPUCopyWithDynamicShapeOpPass());
  pm.addNestedPass<FuncOp>(
      mlir::TFTPU::CreateTPUColocateCompositeResourceOps());
  if (tensorflow::GetMlirCommonFlags()
          ->tf_mlir_enable_tpu_variable_runtime_reformatting_pass) {
    pm.addPass(mlir::TFTPU::CreateTPUVariableRuntimeReformattingPass());
  }
}
void AddNonTPULowerClusterToRuntimeOpsPassPipeline(
    OpPassManager& pm, llvm::StringRef module_name) {
  pm.addPass(mlir::TFDevice::CreateXlaRewritePass());
  pm.addNestedPass<FuncOp>(mlir::createCanonicalizerPass());
  pm.addNestedPass<FuncOp>(mlir::createCSEPass());
  pm.addPass(mlir::createSymbolDCEPass());
}
void CreateTPULowerClusterToRuntimeOpsPassPipeline(
    OpPassManager& pm, const StandardPipelineOptions& options) {
  AddTPULowerClusterToRuntimeOpsPassPipeline(pm, "");
}
void CreateNonTPULowerClusterToRuntimeOpsPassPipeline(
    OpPassManager& pm, const StandardPipelineOptions& options) {
  AddNonTPULowerClusterToRuntimeOpsPassPipeline(pm, "");
}
tensorflow::Status RecordIfErrorStatus(const std::string error_prefix,
                                       std::string bridge_type,
                                       tsl::DeviceType device_type,
                                       absl::Status status) {
  if (status.ok()) {
    return status;
  }
  VLOG(2) << error_prefix << " " << status;
  tensorflow::metrics::UpdateTfMlirBridgeFirstPhaseCounter(
      bridge_type,
      mlir::TF::kMlirPh1BridgeCounterV2,
      device_type.type_string(),
      false,
      "failure");
  std::string bridge_subcomponent = "TFXLA_PHASE_ONE_MLIR_TPU_BRIDGE";
  tsl::OkOrSetErrorCounterPayload(
      tensorflow::core::platform::ErrorSourceProto::MLIR_BRIDGE_PHASE_1,
      status);
  if (device_type != DeviceType(DEVICE_TPU_XLA_JIT)) {
    bridge_subcomponent = "TFXLA_PHASE_ONE_MLIR_CPU/GPU_BRIDGE";
  }
  tsl::error_logging::Log(mlir::TF::kBridgeComponent, bridge_subcomponent,
                          status.ToString())
      .IgnoreError();
  return status;
}
absl::Status RunLowerClusterToRuntimeOpsPassPipeline(
    mlir::ModuleOp module, tsl::DeviceType xla_device_type,
    llvm::StringRef module_name) {
  PassManager runtime_lowering(module.getContext());
  ::tensorflow::applyTensorflowAndCLOptions(runtime_lowering);
  if (xla_device_type == DeviceType(DEVICE_TPU_XLA_JIT)) {
    AddTPULowerClusterToRuntimeOpsPassPipeline(runtime_lowering, module_name);
  } else {
    AddNonTPULowerClusterToRuntimeOpsPassPipeline(runtime_lowering,
                                                  module_name);
  }
  mlir::StatusScopedDiagnosticHandler diag_handler(
      module.getContext(), false,
      !VLOG_IS_ON(1));
  if (VLOG_IS_ON(1) ||
      DEBUG_DATA_DUMPER()->ShouldDump(module_name.str(), kDebugGroupMain)) {
    ::tensorflow::DumpMlirOpToFile(
        DEBUG_DATA_DUMPER()->GetDumpFilename(module_name.str(), kDebugGroupMain,
                                             "runtime_lowering_before"),
        module, llvm::StringRef(), &runtime_lowering);
  }
  if (VLOG_IS_ON(2) || DEBUG_DATA_DUMPER()->ShouldDump(
                           module_name.str(), kDebugGroupRuntimeLowering)) {
    EnablePassIRPrinting(runtime_lowering, kDebugGroupRuntimeLowering,
                         module_name);
  }
  LogicalResult result = runtime_lowering.run(module);
  (void)result;
  if (VLOG_IS_ON(1) ||
      DEBUG_DATA_DUMPER()->ShouldDump(module_name.str(), kDebugGroupMain)) {
    ::tensorflow::DumpMlirOpToFile(
        DEBUG_DATA_DUMPER()->GetDumpFilename(module_name.str(), kDebugGroupMain,
                                             "runtime_lowering_after"),
        module, llvm::StringRef(), &runtime_lowering);
  }
  std::string bridge_type = xla_device_type == DeviceType(DEVICE_TPU_XLA_JIT)
                                ? mlir::TF::kMlirPh1BridgeCounterReplicated
                                : mlir::TF::kMlirPh1BridgeCounterNonReplicated;
  auto result_status = diag_handler.ConsumeStatus();
  TF_RETURN_IF_ERROR(
      RecordIfErrorStatus("lower_cluster_to_runtime",
                          bridge_type, xla_device_type, result_status));
  return absl::OkStatus();
}
void RegisterTPULowerClusterToRuntimeOpsPassPipeline() {
  static mlir::PassPipelineRegistration<StandardPipelineOptions> pipeline(
      "tfrt-lower-cluster-to-runtime-ops-tpu",
      "Run all the passes involved after the clustering transformations from "
      "the TF2XLA Bridge. Takes as input a Module with tf_device.cluster ops "
      "and outputs TFRT runtime ops such as TPUCompile. This pipeline is for "
      "TPU.",
      CreateTPULowerClusterToRuntimeOpsPassPipeline);
}
void RegisterNonTPULowerClusterToRuntimeOpsPassPipeline() {
  static mlir::PassPipelineRegistration<StandardPipelineOptions> pipeline(
      "tfrt-lower-cluster-to-runtime-ops-non-tpu",
      "Run all the passes involved after the clustering transformations from "
      "the TF2XLA Bridge. Takes as input a Module with tf_device.cluster ops "
      "and outputs TFRT runtime ops such as XlaLaunch. This is for CPU/GPU",
      CreateNonTPULowerClusterToRuntimeOpsPassPipeline);
}
}  
}  