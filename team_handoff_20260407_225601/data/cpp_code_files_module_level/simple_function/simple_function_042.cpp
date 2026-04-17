#ifndef TENSORFLOW_CORE_KERNELS_TENSOR_MAP_H_
#define TENSORFLOW_CORE_KERNELS_TENSOR_MAP_H_
#include <utility>
#include "absl/container/flat_hash_map.h"
#include "tensorflow/core/framework/tensor.h"
#include "tensorflow/core/framework/tensor_key.h"
#include "tensorflow/core/framework/variant.h"
#include "tensorflow/core/framework/variant_tensor_data.h"
#include "tensorflow/core/lib/core/refcount.h"
namespace tensorflow {
class TensorMap {
 public:
  TensorMap() : tensors_(new Tensors) {}
  ~TensorMap();
  TensorMap(const TensorMap& other) : tensors_(other.tensors_) {
    tensors_->Ref();
  }
  TensorMap(TensorMap&& rhs) : tensors_(rhs.tensors_) {
    rhs.tensors_ = nullptr;
  }
  TensorMap& operator=(const TensorMap& rhs) {
    if (this == &rhs) return *this;
    tensors_->Unref();
    tensors_ = rhs.tensors_;
    tensors_->Ref();
    return *this;
  }
  TensorMap& operator=(TensorMap&& rhs) {
    if (this == &rhs) return *this;
    std::swap(tensors_, rhs.tensors_);
    return *this;
  }
  static const char kTypeName[];
  string TypeName() const { return kTypeName; }
  void Encode(VariantTensorData* data) const;
  bool Decode(const VariantTensorData& data);
  string DebugString() const { return "TensorMap"; }
  absl::flat_hash_map<TensorKey, Tensor>& tensors() {
    return tensors_->values_;
  }
  const absl::flat_hash_map<TensorKey, Tensor>& tensors() const {
    return tensors_->values_;
  }
  TensorMap Copy() const {
    TensorMap out;
    out.tensors_->values_ = tensors_->values_;
    return out;
  }
  bool insert(const TensorKey& key, const Tensor& value) {
    auto r = tensors_->values_.try_emplace(key, value);
    return r.second;
  }
  absl::flat_hash_map<TensorKey, Tensor>::iterator find(TensorKey key) {
    return tensors_->values_.find(key);
  }
  Tensor& lookup(TensorKey key) { return tensors_->values_.find(key)->second; }
  Tensor& operator[](TensorKey& k) { return tensors_->values_[k]; }
  bool replace(const TensorKey& k, const Tensor& v) {
    tensors_->values_[k] = v;
    return true;
  }
  size_t erase(TensorKey key) { return tensors_->values_.erase(key); }
  size_t size() const { return tensors_->values_.size(); }
  std::vector<Tensor> keys() const {
    std::vector<Tensor> keys;
    keys.reserve(tensors_->values_.size());
    absl::flat_hash_map<TensorKey, Tensor>::iterator it =
        tensors_->values_.begin();
    while (it != tensors_->values_.end()) {
      keys.push_back(it->first);
      it++;
    }
    return keys;
  }
  bool RefCountIsOne() const { return tensors_->RefCountIsOne(); }
 private:
  class Tensors : public core::RefCounted {
   public:
    absl::flat_hash_map<TensorKey, Tensor> values_;
  };
  Tensors* tensors_;
};
#if defined(PLATFORM_GOOGLE)
static_assert(Variant::CanInlineType<TensorMap>() || sizeof(void*) < 8,
              "Must be able to inline TensorMap into a Variant");
#endif
}  
#endif  
#include "tensorflow/core/kernels/tensor_map.h"
#include "tensorflow/core/framework/tensor_shape.h"
#include "tensorflow/core/framework/tensor_shape.pb.h"
#include "tensorflow/core/framework/variant_op_registry.h"
#include "tensorflow/core/lib/core/coding.h"
namespace tensorflow {
TensorMap::~TensorMap() {
  if (tensors_) tensors_->Unref();
}
void TensorMap::Encode(VariantTensorData* data) const {
  data->set_type_name(TypeName());
  absl::flat_hash_map<TensorKey, Tensor>::const_iterator map_it =
      tensors().begin();
  while (map_it != tensors().end()) {
    Tensor k = map_it->first;
    Tensor v = map_it->second;
    CHECK_NE(k.dtype(), DT_INVALID);
    CHECK_NE(v.dtype(), DT_INVALID);
    *data->add_tensors() = k;
    *data->add_tensors() = v;
    map_it++;
  }
}
static Status TensorMapDeviceCopy(
    const TensorMap& from, TensorMap* to,
    const UnaryVariantOpRegistry::AsyncTensorDeviceCopyFn& copy) {
  for (const std::pair<TensorKey, Tensor>& p : from.tensors()) {
    TensorKey to_key(p.first.dtype());
    Tensor to_val(p.second.dtype());
    TF_RETURN_IF_ERROR(copy(p.first, &to_key));
    TF_RETURN_IF_ERROR(copy(p.second, &to_val));
    to->tensors().emplace(to_key, to_val);
  }
  return absl::OkStatus();
}
#define REGISTER_LIST_COPY(DIRECTION)                                        \
  INTERNAL_REGISTER_UNARY_VARIANT_DEVICE_COPY_FUNCTION(TensorMap, DIRECTION, \
                                                       TensorMapDeviceCopy)
REGISTER_LIST_COPY(VariantDeviceCopyDirection::HOST_TO_DEVICE);
REGISTER_LIST_COPY(VariantDeviceCopyDirection::DEVICE_TO_HOST);
REGISTER_LIST_COPY(VariantDeviceCopyDirection::DEVICE_TO_DEVICE);
REGISTER_UNARY_VARIANT_DECODE_FUNCTION(TensorMap, TensorMap::kTypeName);
bool TensorMap::Decode(const VariantTensorData& data) {
  std::vector<Tensor>::const_iterator tensors_it = data.tensors().begin();
  while (tensors_it != data.tensors().end()) {
    if (std::next(tensors_it) == data.tensors().end()) {
      return false;
    }
    tensors().emplace(tensors_it[0], tensors_it[1]);
    tensors_it += 2;
  }
  return true;
}
const char TensorMap::kTypeName[] = "tensorflow::TensorMap";
}  