#ifndef ABSL_SYNCHRONIZATION_MUTEX_H_
#define ABSL_SYNCHRONIZATION_MUTEX_H_
#include <atomic>
#include <cstdint>
#include <cstring>
#include <iterator>
#include <string>
#include "absl/base/attributes.h"
#include "absl/base/const_init.h"
#include "absl/base/internal/identity.h"
#include "absl/base/internal/low_level_alloc.h"
#include "absl/base/internal/thread_identity.h"
#include "absl/base/internal/tsan_mutex_interface.h"
#include "absl/base/port.h"
#include "absl/base/thread_annotations.h"
#include "absl/synchronization/internal/kernel_timeout.h"
#include "absl/synchronization/internal/per_thread_sem.h"
#include "absl/time/time.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
class Condition;
struct SynchWaitParams;
class ABSL_LOCKABLE ABSL_ATTRIBUTE_WARN_UNUSED Mutex {
 public:
  Mutex();
  explicit constexpr Mutex(absl::ConstInitType);
  ~Mutex();
  void Lock() ABSL_EXCLUSIVE_LOCK_FUNCTION();
  void Unlock() ABSL_UNLOCK_FUNCTION();
  ABSL_MUST_USE_RESULT bool TryLock() ABSL_EXCLUSIVE_TRYLOCK_FUNCTION(true);
  void AssertHeld() const ABSL_ASSERT_EXCLUSIVE_LOCK();
  void ReaderLock() ABSL_SHARED_LOCK_FUNCTION();
  void ReaderUnlock() ABSL_UNLOCK_FUNCTION();
  ABSL_MUST_USE_RESULT bool ReaderTryLock() ABSL_SHARED_TRYLOCK_FUNCTION(true);
  void AssertReaderHeld() const ABSL_ASSERT_SHARED_LOCK();
  void WriterLock() ABSL_EXCLUSIVE_LOCK_FUNCTION() { this->Lock(); }
  void WriterUnlock() ABSL_UNLOCK_FUNCTION() { this->Unlock(); }
  ABSL_MUST_USE_RESULT bool WriterTryLock()
      ABSL_EXCLUSIVE_TRYLOCK_FUNCTION(true) {
    return this->TryLock();
  }
  void Await(const Condition& cond) {
    AwaitCommon(cond, synchronization_internal::KernelTimeout::Never());
  }
  void LockWhen(const Condition& cond) ABSL_EXCLUSIVE_LOCK_FUNCTION() {
    LockWhenCommon(cond, synchronization_internal::KernelTimeout::Never(),
                   true);
  }
  void ReaderLockWhen(const Condition& cond) ABSL_SHARED_LOCK_FUNCTION() {
    LockWhenCommon(cond, synchronization_internal::KernelTimeout::Never(),
                   false);
  }
  void WriterLockWhen(const Condition& cond) ABSL_EXCLUSIVE_LOCK_FUNCTION() {
    this->LockWhen(cond);
  }
  bool AwaitWithTimeout(const Condition& cond, absl::Duration timeout) {
    return AwaitCommon(cond, synchronization_internal::KernelTimeout{timeout});
  }
  bool AwaitWithDeadline(const Condition& cond, absl::Time deadline) {
    return AwaitCommon(cond, synchronization_internal::KernelTimeout{deadline});
  }
  bool LockWhenWithTimeout(const Condition& cond, absl::Duration timeout)
      ABSL_EXCLUSIVE_LOCK_FUNCTION() {
    return LockWhenCommon(
        cond, synchronization_internal::KernelTimeout{timeout}, true);
  }
  bool ReaderLockWhenWithTimeout(const Condition& cond, absl::Duration timeout)
      ABSL_SHARED_LOCK_FUNCTION() {
    return LockWhenCommon(
        cond, synchronization_internal::KernelTimeout{timeout}, false);
  }
  bool WriterLockWhenWithTimeout(const Condition& cond, absl::Duration timeout)
      ABSL_EXCLUSIVE_LOCK_FUNCTION() {
    return this->LockWhenWithTimeout(cond, timeout);
  }
  bool LockWhenWithDeadline(const Condition& cond, absl::Time deadline)
      ABSL_EXCLUSIVE_LOCK_FUNCTION() {
    return LockWhenCommon(
        cond, synchronization_internal::KernelTimeout{deadline}, true);
  }
  bool ReaderLockWhenWithDeadline(const Condition& cond, absl::Time deadline)
      ABSL_SHARED_LOCK_FUNCTION() {
    return LockWhenCommon(
        cond, synchronization_internal::KernelTimeout{deadline}, false);
  }
  bool WriterLockWhenWithDeadline(const Condition& cond, absl::Time deadline)
      ABSL_EXCLUSIVE_LOCK_FUNCTION() {
    return this->LockWhenWithDeadline(cond, deadline);
  }
  void EnableInvariantDebugging(void (*invariant)(void*), void* arg);
  void EnableDebugLog(const char* name);
  void ForgetDeadlockInfo();
  void AssertNotHeld() const;
  typedef const struct MuHowS* MuHow;
  static void InternalAttemptToUseMutexInFatalSignalHandler();
 private:
  std::atomic<intptr_t> mu_;  
  static void IncrementSynchSem(Mutex* mu, base_internal::PerThreadSynch* w);
  static bool DecrementSynchSem(Mutex* mu, base_internal::PerThreadSynch* w,
                                synchronization_internal::KernelTimeout t);
  void LockSlowLoop(SynchWaitParams* waitp, int flags);
  bool LockSlowWithDeadline(MuHow how, const Condition* cond,
                            synchronization_internal::KernelTimeout t,
                            int flags);
  void LockSlow(MuHow how, const Condition* cond,
                int flags) ABSL_ATTRIBUTE_COLD;
  void UnlockSlow(SynchWaitParams* waitp) ABSL_ATTRIBUTE_COLD;
  bool TryLockSlow();
  bool ReaderTryLockSlow();
  bool AwaitCommon(const Condition& cond,
                   synchronization_internal::KernelTimeout t);
  bool LockWhenCommon(const Condition& cond,
                      synchronization_internal::KernelTimeout t, bool write);
  void TryRemove(base_internal::PerThreadSynch* s);
  void Block(base_internal::PerThreadSynch* s);
  base_internal::PerThreadSynch* Wakeup(base_internal::PerThreadSynch* w);
  void Dtor();
  friend class CondVar;   
  void Trans(MuHow how);  
  void Fer(
      base_internal::PerThreadSynch* w);  
  explicit Mutex(const volatile Mutex* ) {}
  Mutex(const Mutex&) = delete;
  Mutex& operator=(const Mutex&) = delete;
};
class ABSL_SCOPED_LOCKABLE MutexLock {
 public:
  explicit MutexLock(Mutex* mu) ABSL_EXCLUSIVE_LOCK_FUNCTION(mu) : mu_(mu) {
    this->mu_->Lock();
  }
  explicit MutexLock(Mutex* mu, const Condition& cond)
      ABSL_EXCLUSIVE_LOCK_FUNCTION(mu)
      : mu_(mu) {
    this->mu_->LockWhen(cond);
  }
  MutexLock(const MutexLock&) = delete;  
  MutexLock(MutexLock&&) = delete;       
  MutexLock& operator=(const MutexLock&) = delete;
  MutexLock& operator=(MutexLock&&) = delete;
  ~MutexLock() ABSL_UNLOCK_FUNCTION() { this->mu_->Unlock(); }
 private:
  Mutex* const mu_;
};
class ABSL_SCOPED_LOCKABLE ReaderMutexLock {
 public:
  explicit ReaderMutexLock(Mutex* mu) ABSL_SHARED_LOCK_FUNCTION(mu) : mu_(mu) {
    mu->ReaderLock();
  }
  explicit ReaderMutexLock(Mutex* mu, const Condition& cond)
      ABSL_SHARED_LOCK_FUNCTION(mu)
      : mu_(mu) {
    mu->ReaderLockWhen(cond);
  }
  ReaderMutexLock(const ReaderMutexLock&) = delete;
  ReaderMutexLock(ReaderMutexLock&&) = delete;
  ReaderMutexLock& operator=(const ReaderMutexLock&) = delete;
  ReaderMutexLock& operator=(ReaderMutexLock&&) = delete;
  ~ReaderMutexLock() ABSL_UNLOCK_FUNCTION() { this->mu_->ReaderUnlock(); }
 private:
  Mutex* const mu_;
};
class ABSL_SCOPED_LOCKABLE WriterMutexLock {
 public:
  explicit WriterMutexLock(Mutex* mu) ABSL_EXCLUSIVE_LOCK_FUNCTION(mu)
      : mu_(mu) {
    mu->WriterLock();
  }
  explicit WriterMutexLock(Mutex* mu, const Condition& cond)
      ABSL_EXCLUSIVE_LOCK_FUNCTION(mu)
      : mu_(mu) {
    mu->WriterLockWhen(cond);
  }
  WriterMutexLock(const WriterMutexLock&) = delete;
  WriterMutexLock(WriterMutexLock&&) = delete;
  WriterMutexLock& operator=(const WriterMutexLock&) = delete;
  WriterMutexLock& operator=(WriterMutexLock&&) = delete;
  ~WriterMutexLock() ABSL_UNLOCK_FUNCTION() { this->mu_->WriterUnlock(); }
 private:
  Mutex* const mu_;
};
class Condition {
 public:
  Condition(bool (*func)(void*), void* arg);
  template <typename T>
  Condition(bool (*func)(T*), T* arg);
  template <typename T, typename = void>
  Condition(bool (*func)(T*),
            typename absl::internal::type_identity<T>::type* arg);
  template <typename T>
  Condition(T* object,
            bool (absl::internal::type_identity<T>::type::*method)());
  template <typename T>
  Condition(const T* object,
            bool (absl::internal::type_identity<T>::type::*method)() const);
  explicit Condition(const bool* cond);