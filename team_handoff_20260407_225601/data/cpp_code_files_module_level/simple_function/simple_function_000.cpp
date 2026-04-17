#ifndef TENSORFLOW_LITE_EXPERIMENTAL_ACCELERATION_MINI_BENCHMARK_MINI_BENCHMARK_H_
#define TENSORFLOW_LITE_EXPERIMENTAL_ACCELERATION_MINI_BENCHMARK_MINI_BENCHMARK_H_
#include <functional>
#include <memory>
#include <string>
#include <unordered_map>
#include <vector>
#include "absl/synchronization/mutex.h"
#include "tensorflow/lite/acceleration/configuration/configuration_generated.h"
namespace tflite {
namespace acceleration {
class MiniBenchmark {
 public:
  virtual ComputeSettingsT GetBestAcceleration() = 0;
  virtual void TriggerMiniBenchmark() = 0;
  virtual void SetEventTimeoutForTesting(int64_t timeout_us) = 0;
  virtual std::vector<MiniBenchmarkEventT> MarkAndGetEventsToLog() = 0;
  virtual int NumRemainingAccelerationTests() = 0;
  MiniBenchmark() {}
  virtual ~MiniBenchmark() {}
  MiniBenchmark(MiniBenchmark&) = delete;
  MiniBenchmark& operator=(const MiniBenchmark&) = delete;
  MiniBenchmark(MiniBenchmark&&) = delete;
  MiniBenchmark& operator=(const MiniBenchmark&&) = delete;
};
std::unique_ptr<MiniBenchmark> CreateMiniBenchmark(
    const MinibenchmarkSettings& settings, const std::string& model_namespace,
    const std::string& model_id);
class MinibenchmarkImplementationRegistry {
 public:
  using CreatorFunction = std::function<std::unique_ptr<MiniBenchmark>(
      const MinibenchmarkSettings& ,
      const std::string& , const std::string& )>;
  static std::unique_ptr<MiniBenchmark> CreateByName(
      const std::string& name, const MinibenchmarkSettings& settings,
      const std::string& model_namespace, const std::string& model_id);
  struct Register {
    Register(const std::string& name, CreatorFunction creator_function);
  };
 private:
  void RegisterImpl(const std::string& name, CreatorFunction creator_function);
  std::unique_ptr<MiniBenchmark> CreateImpl(
      const std::string& name, const MinibenchmarkSettings& settings,
      const std::string& model_namespace, const std::string& model_id);
  static MinibenchmarkImplementationRegistry* GetSingleton();
  absl::Mutex mutex_;
  std::unordered_map<std::string, CreatorFunction> factories_
      ABSL_GUARDED_BY(mutex_);
};
}  
}  
#define TFLITE_REGISTER_MINI_BENCHMARK_FACTORY_FUNCTION(name, f) \
  static auto* g_tflite_mini_benchmark_##name##_ =               \
      new MinibenchmarkImplementationRegistry::Register(#name, f);
#endif  
#include "tensorflow/lite/experimental/acceleration/mini_benchmark/mini_benchmark.h"
#include <string>
#include <utility>
#include "absl/status/statusor.h"
#include "flatbuffers/flatbuffers.h"  
namespace tflite {
namespace acceleration {
namespace {
class NoopMiniBenchmark : public MiniBenchmark {
 public:
  ComputeSettingsT GetBestAcceleration() override { return ComputeSettingsT(); }
  void TriggerMiniBenchmark() override {}
  void SetEventTimeoutForTesting(int64_t) override {}
  std::vector<MiniBenchmarkEventT> MarkAndGetEventsToLog() override {
    return {};
  }
  int NumRemainingAccelerationTests() override { return -1; }
};
}  
std::unique_ptr<MiniBenchmark> CreateMiniBenchmark(
    const MinibenchmarkSettings& settings, const std::string& model_namespace,
    const std::string& model_id) {
  absl::StatusOr<std::unique_ptr<MiniBenchmark>> s_or_mb =
      MinibenchmarkImplementationRegistry::CreateByName(
          "Impl", settings, model_namespace, model_id);
  if (!s_or_mb.ok()) {
    return std::unique_ptr<MiniBenchmark>(new NoopMiniBenchmark());
  } else {
    return std::move(*s_or_mb);
  }
}
void MinibenchmarkImplementationRegistry::RegisterImpl(
    const std::string& name, CreatorFunction creator_function) {
  absl::MutexLock lock(&mutex_);
  factories_[name] = creator_function;
}
std::unique_ptr<MiniBenchmark> MinibenchmarkImplementationRegistry::CreateImpl(
    const std::string& name, const MinibenchmarkSettings& settings,
    const std::string& model_namespace, const std::string& model_id) {
  absl::MutexLock lock(&mutex_);
  auto it = factories_.find(name);
  return (it != factories_.end())
             ? it->second(settings, model_namespace, model_id)
             : nullptr;
}
MinibenchmarkImplementationRegistry*
MinibenchmarkImplementationRegistry::GetSingleton() {
  static auto* instance = new MinibenchmarkImplementationRegistry();
  return instance;
}
std::unique_ptr<MiniBenchmark>
MinibenchmarkImplementationRegistry::CreateByName(
    const std::string& name, const MinibenchmarkSettings& settings,
    const std::string& model_namespace, const std::string& model_id) {
  auto* const instance = MinibenchmarkImplementationRegistry::GetSingleton();
  return instance->CreateImpl(name, settings, model_namespace, model_id);
}
MinibenchmarkImplementationRegistry::Register::Register(
    const std::string& name, CreatorFunction creator_function) {
  auto* const instance = MinibenchmarkImplementationRegistry::GetSingleton();
  instance->RegisterImpl(name, creator_function);
}
}  
}  