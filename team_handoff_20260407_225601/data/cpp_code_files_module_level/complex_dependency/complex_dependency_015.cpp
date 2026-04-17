#ifndef TENSORFLOW_CORE_PROFILER_CONVERT_REPOSITORY_H_
#define TENSORFLOW_CORE_PROFILER_CONVERT_REPOSITORY_H_
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>
#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "tensorflow/core/platform/env.h"
#include "tensorflow/core/platform/path.h"
#include "tensorflow/core/platform/statusor.h"
#include "tsl/platform/env.h"
#include "tsl/platform/statusor.h"
#include "tsl/profiler/protobuf/xplane.pb.h"
#include "tsl/profiler/utils/file_system_utils.h"
namespace tensorflow {
namespace profiler {
constexpr char kAllHostsIdentifier[] = "ALL_HOSTS";
constexpr char kNoHostIdentifier[] = "NO_HOST";
enum StoredDataType {
  DCN_COLLECTIVE_STATS,
};
static auto* kHostDataSuffixes =
    new std::vector<std::pair<StoredDataType, const char*>>(
        {{StoredDataType::DCN_COLLECTIVE_STATS, ".dcn_collective_stats.pb"}});
class SessionSnapshot {
 public:
  static absl::StatusOr<SessionSnapshot> Create(
      std::vector<std::string> xspace_paths,
      std::optional<std::vector<std::unique_ptr<XSpace>>> xspaces);
  size_t XSpaceSize() const { return xspace_paths_.size(); }
  absl::StatusOr<std::unique_ptr<XSpace>> GetXSpace(size_t index) const;
  absl::StatusOr<std::unique_ptr<XSpace>> GetXSpaceByName(
      absl::string_view name) const;
  std::string GetHostname(size_t index) const;
  absl::string_view GetSessionRunDir() const { return session_run_dir_; }
  bool HasAccessibleRunDir() const { return has_accessible_run_dir_; }
  std::optional<std::string> GetFilePath(absl::string_view toolname,
                                         absl::string_view host) const;
  absl::StatusOr<std::string> GetHostDataFileName(StoredDataType data_type,
                                                  std::string host) const;
  absl::StatusOr<std::optional<std::string>> GetHostDataFilePath(
      StoredDataType data_type, std::string host) const;
  absl::StatusOr<std::pair<bool, std::string>> HasCacheFile(
      StoredDataType data_type) const;
  template <typename T>
  absl::Status WriteBinaryProto(const StoredDataType data_type,
                                const std::string host, T& proto) const {
    TF_ASSIGN_OR_RETURN(std::string filename,
                        GetHostDataFileName(data_type, host));
    std::string filepath =
        tsl::profiler::ProfilerJoinPath(GetSessionRunDir(), filename);
    return tensorflow::WriteBinaryProto(tsl::Env::Default(), filepath, proto);
  }
  template <typename T>
  absl::Status ReadBinaryProto(const StoredDataType data_type,
                               const std::string host, T* proto) const {
    TF_ASSIGN_OR_RETURN(std::optional<std::string> filepath,
                        GetHostDataFilePath(data_type, host));
    if (filepath) {
      return tensorflow::ReadBinaryProto(tsl::Env::Default(), filepath.value(),
                                         proto);
    }
    return absl::NotFoundError(
        absl::StrCat("No binary proto found for ", host, " and ", data_type));
  }
 private:
  SessionSnapshot(std::vector<std::string> xspace_paths,
                  std::optional<std::vector<std::unique_ptr<XSpace>>> xspaces)
      : xspace_paths_(std::move(xspace_paths)),
        has_accessible_run_dir_(!xspaces.has_value()),
        xspaces_(std::move(xspaces)) {
    session_run_dir_ = tensorflow::io::Dirname(xspace_paths_.at(0));
    for (size_t i = 0; i < xspace_paths_.size(); ++i) {
      std::string host_name = GetHostname(i);
      hostname_map_[host_name] = i;
    }
  }
  std::vector<std::string> xspace_paths_;
  absl::string_view session_run_dir_;
  absl::flat_hash_map<std::string , size_t >
      hostname_map_;
  const bool has_accessible_run_dir_;
  mutable std::optional<std::vector<std::unique_ptr<XSpace>>> xspaces_;
};
template <typename T>
absl::Status WriteBinaryProto(const SessionSnapshot& session_snapshot,
                              const StoredDataType data_type,
                              const std::string& host, T& proto) {
  return session_snapshot.WriteBinaryProto(data_type, host, proto);
}
template <typename T>
absl::Status ReadBinaryProto(const SessionSnapshot& session_snapshot,
                             const StoredDataType data_type,
                             const std::string& host, T* proto) {
  return session_snapshot.ReadBinaryProto(data_type, host, proto);
}
}  
}  
#endif  
#include "tensorflow/core/profiler/convert/repository.h"
#include <cstdint>
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>
#include "absl/status/status.h"
#include "absl/strings/match.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/strings/strip.h"
#include "tensorflow/core/platform/env.h"
#include "tensorflow/core/platform/errors.h"
#include "tensorflow/core/platform/path.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/statusor.h"
#include "tsl/platform/errors.h"
#include "tsl/profiler/protobuf/xplane.pb.h"
#include "tsl/profiler/utils/file_system_utils.h"
namespace tensorflow {
namespace profiler {
namespace {
std::string GetHostnameByPath(absl::string_view xspace_path) {
  std::string_view file_name = tensorflow::io::Basename(xspace_path);
  absl::ConsumeSuffix(&file_name, ".xplane.pb");
  return std::string(file_name);
}
}  
absl::StatusOr<SessionSnapshot> SessionSnapshot::Create(
    std::vector<std::string> xspace_paths,
    std::optional<std::vector<std::unique_ptr<XSpace>>> xspaces) {
  if (xspace_paths.empty()) {
    return errors::InvalidArgument("Can not find XSpace path.");
  }
  if (xspaces.has_value()) {
    if (xspaces->size() != xspace_paths.size()) {
      return errors::InvalidArgument(
          "The size of the XSpace paths: ", xspace_paths.size(),
          " is not equal ",
          "to the size of the XSpace proto: ", xspaces->size());
    }
    for (size_t i = 0; i < xspace_paths.size(); ++i) {
      auto host_name = GetHostnameByPath(xspace_paths.at(i));
      if (xspaces->at(i)->hostnames_size() > 0 && !host_name.empty()) {
        if (!absl::StrContains(host_name, xspaces->at(i)->hostnames(0))) {
          return errors::InvalidArgument(
              "The hostname of xspace path and preloaded xpace don't match at "
              "index: ",
              i, ". \nThe host name of xpace path is ", host_name,
              " but the host name of preloaded xpace is ",
              xspaces->at(i)->hostnames(0), ".");
        }
      }
    }
  }
  return SessionSnapshot(std::move(xspace_paths), std::move(xspaces));
}
absl::StatusOr<std::unique_ptr<XSpace>> SessionSnapshot::GetXSpace(
    size_t index) const {
  if (index >= xspace_paths_.size()) {
    return errors::InvalidArgument("Can not get the ", index,
                                   "th XSpace. The total number of XSpace is ",
                                   xspace_paths_.size());
  }
  if (xspaces_.has_value()) {
    if (xspaces_->at(index) == nullptr) {
      return errors::Internal("");
    }
    return std::move(xspaces_->at(index));
  }
  auto xspace_from_file = std::make_unique<XSpace>();
  TF_RETURN_IF_ERROR(tensorflow::ReadBinaryProto(tensorflow::Env::Default(),
                                                 xspace_paths_.at(index),
                                                 xspace_from_file.get()));
  return xspace_from_file;
}
absl::StatusOr<std::unique_ptr<XSpace>> SessionSnapshot::GetXSpaceByName(
    absl::string_view name) const {
  if (auto it = hostname_map_.find(name); it != hostname_map_.end()) {
    return GetXSpace(it->second);
  }
  return errors::InvalidArgument("Can not find the XSpace by name: ", name,
                                 ". The total number of XSpace is ",
                                 xspace_paths_.size());
}
std::string SessionSnapshot::GetHostname(size_t index) const {
  return GetHostnameByPath(xspace_paths_.at(index));
}
std::optional<std::string> SessionSnapshot::GetFilePath(
    absl::string_view toolname, absl::string_view hostname) const {
  if (!has_accessible_run_dir_) return std::nullopt;
  std::string file_name = "";
  if (toolname == "trace_viewer@")
    file_name = absl::StrCat(hostname, ".", "SSTABLE");
  if (!file_name.empty())
    return tensorflow::io::JoinPath(session_run_dir_, file_name);
  return std::nullopt;
}
absl::StatusOr<std::string> SessionSnapshot::GetHostDataFileName(
    const StoredDataType data_type, const std::string host) const {
  for (const auto& format : *kHostDataSuffixes) {
    if (data_type == format.first) return absl::StrCat(host, format.second);
  }
  return absl::InternalError(&"Unknown StoredDataType: "[data_type]);
}
absl::StatusOr<std::optional<std::string>> SessionSnapshot::GetHostDataFilePath(
    const StoredDataType data_type, const std::string host) const {
  std::vector<std::string> results;
  TF_RETURN_IF_ERROR(::tsl::Env::Default()->GetChildren(
      std::string(GetSessionRunDir()), &results));
  TF_ASSIGN_OR_RETURN(std::string filename,
                      GetHostDataFileName(data_type, host));
  for (const std::string& path : results) {
    if (absl::EndsWith(path, filename)) {
      return ::tsl::profiler::ProfilerJoinPath(GetSessionRunDir(), filename);
    }
  }
  return std::nullopt;
}
absl::StatusOr<std::pair<bool, std::string>> SessionSnapshot::HasCacheFile(
    const StoredDataType data_type) const {
  std::optional<std::string> filepath;
  TF_ASSIGN_OR_RETURN(filepath,
                      GetHostDataFilePath(data_type, kNoHostIdentifier));
  if (filepath) {
    return std::pair<bool, std::string>(true, std::string());
  }
  TF_ASSIGN_OR_RETURN(filepath,
                      GetHostDataFilePath(data_type, kAllHostsIdentifier));
  if (filepath) {
    return std::pair<bool, std::string>(true, filepath.value());
  }
  return std::pair<bool, std::string>(false, std::string());
}
}  
}  