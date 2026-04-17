#ifndef XLA_PYTHON_IFRT_DTYPE_H_
#define XLA_PYTHON_IFRT_DTYPE_H_
#include <optional>
#include <ostream>
#include <string>
#include "absl/status/statusor.h"
#include "xla/python/ifrt/dtype.pb.h"
namespace xla {
namespace ifrt {
class DType {
 public:
  enum Kind {
    kInvalid = 0,
    kPred = 1,
    kS2 = 26,
    kS4 = 21,
    kS8 = 2,
    kS16 = 3,
    kS32 = 4,
    kS64 = 5,
    kU2 = 27,
    kU4 = 22,
    kU8 = 6,
    kU16 = 7,
    kU32 = 8,
    kU64 = 9,
    kF16 = 10,
    kF32 = 11,
    kF64 = 12,
    kBF16 = 16,
    kC64 = 15,   
    kC128 = 18,  
    kToken = 17,
    kF8E4M3FN = 20,
    kF8E4M3B11FNUZ = 23,
    kF8E4M3FNUZ = 25,
    kF8E5M2 = 19,
    kF8E5M2FNUZ = 24,
    kString = 99,
  };
  explicit DType(Kind kind) : kind_(kind) {}
  DType(const DType&) = default;
  DType(DType&&) = default;
  DType& operator=(const DType&) = default;
  DType& operator=(DType&&) = default;
  Kind kind() const { return kind_; }
  bool operator==(const DType& other) const { return kind_ == other.kind_; }
  bool operator!=(const DType& other) const { return kind_ != other.kind_; }
  template <typename H>
  friend H AbslHashValue(H h, const DType& value) {
    return H::combine(std::move(h), value.kind());
  }
  std::optional<int> byte_size() const;
  std::optional<int> bit_size() const;
  static absl::StatusOr<DType> FromProto(const DTypeProto& proto);
  DTypeProto ToProto() const;
  std::string DebugString() const;
 private:
  Kind kind_;
};
std::ostream& operator<<(std::ostream& os, const DType& dtype);
}  
}  
#endif  
#include "xla/python/ifrt/dtype.h"
#include <optional>
#include <ostream>
#include <string>
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "xla/python/ifrt/dtype.pb.h"
namespace xla {
namespace ifrt {
std::optional<int> DType::byte_size() const {
  switch (kind_) {
    case kPred:
    case kS8:
    case kU8:
      return 1;
    case kS16:
    case kU16:
    case kF16:
    case kBF16:
      return 2;
    case kS32:
    case kU32:
    case kF32:
      return 4;
    case kS64:
    case kU64:
    case kF64:
    case kC64:
      return 8;
    case kC128:
      return 16;
    default:
      return std::nullopt;
  }
}
std::optional<int> DType::bit_size() const {
  switch (kind_) {
    case kPred:
    case kS8:
    case kU8:
      return 8;
    case kS16:
    case kU16:
    case kF16:
    case kBF16:
      return 16;
    case kS32:
    case kU32:
    case kF32:
      return 32;
    case kS64:
    case kU64:
    case kF64:
    case kC64:
      return 64;
    case kC128:
      return 128;
    default:
      return std::nullopt;
  }
}
absl::StatusOr<DType> DType::FromProto(const DTypeProto& dtype_proto) {
  switch (dtype_proto.kind()) {
    case DTypeProto::KIND_PRED:
      return DType(DType::Kind::kPred);
    case DTypeProto::KIND_TOKEN:
      return DType(DType::Kind::kToken);
#define CASE(X)              \
  case DTypeProto::KIND_##X: \
    return DType(DType::Kind::k##X);
      CASE(S4);
      CASE(S8);
      CASE(S16);
      CASE(S32);
      CASE(S64);
      CASE(U4);
      CASE(U8);
      CASE(U16);
      CASE(U32);
      CASE(U64);
      CASE(F16);
      CASE(F32);
      CASE(F64);
      CASE(BF16);
      CASE(C64);
      CASE(C128);
      CASE(F8E4M3FN);
      CASE(F8E4M3B11FNUZ);
      CASE(F8E4M3FNUZ);
      CASE(F8E5M2);
      CASE(F8E5M2FNUZ);
#undef CASE
    case DTypeProto::KIND_STRING:
      return DType(DType::Kind::kString);
    default:
      return DType(DType::Kind::kInvalid);
  }
}
DTypeProto DType::ToProto() const {
  DTypeProto dtype_proto;
  switch (kind()) {
    case DType::Kind::kPred:
      dtype_proto.set_kind(DTypeProto::KIND_PRED);
      break;
    case DType::Kind::kToken:
      dtype_proto.set_kind(DTypeProto::KIND_TOKEN);
      break;
#define CASE(X)                                 \
  case DType::Kind::k##X:                       \
    dtype_proto.set_kind(DTypeProto::KIND_##X); \
    break;
      CASE(S4);
      CASE(S8);
      CASE(S16);
      CASE(S32);
      CASE(S64);
      CASE(U4);
      CASE(U8);
      CASE(U16);
      CASE(U32);
      CASE(U64);
      CASE(F16);
      CASE(F32);
      CASE(F64);
      CASE(BF16);
      CASE(C64);
      CASE(C128);
      CASE(F8E4M3FN);
      CASE(F8E4M3B11FNUZ);
      CASE(F8E4M3FNUZ);
      CASE(F8E5M2);
      CASE(F8E5M2FNUZ);
#undef CASE
    case DType::Kind::kString:
      dtype_proto.set_kind(DTypeProto::KIND_STRING);
      break;
    default:
      dtype_proto.set_kind(DTypeProto::KIND_UNSPECIFIED);
      break;
  }
  return dtype_proto;
}
std::string DType::DebugString() const {
  switch (kind_) {
    case kInvalid:
      return "INVALID";
    case kPred:
      return "PRED";
    case kS8:
      return "S8";
    case kS16:
      return "S16";
    case kS32:
      return "S32";
    case kS64:
      return "S64";
    case kU8:
      return "U8";
    case kU16:
      return "U16";
    case kU32:
      return "U32";
    case kU64:
      return "U64";
    case kF16:
      return "F16";
    case kF32:
      return "F32";
    case kF64:
      return "F64";
    case kBF16:
      return "BF16";
    case kC64:
      return "C64";
    case kC128:
      return "C128";
    case kToken:
      return "TOKEN";
    case kString:
      return "STRING";
    default:
      return absl::StrCat("UNKNOWN(", static_cast<int>(kind_), ")");
  }
}
std::ostream& operator<<(std::ostream& os, const DType& dtype) {
  return os << dtype.DebugString();
}
}  
}  