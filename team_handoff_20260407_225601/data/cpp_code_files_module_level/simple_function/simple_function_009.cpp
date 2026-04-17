#ifndef TENSORSTORE_INTERNAL_THREAD_POOL_IMPL_H_
#define TENSORSTORE_INTERNAL_THREAD_POOL_IMPL_H_
#include <stddef.h>
#include <cassert>
#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_set.h"
#include "absl/synchronization/mutex.h"
#include "absl/time/time.h"
#include "tensorstore/internal/container/circular_queue.h"
#include "tensorstore/internal/intrusive_ptr.h"
#include "tensorstore/internal/thread/task_provider.h"
namespace tensorstore {
namespace internal_thread_impl {
class SharedThreadPool
    : public internal::AtomicReferenceCount<SharedThreadPool> {
 public:
  SharedThreadPool();
  void NotifyWorkAvailable(internal::IntrusivePtr<TaskProvider>)
      ABSL_LOCKS_EXCLUDED(mutex_);
 private:
  struct Overseer;
  struct Worker;
  internal::IntrusivePtr<TaskProvider> FindActiveTaskProvider()
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(mutex_);
  void StartOverseer() ABSL_EXCLUSIVE_LOCKS_REQUIRED(mutex_);
  void StartWorker(internal::IntrusivePtr<TaskProvider>, absl::Time now)
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(mutex_);
  absl::Mutex mutex_;
  size_t worker_threads_ ABSL_GUARDED_BY(mutex_) = 0;
  size_t idle_threads_ ABSL_GUARDED_BY(mutex_) = 0;
  absl::CondVar overseer_condvar_;
  bool overseer_running_ ABSL_GUARDED_BY(mutex_) = false;
  absl::Time last_thread_start_time_ ABSL_GUARDED_BY(mutex_) =
      absl::InfinitePast();
  absl::Time last_thread_exit_time_ ABSL_GUARDED_BY(mutex_) =
      absl::InfinitePast();
  absl::Time queue_assignment_time_ ABSL_GUARDED_BY(mutex_) =
      absl::InfinitePast();
  absl::flat_hash_set<TaskProvider*> in_queue_ ABSL_GUARDED_BY(mutex_);
  internal_container::CircularQueue<internal::IntrusivePtr<TaskProvider>>
      waiting_ ABSL_GUARDED_BY(mutex_);
};
}  
}  
#endif  
#include "tensorstore/internal/thread/pool_impl.h"
#include <stddef.h>
#include <stdint.h>
#include <algorithm>
#include <cassert>
#include <utility>
#include "absl/base/attributes.h"
#include "absl/base/thread_annotations.h"
#include "absl/log/absl_log.h"
#include "absl/synchronization/mutex.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "tensorstore/internal/intrusive_ptr.h"
#include "tensorstore/internal/log/verbose_flag.h"
#include "tensorstore/internal/metrics/counter.h"
#include "tensorstore/internal/metrics/gauge.h"
#include "tensorstore/internal/thread/task_provider.h"
#include "tensorstore/internal/thread/thread.h"
namespace tensorstore {
namespace internal_thread_impl {
namespace {
constexpr absl::Duration kThreadStartDelay = absl::Milliseconds(5);
constexpr absl::Duration kThreadExitDelay = absl::Milliseconds(5);
constexpr absl::Duration kThreadIdleBeforeExit = absl::Seconds(20);
constexpr absl::Duration kOverseerIdleBeforeExit = absl::Seconds(20);
auto& thread_pool_started = internal_metrics::Counter<int64_t>::New(
    "/tensorstore/thread_pool/started", "Threads started by SharedThreadPool");
auto& thread_pool_active = internal_metrics::Gauge<int64_t>::New(
    "/tensorstore/thread_pool/active",
    "Active threads managed by SharedThreadPool");
auto& thread_pool_task_providers = internal_metrics::Gauge<int64_t>::New(
    "/tensorstore/thread_pool/task_providers",
    "TaskProviders requesting threads from SharedThreadPool");
ABSL_CONST_INIT internal_log::VerboseFlag thread_pool_logging("thread_pool");
}  
SharedThreadPool::SharedThreadPool() : waiting_(128) {
  ABSL_LOG_IF(INFO, thread_pool_logging) << "SharedThreadPool: " << this;
}
void SharedThreadPool::NotifyWorkAvailable(
    internal::IntrusivePtr<TaskProvider> task_provider) {
  absl::MutexLock lock(&mutex_);
  if (in_queue_.insert(task_provider.get()).second) {
    waiting_.push_back(std::move(task_provider));
  }
  if (!overseer_running_) {
    StartOverseer();
  } else {
    overseer_condvar_.Signal();
  }
}
internal::IntrusivePtr<TaskProvider>
SharedThreadPool::FindActiveTaskProvider() {
  for (int i = waiting_.size(); i > 0; i--) {
    internal::IntrusivePtr<TaskProvider> ptr = std::move(waiting_.front());
    waiting_.pop_front();
    auto work = ptr->EstimateThreadsRequired();
    if (work == 0) {
      in_queue_.erase(ptr.get());
      continue;
    }
    if (work == 1) {
      in_queue_.erase(ptr.get());
    } else {
      waiting_.push_back(ptr);
    }
    thread_pool_task_providers.Set(waiting_.size());
    return ptr;
  }
  return nullptr;
}
struct SharedThreadPool::Overseer {
  internal::IntrusivePtr<SharedThreadPool> pool_;
  mutable absl::Time idle_start_time_;
  void operator()() const;
  void OverseerBody();
  absl::Time MaybeStartWorker(absl::Time now)
      ABSL_EXCLUSIVE_LOCKS_REQUIRED(pool_->mutex_);
};
void SharedThreadPool::StartOverseer() {
  assert(!overseer_running_);
  overseer_running_ = true;
  tensorstore::internal::Thread::StartDetached(
      {"ts_pool_overseer"},
      SharedThreadPool::Overseer{
          internal::IntrusivePtr<SharedThreadPool>(this)});
}
void SharedThreadPool::Overseer::operator()() const {
  const_cast<SharedThreadPool::Overseer*>(this)->OverseerBody();
}
void SharedThreadPool::Overseer::OverseerBody() {
  ABSL_LOG_IF(INFO, thread_pool_logging.Level(1)) << "Overseer: " << this;
  absl::Time now = absl::Now();
  idle_start_time_ = now;
  absl::Time deadline = absl::InfinitePast();
  absl::MutexLock lock(&pool_->mutex_);
  while (true) {
    pool_->overseer_condvar_.WaitWithDeadline(&pool_->mutex_, deadline);
    now = absl::Now();
    deadline = MaybeStartWorker(now);
    if (deadline < now) break;
  }
  ABSL_LOG_IF(INFO, thread_pool_logging.Level(1)) << "~Overseer: " << this;
  pool_->overseer_running_ = false;
}
absl::Time SharedThreadPool::Overseer::MaybeStartWorker(absl::Time now) {
  if (pool_->idle_threads_ || pool_->waiting_.empty()) {
    return idle_start_time_ + kOverseerIdleBeforeExit;
  }
  if (now < pool_->last_thread_start_time_ + kThreadStartDelay) {
    return pool_->last_thread_start_time_ + kThreadStartDelay;
  }
  if (now < pool_->queue_assignment_time_ + kThreadStartDelay) {
    return pool_->queue_assignment_time_ + kThreadStartDelay;
  }
  auto task_provider = pool_->FindActiveTaskProvider();
  if (!task_provider) {
    return idle_start_time_ + kOverseerIdleBeforeExit;
  }
  pool_->StartWorker(std::move(task_provider), now);
  idle_start_time_ = now;
  return now + kThreadStartDelay;
}
struct SharedThreadPool::Worker {
  internal::IntrusivePtr<SharedThreadPool> pool_;
  internal::IntrusivePtr<TaskProvider> task_provider_;
  void operator()() const;
  void WorkerBody();
};
void SharedThreadPool::StartWorker(
    internal::IntrusivePtr<TaskProvider> task_provider, absl::Time now) {
  last_thread_start_time_ = now;
  worker_threads_++;
  thread_pool_started.Increment();
  tensorstore::internal::Thread::StartDetached(
      {"ts_pool_worker"}, Worker{internal::IntrusivePtr<SharedThreadPool>(this),
                                 std::move(task_provider)});
}
void SharedThreadPool::Worker::operator()() const {
  const_cast<SharedThreadPool::Worker*>(this)->WorkerBody();
}
void SharedThreadPool::Worker::WorkerBody() {
  struct ScopedIncDec {
    size_t& x_;
    ScopedIncDec(size_t& x) : x_(x) { x_++; }
    ~ScopedIncDec() { x_--; }
  };
  thread_pool_active.Increment();
  ABSL_LOG_IF(INFO, thread_pool_logging.Level(1)) << "Worker: " << this;
  while (true) {
    if (task_provider_) {
      task_provider_->DoWorkOnThread();
      task_provider_ = nullptr;
    }
    ABSL_LOG_IF(INFO, thread_pool_logging.Level(1)) << "Idle: " << this;
    absl::Time now = absl::Now();
    absl::Time deadline = now + kThreadIdleBeforeExit;
    {
      absl::MutexLock lock(&pool_->mutex_);
      ScopedIncDec idle(pool_->idle_threads_);
      while (!task_provider_) {
        bool active = pool_->mutex_.AwaitWithDeadline(
            absl::Condition(
                +[](SharedThreadPool* self) ABSL_EXCLUSIVE_LOCKS_REQUIRED(
                     self->mutex_) { return !self->waiting_.empty(); },
                pool_.get()),
            deadline);
        now = absl::Now();
        if (active) {
          task_provider_ = pool_->FindActiveTaskProvider();
        } else {
          deadline = std::max(deadline,
                              pool_->last_thread_exit_time_ + kThreadExitDelay);
          if (deadline < now) {
            break;
          }
        }
      }
      if (task_provider_) {
        pool_->queue_assignment_time_ = now;
      } else {
        pool_->worker_threads_--;
        pool_->last_thread_exit_time_ = now;
        break;
      }
    }
  }
  thread_pool_active.Decrement();
  ABSL_LOG_IF(INFO, thread_pool_logging.Level(1)) << "~Worker: " << this;
}
}  
}  