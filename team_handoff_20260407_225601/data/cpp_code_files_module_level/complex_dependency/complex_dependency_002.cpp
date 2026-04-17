#ifndef TENSORFLOW_CORE_DATA_SERVICE_GRAPH_REWRITERS_H_
#define TENSORFLOW_CORE_DATA_SERVICE_GRAPH_REWRITERS_H_
#include <cstdint>
#include <string>
#include <vector>
#include "absl/strings/string_view.h"
#include "tensorflow/core/data/service/common.pb.h"
#include "tensorflow/core/framework/dataset_options.pb.h"
#include "tensorflow/core/framework/graph.pb.h"
#include "tensorflow/core/grappler/optimizers/custom_graph_optimizer.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/statusor.h"
#include "tensorflow/core/protobuf/rewriter_config.pb.h"
namespace tensorflow {
namespace data {
class RemoveCompressionMapRewriter {
 public:
  absl::StatusOr<GraphDef> ApplyRemoveCompressionMapRewrite(
      const GraphDef& graph_def);
 private:
  tensorflow::RewriterConfig::CustomGraphOptimizer GetRewriteConfig() const;
};
class AutoShardRewriter {
 public:
  static absl::StatusOr<AutoShardRewriter> Create(const TaskDef& task_def);
  absl::StatusOr<GraphDef> ApplyAutoShardRewrite(const GraphDef& graph_def);
 private:
  AutoShardRewriter(AutoShardPolicy auto_shard_policy, int64_t num_workers,
                    int64_t worker_index);
  tensorflow::RewriterConfig::CustomGraphOptimizer GetRewriteConfig() const;
  const AutoShardPolicy auto_shard_policy_;
  const int64_t num_workers_;
  const int64_t worker_index_;
};
class WorkerIndexResolver {
 public:
  template <class T>
  explicit WorkerIndexResolver(const T& worker_addresses)
      : worker_addresses_(worker_addresses.cbegin(), worker_addresses.cend()) {}
  Status ValidateWorker(absl::string_view worker_address) const;
  void AddWorker(absl::string_view worker_address);
  absl::StatusOr<int64_t> GetWorkerIndex(
      absl::string_view worker_address) const;
 private:
  std::vector<std::string> worker_addresses_;
};
}  
}  
#endif  
#include "tensorflow/core/data/service/graph_rewriters.h"
#include <cstdlib>
#include <iterator>
#include <memory>
#include <string>
#include <unordered_map>
#include <utility>
#include "absl/algorithm/container.h"
#include "absl/strings/match.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "absl/strings/substitute.h"
#include "absl/types/optional.h"
#include "tensorflow/core/data/rewrite_utils.h"
#include "tensorflow/core/data/service/common.h"
#include "tensorflow/core/data/service/common.pb.h"
#include "tensorflow/core/data/service/url.h"
#include "tensorflow/core/framework/dataset_options.pb.h"
#include "tensorflow/core/framework/graph.pb.h"
#include "tensorflow/core/framework/node_def.pb.h"
#include "tensorflow/core/framework/types.pb.h"
#include "tensorflow/core/grappler/clusters/virtual_cluster.h"
#include "tensorflow/core/grappler/grappler_item.h"
#include "tensorflow/core/grappler/grappler_item_builder.h"
#include "tensorflow/core/grappler/optimizers/custom_graph_optimizer.h"
#include "tensorflow/core/grappler/optimizers/data/auto_shard.h"
#include "tensorflow/core/grappler/optimizers/data/graph_utils.h"
#include "tensorflow/core/grappler/optimizers/data/optimizer_base.h"
#include "tensorflow/core/grappler/optimizers/data/remove_compression_map.h"
#include "tensorflow/core/kernels/data/experimental/auto_shard_dataset_op.h"
#include "tensorflow/core/platform/errors.h"
#include "tensorflow/core/platform/statusor.h"
#include "tensorflow/core/protobuf/data_service.pb.h"
#include "tensorflow/core/protobuf/device_properties.pb.h"
#include "tensorflow/core/protobuf/meta_graph.pb.h"
#include "tensorflow/core/protobuf/rewriter_config.pb.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/statusor.h"
namespace tensorflow {
namespace data {
namespace {
using ::tensorflow::data::experimental::AutoShardDatasetOp;
constexpr bool kApplyGeneralGrapplerOptimizations = false;
bool HasDynamicPort(absl::string_view address) {
  URL url(address);
  return url.has_port() && absl::StartsWith(url.port(), "%port") &&
         absl::EndsWith(url.port(), "%");
}
bool ShouldReplaceDynamicPort(absl::string_view config_address,
                              absl::string_view worker_address) {
  URL config_url(config_address), worker_url(worker_address);
  return (!config_url.has_port() || HasDynamicPort(config_address)) &&
         worker_url.has_port() && config_url.host() == worker_url.host();
}
}  
absl::StatusOr<GraphDef>
RemoveCompressionMapRewriter::ApplyRemoveCompressionMapRewrite(
    const GraphDef& graph_def) {
  grappler::RemoveCompressionMap remove_compression_map;
  tensorflow::RewriterConfig::CustomGraphOptimizer config = GetRewriteConfig();
  TF_RETURN_IF_ERROR(remove_compression_map.Init(&config));
  GraphDef input_graph = graph_def;
  TF_ASSIGN_OR_RETURN(std::string dataset_node, GetDatasetNode(input_graph));
  std::unique_ptr<tensorflow::grappler::GrapplerItem> grappler_item =
      GetGrapplerItem(&input_graph, &dataset_node, false,
                      kApplyGeneralGrapplerOptimizations);
  GraphDef rewritten_graph;
  std::unordered_map<std::string, tensorflow::DeviceProperties> device_map;
  tensorflow::grappler::VirtualCluster cluster(device_map);
  grappler::AutoShard::OptimizationStats stats;
  TF_RETURN_IF_ERROR(remove_compression_map.OptimizeAndCollectStats(
      &cluster, *grappler_item, &rewritten_graph, &stats));
  return rewritten_graph;
}
tensorflow::RewriterConfig::CustomGraphOptimizer
RemoveCompressionMapRewriter::GetRewriteConfig() const {
  tensorflow::RewriterConfig::CustomGraphOptimizer config;
  config.set_name("tf-data-service-remove-compression-map");
  return config;
}
absl::StatusOr<AutoShardRewriter> AutoShardRewriter::Create(
    const TaskDef& task_def) {
  TF_ASSIGN_OR_RETURN(
      AutoShardPolicy auto_shard_policy,
      ToAutoShardPolicy(task_def.processing_mode_def().sharding_policy()));
  return AutoShardRewriter(auto_shard_policy, task_def.num_workers(),
                           task_def.worker_index());
}
absl::StatusOr<GraphDef> AutoShardRewriter::ApplyAutoShardRewrite(
    const GraphDef& graph_def) {
  if (auto_shard_policy_ == AutoShardPolicy::OFF) {
    return graph_def;
  }
  VLOG(2) << "Applying auto-shard policy "
          << AutoShardPolicy_Name(auto_shard_policy_)
          << ". Number of workers: " << num_workers_
          << "; worker index: " << worker_index_ << ".";
  grappler::AutoShard autoshard;
  tensorflow::RewriterConfig::CustomGraphOptimizer config = GetRewriteConfig();
  TF_RETURN_IF_ERROR(autoshard.Init(&config));
  GraphDef input_graph = graph_def;
  TF_ASSIGN_OR_RETURN(std::string dataset_node, GetDatasetNode(input_graph));
  std::unique_ptr<tensorflow::grappler::GrapplerItem> grappler_item =
      GetGrapplerItem(&input_graph, &dataset_node, false,
                      kApplyGeneralGrapplerOptimizations);
  GraphDef rewritten_graph;
  std::unordered_map<std::string, tensorflow::DeviceProperties> device_map;
  tensorflow::grappler::VirtualCluster cluster(device_map);
  grappler::AutoShard::OptimizationStats stats;
  TF_RETURN_IF_ERROR(autoshard.OptimizeAndCollectStats(
      &cluster, *grappler_item, &rewritten_graph, &stats));
  return rewritten_graph;
}
AutoShardRewriter::AutoShardRewriter(AutoShardPolicy auto_shard_policy,
                                     int64_t num_workers, int64_t worker_index)
    : auto_shard_policy_(auto_shard_policy),
      num_workers_(num_workers),
      worker_index_(worker_index) {}
tensorflow::RewriterConfig::CustomGraphOptimizer
AutoShardRewriter::GetRewriteConfig() const {
  tensorflow::RewriterConfig::CustomGraphOptimizer config;
  config.set_name("tf-data-service-auto-shard");
  (*config.mutable_parameter_map())[AutoShardDatasetOp::kNumWorkers].set_i(
      num_workers_);
  (*config.mutable_parameter_map())[AutoShardDatasetOp::kIndex].set_i(
      worker_index_);
  (*config.mutable_parameter_map())[AutoShardDatasetOp::kAutoShardPolicy].set_i(
      auto_shard_policy_);
  (*config.mutable_parameter_map())[AutoShardDatasetOp::kNumReplicas].set_i(1);
  return config;
}
Status WorkerIndexResolver::ValidateWorker(
    absl::string_view worker_address) const {
  if (worker_addresses_.empty()) {
    return absl::OkStatus();
  }
  for (absl::string_view config_address : worker_addresses_) {
    if (config_address == worker_address ||
        ShouldReplaceDynamicPort(config_address, worker_address)) {
      return absl::OkStatus();
    }
  }
  return errors::FailedPrecondition(absl::Substitute(
      "Failed to assign an index for worker $0. Configured workers list: [$1]. "
      "The worker's address is not configured, or other workers are already "
      "running at the configured host. If your worker has restarted, make sure "
      "it runs at the same address and port.",
      worker_address, absl::StrJoin(worker_addresses_, ", ")));
}
void WorkerIndexResolver::AddWorker(absl::string_view worker_address) {
  for (std::string& config_address : worker_addresses_) {
    if (config_address == worker_address) {
      return;
    }
    if (ShouldReplaceDynamicPort(config_address, worker_address)) {
      config_address = std::string(worker_address);
      return;
    }
  }
}
absl::StatusOr<int64_t> WorkerIndexResolver::GetWorkerIndex(
    absl::string_view worker_address) const {
  const auto it = absl::c_find(worker_addresses_, worker_address);
  if (it == worker_addresses_.cend()) {
    return errors::NotFound(absl::Substitute(
        "Failed to shard dataset in tf.data service: Worker $0 is not in the "
        "workers list. Got workers list $1.",
        worker_address, absl::StrJoin(worker_addresses_, ",")));
  }
  return std::distance(worker_addresses_.cbegin(), it);
}
}  
}  