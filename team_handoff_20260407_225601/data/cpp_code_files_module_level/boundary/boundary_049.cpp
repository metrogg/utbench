#ifndef TENSORFLOW_CORE_GRAPPLER_OPTIMIZERS_INFERENCE_BATCH_OP_REWRITER_H_
#define TENSORFLOW_CORE_GRAPPLER_OPTIMIZERS_INFERENCE_BATCH_OP_REWRITER_H_
#include "tensorflow/core/grappler/grappler_item.h"
#include "tensorflow/core/grappler/optimizers/custom_graph_optimizer.h"
#include "tensorflow/core/grappler/optimizers/inference/batch_op_rewriter.pb.h"
namespace tensorflow {
namespace grappler {
constexpr char kEnableAdaptiveSchedulerAttr[] = "_enable_adaptive_scheduler";
constexpr char kMinInflightBatchesAttr[] = "_min_inflight_batches";
constexpr char kInitialInflightBatchesAttr[] = "_initial_inflight_batches";
constexpr char kMaxInflightBatchesAttr[] = "_max_inflight_batches";
constexpr char kBatchesToAverageOverAttr[] = "_batches_to_average_over";
constexpr char kFullBatchSchedulingBoostMicros[] =
    "_full_batch_scheduling_boost_micros";  
constexpr int64_t kMinInflightBatches = 16;
constexpr int64_t kInitialInflightBatches = 16;
constexpr int64_t kBatchesToAverageOver = 10;
constexpr int64_t kMaxInflightBatches = 64;
using ::tensorflow::serving::BatchOpRewriteConfig;
class BatchOpRewriter : public ::tensorflow::grappler::CustomGraphOptimizer {
 public:
  ::tensorflow::Status Init(
      const ::tensorflow::RewriterConfig_CustomGraphOptimizer* config) override;
  std::string name() const override { return "batch_op_rewriter"; }
  bool UsesFunctionLibrary() const override { return false; }
  ::tensorflow::Status Optimize(
      ::tensorflow::grappler::Cluster* cluster,
      const ::tensorflow::grappler::GrapplerItem& item,
      ::tensorflow::GraphDef* optimized_graph) override;
 private:
  BatchOpRewriteConfig config_;
};
}  
}  
#endif  
#include "tensorflow/core/grappler/optimizers/inference/batch_op_rewriter.h"
#include <functional>
#include <string>
#include "google/protobuf/wrappers.pb.h"
#include "google/protobuf/map.h"
#include "google/protobuf/repeated_field.h"
#include "absl/status/status.h"
#include "absl/strings/escaping.h"
#include "absl/strings/str_cat.h"
#include "tensorflow/core/framework/attr_value.pb.h"
#include "tensorflow/core/framework/function.pb.h"
#include "tensorflow/core/framework/node_def.pb.h"
#include "tensorflow/core/framework/tensor.pb.h"
#include "tensorflow/core/grappler/optimizers/custom_graph_optimizer_registry.h"
#include "tensorflow/core/grappler/optimizers/inference/batch_op_rewriter.pb.h"
#include "tensorflow/core/lib/core/errors.h"
#include "tensorflow/tools/graph_transforms/transform_utils.h"
namespace tensorflow {
namespace grappler {
namespace {
constexpr char kBatchFunction[] = "BatchFunction";
constexpr char kBatchOpRewriteConfigParamKey[] = "batch_op_rewrite_config";
constexpr char kNumBatchThreadsAttr[] = "num_batch_threads";
constexpr char kMaxBatchSizeAttr[] = "max_batch_size";
constexpr char kBatchTimeoutMicrosAttr[] = "batch_timeout_micros";
constexpr char kAllowedBatchSizesAttr[] = "allowed_batch_sizes";
constexpr char kMaxEnqueuedBatchesAttr[] = "max_enqueued_batches";
constexpr char kEnableLargeBatchSplitting[] = "enable_large_batch_splitting";
constexpr int64 kBoostMicrosNotSet = -1;
using BatchOpRewriteFunction = std::function<void(NodeDef* batch_op)>;
}  
using ::tensorflow::GraphDef;
using ::tensorflow::NodeDef;
using ::tensorflow::Status;
using ::tensorflow::grappler::Cluster;
using ::tensorflow::grappler::GrapplerItem;
namespace {
struct AdaptiveBatchSchedulerParams {
  int32 initial_inflight_batches;
  int32 min_inflight_batches;
  int32 max_inflight_batches;
  int32 batches_to_average_over;
  int64_t full_batch_scheduling_boost_micros;
};
AdaptiveBatchSchedulerParams GetAdaptiveBatchSchedulerParams(
    const BatchOpRewriteConfig::AdaptiveBatchSchedulerOption& option) {
  AdaptiveBatchSchedulerParams params;
  params.min_inflight_batches =
      option.has_min_inflight_batches_limit()
          ? option.min_inflight_batches_limit().value()
          : kMinInflightBatches;
  params.initial_inflight_batches =
      option.has_initial_inflight_batches_limit()
          ? option.initial_inflight_batches_limit().value()
          : kInitialInflightBatches;
  params.max_inflight_batches =
      option.has_max_inflight_batches_limit()
          ? option.max_inflight_batches_limit().value()
          : kMaxInflightBatches;
  params.batches_to_average_over =
      option.has_batches_to_average_over()
          ? option.batches_to_average_over().value()
          : kBatchesToAverageOver;
  params.full_batch_scheduling_boost_micros =
      option.has_full_batch_scheduling_boost_micros()
          ? option.full_batch_scheduling_boost_micros().value()
          : kBoostMicrosNotSet;
  return params;
}
void SetNodeAttrs(const AdaptiveBatchSchedulerParams& params, NodeDef* node) {
  ::tensorflow::graph_transforms::SetNodeAttr(kEnableAdaptiveSchedulerAttr,
                                              true, node);
  ::tensorflow::graph_transforms::SetNodeAttr(
      kMaxInflightBatchesAttr, params.max_inflight_batches, node);
  ::tensorflow::graph_transforms::SetNodeAttr(
      kMinInflightBatchesAttr, params.min_inflight_batches, node);
  ::tensorflow::graph_transforms::SetNodeAttr(
      kInitialInflightBatchesAttr, params.initial_inflight_batches, node);
  ::tensorflow::graph_transforms::SetNodeAttr(
      kBatchesToAverageOverAttr, params.batches_to_average_over, node);
  if (params.full_batch_scheduling_boost_micros != -1) {
    ::tensorflow::graph_transforms::SetNodeAttr(
        kFullBatchSchedulingBoostMicros,
        params.full_batch_scheduling_boost_micros, node);
  }
}
void UpdateBatchOps(GraphDef* graph, BatchOpRewriteFunction rewrite_fn) {
  for (int i = 0; i < graph->node_size(); ++i) {
    NodeDef* node = graph->mutable_node(i);
    if (node->op() == kBatchFunction) {
      rewrite_fn(node);
    }
  }
  for (int i = 0; i < graph->library().function_size(); i++) {
    FunctionDef* function_def = graph->mutable_library()->mutable_function(i);
    for (int j = 0; j < function_def->node_def_size(); j++) {
      NodeDef* node = function_def->mutable_node_def(j);
      if (node->op() == kBatchFunction) {
        rewrite_fn(node);
      }
    }
  }
}
}  
Status BatchOpRewriter::Init(
    const ::tensorflow::RewriterConfig_CustomGraphOptimizer* config) {
  if (config->parameter_map().find(kBatchOpRewriteConfigParamKey) ==
      config->parameter_map().end()) {
    return absl::InternalError(
        "batch_op_rewrite_config param must be set in the rewriter config "
        "with a serialized/encoded BatchOpRewriteConfig.");
  }
  const auto& params =
      config->parameter_map().at(kBatchOpRewriteConfigParamKey);
  std::string unencoded;
  if (params.s().empty()) {
    VLOG(2) << "Empty batch-op rewrite config";
    return absl::OkStatus();
  }
  if (!absl::Base64Unescape(params.s(), &unencoded)) {
    return absl::InternalError(
        "Failed to unencode batch_op_rewrite_config from params.");
  }
  if (!config_.ParseFromString(unencoded)) {
    return absl::InternalError(
        "Failed to parse batch_op_rewrite_config from params.");
  }
  VLOG(2) << "BatchOp Rewrite config is " << config_.DebugString();
  return absl::OkStatus();
}
Status BatchOpRewriter::Optimize(Cluster* cluster, const GrapplerItem& item,
                                 GraphDef* optimized_graph) {
  VLOG(2) << "Running BatchOp Rewriter";
  *optimized_graph = item.graph;
  bool asbs_overridden = false;
  if (config_proto_.has_experimental() &&
      config_proto_.experimental().has_session_metadata()) {
    const string model_name =
        config_proto_.experimental().session_metadata().name();
    if (!config_.model_scheduler_options().empty()) {
      return absl::InvalidArgumentError(
          "model_scheduler_options is deprecated. Please use the "
          "adaptive_batch_scheduler_option field in batch_options instead.");
    }
    auto model_batch_options = config_.batch_options().find(model_name);
    if (model_batch_options != config_.batch_options().end()) {
      auto& batch_options = model_batch_options->second;
      VLOG(2) << "Rewriting batch_options for " << model_name << " to "
              << batch_options.DebugString();
      if (batch_options.has_adaptive_batch_scheduler_option()) {
        AdaptiveBatchSchedulerParams params = GetAdaptiveBatchSchedulerParams(
            batch_options.adaptive_batch_scheduler_option());
        if ((params.min_inflight_batches > params.max_inflight_batches) ||
            (params.initial_inflight_batches < params.min_inflight_batches) ||
            (params.initial_inflight_batches > params.max_inflight_batches)) {
          return absl::InvalidArgumentError(absl::StrCat(
              "Requires min_inflight_batches <= initial_inflight_batches "
              "and initial_inflight_batches <= max_inflight_batches; Got "
              "{min_inflight_batches : ",
              params.min_inflight_batches,
              ", initial_inflight_batches : ", params.initial_inflight_batches,
              ", max_inflight_batches : ", params.max_inflight_batches, "}."));
        }
        asbs_overridden = true;
        UpdateBatchOps(optimized_graph, [&params](NodeDef* batch_op) {
          SetNodeAttrs(params, batch_op);
        });
      }
      if (config_.enable_adaptive_shared_batching_thread_pool() &&
          !asbs_overridden && batch_options.has_num_batch_threads() &&
          batch_options.num_batch_threads() != 0) {
        return absl::InvalidArgumentError(
            "Unable to enable adapative shared batching because it requires "
            "num_batch_threads=0 but the BatchOpRewriteConfig is also trying "
            "to set num_batch_threads. Set either set "
            "enable_adaptive_shared_batching_thread_pool or num_batch_threads "
            "but not both.");
      }
      UpdateBatchOps(optimized_graph, [&batch_options](NodeDef* batch_op) {
        if (batch_options.has_num_batch_threads()) {
          ::tensorflow::graph_transforms::SetNodeAttr(
              kNumBatchThreadsAttr, batch_options.num_batch_threads(),
              batch_op);
        }
        if (batch_options.has_max_batch_size()) {
          ::tensorflow::graph_transforms::SetNodeAttr(
              kMaxBatchSizeAttr, batch_options.max_batch_size(), batch_op);
        }
        if (batch_options.has_batch_timeout_micros()) {
          ::tensorflow::graph_transforms::SetNodeAttr(
              kBatchTimeoutMicrosAttr, batch_options.batch_timeout_micros(),
              batch_op);
        }
        if (!batch_options.allowed_batch_sizes().empty()) {
          ::tensorflow::graph_transforms::SetNodeAttr(
              kAllowedBatchSizesAttr, batch_options.allowed_batch_sizes(),
              batch_op);
        }
        if (batch_options.has_max_enqueued_batches()) {
          ::tensorflow::graph_transforms::SetNodeAttr(
              kMaxEnqueuedBatchesAttr, batch_options.max_enqueued_batches(),
              batch_op);
        }
        if (batch_options.has_disable_large_batch_splitting()) {
          ::tensorflow::graph_transforms::SetNodeAttr(
              kEnableLargeBatchSplitting,
              !batch_options.disable_large_batch_splitting(), batch_op);
        }
      });
    }
  }
  if (asbs_overridden) {
    return absl::OkStatus();
  }
  if (config_.enable_adaptive_shared_batching_thread_pool()) {
    UpdateBatchOps(optimized_graph, [](NodeDef* batch_op) {
      ::tensorflow::graph_transforms::SetNodeAttr(kNumBatchThreadsAttr, 0,
                                                  batch_op);
    });
  }
  return absl::OkStatus();
}
REGISTER_GRAPH_OPTIMIZER_AS(BatchOpRewriter, "batch_op_rewrite");
}  
}  