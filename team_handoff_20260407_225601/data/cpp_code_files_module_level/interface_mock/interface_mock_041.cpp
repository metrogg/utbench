#ifndef TENSORFLOW_COMPILER_MLIR_LITE_STABLEHLO_TRANSFORMS_LEGALIZE_HLO_CONVERSIONS_REDUCE_H_
#define TENSORFLOW_COMPILER_MLIR_LITE_STABLEHLO_TRANSFORMS_LEGALIZE_HLO_CONVERSIONS_REDUCE_H_
#include <cstdint>
#include <optional>
#include "llvm/ADT/APInt.h"
#include "llvm/ADT/StringRef.h"
#include "mlir/Dialect/Arith/IR/Arith.h"  
#include "mlir/IR/BuiltinAttributes.h"  
#include "mlir/IR/BuiltinTypeInterfaces.h"  
#include "mlir/IR/BuiltinTypes.h"  
#include "mlir/IR/Matchers.h"  
#include "mlir/IR/Operation.h"  
#include "mlir/IR/Value.h"  
#include "mlir/Support/LLVM.h"  
#include "mlir/Support/LogicalResult.h"  
#include "mlir/Transforms/DialectConversion.h"  
#include "tensorflow/compiler/mlir/lite/stablehlo/transforms/hlo_matchers.h"
#include "xla/mlir_hlo/mhlo/IR/hlo_ops.h"
namespace mlir {
namespace odml {
LogicalResult MatchReduceToArgMinMaxType1(mhlo::ReduceOp reduce_op,
                                          bool is_float, bool is_argmax);
LogicalResult MatchReduceToArgMinMaxType2(mhlo::ReduceOp reduce_op,
                                          bool is_argmax);
template <typename Reduce, typename ArgReduce, typename BooleanReduce,
          bool is_argmax>
class ConvertReduceOpToArgMinMax : public OpConversionPattern<mhlo::ReduceOp> {
 public:
  using OpConversionPattern::OpConversionPattern;
  LogicalResult matchAndRewrite(
      mhlo::ReduceOp reduce_op, OpAdaptor adaptor,
      ConversionPatternRewriter& rewriter) const final {
    if (reduce_op.getInputs().size() != 2) return failure();
    if (reduce_op.getDimensions().getNumElements() != 1) return failure();
    DenseElementsAttr operand_init;
    if (!matchPattern(reduce_op.getInitValues().front(),
                      m_Constant(&operand_init)))
      return failure();
    if (!IsValueInitValue(operand_init)) return failure();
    DenseElementsAttr iota_init;
    if (!matchPattern(reduce_op.getInitValues().back(), m_Constant(&iota_init)))
      return failure();
    if (iota_init.getValues<APInt>()[0] != 0) return failure();
    Value iota = reduce_op.getInputs().back();
    if (!MatchIota(reduce_op.getDimensions(), iota)) return failure();
    const bool is_float = mlir::isa<FloatType>(operand_init.getElementType());
    if (failed(MatchReduceToArgMinMaxType1(reduce_op, is_float, is_argmax)) &&
        failed(MatchReduceToArgMinMaxType2(reduce_op, is_argmax)))
      return rewriter.notifyMatchFailure(
          reduce_op, "Unsupported Reduce -> ArgMax/ArgMin pattern");
    Value operand = reduce_op.getInputs().front();
    int64_t axis = reduce_op.getDimensions().getValues<int64_t>()[0];
    auto dim_type = RankedTensorType::get({1}, rewriter.getI32Type());
    auto reduction_indices = rewriter.create<arith::ConstantOp>(
        reduce_op.getLoc(), dim_type,
        rewriter.getI32TensorAttr({static_cast<int32_t>(axis)}));
    if (!mlir::isa<ShapedType>(operand.getType())) return failure();
    auto operand_type = mlir::cast<ShapedType>(operand.getType());
    if (operand_type.getElementType().isInteger(1)) {
      auto tf_reduce_op = rewriter.create<BooleanReduce>(
          reduce_op.getLoc(), reduce_op->getResult(0).getType(), operand,
          reduction_indices,
          rewriter.getBoolAttr(false));
      auto tf_argreduce_op = rewriter.create<ArgReduce>(
          reduce_op.getLoc(), reduce_op->getResult(1).getType(), operand,
          reduction_indices);
      rewriter.replaceOp(reduce_op, {tf_reduce_op, tf_argreduce_op});
    } else {
      auto tf_reduce_op = rewriter.create<Reduce>(
          reduce_op.getLoc(), reduce_op->getResult(0).getType(), operand,
          reduction_indices,
          rewriter.getBoolAttr(false));
      auto tf_argreduce_op = rewriter.create<ArgReduce>(
          reduce_op.getLoc(), reduce_op->getResult(1).getType(), operand,
          reduction_indices);
      rewriter.replaceOp(reduce_op, {tf_reduce_op, tf_argreduce_op});
    }
    return success();
  }
  virtual bool IsValueInitValue(const DenseElementsAttr& attr) const = 0;
};
std::optional<bool> IsReduceOpLegal(mhlo::ReduceOp reduce_op);
}  
}  
#endif  
#include "tensorflow/compiler/mlir/lite/stablehlo/transforms/legalize_hlo_conversions/reduce.h"
#include <optional>
#include "llvm/Support/Casting.h"
#include "mlir/IR/Block.h"  
#include "mlir/Support/LLVM.h"  
#include "mlir/Support/LogicalResult.h"  
#include "xla/mlir_hlo/mhlo/IR/hlo_ops.h"
namespace mlir {
namespace odml {
LogicalResult MatchReduceToArgMinMaxType2(mhlo::ReduceOp reduce_op,
                                          bool is_argmax) {
  Block& body = reduce_op.getBody().front();
  if (body.getNumArguments() != 4) return failure();
  mhlo::ReturnOp return_op = dyn_cast<mhlo::ReturnOp>(body.back());
  if (!return_op || return_op.getNumOperands() != 2) return failure();
  mhlo::SelectOp value_select = llvm::dyn_cast_or_null<mhlo::SelectOp>(
      return_op.getOperand(0).getDefiningOp());
  if (!value_select || value_select.getOnTrue() != body.getArgument(0) ||
      value_select.getOnFalse() != body.getArgument(2))
    return failure();
  auto compare_direction_included =
      is_argmax ? mhlo::ComparisonDirection::GE : mhlo::ComparisonDirection::LE;
  mhlo::CompareOp value_gt = llvm::dyn_cast_or_null<mhlo::CompareOp>(
      value_select.getOperand(0).getDefiningOp());
  if (!value_gt ||
      value_gt.getComparisonDirection() != compare_direction_included ||
      value_gt.getLhs() != body.getArgument(0) ||
      value_gt.getRhs() != body.getArgument(2))
    return failure();
  mhlo::SelectOp index_select = llvm::dyn_cast_or_null<mhlo::SelectOp>(
      return_op.getOperand(1).getDefiningOp());
  if (!index_select) return failure();
  mhlo::MinOp index_select_min = llvm::dyn_cast_or_null<mhlo::MinOp>(
      index_select.getOnTrue().getDefiningOp());
  if (!index_select_min || index_select_min.getLhs() != body.getArgument(1) ||
      index_select_min.getRhs() != body.getArgument(3))
    return failure();
  mhlo::SelectOp index_select_select = llvm::dyn_cast_or_null<mhlo::SelectOp>(
      index_select.getOnFalse().getDefiningOp());
  if (!index_select_select ||
      index_select_select.getOnTrue() != body.getArgument(1) ||
      index_select_select.getOnFalse() != body.getArgument(3) ||
      index_select_select.getOperand(0).getDefiningOp() != value_gt)
    return failure();
  mhlo::CompareOp value_eq = llvm::dyn_cast_or_null<mhlo::CompareOp>(
      index_select.getOperand(0).getDefiningOp());
  if (!value_eq ||
      value_eq.getComparisonDirection() != mhlo::ComparisonDirection::EQ ||
      value_eq.getLhs() != body.getArgument(0) ||
      value_eq.getRhs() != body.getArgument(2))
    return failure();
  return success();
}
LogicalResult MatchReduceToArgMinMaxType1(mhlo::ReduceOp reduce_op,
                                          bool is_float, bool is_argmax) {
  Block& body = reduce_op.getBody().front();
  if (body.getNumArguments() != 4) return failure();
  mhlo::ReturnOp return_op = dyn_cast<mhlo::ReturnOp>(body.back());
  if (!return_op || return_op.getNumOperands() != 2) return failure();
  mhlo::SelectOp value_select = llvm::dyn_cast_or_null<mhlo::SelectOp>(
      return_op.getOperand(0).getDefiningOp());
  if (!value_select || value_select.getOnTrue() != body.getArgument(0) ||
      value_select.getOnFalse() != body.getArgument(2))
    return failure();
  auto compare_direction =
      is_argmax ? mhlo::ComparisonDirection::GT : mhlo::ComparisonDirection::LT;
  if (is_float) {
    mhlo::OrOp value_or = llvm::dyn_cast_or_null<mhlo::OrOp>(
        value_select.getOperand(0).getDefiningOp());
    if (!value_or) return failure();
    mhlo::CompareOp value_gt = llvm::dyn_cast_or_null<mhlo::CompareOp>(
        value_or.getLhs().getDefiningOp());
    if (!value_gt || value_gt.getComparisonDirection() != compare_direction ||
        value_gt.getLhs() != body.getArgument(0) ||
        value_gt.getRhs() != body.getArgument(2))
      return failure();
    mhlo::CompareOp value_ne = llvm::dyn_cast_or_null<mhlo::CompareOp>(
        value_or.getRhs().getDefiningOp());
    if (!value_ne ||
        value_ne.getComparisonDirection() != mhlo::ComparisonDirection::NE ||
        value_ne.getLhs() != body.getArgument(0) ||
        value_ne.getRhs() != body.getArgument(0))
      return failure();
  } else {
    mhlo::CompareOp value_gt = llvm::dyn_cast_or_null<mhlo::CompareOp>(
        value_select.getOperand(0).getDefiningOp());
    if (!value_gt || value_gt.getComparisonDirection() != compare_direction ||
        value_gt.getLhs() != body.getArgument(0) ||
        value_gt.getRhs() != body.getArgument(2))
      return failure();
  }
  mhlo::SelectOp index_select = llvm::dyn_cast_or_null<mhlo::SelectOp>(
      return_op.getOperand(1).getDefiningOp());
  if (!index_select || index_select.getOnTrue() != body.getArgument(1) ||
      index_select.getOnFalse() != body.getArgument(3))
    return failure();
  mhlo::OrOp index_or = llvm::dyn_cast_or_null<mhlo::OrOp>(
      index_select.getPred().getDefiningOp());
  if (!index_or || index_or.getLhs() != value_select.getPred())
    return failure();
  mhlo::AndOp index_and =
      llvm::dyn_cast_or_null<mhlo::AndOp>(index_or.getRhs().getDefiningOp());
  if (!index_and) return failure();
  mhlo::CompareOp value_eq = llvm::dyn_cast_or_null<mhlo::CompareOp>(
      index_and.getLhs().getDefiningOp());
  if (!value_eq ||
      value_eq.getComparisonDirection() != mhlo::ComparisonDirection::EQ ||
      value_eq.getLhs() != body.getArgument(0) ||
      value_eq.getRhs() != body.getArgument(2))
    return failure();
  mhlo::CompareOp index_lt = llvm::dyn_cast_or_null<mhlo::CompareOp>(
      index_and.getRhs().getDefiningOp());
  if (!index_lt ||
      index_lt.getComparisonDirection() != mhlo::ComparisonDirection::LT ||
      index_lt.getLhs() != body.getArgument(1) ||
      index_lt.getRhs() != body.getArgument(3))
    return failure();
  return success();
}
std::optional<bool> IsReduceOpLegal(mhlo::ReduceOp reduce_op) {
  if (succeeded(MatchReduceToArgMinMaxType1(reduce_op, true, true)) ||
      succeeded(MatchReduceToArgMinMaxType1(reduce_op, false, true)) ||
      succeeded(MatchReduceToArgMinMaxType1(reduce_op, true, false)) ||
      succeeded(MatchReduceToArgMinMaxType1(reduce_op, false, false)) ||
      succeeded(MatchReduceToArgMinMaxType2(reduce_op, false)) ||
      succeeded(MatchReduceToArgMinMaxType2(reduce_op, true))) {
    return false;
  }
  return true;
}
}  
}  