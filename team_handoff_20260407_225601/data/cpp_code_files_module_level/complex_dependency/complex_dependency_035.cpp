#ifndef ABSL_STRINGS_CORD_H_
#define ABSL_STRINGS_CORD_H_
#include <algorithm>
#include <cstddef>
#include <cstdint>
#include <cstring>
#include <iosfwd>
#include <iterator>
#include <string>
#include <type_traits>
#include "absl/base/attributes.h"
#include "absl/base/config.h"
#include "absl/base/internal/endian.h"
#include "absl/base/internal/per_thread_tls.h"
#include "absl/base/macros.h"
#include "absl/base/nullability.h"
#include "absl/base/optimization.h"
#include "absl/base/port.h"
#include "absl/container/inlined_vector.h"
#include "absl/crc/internal/crc_cord_state.h"
#include "absl/functional/function_ref.h"
#include "absl/meta/type_traits.h"
#include "absl/strings/cord_analysis.h"
#include "absl/strings/cord_buffer.h"
#include "absl/strings/internal/cord_data_edge.h"
#include "absl/strings/internal/cord_internal.h"
#include "absl/strings/internal/cord_rep_btree.h"
#include "absl/strings/internal/cord_rep_btree_reader.h"
#include "absl/strings/internal/cord_rep_crc.h"
#include "absl/strings/internal/cordz_functions.h"
#include "absl/strings/internal/cordz_info.h"
#include "absl/strings/internal/cordz_statistics.h"
#include "absl/strings/internal/cordz_update_scope.h"
#include "absl/strings/internal/cordz_update_tracker.h"
#include "absl/strings/internal/resize_uninitialized.h"
#include "absl/strings/internal/string_constant.h"
#include "absl/strings/string_view.h"
#include "absl/types/optional.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
class Cord;
class CordTestPeer;
template <typename Releaser>
Cord MakeCordFromExternal(absl::string_view, Releaser&&);
void CopyCordToString(const Cord& src, absl::Nonnull<std::string*> dst);
void AppendCordToString(const Cord& src, absl::Nonnull<std::string*> dst);
enum class CordMemoryAccounting {
  kTotal,
  kTotalMorePrecise,
  kFairShare,
};
class Cord {
 private:
  template <typename T>
  using EnableIfString =
      absl::enable_if_t<std::is_same<T, std::string>::value, int>;
 public:
  constexpr Cord() noexcept;
  Cord(const Cord& src);
  Cord(Cord&& src) noexcept;
  Cord& operator=(const Cord& x);
  Cord& operator=(Cord&& x) noexcept;
  explicit Cord(absl::string_view src);
  Cord& operator=(absl::string_view src);
  template <typename T, EnableIfString<T> = 0>
  explicit Cord(T&& src);
  template <typename T, EnableIfString<T> = 0>
  Cord& operator=(T&& src);
  ~Cord() {
    if (contents_.is_tree()) DestroyCordSlow();
  }
  template <typename Releaser>
  friend Cord MakeCordFromExternal(absl::string_view data, Releaser&& releaser);
  ABSL_ATTRIBUTE_REINITIALIZES void Clear();
  void Append(const Cord& src);
  void Append(Cord&& src);
  void Append(absl::string_view src);
  template <typename T, EnableIfString<T> = 0>
  void Append(T&& src);
  void Append(CordBuffer buffer);
  CordBuffer GetAppendBuffer(size_t capacity, size_t min_capacity = 16);
  CordBuffer GetCustomAppendBuffer(size_t block_size, size_t capacity,
                                   size_t min_capacity = 16);
  void Prepend(const Cord& src);
  void Prepend(absl::string_view src);
  template <typename T, EnableIfString<T> = 0>
  void Prepend(T&& src);
  void Prepend(CordBuffer buffer);
  void RemovePrefix(size_t n);
  void RemoveSuffix(size_t n);
  Cord Subcord(size_t pos, size_t new_size) const;
  void swap(Cord& other) noexcept;
  friend void swap(Cord& x, Cord& y) noexcept { x.swap(y); }
  size_t size() const;
  bool empty() const;
  size_t EstimatedMemoryUsage(CordMemoryAccounting accounting_method =
                                  CordMemoryAccounting::kTotal) const;
  int Compare(absl::string_view rhs) const;
  int Compare(const Cord& rhs) const;
  bool StartsWith(const Cord& rhs) const;
  bool StartsWith(absl::string_view rhs) const;
  bool EndsWith(absl::string_view rhs) const;
  bool EndsWith(const Cord& rhs) const;
  bool Contains(absl::string_view rhs) const;
  bool Contains(const Cord& rhs) const;
  explicit operator std::string() const;
  friend void CopyCordToString(const Cord& src,
                               absl::Nonnull<std::string*> dst);
  friend void AppendCordToString(const Cord& src,
                                 absl::Nonnull<std::string*> dst);
  class CharIterator;
  class ChunkIterator {
   public:
    using iterator_category = std::input_iterator_tag;
    using value_type = absl::string_view;
    using difference_type = ptrdiff_t;
    using pointer = absl::Nonnull<const value_type*>;
    using reference = value_type;
    ChunkIterator() = default;
    ChunkIterator& operator++();
    ChunkIterator operator++(int);
    bool operator==(const ChunkIterator& other) const;
    bool operator!=(const ChunkIterator& other) const;
    reference operator*() const;
    pointer operator->() const;
    friend class Cord;
    friend class CharIterator;
   private:
    using CordRep = absl::cord_internal::CordRep;
    using CordRepBtree = absl::cord_internal::CordRepBtree;
    using CordRepBtreeReader = absl::cord_internal::CordRepBtreeReader;
    explicit ChunkIterator(absl::Nonnull<cord_internal::CordRep*> tree);
    explicit ChunkIterator(absl::Nonnull<const Cord*> cord);
    void InitTree(absl::Nonnull<cord_internal::CordRep*> tree);
    void RemoveChunkPrefix(size_t n);
    Cord AdvanceAndReadBytes(size_t n);
    void AdvanceBytes(size_t n);
    ChunkIterator& AdvanceBtree();
    void AdvanceBytesBtree(size_t n);
    absl::string_view current_chunk_;
    absl::Nullable<absl::cord_internal::CordRep*> current_leaf_ = nullptr;
    size_t bytes_remaining_ = 0;
    CordRepBtreeReader btree_reader_;
  };
  ChunkIterator chunk_begin() const ABSL_ATTRIBUTE_LIFETIME_BOUND;
  ChunkIterator chunk_end() const ABSL_ATTRIBUTE_LIFETIME_BOUND;
  class ChunkRange {
   public:
    using value_type = absl::string_view;
    using reference = value_type&;
    using const_reference = const value_type&;
    using iterator = ChunkIterator;
    using const_iterator = ChunkIterator;
    explicit ChunkRange(absl::Nonnull<const Cord*> cord) : cord_(cord) {}
    ChunkIterator begin() const;
    ChunkIterator end() const;
   private:
    absl::Nonnull<const Cord*> cord_;
  };
  ChunkRange Chunks() const ABSL_ATTRIBUTE_LIFETIME_BOUND;
  class CharIterator {
   public:
    using iterator_category = std::input_iterator_tag;
    using value_type = char;
    using difference_type = ptrdiff_t;
    using pointer = absl::Nonnull<const char*>;
    using reference = const char&;
    CharIterator() = default;
    CharIterator& operator++();
    CharIterator operator++(int);
    bool operator==(const CharIterator& other) const;
    bool operator!=(const CharIterator& other) const;
    reference operator*() const;
    pointer operator->() const;
    friend Cord;
   private:
    explicit CharIterator(absl::Nonnull<const Cord*> cord)
        : chunk_iterator_(cord) {}
    ChunkIterator chunk_iterator_;
  };
  static Cord AdvanceAndRead(absl::Nonnull<CharIterator*> it, size_t n_bytes);
  static void Advance(absl::Nonnull<CharIterator*> it, size_t n_bytes);
  static absl::string_view ChunkRemaining(const CharIterator& it);
  CharIterator char_begin() const ABSL_ATTRIBUTE_LIFETIME_BOUND;
  CharIterator char_end() const ABSL_ATTRIBUTE_LIFETIME_BOUND;
  class CharRange {
   public:
    using value_type = char;
    using reference = value_type&;
    using const_reference = const value_type&;
    using iterator = CharIterator;
    using const_iterator = CharIterator;
    explicit CharRange(absl::Nonnull<const Cord*> cord) : cord_(cord) {}
    CharIterator begin() const;
    CharIterator end() const;
   private:
    absl::Nonnull<const Cord*> cord_;
  };
  CharRange Chars() const ABSL_ATTRIBUTE_LIFETIME_BOUND;
  char operator[](size_t i) const;
  absl::optional<absl::string_view> TryFlat() const
      ABSL_ATTRIBUTE_LIFETIME_BOUND;
  absl::string_view Flatten() ABSL_ATTRIBUTE_LIFETIME_BOUND;
  CharIterator Find(absl::string_view needle) const;
  CharIterator Find(const absl::Cord& needle) const;
  friend void AbslFormatFlush(absl::Nonnull<absl::Cord*> cord,
                              absl::string_view part) {
    cord->Append(part);
  }
  template <typename Sink>
  friend void AbslStringify(Sink& sink, const absl::Cord& cord) {
    for (absl::string_view chunk : cord.Chunks()) {
      sink.Append(chunk);
    }
  }