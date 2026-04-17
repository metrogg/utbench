#ifndef TENSORFLOW_COMPILER_MLIR_QUANTIZATION_STABLEHLO_UTILS_FILL_QUANTIZATION_OPTIONS_H_
#define TENSORFLOW_COMPILER_MLIR_QUANTIZATION_STABLEHLO_UTILS_FILL_QUANTIZATION_OPTIONS_H_
#include "mlir/Support/LogicalResult.h"  
#include "tensorflow/compiler/mlir/quantization/stablehlo/quantization_options.pb.h"
namespace mlir::quant::stablehlo {
using ::stablehlo::quantization::QuantizationOptions;
QuantizationOptions FillPresetQuantizationOptions(
    QuantizationOptions quantization_options);
LogicalResult GetActivationBitWidth(QuantizationOptions quantization_options,
                                    int* bit_width);
}  
#endif  
#include "mlir/Support/LogicalResult.h"  
#include "tensorflow/compiler/mlir/quantization/stablehlo/quantization_options.pb.h"
namespace mlir::quant::stablehlo {
using ::stablehlo::quantization::CustomQuantizationMethod;
using ::stablehlo::quantization::PresetQuantizationMethod;
using ::stablehlo::quantization::QuantizationComponentSpec;
using ::stablehlo::quantization::QuantizationOptions;
using QuantizationComponent =
    ::stablehlo::quantization::QuantizationComponentSpec_QuantizationComponent;
using BitType = ::stablehlo::quantization::QuantizationComponentSpec_BitType;
using BitWidth = ::stablehlo::quantization::QuantizationComponentSpec_BitWidth;
void SetQuantizationComponentSpec(QuantizationComponentSpec* spec,
                                  const QuantizationComponent& component,
                                  const BitType bit_type,
                                  const BitWidth bit_width) {
  spec->set_quantization_component(component);
  spec->set_bit_type(bit_type);
  spec->set_bit_width(bit_width);
}
::stablehlo::quantization::QuantizationOptions FillPresetQuantizationOptions(
    ::stablehlo::quantization::QuantizationOptions quantization_options_) {
  CustomQuantizationMethod custom_method =
      quantization_options_.quantization_method().custom_quantization_method();
  QuantizationComponentSpec *activation_component, *weight_component,
      *bias_component;
  const auto preset_method = quantization_options_.quantization_method()
                                 .preset_quantization_method()
                                 .preset_method();
  if (!preset_method) return quantization_options_;
  switch (preset_method) {
    case PresetQuantizationMethod::FLOAT16:
      weight_component = custom_method.add_quantization_component_spec();
      SetQuantizationComponentSpec(weight_component,
                                   QuantizationComponentSpec::COMPONENT_WEIGHT,
                                   QuantizationComponentSpec::BIT_TYPE_FLOAT,
                                   QuantizationComponentSpec::BIT_WIDTH_16);
      bias_component = custom_method.add_quantization_component_spec();
      SetQuantizationComponentSpec(bias_component,
                                   QuantizationComponentSpec::COMPONENT_BIAS,
                                   QuantizationComponentSpec::BIT_TYPE_FLOAT,
                                   QuantizationComponentSpec::BIT_WIDTH_16);
      break;
    case PresetQuantizationMethod::WEIGHT_ONLY:
      weight_component = custom_method.add_quantization_component_spec();
      SetQuantizationComponentSpec(weight_component,
                                   QuantizationComponentSpec::COMPONENT_WEIGHT,
                                   QuantizationComponentSpec::BIT_TYPE_INT,
                                   QuantizationComponentSpec::BIT_WIDTH_8);
      break;
    case PresetQuantizationMethod::POST_TRAINING_QUANTIZATION_STATIC_RANGE:
      activation_component = custom_method.add_quantization_component_spec();
      SetQuantizationComponentSpec(
          activation_component, QuantizationComponentSpec::COMPONENT_ACTIVATION,
          QuantizationComponentSpec::BIT_TYPE_INT,
          QuantizationComponentSpec::BIT_WIDTH_8);
      weight_component = custom_method.add_quantization_component_spec();
      SetQuantizationComponentSpec(weight_component,
                                   QuantizationComponentSpec::COMPONENT_WEIGHT,
                                   QuantizationComponentSpec::BIT_TYPE_INT,
                                   QuantizationComponentSpec::BIT_WIDTH_8);
      bias_component = custom_method.add_quantization_component_spec();
      SetQuantizationComponentSpec(bias_component,
                                   QuantizationComponentSpec::COMPONENT_BIAS,
                                   QuantizationComponentSpec::BIT_TYPE_INT,
                                   QuantizationComponentSpec::BIT_WIDTH_32);
      break;
    default:
      break;
  }
  *quantization_options_.mutable_quantization_method()
       ->mutable_custom_quantization_method() = custom_method;
  return quantization_options_;
}
LogicalResult GetActivationBitWidth(QuantizationOptions quantization_options,
                                    int* bit_width) {
  CustomQuantizationMethod custom_method =
      quantization_options.quantization_method().custom_quantization_method();
  for (const auto& component : custom_method.quantization_component_spec()) {
    if (component.quantization_component() ==
        QuantizationComponentSpec::COMPONENT_ACTIVATION) {
      switch (component.bit_width()) {
        case QuantizationComponentSpec::BIT_WIDTH_4:
          *bit_width = 4;
          return success();
          break;
        case QuantizationComponentSpec::BIT_WIDTH_8:
          *bit_width = 8;
          return success();
          break;
        case QuantizationComponentSpec::BIT_WIDTH_16:
          *bit_width = 16;
          return success();
          break;
        case QuantizationComponentSpec::BIT_WIDTH_32:
          *bit_width = 32;
          return success();
          break;
        default:
          break;
      }
    }
  }
  return failure();
}
}  