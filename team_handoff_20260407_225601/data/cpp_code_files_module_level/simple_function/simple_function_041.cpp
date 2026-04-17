#ifndef TENSORFLOW_CORE_KERNELS_SPECTROGRAM_H_
#define TENSORFLOW_CORE_KERNELS_SPECTROGRAM_H_
#include <complex>
#include <deque>
#include <vector>
#include "third_party/fft2d/fft.h"
#include "tensorflow/core/framework/op_kernel.h"
#include "tensorflow/core/framework/tensor.h"
namespace tensorflow {
class Spectrogram {
 public:
  Spectrogram() : initialized_(false) {}
  ~Spectrogram() {}
  bool Initialize(int window_length, int step_length);
  bool Initialize(const std::vector<double>& window, int step_length);
  bool Reset();
  template <class InputSample, class OutputSample>
  bool ComputeComplexSpectrogram(
      const std::vector<InputSample>& input,
      std::vector<std::vector<std::complex<OutputSample>>>* output);
  template <class InputSample, class OutputSample>
  bool ComputeSquaredMagnitudeSpectrogram(
      const std::vector<InputSample>& input,
      std::vector<std::vector<OutputSample>>* output);
  const std::vector<double>& GetWindow() const { return window_; }
  int output_frequency_channels() const { return output_frequency_channels_; }
 private:
  template <class InputSample>
  bool GetNextWindowOfSamples(const std::vector<InputSample>& input,
                              int* input_start);
  void ProcessCoreFFT();
  int fft_length_;
  int output_frequency_channels_;
  int window_length_;
  int step_length_;
  bool initialized_;
  int samples_to_next_step_;
  std::vector<double> window_;
  std::vector<double> fft_input_output_;
  std::deque<double> input_queue_;
  std::vector<int> fft_integer_working_area_;
  std::vector<double> fft_double_working_area_;
  Spectrogram(const Spectrogram&) = delete;
  void operator=(const Spectrogram&) = delete;
};
}  
#endif  
#include "tensorflow/core/kernels/spectrogram.h"
#include <math.h>
#include "third_party/fft2d/fft.h"
#include "tensorflow/core/lib/core/bits.h"
namespace tensorflow {
using std::complex;
namespace {
void GetPeriodicHann(int window_length, std::vector<double>* window) {
  const double pi = std::atan(1) * 4;
  window->resize(window_length);
  for (int i = 0; i < window_length; ++i) {
    (*window)[i] = 0.5 - 0.5 * cos((2 * pi * i) / window_length);
  }
}
}  
bool Spectrogram::Initialize(int window_length, int step_length) {
  std::vector<double> window;
  GetPeriodicHann(window_length, &window);
  return Initialize(window, step_length);
}
bool Spectrogram::Initialize(const std::vector<double>& window,
                             int step_length) {
  window_length_ = window.size();
  window_ = window;  
  if (window_length_ < 2) {
    LOG(ERROR) << "Window length too short.";
    initialized_ = false;
    return false;
  }
  step_length_ = step_length;
  if (step_length_ < 1) {
    LOG(ERROR) << "Step length must be positive.";
    initialized_ = false;
    return false;
  }
  fft_length_ = NextPowerOfTwo(window_length_);
  CHECK(fft_length_ >= window_length_);
  output_frequency_channels_ = 1 + fft_length_ / 2;
  fft_input_output_.resize(fft_length_ + 2);
  int half_fft_length = fft_length_ / 2;
  fft_double_working_area_.resize(half_fft_length);
  fft_integer_working_area_.resize(2 + static_cast<int>(sqrt(half_fft_length)));
  initialized_ = true;
  if (!Reset()) {
    LOG(ERROR) << "Failed to Reset()";
    return false;
  }
  return true;
}
bool Spectrogram::Reset() {
  if (!initialized_) {
    LOG(ERROR) << "Initialize() has to be called, before Reset().";
    return false;
  }
  std::fill(fft_double_working_area_.begin(), fft_double_working_area_.end(),
            0.0);
  std::fill(fft_integer_working_area_.begin(), fft_integer_working_area_.end(),
            0);
  fft_integer_working_area_[0] = 0;
  input_queue_.clear();
  samples_to_next_step_ = window_length_;
  return true;
}
template <class InputSample, class OutputSample>
bool Spectrogram::ComputeComplexSpectrogram(
    const std::vector<InputSample>& input,
    std::vector<std::vector<complex<OutputSample>>>* output) {
  if (!initialized_) {
    LOG(ERROR) << "ComputeComplexSpectrogram() called before successful call "
               << "to Initialize().";
    return false;
  }
  CHECK(output);
  output->clear();
  int input_start = 0;
  while (GetNextWindowOfSamples(input, &input_start)) {
    DCHECK_EQ(input_queue_.size(), window_length_);
    ProcessCoreFFT();  
    output->resize(output->size() + 1);
    auto& spectrogram_slice = output->back();
    spectrogram_slice.resize(output_frequency_channels_);
    for (int i = 0; i < output_frequency_channels_; ++i) {
      spectrogram_slice[i] = complex<OutputSample>(
          fft_input_output_[2 * i], fft_input_output_[2 * i + 1]);
    }
  }
  return true;
}
template bool Spectrogram::ComputeComplexSpectrogram(
    const std::vector<float>& input, std::vector<std::vector<complex<float>>>*);
template bool Spectrogram::ComputeComplexSpectrogram(
    const std::vector<double>& input,
    std::vector<std::vector<complex<float>>>*);
template bool Spectrogram::ComputeComplexSpectrogram(
    const std::vector<float>& input,
    std::vector<std::vector<complex<double>>>*);
template bool Spectrogram::ComputeComplexSpectrogram(
    const std::vector<double>& input,
    std::vector<std::vector<complex<double>>>*);
template <class InputSample, class OutputSample>
bool Spectrogram::ComputeSquaredMagnitudeSpectrogram(
    const std::vector<InputSample>& input,
    std::vector<std::vector<OutputSample>>* output) {
  if (!initialized_) {
    LOG(ERROR) << "ComputeSquaredMagnitudeSpectrogram() called before "
               << "successful call to Initialize().";
    return false;
  }
  CHECK(output);
  output->clear();
  int input_start = 0;
  while (GetNextWindowOfSamples(input, &input_start)) {
    DCHECK_EQ(input_queue_.size(), window_length_);
    ProcessCoreFFT();  
    output->resize(output->size() + 1);
    auto& spectrogram_slice = output->back();
    spectrogram_slice.resize(output_frequency_channels_);
    for (int i = 0; i < output_frequency_channels_; ++i) {
      const double re = fft_input_output_[2 * i];
      const double im = fft_input_output_[2 * i + 1];
      spectrogram_slice[i] = re * re + im * im;
    }
  }
  return true;
}
template bool Spectrogram::ComputeSquaredMagnitudeSpectrogram(
    const std::vector<float>& input, std::vector<std::vector<float>>*);
template bool Spectrogram::ComputeSquaredMagnitudeSpectrogram(
    const std::vector<double>& input, std::vector<std::vector<float>>*);
template bool Spectrogram::ComputeSquaredMagnitudeSpectrogram(
    const std::vector<float>& input, std::vector<std::vector<double>>*);
template bool Spectrogram::ComputeSquaredMagnitudeSpectrogram(
    const std::vector<double>& input, std::vector<std::vector<double>>*);
template <class InputSample>
bool Spectrogram::GetNextWindowOfSamples(const std::vector<InputSample>& input,
                                         int* input_start) {
  auto input_it = input.begin() + *input_start;
  int input_remaining = input.end() - input_it;
  if (samples_to_next_step_ > input_remaining) {
    input_queue_.insert(input_queue_.end(), input_it, input.end());
    *input_start += input_remaining;  
    samples_to_next_step_ -= input_remaining;
    return false;  
  } else {
    input_queue_.insert(input_queue_.end(), input_it,
                        input_it + samples_to_next_step_);
    *input_start += samples_to_next_step_;
    input_queue_.erase(
        input_queue_.begin(),
        input_queue_.begin() + input_queue_.size() - window_length_);
    DCHECK_EQ(window_length_, input_queue_.size());
    samples_to_next_step_ = step_length_;  
    return true;  
  }
}
void Spectrogram::ProcessCoreFFT() {
  for (int j = 0; j < window_length_; ++j) {
    fft_input_output_[j] = input_queue_[j] * window_[j];
  }
  for (int j = window_length_; j < fft_length_; ++j) {
    fft_input_output_[j] = 0.0;
  }
  const int kForwardFFT = 1;  
  rdft(fft_length_, kForwardFFT, &fft_input_output_[0],
       &fft_integer_working_area_[0], &fft_double_working_area_[0]);
  fft_input_output_[fft_length_] = fft_input_output_[1];
  fft_input_output_[fft_length_ + 1] = 0;
  fft_input_output_[1] = 0;
}
}  