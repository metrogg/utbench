#ifndef QUICHE_QUIC_TOOLS_QUIC_SIMPLE_SERVER_SESSION_H_
#define QUICHE_QUIC_TOOLS_QUIC_SIMPLE_SERVER_SESSION_H_
#include <stdint.h>
#include <list>
#include <memory>
#include <set>
#include <string>
#include <utility>
#include <vector>
#include "quiche/quic/core/http/quic_server_session_base.h"
#include "quiche/quic/core/http/quic_spdy_session.h"
#include "quiche/quic/core/quic_crypto_server_stream_base.h"
#include "quiche/quic/core/quic_packets.h"
#include "quiche/quic/tools/quic_backend_response.h"
#include "quiche/quic/tools/quic_simple_server_backend.h"
#include "quiche/quic/tools/quic_simple_server_stream.h"
#include "quiche/spdy/core/http2_header_block.h"
namespace quic {
namespace test {
class QuicSimpleServerSessionPeer;
}  
class QuicSimpleServerSession : public QuicServerSessionBase {
 public:
  QuicSimpleServerSession(const QuicConfig& config,
                          const ParsedQuicVersionVector& supported_versions,
                          QuicConnection* connection,
                          QuicSession::Visitor* visitor,
                          QuicCryptoServerStreamBase::Helper* helper,
                          const QuicCryptoServerConfig* crypto_config,
                          QuicCompressedCertsCache* compressed_certs_cache,
                          QuicSimpleServerBackend* quic_simple_server_backend);
  QuicSimpleServerSession(const QuicSimpleServerSession&) = delete;
  QuicSimpleServerSession& operator=(const QuicSimpleServerSession&) = delete;
  ~QuicSimpleServerSession() override;
  void OnStreamFrame(const QuicStreamFrame& frame) override;
 protected:
  QuicSpdyStream* CreateIncomingStream(QuicStreamId id) override;
  QuicSpdyStream* CreateIncomingStream(PendingStream* pending) override;
  QuicSpdyStream* CreateOutgoingBidirectionalStream() override;
  QuicSimpleServerStream* CreateOutgoingUnidirectionalStream() override;
  std::unique_ptr<QuicCryptoServerStreamBase> CreateQuicCryptoServerStream(
      const QuicCryptoServerConfig* crypto_config,
      QuicCompressedCertsCache* compressed_certs_cache) override;
  QuicStream* ProcessBidirectionalPendingStream(
      PendingStream* pending) override;
  QuicSimpleServerBackend* server_backend() {
    return quic_simple_server_backend_;
  }
  WebTransportHttp3VersionSet LocallySupportedWebTransportVersions()
      const override {
    return quic_simple_server_backend_->SupportsWebTransport()
               ? kDefaultSupportedWebTransportVersions
               : WebTransportHttp3VersionSet();
  }
  HttpDatagramSupport LocalHttpDatagramSupport() override {
    if (ShouldNegotiateWebTransport()) {
      return HttpDatagramSupport::kRfcAndDraft04;
    }
    return QuicServerSessionBase::LocalHttpDatagramSupport();
  }
 private:
  friend class test::QuicSimpleServerSessionPeer;
  QuicSimpleServerBackend* quic_simple_server_backend_;  
};
}  
#endif  
#include "quiche/quic/tools/quic_simple_server_session.h"
#include <memory>
#include <utility>
#include "absl/memory/memory.h"
#include "quiche/quic/core/http/quic_server_initiated_spdy_stream.h"
#include "quiche/quic/core/http/quic_spdy_session.h"
#include "quiche/quic/core/quic_connection.h"
#include "quiche/quic/core/quic_stream_priority.h"
#include "quiche/quic/core/quic_types.h"
#include "quiche/quic/core/quic_utils.h"
#include "quiche/quic/platform/api/quic_flags.h"
#include "quiche/quic/platform/api/quic_logging.h"
#include "quiche/quic/tools/quic_simple_server_stream.h"
namespace quic {
QuicSimpleServerSession::QuicSimpleServerSession(
    const QuicConfig& config, const ParsedQuicVersionVector& supported_versions,
    QuicConnection* connection, QuicSession::Visitor* visitor,
    QuicCryptoServerStreamBase::Helper* helper,
    const QuicCryptoServerConfig* crypto_config,
    QuicCompressedCertsCache* compressed_certs_cache,
    QuicSimpleServerBackend* quic_simple_server_backend)
    : QuicServerSessionBase(config, supported_versions, connection, visitor,
                            helper, crypto_config, compressed_certs_cache),
      quic_simple_server_backend_(quic_simple_server_backend) {
  QUICHE_DCHECK(quic_simple_server_backend_);
  set_max_streams_accepted_per_loop(5u);
}
QuicSimpleServerSession::~QuicSimpleServerSession() { DeleteConnection(); }
std::unique_ptr<QuicCryptoServerStreamBase>
QuicSimpleServerSession::CreateQuicCryptoServerStream(
    const QuicCryptoServerConfig* crypto_config,
    QuicCompressedCertsCache* compressed_certs_cache) {
  return CreateCryptoServerStream(crypto_config, compressed_certs_cache, this,
                                  stream_helper());
}
void QuicSimpleServerSession::OnStreamFrame(const QuicStreamFrame& frame) {
  if (!IsIncomingStream(frame.stream_id) && !WillNegotiateWebTransport()) {
    QUIC_LOG(WARNING) << "Client shouldn't send data on server push stream";
    connection()->CloseConnection(
        QUIC_INVALID_STREAM_ID, "Client sent data on server push stream",
        ConnectionCloseBehavior::SEND_CONNECTION_CLOSE_PACKET);
    return;
  }
  QuicSpdySession::OnStreamFrame(frame);
}
QuicSpdyStream* QuicSimpleServerSession::CreateIncomingStream(QuicStreamId id) {
  if (!ShouldCreateIncomingStream(id)) {
    return nullptr;
  }
  QuicSpdyStream* stream = new QuicSimpleServerStream(
      id, this, BIDIRECTIONAL, quic_simple_server_backend_);
  ActivateStream(absl::WrapUnique(stream));
  return stream;
}
QuicSpdyStream* QuicSimpleServerSession::CreateIncomingStream(
    PendingStream* pending) {
  QuicSpdyStream* stream =
      new QuicSimpleServerStream(pending, this, quic_simple_server_backend_);
  ActivateStream(absl::WrapUnique(stream));
  return stream;
}
QuicSpdyStream* QuicSimpleServerSession::CreateOutgoingBidirectionalStream() {
  if (!WillNegotiateWebTransport()) {
    QUIC_BUG(QuicSimpleServerSession CreateOutgoingBidirectionalStream without
                 WebTransport support)
        << "QuicSimpleServerSession::CreateOutgoingBidirectionalStream called "
           "in a session without WebTransport support.";
    return nullptr;
  }
  if (!ShouldCreateOutgoingBidirectionalStream()) {
    return nullptr;
  }
  QuicServerInitiatedSpdyStream* stream = new QuicServerInitiatedSpdyStream(
      GetNextOutgoingBidirectionalStreamId(), this, BIDIRECTIONAL);
  ActivateStream(absl::WrapUnique(stream));
  return stream;
}
QuicSimpleServerStream*
QuicSimpleServerSession::CreateOutgoingUnidirectionalStream() {
  if (!ShouldCreateOutgoingUnidirectionalStream()) {
    return nullptr;
  }
  QuicSimpleServerStream* stream = new QuicSimpleServerStream(
      GetNextOutgoingUnidirectionalStreamId(), this, WRITE_UNIDIRECTIONAL,
      quic_simple_server_backend_);
  ActivateStream(absl::WrapUnique(stream));
  return stream;
}
QuicStream* QuicSimpleServerSession::ProcessBidirectionalPendingStream(
    PendingStream* pending) {
  QUICHE_DCHECK(IsEncryptionEstablished());
  return CreateIncomingStream(pending);
}
}  