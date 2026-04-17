#ifndef TENSORFLOW_LITE_TOOLS_EVALUATION_EVALUATION_DELEGATE_PROVIDER_H_
#define TENSORFLOW_LITE_TOOLS_EVALUATION_EVALUATION_DELEGATE_PROVIDER_H_
#include <string>
#include <unordered_map>
#include <vector>
#include "tensorflow/lite/tools/command_line_flags.h"
#include "tensorflow/lite/tools/delegates/delegate_provider.h"
#include "tensorflow/lite/tools/evaluation/proto/evaluation_stages.pb.h"
#include "tensorflow/lite/tools/evaluation/utils.h"
#include "tensorflow/lite/tools/tool_params.h"
namespace tflite {
namespace evaluation {
using ProvidedDelegateList = tflite::tools::ProvidedDelegateList;
class DelegateProviders {
 public:
  DelegateProviders();
  std::vector<Flag> GetFlags();
  bool InitFromCmdlineArgs(int* argc, const char** argv);
  const tools::ToolParams& GetAllParams() const { return params_; }
  tools::ToolParams GetAllParams(const TfliteInferenceParams& params) const;
  TfLiteDelegatePtr CreateDelegate(const std::string& name) const;
  std::vector<ProvidedDelegateList::ProvidedDelegate> CreateAllDelegates()
      const {
    return delegate_list_util_.CreateAllRankedDelegates();
  }
  std::vector<ProvidedDelegateList::ProvidedDelegate> CreateAllDelegates(
      const TfliteInferenceParams& params) const {
    auto converted = GetAllParams(params);
    ProvidedDelegateList util(&converted);
    return util.CreateAllRankedDelegates();
  }
 private:
  tools::ToolParams params_;
  ProvidedDelegateList delegate_list_util_;
  const std::unordered_map<std::string, int> delegates_map_;
};
TfliteInferenceParams::Delegate ParseStringToDelegateType(
    const std::string& val);
TfLiteDelegatePtr CreateTfLiteDelegate(const TfliteInferenceParams& params,
                                       std::string* error_msg = nullptr);
}  
}  
#endif  
#include "tensorflow/lite/tools/evaluation/evaluation_delegate_provider.h"
#include <string>
#include "tensorflow/lite/c/c_api_types.h"
#include "tensorflow/lite/tools/command_line_flags.h"
#include "tensorflow/lite/tools/evaluation/proto/evaluation_stages.pb.h"
#include "tensorflow/lite/tools/evaluation/utils.h"
#include "tensorflow/lite/tools/logging.h"
#include "tensorflow/lite/tools/tool_params.h"
namespace tflite {
namespace evaluation {
namespace {
constexpr char kNnapiDelegate[] = "nnapi";
constexpr char kGpuDelegate[] = "gpu";
constexpr char kHexagonDelegate[] = "hexagon";
constexpr char kXnnpackDelegate[] = "xnnpack";
constexpr char kCoremlDelegate[] = "coreml";
}  
TfliteInferenceParams::Delegate ParseStringToDelegateType(
    const std::string& val) {
  if (val == kNnapiDelegate) return TfliteInferenceParams::NNAPI;
  if (val == kGpuDelegate) return TfliteInferenceParams::GPU;
  if (val == kHexagonDelegate) return TfliteInferenceParams::HEXAGON;
  if (val == kXnnpackDelegate) return TfliteInferenceParams::XNNPACK;
  if (val == kCoremlDelegate) return TfliteInferenceParams::COREML;
  return TfliteInferenceParams::NONE;
}
TfLiteDelegatePtr CreateTfLiteDelegate(const TfliteInferenceParams& params,
                                       std::string* error_msg) {
  const auto type = params.delegate();
  switch (type) {
    case TfliteInferenceParams::NNAPI: {
      auto p = CreateNNAPIDelegate();
      if (!p && error_msg) *error_msg = "NNAPI not supported";
      return p;
    }
    case TfliteInferenceParams::GPU: {
      auto p = CreateGPUDelegate();
      if (!p && error_msg) *error_msg = "GPU delegate not supported.";
      return p;
    }
    case TfliteInferenceParams::HEXAGON: {
      auto p = CreateHexagonDelegate("",
                                     false);
      if (!p && error_msg) {
        *error_msg =
            "Hexagon delegate is not supported on the platform or required "
            "libraries are missing.";
      }
      return p;
    }
    case TfliteInferenceParams::XNNPACK: {
      auto p = CreateXNNPACKDelegate(params.num_threads(), false);
      if (!p && error_msg) *error_msg = "XNNPACK delegate not supported.";
      return p;
    }
    case TfliteInferenceParams::COREML: {
      auto p = CreateCoreMlDelegate();
      if (!p && error_msg) *error_msg = "CoreML delegate not supported.";
      return p;
    }
    case TfliteInferenceParams::NONE:
      return TfLiteDelegatePtr(nullptr, [](TfLiteDelegate*) {});
    default:
      if (error_msg) {
        *error_msg = "Creation of delegate type: " +
                     TfliteInferenceParams::Delegate_Name(type) +
                     " not supported yet.";
      }
      return TfLiteDelegatePtr(nullptr, [](TfLiteDelegate*) {});
  }
}
DelegateProviders::DelegateProviders()
    : delegate_list_util_(&params_),
      delegates_map_([=]() -> std::unordered_map<std::string, int> {
        std::unordered_map<std::string, int> delegates_map;
        const auto& providers = delegate_list_util_.providers();
        for (int i = 0; i < providers.size(); ++i) {
          delegates_map[providers[i]->GetName()] = i;
        }
        return delegates_map;
      }()) {
  delegate_list_util_.AddAllDelegateParams();
}
std::vector<Flag> DelegateProviders::GetFlags() {
  std::vector<Flag> flags;
  delegate_list_util_.AppendCmdlineFlags(flags);
  return flags;
}
bool DelegateProviders::InitFromCmdlineArgs(int* argc, const char** argv) {
  std::vector<Flag> flags = GetFlags();
  bool parse_result = Flags::Parse(argc, argv, flags);
  if (!parse_result || params_.Get<bool>("help")) {
    std::string usage = Flags::Usage(argv[0], flags);
    TFLITE_LOG(ERROR) << usage;
    parse_result = false;
  }
  return parse_result;
}
TfLiteDelegatePtr DelegateProviders::CreateDelegate(
    const std::string& name) const {
  const auto it = delegates_map_.find(name);
  if (it == delegates_map_.end()) {
    return TfLiteDelegatePtr(nullptr, [](TfLiteDelegate*) {});
  }
  const auto& providers = delegate_list_util_.providers();
  return providers[it->second]->CreateTfLiteDelegate(params_);
}
tools::ToolParams DelegateProviders::GetAllParams(
    const TfliteInferenceParams& params) const {
  tools::ToolParams tool_params;
  tool_params.Merge(params_,  false);
  if (params.has_num_threads()) {
    tool_params.Set<int32_t>("num_threads", params.num_threads());
  }
  const auto type = params.delegate();
  switch (type) {
    case TfliteInferenceParams::NNAPI:
      if (tool_params.HasParam("use_nnapi")) {
        tool_params.Set<bool>("use_nnapi", true);
      }
      break;
    case TfliteInferenceParams::GPU:
      if (tool_params.HasParam("use_gpu")) {
        tool_params.Set<bool>("use_gpu", true);
      }
      break;
    case TfliteInferenceParams::HEXAGON:
      if (tool_params.HasParam("use_hexagon")) {
        tool_params.Set<bool>("use_hexagon", true);
      }
      break;
    case TfliteInferenceParams::XNNPACK:
      if (tool_params.HasParam("use_xnnpack")) {
        tool_params.Set<bool>("use_xnnpack", true);
      }
      if (tool_params.HasParam("xnnpack_force_fp16")) {
        tool_params.Set<bool>("xnnpack_force_fp16", true);
      }
      break;
    case TfliteInferenceParams::COREML:
      if (tool_params.HasParam("use_coreml")) {
        tool_params.Set<bool>("use_coreml", true);
      }
      break;
    default:
      break;
  }
  return tool_params;
}
}  
}  