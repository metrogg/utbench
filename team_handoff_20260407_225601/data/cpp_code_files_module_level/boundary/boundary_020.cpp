#ifndef QUICHE_QUIC_CORE_QUIC_DATA_WRITER_H_
#define QUICHE_QUIC_CORE_QUIC_DATA_WRITER_H_
#include <cstddef>
#include <cstdint>
#include "absl/strings/string_view.h"
#include "quiche/quic/core/crypto/quic_random.h"
#include "quiche/quic/core/quic_types.h"
#include "quiche/quic/platform/api/quic_export.h"
#include "quiche/common/quiche_data_writer.h"
#include "quiche/common/quiche_endian.h"
namespace quic {
class QUICHE_EXPORT QuicDataWriter : public quiche::QuicheDataWriter {
 public:
  QuicDataWriter(size_t size, char* buffer);
  QuicDataWriter(size_t size, char* buffer, quiche::Endianness endianness);
  QuicDataWriter(const QuicDataWriter&) = delete;
  QuicDataWriter& operator=(const QuicDataWriter&) = delete;
  ~QuicDataWriter();
  bool WriteUFloat16(uint64_t value);
  bool WriteConnectionId(QuicConnectionId connection_id);
  bool WriteLengthPrefixedConnectionId(QuicConnectionId connection_id);
  bool WriteRandomBytes(QuicRandom* random, size_t length);
  bool WriteInsecureRandomBytes(QuicRandom* random, size_t length);
};
}  
#endif  
#include "quiche/quic/core/quic_data_writer.h"
#include <algorithm>
#include <limits>
#include "absl/strings/string_view.h"
#include "quiche/quic/core/crypto/quic_random.h"
#include "quiche/quic/core/quic_constants.h"
#include "quiche/quic/platform/api/quic_bug_tracker.h"
#include "quiche/quic/platform/api/quic_flags.h"
#include "quiche/common/quiche_endian.h"
namespace quic {
QuicDataWriter::QuicDataWriter(size_t size, char* buffer)
    : quiche::QuicheDataWriter(size, buffer) {}
QuicDataWriter::QuicDataWriter(size_t size, char* buffer,
                               quiche::Endianness endianness)
    : quiche::QuicheDataWriter(size, buffer, endianness) {}
QuicDataWriter::~QuicDataWriter() {}
bool QuicDataWriter::WriteUFloat16(uint64_t value) {
  uint16_t result;
  if (value < (UINT64_C(1) << kUFloat16MantissaEffectiveBits)) {
    result = static_cast<uint16_t>(value);
  } else if (value >= kUFloat16MaxValue) {
    result = std::numeric_limits<uint16_t>::max();
  } else {
    uint16_t exponent = 0;
    for (uint16_t offset = 16; offset > 0; offset /= 2) {
      if (value >= (UINT64_C(1) << (kUFloat16MantissaBits + offset))) {
        exponent += offset;
        value >>= offset;
      }
    }
    QUICHE_DCHECK_GE(exponent, 1);
    QUICHE_DCHECK_LE(exponent, kUFloat16MaxExponent);
    QUICHE_DCHECK_GE(value, UINT64_C(1) << kUFloat16MantissaBits);
    QUICHE_DCHECK_LT(value, UINT64_C(1) << kUFloat16MantissaEffectiveBits);
    result = static_cast<uint16_t>(value + (exponent << kUFloat16MantissaBits));
  }
  if (endianness() == quiche::NETWORK_BYTE_ORDER) {
    result = quiche::QuicheEndian::HostToNet16(result);
  }
  return WriteBytes(&result, sizeof(result));
}
bool QuicDataWriter::WriteConnectionId(QuicConnectionId connection_id) {
  if (connection_id.IsEmpty()) {
    return true;
  }
  return WriteBytes(connection_id.data(), connection_id.length());
}
bool QuicDataWriter::WriteLengthPrefixedConnectionId(
    QuicConnectionId connection_id) {
  return WriteUInt8(connection_id.length()) && WriteConnectionId(connection_id);
}
bool QuicDataWriter::WriteRandomBytes(QuicRandom* random, size_t length) {
  char* dest = BeginWrite(length);
  if (!dest) {
    return false;
  }
  random->RandBytes(dest, length);
  IncreaseLength(length);
  return true;
}
bool QuicDataWriter::WriteInsecureRandomBytes(QuicRandom* random,
                                              size_t length) {
  char* dest = BeginWrite(length);
  if (!dest) {
    return false;
  }
  random->InsecureRandBytes(dest, length);
  IncreaseLength(length);
  return true;
}
}  