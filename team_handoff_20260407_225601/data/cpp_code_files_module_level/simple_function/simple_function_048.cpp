#ifndef TENSORFLOW_EXAMPLES_SPEECH_COMMANDS_RECOGNIZE_COMMANDS_H_
#define TENSORFLOW_EXAMPLES_SPEECH_COMMANDS_RECOGNIZE_COMMANDS_H_
#include <deque>
#include <unordered_set>
#include <vector>
#include "tensorflow/core/framework/tensor.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/types.h"
namespace tensorflow {
class RecognizeCommands {
 public:
  explicit RecognizeCommands(const std::vector<string>& labels,
                             int32_t average_window_duration_ms = 1000,
                             float detection_threshold = 0.2,
                             int32_t suppression_ms = 500,
                             int32_t minimum_count = 3);
  Status ProcessLatestResults(const Tensor& latest_results,
                              const int64_t current_time_ms,
                              string* found_command, float* score,
                              bool* is_new_command);
 private:
  std::vector<string> labels_;
  int32 average_window_duration_ms_;
  float detection_threshold_;
  int32 suppression_ms_;
  int32 minimum_count_;
  std::deque<std::pair<int64_t, Tensor>> previous_results_;
  string previous_top_label_;
  int64_t labels_count_;
  int64_t previous_top_label_time_;
};
}  
#endif  
#include "tensorflow/examples/speech_commands/recognize_commands.h"
#include "absl/status/status.h"
#include "tensorflow/core/framework/tensor.h"
#include "tensorflow/core/platform/errors.h"
#include "tensorflow/core/platform/status.h"
#include "tensorflow/core/platform/types.h"
namespace tensorflow {
RecognizeCommands::RecognizeCommands(const std::vector<string>& labels,
                                     int32_t average_window_duration_ms,
                                     float detection_threshold,
                                     int32_t suppression_ms,
                                     int32_t minimum_count)
    : labels_(labels),
      average_window_duration_ms_(average_window_duration_ms),
      detection_threshold_(detection_threshold),
      suppression_ms_(suppression_ms),
      minimum_count_(minimum_count) {
  labels_count_ = labels.size();
  previous_top_label_ = "_silence_";
  previous_top_label_time_ = std::numeric_limits<int64_t>::min();
}
Status RecognizeCommands::ProcessLatestResults(const Tensor& latest_results,
                                               const int64_t current_time_ms,
                                               string* found_command,
                                               float* score,
                                               bool* is_new_command) {
  if (latest_results.NumElements() != labels_count_) {
    return errors::InvalidArgument(
        "The results for recognition should contain ", labels_count_,
        " elements, but there are ", latest_results.NumElements());
  }
  if ((!previous_results_.empty()) &&
      (current_time_ms < previous_results_.front().first)) {
    return errors::InvalidArgument(
        "Results must be fed in increasing time order, but received a "
        "timestamp of ",
        current_time_ms, " that was earlier than the previous one of ",
        previous_results_.front().first);
  }
  previous_results_.push_back({current_time_ms, latest_results});
  const int64_t time_limit = current_time_ms - average_window_duration_ms_;
  while (previous_results_.front().first < time_limit) {
    previous_results_.pop_front();
  }
  const int64_t how_many_results = previous_results_.size();
  const int64_t earliest_time = previous_results_.front().first;
  const int64_t samples_duration = current_time_ms - earliest_time;
  if ((how_many_results < minimum_count_) ||
      (samples_duration < (average_window_duration_ms_ / 4))) {
    *found_command = previous_top_label_;
    *score = 0.0f;
    *is_new_command = false;
    return absl::OkStatus();
  }
  std::vector<float> average_scores(labels_count_);
  for (const auto& previous_result : previous_results_) {
    const Tensor& scores_tensor = previous_result.second;
    auto scores_flat = scores_tensor.flat<float>();
    for (int i = 0; i < scores_flat.size(); ++i) {
      average_scores[i] += scores_flat(i) / how_many_results;
    }
  }
  std::vector<std::pair<int, float>> sorted_average_scores;
  sorted_average_scores.reserve(labels_count_);
  for (int i = 0; i < labels_count_; ++i) {
    sorted_average_scores.push_back(
        std::pair<int, float>({i, average_scores[i]}));
  }
  std::sort(sorted_average_scores.begin(), sorted_average_scores.end(),
            [](const std::pair<int, float>& left,
               const std::pair<int, float>& right) {
              return left.second > right.second;
            });
  const int current_top_index = sorted_average_scores[0].first;
  const string current_top_label = labels_[current_top_index];
  const float current_top_score = sorted_average_scores[0].second;
  int64_t time_since_last_top;
  if ((previous_top_label_ == "_silence_") ||
      (previous_top_label_time_ == std::numeric_limits<int64_t>::min())) {
    time_since_last_top = std::numeric_limits<int64_t>::max();
  } else {
    time_since_last_top = current_time_ms - previous_top_label_time_;
  }
  if ((current_top_score > detection_threshold_) &&
      (current_top_label != previous_top_label_) &&
      (time_since_last_top > suppression_ms_)) {
    previous_top_label_ = current_top_label;
    previous_top_label_time_ = current_time_ms;
    *is_new_command = true;
  } else {
    *is_new_command = false;
  }
  *found_command = current_top_label;
  *score = current_top_score;
  return absl::OkStatus();
}
}  