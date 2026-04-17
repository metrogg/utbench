#ifndef TENSORFLOW_CORE_TFRT_RUNTIME_TF_THREADPOOL_CONCURRENT_WORK_QUEUE_H_
#define TENSORFLOW_CORE_TFRT_RUNTIME_TF_THREADPOOL_CONCURRENT_WORK_QUEUE_H_
#include <memory>
#include <optional>
#include <string>
#include "tensorflow/core/platform/cpu_info.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/threadpool_interface.h"
#include "tensorflow/core/tfrt/runtime/work_queue_interface.h"
#include "tfrt/host_context/async_value.h"  
#include "tfrt/host_context/concurrent_work_queue.h"  
#include "tfrt/host_context/execution_context.h"  
#include "tfrt/host_context/task_function.h"  
#include "tfrt/support/forward_decls.h"  
namespace tensorflow {
namespace tfrt_stub {
class TfThreadPoolWorkQueue : public WorkQueueInterface {
 public:
  TfThreadPoolWorkQueue(
      tensorflow::thread::ThreadPoolInterface* intra_op_threadpool,
      tensorflow::thread::ThreadPoolInterface* inter_op_threadpool)
      : TfThreadPoolWorkQueue(0, intra_op_threadpool,
                              inter_op_threadpool) {}
  TfThreadPoolWorkQueue(
      int64_t id, tensorflow::thread::ThreadPoolInterface* intra_op_threadpool,
      tensorflow::thread::ThreadPoolInterface* inter_op_threadpool)
      : WorkQueueInterface(id, intra_op_threadpool),
        intra_op_threadpool_(intra_op_threadpool),
        inter_op_threadpool_(inter_op_threadpool) {}
  absl::StatusOr<std::unique_ptr<WorkQueueInterface>> InitializeRequest(
      int64_t request_id) const override;
  int GetParallelismLevel() const override {
    return inter_op_threadpool_->NumThreads();
  }
  std::string name() const override { return "TfThreadPoolWorkQueue"; }
  void AddTask(tfrt::TaskFunction work) override;
  std::optional<tfrt::TaskFunction> AddBlockingTask(
      tfrt::TaskFunction work, bool allow_queuing) override;
  ABSL_DEPRECATED("Use the destructor instead.")
  void Quiesce() override;
  void Await(
      tfrt::ArrayRef<::tfrt::RCReference<::tfrt::AsyncValue>> values) override;
  bool IsInWorkerThread() const override;
 private:
  tensorflow::thread::ThreadPoolInterface* intra_op_threadpool_ = nullptr;
  tensorflow::thread::ThreadPoolInterface* inter_op_threadpool_ = nullptr;
};
std::unique_ptr<TfThreadPoolWorkQueue> CreateDefaultTfThreadPoolWorkQueue(
    int num_inter_op_threads, int num_intra_op_threads);
}  
}  
#endif  
#include "tensorflow/core/tfrt/runtime/tf_threadpool_concurrent_work_queue.h"
#include <memory>
#include <optional>
#include <utility>
#include "tensorflow/core/platform/errors.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/threadpool.h"
#include "tensorflow/core/platform/threadpool_interface.h"
#include "tensorflow/core/tfrt/utils/thread_pool.h"
#include "tfrt/host_context/async_value.h"  
#include "tfrt/host_context/execution_context.h"  
#include "tfrt/host_context/task_function.h"  
#include "tfrt/support/forward_decls.h"  
#include "tfrt/support/latch.h"  
namespace tensorflow {
namespace tfrt_stub {
using ::tensorflow::thread::ThreadPoolInterface;
absl::StatusOr<std::unique_ptr<WorkQueueInterface>>
TfThreadPoolWorkQueue::InitializeRequest(int64_t request_id) const {
  return {std::make_unique<TfThreadPoolWorkQueue>(
      request_id, intra_op_threadpool_, inter_op_threadpool_)};
}
void TfThreadPoolWorkQueue::AddTask(tfrt::TaskFunction work) {
  auto* copy = new tfrt::TaskFunction(
      tensorflow::tfrt_stub::WrapWork(id(), "inter", std::move(work)));
  inter_op_threadpool_->Schedule([copy] {
    (*copy)();
    delete copy;
  });
}
std::optional<tfrt::TaskFunction> TfThreadPoolWorkQueue::AddBlockingTask(
    tfrt::TaskFunction work, bool allow_queuing) {
  AddTask(std::move(work));
  return std::nullopt;
}
void TfThreadPoolWorkQueue::Quiesce() {
}
void TfThreadPoolWorkQueue::Await(
    tfrt::ArrayRef<tfrt::RCReference<tfrt::AsyncValue>> values) {
  tfrt::latch values_remaining(values.size());
  for (auto& value : values) {
    value->AndThen([&values_remaining]() { values_remaining.count_down(); });
  }
  values_remaining.wait();
}
bool TfThreadPoolWorkQueue::IsInWorkerThread() const {
  return true;
}
std::unique_ptr<TfThreadPoolWorkQueue> CreateDefaultTfThreadPoolWorkQueue(
    int num_inter_op_threads, int num_intra_op_threads) {
  struct ThreadPools {
    TfThreadPool inter_op_threadpool;
    TfThreadPool intra_op_threadpool;
    ThreadPools(int num_inter_op_threads, int num_intra_op_threads)
        : inter_op_threadpool("default_work_queue_inter", num_inter_op_threads),
          intra_op_threadpool("default_work_queue_intra",
                              num_intra_op_threads) {}
  };
  class Wrapper : public TfThreadPoolWorkQueue {
   public:
    explicit Wrapper(std::unique_ptr<ThreadPools> thread_pools)
        : TfThreadPoolWorkQueue(
              &thread_pools->intra_op_threadpool,
              &thread_pools->inter_op_threadpool),
          thread_pools_(std::move(thread_pools)) {}
    ~Wrapper() override = default;
   private:
    std::unique_ptr<ThreadPools> thread_pools_;
  };
  return std::make_unique<Wrapper>(std::make_unique<ThreadPools>(
      num_inter_op_threads, num_intra_op_threads));
}
}  
}  