#ifndef XLA_SERVICE_GPU_MODEL_AFFINE_MAP_PRINTER_H_
#define XLA_SERVICE_GPU_MODEL_AFFINE_MAP_PRINTER_H_
#include <cstdint>
#include <ostream>
#include <string>
#include <string_view>
#include "absl/types/span.h"
#include "llvm/ADT/DenseMap.h"
#include "llvm/ADT/StringRef.h"
#include "llvm/ADT/Twine.h"
#include "mlir/IR/AffineExpr.h"  
#include "mlir/IR/AffineMap.h"  
namespace xla {
namespace gpu {
class AffineMapPrinter {
 public:
  AffineMapPrinter() = default;
  AffineMapPrinter(AffineMapPrinter&& other) = default;
  AffineMapPrinter& operator=(AffineMapPrinter&& other) = default;
  AffineMapPrinter(absl::Span<const std::string_view> dim_names,
                   absl::Span<const std::string_view> symbol_names);
  void SetSymbolName(int64_t symbol_id, llvm::StringRef name);
  void SetDimensionName(int64_t dim_id, llvm::StringRef name);
  std::string GetSymbolName(int64_t symbol_id) const;
  std::string GetDimensionName(int64_t dim_id) const;
  void Print(std::ostream& out, mlir::AffineMap affine_map) const;
  std::string ToString(mlir::AffineMap affine_map) const;
  void Print(std::ostream& out, mlir::AffineExpr affine_expr) const;
  std::string ToString(mlir::AffineExpr affine_expr) const;
 private:
  void PrintExprImpl(mlir::AffineExpr affine_expr, bool add_parentheses,
                     llvm::raw_ostream& os) const;
  llvm::DenseMap<unsigned, std::string> dim_id_to_name_;
  llvm::DenseMap<unsigned, std::string> symbol_id_to_name_;
};
}  
}  
#endif  
#include "xla/service/gpu/model/affine_map_printer.h"
#include <cstdint>
#include <ostream>
#include <string>
#include <string_view>
#include "absl/strings/str_cat.h"
#include "absl/types/span.h"
#include "llvm/ADT/STLExtras.h"
#include "llvm/Support/raw_ostream.h"
#include "mlir/IR/AffineExpr.h"  
#include "mlir/IR/AffineMap.h"  
#include "mlir/Support/LLVM.h"  
namespace xla {
namespace gpu {
namespace {
using mlir::AffineBinaryOpExpr;
using mlir::AffineConstantExpr;
using mlir::AffineDimExpr;
using mlir::AffineExpr;
using mlir::AffineExprKind;
using mlir::AffineMap;
using mlir::AffineSymbolExpr;
}  
AffineMapPrinter::AffineMapPrinter(
    absl::Span<const std::string_view> dim_names,
    absl::Span<const std::string_view> symbol_names) {
  dim_id_to_name_.reserve(dim_names.size());
  for (const auto& [index, name] : llvm::enumerate(dim_names)) {
    dim_id_to_name_[index] = name;
  }
  symbol_id_to_name_.reserve(symbol_names.size());
  for (const auto& [index, name] : llvm::enumerate(symbol_names)) {
    symbol_id_to_name_[index] = name;
  }
}
void AffineMapPrinter::Print(std::ostream& out, AffineMap affine_map) const {
  out << ToString(affine_map);
}
std::string AffineMapPrinter::ToString(AffineMap affine_map) const {
  std::string s;
  llvm::raw_string_ostream ss(s);
  if (dim_id_to_name_.empty() && symbol_id_to_name_.empty()) {
    affine_map.print(ss);
    return s;
  }
  int dim_count = affine_map.getNumDims();
  ss << '(';
  for (int i = 0; i < dim_count - 1; ++i) {
    ss << GetDimensionName(i) << ", ";
  }
  if (dim_count >= 1) {
    ss << GetDimensionName(dim_count - 1);
  }
  ss << ')';
  int symbol_count = affine_map.getNumSymbols();
  if (symbol_count != 0) {
    ss << '[';
    for (unsigned i = 0; i < symbol_count - 1; ++i) {
      ss << GetSymbolName(i) << ", ";
    }
    if (affine_map.getNumSymbols() >= 1) {
      ss << GetSymbolName(symbol_count - 1);
    }
    ss << ']';
  }
  ss << " -> (";
  llvm::interleaveComma(affine_map.getResults(), ss, [&](AffineExpr expr) {
    PrintExprImpl(expr, false, ss);
  });
  ss << ')';
  return s;
}
void AffineMapPrinter::Print(std::ostream& out,
                             mlir::AffineExpr affine_expr) const {
  out << ToString(affine_expr);
}
std::string AffineMapPrinter::ToString(mlir::AffineExpr affine_expr) const {
  std::string s;
  llvm::raw_string_ostream ss(s);
  PrintExprImpl(affine_expr, false, ss);
  return s;
}
void AffineMapPrinter::PrintExprImpl(const mlir::AffineExpr affine_expr,
                                     bool add_parentheses,
                                     llvm::raw_ostream& os) const {
  const char* binopSpelling = nullptr;
  switch (affine_expr.getKind()) {
    case AffineExprKind::SymbolId: {
      unsigned symbol_id =
          mlir::cast<AffineSymbolExpr>(affine_expr).getPosition();
      os << GetSymbolName(symbol_id);
      return;
    }
    case AffineExprKind::DimId: {
      unsigned dim_id = mlir::cast<AffineDimExpr>(affine_expr).getPosition();
      os << GetDimensionName(dim_id);
      return;
    }
    case AffineExprKind::Constant:
      os << mlir::cast<AffineConstantExpr>(affine_expr).getValue();
      return;
    case AffineExprKind::Add:
      binopSpelling = " + ";
      break;
    case AffineExprKind::Mul:
      binopSpelling = " * ";
      break;
    case AffineExprKind::FloorDiv:
      binopSpelling = " floordiv ";
      break;
    case AffineExprKind::CeilDiv:
      binopSpelling = " ceildiv ";
      break;
    case AffineExprKind::Mod:
      binopSpelling = " mod ";
      break;
  }
  auto binOp = mlir::cast<AffineBinaryOpExpr>(affine_expr);
  AffineExpr lhsExpr = binOp.getLHS();
  AffineExpr rhsExpr = binOp.getRHS();
  if (binOp.getKind() != AffineExprKind::Add) {
    if (add_parentheses) {
      os << '(';
    }
    auto rhsConst = mlir::dyn_cast<AffineConstantExpr>(rhsExpr);
    if (rhsConst && binOp.getKind() == AffineExprKind::Mul &&
        rhsConst.getValue() == -1) {
      os << "-";
      PrintExprImpl(lhsExpr, true, os);
      if (add_parentheses) {
        os << ')';
      }
      return;
    }
    PrintExprImpl(lhsExpr, true, os);
    os << binopSpelling;
    PrintExprImpl(rhsExpr, true, os);
    if (add_parentheses) {
      os << ')';
    }
    return;
  }
  if (add_parentheses) {
    os << '(';
  }
  if (auto rhs = mlir::dyn_cast<AffineBinaryOpExpr>(rhsExpr)) {
    if (rhs.getKind() == AffineExprKind::Mul) {
      AffineExpr rrhsExpr = rhs.getRHS();
      if (auto rrhs = mlir::dyn_cast<AffineConstantExpr>(rrhsExpr)) {
        if (rrhs.getValue() == -1) {
          PrintExprImpl(lhsExpr, false, os);
          os << " - ";
          if (rhs.getLHS().getKind() == AffineExprKind::Add) {
            PrintExprImpl(rhs.getLHS(), true, os);
          } else {
            PrintExprImpl(rhs.getLHS(), false, os);
          }
          if (add_parentheses) {
            os << ')';
          }
          return;
        }
        if (rrhs.getValue() < -1) {
          PrintExprImpl(lhsExpr, false, os);
          os << " - ";
          PrintExprImpl(rhs.getLHS(), true, os);
          os << " * " << -rrhs.getValue();
          if (add_parentheses) {
            os << ')';
          }
          return;
        }
      }
    }
  }
  if (auto rhsConst = mlir::dyn_cast<AffineConstantExpr>(rhsExpr)) {
    if (rhsConst.getValue() < 0) {
      PrintExprImpl(lhsExpr, false, os);
      os << " - " << -rhsConst.getValue();
      if (add_parentheses) {
        os << ')';
      }
      return;
    }
  }
  PrintExprImpl(lhsExpr, false, os);
  os << " + ";
  PrintExprImpl(rhsExpr, false, os);
  if (add_parentheses) {
    os << ')';
  }
}
void AffineMapPrinter::SetSymbolName(int64_t symbol_id, llvm::StringRef name) {
  symbol_id_to_name_[symbol_id] = name;
}
void AffineMapPrinter::SetDimensionName(int64_t dim_id, llvm::StringRef name) {
  dim_id_to_name_[dim_id] = name;
}
std::string AffineMapPrinter::GetSymbolName(int64_t symbol_id) const {
  auto it = symbol_id_to_name_.find(symbol_id);
  if (it == symbol_id_to_name_.end()) {
    return absl::StrCat("s", symbol_id);
  }
  return it->second;
}
std::string AffineMapPrinter::GetDimensionName(int64_t dim_id) const {
  auto it = dim_id_to_name_.find(dim_id);
  if (it == dim_id_to_name_.end()) {
    return absl::StrCat("d", dim_id);
  }
  return it->second;
}
}  
}  