#ifndef QUICHE_HTTP2_HPACK_DECODER_HPACK_DECODER_STATE_H_
#define QUICHE_HTTP2_HPACK_DECODER_HPACK_DECODER_STATE_H_
#include <stddef.h>
#include <cstdint>
#include "absl/strings/string_view.h"
#include "quiche/http2/hpack/decoder/hpack_decoder_listener.h"
#include "quiche/http2/hpack/decoder/hpack_decoder_string_buffer.h"
#include "quiche/http2/hpack/decoder/hpack_decoder_tables.h"
#include "quiche/http2/hpack/decoder/hpack_decoding_error.h"
#include "quiche/http2/hpack/decoder/hpack_whole_entry_listener.h"
#include "quiche/http2/hpack/http2_hpack_constants.h"
#include "quiche/common/platform/api/quiche_export.h"
namespace http2 {
namespace test {
class HpackDecoderStatePeer;
}  
class QUICHE_EXPORT HpackDecoderState : public HpackWholeEntryListener {
 public:
  explicit HpackDecoderState(HpackDecoderListener* listener);
  ~HpackDecoderState() override;
  HpackDecoderState(const HpackDecoderState&) = delete;
  HpackDecoderState& operator=(const HpackDecoderState&) = delete;
  HpackDecoderListener* listener() const { return listener_; }
  void ApplyHeaderTableSizeSetting(uint32_t max_header_table_size);
  size_t GetCurrentHeaderTableSizeSetting() const {
    return final_header_table_size_;
  }
  void OnHeaderBlockStart();
  void OnIndexedHeader(size_t index) override;
  void OnNameIndexAndLiteralValue(
      HpackEntryType entry_type, size_t name_index,
      HpackDecoderStringBuffer* value_buffer) override;
  void OnLiteralNameAndValue(HpackEntryType entry_type,
                             HpackDecoderStringBuffer* name_buffer,
                             HpackDecoderStringBuffer* value_buffer) override;
  void OnDynamicTableSizeUpdate(size_t size) override;
  void OnHpackDecodeError(HpackDecodingError error) override;
  void OnHeaderBlockEnd();
  HpackDecodingError error() const { return error_; }
  size_t GetDynamicTableSize() const {
    return decoder_tables_.current_header_table_size();
  }
  const HpackDecoderTables& decoder_tables_for_test() const {
    return decoder_tables_;
  }
 private:
  friend class test::HpackDecoderStatePeer;
  void ReportError(HpackDecodingError error);
  HpackDecoderTables decoder_tables_;
  HpackDecoderListener* listener_;
  uint32_t final_header_table_size_;
  uint32_t lowest_header_table_size_;
  bool require_dynamic_table_size_update_;
  bool allow_dynamic_table_size_update_;
  bool saw_dynamic_table_size_update_;
  HpackDecodingError error_;
};
}  
#endif  
#include "quiche/http2/hpack/decoder/hpack_decoder_state.h"
#include <string>
#include <utility>
#include "quiche/http2/http2_constants.h"
#include "quiche/common/platform/api/quiche_logging.h"
namespace http2 {
namespace {
std::string ExtractString(HpackDecoderStringBuffer* string_buffer) {
  if (string_buffer->IsBuffered()) {
    return string_buffer->ReleaseString();
  } else {
    auto result = std::string(string_buffer->str());
    string_buffer->Reset();
    return result;
  }
}
}  
HpackDecoderState::HpackDecoderState(HpackDecoderListener* listener)
    : listener_(listener),
      final_header_table_size_(Http2SettingsInfo::DefaultHeaderTableSize()),
      lowest_header_table_size_(final_header_table_size_),
      require_dynamic_table_size_update_(false),
      allow_dynamic_table_size_update_(true),
      saw_dynamic_table_size_update_(false),
      error_(HpackDecodingError::kOk) {
  QUICHE_CHECK(listener_);
}
HpackDecoderState::~HpackDecoderState() = default;
void HpackDecoderState::ApplyHeaderTableSizeSetting(
    uint32_t header_table_size) {
  QUICHE_DVLOG(2) << "HpackDecoderState::ApplyHeaderTableSizeSetting("
                  << header_table_size << ")";
  QUICHE_DCHECK_LE(lowest_header_table_size_, final_header_table_size_);
  if (header_table_size < lowest_header_table_size_) {
    lowest_header_table_size_ = header_table_size;
  }
  final_header_table_size_ = header_table_size;
  QUICHE_DVLOG(2) << "low water mark: " << lowest_header_table_size_;
  QUICHE_DVLOG(2) << "final limit: " << final_header_table_size_;
}
void HpackDecoderState::OnHeaderBlockStart() {
  QUICHE_DVLOG(2) << "HpackDecoderState::OnHeaderBlockStart";
  QUICHE_DCHECK(error_ == HpackDecodingError::kOk)
      << HpackDecodingErrorToString(error_);
  QUICHE_DCHECK_LE(lowest_header_table_size_, final_header_table_size_);
  allow_dynamic_table_size_update_ = true;
  saw_dynamic_table_size_update_ = false;
  require_dynamic_table_size_update_ =
      (lowest_header_table_size_ <
           decoder_tables_.current_header_table_size() ||
       final_header_table_size_ < decoder_tables_.header_table_size_limit());
  QUICHE_DVLOG(2) << "HpackDecoderState::OnHeaderListStart "
                  << "require_dynamic_table_size_update_="
                  << require_dynamic_table_size_update_;
  listener_->OnHeaderListStart();
}
void HpackDecoderState::OnIndexedHeader(size_t index) {
  QUICHE_DVLOG(2) << "HpackDecoderState::OnIndexedHeader: " << index;
  if (error_ != HpackDecodingError::kOk) {
    return;
  }
  if (require_dynamic_table_size_update_) {
    ReportError(HpackDecodingError::kMissingDynamicTableSizeUpdate);
    return;
  }
  allow_dynamic_table_size_update_ = false;
  const HpackStringPair* entry = decoder_tables_.Lookup(index);
  if (entry != nullptr) {
    listener_->OnHeader(entry->name, entry->value);
  } else {
    ReportError(HpackDecodingError::kInvalidIndex);
  }
}
void HpackDecoderState::OnNameIndexAndLiteralValue(
    HpackEntryType entry_type, size_t name_index,
    HpackDecoderStringBuffer* value_buffer) {
  QUICHE_DVLOG(2) << "HpackDecoderState::OnNameIndexAndLiteralValue "
                  << entry_type << ", " << name_index << ", "
                  << value_buffer->str();
  if (error_ != HpackDecodingError::kOk) {
    return;
  }
  if (require_dynamic_table_size_update_) {
    ReportError(HpackDecodingError::kMissingDynamicTableSizeUpdate);
    return;
  }
  allow_dynamic_table_size_update_ = false;
  const HpackStringPair* entry = decoder_tables_.Lookup(name_index);
  if (entry != nullptr) {
    std::string value(ExtractString(value_buffer));
    listener_->OnHeader(entry->name, value);
    if (entry_type == HpackEntryType::kIndexedLiteralHeader) {
      decoder_tables_.Insert(entry->name, std::move(value));
    }
  } else {
    ReportError(HpackDecodingError::kInvalidNameIndex);
  }
}
void HpackDecoderState::OnLiteralNameAndValue(
    HpackEntryType entry_type, HpackDecoderStringBuffer* name_buffer,
    HpackDecoderStringBuffer* value_buffer) {
  QUICHE_DVLOG(2) << "HpackDecoderState::OnLiteralNameAndValue " << entry_type
                  << ", " << name_buffer->str() << ", " << value_buffer->str();
  if (error_ != HpackDecodingError::kOk) {
    return;
  }
  if (require_dynamic_table_size_update_) {
    ReportError(HpackDecodingError::kMissingDynamicTableSizeUpdate);
    return;
  }
  allow_dynamic_table_size_update_ = false;
  std::string name(ExtractString(name_buffer));
  std::string value(ExtractString(value_buffer));
  listener_->OnHeader(name, value);
  if (entry_type == HpackEntryType::kIndexedLiteralHeader) {
    decoder_tables_.Insert(std::move(name), std::move(value));
  }
}
void HpackDecoderState::OnDynamicTableSizeUpdate(size_t size_limit) {
  QUICHE_DVLOG(2) << "HpackDecoderState::OnDynamicTableSizeUpdate "
                  << size_limit << ", required="
                  << (require_dynamic_table_size_update_ ? "true" : "false")
                  << ", allowed="
                  << (allow_dynamic_table_size_update_ ? "true" : "false");
  if (error_ != HpackDecodingError::kOk) {
    return;
  }
  QUICHE_DCHECK_LE(lowest_header_table_size_, final_header_table_size_);
  if (!allow_dynamic_table_size_update_) {
    ReportError(HpackDecodingError::kDynamicTableSizeUpdateNotAllowed);
    return;
  }
  if (require_dynamic_table_size_update_) {
    if (size_limit > lowest_header_table_size_) {
      ReportError(HpackDecodingError::
                      kInitialDynamicTableSizeUpdateIsAboveLowWaterMark);
      return;
    }
    require_dynamic_table_size_update_ = false;
  } else if (size_limit > final_header_table_size_) {
    ReportError(
        HpackDecodingError::kDynamicTableSizeUpdateIsAboveAcknowledgedSetting);
    return;
  }
  decoder_tables_.DynamicTableSizeUpdate(size_limit);
  if (saw_dynamic_table_size_update_) {
    allow_dynamic_table_size_update_ = false;
  } else {
    saw_dynamic_table_size_update_ = true;
  }
  lowest_header_table_size_ = final_header_table_size_;
}
void HpackDecoderState::OnHpackDecodeError(HpackDecodingError error) {
  QUICHE_DVLOG(2) << "HpackDecoderState::OnHpackDecodeError "
                  << HpackDecodingErrorToString(error);
  if (error_ == HpackDecodingError::kOk) {
    ReportError(error);
  }
}
void HpackDecoderState::OnHeaderBlockEnd() {
  QUICHE_DVLOG(2) << "HpackDecoderState::OnHeaderBlockEnd";
  if (error_ != HpackDecodingError::kOk) {
    return;
  }
  if (require_dynamic_table_size_update_) {
    ReportError(HpackDecodingError::kMissingDynamicTableSizeUpdate);
  } else {
    listener_->OnHeaderListEnd();
  }
}
void HpackDecoderState::ReportError(HpackDecodingError error) {
  QUICHE_DVLOG(2) << "HpackDecoderState::ReportError is new="
                  << (error_ == HpackDecodingError::kOk ? "true" : "false")
                  << ", error: " << HpackDecodingErrorToString(error);
  if (error_ == HpackDecodingError::kOk) {
    listener_->OnHeaderErrorDetected(HpackDecodingErrorToString(error));
    error_ = error;
  }
}
}  