#ifndef TENSORFLOW_CORE_DATA_SERVICE_SNAPSHOT_SNAPSHOT_CHUNK_PROVIDER_H_
#define TENSORFLOW_CORE_DATA_SERVICE_SNAPSHOT_SNAPSHOT_CHUNK_PROVIDER_H_
#include <cstdint>
#include <functional>
#include <optional>
#include <string>
#include <utility>
#include <vector>
#include "absl/base/thread_annotations.h"
#include "absl/container/btree_set.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "tensorflow/core/framework/dataset.h"
#include "tensorflow/core/framework/tensor.h"
#include "tsl/platform/env.h"
namespace tensorflow {
namespace data {
class SnapshotChunkProvider : public SplitProvider {
 public:
  SnapshotChunkProvider(absl::string_view snapshot_path, tsl::Env* env);
  ~SnapshotChunkProvider() override = default;
  SnapshotChunkProvider(const SnapshotChunkProvider&) = delete;
  SnapshotChunkProvider& operator=(const SnapshotChunkProvider&) = delete;
  absl::Status GetNext(Tensor* split, bool* end_of_splits) override;
  absl::Status Reset() override;
  absl::Status Save(std::function<std::string(std::string)> full_name,
                    IteratorStateWriter* writer) override;
  absl::Status Restore(std::function<std::string(std::string)> full_name,
                       IteratorStateReader* reader) override;
  int64_t Cardinality() const override;
  void Cancel() override;
 private:
  struct SnapshotState {
    SnapshotState() = default;
    explicit SnapshotState(bool snapshot_is_done)
        : snapshot_is_done(snapshot_is_done) {}
    explicit SnapshotState(absl::Status status) : status(std::move(status)) {}
    bool snapshot_is_done = false;
    absl::Status status = absl::OkStatus();
  };
  struct ChunkOrder {
    bool operator()(const std::string& chunk1, const std::string& chunk2) const;
  };
  using OrderedChunkSet = absl::btree_set<std::string, ChunkOrder>;
  static std::string SetToString(const OrderedChunkSet& s);
  static OrderedChunkSet SetFromString(absl::string_view s);
  absl::Status UpdateSnapshot();
  absl::StatusOr<SnapshotState> GetSnapshotState();
  absl::StatusOr<std::vector<std::string>> GetAvailableChunks();
  const std::string snapshot_path_;
  tsl::Env* const env_;
  mutable absl::Mutex mu_;
  OrderedChunkSet chunks_read_ ABSL_GUARDED_BY(mu_);
  OrderedChunkSet chunks_unread_ ABSL_GUARDED_BY(mu_);
  SnapshotState snapshot_state_ ABSL_GUARDED_BY(mu_);
};
}  
}  
#endif  
#include "tensorflow/core/data/service/snapshot/snapshot_chunk_provider.h"
#include <cstdint>
#include <functional>
#include <optional>
#include <string>
#include <tuple>
#include <utility>
#include <vector>
#include "absl/base/thread_annotations.h"
#include "absl/container/btree_set.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_join.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "absl/time/time.h"
#include "tensorflow/core/data/service/snapshot/file_utils.h"
#include "tensorflow/core/data/service/snapshot/path_utils.h"
#include "tensorflow/core/framework/dataset.h"
#include "tensorflow/core/framework/tensor.h"
#include "tensorflow/core/framework/tensor_shape.h"
#include "tensorflow/core/framework/types.pb.h"
#include "tsl/platform/env.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/path.h"
#include "tsl/platform/retrying_utils.h"
#include "tsl/platform/status_to_from_proto.h"
#include "tsl/platform/statusor.h"
#include "tsl/platform/tstring.h"
#include "tsl/protobuf/status.pb.h"
namespace tensorflow {
namespace data {
namespace {
constexpr char kChunksRead[] = "chunks_read";
constexpr absl::string_view kSetElementDelimiter = ",";
Tensor ConvertToTensor(absl::string_view s) {
  Tensor tensor(DT_STRING, TensorShape({}));
  tensor.scalar<tsl::tstring>()() = tsl::tstring(s);
  return tensor;
}
std::string AbsPath(absl::string_view snapshot_path, absl::string_view chunk) {
  return tsl::io::JoinPath(CommittedChunksDirectory(snapshot_path), chunk);
}
void Backoff(int num_retries, tsl::Env* env) {
  if (num_retries >= 1) {  
    absl::Duration retry_backoff = tsl::ComputeRetryBackoff(num_retries - 1);
    env->SleepForMicroseconds(absl::ToInt64Microseconds(retry_backoff));
  }
}
}  
SnapshotChunkProvider::SnapshotChunkProvider(absl::string_view snapshot_path,
                                             tsl::Env* env)
    : snapshot_path_(snapshot_path), env_(env) {}
absl::Status SnapshotChunkProvider::GetNext(Tensor* split, bool* end_of_splits)
    ABSL_LOCKS_EXCLUDED(mu_) {
  for (int num_retries = 0;; ++num_retries) {
    Backoff(num_retries, env_);
    absl::MutexLock l(&mu_);
    TF_RETURN_IF_ERROR(snapshot_state_.status);
    if (!chunks_unread_.empty()) {
      std::string next_chunk = *chunks_unread_.begin();
      chunks_read_.insert(next_chunk);
      chunks_unread_.erase(next_chunk);
      *split = ConvertToTensor(AbsPath(snapshot_path_, next_chunk));
      *end_of_splits = false;
      return absl::OkStatus();
    }
    if (snapshot_state_.snapshot_is_done) {
      *end_of_splits = true;
      return absl::OkStatus();
    }
    TF_RETURN_IF_ERROR(UpdateSnapshot());
  }
}
absl::Status SnapshotChunkProvider::UpdateSnapshot()
    ABSL_EXCLUSIVE_LOCKS_REQUIRED(mu_) {
  TF_ASSIGN_OR_RETURN(snapshot_state_, GetSnapshotState());
  TF_RETURN_IF_ERROR(snapshot_state_.status);
  TF_ASSIGN_OR_RETURN(std::vector<std::string> chunks, GetAvailableChunks());
  for (const std::string& chunk : chunks) {
    if (!chunks_read_.contains(chunk)) {
      chunks_unread_.insert(std::string(chunk));
    }
  }
  return absl::OkStatus();
}
absl::StatusOr<SnapshotChunkProvider::SnapshotState>
SnapshotChunkProvider::GetSnapshotState() {
  std::string error_file_path = SnapshotErrorFilePath(snapshot_path_);
  if (env_->FileExists(error_file_path).ok()) {
    StatusProto status_proto;
    TF_RETURN_IF_ERROR(ReadTextProto(env_, error_file_path, &status_proto));
    absl::Status status = tsl::StatusFromProto(status_proto);
    if (status.ok()) {
      return absl::InternalError(absl::StrCat(
          "Unexpected snapshot ERROR file contains an OK status at ",
          error_file_path, "."));
    }
    return SnapshotState(status);
  }
  return SnapshotState(
      env_->FileExists(SnapshotDoneFilePath(snapshot_path_)).ok());
}
absl::StatusOr<std::vector<std::string>>
SnapshotChunkProvider::GetAvailableChunks() {
  absl::StatusOr<std::vector<std::string>> status_or_chunks =
      GetChildren(CommittedChunksDirectory(snapshot_path_), env_);
  if (status_or_chunks.ok()) {
    return *std::move(status_or_chunks);
  } else if (absl::IsNotFound(status_or_chunks.status())) {
    return std::vector<std::string>{};
  }
  return status_or_chunks.status();
}
absl::Status SnapshotChunkProvider::Reset() {
  absl::MutexLock l(&mu_);
  chunks_read_.clear();
  chunks_unread_.clear();
  return UpdateSnapshot();
}
absl::Status SnapshotChunkProvider::Save(
    std::function<std::string(std::string)> full_name,
    IteratorStateWriter* writer) {
  absl::MutexLock l(&mu_);
  TF_RETURN_IF_ERROR(
      writer->WriteScalar(full_name(kChunksRead), SetToString(chunks_read_)));
  return absl::OkStatus();
}
absl::Status SnapshotChunkProvider::Restore(
    std::function<std::string(std::string)> full_name,
    IteratorStateReader* reader) {
  absl::MutexLock l(&mu_);
  tsl::tstring chunks_read;
  TF_RETURN_IF_ERROR(reader->ReadScalar(full_name(kChunksRead), &chunks_read));
  chunks_read_ = SetFromString(chunks_read);
  return UpdateSnapshot();
}
int64_t SnapshotChunkProvider::Cardinality() const {
  return SnapshotChunksCardinality(snapshot_path_, env_);
}
void SnapshotChunkProvider::Cancel() {
  absl::MutexLock l(&mu_);
  if (snapshot_state_.snapshot_is_done || !snapshot_state_.status.ok()) {
    return;
  }
  snapshot_state_.status = absl::CancelledError(
      absl::StrCat("Cancelled loading tf.data snapshot at ", snapshot_path_));
  VLOG(2) << snapshot_state_.status;
}
std::string SnapshotChunkProvider::SetToString(
    const SnapshotChunkProvider::OrderedChunkSet& s) {
  return absl::StrJoin(s, kSetElementDelimiter);
}
SnapshotChunkProvider::OrderedChunkSet SnapshotChunkProvider::SetFromString(
    absl::string_view s) {
  if (s.empty()) {
    return {};
  }
  std::vector<std::string> split = absl::StrSplit(s, kSetElementDelimiter);
  return OrderedChunkSet(split.begin(), split.end());
}
bool SnapshotChunkProvider::ChunkOrder::operator()(
    const std::string& chunk1, const std::string& chunk2) const {
  absl::StatusOr<std::tuple<int64_t, int64_t, int64_t>> tokens1 =
      ParseChunkFilename(chunk1);
  absl::StatusOr<std::tuple<int64_t, int64_t, int64_t>> tokens2 =
      ParseChunkFilename(chunk2);
  if (!tokens1.status().ok()) {
    LOG_EVERY_N_SEC(ERROR, 60) << "Failed to parse tf.data snapshot chunk file "
                               << chunk1 << ": " << tokens1.status();
    return chunk1 < chunk2;
  }
  if (!tokens2.status().ok()) {
    LOG_EVERY_N_SEC(ERROR, 60) << "Failed to parse tf.data snapshot chunk file "
                               << chunk2 << ": " << tokens2.status();
    return chunk1 < chunk2;
  }
  auto [stream_index1, chunk_index1, num_records1] = *tokens1;
  auto [stream_index2, chunk_index2, num_records2] = *tokens2;
  if (chunk_index1 != chunk_index2) {
    return chunk_index1 < chunk_index2;
  }
  return stream_index1 < stream_index2;
}
}  
}  