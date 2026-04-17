#ifndef ABSL_RANDOM_INTERNAL_POOL_URBG_H_
#define ABSL_RANDOM_INTERNAL_POOL_URBG_H_
#include <cinttypes>
#include <limits>
#include "absl/random/internal/traits.h"
#include "absl/types/span.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
namespace random_internal {
template <typename T>
class RandenPool {
 public:
  using result_type = T;
  static_assert(std::is_unsigned<result_type>::value,
                "RandenPool template argument must be a built-in unsigned "
                "integer type");
  static constexpr result_type(min)() {
    return (std::numeric_limits<result_type>::min)();
  }
  static constexpr result_type(max)() {
    return (std::numeric_limits<result_type>::max)();
  }
  RandenPool() {}
  inline result_type operator()() { return Generate(); }
  static void Fill(absl::Span<result_type> data);
 protected:
  static result_type Generate();
};
extern template class RandenPool<uint8_t>;
extern template class RandenPool<uint16_t>;
extern template class RandenPool<uint32_t>;
extern template class RandenPool<uint64_t>;
template <typename T, size_t kBufferSize>
class PoolURBG {
  using unsigned_type = typename make_unsigned_bits<T>::type;
  using PoolType = RandenPool<unsigned_type>;
  using SpanType = absl::Span<unsigned_type>;
  static constexpr size_t kInitialBuffer = kBufferSize + 1;
  static constexpr size_t kHalfBuffer = kBufferSize / 2;
 public:
  using result_type = T;
  static_assert(std::is_unsigned<result_type>::value,
                "PoolURBG must be parameterized by an unsigned integer type");
  static_assert(kBufferSize > 1,
                "PoolURBG must be parameterized by a buffer-size > 1");
  static_assert(kBufferSize <= 256,
                "PoolURBG must be parameterized by a buffer-size <= 256");
  static constexpr result_type(min)() {
    return (std::numeric_limits<result_type>::min)();
  }
  static constexpr result_type(max)() {
    return (std::numeric_limits<result_type>::max)();
  }
  PoolURBG() : next_(kInitialBuffer) {}
  PoolURBG(const PoolURBG&) : next_(kInitialBuffer) {}
  const PoolURBG& operator=(const PoolURBG&) {
    next_ = kInitialBuffer;
    return *this;
  }
  PoolURBG(PoolURBG&&) = default;
  PoolURBG& operator=(PoolURBG&&) = default;
  inline result_type operator()() {
    if (next_ >= kBufferSize) {
      next_ = (kBufferSize > 2 && next_ > kBufferSize) ? kHalfBuffer : 0;
      PoolType::Fill(SpanType(reinterpret_cast<unsigned_type*>(state_ + next_),
                              kBufferSize - next_));
    }
    return state_[next_++];
  }
 private:
  size_t next_;  
  result_type state_[kBufferSize];
};
}  
ABSL_NAMESPACE_END
}  
#endif  
#include "absl/random/internal/pool_urbg.h"
#include <algorithm>
#include <atomic>
#include <cstdint>
#include <cstring>
#include <iterator>
#include "absl/base/attributes.h"
#include "absl/base/call_once.h"
#include "absl/base/config.h"
#include "absl/base/internal/endian.h"
#include "absl/base/internal/raw_logging.h"
#include "absl/base/internal/spinlock.h"
#include "absl/base/internal/sysinfo.h"
#include "absl/base/internal/unaligned_access.h"
#include "absl/base/optimization.h"
#include "absl/random/internal/randen.h"
#include "absl/random/internal/seed_material.h"
#include "absl/random/seed_gen_exception.h"
using absl::base_internal::SpinLock;
using absl::base_internal::SpinLockHolder;
namespace absl {
ABSL_NAMESPACE_BEGIN
namespace random_internal {
namespace {
class RandenPoolEntry {
 public:
  static constexpr size_t kState = RandenTraits::kStateBytes / sizeof(uint32_t);
  static constexpr size_t kCapacity =
      RandenTraits::kCapacityBytes / sizeof(uint32_t);
  void Init(absl::Span<const uint32_t> data) {
    SpinLockHolder l(&mu_);  
    std::copy(data.begin(), data.end(), std::begin(state_));
    next_ = kState;
  }
  void Fill(uint8_t* out, size_t bytes) ABSL_LOCKS_EXCLUDED(mu_);
  template <typename T>
  inline T Generate() ABSL_LOCKS_EXCLUDED(mu_);
  inline void MaybeRefill() ABSL_EXCLUSIVE_LOCKS_REQUIRED(mu_) {
    if (next_ >= kState) {
      next_ = kCapacity;
      impl_.Generate(state_);
    }
  }
 private:
  uint32_t state_[kState] ABSL_GUARDED_BY(mu_);  
  SpinLock mu_;
  const Randen impl_;
  size_t next_ ABSL_GUARDED_BY(mu_);
};
template <>
inline uint8_t RandenPoolEntry::Generate<uint8_t>() {
  SpinLockHolder l(&mu_);
  MaybeRefill();
  return static_cast<uint8_t>(state_[next_++]);
}
template <>
inline uint16_t RandenPoolEntry::Generate<uint16_t>() {
  SpinLockHolder l(&mu_);
  MaybeRefill();
  return static_cast<uint16_t>(state_[next_++]);
}
template <>
inline uint32_t RandenPoolEntry::Generate<uint32_t>() {
  SpinLockHolder l(&mu_);
  MaybeRefill();
  return state_[next_++];
}
template <>
inline uint64_t RandenPoolEntry::Generate<uint64_t>() {
  SpinLockHolder l(&mu_);
  if (next_ >= kState - 1) {
    next_ = kCapacity;
    impl_.Generate(state_);
  }
  auto p = state_ + next_;
  next_ += 2;
  uint64_t result;
  std::memcpy(&result, p, sizeof(result));
  return result;
}
void RandenPoolEntry::Fill(uint8_t* out, size_t bytes) {
  SpinLockHolder l(&mu_);
  while (bytes > 0) {
    MaybeRefill();
    size_t remaining = (kState - next_) * sizeof(state_[0]);
    size_t to_copy = std::min(bytes, remaining);
    std::memcpy(out, &state_[next_], to_copy);
    out += to_copy;
    bytes -= to_copy;
    next_ += (to_copy + sizeof(state_[0]) - 1) / sizeof(state_[0]);
  }
}
static constexpr size_t kPoolSize = 8;
static absl::once_flag pool_once;
ABSL_CACHELINE_ALIGNED static RandenPoolEntry* shared_pools[kPoolSize];
size_t GetPoolID() {
  static_assert(kPoolSize >= 1,
                "At least one urbg instance is required for PoolURBG");
  ABSL_CONST_INIT static std::atomic<uint64_t> sequence{0};
#ifdef ABSL_HAVE_THREAD_LOCAL
  static thread_local size_t my_pool_id = kPoolSize;
  if (ABSL_PREDICT_FALSE(my_pool_id == kPoolSize)) {
    my_pool_id = (sequence++ % kPoolSize);
  }
  return my_pool_id;
#else
  static pthread_key_t tid_key = [] {
    pthread_key_t tmp_key;
    int err = pthread_key_create(&tmp_key, nullptr);
    if (err) {
      ABSL_RAW_LOG(FATAL, "pthread_key_create failed with %d", err);
    }
    return tmp_key;
  }();
  uintptr_t my_pool_id =
      reinterpret_cast<uintptr_t>(pthread_getspecific(tid_key));
  if (ABSL_PREDICT_FALSE(my_pool_id == 0)) {
    my_pool_id = (sequence++ % kPoolSize) + 1;
    int err = pthread_setspecific(tid_key, reinterpret_cast<void*>(my_pool_id));
    if (err) {
      ABSL_RAW_LOG(FATAL, "pthread_setspecific failed with %d", err);
    }
  }
  return my_pool_id - 1;
#endif
}
RandenPoolEntry* PoolAlignedAlloc() {
  constexpr size_t kAlignment =
      ABSL_CACHELINE_SIZE > 32 ? ABSL_CACHELINE_SIZE : 32;
  uintptr_t x = reinterpret_cast<uintptr_t>(
      new char[sizeof(RandenPoolEntry) + kAlignment]);
  auto y = x % kAlignment;
  void* aligned = reinterpret_cast<void*>(y == 0 ? x : (x + kAlignment - y));
  return new (aligned) RandenPoolEntry();
}
void InitPoolURBG() {
  static constexpr size_t kSeedSize =
      RandenTraits::kStateBytes / sizeof(uint32_t);
  uint32_t seed_material[kPoolSize * kSeedSize];
  if (!random_internal::ReadSeedMaterialFromOSEntropy(
          absl::MakeSpan(seed_material))) {
    random_internal::ThrowSeedGenException();
  }
  for (size_t i = 0; i < kPoolSize; i++) {
    shared_pools[i] = PoolAlignedAlloc();
    shared_pools[i]->Init(
        absl::MakeSpan(&seed_material[i * kSeedSize], kSeedSize));
  }
}
RandenPoolEntry* GetPoolForCurrentThread() {
  absl::call_once(pool_once, InitPoolURBG);
  return shared_pools[GetPoolID()];
}
}  
template <typename T>
typename RandenPool<T>::result_type RandenPool<T>::Generate() {
  auto* pool = GetPoolForCurrentThread();
  return pool->Generate<T>();
}
template <typename T>
void RandenPool<T>::Fill(absl::Span<result_type> data) {
  auto* pool = GetPoolForCurrentThread();
  pool->Fill(reinterpret_cast<uint8_t*>(data.data()),
             data.size() * sizeof(result_type));
}
template class RandenPool<uint8_t>;
template class RandenPool<uint16_t>;
template class RandenPool<uint32_t>;
template class RandenPool<uint64_t>;
}  
ABSL_NAMESPACE_END
}  