func Binary(op syntax.Token, x, y Value) (Value, error) {
	switch op {
	case syntax.PLUS:
		switch x := x.(type) {
		case String:
			if y, ok := y.(String); ok {
				return x + y, nil
			}
		case Int:
			switch y := y.(type) {
			case Int:
				return x.Add(y), nil
			case Float:
				return x.Float() + y, nil
			}
		case Float:
			switch y := y.(type) {
			case Float:
				return x + y, nil
			case Int:
				return x + y.Float(), nil
			}
		case *List:
			if y, ok := y.(*List); ok {
				z := make([]Value, 0, x.Len()+y.Len())
				z = append(z, x.elems...)
				z = append(z, y.elems...)
				return NewList(z), nil
			}
		case Tuple:
			if y, ok := y.(Tuple); ok {
				z := make(Tuple, 0, len(x)+len(y))
				z = append(z, x...)
				z = append(z, y...)
				return z, nil
			}
		case *Dict:
			// Python doesn't have dict+dict, and I can't find
			// it documented for Skylark.  But it is used; see:
			//   tools/build_defs/haskell/def.bzl:448
			// TODO(adonovan): clarify spec; see b/36360157.
			if y, ok := y.(*Dict); ok {
				z := new(Dict)
				for _, item := range x.Items() {
					z.SetKey(item[0], item[1])
				}
				for _, item := range y.Items() {
					z.SetKey(item[0], item[1])
				}
				return z, nil
			}
		}

	case syntax.MINUS:
		switch x := x.(type) {
		case Int:
			switch y := y.(type) {
			case Int:
				return x.Sub(y), nil
			case Float:
				return x.Float() - y, nil
			}
		case Float:
			switch y := y.(type) {
			case Float:
				return x - y, nil
			case Int:
				return x - y.Float(), nil
			}
		}

	case syntax.STAR:
		switch x := x.(type) {
		case Int:
			switch y := y.(type) {
			case Int:
				return x.Mul(y), nil
			case Float:
				return x.Float() * y, nil
			case String:
				if i, err := AsInt32(x); err == nil {
					if i < 1 {
						return String(""), nil
					}
					return String(strings.Repeat(string(y), i)), nil
				}
			case *List:
				if i, err := AsInt32(x); err == nil {
					return NewList(repeat(y.elems, i)), nil
				}
			case Tuple:
				if i, err := AsInt32(x); err == nil {
					return Tuple(repeat([]Value(y), i)), nil
				}
			}
		case Float:
			switch y := y.(type) {
			case Float:
				return x * y, nil
			case Int:
				return x * y.Float(), nil
			}
		case String:
			if y, ok := y.(Int); ok {
				if i, err := AsInt32(y); err == nil {
					if i < 1 {
						return String(""), nil
					}
					return String(strings.Repeat(string(x), i)), nil
				}
			}
		case *List:
			if y, ok := y.(Int); ok {
				if i, err := AsInt32(y); err == nil {
					return NewList(repeat(x.elems, i)), nil
				}
			}
		case Tuple:
			if y, ok := y.(Int); ok {
				if i, err := AsInt32(y); err == nil {
					return Tuple(repeat([]Value(x), i)), nil
				}
			}

		}

	case syntax.SLASH:
		switch x := x.(type) {
		case Int:
			switch y := y.(type) {
			case Int:
				yf := y.Float()
				if yf == 0.0 {
					return nil, fmt.Errorf("real division by zero")
				}
				return x.Float() / yf, nil
			case Float:
				if y == 0.0 {
					return nil, fmt.Errorf("real division by zero")
				}
				return x.Float() / y, nil
			}
		case Float:
			switch y := y.(type) {
			case Float:
				if y == 0.0 {
					return nil, fmt.Errorf("real division by zero")
				}
				return x / y, nil
			case Int:
				yf := y.Float()
				if yf == 0.0 {
					return nil, fmt.Errorf("real division by zero")
				}
				return x / yf, nil
			}
		}

	case syntax.SLASHSLASH:
		switch x := x.(type) {
		case Int:
			switch y := y.(type) {
			case Int:
				if y.Sign() == 0 {
					return nil, fmt.Errorf("floored division by zero")
				}
				return x.Div(y), nil
			case Float:
				if y == 0.0 {
					return nil, fmt.Errorf("floored division by zero")
				}
				return floor((x.Float() / y)), nil
			}
		case Float:
			switch y := y.(type) {
			case Float:
				if y == 0.0 {
					return nil, fmt.Errorf("floored division by zero")
				}
				return floor(x / y), nil
			case Int:
				yf := y.Float()
				if yf == 0.0 {
					return nil, fmt.Errorf("floored division by zero")
				}
				return floor(x / yf), nil
			}
		}

	case syntax.PERCENT:
		switch x := x.(type) {
		case Int:
			switch y := y.(type) {
			case Int:
				if y.Sign() == 0 {
					return nil, fmt.Errorf("integer modulo by zero")
				}
				return x.Mod(y), nil
			case Float:
				if y == 0 {
					return nil, fmt.Errorf("float modulo by zero")
				}
				return x.Float().Mod(y), nil
			}
		case Float:
			switch y := y.(type) {
			case Float:
				if y == 0.0 {
					return nil, fmt.Errorf("float modulo by zero")
				}
				return Float(math.Mod(float64(x), float64(y))), nil
			case Int:
				if y.Sign() == 0 {
					return nil, fmt.Errorf("float modulo by zero")
				}
				return x.Mod(y.Float()), nil
			}
		case String:
			return interpolate(string(x), y)
		}

	case syntax.NOT_IN:
		z, err := Binary(syntax.IN, x, y)
		if err != nil {
			return nil, err
		}
		return !z.Truth(), nil

	case syntax.IN:
		switch y := y.(type) {
		case *List:
			for _, elem := range y.elems {
				if eq, err := Equal(elem, x); err != nil {
					return nil, err
				} else if eq {
					return True, nil
				}
			}
			return False, nil
		case Tuple:
			for _, elem := range y {
				if eq, err := Equal(elem, x); err != nil {
					return nil, err
				} else if eq {
					return True, nil
				}
			}
			return False, nil
		case Mapping: // e.g. dict
			// Ignore error from Get as we cannot distinguish true
			// errors (value cycle, type error) from "key not found".
			_, found, _ := y.Get(x)
			return Bool(found), nil
		case *Set:
			ok, err := y.Has(x)
			return Bool(ok), err
		case String:
			needle, ok := x.(String)
			if !ok {
				return nil, fmt.Errorf("'in <string>' requires string as left operand, not %s", x.Type())
			}
			return Bool(strings.Contains(string(y), string(needle))), nil
		case rangeValue:
			i, err := NumberToInt(x)
			if err != nil {
				return nil, fmt.Errorf("'in <range>' requires integer as left operand, not %s", x.Type())
			}
			return Bool(y.contains(i)), nil
		}

	case syntax.PIPE:
		switch x := x.(type) {
		case Int:
			if y, ok := y.(Int); ok {
				return x.Or(y), nil
			}
		case *Set: // union
			if y, ok := y.(*Set); ok {
				iter := Iterate(y)
				defer iter.Done()
				return x.Union(iter)
			}
		}

	case syntax.AMP:
		switch x := x.(type) {
		case Int:
			if y, ok := y.(Int); ok {
				return x.And(y), nil
			}
		case *Set: // intersection
			if y, ok := y.(*Set); ok {
				set := new(Set)
				if x.Len() > y.Len() {
					x, y = y, x // opt: range over smaller set
				}
				for _, xelem := range x.elems() {
					// Has, Insert cannot fail here.
					if found, _ := y.Has(xelem); found {
						set.Insert(xelem)
					}
				}
				return set, nil
			}
		}

	case syntax.CIRCUMFLEX:
		switch x := x.(type) {
		case Int:
			if y, ok := y.(Int); ok {
				return x.Xor(y), nil
			}
		case *Set: // symmetric difference
			if y, ok := y.(*Set); ok {
				set := new(Set)
				for _, xelem := range x.elems() {
					if found, _ := y.Has(xelem); !found {
						set.Insert(xelem)
					}
				}
				for _, yelem := range y.elems() {
					if found, _ := x.Has(yelem); !found {
						set.Insert(yelem)
					}
				}
				return set, nil
			}
		}

	case syntax.LTLT, syntax.GTGT:
		if x, ok := x.(Int); ok {
			y, err := AsInt32(y)
			if err != nil {
				return nil, err
			}
			if y < 0 {
				return nil, fmt.Errorf("negative shift count: %v", y)
			}
			if op == syntax.LTLT {
				if y >= 512 {
					return nil, fmt.Errorf("shift count too large: %v", y)
				}
				return x.Lsh(uint(y)), nil
			} else {
				return x.Rsh(uint(y)), nil
			}
		}

	default:
		// unknown operator
		goto unknown
	}

	// user-defined types
	if x, ok := x.(HasBinary); ok {
		z, err := x.Binary(op, y, Left)
		if z != nil || err != nil {
			return z, err
		}
	}
	if y, ok := y.(HasBinary); ok {
		z, err := y.Binary(op, x, Right)
		if z != nil || err != nil {
			return z, err
		}
	}

	// unsupported operand types
unknown:
	return nil, fmt.Errorf("unknown binary op: %s %s %s", x.Type(), op, y.Type())
}
