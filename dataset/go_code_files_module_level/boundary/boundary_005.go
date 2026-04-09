func fastpathDecodeTypeSwitch(iv interface{}, d *Decoder) bool {
	var changed bool
	switch v := iv.(type) {

	case []interface{}:
		var v2 []interface{}
		v2, changed = fastpathTV.DecSliceIntfV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]interface{}:
		var v2 []interface{}
		v2, changed = fastpathTV.DecSliceIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case []string:
		var v2 []string
		v2, changed = fastpathTV.DecSliceStringV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]string:
		var v2 []string
		v2, changed = fastpathTV.DecSliceStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case []float32:
		var v2 []float32
		v2, changed = fastpathTV.DecSliceFloat32V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]float32:
		var v2 []float32
		v2, changed = fastpathTV.DecSliceFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case []float64:
		var v2 []float64
		v2, changed = fastpathTV.DecSliceFloat64V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]float64:
		var v2 []float64
		v2, changed = fastpathTV.DecSliceFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint:
		var v2 []uint
		v2, changed = fastpathTV.DecSliceUintV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint:
		var v2 []uint
		v2, changed = fastpathTV.DecSliceUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint16:
		var v2 []uint16
		v2, changed = fastpathTV.DecSliceUint16V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint16:
		var v2 []uint16
		v2, changed = fastpathTV.DecSliceUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint32:
		var v2 []uint32
		v2, changed = fastpathTV.DecSliceUint32V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint32:
		var v2 []uint32
		v2, changed = fastpathTV.DecSliceUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uint64:
		var v2 []uint64
		v2, changed = fastpathTV.DecSliceUint64V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uint64:
		var v2 []uint64
		v2, changed = fastpathTV.DecSliceUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case []uintptr:
		var v2 []uintptr
		v2, changed = fastpathTV.DecSliceUintptrV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]uintptr:
		var v2 []uintptr
		v2, changed = fastpathTV.DecSliceUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case []int:
		var v2 []int
		v2, changed = fastpathTV.DecSliceIntV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int:
		var v2 []int
		v2, changed = fastpathTV.DecSliceIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case []int8:
		var v2 []int8
		v2, changed = fastpathTV.DecSliceInt8V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int8:
		var v2 []int8
		v2, changed = fastpathTV.DecSliceInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case []int16:
		var v2 []int16
		v2, changed = fastpathTV.DecSliceInt16V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int16:
		var v2 []int16
		v2, changed = fastpathTV.DecSliceInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case []int32:
		var v2 []int32
		v2, changed = fastpathTV.DecSliceInt32V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int32:
		var v2 []int32
		v2, changed = fastpathTV.DecSliceInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case []int64:
		var v2 []int64
		v2, changed = fastpathTV.DecSliceInt64V(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]int64:
		var v2 []int64
		v2, changed = fastpathTV.DecSliceInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case []bool:
		var v2 []bool
		v2, changed = fastpathTV.DecSliceBoolV(v, false, d)
		if changed && len(v) > 0 && len(v2) > 0 && !(len(v2) == len(v) && &v2[0] == &v[0]) {
			copy(v, v2)
		}
	case *[]bool:
		var v2 []bool
		v2, changed = fastpathTV.DecSliceBoolV(*v, true, d)
		if changed {
			*v = v2
		}

	case map[interface{}]interface{}:
		fastpathTV.DecMapIntfIntfV(v, false, d)
	case *map[interface{}]interface{}:
		var v2 map[interface{}]interface{}
		v2, changed = fastpathTV.DecMapIntfIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]string:
		fastpathTV.DecMapIntfStringV(v, false, d)
	case *map[interface{}]string:
		var v2 map[interface{}]string
		v2, changed = fastpathTV.DecMapIntfStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]uint:
		fastpathTV.DecMapIntfUintV(v, false, d)
	case *map[interface{}]uint:
		var v2 map[interface{}]uint
		v2, changed = fastpathTV.DecMapIntfUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]uint8:
		fastpathTV.DecMapIntfUint8V(v, false, d)
	case *map[interface{}]uint8:
		var v2 map[interface{}]uint8
		v2, changed = fastpathTV.DecMapIntfUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]uint16:
		fastpathTV.DecMapIntfUint16V(v, false, d)
	case *map[interface{}]uint16:
		var v2 map[interface{}]uint16
		v2, changed = fastpathTV.DecMapIntfUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]uint32:
		fastpathTV.DecMapIntfUint32V(v, false, d)
	case *map[interface{}]uint32:
		var v2 map[interface{}]uint32
		v2, changed = fastpathTV.DecMapIntfUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]uint64:
		fastpathTV.DecMapIntfUint64V(v, false, d)
	case *map[interface{}]uint64:
		var v2 map[interface{}]uint64
		v2, changed = fastpathTV.DecMapIntfUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]uintptr:
		fastpathTV.DecMapIntfUintptrV(v, false, d)
	case *map[interface{}]uintptr:
		var v2 map[interface{}]uintptr
		v2, changed = fastpathTV.DecMapIntfUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]int:
		fastpathTV.DecMapIntfIntV(v, false, d)
	case *map[interface{}]int:
		var v2 map[interface{}]int
		v2, changed = fastpathTV.DecMapIntfIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]int8:
		fastpathTV.DecMapIntfInt8V(v, false, d)
	case *map[interface{}]int8:
		var v2 map[interface{}]int8
		v2, changed = fastpathTV.DecMapIntfInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]int16:
		fastpathTV.DecMapIntfInt16V(v, false, d)
	case *map[interface{}]int16:
		var v2 map[interface{}]int16
		v2, changed = fastpathTV.DecMapIntfInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]int32:
		fastpathTV.DecMapIntfInt32V(v, false, d)
	case *map[interface{}]int32:
		var v2 map[interface{}]int32
		v2, changed = fastpathTV.DecMapIntfInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]int64:
		fastpathTV.DecMapIntfInt64V(v, false, d)
	case *map[interface{}]int64:
		var v2 map[interface{}]int64
		v2, changed = fastpathTV.DecMapIntfInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]float32:
		fastpathTV.DecMapIntfFloat32V(v, false, d)
	case *map[interface{}]float32:
		var v2 map[interface{}]float32
		v2, changed = fastpathTV.DecMapIntfFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]float64:
		fastpathTV.DecMapIntfFloat64V(v, false, d)
	case *map[interface{}]float64:
		var v2 map[interface{}]float64
		v2, changed = fastpathTV.DecMapIntfFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[interface{}]bool:
		fastpathTV.DecMapIntfBoolV(v, false, d)
	case *map[interface{}]bool:
		var v2 map[interface{}]bool
		v2, changed = fastpathTV.DecMapIntfBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]interface{}:
		fastpathTV.DecMapStringIntfV(v, false, d)
	case *map[string]interface{}:
		var v2 map[string]interface{}
		v2, changed = fastpathTV.DecMapStringIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]string:
		fastpathTV.DecMapStringStringV(v, false, d)
	case *map[string]string:
		var v2 map[string]string
		v2, changed = fastpathTV.DecMapStringStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint:
		fastpathTV.DecMapStringUintV(v, false, d)
	case *map[string]uint:
		var v2 map[string]uint
		v2, changed = fastpathTV.DecMapStringUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint8:
		fastpathTV.DecMapStringUint8V(v, false, d)
	case *map[string]uint8:
		var v2 map[string]uint8
		v2, changed = fastpathTV.DecMapStringUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint16:
		fastpathTV.DecMapStringUint16V(v, false, d)
	case *map[string]uint16:
		var v2 map[string]uint16
		v2, changed = fastpathTV.DecMapStringUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint32:
		fastpathTV.DecMapStringUint32V(v, false, d)
	case *map[string]uint32:
		var v2 map[string]uint32
		v2, changed = fastpathTV.DecMapStringUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uint64:
		fastpathTV.DecMapStringUint64V(v, false, d)
	case *map[string]uint64:
		var v2 map[string]uint64
		v2, changed = fastpathTV.DecMapStringUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]uintptr:
		fastpathTV.DecMapStringUintptrV(v, false, d)
	case *map[string]uintptr:
		var v2 map[string]uintptr
		v2, changed = fastpathTV.DecMapStringUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]int:
		fastpathTV.DecMapStringIntV(v, false, d)
	case *map[string]int:
		var v2 map[string]int
		v2, changed = fastpathTV.DecMapStringIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]int8:
		fastpathTV.DecMapStringInt8V(v, false, d)
	case *map[string]int8:
		var v2 map[string]int8
		v2, changed = fastpathTV.DecMapStringInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]int16:
		fastpathTV.DecMapStringInt16V(v, false, d)
	case *map[string]int16:
		var v2 map[string]int16
		v2, changed = fastpathTV.DecMapStringInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]int32:
		fastpathTV.DecMapStringInt32V(v, false, d)
	case *map[string]int32:
		var v2 map[string]int32
		v2, changed = fastpathTV.DecMapStringInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]int64:
		fastpathTV.DecMapStringInt64V(v, false, d)
	case *map[string]int64:
		var v2 map[string]int64
		v2, changed = fastpathTV.DecMapStringInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]float32:
		fastpathTV.DecMapStringFloat32V(v, false, d)
	case *map[string]float32:
		var v2 map[string]float32
		v2, changed = fastpathTV.DecMapStringFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]float64:
		fastpathTV.DecMapStringFloat64V(v, false, d)
	case *map[string]float64:
		var v2 map[string]float64
		v2, changed = fastpathTV.DecMapStringFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[string]bool:
		fastpathTV.DecMapStringBoolV(v, false, d)
	case *map[string]bool:
		var v2 map[string]bool
		v2, changed = fastpathTV.DecMapStringBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]interface{}:
		fastpathTV.DecMapFloat32IntfV(v, false, d)
	case *map[float32]interface{}:
		var v2 map[float32]interface{}
		v2, changed = fastpathTV.DecMapFloat32IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]string:
		fastpathTV.DecMapFloat32StringV(v, false, d)
	case *map[float32]string:
		var v2 map[float32]string
		v2, changed = fastpathTV.DecMapFloat32StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]uint:
		fastpathTV.DecMapFloat32UintV(v, false, d)
	case *map[float32]uint:
		var v2 map[float32]uint
		v2, changed = fastpathTV.DecMapFloat32UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]uint8:
		fastpathTV.DecMapFloat32Uint8V(v, false, d)
	case *map[float32]uint8:
		var v2 map[float32]uint8
		v2, changed = fastpathTV.DecMapFloat32Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]uint16:
		fastpathTV.DecMapFloat32Uint16V(v, false, d)
	case *map[float32]uint16:
		var v2 map[float32]uint16
		v2, changed = fastpathTV.DecMapFloat32Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]uint32:
		fastpathTV.DecMapFloat32Uint32V(v, false, d)
	case *map[float32]uint32:
		var v2 map[float32]uint32
		v2, changed = fastpathTV.DecMapFloat32Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]uint64:
		fastpathTV.DecMapFloat32Uint64V(v, false, d)
	case *map[float32]uint64:
		var v2 map[float32]uint64
		v2, changed = fastpathTV.DecMapFloat32Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]uintptr:
		fastpathTV.DecMapFloat32UintptrV(v, false, d)
	case *map[float32]uintptr:
		var v2 map[float32]uintptr
		v2, changed = fastpathTV.DecMapFloat32UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]int:
		fastpathTV.DecMapFloat32IntV(v, false, d)
	case *map[float32]int:
		var v2 map[float32]int
		v2, changed = fastpathTV.DecMapFloat32IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]int8:
		fastpathTV.DecMapFloat32Int8V(v, false, d)
	case *map[float32]int8:
		var v2 map[float32]int8
		v2, changed = fastpathTV.DecMapFloat32Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]int16:
		fastpathTV.DecMapFloat32Int16V(v, false, d)
	case *map[float32]int16:
		var v2 map[float32]int16
		v2, changed = fastpathTV.DecMapFloat32Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]int32:
		fastpathTV.DecMapFloat32Int32V(v, false, d)
	case *map[float32]int32:
		var v2 map[float32]int32
		v2, changed = fastpathTV.DecMapFloat32Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]int64:
		fastpathTV.DecMapFloat32Int64V(v, false, d)
	case *map[float32]int64:
		var v2 map[float32]int64
		v2, changed = fastpathTV.DecMapFloat32Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]float32:
		fastpathTV.DecMapFloat32Float32V(v, false, d)
	case *map[float32]float32:
		var v2 map[float32]float32
		v2, changed = fastpathTV.DecMapFloat32Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]float64:
		fastpathTV.DecMapFloat32Float64V(v, false, d)
	case *map[float32]float64:
		var v2 map[float32]float64
		v2, changed = fastpathTV.DecMapFloat32Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float32]bool:
		fastpathTV.DecMapFloat32BoolV(v, false, d)
	case *map[float32]bool:
		var v2 map[float32]bool
		v2, changed = fastpathTV.DecMapFloat32BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]interface{}:
		fastpathTV.DecMapFloat64IntfV(v, false, d)
	case *map[float64]interface{}:
		var v2 map[float64]interface{}
		v2, changed = fastpathTV.DecMapFloat64IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]string:
		fastpathTV.DecMapFloat64StringV(v, false, d)
	case *map[float64]string:
		var v2 map[float64]string
		v2, changed = fastpathTV.DecMapFloat64StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]uint:
		fastpathTV.DecMapFloat64UintV(v, false, d)
	case *map[float64]uint:
		var v2 map[float64]uint
		v2, changed = fastpathTV.DecMapFloat64UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]uint8:
		fastpathTV.DecMapFloat64Uint8V(v, false, d)
	case *map[float64]uint8:
		var v2 map[float64]uint8
		v2, changed = fastpathTV.DecMapFloat64Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]uint16:
		fastpathTV.DecMapFloat64Uint16V(v, false, d)
	case *map[float64]uint16:
		var v2 map[float64]uint16
		v2, changed = fastpathTV.DecMapFloat64Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]uint32:
		fastpathTV.DecMapFloat64Uint32V(v, false, d)
	case *map[float64]uint32:
		var v2 map[float64]uint32
		v2, changed = fastpathTV.DecMapFloat64Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]uint64:
		fastpathTV.DecMapFloat64Uint64V(v, false, d)
	case *map[float64]uint64:
		var v2 map[float64]uint64
		v2, changed = fastpathTV.DecMapFloat64Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]uintptr:
		fastpathTV.DecMapFloat64UintptrV(v, false, d)
	case *map[float64]uintptr:
		var v2 map[float64]uintptr
		v2, changed = fastpathTV.DecMapFloat64UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]int:
		fastpathTV.DecMapFloat64IntV(v, false, d)
	case *map[float64]int:
		var v2 map[float64]int
		v2, changed = fastpathTV.DecMapFloat64IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]int8:
		fastpathTV.DecMapFloat64Int8V(v, false, d)
	case *map[float64]int8:
		var v2 map[float64]int8
		v2, changed = fastpathTV.DecMapFloat64Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]int16:
		fastpathTV.DecMapFloat64Int16V(v, false, d)
	case *map[float64]int16:
		var v2 map[float64]int16
		v2, changed = fastpathTV.DecMapFloat64Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]int32:
		fastpathTV.DecMapFloat64Int32V(v, false, d)
	case *map[float64]int32:
		var v2 map[float64]int32
		v2, changed = fastpathTV.DecMapFloat64Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]int64:
		fastpathTV.DecMapFloat64Int64V(v, false, d)
	case *map[float64]int64:
		var v2 map[float64]int64
		v2, changed = fastpathTV.DecMapFloat64Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]float32:
		fastpathTV.DecMapFloat64Float32V(v, false, d)
	case *map[float64]float32:
		var v2 map[float64]float32
		v2, changed = fastpathTV.DecMapFloat64Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]float64:
		fastpathTV.DecMapFloat64Float64V(v, false, d)
	case *map[float64]float64:
		var v2 map[float64]float64
		v2, changed = fastpathTV.DecMapFloat64Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[float64]bool:
		fastpathTV.DecMapFloat64BoolV(v, false, d)
	case *map[float64]bool:
		var v2 map[float64]bool
		v2, changed = fastpathTV.DecMapFloat64BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]interface{}:
		fastpathTV.DecMapUintIntfV(v, false, d)
	case *map[uint]interface{}:
		var v2 map[uint]interface{}
		v2, changed = fastpathTV.DecMapUintIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]string:
		fastpathTV.DecMapUintStringV(v, false, d)
	case *map[uint]string:
		var v2 map[uint]string
		v2, changed = fastpathTV.DecMapUintStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint:
		fastpathTV.DecMapUintUintV(v, false, d)
	case *map[uint]uint:
		var v2 map[uint]uint
		v2, changed = fastpathTV.DecMapUintUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint8:
		fastpathTV.DecMapUintUint8V(v, false, d)
	case *map[uint]uint8:
		var v2 map[uint]uint8
		v2, changed = fastpathTV.DecMapUintUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint16:
		fastpathTV.DecMapUintUint16V(v, false, d)
	case *map[uint]uint16:
		var v2 map[uint]uint16
		v2, changed = fastpathTV.DecMapUintUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint32:
		fastpathTV.DecMapUintUint32V(v, false, d)
	case *map[uint]uint32:
		var v2 map[uint]uint32
		v2, changed = fastpathTV.DecMapUintUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uint64:
		fastpathTV.DecMapUintUint64V(v, false, d)
	case *map[uint]uint64:
		var v2 map[uint]uint64
		v2, changed = fastpathTV.DecMapUintUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]uintptr:
		fastpathTV.DecMapUintUintptrV(v, false, d)
	case *map[uint]uintptr:
		var v2 map[uint]uintptr
		v2, changed = fastpathTV.DecMapUintUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]int:
		fastpathTV.DecMapUintIntV(v, false, d)
	case *map[uint]int:
		var v2 map[uint]int
		v2, changed = fastpathTV.DecMapUintIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]int8:
		fastpathTV.DecMapUintInt8V(v, false, d)
	case *map[uint]int8:
		var v2 map[uint]int8
		v2, changed = fastpathTV.DecMapUintInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]int16:
		fastpathTV.DecMapUintInt16V(v, false, d)
	case *map[uint]int16:
		var v2 map[uint]int16
		v2, changed = fastpathTV.DecMapUintInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]int32:
		fastpathTV.DecMapUintInt32V(v, false, d)
	case *map[uint]int32:
		var v2 map[uint]int32
		v2, changed = fastpathTV.DecMapUintInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]int64:
		fastpathTV.DecMapUintInt64V(v, false, d)
	case *map[uint]int64:
		var v2 map[uint]int64
		v2, changed = fastpathTV.DecMapUintInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]float32:
		fastpathTV.DecMapUintFloat32V(v, false, d)
	case *map[uint]float32:
		var v2 map[uint]float32
		v2, changed = fastpathTV.DecMapUintFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]float64:
		fastpathTV.DecMapUintFloat64V(v, false, d)
	case *map[uint]float64:
		var v2 map[uint]float64
		v2, changed = fastpathTV.DecMapUintFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint]bool:
		fastpathTV.DecMapUintBoolV(v, false, d)
	case *map[uint]bool:
		var v2 map[uint]bool
		v2, changed = fastpathTV.DecMapUintBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]interface{}:
		fastpathTV.DecMapUint8IntfV(v, false, d)
	case *map[uint8]interface{}:
		var v2 map[uint8]interface{}
		v2, changed = fastpathTV.DecMapUint8IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]string:
		fastpathTV.DecMapUint8StringV(v, false, d)
	case *map[uint8]string:
		var v2 map[uint8]string
		v2, changed = fastpathTV.DecMapUint8StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint:
		fastpathTV.DecMapUint8UintV(v, false, d)
	case *map[uint8]uint:
		var v2 map[uint8]uint
		v2, changed = fastpathTV.DecMapUint8UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint8:
		fastpathTV.DecMapUint8Uint8V(v, false, d)
	case *map[uint8]uint8:
		var v2 map[uint8]uint8
		v2, changed = fastpathTV.DecMapUint8Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint16:
		fastpathTV.DecMapUint8Uint16V(v, false, d)
	case *map[uint8]uint16:
		var v2 map[uint8]uint16
		v2, changed = fastpathTV.DecMapUint8Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint32:
		fastpathTV.DecMapUint8Uint32V(v, false, d)
	case *map[uint8]uint32:
		var v2 map[uint8]uint32
		v2, changed = fastpathTV.DecMapUint8Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uint64:
		fastpathTV.DecMapUint8Uint64V(v, false, d)
	case *map[uint8]uint64:
		var v2 map[uint8]uint64
		v2, changed = fastpathTV.DecMapUint8Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]uintptr:
		fastpathTV.DecMapUint8UintptrV(v, false, d)
	case *map[uint8]uintptr:
		var v2 map[uint8]uintptr
		v2, changed = fastpathTV.DecMapUint8UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]int:
		fastpathTV.DecMapUint8IntV(v, false, d)
	case *map[uint8]int:
		var v2 map[uint8]int
		v2, changed = fastpathTV.DecMapUint8IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]int8:
		fastpathTV.DecMapUint8Int8V(v, false, d)
	case *map[uint8]int8:
		var v2 map[uint8]int8
		v2, changed = fastpathTV.DecMapUint8Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]int16:
		fastpathTV.DecMapUint8Int16V(v, false, d)
	case *map[uint8]int16:
		var v2 map[uint8]int16
		v2, changed = fastpathTV.DecMapUint8Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]int32:
		fastpathTV.DecMapUint8Int32V(v, false, d)
	case *map[uint8]int32:
		var v2 map[uint8]int32
		v2, changed = fastpathTV.DecMapUint8Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]int64:
		fastpathTV.DecMapUint8Int64V(v, false, d)
	case *map[uint8]int64:
		var v2 map[uint8]int64
		v2, changed = fastpathTV.DecMapUint8Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]float32:
		fastpathTV.DecMapUint8Float32V(v, false, d)
	case *map[uint8]float32:
		var v2 map[uint8]float32
		v2, changed = fastpathTV.DecMapUint8Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]float64:
		fastpathTV.DecMapUint8Float64V(v, false, d)
	case *map[uint8]float64:
		var v2 map[uint8]float64
		v2, changed = fastpathTV.DecMapUint8Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint8]bool:
		fastpathTV.DecMapUint8BoolV(v, false, d)
	case *map[uint8]bool:
		var v2 map[uint8]bool
		v2, changed = fastpathTV.DecMapUint8BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]interface{}:
		fastpathTV.DecMapUint16IntfV(v, false, d)
	case *map[uint16]interface{}:
		var v2 map[uint16]interface{}
		v2, changed = fastpathTV.DecMapUint16IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]string:
		fastpathTV.DecMapUint16StringV(v, false, d)
	case *map[uint16]string:
		var v2 map[uint16]string
		v2, changed = fastpathTV.DecMapUint16StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]uint:
		fastpathTV.DecMapUint16UintV(v, false, d)
	case *map[uint16]uint:
		var v2 map[uint16]uint
		v2, changed = fastpathTV.DecMapUint16UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]uint8:
		fastpathTV.DecMapUint16Uint8V(v, false, d)
	case *map[uint16]uint8:
		var v2 map[uint16]uint8
		v2, changed = fastpathTV.DecMapUint16Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]uint16:
		fastpathTV.DecMapUint16Uint16V(v, false, d)
	case *map[uint16]uint16:
		var v2 map[uint16]uint16
		v2, changed = fastpathTV.DecMapUint16Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]uint32:
		fastpathTV.DecMapUint16Uint32V(v, false, d)
	case *map[uint16]uint32:
		var v2 map[uint16]uint32
		v2, changed = fastpathTV.DecMapUint16Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]uint64:
		fastpathTV.DecMapUint16Uint64V(v, false, d)
	case *map[uint16]uint64:
		var v2 map[uint16]uint64
		v2, changed = fastpathTV.DecMapUint16Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]uintptr:
		fastpathTV.DecMapUint16UintptrV(v, false, d)
	case *map[uint16]uintptr:
		var v2 map[uint16]uintptr
		v2, changed = fastpathTV.DecMapUint16UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]int:
		fastpathTV.DecMapUint16IntV(v, false, d)
	case *map[uint16]int:
		var v2 map[uint16]int
		v2, changed = fastpathTV.DecMapUint16IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]int8:
		fastpathTV.DecMapUint16Int8V(v, false, d)
	case *map[uint16]int8:
		var v2 map[uint16]int8
		v2, changed = fastpathTV.DecMapUint16Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]int16:
		fastpathTV.DecMapUint16Int16V(v, false, d)
	case *map[uint16]int16:
		var v2 map[uint16]int16
		v2, changed = fastpathTV.DecMapUint16Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]int32:
		fastpathTV.DecMapUint16Int32V(v, false, d)
	case *map[uint16]int32:
		var v2 map[uint16]int32
		v2, changed = fastpathTV.DecMapUint16Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]int64:
		fastpathTV.DecMapUint16Int64V(v, false, d)
	case *map[uint16]int64:
		var v2 map[uint16]int64
		v2, changed = fastpathTV.DecMapUint16Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]float32:
		fastpathTV.DecMapUint16Float32V(v, false, d)
	case *map[uint16]float32:
		var v2 map[uint16]float32
		v2, changed = fastpathTV.DecMapUint16Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]float64:
		fastpathTV.DecMapUint16Float64V(v, false, d)
	case *map[uint16]float64:
		var v2 map[uint16]float64
		v2, changed = fastpathTV.DecMapUint16Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint16]bool:
		fastpathTV.DecMapUint16BoolV(v, false, d)
	case *map[uint16]bool:
		var v2 map[uint16]bool
		v2, changed = fastpathTV.DecMapUint16BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]interface{}:
		fastpathTV.DecMapUint32IntfV(v, false, d)
	case *map[uint32]interface{}:
		var v2 map[uint32]interface{}
		v2, changed = fastpathTV.DecMapUint32IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]string:
		fastpathTV.DecMapUint32StringV(v, false, d)
	case *map[uint32]string:
		var v2 map[uint32]string
		v2, changed = fastpathTV.DecMapUint32StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]uint:
		fastpathTV.DecMapUint32UintV(v, false, d)
	case *map[uint32]uint:
		var v2 map[uint32]uint
		v2, changed = fastpathTV.DecMapUint32UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]uint8:
		fastpathTV.DecMapUint32Uint8V(v, false, d)
	case *map[uint32]uint8:
		var v2 map[uint32]uint8
		v2, changed = fastpathTV.DecMapUint32Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]uint16:
		fastpathTV.DecMapUint32Uint16V(v, false, d)
	case *map[uint32]uint16:
		var v2 map[uint32]uint16
		v2, changed = fastpathTV.DecMapUint32Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]uint32:
		fastpathTV.DecMapUint32Uint32V(v, false, d)
	case *map[uint32]uint32:
		var v2 map[uint32]uint32
		v2, changed = fastpathTV.DecMapUint32Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]uint64:
		fastpathTV.DecMapUint32Uint64V(v, false, d)
	case *map[uint32]uint64:
		var v2 map[uint32]uint64
		v2, changed = fastpathTV.DecMapUint32Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]uintptr:
		fastpathTV.DecMapUint32UintptrV(v, false, d)
	case *map[uint32]uintptr:
		var v2 map[uint32]uintptr
		v2, changed = fastpathTV.DecMapUint32UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]int:
		fastpathTV.DecMapUint32IntV(v, false, d)
	case *map[uint32]int:
		var v2 map[uint32]int
		v2, changed = fastpathTV.DecMapUint32IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]int8:
		fastpathTV.DecMapUint32Int8V(v, false, d)
	case *map[uint32]int8:
		var v2 map[uint32]int8
		v2, changed = fastpathTV.DecMapUint32Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]int16:
		fastpathTV.DecMapUint32Int16V(v, false, d)
	case *map[uint32]int16:
		var v2 map[uint32]int16
		v2, changed = fastpathTV.DecMapUint32Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]int32:
		fastpathTV.DecMapUint32Int32V(v, false, d)
	case *map[uint32]int32:
		var v2 map[uint32]int32
		v2, changed = fastpathTV.DecMapUint32Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]int64:
		fastpathTV.DecMapUint32Int64V(v, false, d)
	case *map[uint32]int64:
		var v2 map[uint32]int64
		v2, changed = fastpathTV.DecMapUint32Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]float32:
		fastpathTV.DecMapUint32Float32V(v, false, d)
	case *map[uint32]float32:
		var v2 map[uint32]float32
		v2, changed = fastpathTV.DecMapUint32Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]float64:
		fastpathTV.DecMapUint32Float64V(v, false, d)
	case *map[uint32]float64:
		var v2 map[uint32]float64
		v2, changed = fastpathTV.DecMapUint32Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint32]bool:
		fastpathTV.DecMapUint32BoolV(v, false, d)
	case *map[uint32]bool:
		var v2 map[uint32]bool
		v2, changed = fastpathTV.DecMapUint32BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]interface{}:
		fastpathTV.DecMapUint64IntfV(v, false, d)
	case *map[uint64]interface{}:
		var v2 map[uint64]interface{}
		v2, changed = fastpathTV.DecMapUint64IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]string:
		fastpathTV.DecMapUint64StringV(v, false, d)
	case *map[uint64]string:
		var v2 map[uint64]string
		v2, changed = fastpathTV.DecMapUint64StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint:
		fastpathTV.DecMapUint64UintV(v, false, d)
	case *map[uint64]uint:
		var v2 map[uint64]uint
		v2, changed = fastpathTV.DecMapUint64UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint8:
		fastpathTV.DecMapUint64Uint8V(v, false, d)
	case *map[uint64]uint8:
		var v2 map[uint64]uint8
		v2, changed = fastpathTV.DecMapUint64Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint16:
		fastpathTV.DecMapUint64Uint16V(v, false, d)
	case *map[uint64]uint16:
		var v2 map[uint64]uint16
		v2, changed = fastpathTV.DecMapUint64Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint32:
		fastpathTV.DecMapUint64Uint32V(v, false, d)
	case *map[uint64]uint32:
		var v2 map[uint64]uint32
		v2, changed = fastpathTV.DecMapUint64Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uint64:
		fastpathTV.DecMapUint64Uint64V(v, false, d)
	case *map[uint64]uint64:
		var v2 map[uint64]uint64
		v2, changed = fastpathTV.DecMapUint64Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]uintptr:
		fastpathTV.DecMapUint64UintptrV(v, false, d)
	case *map[uint64]uintptr:
		var v2 map[uint64]uintptr
		v2, changed = fastpathTV.DecMapUint64UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]int:
		fastpathTV.DecMapUint64IntV(v, false, d)
	case *map[uint64]int:
		var v2 map[uint64]int
		v2, changed = fastpathTV.DecMapUint64IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]int8:
		fastpathTV.DecMapUint64Int8V(v, false, d)
	case *map[uint64]int8:
		var v2 map[uint64]int8
		v2, changed = fastpathTV.DecMapUint64Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]int16:
		fastpathTV.DecMapUint64Int16V(v, false, d)
	case *map[uint64]int16:
		var v2 map[uint64]int16
		v2, changed = fastpathTV.DecMapUint64Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]int32:
		fastpathTV.DecMapUint64Int32V(v, false, d)
	case *map[uint64]int32:
		var v2 map[uint64]int32
		v2, changed = fastpathTV.DecMapUint64Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]int64:
		fastpathTV.DecMapUint64Int64V(v, false, d)
	case *map[uint64]int64:
		var v2 map[uint64]int64
		v2, changed = fastpathTV.DecMapUint64Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]float32:
		fastpathTV.DecMapUint64Float32V(v, false, d)
	case *map[uint64]float32:
		var v2 map[uint64]float32
		v2, changed = fastpathTV.DecMapUint64Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]float64:
		fastpathTV.DecMapUint64Float64V(v, false, d)
	case *map[uint64]float64:
		var v2 map[uint64]float64
		v2, changed = fastpathTV.DecMapUint64Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uint64]bool:
		fastpathTV.DecMapUint64BoolV(v, false, d)
	case *map[uint64]bool:
		var v2 map[uint64]bool
		v2, changed = fastpathTV.DecMapUint64BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]interface{}:
		fastpathTV.DecMapUintptrIntfV(v, false, d)
	case *map[uintptr]interface{}:
		var v2 map[uintptr]interface{}
		v2, changed = fastpathTV.DecMapUintptrIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]string:
		fastpathTV.DecMapUintptrStringV(v, false, d)
	case *map[uintptr]string:
		var v2 map[uintptr]string
		v2, changed = fastpathTV.DecMapUintptrStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]uint:
		fastpathTV.DecMapUintptrUintV(v, false, d)
	case *map[uintptr]uint:
		var v2 map[uintptr]uint
		v2, changed = fastpathTV.DecMapUintptrUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]uint8:
		fastpathTV.DecMapUintptrUint8V(v, false, d)
	case *map[uintptr]uint8:
		var v2 map[uintptr]uint8
		v2, changed = fastpathTV.DecMapUintptrUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]uint16:
		fastpathTV.DecMapUintptrUint16V(v, false, d)
	case *map[uintptr]uint16:
		var v2 map[uintptr]uint16
		v2, changed = fastpathTV.DecMapUintptrUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]uint32:
		fastpathTV.DecMapUintptrUint32V(v, false, d)
	case *map[uintptr]uint32:
		var v2 map[uintptr]uint32
		v2, changed = fastpathTV.DecMapUintptrUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]uint64:
		fastpathTV.DecMapUintptrUint64V(v, false, d)
	case *map[uintptr]uint64:
		var v2 map[uintptr]uint64
		v2, changed = fastpathTV.DecMapUintptrUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]uintptr:
		fastpathTV.DecMapUintptrUintptrV(v, false, d)
	case *map[uintptr]uintptr:
		var v2 map[uintptr]uintptr
		v2, changed = fastpathTV.DecMapUintptrUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]int:
		fastpathTV.DecMapUintptrIntV(v, false, d)
	case *map[uintptr]int:
		var v2 map[uintptr]int
		v2, changed = fastpathTV.DecMapUintptrIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]int8:
		fastpathTV.DecMapUintptrInt8V(v, false, d)
	case *map[uintptr]int8:
		var v2 map[uintptr]int8
		v2, changed = fastpathTV.DecMapUintptrInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]int16:
		fastpathTV.DecMapUintptrInt16V(v, false, d)
	case *map[uintptr]int16:
		var v2 map[uintptr]int16
		v2, changed = fastpathTV.DecMapUintptrInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]int32:
		fastpathTV.DecMapUintptrInt32V(v, false, d)
	case *map[uintptr]int32:
		var v2 map[uintptr]int32
		v2, changed = fastpathTV.DecMapUintptrInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]int64:
		fastpathTV.DecMapUintptrInt64V(v, false, d)
	case *map[uintptr]int64:
		var v2 map[uintptr]int64
		v2, changed = fastpathTV.DecMapUintptrInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]float32:
		fastpathTV.DecMapUintptrFloat32V(v, false, d)
	case *map[uintptr]float32:
		var v2 map[uintptr]float32
		v2, changed = fastpathTV.DecMapUintptrFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]float64:
		fastpathTV.DecMapUintptrFloat64V(v, false, d)
	case *map[uintptr]float64:
		var v2 map[uintptr]float64
		v2, changed = fastpathTV.DecMapUintptrFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[uintptr]bool:
		fastpathTV.DecMapUintptrBoolV(v, false, d)
	case *map[uintptr]bool:
		var v2 map[uintptr]bool
		v2, changed = fastpathTV.DecMapUintptrBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]interface{}:
		fastpathTV.DecMapIntIntfV(v, false, d)
	case *map[int]interface{}:
		var v2 map[int]interface{}
		v2, changed = fastpathTV.DecMapIntIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]string:
		fastpathTV.DecMapIntStringV(v, false, d)
	case *map[int]string:
		var v2 map[int]string
		v2, changed = fastpathTV.DecMapIntStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint:
		fastpathTV.DecMapIntUintV(v, false, d)
	case *map[int]uint:
		var v2 map[int]uint
		v2, changed = fastpathTV.DecMapIntUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint8:
		fastpathTV.DecMapIntUint8V(v, false, d)
	case *map[int]uint8:
		var v2 map[int]uint8
		v2, changed = fastpathTV.DecMapIntUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint16:
		fastpathTV.DecMapIntUint16V(v, false, d)
	case *map[int]uint16:
		var v2 map[int]uint16
		v2, changed = fastpathTV.DecMapIntUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint32:
		fastpathTV.DecMapIntUint32V(v, false, d)
	case *map[int]uint32:
		var v2 map[int]uint32
		v2, changed = fastpathTV.DecMapIntUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uint64:
		fastpathTV.DecMapIntUint64V(v, false, d)
	case *map[int]uint64:
		var v2 map[int]uint64
		v2, changed = fastpathTV.DecMapIntUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]uintptr:
		fastpathTV.DecMapIntUintptrV(v, false, d)
	case *map[int]uintptr:
		var v2 map[int]uintptr
		v2, changed = fastpathTV.DecMapIntUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]int:
		fastpathTV.DecMapIntIntV(v, false, d)
	case *map[int]int:
		var v2 map[int]int
		v2, changed = fastpathTV.DecMapIntIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]int8:
		fastpathTV.DecMapIntInt8V(v, false, d)
	case *map[int]int8:
		var v2 map[int]int8
		v2, changed = fastpathTV.DecMapIntInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]int16:
		fastpathTV.DecMapIntInt16V(v, false, d)
	case *map[int]int16:
		var v2 map[int]int16
		v2, changed = fastpathTV.DecMapIntInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]int32:
		fastpathTV.DecMapIntInt32V(v, false, d)
	case *map[int]int32:
		var v2 map[int]int32
		v2, changed = fastpathTV.DecMapIntInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]int64:
		fastpathTV.DecMapIntInt64V(v, false, d)
	case *map[int]int64:
		var v2 map[int]int64
		v2, changed = fastpathTV.DecMapIntInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]float32:
		fastpathTV.DecMapIntFloat32V(v, false, d)
	case *map[int]float32:
		var v2 map[int]float32
		v2, changed = fastpathTV.DecMapIntFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]float64:
		fastpathTV.DecMapIntFloat64V(v, false, d)
	case *map[int]float64:
		var v2 map[int]float64
		v2, changed = fastpathTV.DecMapIntFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int]bool:
		fastpathTV.DecMapIntBoolV(v, false, d)
	case *map[int]bool:
		var v2 map[int]bool
		v2, changed = fastpathTV.DecMapIntBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]interface{}:
		fastpathTV.DecMapInt8IntfV(v, false, d)
	case *map[int8]interface{}:
		var v2 map[int8]interface{}
		v2, changed = fastpathTV.DecMapInt8IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]string:
		fastpathTV.DecMapInt8StringV(v, false, d)
	case *map[int8]string:
		var v2 map[int8]string
		v2, changed = fastpathTV.DecMapInt8StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]uint:
		fastpathTV.DecMapInt8UintV(v, false, d)
	case *map[int8]uint:
		var v2 map[int8]uint
		v2, changed = fastpathTV.DecMapInt8UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]uint8:
		fastpathTV.DecMapInt8Uint8V(v, false, d)
	case *map[int8]uint8:
		var v2 map[int8]uint8
		v2, changed = fastpathTV.DecMapInt8Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]uint16:
		fastpathTV.DecMapInt8Uint16V(v, false, d)
	case *map[int8]uint16:
		var v2 map[int8]uint16
		v2, changed = fastpathTV.DecMapInt8Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]uint32:
		fastpathTV.DecMapInt8Uint32V(v, false, d)
	case *map[int8]uint32:
		var v2 map[int8]uint32
		v2, changed = fastpathTV.DecMapInt8Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]uint64:
		fastpathTV.DecMapInt8Uint64V(v, false, d)
	case *map[int8]uint64:
		var v2 map[int8]uint64
		v2, changed = fastpathTV.DecMapInt8Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]uintptr:
		fastpathTV.DecMapInt8UintptrV(v, false, d)
	case *map[int8]uintptr:
		var v2 map[int8]uintptr
		v2, changed = fastpathTV.DecMapInt8UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]int:
		fastpathTV.DecMapInt8IntV(v, false, d)
	case *map[int8]int:
		var v2 map[int8]int
		v2, changed = fastpathTV.DecMapInt8IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]int8:
		fastpathTV.DecMapInt8Int8V(v, false, d)
	case *map[int8]int8:
		var v2 map[int8]int8
		v2, changed = fastpathTV.DecMapInt8Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]int16:
		fastpathTV.DecMapInt8Int16V(v, false, d)
	case *map[int8]int16:
		var v2 map[int8]int16
		v2, changed = fastpathTV.DecMapInt8Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]int32:
		fastpathTV.DecMapInt8Int32V(v, false, d)
	case *map[int8]int32:
		var v2 map[int8]int32
		v2, changed = fastpathTV.DecMapInt8Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]int64:
		fastpathTV.DecMapInt8Int64V(v, false, d)
	case *map[int8]int64:
		var v2 map[int8]int64
		v2, changed = fastpathTV.DecMapInt8Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]float32:
		fastpathTV.DecMapInt8Float32V(v, false, d)
	case *map[int8]float32:
		var v2 map[int8]float32
		v2, changed = fastpathTV.DecMapInt8Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]float64:
		fastpathTV.DecMapInt8Float64V(v, false, d)
	case *map[int8]float64:
		var v2 map[int8]float64
		v2, changed = fastpathTV.DecMapInt8Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int8]bool:
		fastpathTV.DecMapInt8BoolV(v, false, d)
	case *map[int8]bool:
		var v2 map[int8]bool
		v2, changed = fastpathTV.DecMapInt8BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]interface{}:
		fastpathTV.DecMapInt16IntfV(v, false, d)
	case *map[int16]interface{}:
		var v2 map[int16]interface{}
		v2, changed = fastpathTV.DecMapInt16IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]string:
		fastpathTV.DecMapInt16StringV(v, false, d)
	case *map[int16]string:
		var v2 map[int16]string
		v2, changed = fastpathTV.DecMapInt16StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]uint:
		fastpathTV.DecMapInt16UintV(v, false, d)
	case *map[int16]uint:
		var v2 map[int16]uint
		v2, changed = fastpathTV.DecMapInt16UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]uint8:
		fastpathTV.DecMapInt16Uint8V(v, false, d)
	case *map[int16]uint8:
		var v2 map[int16]uint8
		v2, changed = fastpathTV.DecMapInt16Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]uint16:
		fastpathTV.DecMapInt16Uint16V(v, false, d)
	case *map[int16]uint16:
		var v2 map[int16]uint16
		v2, changed = fastpathTV.DecMapInt16Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]uint32:
		fastpathTV.DecMapInt16Uint32V(v, false, d)
	case *map[int16]uint32:
		var v2 map[int16]uint32
		v2, changed = fastpathTV.DecMapInt16Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]uint64:
		fastpathTV.DecMapInt16Uint64V(v, false, d)
	case *map[int16]uint64:
		var v2 map[int16]uint64
		v2, changed = fastpathTV.DecMapInt16Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]uintptr:
		fastpathTV.DecMapInt16UintptrV(v, false, d)
	case *map[int16]uintptr:
		var v2 map[int16]uintptr
		v2, changed = fastpathTV.DecMapInt16UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]int:
		fastpathTV.DecMapInt16IntV(v, false, d)
	case *map[int16]int:
		var v2 map[int16]int
		v2, changed = fastpathTV.DecMapInt16IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]int8:
		fastpathTV.DecMapInt16Int8V(v, false, d)
	case *map[int16]int8:
		var v2 map[int16]int8
		v2, changed = fastpathTV.DecMapInt16Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]int16:
		fastpathTV.DecMapInt16Int16V(v, false, d)
	case *map[int16]int16:
		var v2 map[int16]int16
		v2, changed = fastpathTV.DecMapInt16Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]int32:
		fastpathTV.DecMapInt16Int32V(v, false, d)
	case *map[int16]int32:
		var v2 map[int16]int32
		v2, changed = fastpathTV.DecMapInt16Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]int64:
		fastpathTV.DecMapInt16Int64V(v, false, d)
	case *map[int16]int64:
		var v2 map[int16]int64
		v2, changed = fastpathTV.DecMapInt16Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]float32:
		fastpathTV.DecMapInt16Float32V(v, false, d)
	case *map[int16]float32:
		var v2 map[int16]float32
		v2, changed = fastpathTV.DecMapInt16Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]float64:
		fastpathTV.DecMapInt16Float64V(v, false, d)
	case *map[int16]float64:
		var v2 map[int16]float64
		v2, changed = fastpathTV.DecMapInt16Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int16]bool:
		fastpathTV.DecMapInt16BoolV(v, false, d)
	case *map[int16]bool:
		var v2 map[int16]bool
		v2, changed = fastpathTV.DecMapInt16BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]interface{}:
		fastpathTV.DecMapInt32IntfV(v, false, d)
	case *map[int32]interface{}:
		var v2 map[int32]interface{}
		v2, changed = fastpathTV.DecMapInt32IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]string:
		fastpathTV.DecMapInt32StringV(v, false, d)
	case *map[int32]string:
		var v2 map[int32]string
		v2, changed = fastpathTV.DecMapInt32StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]uint:
		fastpathTV.DecMapInt32UintV(v, false, d)
	case *map[int32]uint:
		var v2 map[int32]uint
		v2, changed = fastpathTV.DecMapInt32UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]uint8:
		fastpathTV.DecMapInt32Uint8V(v, false, d)
	case *map[int32]uint8:
		var v2 map[int32]uint8
		v2, changed = fastpathTV.DecMapInt32Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]uint16:
		fastpathTV.DecMapInt32Uint16V(v, false, d)
	case *map[int32]uint16:
		var v2 map[int32]uint16
		v2, changed = fastpathTV.DecMapInt32Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]uint32:
		fastpathTV.DecMapInt32Uint32V(v, false, d)
	case *map[int32]uint32:
		var v2 map[int32]uint32
		v2, changed = fastpathTV.DecMapInt32Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]uint64:
		fastpathTV.DecMapInt32Uint64V(v, false, d)
	case *map[int32]uint64:
		var v2 map[int32]uint64
		v2, changed = fastpathTV.DecMapInt32Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]uintptr:
		fastpathTV.DecMapInt32UintptrV(v, false, d)
	case *map[int32]uintptr:
		var v2 map[int32]uintptr
		v2, changed = fastpathTV.DecMapInt32UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]int:
		fastpathTV.DecMapInt32IntV(v, false, d)
	case *map[int32]int:
		var v2 map[int32]int
		v2, changed = fastpathTV.DecMapInt32IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]int8:
		fastpathTV.DecMapInt32Int8V(v, false, d)
	case *map[int32]int8:
		var v2 map[int32]int8
		v2, changed = fastpathTV.DecMapInt32Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]int16:
		fastpathTV.DecMapInt32Int16V(v, false, d)
	case *map[int32]int16:
		var v2 map[int32]int16
		v2, changed = fastpathTV.DecMapInt32Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]int32:
		fastpathTV.DecMapInt32Int32V(v, false, d)
	case *map[int32]int32:
		var v2 map[int32]int32
		v2, changed = fastpathTV.DecMapInt32Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]int64:
		fastpathTV.DecMapInt32Int64V(v, false, d)
	case *map[int32]int64:
		var v2 map[int32]int64
		v2, changed = fastpathTV.DecMapInt32Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]float32:
		fastpathTV.DecMapInt32Float32V(v, false, d)
	case *map[int32]float32:
		var v2 map[int32]float32
		v2, changed = fastpathTV.DecMapInt32Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]float64:
		fastpathTV.DecMapInt32Float64V(v, false, d)
	case *map[int32]float64:
		var v2 map[int32]float64
		v2, changed = fastpathTV.DecMapInt32Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int32]bool:
		fastpathTV.DecMapInt32BoolV(v, false, d)
	case *map[int32]bool:
		var v2 map[int32]bool
		v2, changed = fastpathTV.DecMapInt32BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]interface{}:
		fastpathTV.DecMapInt64IntfV(v, false, d)
	case *map[int64]interface{}:
		var v2 map[int64]interface{}
		v2, changed = fastpathTV.DecMapInt64IntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]string:
		fastpathTV.DecMapInt64StringV(v, false, d)
	case *map[int64]string:
		var v2 map[int64]string
		v2, changed = fastpathTV.DecMapInt64StringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint:
		fastpathTV.DecMapInt64UintV(v, false, d)
	case *map[int64]uint:
		var v2 map[int64]uint
		v2, changed = fastpathTV.DecMapInt64UintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint8:
		fastpathTV.DecMapInt64Uint8V(v, false, d)
	case *map[int64]uint8:
		var v2 map[int64]uint8
		v2, changed = fastpathTV.DecMapInt64Uint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint16:
		fastpathTV.DecMapInt64Uint16V(v, false, d)
	case *map[int64]uint16:
		var v2 map[int64]uint16
		v2, changed = fastpathTV.DecMapInt64Uint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint32:
		fastpathTV.DecMapInt64Uint32V(v, false, d)
	case *map[int64]uint32:
		var v2 map[int64]uint32
		v2, changed = fastpathTV.DecMapInt64Uint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uint64:
		fastpathTV.DecMapInt64Uint64V(v, false, d)
	case *map[int64]uint64:
		var v2 map[int64]uint64
		v2, changed = fastpathTV.DecMapInt64Uint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]uintptr:
		fastpathTV.DecMapInt64UintptrV(v, false, d)
	case *map[int64]uintptr:
		var v2 map[int64]uintptr
		v2, changed = fastpathTV.DecMapInt64UintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]int:
		fastpathTV.DecMapInt64IntV(v, false, d)
	case *map[int64]int:
		var v2 map[int64]int
		v2, changed = fastpathTV.DecMapInt64IntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]int8:
		fastpathTV.DecMapInt64Int8V(v, false, d)
	case *map[int64]int8:
		var v2 map[int64]int8
		v2, changed = fastpathTV.DecMapInt64Int8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]int16:
		fastpathTV.DecMapInt64Int16V(v, false, d)
	case *map[int64]int16:
		var v2 map[int64]int16
		v2, changed = fastpathTV.DecMapInt64Int16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]int32:
		fastpathTV.DecMapInt64Int32V(v, false, d)
	case *map[int64]int32:
		var v2 map[int64]int32
		v2, changed = fastpathTV.DecMapInt64Int32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]int64:
		fastpathTV.DecMapInt64Int64V(v, false, d)
	case *map[int64]int64:
		var v2 map[int64]int64
		v2, changed = fastpathTV.DecMapInt64Int64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]float32:
		fastpathTV.DecMapInt64Float32V(v, false, d)
	case *map[int64]float32:
		var v2 map[int64]float32
		v2, changed = fastpathTV.DecMapInt64Float32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]float64:
		fastpathTV.DecMapInt64Float64V(v, false, d)
	case *map[int64]float64:
		var v2 map[int64]float64
		v2, changed = fastpathTV.DecMapInt64Float64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[int64]bool:
		fastpathTV.DecMapInt64BoolV(v, false, d)
	case *map[int64]bool:
		var v2 map[int64]bool
		v2, changed = fastpathTV.DecMapInt64BoolV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]interface{}:
		fastpathTV.DecMapBoolIntfV(v, false, d)
	case *map[bool]interface{}:
		var v2 map[bool]interface{}
		v2, changed = fastpathTV.DecMapBoolIntfV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]string:
		fastpathTV.DecMapBoolStringV(v, false, d)
	case *map[bool]string:
		var v2 map[bool]string
		v2, changed = fastpathTV.DecMapBoolStringV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]uint:
		fastpathTV.DecMapBoolUintV(v, false, d)
	case *map[bool]uint:
		var v2 map[bool]uint
		v2, changed = fastpathTV.DecMapBoolUintV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]uint8:
		fastpathTV.DecMapBoolUint8V(v, false, d)
	case *map[bool]uint8:
		var v2 map[bool]uint8
		v2, changed = fastpathTV.DecMapBoolUint8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]uint16:
		fastpathTV.DecMapBoolUint16V(v, false, d)
	case *map[bool]uint16:
		var v2 map[bool]uint16
		v2, changed = fastpathTV.DecMapBoolUint16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]uint32:
		fastpathTV.DecMapBoolUint32V(v, false, d)
	case *map[bool]uint32:
		var v2 map[bool]uint32
		v2, changed = fastpathTV.DecMapBoolUint32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]uint64:
		fastpathTV.DecMapBoolUint64V(v, false, d)
	case *map[bool]uint64:
		var v2 map[bool]uint64
		v2, changed = fastpathTV.DecMapBoolUint64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]uintptr:
		fastpathTV.DecMapBoolUintptrV(v, false, d)
	case *map[bool]uintptr:
		var v2 map[bool]uintptr
		v2, changed = fastpathTV.DecMapBoolUintptrV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]int:
		fastpathTV.DecMapBoolIntV(v, false, d)
	case *map[bool]int:
		var v2 map[bool]int
		v2, changed = fastpathTV.DecMapBoolIntV(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]int8:
		fastpathTV.DecMapBoolInt8V(v, false, d)
	case *map[bool]int8:
		var v2 map[bool]int8
		v2, changed = fastpathTV.DecMapBoolInt8V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]int16:
		fastpathTV.DecMapBoolInt16V(v, false, d)
	case *map[bool]int16:
		var v2 map[bool]int16
		v2, changed = fastpathTV.DecMapBoolInt16V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]int32:
		fastpathTV.DecMapBoolInt32V(v, false, d)
	case *map[bool]int32:
		var v2 map[bool]int32
		v2, changed = fastpathTV.DecMapBoolInt32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]int64:
		fastpathTV.DecMapBoolInt64V(v, false, d)
	case *map[bool]int64:
		var v2 map[bool]int64
		v2, changed = fastpathTV.DecMapBoolInt64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]float32:
		fastpathTV.DecMapBoolFloat32V(v, false, d)
	case *map[bool]float32:
		var v2 map[bool]float32
		v2, changed = fastpathTV.DecMapBoolFloat32V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]float64:
		fastpathTV.DecMapBoolFloat64V(v, false, d)
	case *map[bool]float64:
		var v2 map[bool]float64
		v2, changed = fastpathTV.DecMapBoolFloat64V(*v, true, d)
		if changed {
			*v = v2
		}
	case map[bool]bool:
		fastpathTV.DecMapBoolBoolV(v, false, d)
	case *map[bool]bool:
		var v2 map[bool]bool
		v2, changed = fastpathTV.DecMapBoolBoolV(*v, true, d)
		if changed {
			*v = v2
		}
	default:
		_ = v // workaround https://github.com/golang/go/issues/12927 seen in go1.4
		return false
	}
	return true
}
