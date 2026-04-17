#ifndef ABSL_STRINGS_INTERNAL_CORD_REP_BTREE_H_
#define ABSL_STRINGS_INTERNAL_CORD_REP_BTREE_H_
#include <cassert>
#include <cstdint>
#include <iosfwd>
#include "absl/base/config.h"
#include "absl/base/internal/raw_logging.h"
#include "absl/base/optimization.h"
#include "absl/strings/internal/cord_data_edge.h"
#include "absl/strings/internal/cord_internal.h"
#include "absl/strings/internal/cord_rep_flat.h"
#include "absl/strings/string_view.h"
#include "absl/types/span.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
namespace cord_internal {
void SetCordBtreeExhaustiveValidation(bool do_exaustive_validation);
bool IsCordBtreeExhaustiveValidationEnabled();
class CordRepBtreeNavigator;
class CordRepBtree : public CordRep {
 public:
  enum class EdgeType { kFront, kBack };
  static constexpr EdgeType kFront = EdgeType::kFront;
  static constexpr EdgeType kBack = EdgeType::kBack;
  static constexpr size_t kMaxCapacity = 6;
  static constexpr size_t kMaxDepth = 12;
  static constexpr int kMaxHeight = static_cast<int>(kMaxDepth - 1);
  enum Action { kSelf, kCopied, kPopped };
  struct OpResult {
    CordRepBtree* tree;
    Action action;
  };
  struct CopyResult {
    CordRep* edge;
    int height;
  };
  struct Position {
    size_t index;
    size_t n;
  };
  static CordRepBtree* Create(CordRep* rep);
  static void Destroy(CordRepBtree* tree);
  static void Delete(CordRepBtree* tree) { delete tree; }
  using CordRep::Unref;
  static void Unref(absl::Span<CordRep* const> edges);
  static CordRepBtree* Append(CordRepBtree* tree, CordRep* rep);
  static CordRepBtree* Prepend(CordRepBtree* tree, CordRep* rep);
  static CordRepBtree* Append(CordRepBtree* tree, string_view data,
                              size_t extra = 0);
  static CordRepBtree* Prepend(CordRepBtree* tree, string_view data,
                               size_t extra = 0);
  CordRep* SubTree(size_t offset, size_t n);
  static CordRep* RemoveSuffix(CordRepBtree* tree, size_t n);
  char GetCharacter(size_t offset) const;
  bool IsFlat(absl::string_view* fragment) const;
  bool IsFlat(size_t offset, size_t n, absl::string_view* fragment) const;
  Span<char> GetAppendBuffer(size_t size);
  static ExtractResult ExtractAppendBuffer(CordRepBtree* tree,
                                           size_t extra_capacity = 1);
  int height() const { return static_cast<int>(storage[0]); }
  size_t begin() const { return static_cast<size_t>(storage[1]); }
  size_t back() const { return static_cast<size_t>(storage[2]) - 1; }
  size_t end() const { return static_cast<size_t>(storage[2]); }
  size_t index(EdgeType edge) const {
    return edge == kFront ? begin() : back();
  }
  size_t size() const { return end() - begin(); }
  size_t capacity() const { return kMaxCapacity; }
  inline CordRep* Edge(size_t index) const;
  inline CordRep* Edge(EdgeType edge_type) const;
  inline absl::Span<CordRep* const> Edges() const;
  inline absl::Span<CordRep* const> Edges(size_t begin, size_t end) const;
  inline absl::string_view Data(size_t index) const;
  static bool IsValid(const CordRepBtree* tree, bool shallow = false);
  static CordRepBtree* AssertValid(CordRepBtree* tree, bool shallow = true);
  static const CordRepBtree* AssertValid(const CordRepBtree* tree,
                                         bool shallow = true);
  static void Dump(const CordRep* rep, std::ostream& stream);
  static void Dump(const CordRep* rep, absl::string_view label,
                   std::ostream& stream);
  static void Dump(const CordRep* rep, absl::string_view label,
                   bool include_contents, std::ostream& stream);
  template <EdgeType edge_type>
  inline OpResult AddEdge(bool owned, CordRep* edge, size_t delta);
  template <EdgeType edge_type>
  OpResult SetEdge(bool owned, CordRep* edge, size_t delta);
  static CordRepBtree* New(int height = 0);
  static CordRepBtree* New(CordRep* rep);
  static CordRepBtree* New(CordRepBtree* front, CordRepBtree* back);
  static CordRepBtree* Rebuild(CordRepBtree* tree);
 private:
  CordRepBtree() = default;
  ~CordRepBtree() = default;
  inline void InitInstance(int height, size_t begin = 0, size_t end = 0);
  void set_begin(size_t begin) { storage[1] = static_cast<uint8_t>(begin); }
  void set_end(size_t end) { storage[2] = static_cast<uint8_t>(end); }
  size_t sub_fetch_begin(size_t n) {
    storage[1] -= static_cast<uint8_t>(n);
    return storage[1];
  }
  size_t fetch_add_end(size_t n) {
    const uint8_t current = storage[2];
    storage[2] = static_cast<uint8_t>(current + n);
    return current;
  }
  Position IndexOf(size_t offset) const;
  Position IndexBefore(size_t offset) const;
  Position IndexOfLength(size_t n) const;
  Position IndexBefore(Position front, size_t offset) const;
  Position IndexBeyond(size_t offset) const;
  template <EdgeType edge_type>
  static CordRepBtree* NewLeaf(absl::string_view data, size_t extra);
  CordRepBtree* CopyRaw(size_t new_length) const;
  CordRepBtree* Copy() const;
  CordRepBtree* CopyBeginTo(size_t end, size_t new_length) const;
  static CordRepBtree* ConsumeBeginTo(CordRepBtree* tree, size_t end,
                                      size_t new_length);
  CordRepBtree* CopyToEndFrom(size_t begin, size_t new_length) const;
  static CordRep* ExtractFront(CordRepBtree* tree);
  static CordRepBtree* MergeTrees(CordRepBtree* left, CordRepBtree* right);
  static CordRepBtree* CreateSlow(CordRep* rep);
  static CordRepBtree* AppendSlow(CordRepBtree*, CordRep* rep);
  static CordRepBtree* PrependSlow(CordRepBtree*, CordRep* rep);
  static void Rebuild(CordRepBtree** stack, CordRepBtree* tree, bool consume);
  inline void AlignBegin();
  inline void AlignEnd();
  template <EdgeType edge_type>
  inline void Add(CordRep* rep);
  template <EdgeType edge_type>
  inline void Add(absl::Span<CordRep* const>);
  template <EdgeType edge_type>
  absl::string_view AddData(absl::string_view data, size_t extra);
  template <EdgeType edge_type>
  inline void SetEdge(CordRep* edge);
  CopyResult CopyPrefix(size_t n, bool allow_folding = true);
  CopyResult CopySuffix(size_t offset);
  inline OpResult ToOpResult(bool owned);
  template <EdgeType edge_type>
  static CordRepBtree* AddCordRep(CordRepBtree* tree, CordRep* rep);
  template <EdgeType edge_type>
  static CordRepBtree* AddData(CordRepBtree* tree, absl::string_view data,
                               size_t extra = 0);
  template <EdgeType edge_type>
  static CordRepBtree* Merge(CordRepBtree* dst, CordRepBtree* src);
  Span<char> GetAppendBufferSlow(size_t size);
  CordRep* edges_[kMaxCapacity];
  friend class CordRepBtreeTestPeer;
  friend class CordRepBtreeNavigator;
};
inline CordRepBtree* CordRep::btree() {
  assert(IsBtree());
  return static_cast<CordRepBtree*>(this);
}
inline const CordRepBtree* CordRep::btree() const {
  assert(IsBtree());
  return static_cast<const CordRepBtree*>(this);
}
inline void CordRepBtree::InitInstance(int height, size_t begin, size_t end) {
  tag = BTREE;
  storage[0] = static_cast<uint8_t>(height);
  storage[1] = static_cast<uint8_t>(begin);
  storage[2] = static_cast<uint8_t>(end);
}
inline CordRep* CordRepBtree::Edge(size_t index) const {
  assert(index >= begin());
  assert(index < end());
  return edges_[index];
}
inline CordRep* CordRepBtree::Edge(EdgeType edge_type) const {
  return edges_[edge_type == kFront ? begin() : back()];
}
inline absl::Span<CordRep* const> CordRepBtree::Edges() const {
  return {edges_ + begin(), size()};
}
inline absl::Span<CordRep* const> CordRepBtree::Edges(size_t begin,
                                                      size_t end) const {
  assert(begin <= end);
  assert(begin >= this->begin());
  assert(end <= this->end());
  return {edges_ + begin, static_cast<size_t>(end - begin)};
}
inline absl::string_view CordRepBtree::Data(size_t index) const {
  assert(height() == 0);
  return EdgeData(Edge(index));
}
inline CordRepBtree* CordRepBtree::New(int height) {
  CordRepBtree* tree = new CordRepBtree;
  tree->length = 0;
  tree->InitInstance(height);
  return tree;
}
inline CordRepBtree* CordRepBtree::New(CordRep* rep) {
  CordRepBtree* tree = new CordRepBtree;
  int height = rep->IsBtree() ? rep->btree()->height() + 1 : 0;
  tree->length = rep->length;
  tree->InitInstance(height, 0, 1);
  tree->edges_[0] = rep;
  return tree;
}
inline CordRepBtree* CordRepBtree::New(CordRepBtree* front,
                                       CordRepBtree* back) {
  assert(front->height() == back->height());
  CordRepBtree* tree = new CordRepBtree;
  tree->length = front->length + back->length;
  tree->InitInstance(front->height() + 1, 0, 2);
  tree->edges_[0] = front;
  tree->edges_[1] = back;
  return tree;
}
inline void CordRepBtree::Unref(absl::Span<CordRep* const> edges) {
  for (CordRep* edge : edges) {
    if (ABSL_PREDICT_FALSE(!edge->refcount.Decrement())) {
      CordRep::Destroy(edge);
    }
  }
}
inline CordRepBtree* CordRepBtree::CopyRaw(size_t new_length) const {
  CordRepBtree* tree = new CordRepBtree;