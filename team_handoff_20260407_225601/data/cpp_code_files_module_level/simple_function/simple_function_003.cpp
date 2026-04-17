#ifndef TENSORFLOW_CORE_GRAPPLER_COSTS_ANALYTICAL_COST_ESTIMATOR_H_
#define TENSORFLOW_CORE_GRAPPLER_COSTS_ANALYTICAL_COST_ESTIMATOR_H_
#include "tensorflow/core/grappler/costs/cost_estimator.h"
#include "tensorflow/core/grappler/costs/op_level_cost_estimator.h"
#include "tensorflow/core/grappler/costs/virtual_scheduler.h"
#include "tensorflow/core/grappler/grappler_item.h"
#include "tensorflow/core/lib/core/status.h"
namespace tensorflow {
class CostGraphDef;
class GraphDef;
}  
namespace tensorflow {
namespace grappler {
class Cluster;
struct GrapplerItem;
class AnalyticalCostEstimator : public CostEstimator {
 public:
  AnalyticalCostEstimator(Cluster* cluster, bool use_static_shapes,
                          bool use_aggressive_shape_inference);
  AnalyticalCostEstimator(Cluster* cluster,
                          std::unique_ptr<OpLevelCostEstimator> node_estimator,
                          std::unique_ptr<ReadyNodeManager> node_manager,
                          bool use_static_shapes,
                          bool use_aggressive_shape_inference);
  AnalyticalCostEstimator(Cluster* cluster,
                          std::unique_ptr<OpLevelCostEstimator> node_estimator,
                          std::unique_ptr<ReadyNodeManager> node_manager,
                          std::unique_ptr<VirtualPlacer> placer,
                          bool use_static_shapes,
                          bool use_aggressive_shape_inference);
  ~AnalyticalCostEstimator() override {}
  Status Initialize(const GrapplerItem& item) override;
  Status PredictCosts(const GraphDef& optimized_graph,
                      RunMetadata* run_metadata, Costs* cost) const override;
  const VirtualScheduler* GetScheduler() const { return scheduler_.get(); }
 private:
  const GrapplerItem* item_;
  std::unique_ptr<OpLevelCostEstimator> node_estimator_;
  std::unique_ptr<ReadyNodeManager> node_manager_;
  std::unique_ptr<VirtualScheduler> scheduler_;
  bool use_static_shapes_;
  bool use_aggressive_shape_inference_;
};
}  
}  
#endif  
#include "tensorflow/core/grappler/costs/analytical_cost_estimator.h"
#include <limits>
#include <unordered_map>
#include "tensorflow/core/framework/tensor.pb.h"  
#include "tensorflow/core/framework/tensor_shape.pb.h"
#include "tensorflow/core/graph/types.h"
#include "tensorflow/core/grappler/costs/graph_properties.h"
#include "tensorflow/core/grappler/costs/op_performance_data.pb.h"
#include "tensorflow/core/grappler/costs/utils.h"
#include "tensorflow/core/grappler/costs/virtual_placer.h"
#include "tensorflow/core/grappler/costs/virtual_scheduler.h"
#include "tensorflow/core/grappler/grappler_item.h"
#include "tensorflow/core/grappler/op_types.h"
#include "tensorflow/core/grappler/utils.h"
#include "tensorflow/core/util/overflow.h"
namespace tensorflow {
namespace grappler {
namespace {
Status AddCostNode(ReadyNodeManager* node_manager, const OpContext& op_context,
                   int node_id, const Costs& node_costs,
                   gtl::FlatMap<string, CostGraphDef::Node*>* name_to_cost_node,
                   gtl::FlatMap<string, int>* name_to_id,
                   CostGraphDef* cost_graph) {
  const string& op_name = op_context.name;
  auto it = name_to_cost_node->find(op_name);
  CostGraphDef::Node* node;
  if (it != name_to_cost_node->end()) {
    node = it->second;
    node->clear_input_info();
    node->clear_output_info();
  } else {
    node = cost_graph->add_node();
    (*name_to_cost_node)[op_name] = node;
    node->set_name(op_name);
    node->set_id(node_id);
    (*name_to_id)[node->name()] = node->id();
  }
  node->set_device(op_context.device_name);
  node->set_compute_cost(node_costs.execution_time.asMicroSeconds().count());
  node->set_compute_time(node_costs.compute_time.asMicroSeconds().count());
  node->set_memory_time(node_costs.memory_time.asMicroSeconds().count());
  node->set_temporary_memory_size(node_costs.temporary_memory);
  node->set_persistent_memory_size(node_costs.persistent_memory);
  node->set_inaccurate(node_costs.inaccurate);
  for (const string& input : node_manager->GetCurrNode()->input()) {
    int input_port;
    string input_name = ParseNodeName(input, &input_port);
    if (name_to_id->find(input_name) == name_to_id->end()) {
      if (!IsMerge(*node_manager->GetCurrNode()))
        VLOG(1) << "input: " << input
                << " not found for non-Merge node: " << op_name;
      continue;
    }
    if (IsControlInput(input)) {
      node->add_control_input(name_to_id->at(input_name));
    } else {
      auto* input_info = node->add_input_info();
      input_info->set_preceding_node(name_to_id->at(input_name));
      input_info->set_preceding_port(input_port);
    }
  }
  for (const auto& output : op_context.op_info.outputs()) {
    auto output_info = node->add_output_info();
    output_info->set_alias_input_port(-1);
    output_info->set_dtype(output.dtype());
    *output_info->mutable_shape() = output.shape();
    int64_t size = DataTypeSize(output.dtype());
    for (const auto& dim : output.shape().dim()) {
      size = MultiplyWithoutOverflow(size, std::max<int64_t>(1, dim.size()));
      if (size < 0) {
        return errors::InvalidArgument(
            "Integer overflow encountered in dimension size.");
      }
    }
    output_info->set_size(size);
  }
  return absl::OkStatus();
}
}  
AnalyticalCostEstimator::AnalyticalCostEstimator(
    Cluster* cluster, bool use_static_shapes,
    bool use_aggressive_shape_inference)
    : AnalyticalCostEstimator(
          cluster, std::make_unique<OpLevelCostEstimator>(),
          ReadyNodeManagerFactory("FirstReady"), use_static_shapes,
          use_aggressive_shape_inference) {}
AnalyticalCostEstimator::AnalyticalCostEstimator(
    Cluster* cluster, std::unique_ptr<OpLevelCostEstimator> node_estimator,
    std::unique_ptr<ReadyNodeManager> node_manager, bool use_static_shapes,
    bool use_aggressive_shape_inference)
    : node_estimator_(std::move(node_estimator)),
      node_manager_(std::move(node_manager)),
      use_static_shapes_(use_static_shapes),
      use_aggressive_shape_inference_(use_aggressive_shape_inference) {
  scheduler_ = std::make_unique<VirtualScheduler>(
      use_static_shapes_, use_aggressive_shape_inference_, cluster,
      node_manager_.get(),
      std::make_unique<VirtualPlacer>(cluster->GetDevices()));
}
AnalyticalCostEstimator::AnalyticalCostEstimator(
    Cluster* cluster, std::unique_ptr<OpLevelCostEstimator> node_estimator,
    std::unique_ptr<ReadyNodeManager> node_manager,
    std::unique_ptr<VirtualPlacer> placer, bool use_static_shapes,
    bool use_aggressive_shape_inference)
    : node_estimator_(std::move(node_estimator)),
      node_manager_(std::move(node_manager)),
      use_static_shapes_(use_static_shapes),
      use_aggressive_shape_inference_(use_aggressive_shape_inference) {
  scheduler_ = std::make_unique<VirtualScheduler>(
      use_static_shapes_, use_aggressive_shape_inference_, cluster,
      node_manager_.get(), std::move(placer));
}
Status AnalyticalCostEstimator::Initialize(const GrapplerItem& item) {
  item_ = &item;
  return absl::OkStatus();
}
Status AnalyticalCostEstimator::PredictCosts(const GraphDef& optimized_graph,
                                             RunMetadata* run_metadata,
                                             Costs* costs) const {
  std::unique_ptr<GrapplerItem> item_storage;
  const GrapplerItem* item;
  if (&optimized_graph == &item_->graph) {
    item = item_;
  } else {
    GraphDef graph_copy = optimized_graph;
    item_storage = std::make_unique<GrapplerItem>(
        item_->WithGraph(std::move(graph_copy)));
    item = item_storage.get();
  }
  auto status = scheduler_->Init(item);
  if (!status.ok()) {
    if (costs) {
      costs->execution_time = Costs::Duration::max();
    }
    return status;
  }
  gtl::FlatMap<string, CostGraphDef::Node*> name_to_cost_node;
  CostGraphDef* cost_graph = nullptr;
  if (run_metadata) {
    cost_graph = run_metadata->mutable_cost_graph();
    for (auto& node : *cost_graph->mutable_node()) {
      name_to_cost_node[node.name()] = &node;
    }
  }
  std::vector<string> inaccurate_nodes;
  int nodes_executed = 0;
  int node_id = 0;
  gtl::FlatMap<string, int> name_to_id;
  Costs node_costs;
  do {
    ++nodes_executed;
    OpContext op_context = scheduler_->GetCurrNode();
    node_costs = node_estimator_->PredictCosts(op_context);
    if (node_costs.inaccurate) {
      inaccurate_nodes.push_back(op_context.name);
      if (node_costs.num_ops_with_unknown_shapes > 0)
        VLOG(4) << op_context.name << " has "
                << node_costs.num_ops_with_unknown_shapes << " unknown shapes";
    }
    if (cost_graph) {
      Status s =
          AddCostNode(node_manager_.get(), op_context, node_id++, node_costs,
                      &name_to_cost_node, &name_to_id, cost_graph);
      if (!s.ok()) {
        return s;
      }
    }
  } while (scheduler_->MarkCurrNodeExecuted(node_costs));
  VLOG(1) << inaccurate_nodes.size() << " out of " << nodes_executed
          << " nodes have inaccurate time estimation";
  if (VLOG_IS_ON(3)) {
    for (const auto& node : inaccurate_nodes) {
      VLOG(4) << "Node with inaccurate time estimation: " << node;
    }
  }
  if (costs) {
    *costs = scheduler_->Summary(run_metadata);
  } else if (run_metadata) {
    scheduler_->GenerateRunMetadata(run_metadata);
  }
  if (VLOG_IS_ON(1)) {
    bool verbose = VLOG_IS_ON(2);
    if (run_metadata) {
      VLOG(1) << GetStatsStringFromRunMetadata(*run_metadata, verbose);
    } else {
      RunMetadata run_metadata;
      scheduler_->GenerateRunMetadata(&run_metadata);
      VLOG(1) << GetStatsStringFromRunMetadata(run_metadata, verbose);
    }
  }
  return absl::OkStatus();
}
}  
}  