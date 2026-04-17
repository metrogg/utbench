#ifndef TENSORFLOW_LITE_TOOLS_EVALUATION_STAGES_UTILS_IMAGE_METRICS_H_
#define TENSORFLOW_LITE_TOOLS_EVALUATION_STAGES_UTILS_IMAGE_METRICS_H_
#include <stdint.h>
#include <vector>
namespace tflite {
namespace evaluation {
namespace image {
struct Box2D {
  struct Interval {
    float min = 0;
    float max = 0;
    Interval(float x, float y) {
      min = x;
      max = y;
    }
    Interval() {}
  };
  Interval x;
  Interval y;
  static float Length(const Interval& a);
  static float Intersection(const Interval& a, const Interval& b);
  float Area() const;
  float Intersection(const Box2D& other) const;
  float Union(const Box2D& other) const;
  float IoU(const Box2D& other) const;
  float Overlap(const Box2D& other) const;
};
enum IgnoreType {
  kDontIgnore = 0,
  kIgnoreOneMatch = 1,
  kIgnoreAllMatches = 2,
};
struct Detection {
 public:
  bool difficult = false;
  int64_t imgid = 0;
  float score = 0;
  Box2D box;
  IgnoreType ignore = IgnoreType::kDontIgnore;
  Detection() {}
  Detection(bool d, int64_t id, float s, Box2D b)
      : difficult(d), imgid(id), score(s), box(b) {}
  Detection(bool d, int64_t id, float s, Box2D b, IgnoreType i)
      : difficult(d), imgid(id), score(s), box(b), ignore(i) {}
};
struct PR {
  float p = 0;
  float r = 0;
  PR(const float p_, const float r_) : p(p_), r(r_) {}
};
class AveragePrecision {
 public:
  struct Options {
    float iou_threshold = 0.5;
    int num_recall_points = 100;
  };
  AveragePrecision() : AveragePrecision(Options()) {}
  explicit AveragePrecision(const Options& opts) : opts_(opts) {}
  float FromPRCurve(const std::vector<PR>& pr,
                    std::vector<PR>* pr_out = nullptr);
  float FromBoxes(const std::vector<Detection>& groundtruth,
                  const std::vector<Detection>& prediction,
                  std::vector<PR>* pr_out = nullptr);
 private:
  Options opts_;
};
}  
}  
}  
#endif  
#include "tensorflow/lite/tools/evaluation/stages/utils/image_metrics.h"
#include <algorithm>
#include <cmath>
#include "absl/container/flat_hash_map.h"
#include "tensorflow/core/platform/logging.h"
namespace tflite {
namespace evaluation {
namespace image {
float Box2D::Length(const Box2D::Interval& a) {
  return std::max(0.f, a.max - a.min);
}
float Box2D::Intersection(const Box2D::Interval& a, const Box2D::Interval& b) {
  return Length(Interval{std::max(a.min, b.min), std::min(a.max, b.max)});
}
float Box2D::Area() const { return Length(x) * Length(y); }
float Box2D::Intersection(const Box2D& other) const {
  return Intersection(x, other.x) * Intersection(y, other.y);
}
float Box2D::Union(const Box2D& other) const {
  return Area() + other.Area() - Intersection(other);
}
float Box2D::IoU(const Box2D& other) const {
  const float total = Union(other);
  if (total > 0) {
    return Intersection(other) / total;
  } else {
    return 0.0;
  }
}
float Box2D::Overlap(const Box2D& other) const {
  const float intersection = Intersection(other);
  return intersection > 0 ? intersection / Area() : 0.0;
}
float AveragePrecision::FromPRCurve(const std::vector<PR>& pr,
                                    std::vector<PR>* pr_out) {
  float p = 0;
  float sum = 0;
  int r_level = opts_.num_recall_points;
  for (int i = pr.size() - 1; i >= 0; --i) {
    const PR& item = pr[i];
    if (i > 0) {
      if (item.r < pr[i - 1].r) {
        LOG(ERROR) << "recall points are not in order: " << pr[i - 1].r << ", "
                   << item.r;
        return 0;
      }
    }
    while (item.r * opts_.num_recall_points < r_level) {
      const float recall =
          static_cast<float>(r_level) / opts_.num_recall_points;
      if (r_level < 0) {
        LOG(ERROR) << "Number of recall points should be > 0";
        return 0;
      }
      sum += p;
      r_level -= 1;
      if (pr_out != nullptr) {
        pr_out->emplace_back(p, recall);
      }
    }
    p = std::max(p, item.p);
  }
  for (; r_level >= 0; --r_level) {
    const float recall = static_cast<float>(r_level) / opts_.num_recall_points;
    sum += p;
    if (pr_out != nullptr) {
      pr_out->emplace_back(p, recall);
    }
  }
  return sum / (1 + opts_.num_recall_points);
}
float AveragePrecision::FromBoxes(const std::vector<Detection>& groundtruth,
                                  const std::vector<Detection>& prediction,
                                  std::vector<PR>* pr_out) {
  absl::flat_hash_map<int64_t, std::list<Detection>> gt;
  int num_gt = 0;
  for (auto& box : groundtruth) {
    gt[box.imgid].push_back(box);
    if (!box.difficult && box.ignore == kDontIgnore) {
      ++num_gt;
    }
  }
  if (num_gt == 0) {
    return NAN;
  }
  std::vector<Detection> pd = prediction;
  std::sort(pd.begin(), pd.end(), [](const Detection& a, const Detection& b) {
    return a.score > b.score;
  });
  std::vector<PR> pr;
  int correct = 0;
  int num_pd = 0;
  for (int i = 0; i < pd.size(); ++i) {
    const Detection& b = pd[i];
    auto* g = &gt[b.imgid];
    auto best = g->end();
    float best_iou = -INFINITY;
    for (auto it = g->begin(); it != g->end(); ++it) {
      const auto iou = b.box.IoU(it->box);
      if (iou > best_iou) {
        best = it;
        best_iou = iou;
      }
    }
    if ((best != g->end()) && (best_iou >= opts_.iou_threshold)) {
      if (best->difficult) {
        continue;
      }
      switch (best->ignore) {
        case kDontIgnore: {
          ++correct;
          ++num_pd;
          g->erase(best);
          pr.push_back({static_cast<float>(correct) / num_pd,
                        static_cast<float>(correct) / num_gt});
          break;
        }
        case kIgnoreOneMatch: {
          g->erase(best);
          break;
        }
        case kIgnoreAllMatches: {
          break;
        }
      }
    } else {
      ++num_pd;
      pr.push_back({static_cast<float>(correct) / num_pd,
                    static_cast<float>(correct) / num_gt});
    }
  }
  return FromPRCurve(pr, pr_out);
}
}  
}  
}  