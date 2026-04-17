#ifndef THIRD_PARTY_CEL_CPP_EXTENSIONS_PROTOBUF_VALUE_H_
#define THIRD_PARTY_CEL_CPP_EXTENSIONS_PROTOBUF_VALUE_H_
#include <type_traits>
#include <utility>
#include "google/protobuf/duration.pb.h"
#include "google/protobuf/struct.pb.h"
#include "google/protobuf/timestamp.pb.h"
#include "google/protobuf/wrappers.pb.h"
#include "absl/base/attributes.h"
#include "absl/base/nullability.h"
#include "absl/functional/overload.h"
#include "absl/status/statusor.h"
#include "common/value.h"
#include "common/value_factory.h"
#include "common/value_manager.h"
#include "extensions/protobuf/internal/duration_lite.h"
#include "extensions/protobuf/internal/enum.h"
#include "extensions/protobuf/internal/message.h"
#include "extensions/protobuf/internal/struct_lite.h"
#include "extensions/protobuf/internal/timestamp_lite.h"
#include "extensions/protobuf/internal/wrappers_lite.h"
#include "internal/status_macros.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/generated_enum_reflection.h"
namespace cel::extensions {
template <typename T>
std::enable_if_t<protobuf_internal::IsProtoEnum<T>, absl::StatusOr<ValueView>>
ProtoEnumToValue(ValueFactory&, T value,
                 Value& scratch ABSL_ATTRIBUTE_LIFETIME_BOUND
                     ABSL_ATTRIBUTE_UNUSED) {
  if constexpr (std::is_same_v<T, google::protobuf::NullValue>) {
    return NullValueView{};
  }
  return IntValueView{static_cast<int>(value)};
}
template <typename T>
std::enable_if_t<protobuf_internal::IsProtoEnum<T>, absl::StatusOr<Value>>
ProtoEnumToValue(ValueFactory&, T value) {
  if constexpr (std::is_same_v<T, google::protobuf::NullValue>) {
    return NullValue{};
  }
  return IntValue{static_cast<int>(value)};
}
absl::StatusOr<int> ProtoEnumFromValue(
    ValueView value, absl::Nonnull<const google::protobuf::EnumDescriptor*> desc);
template <typename T>
std::enable_if_t<protobuf_internal::IsProtoEnum<T>, absl::StatusOr<T>>
ProtoEnumFromValue(ValueView value) {
  CEL_ASSIGN_OR_RETURN(
      auto enum_value,
      ProtoEnumFromValue(value, google::protobuf::GetEnumDescriptor<T>()));
  return static_cast<T>(enum_value);
}
template <typename T>
std::enable_if_t<protobuf_internal::IsProtoMessage<T>, absl::StatusOr<Value>>
ProtoMessageToValue(ValueManager& value_manager, T&& value) {
  using Tp = std::decay_t<T>;
  if constexpr (std::is_same_v<Tp, google::protobuf::BoolValue>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedBoolValueProto(
                             std::forward<T>(value)));
    return BoolValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::Int32Value>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedInt32ValueProto(
                             std::forward<T>(value)));
    return IntValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::Int64Value>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedInt64ValueProto(
                             std::forward<T>(value)));
    return IntValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::UInt32Value>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedUInt32ValueProto(
                             std::forward<T>(value)));
    return UintValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::UInt64Value>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedUInt64ValueProto(
                             std::forward<T>(value)));
    return UintValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::FloatValue>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedFloatValueProto(
                             std::forward<T>(value)));
    return DoubleValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::DoubleValue>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedDoubleValueProto(
                             std::forward<T>(value)));
    return DoubleValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::BytesValue>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedBytesValueProto(
                             std::forward<T>(value)));
    return BytesValue{std::move(result)};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::StringValue>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedStringValueProto(
                             std::forward<T>(value)));
    return StringValue{std::move(result)};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::Duration>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedDurationProto(
                             std::forward<T>(value)));
    return DurationValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::Timestamp>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::UnwrapGeneratedTimestampProto(
                             std::forward<T>(value)));
    return TimestampValue{result};
  } else if constexpr (std::is_same_v<Tp, google::protobuf::Value>) {
    CEL_ASSIGN_OR_RETURN(
        auto result,
        protobuf_internal::GeneratedValueProtoToJson(std::forward<T>(value)));
    return value_manager.CreateValueFromJson(std::move(result));
  } else if constexpr (std::is_same_v<Tp, google::protobuf::ListValue>) {
    CEL_ASSIGN_OR_RETURN(auto result,
                         protobuf_internal::GeneratedListValueProtoToJson(
                             std::forward<T>(value)));
    return value_manager.CreateListValueFromJsonArray(std::move(result));
  } else if constexpr (std::is_same_v<Tp, google::protobuf::Struct>) {
    CEL_ASSIGN_OR_RETURN(
        auto result,
        protobuf_internal::GeneratedStructProtoToJson(std::forward<T>(value)));
    return value_manager.CreateMapValueFromJsonObject(std::move(result));
  } else {
    auto dispatcher = absl::Overload(
        [&](Tp&& m) {
          return protobuf_internal::ProtoMessageToValueImpl(
              value_manager, &m, sizeof(T), alignof(T),
              &protobuf_internal::ProtoMessageTraits<Tp>::ArenaMoveConstruct,
              &protobuf_internal::ProtoMessageTraits<Tp>::MoveConstruct);
        },
        [&](const Tp& m) {
          return protobuf_internal::ProtoMessageToValueImpl(
              value_manager, &m, sizeof(T), alignof(T),
              &protobuf_internal::ProtoMessageTraits<Tp>::ArenaCopyConstruct,
              &protobuf_internal::ProtoMessageTraits<Tp>::CopyConstruct);
        });
    return dispatcher(std::forward<T>(value));
  }
}
template <typename T>
std::enable_if_t<protobuf_internal::IsProtoMessage<T>,
                 absl::StatusOr<ValueView>>
