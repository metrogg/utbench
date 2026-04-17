#ifndef QUICHE_HTTP2_HPACK_DECODER_HPACK_DECODER_H_
#define QUICHE_HTTP2_HPACK_DECODER_HPACK_DECODER_H_
#include <stddef.h>
#include <cstdint>
#include "quiche/http2/decoder/decode_buffer.h"
#include "quiche/http2/hpack/decoder/hpack_block_decoder.h"
#include "quiche/http2/hpack/decoder/hpack_decoder_listener.h"
#include "quiche/http2/hpack/decoder/hpack_decoder_state.h"
#include "quiche/http2/hpack/decoder/hpack_decoder_tables.h"
#include "quiche/http2/hpack/decoder/hpack_decoding_error.h"
#include "quiche/http2/hpack/decoder/hpack_whole_entry_buffer.h"
#include "quiche/common/platform/api/quiche_export.h"
namespace http2 {
namespace test {
class HpackDecoderPeer;
}  
class QUICHE_EXPORT HpackDecoder {
 public:
  HpackDecoder(HpackDecoderListener* listener, size_t max_string_size);
  virtual ~HpackDecoder();
  HpackDecoder(const HpackDecoder&) = delete;
  HpackDecoder& operator=(const HpackDecoder&) = delete;
  void set_max_string_size_bytes(size_t max_string_size_bytes);
  void ApplyHeaderTableSizeSetting(uint32_t max_header_table_size);
  size_t GetCurrentHeaderTableSizeSetting() const {
    return decoder_state_.GetCurrentHeaderTableSizeSetting();
  }
  bool StartDecodingBlock();
  bool DecodeFragment(DecodeBuffer* db);
  bool EndDecodingBlock();
  bool DetectError();
  size_t GetDynamicTableSize() const {
    return decoder_state_.GetDynamicTableSize();
  }
  HpackDecodingError error() const { return error_; }
 private:
  friend class test::HpackDecoderPeer;
  void ReportError(HpackDecodingError error);
  HpackDecoderState decoder_state_;
  HpackWholeEntryBuffer entry_buffer_;
  HpackBlockDecoder block_decoder_;
  HpackDecodingError error_;
};
}  
#endif  
#include "quiche/http2/hpack/decoder/hpack_decoder.h"
#include "quiche/http2/decoder/decode_status.h"
#include "quiche/common/platform/api/quiche_flag_utils.h"
#include "quiche/common/platform/api/quiche_logging.h"
namespace http2 {
HpackDecoder::HpackDecoder(HpackDecoderListener* listener,
                           size_t max_string_size)
    : decoder_state_(listener),
      entry_buffer_(&decoder_state_, max_string_size),
      block_decoder_(&entry_buffer_),
      error_(HpackDecodingError::kOk) {}
HpackDecoder::~HpackDecoder() = default;
void HpackDecoder::set_max_string_size_bytes(size_t max_string_size_bytes) {
  entry_buffer_.set_max_string_size_bytes(max_string_size_bytes);
}
void HpackDecoder::ApplyHeaderTableSizeSetting(uint32_t max_header_table_size) {
  decoder_state_.ApplyHeaderTableSizeSetting(max_header_table_size);
}
bool HpackDecoder::StartDecodingBlock() {
  QUICHE_DVLOG(3) << "HpackDecoder::StartDecodingBlock, error_detected="
                  << (DetectError() ? "true" : "false");
  if (DetectError()) {
    return false;
  }
  block_decoder_.Reset();
  decoder_state_.OnHeaderBlockStart();
  return true;
}
bool HpackDecoder::DecodeFragment(DecodeBuffer* db) {
  QUICHE_DVLOG(3) << "HpackDecoder::DecodeFragment, error_detected="
                  << (DetectError() ? "true" : "false")
                  << ", size=" << db->Remaining();
  if (DetectError()) {
    QUICHE_CODE_COUNT_N(decompress_failure_3, 3, 23);
    return false;
  }
  DecodeStatus status = block_decoder_.Decode(db);
  if (status == DecodeStatus::kDecodeError) {
    ReportError(block_decoder_.error());
    QUICHE_CODE_COUNT_N(decompress_failure_3, 4, 23);
    return false;
  } else if (DetectError()) {
    QUICHE_CODE_COUNT_N(decompress_failure_3, 5, 23);
    return false;
  }
  QUICHE_DCHECK_EQ(block_decoder_.before_entry(),
                   status == DecodeStatus::kDecodeDone)
      << status;
  if (!block_decoder_.before_entry()) {
    entry_buffer_.BufferStringsIfUnbuffered();
  }
  return true;
}
bool HpackDecoder::EndDecodingBlock() {
  QUICHE_DVLOG(3) << "HpackDecoder::EndDecodingBlock, error_detected="
                  << (DetectError() ? "true" : "false");
  if (DetectError()) {
    QUICHE_CODE_COUNT_N(decompress_failure_3, 6, 23);
    return false;
  }
  if (!block_decoder_.before_entry()) {
    ReportError(HpackDecodingError::kTruncatedBlock);
    QUICHE_CODE_COUNT_N(decompress_failure_3, 7, 23);
    return false;
  }
  decoder_state_.OnHeaderBlockEnd();
  if (DetectError()) {
    QUICHE_CODE_COUNT_N(decompress_failure_3, 8, 23);
    return false;
  }
  return true;
}
bool HpackDecoder::DetectError() {
  if (error_ != HpackDecodingError::kOk) {
    return true;
  }
  if (decoder_state_.error() != HpackDecodingError::kOk) {
    QUICHE_DVLOG(2) << "Error detected in decoder_state_";
    QUICHE_CODE_COUNT_N(decompress_failure_3, 10, 23);
    error_ = decoder_state_.error();
  }
  return error_ != HpackDecodingError::kOk;
}
void HpackDecoder::ReportError(HpackDecodingError error) {
  QUICHE_DVLOG(3) << "HpackDecoder::ReportError is new="
                  << (error_ == HpackDecodingError::kOk ? "true" : "false")
                  << ", error: " << HpackDecodingErrorToString(error);
  if (error_ == HpackDecodingError::kOk) {
    error_ = error;
    decoder_state_.listener()->OnHeaderErrorDetected(
        HpackDecodingErrorToString(error));
  }
}
}  