#ifndef XLA_BACKENDS_INTERPRETER_EXECUTOR_H_
#define XLA_BACKENDS_INTERPRETER_EXECUTOR_H_
#include <cstdint>
#include <memory>
#include "absl/functional/any_invocable.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/types/span.h"
#include "xla/shape.h"
#include "xla/stream_executor/blas.h"
#include "xla/stream_executor/device_description.h"
#include "xla/stream_executor/device_memory.h"
#include "xla/stream_executor/event.h"
#include "xla/stream_executor/host/host_stream.h"
#include "xla/stream_executor/host_memory_allocation.h"
#include "xla/stream_executor/kernel.h"
#include "xla/stream_executor/kernel_spec.h"
#include "xla/stream_executor/launch_dim.h"
#include "xla/stream_executor/memory_allocation.h"
#include "xla/stream_executor/stream.h"
#include "xla/stream_executor/stream_executor.h"
#include "xla/stream_executor/stream_executor_common.h"
#include "xla/xla_data.pb.h"
namespace stream_executor {
namespace interpreter {
class InterpreterStream : public host::HostStream {
 public:
  explicit InterpreterStream(StreamExecutor *executor)
      : host::HostStream(executor) {}
  absl::Status WaitFor(Stream *stream) override {
    return host::HostStream::WaitFor(stream);
  }
  absl::Status WaitFor(Event *event) override {
    return absl::UnimplementedError("Not implemented.");
  }
  absl::Status RecordEvent(Event *event) override {
    return absl::UnimplementedError("Not implemented.");
  }
  absl::Status Memcpy(void *host_dst, const DeviceMemoryBase &gpu_src,
                      uint64_t size) override {
    void *src_mem = const_cast<void *>(gpu_src.opaque());
    EnqueueTask(
        [this, host_dst, src_mem, size]() { memcpy(host_dst, src_mem, size); });
    return BlockUntilDone();
  }
  absl::Status Memcpy(DeviceMemoryBase *gpu_dst, const void *host_src,
                      uint64_t size) override {
    void *dst_mem = gpu_dst->opaque();
    EnqueueTask(
        [this, dst_mem, host_src, size]() { memcpy(dst_mem, host_src, size); });
    return BlockUntilDone();
  }
};
class XlaInterpreterExecutor : public StreamExecutorCommon {
 public:
  XlaInterpreterExecutor(int device_ordinal, Platform *platform)
      : StreamExecutorCommon(platform), device_ordinal_(device_ordinal) {}
  absl::Status Init() override { return absl::OkStatus(); }
  int device_ordinal() const override { return device_ordinal_; };
  absl::Status GetKernel(const MultiKernelLoaderSpec &spec,
                         Kernel *kernel) override {
    return absl::UnimplementedError("Not Implemented");
  }
  absl::Status Launch(Stream *stream, const ThreadDim &thread_dims,
                      const BlockDim &block_dims, const Kernel &kernel,
                      const KernelArgs &args) override {
    return absl::UnimplementedError("Not Implemented");
  }
  DeviceMemoryBase Allocate(uint64_t size, int64_t memory_space) override;
  void Deallocate(DeviceMemoryBase *mem) override;
  absl::StatusOr<std::unique_ptr<MemoryAllocation>> HostMemoryAllocate(
      uint64_t size) override {
    return std::make_unique<HostMemoryAllocation>(new char[size], size, this);
  }
  void HostMemoryDeallocate(void *mem) override {
    delete[] static_cast<char *>(mem);
  }
  absl::Status Memset(Stream *stream, DeviceMemoryBase *location,
                      uint8_t pattern, uint64_t size) override {
    return absl::InternalError("Interpreter can not memset");
  }
  bool SynchronizeAllActivity() override { return true; }
  absl::Status SynchronousMemZero(DeviceMemoryBase *location,
                                  uint64_t size) override {
    return absl::InternalError("Interpreter can not memzero");
  }
  absl::Status SynchronousMemcpy(DeviceMemoryBase *dev_dst,
                                 const void *host_src, uint64_t size) override;
  absl::Status SynchronousMemcpy(void *host_dst,
                                 const DeviceMemoryBase &dev_src,
                                 uint64_t size) override;
  bool HostCallback(Stream *stream,
                    absl::AnyInvocable<absl::Status() &&> callback) override;
  void DeallocateStream(Stream *stream) override {}
  absl::Status BlockHostUntilDone(Stream *stream) override;
  bool DeviceMemoryUsage(int64_t *free, int64_t *total) const override {
    return false;
  }
  absl::StatusOr<std::unique_ptr<DeviceDescription>> CreateDeviceDescription()
      const override {
    return CreateDeviceDescription(0);
  }
  static absl::StatusOr<std::unique_ptr<DeviceDescription>>
  CreateDeviceDescription(int device_ordinal);
  absl::Status EnablePeerAccessTo(StreamExecutor *other) override {
    return absl::OkStatus();
  }
  bool CanEnablePeerAccessTo(StreamExecutor *other) override { return true; }
  absl::StatusOr<std::unique_ptr<Event>> CreateEvent() override {
    return std::make_unique<Event>();
  }
  absl::StatusOr<std::unique_ptr<Stream>> CreateStream(
      std::optional<std::variant<StreamPriority, int>> priority =
          std::nullopt) override {
    return std::make_unique<InterpreterStream>(this);
  }
 private:
  int device_ordinal_;
};
}  
}  
#endif  
#include "xla/backends/interpreter/executor.h"
#include <cstring>
#include <utility>
#include "absl/functional/any_invocable.h"
#include "absl/log/log.h"
#include "xla/status_macros.h"
namespace stream_executor {
namespace interpreter {
host::HostStream *AsExecutorStream(Stream *stream) {
  DCHECK(stream != nullptr);
  return dynamic_cast<host::HostStream *>(stream);
}
DeviceMemoryBase XlaInterpreterExecutor::Allocate(uint64_t size,
                                                  int64_t memory_space) {
  return DeviceMemoryBase(new char[size], size);
}
void XlaInterpreterExecutor::Deallocate(DeviceMemoryBase *mem) {
  delete[] static_cast<char *>(mem->opaque());
}
absl::Status XlaInterpreterExecutor::SynchronousMemcpy(
    DeviceMemoryBase *dev_dst, const void *host_src, uint64_t size) {
  memcpy(dev_dst->opaque(), host_src, size);
  return absl::OkStatus();
}
absl::Status XlaInterpreterExecutor::SynchronousMemcpy(
    void *host_dst, const DeviceMemoryBase &dev_src, uint64_t size) {
  memcpy(host_dst, dev_src.opaque(), size);
  return absl::OkStatus();
}
bool XlaInterpreterExecutor::HostCallback(
    Stream *stream, absl::AnyInvocable<absl::Status() &&> callback) {
  AsExecutorStream(stream)->EnqueueTaskWithStatus(std::move(callback));
  return true;
}
absl::Status XlaInterpreterExecutor::BlockHostUntilDone(Stream *stream) {
  return AsExecutorStream(stream)->BlockUntilDone();
}
absl::StatusOr<std::unique_ptr<DeviceDescription>>
XlaInterpreterExecutor::CreateDeviceDescription(int device_ordinal) {
  internal::DeviceDescriptionBuilder builder;
  builder.set_device_address_bits(64);
  builder.set_name("Interpreter");
  builder.set_device_memory_size(static_cast<uint64_t>(4) * 1024 * 1024 * 1024);
  builder.set_clock_rate_ghz(static_cast<float>(CLOCKS_PER_SEC) / 1e9);
  return builder.Build();
}
}  
}  