ProtoMessageToValue(ValueManager& value_manager, T&& value,
                    Value& scratch ABSL_ATTRIBUTE_LIFETIME_BOUND) {
  CEL_ASSIGN_OR_RETURN(
      scratch, ProtoMessageToValue(value_manager, std::forward<T>(value)));
  return scratch;
}
template <typename T>
std::enable_if_t<protobuf_internal::IsProtoMessage<T>, absl::Status>
ProtoMessageFromValue(ValueView value, T& message,
                      absl::Nonnull<const google::protobuf::DescriptorPool*> pool,
                      absl::Nonnull<google::protobuf::MessageFactory*> factory) {
  return protobuf_internal::ProtoMessageFromValueImpl(value, pool, factory,
                                                      &message);
}
template <typename T>
std::enable_if_t<protobuf_internal::IsProtoMessage<T>, absl::Status>
ProtoMessageFromValue(ValueView value, T& message) {
  return protobuf_internal::ProtoMessageFromValueImpl(value, &message);
}
inline absl::StatusOr<absl::Nonnull<google::protobuf::Message*>> ProtoMessageFromValue(
    ValueView value, absl::Nullable<google::protobuf::Arena*> arena) {
  return protobuf_internal::ProtoMessageFromValueImpl(value, arena);
}
inline absl::StatusOr<absl::Nonnull<google::protobuf::Message*>> ProtoMessageFromValue(
    ValueView value, absl::Nullable<google::protobuf::Arena*> arena,
    absl::Nonnull<const google::protobuf::DescriptorPool*> pool,
    absl::Nonnull<google::protobuf::MessageFactory*> factory) {
  return protobuf_internal::ProtoMessageFromValueImpl(value, pool, factory,
                                                      arena);
}
}  
#endif  
#include "extensions/protobuf/value.h"
#include <limits>
#include "absl/base/nullability.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "common/casting.h"
#include "common/value.h"
#include "google/protobuf/descriptor.h"
namespace cel::extensions {
absl::StatusOr<int> ProtoEnumFromValue(
    ValueView value, absl::Nonnull<const google::protobuf::EnumDescriptor*> desc) {
  if (desc->full_name() == "google.protobuf.NullValue") {
    if (InstanceOf<NullValueView>(value) || InstanceOf<IntValueView>(value)) {
      return 0;
    }
    return TypeConversionError(value.GetTypeName(), desc->full_name())
        .NativeValue();
  }
  if (auto int_value = As<IntValueView>(value); int_value) {
    if (int_value->NativeValue() >= 0 &&
        int_value->NativeValue() <= std::numeric_limits<int>::max()) {
      const auto* value_desc =
          desc->FindValueByNumber(static_cast<int>(int_value->NativeValue()));
      if (value_desc != nullptr) {
        return value_desc->number();
      }
    }
    return absl::NotFoundError(absl::StrCat("enum `", desc->full_name(),
                                            "` has no value with number ",
                                            int_value->NativeValue()));
  }
  return TypeConversionError(value.GetTypeName(), desc->full_name())
      .NativeValue();
}
}  