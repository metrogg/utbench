#ifndef XLA_PYTHON_IFRT_SERDES_H_
#define XLA_PYTHON_IFRT_SERDES_H_
#include <memory>
#include <string>
#include <utility>
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "llvm/Support/Casting.h"
#include "llvm/Support/ExtensibleRTTI.h"
#include "xla/python/ifrt/serdes.pb.h"
#include "tsl/platform/statusor.h"
namespace xla {
namespace ifrt {
struct DeserializeOptions
    : llvm::RTTIExtends<DeserializeOptions, llvm::RTTIRoot> {
  static char ID;  
};
class Serializable : public llvm::RTTIExtends<Serializable, llvm::RTTIRoot> {
 public:
  static char ID;  
  using DeserializeOptions = ::xla::ifrt::DeserializeOptions;
};
class SerDes : public llvm::RTTIExtends<SerDes, llvm::RTTIRoot> {
 public:
  virtual absl::string_view type_name() const = 0;
  virtual absl::StatusOr<std::string> Serialize(Serializable& serializable) = 0;
  virtual absl::StatusOr<std::unique_ptr<Serializable>> Deserialize(
      const std::string& serialized,
      std::unique_ptr<DeserializeOptions> options) = 0;
  static char ID;  
};
void RegisterSerDes(const void* type_id, std::unique_ptr<SerDes> serdes);
template <typename T>
void RegisterSerDes(std::unique_ptr<SerDes> serdes) {
  static_assert(std::is_base_of_v<Serializable, T>,
                "Types must implement `xla::ifrt::Serializable` to have a "
                "serdes implementation");
  RegisterSerDes(T::classID(), std::move(serdes));
}
namespace serdes_internal {
absl::StatusOr<std::unique_ptr<Serializable>> DeserializeUnchecked(
    const Serialized& serialized, std::unique_ptr<DeserializeOptions> options);
}  
absl::StatusOr<Serialized> Serialize(Serializable& serializable);
template <typename InterfaceType>
absl::StatusOr<std::unique_ptr<InterfaceType>> Deserialize(
    const Serialized& serialized,
    std::unique_ptr<typename InterfaceType::DeserializeOptions> options) {
  TF_ASSIGN_OR_RETURN(auto result, serdes_internal::DeserializeUnchecked(
                                       serialized, std::move(options)));
  if (!llvm::isa<InterfaceType>(result.get())) {
    return absl::InternalError(
        "Unexpected Serializable type after deserialization");
  }
  return std::unique_ptr<InterfaceType>(
      static_cast<InterfaceType*>(result.release()));
}
}  
}  
#endif  
#include "xla/python/ifrt/serdes.h"
#include <memory>
#include <string>
#include <utility>
#include "absl/base/thread_annotations.h"
#include "absl/container/flat_hash_map.h"
#include "absl/log/check.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
#include "absl/synchronization/mutex.h"
#include "xla/python/ifrt/serdes.pb.h"
#include "tsl/platform/statusor.h"
namespace xla {
namespace ifrt {
namespace {
struct Registry {
  absl::Mutex mu;
  absl::flat_hash_map<const void*, SerDes*> type_id_to_serdes
      ABSL_GUARDED_BY(mu);
  absl::flat_hash_map<absl::string_view, SerDes*> name_to_serdes
      ABSL_GUARDED_BY(mu);
};
Registry* registry() {
  static auto* r = new Registry();
  return r;
}
}  
char Serializable::ID = 0;
char DeserializeOptions::ID = 0;
char SerDes::ID = 0;
void RegisterSerDes(const void* type_id, std::unique_ptr<SerDes> serdes) {
  Registry* const r = registry();
  absl::MutexLock l(&r->mu);
  CHECK(r->type_id_to_serdes.insert({type_id, serdes.get()}).second)
      << "xla::ifrt::SerDes cannot be registered more than once for the same "
         "type id: "
      << type_id;
  const absl::string_view name = serdes->type_name();
  CHECK(r->name_to_serdes.insert({name, serdes.get()}).second)
      << "xla::ifrt::SerDes cannot be registered more than once for the same "
         "name: "
      << name;
  serdes.release();
}
absl::StatusOr<Serialized> Serialize(Serializable& serializable) {
  SerDes* serdes;
  {
    Registry* const r = registry();
    absl::MutexLock l(&r->mu);
    auto it = r->type_id_to_serdes.find(serializable.dynamicClassID());
    if (it == r->type_id_to_serdes.end()) {
      return absl::UnimplementedError(
          "Serialize call failed. Serializable has no associated SerDes "
          "implementation");
    }
    serdes = it->second;
  }
  TF_ASSIGN_OR_RETURN(std::string data, serdes->Serialize(serializable));
  Serialized proto;
  proto.set_type_name(std::string(serdes->type_name()));
  proto.set_data(std::move(data));
  return proto;
}
namespace serdes_internal {
absl::StatusOr<std::unique_ptr<Serializable>> DeserializeUnchecked(
    const Serialized& serialized, std::unique_ptr<DeserializeOptions> options) {
  SerDes* serdes;
  {
    Registry* const r = registry();
    absl::MutexLock l(&r->mu);
    auto it = r->name_to_serdes.find(serialized.type_name());
    if (it == r->name_to_serdes.end()) {
      return absl::UnimplementedError(absl::StrCat(
          "Deserialize call failed. Serializable has no associated SerDes ",
          "implementation. type_name: ", serialized.type_name()));
    }
    serdes = it->second;
  }
  return serdes->Deserialize(serialized.data(), std::move(options));
}
}  
}  
}  