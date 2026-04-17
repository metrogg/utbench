#ifndef QUICHE_QUIC_CORE_QPACK_QPACK_DECODER_H_
#define QUICHE_QUIC_CORE_QPACK_QPACK_DECODER_H_
#include <cstdint>
#include <memory>
#include <set>
#include "absl/strings/string_view.h"
#include "quiche/quic/core/qpack/qpack_decoder_stream_sender.h"
#include "quiche/quic/core/qpack/qpack_encoder_stream_receiver.h"
#include "quiche/quic/core/qpack/qpack_header_table.h"
#include "quiche/quic/core/qpack/qpack_progressive_decoder.h"
#include "quiche/quic/core/quic_error_codes.h"
#include "quiche/quic/core/quic_types.h"
#include "quiche/quic/platform/api/quic_export.h"
namespace quic {
class QUICHE_EXPORT QpackDecoder
    : public QpackEncoderStreamReceiver::Delegate,
      public QpackProgressiveDecoder::BlockedStreamLimitEnforcer,
      public QpackProgressiveDecoder::DecodingCompletedVisitor {
 public:
  class QUICHE_EXPORT EncoderStreamErrorDelegate {
   public:
    virtual ~EncoderStreamErrorDelegate() {}
    virtual void OnEncoderStreamError(QuicErrorCode error_code,
                                      absl::string_view error_message) = 0;
  };
  QpackDecoder(uint64_t maximum_dynamic_table_capacity,
               uint64_t maximum_blocked_streams,
               EncoderStreamErrorDelegate* encoder_stream_error_delegate);
  ~QpackDecoder() override;
  void OnStreamReset(QuicStreamId stream_id);
  bool OnStreamBlocked(QuicStreamId stream_id) override;
  void OnStreamUnblocked(QuicStreamId stream_id) override;
  void OnDecodingCompleted(QuicStreamId stream_id,
                           uint64_t required_insert_count) override;
  std::unique_ptr<QpackProgressiveDecoder> CreateProgressiveDecoder(
      QuicStreamId stream_id,
      QpackProgressiveDecoder::HeadersHandlerInterface* handler);
  void OnInsertWithNameReference(bool is_static, uint64_t name_index,
                                 absl::string_view value) override;
  void OnInsertWithoutNameReference(absl::string_view name,
                                    absl::string_view value) override;
  void OnDuplicate(uint64_t index) override;
  void OnSetDynamicTableCapacity(uint64_t capacity) override;
  void OnErrorDetected(QuicErrorCode error_code,
                       absl::string_view error_message) override;
  void set_qpack_stream_sender_delegate(QpackStreamSenderDelegate* delegate) {
    decoder_stream_sender_.set_qpack_stream_sender_delegate(delegate);
  }
  QpackStreamReceiver* encoder_stream_receiver() {
    return &encoder_stream_receiver_;
  }
  bool dynamic_table_entry_referenced() const {
    return header_table_.dynamic_table_entry_referenced();
  }
  void FlushDecoderStream();
 private:
  EncoderStreamErrorDelegate* const encoder_stream_error_delegate_;
  QpackEncoderStreamReceiver encoder_stream_receiver_;
  QpackDecoderStreamSender decoder_stream_sender_;
  QpackDecoderHeaderTable header_table_;
  std::set<QuicStreamId> blocked_streams_;
  const uint64_t maximum_blocked_streams_;
  uint64_t known_received_count_;
};
class QUICHE_EXPORT NoopEncoderStreamErrorDelegate
    : public QpackDecoder::EncoderStreamErrorDelegate {
 public:
  ~NoopEncoderStreamErrorDelegate() override = default;
  void OnEncoderStreamError(QuicErrorCode ,
                            absl::string_view ) override {}
};
}  
#endif  
#include "quiche/quic/core/qpack/qpack_decoder.h"
#include <memory>
#include <utility>
#include "absl/strings/string_view.h"
#include "quiche/quic/core/qpack/qpack_index_conversions.h"
#include "quiche/quic/platform/api/quic_flag_utils.h"
#include "quiche/quic/platform/api/quic_flags.h"
#include "quiche/quic/platform/api/quic_logging.h"
namespace quic {
QpackDecoder::QpackDecoder(
    uint64_t maximum_dynamic_table_capacity, uint64_t maximum_blocked_streams,
    EncoderStreamErrorDelegate* encoder_stream_error_delegate)
    : encoder_stream_error_delegate_(encoder_stream_error_delegate),
      encoder_stream_receiver_(this),
      maximum_blocked_streams_(maximum_blocked_streams),
      known_received_count_(0) {
  QUICHE_DCHECK(encoder_stream_error_delegate_);
  header_table_.SetMaximumDynamicTableCapacity(maximum_dynamic_table_capacity);
}
QpackDecoder::~QpackDecoder() {}
void QpackDecoder::OnStreamReset(QuicStreamId stream_id) {
  if (header_table_.maximum_dynamic_table_capacity() > 0) {
    decoder_stream_sender_.SendStreamCancellation(stream_id);
    if (!GetQuicRestartFlag(quic_opport_bundle_qpack_decoder_data5)) {
      decoder_stream_sender_.Flush();
    }
  }
}
bool QpackDecoder::OnStreamBlocked(QuicStreamId stream_id) {
  auto result = blocked_streams_.insert(stream_id);
  QUICHE_DCHECK(result.second);
  return blocked_streams_.size() <= maximum_blocked_streams_;
}
void QpackDecoder::OnStreamUnblocked(QuicStreamId stream_id) {
  size_t result = blocked_streams_.erase(stream_id);
  QUICHE_DCHECK_EQ(1u, result);
}
void QpackDecoder::OnDecodingCompleted(QuicStreamId stream_id,
                                       uint64_t required_insert_count) {
  if (required_insert_count > 0) {
    decoder_stream_sender_.SendHeaderAcknowledgement(stream_id);
    if (known_received_count_ < required_insert_count) {
      known_received_count_ = required_insert_count;
    }
  }
  if (known_received_count_ < header_table_.inserted_entry_count()) {
    decoder_stream_sender_.SendInsertCountIncrement(
        header_table_.inserted_entry_count() - known_received_count_);
    known_received_count_ = header_table_.inserted_entry_count();
  }
  if (!GetQuicRestartFlag(quic_opport_bundle_qpack_decoder_data5)) {
    decoder_stream_sender_.Flush();
  }
}
void QpackDecoder::OnInsertWithNameReference(bool is_static,
                                             uint64_t name_index,
                                             absl::string_view value) {
  if (is_static) {
    auto entry = header_table_.LookupEntry( true, name_index);
    if (!entry) {
      OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_INVALID_STATIC_ENTRY,
                      "Invalid static table entry.");
      return;
    }
    if (!header_table_.EntryFitsDynamicTableCapacity(entry->name(), value)) {
      OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_ERROR_INSERTING_STATIC,
                      "Error inserting entry with name reference.");
      return;
    }
    header_table_.InsertEntry(entry->name(), value);
    return;
  }
  uint64_t absolute_index;
  if (!QpackEncoderStreamRelativeIndexToAbsoluteIndex(
          name_index, header_table_.inserted_entry_count(), &absolute_index)) {
    OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_INSERTION_INVALID_RELATIVE_INDEX,
                    "Invalid relative index.");
    return;
  }
  const QpackEntry* entry =
      header_table_.LookupEntry( false, absolute_index);
  if (!entry) {
    OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_INSERTION_DYNAMIC_ENTRY_NOT_FOUND,
                    "Dynamic table entry not found.");
    return;
  }
  if (!header_table_.EntryFitsDynamicTableCapacity(entry->name(), value)) {
    OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_ERROR_INSERTING_DYNAMIC,
                    "Error inserting entry with name reference.");
    return;
  }
  header_table_.InsertEntry(entry->name(), value);
}
void QpackDecoder::OnInsertWithoutNameReference(absl::string_view name,
                                                absl::string_view value) {
  if (!header_table_.EntryFitsDynamicTableCapacity(name, value)) {
    OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_ERROR_INSERTING_LITERAL,
                    "Error inserting literal entry.");
    return;
  }
  header_table_.InsertEntry(name, value);
}
void QpackDecoder::OnDuplicate(uint64_t index) {
  uint64_t absolute_index;
  if (!QpackEncoderStreamRelativeIndexToAbsoluteIndex(
          index, header_table_.inserted_entry_count(), &absolute_index)) {
    OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_DUPLICATE_INVALID_RELATIVE_INDEX,
                    "Invalid relative index.");
    return;
  }
  const QpackEntry* entry =
      header_table_.LookupEntry( false, absolute_index);
  if (!entry) {
    OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_DUPLICATE_DYNAMIC_ENTRY_NOT_FOUND,
                    "Dynamic table entry not found.");
    return;
  }
  if (!header_table_.EntryFitsDynamicTableCapacity(entry->name(),
                                                   entry->value())) {
    OnErrorDetected(QUIC_INTERNAL_ERROR, "Error inserting duplicate entry.");
    return;
  }
  header_table_.InsertEntry(entry->name(), entry->value());
}
void QpackDecoder::OnSetDynamicTableCapacity(uint64_t capacity) {
  if (!header_table_.SetDynamicTableCapacity(capacity)) {
    OnErrorDetected(QUIC_QPACK_ENCODER_STREAM_SET_DYNAMIC_TABLE_CAPACITY,
                    "Error updating dynamic table capacity.");
  }
}
void QpackDecoder::OnErrorDetected(QuicErrorCode error_code,
                                   absl::string_view error_message) {
  encoder_stream_error_delegate_->OnEncoderStreamError(error_code,
                                                       error_message);
}
std::unique_ptr<QpackProgressiveDecoder> QpackDecoder::CreateProgressiveDecoder(
    QuicStreamId stream_id,
    QpackProgressiveDecoder::HeadersHandlerInterface* handler) {
  return std::make_unique<QpackProgressiveDecoder>(stream_id, this, this,
                                                   &header_table_, handler);
}
void QpackDecoder::FlushDecoderStream() { decoder_stream_sender_.Flush(); }
}  