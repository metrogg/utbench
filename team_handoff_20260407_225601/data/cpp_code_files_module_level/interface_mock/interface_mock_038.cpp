#ifndef TENSORFLOW_CORE_TFRT_RUNTIME_WORK_QUEUE_INTERFACE_H_
#define TENSORFLOW_CORE_TFRT_RUNTIME_WORK_QUEUE_INTERFACE_H_
#include <cstdint>
#include <memory>
#include <string>
#include <utility>
#include "tensorflow/core/platform/context.h"
#include "tensorflow/core/platform/statusor.h"
#include "tensorflow/core/platform/threadpool_interface.h"
#include "tensorflow/core/profiler/lib/connected_traceme.h"
#include "tensorflow/core/profiler/lib/traceme_encode.h"
#include "tfrt/host_context/concurrent_work_queue.h"  
#include "tfrt/support/error_util.h"  
namespace tensorflow {
namespace tfrt_stub {
class WorkQueueInterface : public tfrt::ConcurrentWorkQueue {
 public:
  WorkQueueInterface() = default;
  explicit WorkQueueInterface(int64_t id) : id_(id) {}
  explicit WorkQueueInterface(int64_t id,
                              thread::ThreadPoolInterface* intra_op_threadpool)
      : id_(id), intra_op_threadpool_(intra_op_threadpool) {}
  ~WorkQueueInterface() override = 0;
  int64_t id() const { return id_; }
  thread::ThreadPoolInterface* GetIntraOpThreadPool() const {
    return intra_op_threadpool_;
  }
  ABSL_DEPRECATED("Create the instance directly instead.")
  virtual absl::StatusOr<std::unique_ptr<WorkQueueInterface>> InitializeRequest(
      int64_t request_id) const {
    return {nullptr};
  }
 private:
  int64_t id_ = 0;
  thread::ThreadPoolInterface* intra_op_threadpool_ = nullptr;
};
inline WorkQueueInterface::~WorkQueueInterface() = default;
std::unique_ptr<WorkQueueInterface> WrapDefaultWorkQueue(
    std::unique_ptr<tfrt::ConcurrentWorkQueue> work_queue);
std::unique_ptr<WorkQueueInterface> WrapDefaultWorkQueue(
    std::unique_ptr<tfrt::ConcurrentWorkQueue> work_queue,
    thread::ThreadPoolInterface* intra_thread_pool);
template <typename Callable>
tfrt::TaskFunction WrapWork(int64_t id, absl::string_view name,
                            Callable&& work) {
  tensorflow::Context context(tensorflow::ContextKind::kThread);
  tsl::profiler::TraceMeProducer producer(
      [&]() { return absl::StrCat("producer_", name); },
      tsl::profiler::ContextType::kTfrtExecutor);
  return tfrt::TaskFunction([traceme_id = producer.GetContextId(),
                             name = std::string(name),
                             context = std::move(context),
                             work = std::forward<Callable>(work)]() mutable {
    tsl::profiler::TraceMeConsumer consumer(
        [&]() { return absl::StrCat("consumer_", name); },
        tsl::profiler::ContextType::kTfrtExecutor, traceme_id,
        tsl::profiler::TraceMeLevel::kInfo);
    tensorflow::WithContext wc(context);
    std::forward<Callable>(work)();
  });
}
}  
}  
#endif  
#include "tensorflow/core/tfrt/runtime/work_queue_interface.h"
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include "tfrt/host_context/execution_context.h"  
namespace tensorflow {
namespace tfrt_stub {
namespace {
class DefaultWorkQueueWrapper : public WorkQueueInterface {
 public:
  explicit DefaultWorkQueueWrapper(
      std::unique_ptr<tfrt::ConcurrentWorkQueue> work_queue)
      : WorkQueueInterface(0),
        work_queue_owner_(std::move(work_queue)),
        work_queue_(work_queue_owner_.get()) {}
  DefaultWorkQueueWrapper(std::unique_ptr<tfrt::ConcurrentWorkQueue> work_queue,
                          thread::ThreadPoolInterface* intra_thread_pool)
      : WorkQueueInterface(0, intra_thread_pool),
        work_queue_owner_(std::move(work_queue)),
        work_queue_(work_queue_owner_.get()) {}
  DefaultWorkQueueWrapper(int64_t request_id,
                          tfrt::ConcurrentWorkQueue* work_queue,
                          thread::ThreadPoolInterface* intra_thread_pool)
      : WorkQueueInterface(request_id, intra_thread_pool),
        work_queue_(work_queue) {}
  ~DefaultWorkQueueWrapper() override = default;
 private:
  std::string name() const override { return work_queue_->name(); }
  void AddTask(tfrt::TaskFunction work) override {
    work_queue_->AddTask(WrapWork(id(), "inter", std::move(work)));
  }
  std::optional<tfrt::TaskFunction> AddBlockingTask(
      tfrt::TaskFunction work, bool allow_queuing) override {
    return work_queue_->AddBlockingTask(
        WrapWork(id(), "blocking", std::move(work)), allow_queuing);
  }
  void Await(
      llvm::ArrayRef<tfrt::RCReference<tfrt::AsyncValue>> values) override {
    work_queue_->Await(values);
  }
  void Quiesce() override { work_queue_->Quiesce(); }
  int GetParallelismLevel() const override {
    return work_queue_->GetParallelismLevel();
  }
  bool IsInWorkerThread() const override {
    return work_queue_->IsInWorkerThread();
  }
  absl::StatusOr<std::unique_ptr<WorkQueueInterface>> InitializeRequest(
      int64_t request_id) const override {
    return {std::make_unique<DefaultWorkQueueWrapper>(request_id, work_queue_,
                                                      GetIntraOpThreadPool())};
  }
 private:
  std::unique_ptr<tfrt::ConcurrentWorkQueue> work_queue_owner_;
  tfrt::ConcurrentWorkQueue* work_queue_ = nullptr;
};
}  
std::unique_ptr<WorkQueueInterface> WrapDefaultWorkQueue(
    std::unique_ptr<tfrt::ConcurrentWorkQueue> work_queue) {
  return std::make_unique<DefaultWorkQueueWrapper>(std::move(work_queue));
}
std::unique_ptr<WorkQueueInterface> WrapDefaultWorkQueue(
    std::unique_ptr<tfrt::ConcurrentWorkQueue> work_queue,
    thread::ThreadPoolInterface* intra_thread_pool) {
  return std::make_unique<DefaultWorkQueueWrapper>(std::move(work_queue),
                                                   intra_thread_pool);
}
}  
}  