#ifndef TENSORFLOW_LITE_TOOLS_TOOL_PARAMS_H_
#define TENSORFLOW_LITE_TOOLS_TOOL_PARAMS_H_
#include <memory>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>
namespace tflite {
namespace tools {
template <typename T>
class TypedToolParam;
class ToolParam {
 protected:
  enum class ParamType { TYPE_INT32, TYPE_FLOAT, TYPE_BOOL, TYPE_STRING };
  template <typename T>
  static ParamType GetValueType();
 public:
  template <typename T>
  static std::unique_ptr<ToolParam> Create(const T& default_value,
                                           int position = 0) {
    auto* param = new TypedToolParam<T>(default_value);
    param->SetPosition(position);
    return std::unique_ptr<ToolParam>(param);
  }
  template <typename T>
  TypedToolParam<T>* AsTyped() {
    AssertHasSameType(GetValueType<T>(), type_);
    return static_cast<TypedToolParam<T>*>(this);
  }
  template <typename T>
  const TypedToolParam<T>* AsConstTyped() const {
    AssertHasSameType(GetValueType<T>(), type_);
    return static_cast<const TypedToolParam<T>*>(this);
  }
  virtual ~ToolParam() {}
  explicit ToolParam(ParamType type)
      : has_value_set_(false), position_(0), type_(type) {}
  bool HasValueSet() const { return has_value_set_; }
  int GetPosition() const { return position_; }
  void SetPosition(int position) { position_ = position; }
  virtual void Set(const ToolParam&) {}
  virtual std::unique_ptr<ToolParam> Clone() const = 0;
 protected:
  bool has_value_set_;
  int position_;
 private:
  static void AssertHasSameType(ParamType a, ParamType b);
  const ParamType type_;
};
template <typename T>
class TypedToolParam : public ToolParam {
 public:
  explicit TypedToolParam(const T& value)
      : ToolParam(GetValueType<T>()), value_(value) {}
  void Set(const T& value) {
    value_ = value;
    has_value_set_ = true;
  }
  const T& Get() const { return value_; }
  void Set(const ToolParam& other) override {
    Set(other.AsConstTyped<T>()->Get());
    SetPosition(other.AsConstTyped<T>()->GetPosition());
  }
  std::unique_ptr<ToolParam> Clone() const override {
    return ToolParam::Create<T>(value_, position_);
  }
 private:
  T value_;
};
class ToolParams {
 public:
  void AddParam(const std::string& name, std::unique_ptr<ToolParam> value) {
    params_[name] = std::move(value);
  }
  void RemoveParam(const std::string& name) { params_.erase(name); }
  bool HasParam(const std::string& name) const {
    return params_.find(name) != params_.end();
  }
  bool Empty() const { return params_.empty(); }
  const ToolParam* GetParam(const std::string& name) const {
    const auto& entry = params_.find(name);
    if (entry == params_.end()) return nullptr;
    return entry->second.get();
  }
  template <typename T>
  void Set(const std::string& name, const T& value, int position = 0) {
    AssertParamExists(name);
    params_.at(name)->AsTyped<T>()->Set(value);
    params_.at(name)->AsTyped<T>()->SetPosition(position);
  }
  template <typename T>
  bool HasValueSet(const std::string& name) const {
    AssertParamExists(name);
    return params_.at(name)->AsConstTyped<T>()->HasValueSet();
  }
  template <typename T>
  int GetPosition(const std::string& name) const {
    AssertParamExists(name);
    return params_.at(name)->AsConstTyped<T>()->GetPosition();
  }
  template <typename T>
  T Get(const std::string& name) const {
    AssertParamExists(name);
    return params_.at(name)->AsConstTyped<T>()->Get();
  }
  void Set(const ToolParams& other);
  void Merge(const ToolParams& other, bool overwrite = false);
 private:
  void AssertParamExists(const std::string& name) const;
  std::unordered_map<std::string, std::unique_ptr<ToolParam>> params_;
};
#define LOG_TOOL_PARAM(params, type, name, description, verbose)      \
  do {                                                                \
    TFLITE_MAY_LOG(INFO, (verbose) || params.HasValueSet<type>(name)) \
        << description << ": [" << params.Get<type>(name) << "]";     \
  } while (0)
}  
}  
#endif  
#include "tensorflow/lite/tools/tool_params.h"
#include <string>
#include <unordered_map>
#include <vector>
#include "tensorflow/lite/tools/logging.h"
namespace tflite {
namespace tools {
void ToolParam::AssertHasSameType(ToolParam::ParamType a,
                                  ToolParam::ParamType b) {
  TFLITE_TOOLS_CHECK(a == b) << "Type mismatch while accessing parameter.";
}
template <>
ToolParam::ParamType ToolParam::GetValueType<int32_t>() {
  return ToolParam::ParamType::TYPE_INT32;
}
template <>
ToolParam::ParamType ToolParam::GetValueType<bool>() {
  return ToolParam::ParamType::TYPE_BOOL;
}
template <>
ToolParam::ParamType ToolParam::GetValueType<float>() {
  return ToolParam::ParamType::TYPE_FLOAT;
}
template <>
ToolParam::ParamType ToolParam::GetValueType<std::string>() {
  return ToolParam::ParamType::TYPE_STRING;
}
void ToolParams::AssertParamExists(const std::string& name) const {
  TFLITE_TOOLS_CHECK(HasParam(name)) << name << " was not found.";
}
void ToolParams::Set(const ToolParams& other) {
  for (const auto& param : params_) {
    const ToolParam* other_param = other.GetParam(param.first);
    if (other_param == nullptr) continue;
    param.second->Set(*other_param);
  }
}
void ToolParams::Merge(const ToolParams& other, bool overwrite) {
  for (const auto& one : other.params_) {
    auto it = params_.find(one.first);
    if (it == params_.end()) {
      AddParam(one.first, one.second->Clone());
    } else if (overwrite) {
      it->second->Set(*one.second);
    }
  }
}
}  
}  