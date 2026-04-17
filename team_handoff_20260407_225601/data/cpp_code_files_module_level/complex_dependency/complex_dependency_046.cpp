#ifndef TENSORFLOW_COMPILER_MLIR_LITE_METRICS_ERROR_COLLECTOR_INST_H_
#define TENSORFLOW_COMPILER_MLIR_LITE_METRICS_ERROR_COLLECTOR_INST_H_
#include <memory>
#include <string>
#include <unordered_map>
#include <utility>
#include "mlir/IR/Diagnostics.h"  
#include "mlir/IR/Location.h"  
#include "mlir/IR/MLIRContext.h"  
#include "mlir/IR/Operation.h"  
#include "mlir/Pass/PassInstrumentation.h"  
#include "tensorflow/compiler/mlir/lite/metrics/error_collector.h"
#include "tensorflow/compiler/mlir/lite/metrics/types_util.h"
#include "tensorflow/lite/python/metrics/converter_error_data.pb.h"
namespace mlir {
namespace TFL {
class ErrorCollectorInstrumentation : public PassInstrumentation {
  using ConverterErrorData = tflite::metrics::ConverterErrorData;
  using ErrorCode = ConverterErrorData::ErrorCode;
 public:
  explicit ErrorCollectorInstrumentation(MLIRContext *context);
 private:
  void runBeforePass(Pass *pass, Operation *module) override;
  void runAfterPass(Pass *pass, Operation *module) override;
  void runAfterPassFailed(Pass *pass, Operation *module) override;
  std::unique_ptr<ScopedDiagnosticHandler> handler_;
  std::unordered_map<Location, std::string, LocationHash> loc_to_name_;
  std::string common_error_message_;
  std::string pass_name_;
  ErrorCollector *error_collector_;
};
constexpr char kErrorCodePrefix[] = "Error code: ";
inline InFlightDiagnostic AttachErrorCode(InFlightDiagnostic &&diag,
                                          int error_code) {
  using tflite::metrics::ConverterErrorData;
  diag.attachNote() << kErrorCodePrefix
                    << ConverterErrorData::ErrorCode_Name(error_code);
  return std::move(diag);
}
}  
}  
#endif  
#include "tensorflow/compiler/mlir/lite/metrics/error_collector_inst.h"
#include <memory>
#include <string>
#include <vector>
#include "absl/strings/match.h"
#include "absl/strings/str_split.h"
#include "mlir/IR/Diagnostics.h"  
#include "mlir/IR/Location.h"  
#include "mlir/Pass/Pass.h"  
#include "mlir/Support/LogicalResult.h"  
#include "tensorflow/compiler/mlir/lite/metrics/error_collector.h"
#include "tensorflow/compiler/mlir/lite/metrics/types_util.h"
namespace mlir {
namespace TFL {
namespace {
inline std::string extract_pass_name(const std::string &signature) {
  const std::vector<std::string> &v = absl::StrSplit(signature, "::");
  return v.back();
}
inline std::string extract_op_name_from_error_message(
    const std::string &error_message) {
  int end_pos = error_message.find("' op");
  if ((absl::StartsWith(error_message, "'tf.") ||
       absl::StartsWith(error_message, "'tfl.")) &&
      end_pos != std::string::npos) {
    return error_message.substr(1, end_pos - 1);
  }
  return "";
}
const int kMaxAcceptedNoteSize = 1024;
}  
ErrorCollectorInstrumentation::ErrorCollectorInstrumentation(
    MLIRContext *context)
    : error_collector_(ErrorCollector::GetErrorCollector()) {
  handler_ = std::make_unique<ScopedDiagnosticHandler>(
      context, [this](Diagnostic &diag) {
        if (diag.getSeverity() == DiagnosticSeverity::Error) {
          Location loc = diag.getLocation();
          std::string error_message = diag.str();
          std::string op_name, error_code;
          if (loc_to_name_.count(loc)) {
            op_name = loc_to_name_[loc];
          } else {
            op_name = extract_op_name_from_error_message(diag.str());
          }
          for (const auto &note : diag.getNotes()) {
            const std::string note_str = note.str();
            if (absl::StartsWith(note_str, kErrorCodePrefix)) {
              error_code = note_str.substr(sizeof(kErrorCodePrefix) - 1);
            }
            error_message += "\n";
            if (note_str.size() <= kMaxAcceptedNoteSize) {
              error_message += note_str;
            } else {
              error_message += note_str.substr(0, kMaxAcceptedNoteSize);
              error_message += "...";
            }
          }
          ErrorCode error_code_enum = ConverterErrorData::UNKNOWN;
          bool has_valid_error_code =
              ConverterErrorData::ErrorCode_Parse(error_code, &error_code_enum);
          if (!op_name.empty() || has_valid_error_code) {
            error_collector_->ReportError(NewConverterErrorData(
                pass_name_, error_message, error_code_enum, op_name, loc));
          } else {
            common_error_message_ += diag.str();
            common_error_message_ += "\n";
          }
        }
        return failure();
      });
}
void ErrorCollectorInstrumentation::runBeforePass(Pass *pass,
                                                  Operation *module) {
  auto collectOps = [this](Operation *op) {
    const auto &op_name = op->getName().getStringRef().str();
    if (absl::StartsWith(op_name, "tf.") || absl::StartsWith(op_name, "tfl.")) {
      loc_to_name_.emplace(op->getLoc(), op_name);
    }
  };
  for (auto &region : module->getRegions()) {
    region.walk(collectOps);
  }
  pass_name_ = extract_pass_name(pass->getName().str());
  error_collector_->Clear();
}
void ErrorCollectorInstrumentation::runAfterPass(Pass *pass,
                                                 Operation *module) {
  loc_to_name_.clear();
  pass_name_.clear();
  common_error_message_.clear();
  error_collector_->Clear();
}
void ErrorCollectorInstrumentation::runAfterPassFailed(Pass *pass,
                                                       Operation *module) {
  if (error_collector_->CollectedErrors().empty() &&
      !common_error_message_.empty()) {
    error_collector_->ReportError(NewConverterErrorData(
        pass_name_, common_error_message_, ConverterErrorData::UNKNOWN,
        "", module->getLoc()));
  }
  loc_to_name_.clear();
  pass_name_.clear();
  common_error_message_.clear();
}
}  
}  