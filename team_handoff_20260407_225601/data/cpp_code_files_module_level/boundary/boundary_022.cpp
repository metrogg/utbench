#ifndef TENSORFLOW_TSL_LIB_IO_BUFFERED_INPUTSTREAM_H_
#define TENSORFLOW_TSL_LIB_IO_BUFFERED_INPUTSTREAM_H_
#include <string>
#include "tsl/lib/io/inputstream_interface.h"
#include "tsl/platform/file_system.h"
namespace tsl {
namespace io {
class BufferedInputStream : public InputStreamInterface {
 public:
  BufferedInputStream(InputStreamInterface* input_stream, size_t buffer_bytes,
                      bool owns_input_stream = false);
  BufferedInputStream(RandomAccessFile* file, size_t buffer_bytes);
  ~BufferedInputStream() override;
  absl::Status ReadNBytes(int64_t bytes_to_read, tstring* result) override;
  absl::Status SkipNBytes(int64_t bytes_to_skip) override;
  int64_t Tell() const override;
  absl::Status Seek(int64_t position);
  absl::Status ReadLine(std::string* result);
  absl::Status ReadLine(tstring* result);
  std::string ReadLineAsString();
  absl::Status SkipLine();
  template <typename T>
  absl::Status ReadAll(T* result);
  absl::Status Reset() override;
 private:
  absl::Status FillBuffer();
  template <typename StringType>
  absl::Status ReadLineHelper(StringType* result, bool include_eol);
  InputStreamInterface* input_stream_;  
  size_t size_;                         
  tstring buf_;                         
  size_t pos_ = 0;    
  size_t limit_ = 0;  
  bool owns_input_stream_ = false;
  absl::Status file_status_ = absl::OkStatus();
  BufferedInputStream(const BufferedInputStream&) = delete;
  void operator=(const BufferedInputStream&) = delete;
};
#ifndef SWIG
extern template Status BufferedInputStream::ReadAll<std::string>(
    std::string* result);
extern template Status BufferedInputStream::ReadAll<tstring>(tstring* result);
#endif  
}  
}  
#endif  
#include "tsl/lib/io/buffered_inputstream.h"
#include "absl/status/status.h"
#include "tsl/lib/io/random_inputstream.h"
namespace tsl {
namespace io {
BufferedInputStream::BufferedInputStream(InputStreamInterface* input_stream,
                                         size_t buffer_bytes,
                                         bool owns_input_stream)
    : input_stream_(input_stream),
      size_(buffer_bytes),
      owns_input_stream_(owns_input_stream) {
  buf_.reserve(size_);
}
BufferedInputStream::BufferedInputStream(RandomAccessFile* file,
                                         size_t buffer_bytes)
    : BufferedInputStream(new RandomAccessInputStream(file), buffer_bytes,
                          true) {}
BufferedInputStream::~BufferedInputStream() {
  if (owns_input_stream_) {
    delete input_stream_;
  }
}
absl::Status BufferedInputStream::FillBuffer() {
  if (!file_status_.ok()) {
    pos_ = 0;
    limit_ = 0;
    return file_status_;
  }
  absl::Status s = input_stream_->ReadNBytes(size_, &buf_);
  pos_ = 0;
  limit_ = buf_.size();
  if (!s.ok()) {
    file_status_ = s;
  }
  return s;
}
template <typename StringType>
absl::Status BufferedInputStream::ReadLineHelper(StringType* result,
                                                 bool include_eol) {
  result->clear();
  absl::Status s;
  size_t start_pos = pos_;
  while (true) {
    if (pos_ == limit_) {
      result->append(buf_.data() + start_pos, pos_ - start_pos);
      s = FillBuffer();
      if (limit_ == 0) {
        break;
      }
      start_pos = pos_;
    }
    char c = buf_[pos_];
    if (c == '\n') {
      result->append(buf_.data() + start_pos, pos_ - start_pos);
      if (include_eol) {
        result->append(1, c);
      }
      pos_++;
      return absl::OkStatus();
    }
    if (c == '\r') {
      result->append(buf_.data() + start_pos, pos_ - start_pos);
      start_pos = pos_ + 1;
    }
    pos_++;
  }
  if (absl::IsOutOfRange(s) && !result->empty()) {
    return absl::OkStatus();
  }
  return s;
}
absl::Status BufferedInputStream::ReadNBytes(int64_t bytes_to_read,
                                             tstring* result) {
  if (bytes_to_read < 0) {
    return errors::InvalidArgument("Can't read a negative number of bytes: ",
                                   bytes_to_read);
  }
  result->clear();
  if (pos_ == limit_ && !file_status_.ok() && bytes_to_read > 0) {
    return file_status_;
  }
  result->reserve(bytes_to_read);
  absl::Status s;
  while (result->size() < static_cast<size_t>(bytes_to_read)) {
    if (pos_ == limit_) {
      s = FillBuffer();
      if (limit_ == 0) {
        DCHECK(!s.ok());
        file_status_ = s;
        break;
      }
    }
    const int64_t bytes_to_copy =
        std::min<int64_t>(limit_ - pos_, bytes_to_read - result->size());
    result->insert(result->size(), buf_, pos_, bytes_to_copy);
    pos_ += bytes_to_copy;
  }
  if (absl::IsOutOfRange(s) &&
      (result->size() == static_cast<size_t>(bytes_to_read))) {
    return absl::OkStatus();
  }
  return s;
}
absl::Status BufferedInputStream::SkipNBytes(int64_t bytes_to_skip) {
  if (bytes_to_skip < 0) {
    return errors::InvalidArgument("Can only skip forward, not ",
                                   bytes_to_skip);
  }
  if (pos_ + bytes_to_skip < limit_) {
    pos_ += bytes_to_skip;
  } else {
    absl::Status s = input_stream_->SkipNBytes(bytes_to_skip - (limit_ - pos_));
    pos_ = 0;
    limit_ = 0;
    if (absl::IsOutOfRange(s)) {
      file_status_ = s;
    }
    return s;
  }
  return absl::OkStatus();
}
int64_t BufferedInputStream::Tell() const {
  return input_stream_->Tell() - (limit_ - pos_);
}
absl::Status BufferedInputStream::Seek(int64_t position) {
  if (position < 0) {
    return errors::InvalidArgument("Seeking to a negative position: ",
                                   position);
  }
  const int64_t buf_lower_limit = input_stream_->Tell() - limit_;
  if (position < buf_lower_limit) {
    TF_RETURN_IF_ERROR(Reset());
    return SkipNBytes(position);
  }
  if (position < Tell()) {
    pos_ -= Tell() - position;
    return absl::OkStatus();
  }
  return SkipNBytes(position - Tell());
}
template <typename T>
absl::Status BufferedInputStream::ReadAll(T* result) {
  result->clear();
  absl::Status status;
  while (status.ok()) {
    status = FillBuffer();
    if (limit_ == 0) {
      break;
    }
    result->append(buf_);
    pos_ = limit_;
  }
  if (absl::IsOutOfRange(status)) {
    file_status_ = status;
    return absl::OkStatus();
  }
  return status;
}
template Status BufferedInputStream::ReadAll<std::string>(std::string* result);
template Status BufferedInputStream::ReadAll<tstring>(tstring* result);
absl::Status BufferedInputStream::Reset() {
  TF_RETURN_IF_ERROR(input_stream_->Reset());
  pos_ = 0;
  limit_ = 0;
  file_status_ = absl::OkStatus();
  return absl::OkStatus();
}
absl::Status BufferedInputStream::ReadLine(std::string* result) {
  return ReadLineHelper(result, false);
}
absl::Status BufferedInputStream::ReadLine(tstring* result) {
  return ReadLineHelper(result, false);
}
std::string BufferedInputStream::ReadLineAsString() {
  std::string result;
  ReadLineHelper(&result, true).IgnoreError();
  return result;
}
absl::Status BufferedInputStream::SkipLine() {
  absl::Status s;
  bool skipped = false;
  while (true) {
    if (pos_ == limit_) {
      s = FillBuffer();
      if (limit_ == 0) {
        break;
      }
    }
    char c = buf_[pos_++];
    skipped = true;
    if (c == '\n') {
      return absl::OkStatus();
    }
  }
  if (absl::IsOutOfRange(s) && skipped) {
    return absl::OkStatus();
  }
  return s;
}
}  
}  