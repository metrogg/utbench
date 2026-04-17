#ifndef TENSORFLOW_LITE_TOOLS_EVALUATION_STAGES_INFERENCE_PROFILER_STAGE_H_
#define TENSORFLOW_LITE_TOOLS_EVALUATION_STAGES_INFERENCE_PROFILER_STAGE_H_
#include <stdint.h>
#include <memory>
#include <string>
#include <vector>
#include "xla/tsl/util/stats_calculator.h"
#include "tensorflow/lite/c/c_api_types.h"
#include "tensorflow/lite/tools/evaluation/evaluation_delegate_provider.h"
#include "tensorflow/lite/tools/evaluation/evaluation_stage.h"
#include "tensorflow/lite/tools/evaluation/proto/evaluation_config.pb.h"
#include "tensorflow/lite/tools/evaluation/stages/tflite_inference_stage.h"
namespace tflite {
namespace evaluation {
class InferenceProfilerStage : public EvaluationStage {
 public:
  explicit InferenceProfilerStage(const EvaluationStageConfig& config)
      : EvaluationStage(config) {}
  TfLiteStatus Init() override { return Init(nullptr); }
  TfLiteStatus Init(const DelegateProviders* delegate_providers);
  TfLiteStatus Run() override;
  EvaluationStageMetrics LatestMetrics() override;
 private:
  std::unique_ptr<TfliteInferenceStage> reference_stage_;
  std::unique_ptr<TfliteInferenceStage> test_stage_;
  const TfLiteModelInfo* model_info_;
  std::vector<int64_t> input_num_elements_;
  std::vector<int64_t> output_num_elements_;
  std::vector<tsl::Stat<float>> error_stats_;
  std::vector<std::vector<float>> float_tensors_;
  std::vector<std::vector<int8_t>> int8_tensors_;
  std::vector<std::vector<uint8_t>> uint8_tensors_;
  std::vector<std::vector<uint16_t>> float16_tensors_;
  std::vector<std::vector<int64_t>> int64_tensors_;
};
}  
}  
#endif  
#include "tensorflow/lite/tools/evaluation/stages/inference_profiler_stage.h"
#include <cmath>
#include <limits>
#include <memory>
#include <random>
#include "fp16.h"  
#include "tensorflow/core/platform/logging.h"
#include "tensorflow/lite/c/c_api_types.h"
#include "tensorflow/lite/tools/evaluation/evaluation_delegate_provider.h"
#include "tensorflow/lite/tools/evaluation/proto/evaluation_config.pb.h"
#include "tensorflow/lite/tools/evaluation/proto/evaluation_stages.pb.h"
#include "tensorflow/lite/tools/evaluation/stages/tflite_inference_stage.h"
namespace tflite {
namespace evaluation {
namespace {
constexpr float kGaussianFloatMean = 0.5;
constexpr float kGaussianStdDev = 1.0 / 3;
template <typename T>
void GenerateRandomGaussianData(int64_t num_elements, float min, float max,
                                std::vector<T>* data) {
  data->clear();
  data->reserve(num_elements);
  static std::normal_distribution<double> distribution(kGaussianFloatMean,
                                                       kGaussianStdDev);
  static std::default_random_engine generator;
  for (int i = 0; i < num_elements; ++i) {
    auto rand_n = distribution(generator);
    while (rand_n < 0 || rand_n >= 1) {
      rand_n = distribution(generator);
    }
    auto rand_float = min + (max - min) * static_cast<float>(rand_n);
    data->push_back(static_cast<T>(rand_float));
  }
}
template <typename T>
float CalculateAverageError(T* reference, T* test, int64_t num_elements) {
  float error = 0;
  for (int i = 0; i < num_elements; i++) {
    float test_value = static_cast<float>(test[i]);
    float reference_value = static_cast<float>(reference[i]);
    error += std::abs(test_value - reference_value);
  }
  error /= num_elements;
  return error;
}
}  
TfLiteStatus InferenceProfilerStage::Init(
    const DelegateProviders* delegate_providers) {
  test_stage_ = std::make_unique<TfliteInferenceStage>(config_);
  if (test_stage_->Init(delegate_providers) != kTfLiteOk) return kTfLiteError;
  LOG(INFO) << "Test interpreter has been initialized.";
  EvaluationStageConfig reference_config;
  reference_config.set_name("reference_inference");
  auto* params = reference_config.mutable_specification()
                     ->mutable_tflite_inference_params();
  params->set_model_file_path(
      config_.specification().tflite_inference_params().model_file_path());
  params->set_invocations_per_run(
      config_.specification().tflite_inference_params().invocations_per_run());
  reference_stage_ = std::make_unique<TfliteInferenceStage>(reference_config);
  if (reference_stage_->Init() != kTfLiteOk) return kTfLiteError;
  LOG(INFO) << "Reference interpreter (1 thread on CPU) has been initialized.";
  model_info_ = reference_stage_->GetModelInfo();
  for (int i = 0; i < model_info_->inputs.size(); ++i) {
    const TfLiteType model_input_type = model_info_->inputs[i]->type;
    if (model_input_type == kTfLiteUInt8 || model_input_type == kTfLiteInt8 ||
        model_input_type == kTfLiteInt64 ||
        model_input_type == kTfLiteFloat32 ||
        model_input_type == kTfLiteFloat16) {
    } else {
      LOG(ERROR) << "InferenceProfilerStage only supports "
                    "float16/float32/int8/uint8/int64 "
                    "input types";
      return kTfLiteError;
    }
    auto* input_shape = model_info_->inputs[i]->dims;
    int64_t total_num_elements = 1;
    for (int i = 0; i < input_shape->size; i++) {
      total_num_elements *= input_shape->data[i];
    }
    input_num_elements_.push_back(total_num_elements);
    float_tensors_.emplace_back();
    uint8_tensors_.emplace_back();
    int8_tensors_.emplace_back();
    float16_tensors_.emplace_back();
    int64_tensors_.emplace_back();
  }
  for (int i = 0; i < model_info_->outputs.size(); ++i) {
    const TfLiteType model_output_type = model_info_->outputs[i]->type;
    if (model_output_type == kTfLiteUInt8 || model_output_type == kTfLiteInt8 ||
        model_output_type == kTfLiteFloat32) {
    } else {
      LOG(ERROR) << "InferenceProfilerStage only supports float32/int8/uint8 "
                    "output types";
      return kTfLiteError;
    }
    auto* output_shape = model_info_->outputs[i]->dims;
    int64_t total_num_elements = 1;
    for (int i = 0; i < output_shape->size; i++) {
      total_num_elements *= output_shape->data[i];
    }
    output_num_elements_.push_back(total_num_elements);
    error_stats_.emplace_back();
  }
  return kTfLiteOk;
}
TfLiteStatus InferenceProfilerStage::Run() {
  std::vector<void*> input_ptrs;
  for (int i = 0; i < model_info_->inputs.size(); ++i) {
    const TfLiteType model_input_type = model_info_->inputs[i]->type;
    if (model_input_type == kTfLiteUInt8) {
      GenerateRandomGaussianData(
          input_num_elements_[i], std::numeric_limits<uint8_t>::min(),
          std::numeric_limits<uint8_t>::max(), &uint8_tensors_[i]);
      input_ptrs.push_back(uint8_tensors_[i].data());
    } else if (model_input_type == kTfLiteInt8) {
      GenerateRandomGaussianData(
          input_num_elements_[i], std::numeric_limits<int8_t>::min(),
          std::numeric_limits<int8_t>::max(), &int8_tensors_[i]);
      input_ptrs.push_back(int8_tensors_[i].data());
    } else if (model_input_type == kTfLiteInt64) {
      GenerateRandomGaussianData(
          input_num_elements_[i], std::numeric_limits<int64_t>::min(),
          std::numeric_limits<int64_t>::max(), &int64_tensors_[i]);
      input_ptrs.push_back(int64_tensors_[i].data());
    } else if (model_input_type == kTfLiteFloat32) {
      GenerateRandomGaussianData(input_num_elements_[i], -1, 1,
                                 &(float_tensors_[i]));
      input_ptrs.push_back(float_tensors_[i].data());
    } else if (model_input_type == kTfLiteFloat16) {
      GenerateRandomGaussianData(input_num_elements_[i], -1, 1,
                                 &(float_tensors_[i]));
      for (size_t j = 0; j < float_tensors_[i].size(); j++) {
        float16_tensors_[i][j] =
            fp16_ieee_from_fp32_value(float_tensors_[i][j]);
      }
      input_ptrs.push_back(float16_tensors_[i].data());
    } else {
      LOG(ERROR) << "InferenceProfilerStage only supports "
                    "float16/float32/int8/uint8/int64 "
                    "input types";
      return kTfLiteError;
    }
  }
  test_stage_->SetInputs(input_ptrs);
  reference_stage_->SetInputs(input_ptrs);
  if (test_stage_->Run() != kTfLiteOk) return kTfLiteError;
  if (reference_stage_->Run() != kTfLiteOk) return kTfLiteError;
  for (int i = 0; i < model_info_->outputs.size(); ++i) {
    const TfLiteType model_output_type = model_info_->outputs[i]->type;
    void* reference_ptr = reference_stage_->GetOutputs()->at(i);
    void* test_ptr = test_stage_->GetOutputs()->at(i);
    float output_diff = 0;
    if (model_output_type == kTfLiteUInt8) {
      output_diff = CalculateAverageError(static_cast<uint8_t*>(reference_ptr),
                                          static_cast<uint8_t*>(test_ptr),
                                          output_num_elements_[i]);
    } else if (model_output_type == kTfLiteInt8) {
      output_diff = CalculateAverageError(static_cast<int8_t*>(reference_ptr),
                                          static_cast<int8_t*>(test_ptr),
                                          output_num_elements_[i]);
    } else if (model_output_type == kTfLiteFloat32) {
      output_diff = CalculateAverageError(static_cast<float*>(reference_ptr),
                                          static_cast<float*>(test_ptr),
                                          output_num_elements_[i]);
    }
    error_stats_[i].UpdateStat(output_diff);
  }
  return kTfLiteOk;
}
EvaluationStageMetrics InferenceProfilerStage::LatestMetrics() {
  EvaluationStageMetrics metrics;
  const auto& reference_metrics = reference_stage_->LatestMetrics();
  metrics.set_num_runs(reference_metrics.num_runs());
  auto* inference_profiler_metrics =
      metrics.mutable_process_metrics()->mutable_inference_profiler_metrics();
  *inference_profiler_metrics->mutable_reference_latency() =
      reference_metrics.process_metrics().total_latency();
  *inference_profiler_metrics->mutable_test_latency() =
      test_stage_->LatestMetrics().process_metrics().total_latency();
  for (int i = 0; i < error_stats_.size(); ++i) {
    AccuracyMetrics* diff = inference_profiler_metrics->add_output_errors();
    diff->set_avg_value(error_stats_[i].avg());
    diff->set_std_deviation(error_stats_[i].std_deviation());
    diff->set_min_value(error_stats_[i].min());
    if (error_stats_[i].avg() != 0) {
      diff->set_max_value(error_stats_[i].max());
    } else {
      diff->set_max_value(0);
    }
  }
  return metrics;
}
}  
}  