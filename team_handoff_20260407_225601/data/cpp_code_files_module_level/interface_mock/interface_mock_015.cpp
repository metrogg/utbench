#ifndef TENSORFLOW_COMPILER_MLIR_TFRT_UTILS_EXPORT_H_
#define TENSORFLOW_COMPILER_MLIR_TFRT_UTILS_EXPORT_H_
#include "absl/functional/any_invocable.h"
#include "absl/status/status.h"
#include "mlir/IR/BuiltinOps.h"  
#include "tensorflow/core/framework/function.pb.h"
namespace tensorflow {
absl::Status ExportFunctionDefs(
    mlir::ModuleOp module,
    absl::AnyInvocable<absl::Status(tensorflow::FunctionDef)> callback,
    bool export_tf_original_func_name = true);
}  
#endif  
#include "tensorflow/compiler/mlir/tfrt/utils/export.h"
#include <memory>
#include <utility>
#include "absl/functional/any_invocable.h"
#include "absl/status/status.h"
#include "absl/strings/string_view.h"
#include "mlir/Dialect/Func/IR/FuncOps.h"  
#include "mlir/IR/BuiltinOps.h"  
#include "mlir/Pass/PassManager.h"  
#include "mlir/Support/LogicalResult.h"  
#include "tensorflow/compiler/mlir/tensorflow/transforms/passes.h"
#include "tensorflow/compiler/mlir/tensorflow/translate/mlir_roundtrip_flags.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/error_util.h"
#include "tensorflow/compiler/mlir/tf2xla/api/v1/tf_dialect_to_executor.h"
#include "tensorflow/compiler/mlir/tf2xla/api/v2/tf_executor_to_graph.h"
#include "tensorflow/core/framework/function.pb.h"
#include "tsl/platform/errors.h"
#include "tsl/profiler/lib/traceme.h"
namespace tensorflow {
absl::Status ExportFunctionDefs(
    mlir::ModuleOp module,
    absl::AnyInvocable<absl::Status(tensorflow::FunctionDef)> callback,
    bool export_tf_original_func_name) {
  tsl::profiler::TraceMe traceme([&]() {
    return tsl::profiler::TraceMeEncode(
        "ExportFunctionDefs",
        {{"module_name", absl::string_view(module.getName().value_or("?"))}});
  });
  TF_RETURN_IF_ERROR(
      tensorflow::tf2xla::v1::ExportFromTensorflowDialectToExecutor(module));
  {
    mlir::StatusScopedDiagnosticHandler diag_handler(module.getContext());
    mlir::PassManager pm(module.getContext());
    pm.addPass(mlir::CreateBreakUpIslandsPass());
    if (mlir::failed(pm.run(module))) {
      return diag_handler.ConsumeStatus();
    }
  }
  tensorflow::GraphExportConfig configs;
  configs.export_original_tf_func_name = export_tf_original_func_name;
  for (auto func : module.getOps<mlir::func::FuncOp>()) {
    tensorflow::FunctionDef function_def;
    TF_RETURN_IF_ERROR(
        tensorflow::tf2xla::v2::ConvertMlirFunctionToFunctionLibraryDef(
            func, configs, &function_def));
    TF_RETURN_IF_ERROR(callback(std::move(function_def)));
  }
  return absl::OkStatus();
}
}  