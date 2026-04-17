func Transform(t Transformer, x interface{}) (interface{}, error) {

	if term, ok := x.(*Term); ok {
		return Transform(t, term.Value)
	}

	y, err := t.Transform(x)
	if err != nil {
		return x, err
	}

	if y == nil {
		return nil, nil
	}

	var ok bool
	switch y := y.(type) {
	case *Module:
		p, err := Transform(t, y.Package)
		if err != nil {
			return nil, err
		}
		if y.Package, ok = p.(*Package); !ok {
			return nil, fmt.Errorf("illegal transform: %T != %T", y.Package, p)
		}
		for i := range y.Imports {
			imp, err := Transform(t, y.Imports[i])
			if err != nil {
				return nil, err
			}
			if y.Imports[i], ok = imp.(*Import); !ok {
				return nil, fmt.Errorf("illegal transform: %T != %T", y.Imports[i], imp)
			}
		}
		for i := range y.Rules {
			rule, err := Transform(t, y.Rules[i])
			if err != nil {
				return nil, err
			}
			if y.Rules[i], ok = rule.(*Rule); !ok {
				return nil, fmt.Errorf("illegal transform: %T != %T", y.Rules[i], rule)
			}
		}
		for i := range y.Comments {
			comment, err := Transform(t, y.Comments[i])
			if err != nil {
				return nil, err
			}
			if y.Comments[i], ok = comment.(*Comment); !ok {
				return nil, fmt.Errorf("illegal transform: %T != %T", y.Comments[i], comment)
			}
		}
		return y, nil
	case *Package:
		ref, err := Transform(t, y.Path)
		if err != nil {
			return nil, err
		}
		if y.Path, ok = ref.(Ref); !ok {
			return nil, fmt.Errorf("illegal transform: %T != %T", y.Path, ref)
		}
		return y, nil
	case *Import:
		y.Path, err = transformTerm(t, y.Path)
		if err != nil {
			return nil, err
		}
		if y.Alias, err = transformVar(t, y.Alias); err != nil {
			return nil, err
		}
		return y, nil
	case *Rule:
		if y.Head, err = transformHead(t, y.Head); err != nil {
			return nil, err
		}
		if y.Body, err = transformBody(t, y.Body); err != nil {
			return nil, err
		}
		if y.Else != nil {
			rule, err := Transform(t, y.Else)
			if err != nil {
				return nil, err
			}
			if y.Else, ok = rule.(*Rule); !ok {
				return nil, fmt.Errorf("illegal transform: %T != %T", y.Else, rule)
			}
		}
		return y, nil
	case *Head:
		if y.Name, err = transformVar(t, y.Name); err != nil {
			return nil, err
		}
		if y.Args, err = transformArgs(t, y.Args); err != nil {
			return nil, err
		}
		if y.Key != nil {
			if y.Key, err = transformTerm(t, y.Key); err != nil {
				return nil, err
			}
		}
		if y.Value != nil {
			if y.Value, err = transformTerm(t, y.Value); err != nil {
				return nil, err
			}
		}
		return y, nil
	case Args:
		for i := range y {
			if y[i], err = transformTerm(t, y[i]); err != nil {
				return nil, err
			}
		}
		return y, nil
	case Body:
		for i, e := range y {
			e, err := Transform(t, e)
			if err != nil {
				return nil, err
			}
			if y[i], ok = e.(*Expr); !ok {
				return nil, fmt.Errorf("illegal transform: %T != %T", y[i], e)
			}
		}
		return y, nil
	case *Expr:
		switch ts := y.Terms.(type) {
		case []*Term:
			for i := range ts {
				if ts[i], err = transformTerm(t, ts[i]); err != nil {
					return nil, err
				}
			}
		case *Term:
			if y.Terms, err = transformTerm(t, ts); err != nil {
				return nil, err
			}
		}
		for i, w := range y.With {
			w, err := Transform(t, w)
			if err != nil {
				return nil, err
			}
			if y.With[i], ok = w.(*With); !ok {
				return nil, fmt.Errorf("illegal transform: %T != %T", y.With[i], w)
			}
		}
		return y, nil
	case *With:
		if y.Target, err = transformTerm(t, y.Target); err != nil {
			return nil, err
		}
		if y.Value, err = transformTerm(t, y.Value); err != nil {
			return nil, err
		}
		return y, nil
	case Ref:
		for i, term := range y {
			if y[i], err = transformTerm(t, term); err != nil {
				return nil, err
			}
		}
		return y, nil
	case Object:
		return y.Map(func(k, v *Term) (*Term, *Term, error) {
			k, err := transformTerm(t, k)
			if err != nil {
				return nil, nil, err
			}
			v, err = transformTerm(t, v)
			if err != nil {
				return nil, nil, err
			}
			return k, v, nil
		})
	case Array:
		for i := range y {
			if y[i], err = transformTerm(t, y[i]); err != nil {
				return nil, err
			}
		}
		return y, nil
	case Set:
		y, err = y.Map(func(term *Term) (*Term, error) {
			return transformTerm(t, term)
		})
		if err != nil {
			return nil, err
		}
		return y, nil
	case *ArrayComprehension:
		if y.Term, err = transformTerm(t, y.Term); err != nil {
			return nil, err
		}
		if y.Body, err = transformBody(t, y.Body); err != nil {
			return nil, err
		}
		return y, nil
	case *ObjectComprehension:
		if y.Key, err = transformTerm(t, y.Key); err != nil {
			return nil, err
		}
		if y.Value, err = transformTerm(t, y.Value); err != nil {
			return nil, err
		}
		if y.Body, err = transformBody(t, y.Body); err != nil {
			return nil, err
		}
		return y, nil
	case *SetComprehension:
		if y.Term, err = transformTerm(t, y.Term); err != nil {
			return nil, err
		}
		if y.Body, err = transformBody(t, y.Body); err != nil {
			return nil, err
		}
		return y, nil
	case Call:
		for i := range y {
			if y[i], err = transformTerm(t, y[i]); err != nil {
				return nil, err
			}
		}
		return y, nil
	default:
		return y, nil
	}
}
