#ifndef XLA_MLIR_TOOLS_MLIR_INTERPRETER_FRAMEWORK_INTERPRETER_H_
#define XLA_MLIR_TOOLS_MLIR_INTERPRETER_FRAMEWORK_INTERPRETER_H_
#include <cstdint>
#include <functional>
#include <limits>
#include <memory>
#include <optional>
#include "absl/status/statusor.h"
#include "llvm/ADT/StringRef.h"
#include "mlir/Dialect/Func/IR/FuncOps.h"  
#include "mlir/Support/LLVM.h"  
#include "xla/mlir/tools/mlir_interpreter/framework/interpreter_value.h"
namespace mlir {
namespace interpreter {
class InterpreterScope;
class InterpreterListener {
 public:
  virtual ~InterpreterListener() = default;
  virtual void BeforeOp(ArrayRef<InterpreterValue> args, mlir::Operation*) {}
  virtual void AfterOp(ArrayRef<InterpreterValue> results) {}
  virtual void EnterRegion(ArrayRef<InterpreterValue> args,
                           mlir::Region& region) {}
  virtual void LeaveRegion(ArrayRef<InterpreterValue> terminator_args) {}
};
struct InterpreterStats {
  int64_t heap_size = 0;
  int64_t peak_heap_size = 0;
  int64_t num_allocations = 0;
  int64_t num_deallocations = 0;
};
struct InterpreterOptions {
  InterpreterListener* listener = nullptr;
  std::optional<int64_t> max_steps = std::nullopt;
  bool disable_deallocations = false;
  std::function<void(llvm::StringRef)> error_handler =
      [](llvm::StringRef failure) {
        llvm::errs() << "Interpreter failure: " << failure << "\n";
      };
  InterpreterStats* stats = nullptr;
};
class InterpreterState {
 public:
  InterpreterState(const mlir::SymbolTable& symbols,
                   InterpreterOptions options);
  void Step() {
    if (remaining_steps_ == 0) {
      AddFailure("maximum number of steps exceeded");
      return;
    }
    --remaining_steps_;
  }
  void AddFailure(llvm::StringRef failure);
  bool HasFailure() const { return failed_; }
  void CheckSuccess(LogicalResult result, llvm::StringRef failure) {
    if (!result.succeeded()) {
      AddFailure(failure);
    }
  }
  InterpreterScope* GetTopScope() { return top_scope_; }
  const mlir::SymbolTable& GetSymbols() const { return symbols_; }
  const InterpreterOptions& GetOptions() { return options_; }
 private:
  const mlir::SymbolTable& symbols_;
  InterpreterScope* top_scope_ = nullptr;
  bool failed_ = false;
  InterpreterOptions options_;
  int64_t remaining_steps_ = std::numeric_limits<int64_t>::max();
  friend class InterpreterScope;
  friend class InterpreterScopeStash;
};
class InterpreterSideChannel {
 public:
  virtual ~InterpreterSideChannel() = default;
};
class InterpreterScope {
 public:
  InterpreterScope(InterpreterScope&&) = delete;
  explicit InterpreterScope(InterpreterState& state)
      : state_(state), parent_scope_(state.top_scope_) {
    state.top_scope_ = this;
  }
  ~InterpreterScope();
  void Set(Value v, InterpreterValue iv) { values_[v] = std::move(iv); }
  const InterpreterValue& Get(Value v) {
    auto ret = values_.find(v);
    if (ret == values_.end()) {
      if (!parent_scope_) {
        v.dump();
      }
      assert(parent_scope_ && "value not found");
      return parent_scope_->Get(v);
    }
    return ret->second;
  }
  void Verify() const;
  template <typename T>
  T* GetSideChannel(bool optional = false) {
    for (auto& side_channel : side_channels_) {
      if (auto it = dynamic_cast<T*>(side_channel.get())) {
        return it;
      }
    }
    if (!parent_scope_ && optional) return nullptr;
    assert(parent_scope_ && "side channel not found");
    return parent_scope_->GetSideChannel<T>(optional);
  }
  void SetSideChannel(std::shared_ptr<InterpreterSideChannel> side_channel) {
    side_channels_.push_back(std::move(side_channel));
  }
  InterpreterScope* GetParentScope() const { return parent_scope_; }
 private:
  DenseMap<Value, InterpreterValue> values_;
  SmallVector<std::shared_ptr<InterpreterSideChannel>> side_channels_;
  InterpreterState& state_;
  InterpreterScope* parent_scope_;
  friend class InterpreterScopeStash;
};
SmallVector<InterpreterValue> Interpret(InterpreterState& state, Region& region,
                                        ArrayRef<InterpreterValue> bbargs);
absl::StatusOr<SmallVector<InterpreterValue>> RunInterpreter(
    const mlir::SymbolTable& symbols, mlir::func::FuncOp function,
    ArrayRef<InterpreterValue> args, InterpreterOptions options = {});
}  
}  
#endif  
#include "xla/mlir/tools/mlir_interpreter/framework/interpreter.h"
#include <cassert>
#include <functional>
#include <memory>
#include <optional>
#include <utility>
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "llvm/ADT/STLExtras.h"
#include "llvm/ADT/SmallVector.h"
#include "llvm/ADT/StringRef.h"
#include "llvm/Support/ErrorHandling.h"
#include "llvm/Support/raw_ostream.h"
#include "mlir/Dialect/Bufferization/IR/BufferizableOpInterface.h"  
#include "mlir/IR/OpDefinition.h"  
#include "mlir/IR/Operation.h"  
#include "mlir/IR/Region.h"  
#include "mlir/IR/SymbolTable.h"  
#include "mlir/Support/LLVM.h"  
#include "xla/mlir/tools/mlir_interpreter/framework/interpreter_value.h"
#include "xla/mlir/tools/mlir_interpreter/framework/registration.h"
namespace mlir {
namespace interpreter {
SmallVector<InterpreterValue> Interpret(InterpreterState& state,
                                        Operation& op) {
  auto fn = detail::GetFunction(op.getName().getStringRef());
  if (!fn) {
    llvm::errs() << "Unsupported op: " << op.getName().getStringRef() << "\n";
    op.dump();
    state.AddFailure("unsupported op");
    return {};
  }
  SmallVector<InterpreterValue> operands;
  for (auto operand : op.getOperands()) {
    operands.push_back(state.GetTopScope()->Get(operand));
  }
  state.GetOptions().listener->BeforeOp(operands, &op);
  auto results = fn(operands, &op, state);
  for (auto* scope = state.GetTopScope(); scope != nullptr;
       scope = scope->GetParentScope()) {
    scope->Verify();
  }
  if (state.HasFailure()) {
    llvm::errs() << "Encountered failure while executing " << op << "\n";
  }
  state.GetOptions().listener->AfterOp(results);
  state.Step();
  return results;
}
SmallVector<InterpreterValue> Interpret(InterpreterState& state, Region& region,
                                        ArrayRef<InterpreterValue> bbargs) {
  if (state.HasFailure()) return {};
  assert(region.hasOneBlock() && "expected region to have one block");
  state.GetOptions().listener->EnterRegion(bbargs, region);
  InterpreterScope scope(state);
  auto& block = region.getBlocks().front();
  for (auto [value, interpreter_value] :
       llvm::zip(block.getArguments(), bbargs)) {
    scope.Set(value, interpreter_value);
  }
  std::optional<SmallVector<InterpreterValue>> block_results;
  for (mlir::Operation& op : block) {
    auto results = Interpret(state, op);
    if (state.HasFailure()) return {};
    if (op.hasTrait<OpTrait::IsTerminator>()) {
      assert(!block_results.has_value() && "Expected at most one terminator");
      block_results = results;
    } else {
      if (results.size() != op.getNumResults()) {
        llvm::errs() << "Unexpected number of results while interpreting "
                     << op.getName().getStringRef() << ". Interpreter bug?\n";
        llvm_unreachable("unexpected number of results");
      }
      for (auto [v, iv] : llvm::zip(op.getResults(), results)) {
        scope.Set(v, iv);
      }
    }
  }
  if (!block_results) {
    block_results = SmallVector<InterpreterValue>{};
  }
  state.GetOptions().listener->LeaveRegion(*block_results);
  return *std::move(block_results);
}
InterpreterState::InterpreterState(const mlir::SymbolTable& symbols,
                                   InterpreterOptions options)
    : symbols_(symbols), options_(options) {
  if (!options_.listener) {
    static auto& no_op_listener = *new InterpreterListener();
    this->options_.listener = &no_op_listener;
  }
  if (options_.max_steps) {
    remaining_steps_ = *options_.max_steps;
  }
}
void InterpreterState::AddFailure(llvm::StringRef failure) {
  failed_ = true;
  options_.error_handler(failure);
}
void InterpreterScope::Verify() const {
  for (auto& [_, value] : values_) {
    if (value.IsTensor() && value.GetBuffer() &&
        !value.GetBuffer()->GetFailure().empty()) {
      state_.AddFailure(value.GetBuffer()->GetFailure());
      break;
    }
  }
}
InterpreterScope::~InterpreterScope() {
  Verify();
  state_.top_scope_ = parent_scope_;
}
absl::StatusOr<SmallVector<InterpreterValue>> RunInterpreter(
    const mlir::SymbolTable& symbols, mlir::func::FuncOp function,
    ArrayRef<InterpreterValue> args, InterpreterOptions options) {
  InterpreterState state{symbols, std::move(options)};
  auto results = Interpret(state, function.getBody(), args);
  if (state.HasFailure()) {
    return absl::InvalidArgumentError("Interpreter failed, check error logs");
  }
  return results;
}
}  
}  