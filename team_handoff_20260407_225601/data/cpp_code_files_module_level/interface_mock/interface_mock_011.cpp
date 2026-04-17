#ifndef AROLLA_SERVING_INPLACE_EXPR_COMPILER_H_
#define AROLLA_SERVING_INPLACE_EXPR_COMPILER_H_
#include <functional>
#include <memory>
#include <string>
#include <type_traits>
#include <utility>
#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "arolla/expr/expr_node.h"
#include "arolla/io/input_loader.h"
#include "arolla/io/slot_listener.h"
#include "arolla/io/struct_io.h"
#include "arolla/memory/frame.h"
#include "arolla/qexpr/eval_context.h"
#include "arolla/qexpr/evaluation_engine.h"
#include "arolla/qtype/qtype.h"
#include "arolla/qtype/qtype_traits.h"
#include "arolla/qtype/typed_slot.h"
#include "arolla/util/status_macros_backport.h"
namespace arolla {
namespace inplace_expr_compiler_impl {
using TypedSlotMap = absl::flat_hash_map<std::string, TypedSlot>;
TypedSlotMap CollectInternalSlots(TypedSlot root_slot);
struct IoSlots {
  TypedSlotMap input_slots;
  TypedSlot output_slot;
  TypedSlotMap named_output_slots;
};
absl::StatusOr<IoSlots> CollectIoSlots(QTypePtr qtype,
                                       const CompiledExpr& compiled_expr,
                                       absl::string_view final_output_name);
}  
template <class T>
using InplaceModelFunction = std::function<absl::Status(T&)>;
template <typename T>
absl::StatusOr<InplaceModelFunction<T>> CompileInplaceExprOnStruct(
    const InplaceCompiledExpr& compiled_expr,
    absl::string_view final_output_name) {
  static_assert(
      std::is_standard_layout<T>::value,
      "Data must be standard layout to be used with CompileExprInplace.");
  QTypePtr qtype = GetQType<T>();
  ASSIGN_OR_RETURN(inplace_expr_compiler_impl::IoSlots slots,
                   inplace_expr_compiler_impl::CollectIoSlots(
                       qtype, compiled_expr, final_output_name));
  ASSIGN_OR_RETURN(auto executable, compiled_expr.InplaceBind(
                                        slots.input_slots, slots.output_slot,
                                        slots.named_output_slots));
  return [qtype, executable(std::shared_ptr<BoundExpr>(std::move(executable)))](
             T& input) -> absl::Status {
    FramePtr frame(&input, &qtype->type_layout());
    EvaluationContext ctx;
    executable->Execute(&ctx, frame);
    return ctx.status();
  };
}
template <typename Struct>
absl::StatusOr<InputLoaderPtr<Struct>> CreateStructInputLoader() {
  return StructInputLoader<Struct>::Create(
      inplace_expr_compiler_impl::CollectInternalSlots(
          TypedSlot::UnsafeFromOffset(GetQType<Struct>(), 0)));
}
template <typename Struct>
absl::StatusOr<std::unique_ptr<SlotListener<Struct>>>
CreateStructSlotListener() {
  return StructSlotListener<Struct>::Create(
      inplace_expr_compiler_impl::CollectInternalSlots(
          TypedSlot::UnsafeFromOffset(GetQType<Struct>(), 0)));
}
}  
#endif  
#include "arolla/serving/inplace_expr_compiler.h"
#include <cstddef>
#include <string>
#include <utility>
#include <vector>
#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "arolla/naming/table.h"
#include "arolla/qexpr/evaluation_engine.h"
#include "arolla/qtype/named_field_qtype.h"
#include "arolla/qtype/qtype.h"
#include "arolla/qtype/typed_slot.h"
#include "arolla/util/status_macros_backport.h"
namespace arolla::inplace_expr_compiler_impl {
TypedSlotMap CollectInternalSlots(TypedSlot root_slot) {
  TypedSlotMap result;
  if (GetFieldNames(root_slot.GetType()).empty()) {
    return result;
  }
  std::vector<std::pair<TypedSlot, naming::TablePath>> stack{{root_slot, {}}};
  while (!stack.empty()) {
    auto [slot, table] = stack.back();
    stack.pop_back();
    auto field_names = GetFieldNames(slot.GetType());
    for (size_t i = 0; i < field_names.size(); ++i) {
      const auto& field_name = field_names[i];
      const TypedSlot& field_slot = slot.SubSlot(i);
      result.emplace(table.Column(naming::FieldAccess(field_name)).FullName(),
                     field_slot);
      if (!GetFieldNames(field_slot.GetType()).empty()) {
        stack.emplace_back(field_slot,
                           table.Child(naming::FieldAccess(field_name)));
      }
    }
  }
  return result;
}
namespace {
absl::Status CheckField(QTypePtr qtype, const TypedSlotMap& slot_map,
                        QTypePtr field_qtype, absl::string_view field_name) {
  if (GetFieldNames(qtype).empty()) {
    return absl::FailedPreconditionError(
        absl::StrCat("no registered field names for ", qtype->name(),
                     " in Compile.*ExprOnStructInput"));
  }
  if (!slot_map.contains(field_name)) {
    return absl::FailedPreconditionError(
        absl::StrCat("input `", field_name, "` not found in ", qtype->name(),
                     " in Compile.*ExprOnStructInput"));
  }
  QTypePtr result_type = slot_map.at(field_name).GetType();
  if (result_type != field_qtype) {
    return absl::FailedPreconditionError(absl::StrCat(
        "input `", field_name, "` type mismatch for ", qtype->name(),
        " in Compile.*ExprOnStructInput, expected in struct: ",
        result_type->name(), ", found in expr: ", field_qtype->name()));
  }
  return absl::OkStatus();
}
absl::StatusOr<TypedSlotMap> CollectInputSlots(
    QTypePtr qtype, const TypedSlotMap& struct_slot_map,
    const CompiledExpr& compiled_expr) {
  TypedSlotMap input_slots;
  input_slots.reserve(compiled_expr.input_types().size());
  for (const auto& [name, field_qtype] : compiled_expr.input_types()) {
    RETURN_IF_ERROR(CheckField(qtype, struct_slot_map, field_qtype, name));
    input_slots.emplace(name, struct_slot_map.at(name));
  }
  return input_slots;
}
}  
absl::StatusOr<IoSlots> CollectIoSlots(QTypePtr qtype,
                                       const CompiledExpr& compiled_expr,
                                       absl::string_view final_output_name) {
  TypedSlotMap struct_slot_map =
      CollectInternalSlots(TypedSlot::UnsafeFromOffset(qtype, 0));
  ASSIGN_OR_RETURN(TypedSlotMap input_slots,
                   CollectInputSlots(qtype, struct_slot_map, compiled_expr));
  RETURN_IF_ERROR(CheckField(qtype, struct_slot_map,
                             compiled_expr.output_type(), final_output_name));
  if (compiled_expr.input_types().contains(final_output_name)) {
    return absl::FailedPreconditionError(absl::StrCat(
        final_output_name, " present both as an input and as final output"));
  }
  if (compiled_expr.named_output_types().contains(final_output_name)) {
    return absl::FailedPreconditionError(
        absl::StrCat(final_output_name,
                     " present both as final output and as named output"));
  }
  for (const auto& [name, field_qtype] : compiled_expr.input_types()) {
    if (compiled_expr.named_output_types().contains(name)) {
      return absl::FailedPreconditionError(
          absl::StrCat(name, " present both as an input and as named output"));
    }
  }
  for (const auto& [name, field_qtype] : compiled_expr.named_output_types()) {
    RETURN_IF_ERROR(CheckField(qtype, struct_slot_map, field_qtype, name));
  }
  absl::flat_hash_map<std::string, TypedSlot> named_output_slots;
  named_output_slots.reserve(compiled_expr.named_output_types().size());
  for (const auto& [name, _] : compiled_expr.named_output_types()) {
    named_output_slots.emplace(name, struct_slot_map.at(name));
  }
  return IoSlots{.input_slots = input_slots,
                 .output_slot = struct_slot_map.at(final_output_name),
                 .named_output_slots = named_output_slots};
}
}  