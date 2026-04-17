#ifndef QUICHE_QUIC_CORE_CRYPTO_CERTIFICATE_UTIL_H_
#define QUICHE_QUIC_CORE_CRYPTO_CERTIFICATE_UTIL_H_
#include <string>
#include "absl/strings/string_view.h"
#include "openssl/evp.h"
#include "quiche/quic/core/quic_time.h"
#include "quiche/quic/platform/api/quic_export.h"
namespace quic {
struct QUICHE_EXPORT CertificateTimestamp {
  uint16_t year;
  uint8_t month;
  uint8_t day;
  uint8_t hour;
  uint8_t minute;
  uint8_t second;
};
struct QUICHE_EXPORT CertificateOptions {
  absl::string_view subject;
  uint64_t serial_number;
  CertificateTimestamp validity_start;  
  CertificateTimestamp validity_end;    
};
QUICHE_EXPORT bssl::UniquePtr<EVP_PKEY> MakeKeyPairForSelfSignedCertificate();
QUICHE_EXPORT std::string CreateSelfSignedCertificate(
    EVP_PKEY& key, const CertificateOptions& options);
}  
#endif  
#include "quiche/quic/core/crypto/certificate_util.h"
#include <string>
#include <vector>
#include "absl/strings/str_format.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"
#include "openssl/bn.h"
#include "openssl/bytestring.h"
#include "openssl/digest.h"
#include "openssl/ec_key.h"
#include "openssl/mem.h"
#include "openssl/pkcs7.h"
#include "openssl/pool.h"
#include "openssl/rsa.h"
#include "openssl/stack.h"
#include "quiche/quic/core/crypto/boring_utils.h"
#include "quiche/quic/platform/api/quic_logging.h"
namespace quic {
namespace {
bool AddEcdsa256SignatureAlgorithm(CBB* cbb) {
  static const uint8_t kEcdsaWithSha256[] = {0x2a, 0x86, 0x48, 0xce,
                                             0x3d, 0x04, 0x03, 0x02};
  CBB sequence, oid;
  if (!CBB_add_asn1(cbb, &sequence, CBS_ASN1_SEQUENCE) ||
      !CBB_add_asn1(&sequence, &oid, CBS_ASN1_OBJECT)) {
    return false;
  }
  if (!CBB_add_bytes(&oid, kEcdsaWithSha256, sizeof(kEcdsaWithSha256))) {
    return false;
  }
  return CBB_flush(cbb);
}
bool AddName(CBB* cbb, absl::string_view name) {
  static const uint8_t kCommonName[] = {0x55, 0x04, 0x03};
  static const uint8_t kCountryName[] = {0x55, 0x04, 0x06};
  static const uint8_t kOrganizationName[] = {0x55, 0x04, 0x0a};
  static const uint8_t kOrganizationalUnitName[] = {0x55, 0x04, 0x0b};
  std::vector<std::string> attributes =
      absl::StrSplit(name, ',', absl::SkipEmpty());
  if (attributes.empty()) {
    QUIC_LOG(ERROR) << "Missing DN or wrong format";
    return false;
  }
  CBB rdns;
  if (!CBB_add_asn1(cbb, &rdns, CBS_ASN1_SEQUENCE)) {
    return false;
  }
  for (const std::string& attribute : attributes) {
    std::vector<std::string> parts =
        absl::StrSplit(absl::StripAsciiWhitespace(attribute), '=');
    if (parts.size() != 2) {
      QUIC_LOG(ERROR) << "Wrong DN format at " + attribute;
      return false;
    }
    const std::string& type_string = parts[0];
    const std::string& value_string = parts[1];
    absl::Span<const uint8_t> type_bytes;
    if (type_string == "CN") {
      type_bytes = kCommonName;
    } else if (type_string == "C") {
      type_bytes = kCountryName;
    } else if (type_string == "O") {
      type_bytes = kOrganizationName;
    } else if (type_string == "OU") {
      type_bytes = kOrganizationalUnitName;
    } else {
      QUIC_LOG(ERROR) << "Unrecognized type " + type_string;
      return false;
    }
    CBB rdn, attr, type, value;
    if (!CBB_add_asn1(&rdns, &rdn, CBS_ASN1_SET) ||
        !CBB_add_asn1(&rdn, &attr, CBS_ASN1_SEQUENCE) ||
        !CBB_add_asn1(&attr, &type, CBS_ASN1_OBJECT) ||
        !CBB_add_bytes(&type, type_bytes.data(), type_bytes.size()) ||
        !CBB_add_asn1(&attr, &value,
                      type_string == "C" ? CBS_ASN1_PRINTABLESTRING
                                         : CBS_ASN1_UTF8STRING) ||
        !AddStringToCbb(&value, value_string) || !CBB_flush(&rdns)) {
      return false;
    }
  }
  if (!CBB_flush(cbb)) {
    return false;
  }
  return true;
}
bool CBBAddTime(CBB* cbb, const CertificateTimestamp& timestamp) {
  CBB child;
  std::string formatted_time;
  const bool is_utc_time = (1950 <= timestamp.year && timestamp.year < 2050);
  if (is_utc_time) {
    uint16_t year = timestamp.year - 1900;
    if (year >= 100) {
      year -= 100;
    }
    formatted_time = absl::StrFormat("%02d", year);
    if (!CBB_add_asn1(cbb, &child, CBS_ASN1_UTCTIME)) {
      return false;
    }
  } else {
    formatted_time = absl::StrFormat("%04d", timestamp.year);
    if (!CBB_add_asn1(cbb, &child, CBS_ASN1_GENERALIZEDTIME)) {
      return false;
    }
  }
  absl::StrAppendFormat(&formatted_time, "%02d%02d%02d%02d%02dZ",
                        timestamp.month, timestamp.day, timestamp.hour,
                        timestamp.minute, timestamp.second);
  static const size_t kGeneralizedTimeLength = 15;
  static const size_t kUTCTimeLength = 13;
  QUICHE_DCHECK_EQ(formatted_time.size(),
                   is_utc_time ? kUTCTimeLength : kGeneralizedTimeLength);
  return AddStringToCbb(&child, formatted_time) && CBB_flush(cbb);
}
bool CBBAddExtension(CBB* extensions, absl::Span<const uint8_t> oid,
                     bool critical, absl::Span<const uint8_t> contents) {
  CBB extension, cbb_oid, cbb_contents;
  if (!CBB_add_asn1(extensions, &extension, CBS_ASN1_SEQUENCE) ||
      !CBB_add_asn1(&extension, &cbb_oid, CBS_ASN1_OBJECT) ||
      !CBB_add_bytes(&cbb_oid, oid.data(), oid.size()) ||
      (critical && !CBB_add_asn1_bool(&extension, 1)) ||
      !CBB_add_asn1(&extension, &cbb_contents, CBS_ASN1_OCTETSTRING) ||
      !CBB_add_bytes(&cbb_contents, contents.data(), contents.size()) ||
      !CBB_flush(extensions)) {
    return false;
  }
  return true;
}
bool IsEcdsa256Key(const EVP_PKEY& evp_key) {
  if (EVP_PKEY_id(&evp_key) != EVP_PKEY_EC) {
    return false;
  }
  const EC_KEY* key = EVP_PKEY_get0_EC_KEY(&evp_key);
  if (key == nullptr) {
    return false;
  }
  const EC_GROUP* group = EC_KEY_get0_group(key);
  if (group == nullptr) {
    return false;
  }
  return EC_GROUP_get_curve_name(group) == NID_X9_62_prime256v1;
}
}  
bssl::UniquePtr<EVP_PKEY> MakeKeyPairForSelfSignedCertificate() {
  bssl::UniquePtr<EVP_PKEY_CTX> context(
      EVP_PKEY_CTX_new_id(EVP_PKEY_EC, nullptr));
  if (!context) {
    return nullptr;
  }
  if (EVP_PKEY_keygen_init(context.get()) != 1) {
    return nullptr;
  }
  if (EVP_PKEY_CTX_set_ec_paramgen_curve_nid(context.get(),
                                             NID_X9_62_prime256v1) != 1) {
    return nullptr;
  }
  EVP_PKEY* raw_key = nullptr;
  if (EVP_PKEY_keygen(context.get(), &raw_key) != 1) {
    return nullptr;
  }
  return bssl::UniquePtr<EVP_PKEY>(raw_key);
}
std::string CreateSelfSignedCertificate(EVP_PKEY& key,
                                        const CertificateOptions& options) {
  std::string error;
  if (!IsEcdsa256Key(key)) {
    QUIC_LOG(ERROR) << "CreateSelfSignedCert only accepts ECDSA P-256 keys";
    return error;
  }
  bssl::ScopedCBB cbb;
  CBB tbs_cert, version, validity;
  uint8_t* tbs_cert_bytes;
  size_t tbs_cert_len;
  if (!CBB_init(cbb.get(), 64) ||
      !CBB_add_asn1(cbb.get(), &tbs_cert, CBS_ASN1_SEQUENCE) ||
      !CBB_add_asn1(&tbs_cert, &version,
                    CBS_ASN1_CONTEXT_SPECIFIC | CBS_ASN1_CONSTRUCTED | 0) ||
      !CBB_add_asn1_uint64(&version, 2) ||  
      !CBB_add_asn1_uint64(&tbs_cert, options.serial_number) ||
      !AddEcdsa256SignatureAlgorithm(&tbs_cert) ||  
      !AddName(&tbs_cert, options.subject) ||       
      !CBB_add_asn1(&tbs_cert, &validity, CBS_ASN1_SEQUENCE) ||
      !CBBAddTime(&validity, options.validity_start) ||
      !CBBAddTime(&validity, options.validity_end) ||
      !AddName(&tbs_cert, options.subject) ||      
      !EVP_marshal_public_key(&tbs_cert, &key)) {  
    return error;
  }
  CBB outer_extensions, extensions;
  if (!CBB_add_asn1(&tbs_cert, &outer_extensions,
                    3 | CBS_ASN1_CONTEXT_SPECIFIC | CBS_ASN1_CONSTRUCTED) ||
      !CBB_add_asn1(&outer_extensions, &extensions, CBS_ASN1_SEQUENCE)) {
    return error;
  }
  constexpr uint8_t kKeyUsageOid[] = {0x55, 0x1d, 0x0f};
  constexpr uint8_t kKeyUsageContent[] = {
      0x3,   
      0x2,   
      0x0,   
      0x80,  
  };
  CBBAddExtension(&extensions, kKeyUsageOid, true, kKeyUsageContent);
  if (!CBB_finish(cbb.get(), &tbs_cert_bytes, &tbs_cert_len)) {
    return error;
  }
  bssl::UniquePtr<uint8_t> delete_tbs_cert_bytes(tbs_cert_bytes);
  CBB cert, signature;
  bssl::ScopedEVP_MD_CTX ctx;
  uint8_t* sig_out;
  size_t sig_len;
  uint8_t* cert_bytes;
  size_t cert_len;
  if (!CBB_init(cbb.get(), tbs_cert_len) ||
      !CBB_add_asn1(cbb.get(), &cert, CBS_ASN1_SEQUENCE) ||
      !CBB_add_bytes(&cert, tbs_cert_bytes, tbs_cert_len) ||
      !AddEcdsa256SignatureAlgorithm(&cert) ||
      !CBB_add_asn1(&cert, &signature, CBS_ASN1_BITSTRING) ||
      !CBB_add_u8(&signature, 0 ) ||
      !EVP_DigestSignInit(ctx.get(), nullptr, EVP_sha256(), nullptr, &key) ||
      !EVP_DigestSign(ctx.get(), nullptr, &sig_len, tbs_cert_bytes,
                      tbs_cert_len) ||
      !CBB_reserve(&signature, &sig_out, sig_len) ||
      !EVP_DigestSign(ctx.get(), sig_out, &sig_len, tbs_cert_bytes,
                      tbs_cert_len) ||
      !CBB_did_write(&signature, sig_len) ||
      !CBB_finish(cbb.get(), &cert_bytes, &cert_len)) {
    return error;
  }
  bssl::UniquePtr<uint8_t> delete_cert_bytes(cert_bytes);
  return std::string(reinterpret_cast<char*>(cert_bytes), cert_len);
}
}  