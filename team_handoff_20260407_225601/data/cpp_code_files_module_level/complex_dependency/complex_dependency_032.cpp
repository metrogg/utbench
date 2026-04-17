#ifndef TENSORFLOW_CORE_PROFILER_CONVERT_XPLANE_TO_DCN_COLLECTIVE_STATS_H_
#define TENSORFLOW_CORE_PROFILER_CONVERT_XPLANE_TO_DCN_COLLECTIVE_STATS_H_
#include "tensorflow/core/platform/statusor.h"
#include "tensorflow/core/profiler/convert/repository.h"
#include "tensorflow/core/profiler/protobuf/dcn_slack_analysis.pb.h"
namespace tensorflow {
namespace profiler {
absl::StatusOr<bool> ConvertMultiXSpaceToDcnCollectiveStats(
    const SessionSnapshot& session_snapshot);
absl::StatusOr<bool> HasDcnCollectiveStatsInMultiXSpace(
    const SessionSnapshot& session_snapshot);
absl::StatusOr<DcnSlackAnalysis> GetDcnSlackAnalysisByHostName(
    const SessionSnapshot& session_snapshot, std::string hostname);
}  
}  
#endif  
#include "tensorflow/core/profiler/convert/xplane_to_dcn_collective_stats.h"
#include <memory>
#include <string>
#include <utility>
#include "absl/strings/match.h"
#include "tensorflow/core/platform/statusor.h"
#include "tensorflow/core/profiler/convert/dcn_slack_analysis_combiner.h"
#include "tensorflow/core/profiler/convert/repository.h"
#include "tensorflow/core/profiler/convert/xspace_to_dcn_slack_analysis.h"
#include "tensorflow/core/profiler/protobuf/dcn_slack_analysis.pb.h"
#include "tensorflow/core/profiler/utils/xplane_schema.h"
#include "tensorflow/core/profiler/utils/xplane_utils.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/statusor.h"
#include "tsl/profiler/protobuf/xplane.pb.h"
namespace tensorflow {
namespace profiler {
namespace {
bool HasDcnCollectiveStatsInXSpace(const XSpace& xspace) {
  if (const tensorflow::profiler::XPlane* xplane = FindPlaneWithName(
          xspace, tensorflow::profiler::kHostThreadsPlaneName);
      xplane != nullptr) {
    for (const auto& [_, metadata] : xplane->event_metadata()) {
      if (absl::StartsWith(metadata.name(), "MegaScale:")) {
        return true;
      }
    }
  }
  return false;
}
absl::StatusOr<bool> GetDcnCollectiveStatsFromMultiXSpaceAndSaveToFile(
    const SessionSnapshot& session_snapshot) {
  DcnSlackAnalysisCombiner combiner;
  for (int idx = 0; idx < session_snapshot.XSpaceSize(); idx++) {
    std::string hostname = session_snapshot.GetHostname(idx);
    TF_ASSIGN_OR_RETURN(std::unique_ptr<XSpace> xspace,
                        session_snapshot.GetXSpace(idx));
    if (!HasDcnCollectiveStatsInXSpace(*xspace)) {
      DcnSlackAnalysis dcnSlackAnalysis;
      TF_RETURN_IF_ERROR(WriteBinaryProto(session_snapshot,
                                          StoredDataType::DCN_COLLECTIVE_STATS,
                                          kNoHostIdentifier, dcnSlackAnalysis));
      return false;
    }
    DcnSlackAnalysis dcnSlackAnalysis =
        ConvertXSpaceToDcnSlackAnalysis(*xspace, nullptr, nullptr);
    TF_RETURN_IF_ERROR(WriteBinaryProto(session_snapshot,
                                        StoredDataType::DCN_COLLECTIVE_STATS,
                                        hostname, dcnSlackAnalysis));
    combiner.Combine(dcnSlackAnalysis);
  }
  DcnSlackAnalysis dcnSlackAnalysis = combiner.Finalize();
  TF_RETURN_IF_ERROR(WriteBinaryProto(session_snapshot,
                                      StoredDataType::DCN_COLLECTIVE_STATS,
                                      kAllHostsIdentifier, dcnSlackAnalysis));
  return true;
}
}  
absl::StatusOr<bool> HasDcnCollectiveStatsInMultiXSpace(
    const SessionSnapshot& session_snapshot) {
  std::pair<bool, std::string> hasCacheFile;
  TF_ASSIGN_OR_RETURN(hasCacheFile, session_snapshot.HasCacheFile(
                                        StoredDataType::DCN_COLLECTIVE_STATS));
  if (!hasCacheFile.first) {
    for (int idx = 0; idx < session_snapshot.XSpaceSize(); idx++) {
      std::string hostname = session_snapshot.GetHostname(idx);
      TF_ASSIGN_OR_RETURN(std::unique_ptr<XSpace> xspace,
                          session_snapshot.GetXSpace(idx));
      if (HasDcnCollectiveStatsInXSpace(*xspace)) {
        return true;
      }
    }
    return false;
  }
  if (hasCacheFile.second.empty()) {
    return false;
  } else {
    return true;
  }
}
absl::StatusOr<bool> ConvertMultiXSpaceToDcnCollectiveStats(
    const SessionSnapshot& session_snapshot) {
  std::pair<bool, std::string> hasCacheFile;
  TF_ASSIGN_OR_RETURN(hasCacheFile, session_snapshot.HasCacheFile(
                                        StoredDataType::DCN_COLLECTIVE_STATS));
  if (!hasCacheFile.first) {
    return GetDcnCollectiveStatsFromMultiXSpaceAndSaveToFile(session_snapshot);
  }
  if (hasCacheFile.second.empty()) {
    return false;
  } else {
    return true;
  }
}
absl::StatusOr<DcnSlackAnalysis> GetDcnSlackAnalysisByHostName(
    const SessionSnapshot& session_snapshot, const std::string hostname) {
  TF_ASSIGN_OR_RETURN(bool hasDcnCollectiveStats,
                      ConvertMultiXSpaceToDcnCollectiveStats(session_snapshot));
  DcnSlackAnalysis dcnSlackAnalysis;
  if (hasDcnCollectiveStats) {
    TF_RETURN_IF_ERROR(ReadBinaryProto(session_snapshot,
                                       StoredDataType::DCN_COLLECTIVE_STATS,
                                       hostname, &dcnSlackAnalysis));
  }
  return dcnSlackAnalysis;
}
}  
}  