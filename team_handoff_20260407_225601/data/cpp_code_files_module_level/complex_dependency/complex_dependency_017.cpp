#ifndef XLA_SERVICE_GPU_RUNTIME_CONVOLUTION_THUNK_H_
#define XLA_SERVICE_GPU_RUNTIME_CONVOLUTION_THUNK_H_
#include <cstdint>
#include <memory>
#include <vector>
#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_map.h"
#include "absl/container/inlined_vector.h"
#include "absl/status/status.h"
#include "absl/synchronization/mutex.h"
#include "absl/types/span.h"
#include "xla/service/buffer_assignment.h"
#include "xla/service/gpu/gpu_conv_runner.h"
#include "xla/service/gpu/runtime/thunk.h"
#include "xla/stream_executor/dnn.h"
#include "xla/stream_executor/stream_executor.h"
namespace xla {
namespace gpu {
class ConvolutionThunk : public Thunk {
 public:
  ConvolutionThunk(ThunkInfo thunk_info, GpuConvConfig config,
                   std::vector<BufferAllocation::Slice> operand_slices,
                   std::vector<BufferAllocation::Slice> result_slices,
                   BufferAllocation::Slice scratch_slice);
  ConvolutionThunk(const ConvolutionThunk&) = delete;
  ConvolutionThunk& operator=(const ConvolutionThunk&) = delete;
  absl::Status ExecuteOnStream(const ExecuteParams& params) override;
 private:
  std::vector<BufferAllocation::Slice> operand_buffers_;
  std::vector<BufferAllocation::Slice> result_buffers_;
  BufferAllocation::Slice scratch_buffer_;
  GenericConvRunner& GetOrCreateRunner(const stream_executor::Stream* stream,
                                       bool* runner_created);
  const GpuConvConfig config_;
  absl::Mutex mu_;
  absl::flat_hash_map<const stream_executor::Stream*,
                      std::unique_ptr<GenericConvRunner>>
      runner_cache_ ABSL_GUARDED_BY(mu_);
};
class ConvolutionReorderThunk : public Thunk {
 public:
  ConvolutionReorderThunk(
      ThunkInfo thunk_info, absl::Span<int64_t> filter_nchw,
      absl::InlinedVector<BufferAllocation::Slice, 2> operand_slices,
      absl::InlinedVector<BufferAllocation::Slice, 2> result_slices);
  ConvolutionReorderThunk(const ConvolutionReorderThunk&) = delete;
  ConvolutionReorderThunk& operator=(const ConvolutionReorderThunk&) = delete;
  absl::Status ExecuteOnStream(const ExecuteParams& params) override;
 private:
  static se::dnn::FilterDescriptor CreateFilterDescriptor(
      absl::Span<int64_t> filter_nchw);
  const se::dnn::FilterDescriptor filter_descriptor_;
  absl::InlinedVector<BufferAllocation::Slice, 2> operand_buffers_;
  absl::InlinedVector<BufferAllocation::Slice, 2> result_buffers_;
};
}  
}  
#endif  
#include "xla/service/gpu/runtime/convolution_thunk.h"
#include <cstdint>
#include <memory>
#include <optional>
#include <utility>
#include <vector>
#include "absl/container/inlined_vector.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/synchronization/mutex.h"
#include "absl/types/span.h"
#include "xla/service/buffer_assignment.h"
#if TENSORFLOW_USE_ROCM
#include "xla/service/gpu/stream_executor_util.h"
#endif  
#include "xla/service/gpu/gpu_conv_runner.h"
#include "xla/service/gpu/runtime/thunk.h"
#include "xla/stream_executor/device_memory.h"
#include "xla/stream_executor/dnn.h"
#include "xla/stream_executor/scratch_allocator.h"
#include "xla/stream_executor/stream_executor.h"
#include "xla/util.h"
#include "tsl/platform/errors.h"
namespace xla {
namespace gpu {
ConvolutionThunk::ConvolutionThunk(
    ThunkInfo thunk_info, GpuConvConfig config,
    std::vector<BufferAllocation::Slice> operand_slices,
    std::vector<BufferAllocation::Slice> result_slices,
    BufferAllocation::Slice scratch_slice)
    : Thunk(Kind::kConvolution, thunk_info),
      operand_buffers_(std::move(operand_slices)),
      result_buffers_(std::move(result_slices)),
      scratch_buffer_(scratch_slice),
      config_(std::move(config)) {}
GenericConvRunner& ConvolutionThunk::GetOrCreateRunner(
    const stream_executor::Stream* stream, bool* runner_created) {
  absl::MutexLock lock(&mu_);
  auto it = runner_cache_.find(stream);
  *runner_created = (it == runner_cache_.end());
  if (*runner_created) {
    it = runner_cache_
             .insert({stream, std::make_unique<GenericConvRunner>(config_)})
             .first;
  }
  return *it->second;
}
absl::Status ConvolutionThunk::ExecuteOnStream(const ExecuteParams& params) {
  const auto& buffer_allocations = *params.buffer_allocations;
  std::vector<se::DeviceMemoryBase> operand_se_buffers, result_se_buffers;
  operand_se_buffers.reserve(operand_buffers_.size());
  for (BufferAllocation::Slice buffer : operand_buffers_) {
    operand_se_buffers.push_back(buffer_allocations.GetDeviceAddress(buffer));
  }
  result_se_buffers.reserve(result_buffers_.size());
  for (BufferAllocation::Slice buffer : result_buffers_) {
    result_se_buffers.push_back(buffer_allocations.GetDeviceAddress(buffer));
  }
  se::DeviceMemoryBase scratch =
      buffer_allocations.GetDeviceAddress(scratch_buffer_);
  bool runner_created = false;
  RunConvOptions opts;
  opts.runner_cache = &GetOrCreateRunner(params.stream, &runner_created);
#if TENSORFLOW_USE_ROCM
  if (runner_created) {
    TF_ASSIGN_OR_RETURN(
        GpuConvParams conv_params,
        GetGpuConvParams(config_, operand_se_buffers, result_se_buffers));
    TF_ASSIGN_OR_RETURN(se::dnn::ConvolutionKind kind,
                        GetDNNConvKindFromCudnnConvKind(config_.kind));
    TF_ASSIGN_OR_RETURN(se::dnn::DataType input_type,
                        GetDNNDataTypeFromPrimitiveType(config_.input_type));
    TF_ASSIGN_OR_RETURN(se::dnn::DataType output_type,
                        GetDNNDataTypeFromPrimitiveType(config_.output_type));
    TF_ASSIGN_OR_RETURN(auto dnn,
                        se::dnn::internal::GetDnnFromStream(params.stream));
    se::OwningScratchAllocator<> scratch_allocator(
        buffer_allocations.device_ordinal(),
        buffer_allocations.memory_allocator());
    std::vector<se::dnn::ProfileResult> profile_results;
    dnn->GetMIOpenConvolveAlgorithms(
        kind, input_type, output_type, params.stream, config_.input_descriptor,
        conv_params.input_buf, config_.filter_descriptor,
        conv_params.filter_buf, config_.output_descriptor,
        conv_params.output_buf, config_.conv_desc, &scratch_allocator,
        &profile_results);
  }
#endif  
  TF_RETURN_IF_ERROR(RunGpuConv(config_, absl::MakeSpan(operand_se_buffers),
                                absl::MakeSpan(result_se_buffers), scratch,
                                params.stream, opts));
  if (!params.stream->ok()) {
    return Internal("ConvolutionThunk::ExecuteOnStream failed.");
  }
  return absl::OkStatus();
}
ConvolutionReorderThunk::ConvolutionReorderThunk(
    ThunkInfo thunk_info, absl::Span<int64_t> filter_nchw,
    absl::InlinedVector<BufferAllocation::Slice, 2> operand_slices,
    absl::InlinedVector<BufferAllocation::Slice, 2> result_slices)
    : Thunk(Kind::kConvolutionReorder, thunk_info),
      filter_descriptor_(CreateFilterDescriptor(filter_nchw)),
      operand_buffers_(operand_slices),
      result_buffers_(result_slices) {}
absl::Status ConvolutionReorderThunk::ExecuteOnStream(
    const ExecuteParams& params) {
  bool has_bias = operand_buffers_.size() > 1;
  CHECK_EQ(operand_buffers_.size(), result_buffers_.size());
  const auto& buffer_allocations = *params.buffer_allocations;
  auto filter_input = se::DeviceMemory<int8_t>(
      buffer_allocations.GetDeviceAddress(operand_buffers_[0]));
  auto filter_output = se::DeviceMemory<int8_t>(
      buffer_allocations.GetDeviceAddress(result_buffers_[0]));
  auto bias_input =
      has_bias ? std::make_optional(se::DeviceMemory<float>(
                     buffer_allocations.GetDeviceAddress(operand_buffers_[1])))
               : std::nullopt;
  auto bias_output =
      has_bias ? std::make_optional(se::DeviceMemory<float>(
                     buffer_allocations.GetDeviceAddress(result_buffers_[1])))
               : std::nullopt;
  auto dnn = params.stream->parent()->AsDnn();
  if (dnn == nullptr) {
    return absl::InternalError("No DNN for stream.");
  }
  return dnn->CudnnReorderConvolutionFilterAndBias(
      params.stream, filter_descriptor_, filter_input, &filter_output,
      std::move(bias_input), std::move(bias_output));
}
se::dnn::FilterDescriptor ConvolutionReorderThunk::CreateFilterDescriptor(
    absl::Span<int64_t> filter_nchw) {
  CHECK_EQ(filter_nchw.size(), 4);
  se::dnn::FilterDescriptor filter_desc(2);
  filter_desc.set_layout(se::dnn::FilterLayout::kOutputInputYX32);
  filter_desc.set_output_feature_map_count(filter_nchw[0]);
  filter_desc.set_input_feature_map_count(filter_nchw[1]);
  filter_desc.set_input_filter_height(filter_nchw[2]);
  filter_desc.set_input_filter_width(filter_nchw[3]);
  return filter_desc;
}
}  
}  