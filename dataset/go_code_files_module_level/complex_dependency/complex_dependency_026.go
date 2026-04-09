func NewValue(goVal interface{}) Value {

	switch val := goVal.(type) {
	case nil:
		return NilValueVal
	case Value:
		return val
	case float64:
		return NewNumberValue(val)
	case float32:
		return NewNumberValue(float64(val))
	case *float64:
		if val == nil {
			return NewNumberNil()
		}
		return NewNumberValue(*val)
	case *float32:
		if val == nil {
			return NewNumberNil()
		}
		return NewNumberValue(float64(*val))
	case int8:
		return NewIntValue(int64(val))
	case *int8:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case int16:
		return NewIntValue(int64(val))
	case *int16:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case int:
		return NewIntValue(int64(val))
	case *int:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case int32:
		return NewIntValue(int64(val))
	case *int32:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case int64:
		return NewIntValue(int64(val))
	case *int64:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case uint8:
		return NewIntValue(int64(val))
	case *uint8:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case uint32:
		return NewIntValue(int64(val))
	case *uint32:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case uint64:
		return NewIntValue(int64(val))
	case *uint64:
		if val != nil {
			return NewIntValue(int64(*val))
		}
		return NewIntValue(0)
	case string:
		// should we return Nil?
		// if val == "null" || val == "NULL" {}
		return NewStringValue(val)
	case []string:
		return NewStringsValue(val)
	// case []uint8:
	// 	return NewByteSliceValue([]byte(val))
	case []byte:
		return NewByteSliceValue(val)
	case json.RawMessage:
		return NewJsonValue(val)
	case bool:
		return NewBoolValue(val)
	case time.Time:
		return NewTimeValue(val)
	case *time.Time:
		return NewTimeValue(*val)
	case map[string]interface{}:
		return NewMapValue(val)
	case map[string]string:
		return NewMapStringValue(val)
	case map[string]float64:
		return NewMapNumberValue(val)
	case map[string]int64:
		return NewMapIntValue(val)
	case map[string]bool:
		return NewMapBoolValue(val)
	case map[string]int:
		nm := make(map[string]int64, len(val))
		for k, v := range val {
			nm[k] = int64(v)
		}
		return NewMapIntValue(nm)
	case map[string]time.Time:
		return NewMapTimeValue(val)
	case []interface{}:
		if len(val) > 0 {
			switch val[0].(type) {
			case string:
				vals := make([]string, len(val))
				for i, v := range val {
					if sv, ok := v.(string); ok {
						vals[i] = sv
					} else {
						vs := make([]Value, len(val))
						for i, v := range val {
							vs[i] = NewValue(v)
						}
						return NewSliceValues(vs)
					}
				}
				return NewStringsValue(vals)
			}
		}
		vals := make([]Value, len(val))
		for i, v := range val {
			vals[i] = NewValue(v)
		}
		return NewSliceValues(vals)
	default:
		if err, isErr := val.(error); isErr {
			return NewErrorValue(err)
		}
		return NewStructValue(val)
	}
}
