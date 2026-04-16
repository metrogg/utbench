func typeUnmarshaler(t reflect.Type, tags string) unmarshaler {
	tagArray := strings.Split(tags, ",")
	encoding := tagArray[0]
	name := "unknown"
	proto3 := false
	validateUTF8 := true
	for _, tag := range tagArray[3:] {
		if strings.HasPrefix(tag, "name=") {
			name = tag[5:]
		}
		if tag == "proto3" {
			proto3 = true
		}
	}
	validateUTF8 = validateUTF8 && proto3

	// Figure out packaging (pointer, slice, or both)
	slice := false
	pointer := false
	if t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8 {
		slice = true
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		pointer = true
		t = t.Elem()
	}

	// We'll never have both pointer and slice for basic types.
	if pointer && slice && t.Kind() != reflect.Struct {
		panic("both pointer and slice for basic type in " + t.Name())
	}

	switch t.Kind() {
	case reflect.Bool:
		if pointer {
			return unmarshalBoolPtr
		}
		if slice {
			return unmarshalBoolSlice
		}
		return unmarshalBoolValue
	case reflect.Int32:
		switch encoding {
		case "fixed32":
			if pointer {
				return unmarshalFixedS32Ptr
			}
			if slice {
				return unmarshalFixedS32Slice
			}
			return unmarshalFixedS32Value
		case "varint":
			// this could be int32 or enum
			if pointer {
				return unmarshalInt32Ptr
			}
			if slice {
				return unmarshalInt32Slice
			}
			return unmarshalInt32Value
		case "zigzag32":
			if pointer {
				return unmarshalSint32Ptr
			}
			if slice {
				return unmarshalSint32Slice
			}
			return unmarshalSint32Value
		}
	case reflect.Int64:
		switch encoding {
		case "fixed64":
			if pointer {
				return unmarshalFixedS64Ptr
			}
			if slice {
				return unmarshalFixedS64Slice
			}
			return unmarshalFixedS64Value
		case "varint":
			if pointer {
				return unmarshalInt64Ptr
			}
			if slice {
				return unmarshalInt64Slice
			}
			return unmarshalInt64Value
		case "zigzag64":
			if pointer {
				return unmarshalSint64Ptr
			}
			if slice {
				return unmarshalSint64Slice
			}
			return unmarshalSint64Value
		}
	case reflect.Uint32:
		switch encoding {
		case "fixed32":
			if pointer {
				return unmarshalFixed32Ptr
			}
			if slice {
				return unmarshalFixed32Slice
			}
			return unmarshalFixed32Value
		case "varint":
			if pointer {
				return unmarshalUint32Ptr
			}
			if slice {
				return unmarshalUint32Slice
			}
			return unmarshalUint32Value
		}
	case reflect.Uint64:
		switch encoding {
		case "fixed64":
			if pointer {
				return unmarshalFixed64Ptr
			}
			if slice {
				return unmarshalFixed64Slice
			}
			return unmarshalFixed64Value
		case "varint":
			if pointer {
				return unmarshalUint64Ptr
			}
			if slice {
				return unmarshalUint64Slice
			}
			return unmarshalUint64Value
		}
	case reflect.Float32:
		if pointer {
			return unmarshalFloat32Ptr
		}
		if slice {
			return unmarshalFloat32Slice
		}
		return unmarshalFloat32Value
	case reflect.Float64:
		if pointer {
			return unmarshalFloat64Ptr
		}
		if slice {
			return unmarshalFloat64Slice
		}
		return unmarshalFloat64Value
	case reflect.Map:
		panic("map type in typeUnmarshaler in " + t.Name())
	case reflect.Slice:
		if pointer {
			panic("bad pointer in slice case in " + t.Name())
		}
		if slice {
			return unmarshalBytesSlice
		}
		return unmarshalBytesValue
	case reflect.String:
		if validateUTF8 {
			if pointer {
				return unmarshalUTF8StringPtr
			}
			if slice {
				return unmarshalUTF8StringSlice
			}
			return unmarshalUTF8StringValue
		}
		if pointer {
			return unmarshalStringPtr
		}
		if slice {
			return unmarshalStringSlice
		}
		return unmarshalStringValue
	case reflect.Struct:
		// message or group field
		if !pointer {
			panic(fmt.Sprintf("message/group field %s:%s without pointer", t, encoding))
		}
		switch encoding {
		case "bytes":
			if slice {
				return makeUnmarshalMessageSlicePtr(getUnmarshalInfo(t), name)
			}
			return makeUnmarshalMessagePtr(getUnmarshalInfo(t), name)
		case "group":
			if slice {
				return makeUnmarshalGroupSlicePtr(getUnmarshalInfo(t), name)
			}
			return makeUnmarshalGroupPtr(getUnmarshalInfo(t), name)
		}
	}
	panic(fmt.Sprintf("unmarshaler not found type:%s encoding:%s", t, encoding))
}
