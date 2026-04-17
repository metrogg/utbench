#ifndef TENSORFLOW_TSL_LIB_IO_TABLE_H_
#define TENSORFLOW_TSL_LIB_IO_TABLE_H_
#include <stdint.h>
#include "tsl/lib/io/iterator.h"
namespace tsl {
class RandomAccessFile;
namespace table {
struct Options;
class Table {
 public:
  static absl::Status Open(const Options& options, tsl::RandomAccessFile* file,
                           uint64 file_size, Table** table);
  ~Table();
  Iterator* NewIterator() const;
  uint64 ApproximateOffsetOf(const StringPiece& key) const;
 private:
  struct Rep;
  Rep* rep_;
  explicit Table(Rep* rep) { rep_ = rep; }
  static Iterator* BlockReader(void*, const StringPiece&);
  absl::Status InternalGet(const StringPiece& key, void* arg,
                           void (*handle_result)(void* arg,
                                                 const StringPiece& k,
                                                 const StringPiece& v));
  Table(const Table&);
  void operator=(const Table&);
};
}  
}  
#endif  
#include "tsl/lib/io/table.h"
#include "tsl/lib/io/block.h"
#include "tsl/lib/io/cache.h"
#include "tsl/lib/io/format.h"
#include "tsl/lib/io/table_options.h"
#include "tsl/lib/io/two_level_iterator.h"
#include "tsl/platform/coding.h"
#include "tsl/platform/env.h"
#include "tsl/platform/errors.h"
namespace tsl {
namespace table {
struct Table::Rep {
  ~Rep() { delete index_block; }
  Options options;
  absl::Status status;
  RandomAccessFile* file;
  uint64 cache_id;
  BlockHandle metaindex_handle;  
  Block* index_block;
};
absl::Status Table::Open(const Options& options, RandomAccessFile* file,
                         uint64 size, Table** table) {
  *table = nullptr;
  if (size < Footer::kEncodedLength) {
    return errors::DataLoss("file is too short to be an sstable");
  }
  char footer_space[Footer::kEncodedLength];
  StringPiece footer_input;
  absl::Status s =
      file->Read(size - Footer::kEncodedLength, Footer::kEncodedLength,
                 &footer_input, footer_space);
  if (!s.ok()) return s;
  Footer footer;
  s = footer.DecodeFrom(&footer_input);
  if (!s.ok()) return s;
  BlockContents contents;
  Block* index_block = nullptr;
  if (s.ok()) {
    s = ReadBlock(file, footer.index_handle(), &contents);
  }
  if (s.ok()) {
    index_block = new Block(contents);
    Rep* rep = new Table::Rep;
    rep->options = options;
    rep->file = file;
    rep->metaindex_handle = footer.metaindex_handle();
    rep->index_block = index_block;
    rep->cache_id = (options.block_cache ? options.block_cache->NewId() : 0);
    *table = new Table(rep);
  } else {
    if (index_block) delete index_block;
  }
  return s;
}
Table::~Table() { delete rep_; }
static void DeleteBlock(void* arg, void* ignored) {
  delete reinterpret_cast<Block*>(arg);
}
static void DeleteCachedBlock(const absl::string_view&, void* value) {
  Block* block = reinterpret_cast<Block*>(value);
  delete block;
}
static void ReleaseBlock(void* arg, void* h) {
  Cache* cache = reinterpret_cast<Cache*>(arg);
  Cache::Handle* handle = reinterpret_cast<Cache::Handle*>(h);
  cache->Release(handle);
}
Iterator* Table::BlockReader(void* arg, const StringPiece& index_value) {
  Table* table = reinterpret_cast<Table*>(arg);
  Cache* block_cache = table->rep_->options.block_cache;
  Block* block = nullptr;
  Cache::Handle* cache_handle = nullptr;
  BlockHandle handle;
  StringPiece input = index_value;
  absl::Status s = handle.DecodeFrom(&input);
  if (s.ok()) {
    BlockContents contents;
    if (block_cache != nullptr) {
      char cache_key_buffer[16];
      core::EncodeFixed64(cache_key_buffer, table->rep_->cache_id);
      core::EncodeFixed64(cache_key_buffer + 8, handle.offset());
      absl::string_view key(cache_key_buffer, sizeof(cache_key_buffer));
      cache_handle = block_cache->Lookup(key);
      if (cache_handle != nullptr) {
        block = reinterpret_cast<Block*>(block_cache->Value(cache_handle));
      } else {
        s = ReadBlock(table->rep_->file, handle, &contents);
        if (s.ok()) {
          block = new Block(contents);
          cache_handle = block_cache->Insert(key, block, block->size(),
                                             &DeleteCachedBlock);
        }
      }
    } else {
      s = ReadBlock(table->rep_->file, handle, &contents);
      if (s.ok()) {
        block = new Block(contents);
      }
    }
  }
  Iterator* iter;
  if (block != nullptr) {
    iter = block->NewIterator();
    if (cache_handle == nullptr) {
      iter->RegisterCleanup(&DeleteBlock, block, nullptr);
    } else {
      iter->RegisterCleanup(&ReleaseBlock, block_cache, cache_handle);
    }
  } else {
    iter = NewErrorIterator(s);
  }
  return iter;
}
Iterator* Table::NewIterator() const {
  return NewTwoLevelIterator(rep_->index_block->NewIterator(),
                             &Table::BlockReader, const_cast<Table*>(this));
}
absl::Status Table::InternalGet(const StringPiece& k, void* arg,
                                void (*saver)(void*, const StringPiece&,
                                              const StringPiece&)) {
  absl::Status s;
  Iterator* iiter = rep_->index_block->NewIterator();
  iiter->Seek(k);
  if (iiter->Valid()) {
    Iterator* block_iter = BlockReader(this, iiter->value());
    block_iter->Seek(k);
    if (block_iter->Valid()) {
      (*saver)(arg, block_iter->key(), block_iter->value());
    }
    s = block_iter->status();
    delete block_iter;
  }
  if (s.ok()) {
    s = iiter->status();
  }
  delete iiter;
  return s;
}
uint64 Table::ApproximateOffsetOf(const StringPiece& key) const {
  Iterator* index_iter = rep_->index_block->NewIterator();
  index_iter->Seek(key);
  uint64 result;
  if (index_iter->Valid()) {
    BlockHandle handle;
    StringPiece input = index_iter->value();
    absl::Status s = handle.DecodeFrom(&input);
    if (s.ok()) {
      result = handle.offset();
    } else {
      result = rep_->metaindex_handle.offset();
    }
  } else {
    result = rep_->metaindex_handle.offset();
  }
  delete index_iter;
  return result;
}
}  
}  