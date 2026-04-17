#ifndef TENSORFLOW_CORE_KERNELS_MFCC_MEL_FILTERBANK_H_
#define TENSORFLOW_CORE_KERNELS_MFCC_MEL_FILTERBANK_H_
#include <vector>
#include "tensorflow/core/framework/op_kernel.h"
namespace tensorflow {
class MfccMelFilterbank {
 public:
  MfccMelFilterbank();
  bool Initialize(int input_length,  
                  double input_sample_rate, int output_channel_count,
                  double lower_frequency_limit, double upper_frequency_limit);
  void Compute(const std::vector<double>& input,
               std::vector<double>* output) const;
 private:
  double FreqToMel(double freq) const;
  bool initialized_;
  int num_channels_;
  double sample_rate_;
  int input_length_;
  std::vector<double> center_frequencies_;  
  std::vector<double> weights_;  
  std::vector<int> band_mapper_;
  int start_index_;  
  int end_index_;    
  MfccMelFilterbank(const MfccMelFilterbank&) = delete;
  void operator=(const MfccMelFilterbank&) = delete;
};
}  
#endif  
#include "tensorflow/core/kernels/mfcc_mel_filterbank.h"
#include <math.h>
#include <limits>
#include "tensorflow/core/platform/logging.h"
namespace tensorflow {
MfccMelFilterbank::MfccMelFilterbank() : initialized_(false) {}
bool MfccMelFilterbank::Initialize(int input_length, double input_sample_rate,
                                   int output_channel_count,
                                   double lower_frequency_limit,
                                   double upper_frequency_limit) {
  num_channels_ = output_channel_count;
  sample_rate_ = input_sample_rate;
  input_length_ = input_length;
  if (num_channels_ < 1) {
    LOG(ERROR) << "Number of filterbank channels must be positive.";
    return false;
  }
  if (sample_rate_ <= 0) {
    LOG(ERROR) << "Sample rate must be positive.";
    return false;
  }
  if (input_length < 2) {
    LOG(ERROR) << "Input length must greater than 1.";
    return false;
  }
  if (lower_frequency_limit < 0) {
    LOG(ERROR) << "Lower frequency limit must be nonnegative.";
    return false;
  }
  if (upper_frequency_limit <= lower_frequency_limit) {
    LOG(ERROR) << "Upper frequency limit must be greater than "
               << "lower frequency limit.";
    return false;
  }
  std::size_t center_frequencies_size = std::size_t(num_channels_) + 1;
  if (center_frequencies_size >= std::numeric_limits<int>::max() ||
      center_frequencies_size > center_frequencies_.max_size()) {
    LOG(ERROR) << "Number of filterbank channels must be less than "
               << std::numeric_limits<int>::max()
               << " and less than or equal to "
               << center_frequencies_.max_size();
    return false;
  }
  center_frequencies_.resize(center_frequencies_size);
  const double mel_low = FreqToMel(lower_frequency_limit);
  const double mel_hi = FreqToMel(upper_frequency_limit);
  const double mel_span = mel_hi - mel_low;
  const double mel_spacing = mel_span / static_cast<double>(num_channels_ + 1);
  for (int i = 0; i < num_channels_ + 1; ++i) {
    center_frequencies_[i] = mel_low + (mel_spacing * (i + 1));
  }
  const double hz_per_sbin =
      0.5 * sample_rate_ / static_cast<double>(input_length_ - 1);
  start_index_ = static_cast<int>(1.5 + (lower_frequency_limit / hz_per_sbin));
  end_index_ = static_cast<int>(upper_frequency_limit / hz_per_sbin);
  band_mapper_.resize(input_length_);
  int channel = 0;
  for (int i = 0; i < input_length_; ++i) {
    double melf = FreqToMel(i * hz_per_sbin);
    if ((i < start_index_) || (i > end_index_)) {
      band_mapper_[i] = -2;  
    } else {
      while ((channel < num_channels_) &&
             (center_frequencies_[channel] < melf)) {
        ++channel;
      }
      band_mapper_[i] = channel - 1;  
    }
  }
  weights_.resize(input_length_);
  for (int i = 0; i < input_length_; ++i) {
    channel = band_mapper_[i];
    if ((i < start_index_) || (i > end_index_)) {
      weights_[i] = 0.0;
    } else {
      if (channel >= 0) {
        weights_[i] =
            (center_frequencies_[channel + 1] - FreqToMel(i * hz_per_sbin)) /
            (center_frequencies_[channel + 1] - center_frequencies_[channel]);
      } else {
        weights_[i] = (center_frequencies_[0] - FreqToMel(i * hz_per_sbin)) /
                      (center_frequencies_[0] - mel_low);
      }
    }
  }
  std::vector<int> bad_channels;
  for (int c = 0; c < num_channels_; ++c) {
    float band_weights_sum = 0.0;
    for (int i = 0; i < input_length_; ++i) {
      if (band_mapper_[i] == c - 1) {
        band_weights_sum += (1.0 - weights_[i]);
      } else if (band_mapper_[i] == c) {
        band_weights_sum += weights_[i];
      }
    }
    if (band_weights_sum < 0.5) {
      bad_channels.push_back(c);
    }
  }
  if (!bad_channels.empty()) {
    LOG(ERROR) << "Missing " << bad_channels.size() << " bands "
               << " starting at " << bad_channels[0]
               << " in mel-frequency design. "
               << "Perhaps too many channels or "
               << "not enough frequency resolution in spectrum. ("
               << "input_length: " << input_length
               << " input_sample_rate: " << input_sample_rate
               << " output_channel_count: " << output_channel_count
               << " lower_frequency_limit: " << lower_frequency_limit
               << " upper_frequency_limit: " << upper_frequency_limit;
  }
  initialized_ = true;
  return true;
}
void MfccMelFilterbank::Compute(const std::vector<double> &input,
                                std::vector<double> *output) const {
  if (!initialized_) {
    LOG(ERROR) << "Mel Filterbank not initialized.";
    return;
  }
  if (input.size() <= end_index_) {
    LOG(ERROR) << "Input too short to compute filterbank";
    return;
  }
  output->assign(num_channels_, 0.0);
  for (int i = start_index_; i <= end_index_; i++) {  
    double spec_val = sqrt(input[i]);
    double weighted = spec_val * weights_[i];
    int channel = band_mapper_[i];
    if (channel >= 0)
      (*output)[channel] += weighted;  
    channel++;
    if (channel < num_channels_)
      (*output)[channel] += spec_val - weighted;  
  }
}
double MfccMelFilterbank::FreqToMel(double freq) const {
  return 1127.0 * log1p(freq / 700.0);
}
}  