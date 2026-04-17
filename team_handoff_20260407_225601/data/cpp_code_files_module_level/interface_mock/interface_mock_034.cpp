#ifndef I18N_ADDRESSINPUT_ONDEMAND_SUPPLY_TASK_H_
#define I18N_ADDRESSINPUT_ONDEMAND_SUPPLY_TASK_H_
#include <libaddressinput/supplier.h>
#include <map>
#include <memory>
#include <set>
#include <string>
#include "retriever.h"
namespace i18n {
namespace addressinput {
class LookupKey;
class Rule;
class OndemandSupplyTask {
 public:
  OndemandSupplyTask(const OndemandSupplyTask&) = delete;
  OndemandSupplyTask& operator=(const OndemandSupplyTask&) = delete;
  OndemandSupplyTask(const LookupKey& lookup_key,
                     std::map<std::string, const Rule*>* rules,
                     const Supplier::Callback& supplied);
  ~OndemandSupplyTask();
  void Queue(const std::string& key);
  void Retrieve(const Retriever& retriever);
  Supplier::RuleHierarchy hierarchy_;
 private:
  void Load(bool success, const std::string& key, const std::string& data);
  void Loaded();
  std::set<std::string> pending_;
  const LookupKey& lookup_key_;
  std::map<std::string, const Rule*>* const rule_cache_;
  const Supplier::Callback& supplied_;
  const std::unique_ptr<const Retriever::Callback> retrieved_;
  bool success_;
};
}  
}  
#endif  
#include "ondemand_supply_task.h"
#include <libaddressinput/address_field.h>
#include <libaddressinput/callback.h>
#include <libaddressinput/supplier.h>
#include <algorithm>
#include <cassert>
#include <cstddef>
#include <map>
#include <string>
#include "lookup_key.h"
#include "retriever.h"
#include "rule.h"
#include "util/size.h"
namespace i18n {
namespace addressinput {
OndemandSupplyTask::OndemandSupplyTask(
    const LookupKey& lookup_key,
    std::map<std::string, const Rule*>* rules,
    const Supplier::Callback& supplied)
    : hierarchy_(),
      pending_(),
      lookup_key_(lookup_key),
      rule_cache_(rules),
      supplied_(supplied),
      retrieved_(BuildCallback(this, &OndemandSupplyTask::Load)),
      success_(true) {
  assert(rule_cache_ != nullptr);
  assert(retrieved_ != nullptr);
}
OndemandSupplyTask::~OndemandSupplyTask() = default;
void OndemandSupplyTask::Queue(const std::string& key) {
  assert(pending_.find(key) == pending_.end());
  pending_.insert(key);
}
void OndemandSupplyTask::Retrieve(const Retriever& retriever) {
  if (pending_.empty()) {
    Loaded();
  } else {
    bool done = false;
    for (auto it = pending_.begin(); !done;) {
      const std::string& key = *it++;
      done = it == pending_.end();
      retriever.Retrieve(key, *retrieved_);
    }
  }
}
void OndemandSupplyTask::Load(bool success,
                              const std::string& key,
                              const std::string& data) {
  size_t depth = std::count(key.begin(), key.end(), '/') - 1;
  assert(depth < size(LookupKey::kHierarchy));
  size_t status = pending_.erase(key);
  assert(status == 1);  
  (void)status;  
  if (success) {
    if (data != "{}") {
      auto* rule = new Rule;
      if (LookupKey::kHierarchy[depth] == COUNTRY) {
        rule->CopyFrom(Rule::GetDefault());
      }
      if (rule->ParseSerializedRule(data)) {
        auto result = rule_cache_->emplace(rule->GetId(), rule);
        if (!result.second) {  
          delete rule;
        }
        hierarchy_.rule[depth] = result.first->second;
      } else {
        delete rule;
        success_ = false;
      }
    }
  } else {
    success_ = false;
  }
  if (pending_.empty()) {
    Loaded();
  }
}
void OndemandSupplyTask::Loaded() {
  supplied_(success_, lookup_key_, hierarchy_);
  delete this;
}
}  
}  