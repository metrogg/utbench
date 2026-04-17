#ifndef TENSORFLOW_CORE_COMMON_RUNTIME_NEXT_PLUGGABLE_DEVICE_C_PLUGIN_COORDINATION_SERVICE_AGENT_H_
#define TENSORFLOW_CORE_COMMON_RUNTIME_NEXT_PLUGGABLE_DEVICE_C_PLUGIN_COORDINATION_SERVICE_AGENT_H_
#include <cstdint>
#include <string>
#include <string_view>
#include "absl/time/time.h"
#include "tensorflow/c/experimental/next_pluggable_device/c_api.h"
#include "tensorflow/c/kernels_experimental.h"
#include "tensorflow/core/common_runtime/next_pluggable_device/plugin_coordination_service_agent.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/statusor.h"
namespace tensorflow {
class CPluginCoordinationServiceAgent : public PluginCoordinationServiceAgent {
 public:
  explicit CPluginCoordinationServiceAgent(void* agent)
      : agent_(reinterpret_cast<TF_CoordinationServiceAgent*>(agent)) {}
  bool IsInitialized() const override {
    if (agent_ == nullptr) return false;
    return TF_CoordinationServiceIsInitialized(agent_);
  }
  Status InsertKeyValue(std::string_view key, std::string_view value) override;
  absl::StatusOr<std::string> GetKeyValue(std::string_view key) override;
  absl::StatusOr<std::string> GetKeyValue(std::string_view key,
                                          absl::Duration timeout) override;
  absl::StatusOr<std::string> TryGetKeyValue(std::string_view key) override;
  Status DeleteKeyValue(std::string_view key) override;
 private:
  TF_CoordinationServiceAgent* agent_;
};
}  
#endif  
#include "tensorflow/core/common_runtime/next_pluggable_device/c_plugin_coordination_service_agent.h"
#include <string>
#include <string_view>
#include "absl/time/time.h"
#include "tensorflow/c/experimental/next_pluggable_device/c_api.h"
#include "tensorflow/c/tf_buffer.h"
#include "tensorflow/c/tf_status.h"
#include "tensorflow/c/tf_status_helper.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/statusor.h"
namespace tensorflow {
namespace {
absl::StatusOr<std::string> ProcessGetKeyValueResult(TF_Buffer* result_buf,
                                                     TF_Status* status) {
  if (TF_GetCode(status) != TF_OK) {
    return StatusFromTF_Status(status);
  } else {
    std::string result{static_cast<const char*>(result_buf->data),
                       result_buf->length};
    TF_DeleteBuffer(result_buf);
    return result;
  }
}
}  
Status CPluginCoordinationServiceAgent::InsertKeyValue(std::string_view key,
                                                       std::string_view value) {
  TF_StatusPtr c_status_ptr(TF_NewStatus());
  TF_Status* status = c_status_ptr.get();
  TF_CoordinationServiceInsertKeyValue(key.data(), key.size(), value.data(),
                                       value.size(), agent_, status);
  return StatusFromTF_Status(status);
}
absl::StatusOr<std::string> CPluginCoordinationServiceAgent::GetKeyValue(
    std::string_view key) {
  TF_StatusPtr c_status_ptr(TF_NewStatus());
  TF_Status* status = c_status_ptr.get();
  TF_Buffer* result_buf =
      TF_CoordinationServiceGetKeyValue(key.data(), key.size(), agent_, status);
  return ProcessGetKeyValueResult(result_buf, status);
}
absl::StatusOr<std::string> CPluginCoordinationServiceAgent::GetKeyValue(
    std::string_view key, absl::Duration timeout) {
  TF_StatusPtr c_status_ptr(TF_NewStatus());
  TF_Status* status = c_status_ptr.get();
  TF_Buffer* result_buf = TF_CoordinationServiceGetKeyValueWithTimeout(
      key.data(), key.size(), absl::ToInt64Seconds(timeout), agent_, status);
  return ProcessGetKeyValueResult(result_buf, status);
}
absl::StatusOr<std::string> CPluginCoordinationServiceAgent::TryGetKeyValue(
    std::string_view key) {
  TF_StatusPtr c_status_ptr(TF_NewStatus());
  TF_Status* status = c_status_ptr.get();
  TF_Buffer* result_buf = TF_CoordinationServiceTryGetKeyValue(
      key.data(), key.size(), agent_, status);
  return ProcessGetKeyValueResult(result_buf, status);
}
Status CPluginCoordinationServiceAgent::DeleteKeyValue(std::string_view key) {
  TF_StatusPtr c_status_ptr(TF_NewStatus());
  TF_Status* status = c_status_ptr.get();
  TF_CoordinationServiceDeleteKeyValue(key.data(), key.size(), agent_, status);
  return StatusFromTF_Status(status);
}
}  