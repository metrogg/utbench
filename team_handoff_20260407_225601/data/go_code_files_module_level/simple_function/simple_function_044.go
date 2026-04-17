func (b *builder) processFunctionNode(root *functionNode) (query, error) {
	var qyOutput query
	switch root.FuncName {
	case "starts-with":
		arg1, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		arg2, err := b.processNode(root.Args[1])
		if err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: b.firstInput, Func: startwithFunc(arg1, arg2)}
	case "ends-with":
		arg1, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		arg2, err := b.processNode(root.Args[1])
		if err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: b.firstInput, Func: endwithFunc(arg1, arg2)}
	case "contains":
		arg1, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		arg2, err := b.processNode(root.Args[1])
		if err != nil {
			return nil, err
		}

		qyOutput = &functionQuery{Input: b.firstInput, Func: containsFunc(arg1, arg2)}
	case "substring":
		//substring( string , start [, length] )
		if len(root.Args) < 2 {
			return nil, errors.New("xpath: substring function must have at least two parameter")
		}
		var (
			arg1, arg2, arg3 query
			err              error
		)
		if arg1, err = b.processNode(root.Args[0]); err != nil {
			return nil, err
		}
		if arg2, err = b.processNode(root.Args[1]); err != nil {
			return nil, err
		}
		if len(root.Args) == 3 {
			if arg3, err = b.processNode(root.Args[2]); err != nil {
				return nil, err
			}
		}
		qyOutput = &functionQuery{Input: b.firstInput, Func: substringFunc(arg1, arg2, arg3)}
	case "substring-before", "substring-after":
		//substring-xxxx( haystack, needle )
		if len(root.Args) != 2 {
			return nil, errors.New("xpath: substring-before function must have two parameters")
		}
		var (
			arg1, arg2 query
			err        error
		)
		if arg1, err = b.processNode(root.Args[0]); err != nil {
			return nil, err
		}
		if arg2, err = b.processNode(root.Args[1]); err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{
			Input: b.firstInput,
			Func:  substringIndFunc(arg1, arg2, root.FuncName == "substring-after"),
		}
	case "string-length":
		// string-length( [string] )
		if len(root.Args) < 1 {
			return nil, errors.New("xpath: string-length function must have at least one parameter")
		}
		arg1, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: b.firstInput, Func: stringLengthFunc(arg1)}
	case "normalize-space":
		if len(root.Args) == 0 {
			return nil, errors.New("xpath: normalize-space function must have at least one parameter")
		}
		argQuery, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: argQuery, Func: normalizespaceFunc}
	case "translate":
		//translate( string , string, string )
		if len(root.Args) != 3 {
			return nil, errors.New("xpath: translate function must have three parameters")
		}
		var (
			arg1, arg2, arg3 query
			err              error
		)
		if arg1, err = b.processNode(root.Args[0]); err != nil {
			return nil, err
		}
		if arg2, err = b.processNode(root.Args[1]); err != nil {
			return nil, err
		}
		if arg3, err = b.processNode(root.Args[2]); err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: b.firstInput, Func: translateFunc(arg1, arg2, arg3)}
	case "not":
		if len(root.Args) == 0 {
			return nil, errors.New("xpath: not function must have at least one parameter")
		}
		argQuery, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: argQuery, Func: notFunc}
	case "name", "local-name", "namespace-uri":
		inp := b.firstInput
		if len(root.Args) > 1 {
			return nil, fmt.Errorf("xpath: %s function must have at most one parameter", root.FuncName)
		}
		if len(root.Args) == 1 {
			argQuery, err := b.processNode(root.Args[0])
			if err != nil {
				return nil, err
			}
			inp = argQuery
		}
		f := &functionQuery{Input: inp}
		switch root.FuncName {
		case "name":
			f.Func = nameFunc
		case "local-name":
			f.Func = localNameFunc
		case "namespace-uri":
			f.Func = namespaceFunc
		}
		qyOutput = f
	case "true", "false":
		val := root.FuncName == "true"
		qyOutput = &functionQuery{
			Input: b.firstInput,
			Func: func(_ query, _ iterator) interface{} {
				return val
			},
		}
	case "last":
		qyOutput = &functionQuery{Input: b.firstInput, Func: lastFunc}
	case "position":
		qyOutput = &functionQuery{Input: b.firstInput, Func: positionFunc}
	case "boolean", "number", "string":
		inp := b.firstInput
		if len(root.Args) > 1 {
			return nil, fmt.Errorf("xpath: %s function must have at most one parameter", root.FuncName)
		}
		if len(root.Args) == 1 {
			argQuery, err := b.processNode(root.Args[0])
			if err != nil {
				return nil, err
			}
			inp = argQuery
		}
		f := &functionQuery{Input: inp}
		switch root.FuncName {
		case "boolean":
			f.Func = booleanFunc
		case "string":
			f.Func = stringFunc
		case "number":
			f.Func = numberFunc
		}
		qyOutput = f
	case "count":
		//if b.firstInput == nil {
		//	return nil, errors.New("xpath: expression must evaluate to node-set")
		//}
		if len(root.Args) == 0 {
			return nil, fmt.Errorf("xpath: count(node-sets) function must with have parameters node-sets")
		}
		argQuery, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: argQuery, Func: countFunc}
	case "sum":
		if len(root.Args) == 0 {
			return nil, fmt.Errorf("xpath: sum(node-sets) function must with have parameters node-sets")
		}
		argQuery, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		qyOutput = &functionQuery{Input: argQuery, Func: sumFunc}
	case "ceiling", "floor", "round":
		if len(root.Args) == 0 {
			return nil, fmt.Errorf("xpath: ceiling(node-sets) function must with have parameters node-sets")
		}
		argQuery, err := b.processNode(root.Args[0])
		if err != nil {
			return nil, err
		}
		f := &functionQuery{Input: argQuery}
		switch root.FuncName {
		case "ceiling":
			f.Func = ceilingFunc
		case "floor":
			f.Func = floorFunc
		case "round":
			f.Func = roundFunc
		}
		qyOutput = f
	case "concat":
		if len(root.Args) < 2 {
			return nil, fmt.Errorf("xpath: concat() must have at least two arguments")
		}
		var args []query
		for _, v := range root.Args {
			q, err := b.processNode(v)
			if err != nil {
				return nil, err
			}
			args = append(args, q)
		}
		qyOutput = &functionQuery{Input: b.firstInput, Func: concatFunc(args...)}
	default:
		return nil, fmt.Errorf("not yet support this function %s()", root.FuncName)
	}
	return qyOutput, nil
}
