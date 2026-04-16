func (s *state) evalExpr(exp parse.Expr) (v Value, e error) {
	switch exp := exp.(type) {
	case *parse.NullExpr:
		return nil, nil
	case *parse.BoolExpr:
		return exp.Value, nil
	case *parse.NameExpr:
		if val, ok := s.scope.Get(exp.Name); ok {
			v = val
		} else {
			e = errors.New("undefined variable \"" + exp.Name + "\"")
		}
	case *parse.NumberExpr:
		num, err := strconv.ParseFloat(exp.Value, 64)
		if err != nil {
			return nil, err
		}
		return num, nil
	case *parse.StringExpr:
		return exp.Text, nil
	case *parse.GroupExpr:
		return s.evalExpr(exp.X)
	case *parse.UnaryExpr:
		in, err := s.evalExpr(exp.X)
		if err != nil {
			return nil, err
		}
		switch exp.Op {
		case parse.OpUnaryNot:
			return !CoerceBool(in), nil
		case parse.OpUnaryPositive:
			// no-op, +1 = 1, +(-1) = -1, +(false) = 0
			return CoerceNumber(in), nil
		case parse.OpUnaryNegative:
			return -CoerceNumber(in), nil
		}
	case *parse.BinaryExpr:
		left, err := s.evalExpr(exp.Left)
		if err != nil {
			return nil, err
		}
		right, err := s.evalExpr(exp.Right)
		if err != nil {
			return nil, err
		}
		switch exp.Op {
		case parse.OpBinaryAdd:
			return CoerceNumber(left) + CoerceNumber(right), nil
		case parse.OpBinarySubtract:
			return CoerceNumber(left) - CoerceNumber(right), nil
		case parse.OpBinaryMultiply:
			return CoerceNumber(left) * CoerceNumber(right), nil
		case parse.OpBinaryDivide:
			return CoerceNumber(left) / CoerceNumber(right), nil
		case parse.OpBinaryFloorDiv:
			return math.Floor(CoerceNumber(left) / CoerceNumber(right)), nil
		case parse.OpBinaryModulo:
			return float64(int(CoerceNumber(left)) % int(CoerceNumber(right))), nil
		case parse.OpBinaryPower:
			return math.Pow(CoerceNumber(left), CoerceNumber(right)), nil
		case parse.OpBinaryConcat:
			return CoerceString(left) + CoerceString(right), nil
		case parse.OpBinaryEndsWith:
			return strings.HasSuffix(CoerceString(left), CoerceString(right)), nil
		case parse.OpBinaryStartsWith:
			return strings.HasPrefix(CoerceString(left), CoerceString(right)), nil
		case parse.OpBinaryIn:
			return Contains(right, left)
		case parse.OpBinaryNotIn:
			res, err := Contains(right, left)
			if err != nil {
				return false, err
			}
			return !res, nil
		case parse.OpBinaryIs:
			if fn, ok := right.(func(v Value) bool); ok {
				return fn(left), nil
			}
			return nil, errors.New("right operand was of unexpected type")
		case parse.OpBinaryIsNot:
			if fn, ok := right.(func(v Value) bool); ok {
				return !fn(left), nil
			}
			return nil, errors.New("right operand was of unexpected type")
		case parse.OpBinaryMatches:
			reg, err := regexp.Compile(CoerceString(right))
			if err != nil {
				return nil, err
			}
			return reg.MatchString(CoerceString(left)), nil
		case parse.OpBinaryEqual:
			return Equal(left, right), nil
		case parse.OpBinaryNotEqual:
			return !Equal(left, right), nil
		case parse.OpBinaryGreaterEqual:
			return CoerceNumber(left) >= CoerceNumber(right), nil
		case parse.OpBinaryGreaterThan:
			return CoerceNumber(left) > CoerceNumber(right), nil
		case parse.OpBinaryLessEqual:
			return CoerceNumber(left) <= CoerceNumber(right), nil
		case parse.OpBinaryLessThan:
			return CoerceNumber(left) < CoerceNumber(right), nil
		case parse.OpBinaryRange:
			l, r := CoerceNumber(left), CoerceNumber(right)
			res := make([]float64, uint(math.Ceil(r-l))+1)
			for i, k := 0, l; k <= r; i, k = i+1, k+1 {
				res[i] = k
			}
			return res, nil
		case parse.OpBinaryBitwiseAnd:
			return int(CoerceNumber(left)) & int(CoerceNumber(right)), nil
		case parse.OpBinaryBitwiseOr:
			return int(CoerceNumber(left)) | int(CoerceNumber(right)), nil
		case parse.OpBinaryBitwiseXor:
			return int(CoerceNumber(left)) ^ int(CoerceNumber(right)), nil
		case parse.OpBinaryAnd:
			return CoerceBool(left) && CoerceBool(right), nil
		case parse.OpBinaryOr:
			return CoerceBool(left) && CoerceBool(right), nil
		}
	case *parse.FuncExpr:
		return s.evalFunction(exp)
	case *parse.FilterExpr:
		return s.evalFilter(exp)
	case *parse.GetAttrExpr:
		c, err := s.evalExpr(exp.Cont)
		if err != nil {
			return nil, err
		}
		k, err := s.evalExpr(exp.Attr)
		if err != nil {
			return nil, err
		}
		exargs := exp.Args
		args := make([]Value, len(exargs))
		for k, e := range exargs {
			v, err := s.evalExpr(e)
			if err != nil {
				return nil, err
			}
			args[k] = v
		}
		if set, ok := c.(macroSet); ok {
			if macro, ok := set.defs[CoerceString(k)]; ok {
				return s.callMacro(macro, args...)
			}
			return nil, errors.New("undefined macro: " + CoerceString(k))
		}
		v, err = GetAttr(c, k, args...)
		if err != nil {
			return nil, err
		}
	case *parse.TestExpr:
		if tfn, ok := s.env.Tests[exp.Name]; ok {
			eargs := exp.Args
			args := make([]Value, len(eargs))
			for i, e := range eargs {
				v, err := s.evalExpr(e)
				if err != nil {
					return nil, err
				}
				args[i] = v
			}
			return func(v Value) bool {
				return tfn(s, v, args...)
			}, nil
		}
		return nil, fmt.Errorf(`unknown test "%v"`, exp.Name)
	case *parse.TernaryIfExpr:
		cond, err := s.evalExpr(exp.Cond)
		if err != nil {
			return nil, err
		}
		if CoerceBool(cond) == true {
			return s.evalExpr(exp.TrueX)
		}
		return s.evalExpr(exp.FalseX)

	case *parse.HashExpr:
		vals := make(map[string]Value)
		for _, v := range exp.Elements {
			var key Value
			var err error
			if k, ok := v.Key.(*parse.NameExpr); ok {
				key = k.Name
			} else {
				key, err = s.evalExpr(v.Key)
				if err != nil {
					return nil, err
				}
			}
			val, err := s.evalExpr(v.Value)
			if err != nil {
				return nil, err
			}
			vals[CoerceString(key)] = val
		}
		return vals, nil

	case *parse.ArrayExpr:
		vals := make([]Value, len(exp.Elements))
		for i, v := range exp.Elements {
			val, err := s.evalExpr(v)
			if err != nil {
				return nil, err
			}
			vals[i] = val
		}
		return vals, nil
	}

	return v, nil
}
