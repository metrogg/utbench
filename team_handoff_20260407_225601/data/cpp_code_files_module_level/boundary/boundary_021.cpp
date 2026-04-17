#ifndef TENSORFLOW_COMPILER_MLIR_LITE_STABLEHLO_TRANSFORMS_LEGALIZE_HLO_CONVERSIONS_CUSTOM_CALL_H_
#define TENSORFLOW_COMPILER_MLIR_LITE_STABLEHLO_TRANSFORMS_LEGALIZE_HLO_CONVERSIONS_CUSTOM_CALL_H_
#include <optional>
#include "mlir/Transforms/DialectConversion.h"  
#include "xla/mlir_hlo/mhlo/IR/hlo_ops.h"
namespace mlir {
namespace odml {
class ConvertCustomCallOp : public OpConversionPattern<mhlo::CustomCallOp> {
 public:
  using OpConversionPattern::OpConversionPattern;
  LogicalResult matchAndRewrite(
      mhlo::CustomCallOp mhlo_custom_call, OpAdaptor adaptor,
      ConversionPatternRewriter& rewriter) const final;
};
std::optional<bool> IsCustomCallLegal(mhlo::CustomCallOp op);
}  
}  
#endif  
#include "tensorflow/compiler/mlir/lite/stablehlo/transforms/legalize_hlo_conversions/custom_call.h"
#include <optional>
#include "mlir/IR/BuiltinAttributes.h"  
#include "mlir/Support/LLVM.h"  
#include "mlir/Support/LogicalResult.h"  
#include "mlir/Transforms/DialectConversion.h"  
#include "tensorflow/compiler/mlir/lite/ir/tfl_ops.h"  
#include "xla/mlir_hlo/mhlo/IR/hlo_ops.h"
namespace mlir {
namespace odml {
LogicalResult ConvertCustomCallOp::matchAndRewrite(
    mhlo::CustomCallOp mhlo_custom_call, OpAdaptor adaptor,
    ConversionPatternRewriter& rewriter) const {
  auto tfl_custom = rewriter.create<TFL::CustomOp>(
      mhlo_custom_call.getLoc(), mhlo_custom_call.getResultTypes(),
      mhlo_custom_call.getInputs());
  tfl_custom.setCustomCodeAttr(
      rewriter.getStringAttr(mhlo_custom_call.getCallTargetName()));
  if (auto bc = mhlo_custom_call.getBackendConfig()) {
    if (auto stringattr = mlir::dyn_cast_or_null<mlir::StringAttr>(*bc)) {
      tfl_custom.setCustomOptionAttr(
          TFL::ConstBytesAttr::get(rewriter.getContext(), stringattr));
    }
  } else {
    tfl_custom.setCustomOptionAttr(
        TFL::ConstBytesAttr::get(rewriter.getContext(), ""));
  }
  rewriter.replaceOp(mhlo_custom_call, tfl_custom);
  return success();
}
std::optional<bool> IsCustomCallLegal(mhlo::CustomCallOp op) {
  if (op.getCallTargetName().starts_with("custom_call.")) {
    auto bc = op.getBackendConfig();
    if (!bc || mlir::isa<mlir::StringAttr>(*bc)) {
      return false;
    }
  }
  return true;
}
}  
}  