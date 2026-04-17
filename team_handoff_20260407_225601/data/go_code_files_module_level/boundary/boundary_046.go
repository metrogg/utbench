func typeMarshaler(t reflect.Type, tags []string, nozero, oneof bool) (sizer, marshaler) {
	encoding := tags[0]

	pointer := false
	slice := false
	if t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8 {
		slice = true
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		pointer = true
		t = t.Elem()
	}

	packed := false
	proto3 := false
	validateUTF8 := true
	for i := 2; i < len(tags); i++ {
		if tags[i] == "packed" {
			packed = true
		}
		if tags[i] == "proto3" {
			proto3 = true
		}
	}
	validateUTF8 = validateUTF8 && proto3

	switch t.Kind() {
	case reflect.Bool:
		if pointer {
			return sizeBoolPtr, appendBoolPtr
		}
		if slice {
			if packed {
				return sizeBoolPackedSlice, appendBoolPackedSlice
			}
			return sizeBoolSlice, appendBoolSlice
		}
		if nozero {
			return sizeBoolValueNoZero, appendBoolValueNoZero
		}
		return sizeBoolValue, appendBoolValue
	case reflect.Uint32:
		switch encoding {
		case "fixed32":
			if pointer {
				return sizeFixed32Ptr, appendFixed32Ptr
			}
			if slice {
				if packed {
					return sizeFixed32PackedSlice, appendFixed32PackedSlice
				}
				return sizeFixed32Slice, appendFixed32Slice
			}
			if nozero {
				return sizeFixed32ValueNoZero, appendFixed32ValueNoZero
			}
			return sizeFixed32Value, appendFixed32Value
		case "varint":
			if pointer {
				return sizeVarint32Ptr, appendVarint32Ptr
			}
			if slice {
				if packed {
					return sizeVarint32PackedSlice, appendVarint32PackedSlice
				}
				return sizeVarint32Slice, appendVarint32Slice
			}
			if nozero {
				return sizeVarint32ValueNoZero, appendVarint32ValueNoZero
			}
			return sizeVarint32Value, appendVarint32Value
		}
	case reflect.Int32:
		switch encoding {
		case "fixed32":
			if pointer {
				return sizeFixedS32Ptr, appendFixedS32Ptr
			}
			if slice {
				if packed {
					return sizeFixedS32PackedSlice, appendFixedS32PackedSlice
				}
				return sizeFixedS32Slice, appendFixedS32Slice
			}
			if nozero {
				return sizeFixedS32ValueNoZero, appendFixedS32ValueNoZero
			}
			return sizeFixedS32Value, appendFixedS32Value
		case "varint":
			if pointer {
				return sizeVarintS32Ptr, appendVarintS32Ptr
			}
			if slice {
				if packed {
					return sizeVarintS32PackedSlice, appendVarintS32PackedSlice
				}
				return sizeVarintS32Slice, appendVarintS32Slice
			}
			if nozero {
				return sizeVarintS32ValueNoZero, appendVarintS32ValueNoZero
			}
			return sizeVarintS32Value, appendVarintS32Value
		case "zigzag32":
			if pointer {
				return sizeZigzag32Ptr, appendZigzag32Ptr
			}
			if slice {
				if packed {
					return sizeZigzag32PackedSlice, appendZigzag32PackedSlice
				}
				return sizeZigzag32Slice, appendZigzag32Slice
			}
			if nozero {
				return sizeZigzag32ValueNoZero, appendZigzag32ValueNoZero
			}
			return sizeZigzag32Value, appendZigzag32Value
		}
	case reflect.Uint64:
		switch encoding {
		case "fixed64":
			if pointer {
				return sizeFixed64Ptr, appendFixed64Ptr
			}
			if slice {
				if packed {
					return sizeFixed64PackedSlice, appendFixed64PackedSlice
				}
				return sizeFixed64Slice, appendFixed64Slice
			}
			if nozero {
				return sizeFixed64ValueNoZero, appendFixed64ValueNoZero
			}
			return sizeFixed64Value, appendFixed64Value
		case "varint":
			if pointer {
				return sizeVarint64Ptr, appendVarint64Ptr
			}
			if slice {
				if packed {
					return sizeVarint64PackedSlice, appendVarint64PackedSlice
				}
				return sizeVarint64Slice, appendVarint64Slice
			}
			if nozero {
				return sizeVarint64ValueNoZero, appendVarint64ValueNoZero
			}
			return sizeVarint64Value, appendVarint64Value
		}
	case reflect.Int64:
		switch encoding {
		case "fixed64":
			if pointer {
				return sizeFixedS64Ptr, appendFixedS64Ptr
			}
			if slice {
				if packed {
					return sizeFixedS64PackedSlice, appendFixedS64PackedSlice
				}
				return sizeFixedS64Slice, appendFixedS64Slice
			}
			if nozero {
				return sizeFixedS64ValueNoZero, appendFixedS64ValueNoZero
			}
			return sizeFixedS64Value, appendFixedS64Value
		case "varint":
			if pointer {
				return sizeVarintS64Ptr, appendVarintS64Ptr
			}
			if slice {
				if packed {
					return sizeVarintS64PackedSlice, appendVarintS64PackedSlice
				}
				return sizeVarintS64Slice, appendVarintS64Slice
			}
			if nozero {
				return sizeVarintS64ValueNoZero, appendVarintS64ValueNoZero
			}
			return sizeVarintS64Value, appendVarintS64Value
		case "zigzag64":
			if pointer {
				return sizeZigzag64Ptr, appendZigzag64Ptr
			}
			if slice {
				if packed {
					return sizeZigzag64PackedSlice, appendZigzag64PackedSlice
				}
				return sizeZigzag64Slice, appendZigzag64Slice
			}
			if nozero {
				return sizeZigzag64ValueNoZero, appendZigzag64ValueNoZero
			}
			return sizeZigzag64Value, appendZigzag64Value
		}
	case reflect.Float32:
		if pointer {
			return sizeFloat32Ptr, appendFloat32Ptr
		}
		if slice {
			if packed {
				return sizeFloat32PackedSlice, appendFloat32PackedSlice
			}
			return sizeFloat32Slice, appendFloat32Slice
		}
		if nozero {
			return sizeFloat32ValueNoZero, appendFloat32ValueNoZero
		}
		return sizeFloat32Value, appendFloat32Value
	case reflect.Float64:
		if pointer {
			return sizeFloat64Ptr, appendFloat64Ptr
		}
		if slice {
			if packed {
				return sizeFloat64PackedSlice, appendFloat64PackedSlice
			}
			return sizeFloat64Slice, appendFloat64Slice
		}
		if nozero {
			return sizeFloat64ValueNoZero, appendFloat64ValueNoZero
		}
		return sizeFloat64Value, appendFloat64Value
	case reflect.String:
		if validateUTF8 {
			if pointer {
				return sizeStringPtr, appendUTF8StringPtr
			}
			if slice {
				return sizeStringSlice, appendUTF8StringSlice
			}
			if nozero {
				return sizeStringValueNoZero, appendUTF8StringValueNoZero
			}
			return sizeStringValue, appendUTF8StringValue
		}
		if pointer {
			return sizeStringPtr, appendStringPtr
		}
		if slice {
			return sizeStringSlice, appendStringSlice
		}
		if nozero {
			return sizeStringValueNoZero, appendStringValueNoZero
		}
		return sizeStringValue, appendStringValue
	case reflect.Slice:
		if slice {
			return sizeBytesSlice, appendBytesSlice
		}
		if oneof {
			// Oneof bytes field may also have "proto3" tag.
			// We want to marshal it as a oneof field. Do this
			// check before the proto3 check.
			return sizeBytesOneof, appendBytesOneof
		}
		if proto3 {
			return sizeBytes3, appendBytes3
		}
		return sizeBytes, appendBytes
	case reflect.Struct:
		switch encoding {
		case "group":
			if slice {
				return makeGroupSliceMarshaler(getMarshalInfo(t))
			}
			return makeGroupMarshaler(getMarshalInfo(t))
		case "bytes":
			if slice {
				return makeMessageSliceMarshaler(getMarshalInfo(t))
			}
			return makeMessageMarshaler(getMarshalInfo(t))
		}
	}
	panic(fmt.Sprintf("unknown or mismatched type: type: %v, wire type: %v", t, encoding))
}
