#ifndef TENSORFLOW_CORE_FUNCTION_RUNTIME_CLIENT_RUNTIME_CLIENT_H_
#define TENSORFLOW_CORE_FUNCTION_RUNTIME_CLIENT_RUNTIME_CLIENT_H_
#include <vector>
#include "absl/types/span.h"
#include "tensorflow/c/eager/abstract_tensor_handle.h"
#include "tensorflow/c/eager/immediate_execution_tensor_handle.h"
#include "tensorflow/core/common_runtime/eager/context.h"
#include "tensorflow/core/framework/function.pb.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/statusor.h"
#include "tensorflow/core/platform/stringpiece.h"
namespace tensorflow {
namespace core {
namespace function {
struct OpaqueTfgGraphFuncOp;
struct OpaqueTfFuncOp;
EagerContext& GlobalPythonEagerContext();
EagerContext& GlobalEagerContext();
using ReturnValues = std::vector<ImmediateTensorHandlePtr>;
class Runtime {
 public:
  explicit Runtime(EagerContext& eager_ctx) : eager_ctx_(eager_ctx) {}
  enum class Dialect {
    TFG,
    TF,
  };
  absl::StatusOr<FunctionDef> GetFunctionProto(StringPiece name);
  Status CreateFunction(const FunctionDef& fdef);
  Status CreateFunction(OpaqueTfgGraphFuncOp* fop);
  Status CreateFunction(OpaqueTfFuncOp* fop);
  Status TransformFunction(StringPiece name, StringPiece pipeline_name,
                           Dialect dialect = Dialect::TFG);
  absl::StatusOr<ReturnValues> CallFunction(
      StringPiece name, absl::Span<AbstractTensorHandle* const> args);
 private:
  EagerContext& eager_ctx_;
};
}  
}  
}  
#endif  
#include "tensorflow/core/function/runtime_client/runtime_client.h"
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/strings/match.h"
#include "absl/strings/str_cat.h"
#include "absl/types/span.h"
#include "mlir/IR/MLIRContext.h"  
#include "mlir/Pass/PassManager.h"  
#include "mlir/Pass/PassRegistry.h"  
#include "mlir/Support/LogicalResult.h"  
#include "tensorflow/c/eager/abstract_tensor_handle.h"
#include "tensorflow/c/eager/immediate_execution_context.h"
#include "tensorflow/c/eager/immediate_execution_operation.h"
#include "tensorflow/c/eager/immediate_execution_tensor_handle.h"
#if !defined(DISABLE_MLIR)
#include "tensorflow/compiler/mlir/python/mlir.h"
#endif
#include "tensorflow/compiler/mlir/tensorflow/translate/export_graphdef.h"
#include "tensorflow/compiler/mlir/tensorflow/translate/import_model.h"
#include "tensorflow/compiler/mlir/tensorflow/utils/error_util.h"
#include "tensorflow/compiler/mlir/tf2xla/api/v2/tf_executor_to_graph.h"
#include "tensorflow/core/common_runtime/device_mgr.h"
#include "tensorflow/core/common_runtime/eager/context.h"
#include "tensorflow/core/common_runtime/function_def_utils.h"
#include "tensorflow/core/framework/device.h"
#include "tensorflow/core/framework/device_factory.h"
#include "tensorflow/core/framework/function.pb.h"
#include "tensorflow/core/framework/graph.pb.h"
#include "tensorflow/core/framework/op_def.pb.h"
#include "tensorflow/core/ir/importexport/graphdef_export.h"
#include "tensorflow/core/ir/importexport/graphdef_import.h"
#include "tensorflow/core/ir/ops.h"
#include "tensorflow/core/platform/errors.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/statusor.h"
#include "tensorflow/core/platform/stringpiece.h"
#include "tensorflow/core/platform/types.h"
#include "tensorflow/core/protobuf/error_codes.pb.h"
#include "tensorflow/core/public/session_options.h"
namespace tensorflow {
namespace core {
namespace function {
EagerContext& GlobalEagerContext() {
  static EagerContext* global_ctx = []() {
    SessionOptions opts;
    std::vector<std::unique_ptr<Device>> devices;
    Status&& device_init_status = DeviceFactory::AddDevices(
        opts, "/job:localhost/replica:0/task:0", &devices);
    CHECK(device_init_status.ok());  
    return new EagerContext(
        opts, ContextDevicePlacementPolicy::DEVICE_PLACEMENT_SILENT,
        false,
        new DynamicDeviceMgr(std::move(devices)),
        true,
        nullptr,
        nullptr,
        nullptr,
        true);
  }();
  return *global_ctx;
}
EagerContext& GlobalPythonEagerContext() {
  EagerContext* ctx = reinterpret_cast<EagerContext*>(GetCEagerContext());
  DCHECK(ctx) << "The Python eager context must be initialized first.";
  return *ctx;
}
absl::StatusOr<FunctionDef> Runtime::GetFunctionProto(StringPiece name) {
  EagerContext& ctx = this->eager_ctx_;
  const FunctionDef* f = ctx.FindFunctionDef(std::string(name));
  if (f == nullptr) {
    return Status(absl::StatusCode::kInvalidArgument,
                  absl::StrCat("Could not find an attribute for key ", name));
  }
  return *f;
}
Status Runtime::CreateFunction(const FunctionDef& fdef) {
  const auto& fname = fdef.signature().name();
  if (this->eager_ctx_.FindFunctionByName(fname)) {
    TF_RETURN_WITH_CONTEXT_IF_ERROR(this->eager_ctx_.RemoveFunction(fname),
                                    "removing function ", fname);
  }
  return this->eager_ctx_.AddFunctionDef(fdef);
}
Status Runtime::CreateFunction(OpaqueTfgGraphFuncOp* fop) {
  mlir::tfg::GraphFuncOp fop_proper =
      *reinterpret_cast<mlir::tfg::GraphFuncOp*>(fop);
  return mlir::tfg::ConvertToFunctionDef(fop_proper,
                                         *this->eager_ctx_.FuncLibDef());
}
Status Runtime::CreateFunction(OpaqueTfFuncOp* fop) {
  mlir::func::FuncOp fop_proper = *reinterpret_cast<mlir::func::FuncOp*>(fop);
  const auto& fname = fop_proper.getName().str();
  GraphExportConfig config;
  FunctionDef fdef;
  TF_RETURN_WITH_CONTEXT_IF_ERROR(
      tf2xla::v2::ConvertMlirFunctionToFunctionLibraryDef(fop_proper, config,
                                                          &fdef),
      "creating function ", fname);
  return CreateFunction(fdef);
}
Status Runtime::TransformFunction(StringPiece name, StringPiece pipeline_name,
                                  Dialect dialect) {
  mlir::MLIRContext ctx;
  mlir::PassManager pm(&ctx);
  std::string error;
  llvm::raw_string_ostream error_stream(error);
  if (mlir::failed(mlir::parsePassPipeline(std::string(pipeline_name), pm,
                                           error_stream))) {
    return Status(absl::StatusCode::kInvalidArgument,
                  absl::StrCat("locating pass pipeline ", pipeline_name, ": ",
                               error_stream.str()));
  }
  auto fn = GetFunctionProto(name);
  TF_RETURN_WITH_CONTEXT_IF_ERROR(fn.status(), "loading function ", name);
  GraphDef graph;
  *graph.mutable_library()->add_function() = *fn;
  tensorflow::GraphDebugInfo debug_info;
  if (dialect == Dialect::TFG) {
    auto mlir_fn = mlir::tfg::ImportGraphDef(&ctx, debug_info, graph);
    TF_RETURN_WITH_CONTEXT_IF_ERROR(mlir_fn.status(), "importing function ",
                                    name);
    mlir::StatusScopedDiagnosticHandler diagnostics_handler(&ctx);
    if (failed(pm.run(mlir_fn->get()))) {
      return diagnostics_handler.Combine(
          Status(absl::StatusCode::kInvalidArgument,
                 absl::StrCat("running pass pipeline ", pipeline_name, ": ")));
    }
    for (auto fn : mlir_fn->get().getBody()->getOps<mlir::tfg::GraphFuncOp>()) {
      TF_RETURN_WITH_CONTEXT_IF_ERROR(
          CreateFunction(reinterpret_cast<OpaqueTfgGraphFuncOp*>(&fn)),
          absl::StrCat("updating function ", fn.getName().str()));
    }
    return absl::OkStatus();
  }
  if (dialect == Dialect::TF) {
    Status status;
    FunctionLibraryDefinition& flib_def = *this->eager_ctx_.FuncLibDef();
    std::unique_ptr<FunctionBody> fbody;
    status = FunctionDefToBodyHelper(*fn, AttrSlice(), &flib_def, &fbody);
    TF_RETURN_WITH_CONTEXT_IF_ERROR(status, "importing function ", name);
    auto mlir_fn = ConvertFunctionToMlir(fbody.get(), flib_def, &ctx);
    TF_RETURN_WITH_CONTEXT_IF_ERROR(mlir_fn.status(), "importing function ",
                                    name);
    mlir::StatusScopedDiagnosticHandler diagnostics_handler(&ctx);
    if (failed(pm.run(mlir_fn->get()))) {
      return diagnostics_handler.Combine(
          Status(absl::StatusCode::kInvalidArgument,
                 absl::StrCat("running pass pipeline ", pipeline_name, ": ")));
    }
    for (auto fn : mlir_fn->get().getBody()->getOps<mlir::func::FuncOp>()) {
      TF_RETURN_WITH_CONTEXT_IF_ERROR(
          CreateFunction(reinterpret_cast<OpaqueTfFuncOp*>(&fn)),
          absl::StrCat("updating function ", fn.getName().str()));
    }
    return absl::OkStatus();
  }
  return Status(
      absl::StatusCode::kInvalidArgument,
      absl::StrCat("Unsupported dialect: ", dialect,
                   ". Supported dialects are Dialect::TFG and Dialect::TF."));
}
absl::StatusOr<ReturnValues> Runtime::CallFunction(
    StringPiece name, absl::Span<AbstractTensorHandle* const> args) {
  EagerContext& ctx = this->eager_ctx_;
  ImmediateOpPtr op(ctx.CreateOperation());
  TF_RETURN_WITH_CONTEXT_IF_ERROR(op->Reset(name.data(), nullptr),
                                  "initializing call op for ", name);
  TF_RETURN_WITH_CONTEXT_IF_ERROR(op->AddInputList(args),
                                  "preparing call args for ", name);
  const FunctionDef* fn_def = ctx.GetFunctionDef(string(name));
  int num_retvals = fn_def->signature().output_arg_size();
  int actual_retvals = num_retvals;
  std::vector<ImmediateExecutionTensorHandle*> retvals(num_retvals);
  TF_RETURN_WITH_CONTEXT_IF_ERROR(
      op->Execute(absl::MakeSpan(
                      reinterpret_cast<AbstractTensorHandle**>(retvals.data()),
                      num_retvals),
                  &actual_retvals),
      "executing call op for ", name);
  DCHECK(num_retvals == actual_retvals);
  ReturnValues final_returns;
  for (const auto& r : retvals) {
    final_returns.emplace_back(ImmediateTensorHandlePtr(r));
  }
  return final_returns;
}
}  
}  
}  