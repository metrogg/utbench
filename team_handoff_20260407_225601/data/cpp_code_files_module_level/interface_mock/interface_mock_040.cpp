#ifndef QUICHE_QUIC_CORE_HTTP_QUIC_SPDY_CLIENT_SESSION_H_
#define QUICHE_QUIC_CORE_HTTP_QUIC_SPDY_CLIENT_SESSION_H_
#include <memory>
#include <string>
#include "quiche/quic/core/http/quic_spdy_client_session_base.h"
#include "quiche/quic/core/http/quic_spdy_client_stream.h"
#include "quiche/quic/core/quic_crypto_client_stream.h"
#include "quiche/quic/core/quic_packets.h"
namespace quic {
class QuicConnection;
class QuicServerId;
class QUICHE_EXPORT QuicSpdyClientSession : public QuicSpdyClientSessionBase {
 public:
  QuicSpdyClientSession(const QuicConfig& config,
                        const ParsedQuicVersionVector& supported_versions,
                        QuicConnection* connection,
                        const QuicServerId& server_id,
                        QuicCryptoClientConfig* crypto_config);
  QuicSpdyClientSession(const QuicConfig& config,
                        const ParsedQuicVersionVector& supported_versions,
                        QuicConnection* connection,
                        QuicSession::Visitor* visitor,
                        const QuicServerId& server_id,
                        QuicCryptoClientConfig* crypto_config);
  QuicSpdyClientSession(const QuicSpdyClientSession&) = delete;
  QuicSpdyClientSession& operator=(const QuicSpdyClientSession&) = delete;
  ~QuicSpdyClientSession() override;
  void Initialize() override;
  QuicSpdyClientStream* CreateOutgoingBidirectionalStream() override;
  QuicSpdyClientStream* CreateOutgoingUnidirectionalStream() override;
  QuicCryptoClientStreamBase* GetMutableCryptoStream() override;
  const QuicCryptoClientStreamBase* GetCryptoStream() const override;
  void OnProofValid(const QuicCryptoClientConfig::CachedState& cached) override;
  void OnProofVerifyDetailsAvailable(
      const ProofVerifyDetails& verify_details) override;
  virtual void CryptoConnect();
  int GetNumSentClientHellos() const;
  bool ResumptionAttempted() const;
  bool IsResumption() const;
  bool EarlyDataAccepted() const;
  bool ReceivedInchoateReject() const;
  int GetNumReceivedServerConfigUpdates() const;
  using QuicSession::CanOpenNextOutgoingBidirectionalStream;
  void set_respect_goaway(bool respect_goaway) {
    respect_goaway_ = respect_goaway;
  }
  QuicSSLConfig GetSSLConfig() const override {
    return crypto_config_->ssl_config();
  }
 protected:
  QuicSpdyStream* CreateIncomingStream(QuicStreamId id) override;
  QuicSpdyStream* CreateIncomingStream(PendingStream* pending) override;
  bool ShouldCreateOutgoingBidirectionalStream() override;
  bool ShouldCreateOutgoingUnidirectionalStream() override;
  bool ShouldCreateIncomingStream(QuicStreamId id) override;
  virtual std::unique_ptr<QuicCryptoClientStreamBase> CreateQuicCryptoStream();
  virtual std::unique_ptr<QuicSpdyClientStream> CreateClientStream();
  const QuicServerId& server_id() const { return server_id_; }
  QuicCryptoClientConfig* crypto_config() { return crypto_config_; }
 private:
  std::unique_ptr<QuicCryptoClientStreamBase> crypto_stream_;
  QuicServerId server_id_;
  QuicCryptoClientConfig* crypto_config_;
  bool respect_goaway_;
};
}  
#endif  
#include "quiche/quic/core/http/quic_spdy_client_session.h"
#include <memory>
#include <string>
#include <utility>
#include "absl/memory/memory.h"
#include "quiche/quic/core/crypto/crypto_protocol.h"
#include "quiche/quic/core/http/quic_server_initiated_spdy_stream.h"
#include "quiche/quic/core/http/quic_spdy_client_stream.h"
#include "quiche/quic/core/http/spdy_utils.h"
#include "quiche/quic/core/quic_server_id.h"
#include "quiche/quic/core/quic_utils.h"
#include "quiche/quic/platform/api/quic_bug_tracker.h"
#include "quiche/quic/platform/api/quic_flag_utils.h"
#include "quiche/quic/platform/api/quic_flags.h"
#include "quiche/quic/platform/api/quic_logging.h"
namespace quic {
QuicSpdyClientSession::QuicSpdyClientSession(
    const QuicConfig& config, const ParsedQuicVersionVector& supported_versions,
    QuicConnection* connection, const QuicServerId& server_id,
    QuicCryptoClientConfig* crypto_config)
    : QuicSpdyClientSession(config, supported_versions, connection, nullptr,
                            server_id, crypto_config) {}
QuicSpdyClientSession::QuicSpdyClientSession(
    const QuicConfig& config, const ParsedQuicVersionVector& supported_versions,
    QuicConnection* connection, QuicSession::Visitor* visitor,
    const QuicServerId& server_id, QuicCryptoClientConfig* crypto_config)
    : QuicSpdyClientSessionBase(connection, visitor, config,
                                supported_versions),
      server_id_(server_id),
      crypto_config_(crypto_config),
      respect_goaway_(true) {}
QuicSpdyClientSession::~QuicSpdyClientSession() = default;
void QuicSpdyClientSession::Initialize() {
  crypto_stream_ = CreateQuicCryptoStream();
  QuicSpdyClientSessionBase::Initialize();
}
void QuicSpdyClientSession::OnProofValid(
    const QuicCryptoClientConfig::CachedState& ) {}
void QuicSpdyClientSession::OnProofVerifyDetailsAvailable(
    const ProofVerifyDetails& ) {}
bool QuicSpdyClientSession::ShouldCreateOutgoingBidirectionalStream() {
  if (!crypto_stream_->encryption_established()) {
    QUIC_DLOG(INFO) << "Encryption not active so no outgoing stream created.";
    QUIC_CODE_COUNT(
        quic_client_fails_to_create_stream_encryption_not_established);
    return false;
  }
  if (goaway_received() && respect_goaway_) {
    QUIC_DLOG(INFO) << "Failed to create a new outgoing stream. "
                    << "Already received goaway.";
    QUIC_CODE_COUNT(quic_client_fails_to_create_stream_goaway_received);
    return false;
  }
  return CanOpenNextOutgoingBidirectionalStream();
}
bool QuicSpdyClientSession::ShouldCreateOutgoingUnidirectionalStream() {
  QUIC_BUG(quic_bug_10396_1)
      << "Try to create outgoing unidirectional client data streams";
  return false;
}
QuicSpdyClientStream*
QuicSpdyClientSession::CreateOutgoingBidirectionalStream() {
  if (!ShouldCreateOutgoingBidirectionalStream()) {
    return nullptr;
  }
  std::unique_ptr<QuicSpdyClientStream> stream = CreateClientStream();
  QuicSpdyClientStream* stream_ptr = stream.get();
  ActivateStream(std::move(stream));
  return stream_ptr;
}
QuicSpdyClientStream*
QuicSpdyClientSession::CreateOutgoingUnidirectionalStream() {
  QUIC_BUG(quic_bug_10396_2)
      << "Try to create outgoing unidirectional client data streams";
  return nullptr;
}
std::unique_ptr<QuicSpdyClientStream>
QuicSpdyClientSession::CreateClientStream() {
  return std::make_unique<QuicSpdyClientStream>(
      GetNextOutgoingBidirectionalStreamId(), this, BIDIRECTIONAL);
}
QuicCryptoClientStreamBase* QuicSpdyClientSession::GetMutableCryptoStream() {
  return crypto_stream_.get();
}
const QuicCryptoClientStreamBase* QuicSpdyClientSession::GetCryptoStream()
    const {
  return crypto_stream_.get();
}
void QuicSpdyClientSession::CryptoConnect() {
  QUICHE_DCHECK(flow_controller());
  crypto_stream_->CryptoConnect();
}
int QuicSpdyClientSession::GetNumSentClientHellos() const {
  return crypto_stream_->num_sent_client_hellos();
}
bool QuicSpdyClientSession::ResumptionAttempted() const {
  return crypto_stream_->ResumptionAttempted();
}
bool QuicSpdyClientSession::IsResumption() const {
  return crypto_stream_->IsResumption();
}
bool QuicSpdyClientSession::EarlyDataAccepted() const {
  return crypto_stream_->EarlyDataAccepted();
}
bool QuicSpdyClientSession::ReceivedInchoateReject() const {
  return crypto_stream_->ReceivedInchoateReject();
}
int QuicSpdyClientSession::GetNumReceivedServerConfigUpdates() const {
  return crypto_stream_->num_scup_messages_received();
}
bool QuicSpdyClientSession::ShouldCreateIncomingStream(QuicStreamId id) {
  if (!connection()->connected()) {
    QUIC_BUG(quic_bug_10396_3)
        << "ShouldCreateIncomingStream called when disconnected";
    return false;
  }
  if (goaway_received() && respect_goaway_) {
    QUIC_DLOG(INFO) << "Failed to create a new outgoing stream. "
                    << "Already received goaway.";
    return false;
  }
  if (QuicUtils::IsClientInitiatedStreamId(transport_version(), id)) {
    QUIC_BUG(quic_bug_10396_4)
        << "ShouldCreateIncomingStream called with client initiated "
           "stream ID.";
    return false;
  }
  if (QuicUtils::IsClientInitiatedStreamId(transport_version(), id)) {
    QUIC_LOG(WARNING) << "Received invalid push stream id " << id;
    connection()->CloseConnection(
        QUIC_INVALID_STREAM_ID,
        "Server created non write unidirectional stream",
        ConnectionCloseBehavior::SEND_CONNECTION_CLOSE_PACKET);
    return false;
  }
  if (VersionHasIetfQuicFrames(transport_version()) &&
      QuicUtils::IsBidirectionalStreamId(id, version()) &&
      !WillNegotiateWebTransport()) {
    connection()->CloseConnection(
        QUIC_HTTP_SERVER_INITIATED_BIDIRECTIONAL_STREAM,
        "Server created bidirectional stream.",
        ConnectionCloseBehavior::SEND_CONNECTION_CLOSE_PACKET);
    return false;
  }
  return true;
}
QuicSpdyStream* QuicSpdyClientSession::CreateIncomingStream(
    PendingStream* pending) {
  QuicSpdyStream* stream = new QuicSpdyClientStream(pending, this);
  ActivateStream(absl::WrapUnique(stream));
  return stream;
}
QuicSpdyStream* QuicSpdyClientSession::CreateIncomingStream(QuicStreamId id) {
  if (!ShouldCreateIncomingStream(id)) {
    return nullptr;
  }
  QuicSpdyStream* stream;
  if (version().UsesHttp3() &&
      QuicUtils::IsBidirectionalStreamId(id, version())) {
    QUIC_BUG_IF(QuicServerInitiatedSpdyStream but no WebTransport support,
                !WillNegotiateWebTransport())
        << "QuicServerInitiatedSpdyStream created but no WebTransport support";
    stream = new QuicServerInitiatedSpdyStream(id, this, BIDIRECTIONAL);
  } else {
    stream = new QuicSpdyClientStream(id, this, READ_UNIDIRECTIONAL);
  }
  ActivateStream(absl::WrapUnique(stream));
  return stream;
}
std::unique_ptr<QuicCryptoClientStreamBase>
QuicSpdyClientSession::CreateQuicCryptoStream() {
  return std::make_unique<QuicCryptoClientStream>(
      server_id_, this,
      crypto_config_->proof_verifier()->CreateDefaultContext(), crypto_config_,
      this,  version().UsesHttp3());
}
}  