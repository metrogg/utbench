#ifndef AROLLA_IO_STRUCT_IO_H_
#define AROLLA_IO_STRUCT_IO_H_
#include <cstddef>
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/base/nullability.h"
#include "absl/container/flat_hash_map.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "arolla/io/input_loader.h"
#include "arolla/io/slot_listener.h"
#include "arolla/memory/frame.h"
#include "arolla/memory/raw_buffer_factory.h"
#include "arolla/qtype/qtype.h"
#include "arolla/qtype/typed_slot.h"
#include "arolla/util/status_macros_backport.h"
namespace arolla {
namespace struct_io_impl {
class StructIO {
 public:
  explicit StructIO(
      const absl::flat_hash_map<std::string, TypedSlot>& struct_slots,
      const absl::flat_hash_map<std::string, TypedSlot>& frame_slots);
  void CopyStructToFrame(const void* struct_ptr, FramePtr frame) const;
  void CopyFrameToStruct(ConstFramePtr frame, void* struct_ptr) const;
 private:
  using Offsets = std::vector<std::pair<size_t, size_t>>;
  Offsets offsets_bool_;
  Offsets offsets_32bits_;
  Offsets offsets_64bits_;
  absl::flat_hash_map<QTypePtr, Offsets> offsets_other_;
};
std::vector<std::string> SuggestAvailableNames(
    const absl::flat_hash_map<std::string, TypedSlot>& slots);
absl::Status ValidateStructSlots(
    const absl::flat_hash_map<std::string, TypedSlot>& slots,
    size_t struct_size);
}  
template <typename T>
class StructInputLoader final : public InputLoader<T> {
 public:
  static absl::StatusOr<InputLoaderPtr<T>> Create(
      absl::flat_hash_map<std::string, TypedSlot> struct_slots) {
    RETURN_IF_ERROR(
        struct_io_impl::ValidateStructSlots(struct_slots, sizeof(T)));
    return InputLoaderPtr<T>(new StructInputLoader(std::move(struct_slots)));
  }
  absl::Nullable<const QType*> GetQTypeOf(absl::string_view name) const final {
    auto it = struct_slots_.find(name);
    return it != struct_slots_.end() ? it->second.GetType() : nullptr;
  }
  std::vector<std::string> SuggestAvailableNames() const final {
    return struct_io_impl::SuggestAvailableNames(struct_slots_);
  }
 private:
  explicit StructInputLoader(
      absl::flat_hash_map<std::string, TypedSlot> struct_slots)
      : struct_slots_(std::move(struct_slots)) {}
  absl::StatusOr<BoundInputLoader<T>> BindImpl(
      const absl::flat_hash_map<std::string, TypedSlot>& slots) const final {
    return BoundInputLoader<T>(
        [io = struct_io_impl::StructIO(struct_slots_, slots)](
            const T& input, FramePtr frame, RawBufferFactory*) -> absl::Status {
          io.CopyStructToFrame(&input, frame);
          return absl::OkStatus();
        });
  }
  absl::flat_hash_map<std::string, TypedSlot> struct_slots_;
};
template <typename T>
class StructSlotListener final : public SlotListener<T> {
 public:
  static absl::StatusOr<std::unique_ptr<SlotListener<T>>> Create(
      absl::flat_hash_map<std::string, TypedSlot> struct_slots) {
    RETURN_IF_ERROR(
        struct_io_impl::ValidateStructSlots(struct_slots, sizeof(T)));
    return std::unique_ptr<SlotListener<T>>(
        new StructSlotListener(std::move(struct_slots)));
  }
  absl::Nullable<const QType*> GetQTypeOf(
      absl::string_view name, absl::Nullable<const QType*>) const final {
    auto it = struct_slots_.find(name);
    return it != struct_slots_.end() ? it->second.GetType() : nullptr;
  }
  std::vector<std::string> SuggestAvailableNames() const final {
    return struct_io_impl::SuggestAvailableNames(struct_slots_);
  }
 private:
  explicit StructSlotListener(
      absl::flat_hash_map<std::string, TypedSlot> struct_slots)
      : struct_slots_(std::move(struct_slots)) {}
  absl::StatusOr<BoundSlotListener<T>> BindImpl(
      const absl::flat_hash_map<std::string, TypedSlot>& slots) const final {
    return BoundSlotListener<T>(
        [io = struct_io_impl::StructIO(struct_slots_, slots)](
            ConstFramePtr frame, T* output) -> absl::Status {
          io.CopyFrameToStruct(frame, output);
          return absl::OkStatus();
        });
  }
  absl::flat_hash_map<std::string, TypedSlot> struct_slots_;
};
}  
#endif  
#include "arolla/io/struct_io.h"
#include <algorithm>
#include <cstddef>
#include <cstdint>
#include <cstring>
#include <string>
#include <type_traits>
#include <utility>
#include <vector>
#include "absl/algorithm/container.h"
#include "absl/container/flat_hash_map.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "arolla/memory/frame.h"
#include "arolla/memory/optional_value.h"
#include "arolla/qtype/base_types.h"
#include "arolla/qtype/optional_qtype.h"
#include "arolla/qtype/qtype.h"
#include "arolla/qtype/qtype_traits.h"
#include "arolla/qtype/typed_slot.h"
namespace arolla::struct_io_impl {
std::vector<std::string> SuggestAvailableNames(
    const absl::flat_hash_map<std::string, TypedSlot>& slots) {
  std::vector<std::string> names;
  names.reserve(slots.size());
  for (const auto& [name, _] : slots) {
    names.emplace_back(name);
  }
  return names;
}
absl::Status ValidateStructSlots(
    const absl::flat_hash_map<std::string, TypedSlot>& slots,
    size_t struct_size) {
  for (const auto& [name, slot] : slots) {
    if (slot.byte_offset() + slot.GetType()->type_layout().AllocSize() >
        struct_size) {
      return absl::InvalidArgumentError(
          absl::StrCat("slot '", name, "' is not within the struct"));
    }
  }
  return absl::OkStatus();
}
StructIO::StructIO(
    const absl::flat_hash_map<std::string, TypedSlot>& struct_slots,
    const absl::flat_hash_map<std::string, TypedSlot>& frame_slots) {
  QTypePtr b = GetQType<bool>();
  std::vector<QTypePtr> types32{GetQType<float>(), GetQType<int32_t>()};
  std::vector<QTypePtr> types64{GetQType<double>(), GetQType<int64_t>(),
                                GetQType<uint64_t>(), GetOptionalQType<float>(),
                                GetOptionalQType<int32_t>()};
  static_assert(sizeof(OptionalValue<float>) == 8);
  static_assert(sizeof(OptionalValue<int32_t>) == 8);
  static_assert(std::is_trivially_copyable_v<OptionalValue<float>>);
  static_assert(std::is_trivially_copyable_v<OptionalValue<int32_t>>);
  for (const auto& [name, frame_slot] : frame_slots) {
    QTypePtr t = frame_slot.GetType();
    size_t struct_offset = struct_slots.at(name).byte_offset();
    size_t frame_offset = frame_slot.byte_offset();
    if (t == b) {
      offsets_bool_.emplace_back(struct_offset, frame_offset);
    } else if (absl::c_find(types32, t) != types32.end()) {
      DCHECK_EQ(t->type_layout().AllocSize(), 4);
      offsets_32bits_.emplace_back(struct_offset, frame_offset);
    } else if (absl::c_find(types64, t) != types64.end()) {
      DCHECK_EQ(t->type_layout().AllocSize(), 8);
      offsets_64bits_.emplace_back(struct_offset, frame_offset);
    } else {
      offsets_other_[t].emplace_back(struct_offset, frame_offset);
    }
  }
  std::sort(offsets_bool_.begin(), offsets_bool_.end());
  std::sort(offsets_32bits_.begin(), offsets_32bits_.end());
  std::sort(offsets_64bits_.begin(), offsets_64bits_.end());
  for (auto& [_, v] : offsets_other_) {
    std::sort(v.begin(), v.end());
  }
}
void StructIO::CopyStructToFrame(const void* struct_ptr, FramePtr frame) const {
  const char* src_base = reinterpret_cast<const char*>(struct_ptr);
  for (const auto& [src, dst] : offsets_bool_) {
    std::memcpy(frame.GetRawPointer(dst), src_base + src, sizeof(bool));
  }
  for (const auto& [src, dst] : offsets_32bits_) {
    std::memcpy(frame.GetRawPointer(dst), src_base + src, 4);
  }
  for (const auto& [src, dst] : offsets_64bits_) {
    std::memcpy(frame.GetRawPointer(dst), src_base + src, 8);
  }
  for (const auto& [t, offsets] : offsets_other_) {
    for (const auto& [src, dst] : offsets) {
      t->UnsafeCopy(src_base + src, frame.GetRawPointer(dst));
    }
  }
}
void StructIO::CopyFrameToStruct(ConstFramePtr frame, void* struct_ptr) const {
  char* dst_base = reinterpret_cast<char*>(struct_ptr);
  for (const auto& [dst, src] : offsets_bool_) {
    std::memcpy(dst_base + dst, frame.GetRawPointer(src), sizeof(bool));
  }
  for (const auto& [dst, src] : offsets_32bits_) {
    std::memcpy(dst_base + dst, frame.GetRawPointer(src), 4);
  }
  for (const auto& [dst, src] : offsets_64bits_) {
    std::memcpy(dst_base + dst, frame.GetRawPointer(src), 8);
  }
  for (const auto& [t, offsets] : offsets_other_) {
    for (const auto& [dst, src] : offsets) {
      t->UnsafeCopy(frame.GetRawPointer(src), dst_base + dst);
    }
  }
}
}  