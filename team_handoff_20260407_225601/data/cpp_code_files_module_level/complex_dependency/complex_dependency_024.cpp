#ifndef TENSORFLOW_CORE_DATA_SERVICE_SNAPSHOT_PREFETCHED_SPLIT_PROVIDER_H_
#define TENSORFLOW_CORE_DATA_SERVICE_SNAPSHOT_PREFETCHED_SPLIT_PROVIDER_H_
#include <cstddef>
#include <memory>
#include <optional>
#include <string>
#include "absl/base/thread_annotations.h"
#include "absl/container/btree_set.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/synchronization/mutex.h"
#include "tensorflow/core/framework/dataset.h"
#include "tensorflow/core/framework/tensor.h"
#include "tsl/platform/env.h"
#include "tsl/platform/threadpool.h"
namespace tensorflow {
namespace data {
class PrefetchedSplitProvider {
 public:
  explicit PrefetchedSplitProvider(
      std::unique_ptr<SplitProvider> split_provider,
      const std::string& directory, tsl::Env* env,
      size_t num_write_threads = 20, size_t buffer_size_per_thread = 5);
  virtual ~PrefetchedSplitProvider();
  PrefetchedSplitProvider(const PrefetchedSplitProvider&) = delete;
  PrefetchedSplitProvider& operator=(const PrefetchedSplitProvider&) = delete;
  absl::StatusOr<std::optional<Tensor>> GetNext(const std::string& split_path);
  absl::Status Reset();
  void Cancel();
 private:
  struct SplitAndIndex {
    Tensor split;
    size_t index = 0;
    std::string SplitPath(const std::string& directory) const {
      return tsl::io::JoinPath(directory,
                               absl::StrCat("split_", index, ".tfrecord"));
    }
    friend bool operator<(const SplitAndIndex& lhs, const SplitAndIndex& rhs) {
      return lhs.index < rhs.index;
    }
  };
  absl::Status InitDirs();
  std::unique_ptr<tsl::thread::ThreadPool> RunPrefetchThreads();
  void PrefetchLoop();
  bool ShouldPrefetchSplit() const;
  absl::StatusOr<bool> PrefetchSplit();
  absl::StatusOr<std::optional<SplitAndIndex>> GetSplitFromProvider();
  void UpdateStatus(absl::Status status);
  tsl::Env* const env_;
  const std::string directory_;
  const size_t num_write_threads_;
  const size_t buffer_size_;
  mutable absl::Mutex mu_;
  mutable absl::CondVar ready_to_push_;
  mutable absl::CondVar ready_to_pop_;
  std::unique_ptr<SplitProvider> split_provider_;
  absl::Status status_ ABSL_GUARDED_BY(mu_);
  bool reset_ ABSL_GUARDED_BY(mu_) = false;
  size_t split_index_to_read_ ABSL_GUARDED_BY(mu_) = 0;
  size_t split_index_to_write_ ABSL_GUARDED_BY(mu_) = 0;
  size_t finished_threads_ ABSL_GUARDED_BY(mu_) = 0;
  absl::btree_set<SplitAndIndex> buffer_ ABSL_GUARDED_BY(mu_);
  std::unique_ptr<tsl::thread::ThreadPool> thread_pool_ ABSL_GUARDED_BY(mu_);
};
}  
}  
#endif  
#include "tensorflow/core/data/service/snapshot/prefetched_split_provider.h"
#include <cstddef>
#include <cstdint>
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include "absl/base/thread_annotations.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/synchronization/mutex.h"
#include "tensorflow/core/data/service/snapshot/file_utils.h"
#include "tensorflow/core/framework/dataset.h"
#include "tensorflow/core/framework/tensor.h"
#include "tsl/lib/io/compression.h"
#include "tsl/platform/env.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/path.h"
#include "tsl/platform/statusor.h"
#include "tsl/platform/threadpool.h"
namespace tensorflow {
namespace data {
PrefetchedSplitProvider::PrefetchedSplitProvider(
    std::unique_ptr<SplitProvider> split_provider, const std::string& directory,
    tsl::Env* env, size_t num_write_threads, size_t buffer_size_per_thread)
    : env_(env),
      directory_(directory),
      num_write_threads_(num_write_threads),
      buffer_size_(num_write_threads_ * buffer_size_per_thread),
      split_provider_(std::move(split_provider)) {
  absl::Status status = InitDirs();
  if (!status.ok()) {
    UpdateStatus(std::move(status));
    return;
  }
  absl::MutexLock l(&mu_);
  thread_pool_ = RunPrefetchThreads();
}
PrefetchedSplitProvider::~PrefetchedSplitProvider() { Cancel(); }
absl::StatusOr<std::optional<Tensor>> PrefetchedSplitProvider::GetNext(
    const std::string& split_path) ABSL_LOCKS_EXCLUDED(mu_) {
  absl::MutexLock l(&mu_);
  while (status_.ok() &&
         (buffer_.empty() || buffer_.begin()->index != split_index_to_read_) &&
         (finished_threads_ < num_write_threads_ || reset_)) {
    ready_to_pop_.Wait(&mu_);
  }
  TF_RETURN_IF_ERROR(status_);
  if (buffer_.empty()) {
    return std::nullopt;
  }
  if (buffer_.begin()->index != split_index_to_read_) {
    return absl::InternalError(absl::StrCat(
        "Failed to get tf.data snapshot split. Expected split ",
        split_index_to_read_, ", got split ", buffer_.begin()->index,
        ". This is likely a tf.data bug."));
  }
  auto it = buffer_.begin();
  SplitAndIndex split = std::move(*it);
  buffer_.erase(it);
  TF_RETURN_IF_ERROR(env_->RenameFile(split.SplitPath(directory_), split_path));
  ++split_index_to_read_;
  ready_to_push_.Signal();
  return std::move(split.split);
}
std::unique_ptr<tsl::thread::ThreadPool>
PrefetchedSplitProvider::RunPrefetchThreads() {
  auto thread_pool = std::make_unique<tsl::thread::ThreadPool>(
      env_, tsl::ThreadOptions{}, "tf_data_prefetch_splits_thread",
      num_write_threads_);
  for (size_t i = 0; i < num_write_threads_; ++i) {
    thread_pool->Schedule([this]() { PrefetchLoop(); });
  }
  return thread_pool;
}
void PrefetchedSplitProvider::PrefetchLoop() ABSL_LOCKS_EXCLUDED(mu_) {
  while (ShouldPrefetchSplit()) {
    absl::StatusOr<bool> has_next = PrefetchSplit();
    if (!has_next.status().ok()) {
      UpdateStatus(has_next.status());
      break;
    }
    if (!*has_next) {
      break;
    }
  }
  absl::MutexLock l(&mu_);
  if (++finished_threads_ >= num_write_threads_) {
    ready_to_pop_.SignalAll();
  }
}
bool PrefetchedSplitProvider::ShouldPrefetchSplit() const
    ABSL_LOCKS_EXCLUDED(mu_) {
  absl::MutexLock l(&mu_);
  return status_.ok() && !reset_;
}
absl::StatusOr<bool> PrefetchedSplitProvider::PrefetchSplit()
    ABSL_LOCKS_EXCLUDED(mu_) {
  TF_ASSIGN_OR_RETURN(std::optional<SplitAndIndex> split,
                      GetSplitFromProvider());
  if (!split.has_value()) {
    return false;
  }
  TF_RETURN_IF_ERROR(
      AtomicallyWriteTFRecords(split->SplitPath(directory_), {split->split},
                               tsl::io::compression::kNone, env_));
  absl::MutexLock l(&mu_);
  buffer_.insert(std::move(*split));
  ready_to_pop_.Signal();
  return true;
}
absl::StatusOr<std::optional<PrefetchedSplitProvider::SplitAndIndex>>
PrefetchedSplitProvider::GetSplitFromProvider() ABSL_LOCKS_EXCLUDED(mu_) {
  absl::MutexLock l(&mu_);
  while (status_.ok() && buffer_.size() >= buffer_size_ && !reset_) {
    ready_to_push_.Wait(&mu_);
  }
  TF_RETURN_IF_ERROR(status_);
  if (reset_) {
    return std::nullopt;
  }
  Tensor split;
  bool end_of_splits = false;
  TF_RETURN_IF_ERROR(split_provider_->GetNext(&split, &end_of_splits));
  if (end_of_splits) {
    return std::nullopt;
  }
  return SplitAndIndex{split, split_index_to_write_++};
}
absl::Status PrefetchedSplitProvider::Reset() ABSL_LOCKS_EXCLUDED(mu_) {
  std::unique_ptr<tsl::thread::ThreadPool> thread_pool;
  {
    absl::MutexLock l(&mu_);
    reset_ = true;
    ready_to_push_.SignalAll();
    ready_to_pop_.SignalAll();
    thread_pool = std::move(thread_pool_);
  }
  thread_pool.reset();
  TF_RETURN_IF_ERROR(split_provider_->Reset());
  absl::MutexLock l(&mu_);
  TF_RETURN_IF_ERROR(status_);
  reset_ = false;
  split_index_to_read_ = 0;
  split_index_to_write_ = 0;
  finished_threads_ = 0;
  buffer_.clear();
  TF_RETURN_IF_ERROR(InitDirs());
  thread_pool_ = RunPrefetchThreads();
  return absl::OkStatus();
}
void PrefetchedSplitProvider::Cancel() {
  UpdateStatus(
      absl::CancelledError("tf.data prefetched split provider is shut down."));
  std::unique_ptr<tsl::thread::ThreadPool> thread_pool;
  {
    absl::MutexLock l(&mu_);
    thread_pool = std::move(thread_pool_);
  }
}
absl::Status PrefetchedSplitProvider::InitDirs() {
  if (env_->FileExists(directory_).ok()) {
    int64_t undeleted_files, undeleted_dirs;
    TF_RETURN_IF_ERROR(
        env_->DeleteRecursively(directory_, &undeleted_files, &undeleted_dirs));
  }
  return env_->RecursivelyCreateDir(directory_);
}
void PrefetchedSplitProvider::UpdateStatus(absl::Status status)
    ABSL_LOCKS_EXCLUDED(mu_) {
  if (status.ok()) {
    return;
  }
  absl::MutexLock l(&mu_);
  status_.Update(std::move(status));
  ready_to_push_.SignalAll();
  ready_to_pop_.SignalAll();
}
}  
}  