#ifndef XLA_SERVICE_GPU_RUNTIME_FOR_ALL_THUNKS_H_
#define XLA_SERVICE_GPU_RUNTIME_FOR_ALL_THUNKS_H_
#include "absl/functional/function_ref.h"
#include "xla/service/gpu/runtime/thunk.h"
namespace xla::gpu {
void ForAllThunks(absl::FunctionRef<void(const Thunk*)> fn, const Thunk* thunk);
void ForAllThunks(absl::FunctionRef<void(const Thunk*)> fn,
                  const ThunkSequence* thunks);
}  
#endif  
#include "xla/service/gpu/runtime/for_all_thunks.h"
#include <memory>
#include <optional>
#include "absl/functional/function_ref.h"
#include "xla/service/gpu/runtime/command_buffer_thunk.h"
#include "xla/service/gpu/runtime/conditional_thunk.h"
#include "xla/service/gpu/runtime/dynamic_slice_thunk.h"
#include "xla/service/gpu/runtime/sequential_thunk.h"
#include "xla/service/gpu/runtime/thunk.h"
#include "xla/service/gpu/runtime/while_thunk.h"
#include "tsl/platform/casts.h"
namespace xla::gpu {
void ForAllThunks(absl::FunctionRef<void(const Thunk*)> fn,
                  const Thunk* thunk) {
  fn(thunk);
  switch (thunk->kind()) {
    case Thunk::kAddressComputation:
      ForAllThunks(fn, tensorflow::down_cast<const DynamicSliceThunk*>(thunk)
                           ->embedded_thunk());
      break;
    case Thunk::kCommandBuffer:
      if (const std::unique_ptr<SequentialThunk>& sequence =
              tensorflow::down_cast<const CommandBufferThunk*>(thunk)->thunks();
          sequence != nullptr) {
        ForAllThunks(fn, sequence.get());
      }
      break;
    case Thunk::kConditional:
      for (const std::unique_ptr<SequentialThunk>& branch :
           tensorflow::down_cast<const ConditionalThunk*>(thunk)
               ->branch_thunks()) {
        ForAllThunks(fn, branch.get());
      }
      break;
    case Thunk::kSequential:
      ForAllThunks(
          fn, &tensorflow::down_cast<const SequentialThunk*>(thunk)->thunks());
      break;
    case Thunk::kWhile:
      ForAllThunks(fn, tensorflow::down_cast<const WhileThunk*>(thunk)
                           ->condition_thunk_sequence());
      ForAllThunks(fn, tensorflow::down_cast<const WhileThunk*>(thunk)
                           ->body_thunk_sequence());
      break;
    case Thunk::kCholesky:
    case Thunk::kConvolution:
    case Thunk::kConvolutionReorder:
    case Thunk::kCopy:
    case Thunk::kCopyDone:
    case Thunk::kCubSort:
    case Thunk::kCublasLtMatmul:
    case Thunk::kCustomCall:
    case Thunk::kCustomKernel:
    case Thunk::kCuDnn:
    case Thunk::kFft:
    case Thunk::kFusedMHA:
    case Thunk::kGemm:
    case Thunk::kInfeed:
    case Thunk::kKernel:
    case Thunk::kMemset32BitValue:
    case Thunk::kMemzero:
    case Thunk::kNcclAllGather:
    case Thunk::kNcclAllGatherStart:
    case Thunk::kNcclAllGatherDone:
    case Thunk::kNcclAllReduce:
    case Thunk::kNcclAllReduceStart:
    case Thunk::kNcclAllReduceDone:
    case Thunk::kNcclCollectiveBroadcast:
    case Thunk::kNcclCollectiveBroadcastStart:
    case Thunk::kNcclCollectiveBroadcastDone:
    case Thunk::kNcclCollectivePermute:
    case Thunk::kNcclCollectivePermuteStart:
    case Thunk::kNcclCollectivePermuteDone:
    case Thunk::kNcclReduceScatter:
    case Thunk::kNcclReduceScatterStart:
    case Thunk::kNcclReduceScatterDone:
    case Thunk::kNcclAllToAll:
    case Thunk::kNcclAllToAllStart:
    case Thunk::kNcclAllToAllDone:
    case Thunk::kNcclSend:
    case Thunk::kNcclSendDone:
    case Thunk::kNcclRecv:
    case Thunk::kNcclRecvDone:
    case Thunk::kNorm:
    case Thunk::kOutfeed:
    case Thunk::kPartitionId:
    case Thunk::kRecv:
    case Thunk::kRecvDone:
    case Thunk::kReplicaId:
    case Thunk::kSend:
    case Thunk::kSendDone:
    case Thunk::kTriangularSolve:
    case Thunk::kWaitForStreams:
      break;
  }
}
void ForAllThunks(absl::FunctionRef<void(const Thunk*)> fn,
                  const ThunkSequence* thunks) {
  for (const std::unique_ptr<Thunk>& thunk : *thunks) {
    ForAllThunks(fn, thunk.get());
  }
}
}  