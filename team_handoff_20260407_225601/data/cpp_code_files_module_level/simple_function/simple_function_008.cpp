#ifndef XLA_SERVICE_FLATTEN_CALL_GRAPH_H_
#define XLA_SERVICE_FLATTEN_CALL_GRAPH_H_
#include "absl/status/statusor.h"
#include "xla/service/hlo_pass_interface.h"
namespace xla {
class FlattenCallGraph : public HloModulePass {
 public:
  absl::string_view name() const override { return "flatten-call-graph"; }
  using HloPassInterface::Run;
  absl::StatusOr<bool> Run(
      HloModule* module,
      const absl::flat_hash_set<absl::string_view>& execution_threads) override;
};
}  
#endif  
#include "xla/service/flatten_call_graph.h"
#include <memory>
#include <vector>
#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "xla/hlo/ir/hlo_computation.h"
#include "xla/hlo/ir/hlo_instruction.h"
#include "xla/hlo/ir/hlo_module.h"
#include "xla/hlo/ir/hlo_opcode.h"
#include "xla/hlo/utils/hlo_query.h"
#include "xla/service/call_graph.h"
#include "xla/util.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/logging.h"
namespace xla {
namespace {
void ReplaceCalledComputation(HloInstruction* instruction,
                              HloComputation* computation,
                              HloComputation* new_computation) {
  switch (instruction->opcode()) {
    case HloOpcode::kWhile: {
      if (computation == instruction->while_condition()) {
        instruction->set_while_condition(new_computation);
      } else {
        CHECK_EQ(computation, instruction->while_body());
        instruction->set_while_body(new_computation);
      }
      break;
    }
    case HloOpcode::kCall: {
      CHECK_EQ(instruction->to_apply(), computation);
      instruction->set_to_apply(new_computation);
      break;
    }
    case HloOpcode::kConditional: {
      for (int b = 0; b < instruction->branch_count(); ++b) {
        if (b == instruction->branch_count() - 1) {
          CHECK_EQ(computation, instruction->branch_computation(b));
        }
        if (computation == instruction->branch_computation(b)) {
          instruction->set_branch_computation(b, new_computation);
          break;
        }
      }
      break;
    }
    default:
      LOG(FATAL) << "unexpected opcode: " << instruction->opcode();
  }
}
absl::Status FlattenNode(const CallGraphNode& node) {
  HloComputation* computation = node.computation();
  HloModule* module = computation->parent();
  for (int i = 0; i < node.caller_callsites().size(); ++i) {
    CallSite call_site = node.caller_callsites()[i];
    if (call_site.context() == CallContext::kEmbedded) {
      continue;
    }
    CHECK_EQ(call_site.context(), CallContext::kControlFlow);
    if (node.context() != CallContext::kBoth && i == 0) {
      continue;
    }
    if (computation->IsAsyncComputation()) {
      continue;
    }
    HloComputation* clone =
        module->AddEmbeddedComputation(computation->Clone());
    ReplaceCalledComputation(call_site.instruction(), computation, clone);
    std::vector<HloComputation*> worklist;
    worklist.push_back(clone);
    while (!worklist.empty()) {
      auto current = worklist.back();
      worklist.pop_back();
      for (auto* instruction : current->instructions()) {
        if (GetInstructionCallContext(instruction->opcode()) !=
            CallContext::kControlFlow) {
          continue;
        }
        for (auto callee : instruction->called_computations()) {
          HloComputation* callee_clone =
              module->AddEmbeddedComputation(callee->Clone());
          ReplaceCalledComputation(instruction, callee, callee_clone);
          worklist.push_back(callee_clone);
        }
      }
    }
  }
  return absl::OkStatus();
}
absl::Status AnnotateNode(const CallGraphNode& node) {
  for (auto& callsite : node.callsites()) {
    HloInstruction* instruction = callsite.instruction();
    if (instruction->opcode() == HloOpcode::kFusion) {
      for (HloComputation* computation : instruction->called_computations()) {
        computation->SetFusionInstruction(instruction);
      }
    } else if (instruction->opcode() == HloOpcode::kCustomCall) {
      for (HloComputation* computation : instruction->called_computations()) {
        computation->SetCustomCallInstruction(instruction);
      }
    } else if (hlo_query::IsCollectiveCommunicationOp(instruction->opcode())) {
      for (HloComputation* computation : instruction->called_computations()) {
        computation->SetCollectiveCallInstruction(instruction);
      }
    } else if (instruction->opcode() == HloOpcode::kWhile) {
      instruction->while_body()->SetWhileCallInstruction(instruction);
    } else if (instruction->opcode() == HloOpcode::kConditional) {
      for (HloComputation* branch : instruction->branch_computations()) {
        branch->SetConditionalCallInstruction(instruction);
      }
    }
  }
  return absl::OkStatus();
}
}  
absl::StatusOr<bool> FlattenCallGraph::Run(
    HloModule* module,
    const absl::flat_hash_set<absl::string_view>& execution_threads) {
  XLA_VLOG_LINES(3, "Before flatten call graph:\n" + module->ToString());
  {  
    std::unique_ptr<CallGraph> call_graph =
        CallGraph::Build(module, execution_threads);
    TF_RETURN_IF_ERROR(call_graph->VisitNodes(FlattenNode));
  }
  {  
    std::unique_ptr<CallGraph> call_graph =
        CallGraph::Build(module, execution_threads);
    TF_RETURN_IF_ERROR(call_graph->VisitNodes(AnnotateNode));
  }
  XLA_VLOG_LINES(3, "After flatten call graph:\n" + module->ToString());
  return true;
}
}  