#ifndef QUICHE_QUIC_CORE_CRYPTO_CLIENT_PROOF_SOURCE_H_
#define QUICHE_QUIC_CORE_CRYPTO_CLIENT_PROOF_SOURCE_H_
#include <memory>
#include "absl/container/flat_hash_map.h"
#include "quiche/quic/core/crypto/certificate_view.h"
#include "quiche/quic/core/crypto/proof_source.h"
namespace quic {
class QUICHE_EXPORT ClientProofSource {
 public:
  using Chain = ProofSource::Chain;
  virtual ~ClientProofSource() {}
  struct QUICHE_EXPORT CertAndKey {
    CertAndKey(quiche::QuicheReferenceCountedPointer<Chain> chain,
               CertificatePrivateKey private_key)
        : chain(std::move(chain)), private_key(std::move(private_key)) {}
    quiche::QuicheReferenceCountedPointer<Chain> chain;
    CertificatePrivateKey private_key;
  };
  virtual std::shared_ptr<const CertAndKey> GetCertAndKey(
      absl::string_view server_hostname) const = 0;
};
class QUICHE_EXPORT DefaultClientProofSource : public ClientProofSource {
 public:
  ~DefaultClientProofSource() override {}
  bool AddCertAndKey(std::vector<std::string> server_hostnames,
                     quiche::QuicheReferenceCountedPointer<Chain> chain,
                     CertificatePrivateKey private_key);
  std::shared_ptr<const CertAndKey> GetCertAndKey(
      absl::string_view hostname) const override;
 private:
  std::shared_ptr<const CertAndKey> LookupExact(
      absl::string_view map_key) const;
  absl::flat_hash_map<std::string, std::shared_ptr<CertAndKey>> cert_and_keys_;
};
}  
#endif  
#include "quiche/quic/core/crypto/client_proof_source.h"
#include <memory>
#include <string>
#include <utility>
#include <vector>
#include "absl/strings/match.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/string_view.h"
namespace quic {
bool DefaultClientProofSource::AddCertAndKey(
    std::vector<std::string> server_hostnames,
    quiche::QuicheReferenceCountedPointer<Chain> chain,
    CertificatePrivateKey private_key) {
  if (!ValidateCertAndKey(chain, private_key)) {
    return false;
  }
  auto cert_and_key =
      std::make_shared<CertAndKey>(std::move(chain), std::move(private_key));
  for (const std::string& domain : server_hostnames) {
    cert_and_keys_[domain] = cert_and_key;
  }
  return true;
}
std::shared_ptr<const ClientProofSource::CertAndKey>
DefaultClientProofSource::GetCertAndKey(absl::string_view hostname) const {
  if (std::shared_ptr<const CertAndKey> result = LookupExact(hostname);
      result || hostname == "*") {
    return result;
  }
  if (hostname.size() > 1 && !absl::StartsWith(hostname, "*.")) {
    auto dot_pos = hostname.find('.');
    if (dot_pos != std::string::npos) {
      std::string wildcard = absl::StrCat("*", hostname.substr(dot_pos));
      std::shared_ptr<const CertAndKey> result = LookupExact(wildcard);
      if (result != nullptr) {
        return result;
      }
    }
  }
  return LookupExact("*");
}
std::shared_ptr<const ClientProofSource::CertAndKey>
DefaultClientProofSource::LookupExact(absl::string_view map_key) const {
  const auto it = cert_and_keys_.find(map_key);
  QUIC_DVLOG(1) << "LookupExact(" << map_key
                << ") found:" << (it != cert_and_keys_.end());
  if (it != cert_and_keys_.end()) {
    return it->second;
  }
  return nullptr;
}
}  