Actually, the default case only triggers when `constant_kind_case()` returns something not in the switch. Setting `int64_value` doesn't change the case type - it just populates the field within the `kInt64Value` case. To properly trigger the default, I'd need to use a proto value that's outside the expected enum range, which isn't straightforward with protobuf.

The safest approach is either leaving the constant kind unset or verifying the behavior with an actual unknown enum value if the proto allows it. Since I can't easily construct an unknown case, I should test edge cases like maximum/minimum int64 values, large uint64 values, and precision loss in doubles instead.

For the final test suite, I'll focus on what's actually testable: round-trip conversions for each type, zero values, max/min boundaries, empty strings, null bytes, and a case that returns InvalidArgumentError. I'll also verify that Duration and Timestamp encode correctly by checking the proto fields they produce. When decoding, I need to ensure the values round-trip properly through the CEL types rather than just checking for proto presence.

For the round-trip conversion tests, I'll encode a Duration or Timestamp to proto format, then decode it back and verify it matches the original value. This validates both the encoding and decoding paths work correctly together.

I should also verify the function signatures match the header declarations. The parameter order looks consistent across both the declaration and implementation, so that's good to confirm.

Now I need to finalize the test structure. The `absl::visit` call on the Constant object will extract the variant from `kind()` and match it against the handlers I've defined, so I need to make sure my test cases cover all the variant types that could be stored in a Constant.

For boundary conditions, I'm considering what edge cases would stress the conversion logic—particularly around numeric precision and overflow scenarios with the various integer and floating-point types.
</think>