#ifndef XLA_SERVICE_CPU_XFEED_MANAGER_H_
#define XLA_SERVICE_CPU_XFEED_MANAGER_H_
#include <deque>
#include "absl/status/statusor.h"
#include "absl/types/span.h"
#include "xla/shape.h"
#include "xla/types.h"
#include "xla/xla_data.pb.h"
namespace xla {
namespace cpu {
namespace runtime {
class XfeedBuffer {
 public:
  virtual ~XfeedBuffer() = default;
  virtual int32_t length() = 0;
  virtual void* data() = 0;
  virtual void Done(absl::StatusOr<Shape> shape) = 0;
};
class XfeedQueueManager {
 public:
  XfeedQueueManager(std::string queue_name) : queue_name_(queue_name) {}
  void Reset();
  void EnqueueBuffersAtomically(absl::Span<XfeedBuffer* const> buffers);
  XfeedBuffer* BlockingDequeueBuffer();
  void ReleaseCurrentBuffer(int32_t length, void* data,
                            absl::StatusOr<Shape> shape);
 private:
  const std::string queue_name_;
  absl::Mutex mu_;
  absl::CondVar cv_;
  std::deque<XfeedBuffer*> enqueued_buffers_;
  XfeedBuffer* current_buffer_ = nullptr;
};
class XfeedManager {
 public:
  XfeedManager() = default;
  void Reset();
  XfeedQueueManager* infeed() { return &infeed_; }
  XfeedQueueManager* outfeed() { return &outfeed_; }
 private:
  XfeedQueueManager infeed_ = {"infeed"};
  XfeedQueueManager outfeed_ = {"outfeed"};
};
int64_t GetByteSizeRequirement(const Shape& shape, int64_t pointer_size);
}  
}  
}  
#endif  
#include "xla/service/cpu/xfeed_manager.h"
#include "xla/shape_util.h"
#include "tsl/platform/logging.h"
namespace xla {
namespace cpu {
namespace runtime {
void XfeedManager::Reset() {
  infeed()->Reset();
  outfeed()->Reset();
}
void XfeedQueueManager::Reset() {
  absl::MutexLock l(&mu_);
  CHECK(current_buffer_ == nullptr);
  for (auto buffer : enqueued_buffers_) {
    buffer->Done(ShapeUtil::MakeNil());
  }
  enqueued_buffers_.clear();
}
void XfeedQueueManager::EnqueueBuffersAtomically(
    absl::Span<XfeedBuffer* const> buffers) {
  absl::MutexLock l(&mu_);
  bool was_empty = enqueued_buffers_.empty();
  for (XfeedBuffer* b : buffers) {
    VLOG(3) << "Enqueueing " << queue_name_ << " buffer (of " << buffers.size()
            << " buffers) with length: " << b->length();
    enqueued_buffers_.push_back(b);
  }
  if (was_empty && !buffers.empty()) {
    cv_.Signal();
  }
}
XfeedBuffer* XfeedQueueManager::BlockingDequeueBuffer() {
  absl::MutexLock l(&mu_);
  VLOG(3) << "Waiting for an available buffer.";
  while (enqueued_buffers_.empty()) {
    cv_.Wait(&mu_);
  }
  VLOG(3) << "A buffer is available!";
  CHECK(current_buffer_ == nullptr);
  current_buffer_ = enqueued_buffers_.front();
  enqueued_buffers_.pop_front();
  return current_buffer_;
}
void XfeedQueueManager::ReleaseCurrentBuffer(int32_t length, void* data,
                                             absl::StatusOr<Shape> shape) {
  VLOG(3) << "Releasing buffer with shape: "
          << (shape.ok() ? ShapeUtil::HumanString(shape.value())
                         : "<error status>");
  absl::MutexLock l(&mu_);
  CHECK(current_buffer_ != nullptr);
  CHECK_EQ(length, current_buffer_->length());
  CHECK_EQ(data, current_buffer_->data());
  current_buffer_->Done(std::move(shape));
  current_buffer_ = nullptr;
}
int64_t GetByteSizeRequirement(const Shape& shape, int64_t pointer_size) {
  if (shape.IsTuple() || shape.is_static()) {
    return ShapeUtil::ByteSizeOf(shape, pointer_size);
  }
  int64_t metadata_size = sizeof(int32_t) * shape.dimensions_size();
  return ShapeUtil::ByteSizeOf(shape, pointer_size) + metadata_size;
}
}  
}  
}  