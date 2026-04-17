#ifndef TENSORFLOW_CORE_TFRT_FALLBACK_FALLBACK_STATE_H_
#define TENSORFLOW_CORE_TFRT_FALLBACK_FALLBACK_STATE_H_
#include <memory>
#include <vector>
#include "tensorflow/core/common_runtime/device_mgr.h"
#include "tensorflow/core/common_runtime/device_set.h"
#include "tensorflow/core/common_runtime/graph_execution_state.h"
#include "tensorflow/core/common_runtime/process_function_library_runtime.h"
#include "tensorflow/core/framework/device.h"
#include "tensorflow/core/framework/graph.pb.h"
#include "tensorflow/core/public/session_options.h"
namespace tensorflow {
namespace tfrt_stub {
class FallbackState {
 public:
  static absl::StatusOr<std::unique_ptr<FallbackState>> Create(
      const SessionOptions &session_options,
      const tensorflow::FunctionDefLibrary &fdef_lib);
  static absl::StatusOr<std::unique_ptr<FallbackState>> CreateWithCpuDevice(
      const SessionOptions &session_options,
      const tensorflow::FunctionDefLibrary &fdef_lib);
  static absl::StatusOr<std::unique_ptr<FallbackState>> CreateWithMockGpuDevice(
      const SessionOptions &session_options,
      const tensorflow::FunctionDefLibrary &fdef_lib);
  FallbackState(const SessionOptions &session_options,
                std::vector<std::unique_ptr<Device>> devices,
                const tensorflow::FunctionDefLibrary &fdef_lib);
  absl::StatusOr<std::unique_ptr<GraphExecutionState>>
  CreateGraphExecutionState(GraphDef graph_def, bool run_placer = true) const;
  Status AddFunctionDef(const FunctionDef &func_def);
  const SessionOptions &session_options() const { return session_options_; }
  const DeviceMgr &device_manager() const { return device_manager_; }
  DeviceMgr &device_manager() { return device_manager_; }
  const DeviceSet &device_set() const { return device_set_; }
  const ProcessFunctionLibraryRuntime &process_function_library_runtime()
      const {
    return pflr_;
  }
  const FunctionLibraryDefinition &func_lib_def() const {
    return func_lib_def_;
  }
 private:
  SessionOptions session_options_;
  StaticDeviceMgr device_manager_;
  DeviceSet device_set_;
  FunctionLibraryDefinition func_lib_def_;
  ProcessFunctionLibraryRuntime pflr_;
};
}  
}  
#endif  
#include "tensorflow/core/tfrt/fallback/fallback_state.h"
#include <cstddef>
#include <cstdint>
#include <memory>
#include <utility>
#include <vector>
#include "absl/strings/string_view.h"
#include "absl/strings/strip.h"
#include "tensorflow/core/common_runtime/device_mgr.h"
#include "tensorflow/core/common_runtime/rendezvous_mgr.h"
#include "tensorflow/core/framework/device_attributes.pb.h"
#include "tensorflow/core/framework/device_factory.h"
#include "tensorflow/core/framework/types.h"
#include "tensorflow/core/graph/types.h"
#include "tensorflow/core/platform/strcat.h"
#include "tensorflow/core/platform/types.h"
#include "tensorflow/core/public/version.h"
#include "tensorflow/core/tpu/virtual_device.h"
namespace tensorflow {
namespace tfrt_stub {
namespace {
string DeviceName(absl::string_view name_prefix, absl::string_view device_type,
                  int32_t task_id, size_t device_id) {
  return strings::StrCat(absl::StripSuffix(name_prefix, "0"), task_id,
                         "/device:", device_type, ":", device_id);
}
DeviceAttributes BuildDeviceAttributes(absl::string_view name_prefix,
                                       const char *device_type, int32_t task_id,
                                       size_t device_id) {
  const DeviceAttributes attrs = Device::BuildDeviceAttributes(
      DeviceName(name_prefix, device_type, task_id, device_id),
      DeviceType(device_type), Bytes(16ULL << 30), DeviceLocality(),
      strings::StrCat("device: ", device_type, " device"));
  return attrs;
}
}  
absl::StatusOr<std::unique_ptr<FallbackState>> FallbackState::Create(
    const SessionOptions &session_options,
    const tensorflow::FunctionDefLibrary &fdef_lib) {
  std::vector<std::unique_ptr<Device>> devices;
  TF_RETURN_IF_ERROR(DeviceFactory::AddDevices(
      session_options, "/job:localhost/replica:0/task:0", &devices));
  return std::make_unique<FallbackState>(session_options, std::move(devices),
                                         fdef_lib);
}
absl::StatusOr<std::unique_ptr<FallbackState>>
FallbackState::CreateWithCpuDevice(
    const SessionOptions &session_options,
    const tensorflow::FunctionDefLibrary &fdef_lib) {
  std::vector<std::unique_ptr<Device>> devices;
  TF_RETURN_IF_ERROR(DeviceFactory::AddCpuDevices(
      session_options, "/job:localhost/replica:0/task:0", &devices));
  return std::make_unique<FallbackState>(session_options, std::move(devices),
                                         fdef_lib);
}
absl::StatusOr<std::unique_ptr<FallbackState>>
FallbackState::CreateWithMockGpuDevice(
    const SessionOptions &session_options,
    const tensorflow::FunctionDefLibrary &fdef_lib) {
  std::vector<std::unique_ptr<Device>> devices;
  TF_RETURN_IF_ERROR(DeviceFactory::AddCpuDevices(
      session_options, "/job:localhost/replica:0/task:0", &devices));
  auto device_attrs =
      BuildDeviceAttributes("/job:localhost/replica:0/task:0", "GPU", 0, 0);
  devices.push_back(
      std::make_unique<VirtualDevice>(session_options.env, device_attrs));
  return std::make_unique<FallbackState>(session_options, std::move(devices),
                                         fdef_lib);
}
FallbackState::FallbackState(const SessionOptions &session_options,
                             std::vector<std::unique_ptr<Device>> devices,
                             const tensorflow::FunctionDefLibrary &fdef_lib)
    : session_options_(session_options),
      device_manager_(std::move(devices)),
      func_lib_def_(OpRegistry::Global(), fdef_lib),
      pflr_(&device_manager_, session_options.env, &session_options.config,
            TF_GRAPH_DEF_VERSION, &func_lib_def_,
            session_options.config.graph_options().optimizer_options(),
            nullptr, nullptr,
            nullptr,
            Rendezvous::Factory{[](const int64, const DeviceMgr *device_mgr,
                                   tsl::core::RefCountPtr<Rendezvous> *r) {
              *r = tsl::core::RefCountPtr<Rendezvous>(
                  new IntraProcessRendezvous(device_mgr));
              return absl::OkStatus();
            }}) {
  for (auto *d : device_manager_.ListDevices()) {
    device_set_.AddDevice(d);
  }
  device_set_.set_client_device(device_manager_.HostCPU());
}
absl::StatusOr<std::unique_ptr<GraphExecutionState>>
FallbackState::CreateGraphExecutionState(GraphDef graph_def,
                                         bool run_placer) const {
  GraphExecutionStateOptions options;
  options.device_set = &device_set_;
  options.session_options = &session_options_;
  options.session_handle = "tfrt_fallback_handle";
  options.run_placer = run_placer;
  std::unique_ptr<GraphExecutionState> execution_state;
  TF_RETURN_IF_ERROR(GraphExecutionState::MakeForBaseGraph(
      std::move(graph_def), options, &execution_state));
  return execution_state;
}
Status FallbackState::AddFunctionDef(const FunctionDef &func_def) {
  return func_lib_def_.AddFunctionDef(func_def);
}
}  
}  