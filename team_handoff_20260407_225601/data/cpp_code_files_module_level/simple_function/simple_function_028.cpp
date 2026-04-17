#ifndef TENSORFLOW_LITE_CORE_ASYNC_INTEROP_RECONCILE_FNS_H_
#define TENSORFLOW_LITE_CORE_ASYNC_INTEROP_RECONCILE_FNS_H_
#include "tensorflow/lite/core/async/interop/attribute_map_internal.h"
#include "tensorflow/lite/core/async/interop/c/types.h"
namespace tflite {
namespace interop {
bool ReconcileGeneralAttributeKeys(TfLiteAttrMapType type,
                                   const AttributeMap::ContainerT* lhs,
                                   const AttributeMap::ContainerT* rhs,
                                   AttributeMap::ContainerT* merged,
                                   AttributeMap::ContainerT* conflict);
bool CheckGeneralAttributeKeysCoverage(TfLiteAttrMapType type,
                                       const AttributeMap::ContainerT* lhs,
                                       const AttributeMap::ContainerT* rhs,
                                       AttributeMap::ContainerT* conflict);
}  
}  
#endif  
#include "tensorflow/lite/core/async/interop/reconcile_fns.h"
#include <algorithm>
#include <cstddef>
#include <cstdint>
#include <iterator>
#include <set>
#include "tensorflow/lite/core/async/interop/attribute_map_internal.h"
#include "tensorflow/lite/core/async/interop/c/types.h"
namespace tflite {
namespace interop {
namespace {
template <typename T>
T gcd(T x, T y) {
  while (y) {
    auto m = x % y;
    x = y;
    y = m;
  }
  return x;
}
template <typename T>
T lcm(T x, T y) {
  return x / gcd(x, y) * y;
}
void ReconcileAlignment(size_t l, size_t r, AttributeMap::ContainerT* merged) {
  merged->insert_or_assign(static_cast<size_t>(kTfLiteBufferAttrKeyAlignment),
                           lcm(l, r));
}
void ReconcilePadding(size_t l, size_t r, AttributeMap::ContainerT* merged) {
  merged->insert_or_assign(static_cast<size_t>(kTfLiteBufferAttrKeyPadding),
                           lcm(l, r));
}
bool CheckMultiples(size_t l, size_t r) { return l % r == 0; }
void ReconcileSize(size_t l, size_t r, AttributeMap::ContainerT* merged) {
  merged->insert_or_assign(static_cast<size_t>(kTfLiteBufferAttrKeySize),
                           std::max(l, r));
}
bool CheckSize(size_t l, size_t r) { return l >= r; }
}  
bool ReconcileGeneralAttributeKeys(TfLiteAttrMapType type,
                                   const AttributeMap::ContainerT* lhs,
                                   const AttributeMap::ContainerT* rhs,
                                   AttributeMap::ContainerT* merged,
                                   AttributeMap::ContainerT* conflict) {
  if (lhs == nullptr || rhs == nullptr || merged == nullptr) return false;
  bool ret = true;
  std::set<uint32_t> keys;
  std::transform(lhs->begin(), lhs->end(), std::inserter(keys, keys.end()),
                 [](auto pair) { return pair.first; });
  std::transform(rhs->begin(), rhs->end(), std::inserter(keys, keys.end()),
                 [](auto pair) { return pair.first; });
  for (auto k : keys) {
    const auto l = lhs->find(k);
    const auto r = rhs->find(k);
    if (l == lhs->end() || l->second.GetPtr() == nullptr) {
      merged->insert_or_assign(k, r->second);
      continue;
    }
    if (r == rhs->end() || r->second.GetPtr() == nullptr) {
      merged->insert_or_assign(k, l->second);
      continue;
    }
    if (type == kTfLiteAttrMapTypeBuffer) {
      switch (static_cast<TfLiteBufferAttrKey>(k)) {
        case kTfLiteBufferAttrKeySize:
          ReconcileSize(*l->second.Get<size_t>(), *r->second.Get<size_t>(),
                        merged);
          break;
        case kTfLiteBufferAttrKeyAlignment:
          ReconcileAlignment(*l->second.Get<size_t>(), *r->second.Get<size_t>(),
                             merged);
          break;
        case kTfLiteBufferAttrKeyPadding:
          ReconcilePadding(*l->second.Get<size_t>(), *r->second.Get<size_t>(),
                           merged);
          break;
        default:
          if (l->second == r->second) {
            merged->insert_or_assign(k, l->second);
          } else {
            ret = false;
            if (conflict) conflict->insert_or_assign(k, r->second);
          }
      }
    } else {
      if (l->second == r->second) {
        merged->insert_or_assign(k, l->second);
      } else {
        ret = false;
        if (conflict) conflict->insert_or_assign(k, r->second);
      }
    }
  }
  return ret;
}
bool CheckGeneralAttributeKeysCoverage(TfLiteAttrMapType type,
                                       const AttributeMap::ContainerT* lhs,
                                       const AttributeMap::ContainerT* rhs,
                                       AttributeMap::ContainerT* conflict) {
  if (lhs == nullptr || rhs == nullptr) return false;
  bool ret = true;
  std::set<uint32_t> keys;
  std::transform(lhs->begin(), lhs->end(), std::inserter(keys, keys.end()),
                 [](auto pair) { return pair.first; });
  std::transform(rhs->begin(), rhs->end(), std::inserter(keys, keys.end()),
                 [](auto pair) { return pair.first; });
  for (auto k : keys) {
    bool has_conflict = false;
    const auto l = lhs->find(k);
    const auto r = rhs->find(k);
    if (r == rhs->end() || r->second.GetPtr() == nullptr) {
      continue;
    } else if (l == lhs->end() || l->second.GetPtr() == nullptr) {
      has_conflict = true;
    } else {
      if (type == kTfLiteAttrMapTypeBuffer) {
        switch (static_cast<TfLiteBufferAttrKey>(k)) {
          case kTfLiteBufferAttrKeySize:
            has_conflict |=
                !CheckSize(*l->second.Get<size_t>(), *r->second.Get<size_t>());
            break;
          case kTfLiteBufferAttrKeyAlignment:
            has_conflict |= !CheckMultiples(*l->second.Get<size_t>(),
                                            *r->second.Get<size_t>());
            break;
          case kTfLiteBufferAttrKeyPadding:
            has_conflict |=
                !CheckSize(*l->second.Get<size_t>(), *r->second.Get<size_t>());
            break;
          default:
            if (l->second != r->second) {
              has_conflict = true;
            }
        }
      } else {
        if (l->second != r->second) {
          has_conflict = true;
        }
      }
    }
    if (has_conflict) {
      if (conflict != nullptr) conflict->insert_or_assign(k, r->second);
      ret = false;
    }
  }
  return ret;
}
}  
}  