#ifndef TENSORFLOW_CORE_DISTRIBUTED_RUNTIME_COORDINATION_COORDINATION_SERVICE_BARRIER_PROXY_H_
#define TENSORFLOW_CORE_DISTRIBUTED_RUNTIME_COORDINATION_COORDINATION_SERVICE_BARRIER_PROXY_H_
#include <cstdint>
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/container/flat_hash_map.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "xla/tsl/distributed_runtime/coordination/coordination_service_agent.h"
#include "tensorflow/core/platform/macros.h"
#include "tensorflow/core/platform/mutex.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/thread_annotations.h"
#include "tsl/protobuf/coordination_service.pb.h"
namespace tensorflow {
class BarrierProxy {
 public:
  BarrierProxy(const BarrierProxy&) = delete;
  void operator=(const BarrierProxy&) = delete;
  BarrierProxy(tsl::CoordinationServiceAgent* agent,
               std::vector<CoordinatedTask> tasks, int num_local_threads,
               absl::string_view key, absl::Duration timeout)
      : key_(key),
        agent_(agent),
        tasks_(std::move(tasks)),
        timeout_(timeout),
        num_local_threads_(num_local_threads) {}
  ~BarrierProxy() = default;
  std::pair<Status, bool> Wait();
 private:
  const std::string key_;
  tsl::CoordinationServiceAgent* agent_;
  const std::vector<CoordinatedTask> tasks_;
  absl::Duration timeout_;
  mutex mu_;
  condition_variable cv_ TF_GUARDED_BY(mu_);
  const int num_local_threads_;
  int num_entered_ TF_GUARDED_BY(mu_) = 0;
  int num_to_exit_ TF_GUARDED_BY(mu_) = 0;
  Status status_ TF_GUARDED_BY(mu_);
  bool status_set_ TF_GUARDED_BY(mu_) = false;
};
class BarrierProxyManager {
 public:
  BarrierProxyManager(const BarrierProxyManager&) = delete;
  void operator=(const BarrierProxyManager&) = delete;
  BarrierProxyManager() = default;
  ~BarrierProxyManager() = default;
  Status Wait(tsl::CoordinationServiceAgent* agent,
              const std::vector<CoordinatedTask>& tasks, int num_local_threads,
              absl::string_view key, absl::Duration timeout);
  size_t size() const;
 private:
  mutable mutex mu_;
  absl::flat_hash_map<std::string, std::shared_ptr<BarrierProxy>> barriers_
      TF_GUARDED_BY(mu_);
};
}  
#endif  
#include "tensorflow/core/distributed_runtime/coordination/coordination_service_barrier_proxy.h"
#include <memory>
#include <string>
#include <tuple>
#include <utility>
#include <vector>
#include "absl/log/check.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "xla/tsl/distributed_runtime/coordination/coordination_service_agent.h"
#include "tensorflow/core/platform/mutex.h"
#include "tensorflow/core/platform/status.h"
#include "tsl/profiler/lib/traceme.h"
#include "tsl/profiler/lib/traceme_encode.h"
#include "tsl/protobuf/coordination_service.pb.h"
namespace tensorflow {
std::pair<Status, bool> BarrierProxy::Wait() {
  mutex_lock l(mu_);
  if (status_set_) {
    return std::make_pair(
        absl::FailedPreconditionError(absl::StrCat(
            "The barrier has already passed or timed out. key=", key_)),
        false);
  }
  if (num_entered_ >= num_local_threads_) {
    return std::make_pair(absl::FailedPreconditionError(absl::StrCat(
                              "Wait() called too many (>", num_local_threads_,
                              ") times. key=", key_)),
                          false);
  }
  ++num_entered_;
  ++num_to_exit_;
  VLOG(1) << "BarrierProxy " << key_ << " enter: num_entered_=" << num_entered_
          << ", num_to_exit_=" << num_to_exit_;
  if (num_entered_ == num_local_threads_) {
    if (tasks_.size() != 1) {
      tsl::profiler::TraceMe traceme("BarrierProxy::Wait::WaitAtBarrier");
      status_ = agent_->WaitAtBarrier(key_, timeout_, tasks_);
    } else {
      status_ = absl::OkStatus();
    }
    status_set_ = true;
    cv_.notify_all();
  } else if (WaitForMilliseconds(&l, &cv_, timeout_ / absl::Milliseconds(1)) ==
             kCond_Timeout) {
    if (!status_set_) {
      if (tasks_.size() != 1) {
        agent_->CancelBarrier(key_).IgnoreError();
      }
      status_ = absl::DeadlineExceededError(
          absl::StrCat("BarrierProxy timeout: key=", key_));
      status_set_ = true;
      cv_.notify_all();
    }
  } else {
    CHECK(status_set_);  
  }
  --num_to_exit_;
  VLOG(1) << "BarrierProxy " << key_ << " enter: num_entered_=" << num_entered_
          << ", num_to_exit=" << num_to_exit_;
  return std::make_pair(status_, num_to_exit_ == 0);
}
size_t BarrierProxyManager::size() const {
  mutex_lock l(mu_);
  return barriers_.size();
}
Status BarrierProxyManager::Wait(tsl::CoordinationServiceAgent* agent,
                                 const std::vector<CoordinatedTask>& tasks,
                                 int num_local_threads, absl::string_view key,
                                 absl::Duration timeout) {
  if (tasks.size() == 1 && num_local_threads <= 1) return absl::OkStatus();
  tsl::profiler::TraceMe traceme([&] {
    return tsl::profiler::TraceMeEncode(
        "BarrierProxyManager::Wait",
        {
            {"num_tasks", tasks.size()},
            {"num_local_threads", num_local_threads},
        });
  });
  std::shared_ptr<BarrierProxy> barrier;
  {
    mutex_lock l(mu_);
    auto [iter, inserted] = barriers_.try_emplace(key);
    if (inserted) {
      iter->second = std::make_shared<BarrierProxy>(
          agent, tasks, num_local_threads, key, timeout);
      VLOG(1) << "BarrierProxy key=" << key << " created.";
    }
    barrier = iter->second;
  }
  CHECK(barrier);  
  auto [status, last_exit] = barrier->Wait();
  if (last_exit) {
    mutex_lock l(mu_);
    barriers_.erase(key);
    VLOG(1) << "BarrierProxy key=" << key << " removed.";
  }
  return status;
}
}  