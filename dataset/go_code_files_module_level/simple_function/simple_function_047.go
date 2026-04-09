func typecheck1(cfg *TypeConfig, f interface{}, typeof map[interface{}]string, assign map[string][]interface{}) {
	// set sets the type of n to typ.
	// If isDecl is true, n is being declared.
	set := func(n ast.Expr, typ string, isDecl bool) {
		if typeof[n] != "" || typ == "" {
			if typeof[n] != typ {
				assign[typ] = append(assign[typ], n)
			}
			return
		}
		typeof[n] = typ

		// If we obtained typ from the declaration of x
		// propagate the type to all the uses.
		// The !isDecl case is a cheat here, but it makes
		// up in some cases for not paying attention to
		// struct fields.  The real type checker will be
		// more accurate so we won't need the cheat.
		if id, ok := n.(*ast.Ident); ok && id.Obj != nil && (isDecl || typeof[id.Obj] == "") {
			typeof[id.Obj] = typ
		}
	}

	// Type-check an assignment lhs = rhs.
	// If isDecl is true, this is := so we can update
	// the types of the objects that lhs refers to.
	typecheckAssign := func(lhs, rhs []ast.Expr, isDecl bool) {
		if len(lhs) > 1 && len(rhs) == 1 {
			if _, ok := rhs[0].(*ast.CallExpr); ok {
				t := split(typeof[rhs[0]])
				// Lists should have same length but may not; pair what can be paired.
				for i := 0; i < len(lhs) && i < len(t); i++ {
					set(lhs[i], t[i], isDecl)
				}
				return
			}
		}
		if len(lhs) == 1 && len(rhs) == 2 {
			// x = y, ok
			rhs = rhs[:1]
		} else if len(lhs) == 2 && len(rhs) == 1 {
			// x, ok = y
			lhs = lhs[:1]
		}

		// Match as much as we can.
		for i := 0; i < len(lhs) && i < len(rhs); i++ {
			x, y := lhs[i], rhs[i]
			if typeof[y] != "" {
				set(x, typeof[y], isDecl)
			} else {
				set(y, typeof[x], false)
			}
		}
	}

	expand := func(s string) string {
		typ := cfg.Type[s]
		if typ != nil && typ.Def != "" {
			return typ.Def
		}
		return s
	}

	// The main type check is a recursive algorithm implemented
	// by walkBeforeAfter(n, before, after).
	// Most of it is bottom-up, but in a few places we need
	// to know the type of the function we are checking.
	// The before function records that information on
	// the curfn stack.
	var curfn []*ast.FuncType

	before := func(n interface{}) {
		// push function type on stack
		switch n := n.(type) {
		case *ast.FuncDecl:
			curfn = append(curfn, n.Type)
		case *ast.FuncLit:
			curfn = append(curfn, n.Type)
		}
	}

	// After is the real type checker.
	after := func(n interface{}) {
		if n == nil {
			return
		}
		if false && reflect.TypeOf(n).Kind() == reflect.Ptr { // debugging trace
			defer func() {
				if t := typeof[n]; t != "" {
					pos := fset.Position(n.(ast.Node).Pos())
					fmt.Fprintf(os.Stderr, "%s: typeof[%s] = %s\n", pos, gofmt(n), t)
				}
			}()
		}

		switch n := n.(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			// pop function type off stack
			curfn = curfn[:len(curfn)-1]

		case *ast.FuncType:
			typeof[n] = mkType(joinFunc(split(typeof[n.Params]), split(typeof[n.Results])))

		case *ast.FieldList:
			// Field list is concatenation of sub-lists.
			t := ""
			for _, field := range n.List {
				if t != "" {
					t += ", "
				}
				t += typeof[field]
			}
			typeof[n] = t

		case *ast.Field:
			// Field is one instance of the type per name.
			all := ""
			t := typeof[n.Type]
			if !isType(t) {
				// Create a type, because it is typically *T or *p.T
				// and we might care about that type.
				t = mkType(gofmt(n.Type))
				typeof[n.Type] = t
			}
			t = getType(t)
			if len(n.Names) == 0 {
				all = t
			} else {
				for _, id := range n.Names {
					if all != "" {
						all += ", "
					}
					all += t
					typeof[id.Obj] = t
					typeof[id] = t
				}
			}
			typeof[n] = all

		case *ast.ValueSpec:
			// var declaration.  Use type if present.
			if n.Type != nil {
				t := typeof[n.Type]
				if !isType(t) {
					t = mkType(gofmt(n.Type))
					typeof[n.Type] = t
				}
				t = getType(t)
				for _, id := range n.Names {
					set(id, t, true)
				}
			}
			// Now treat same as assignment.
			typecheckAssign(makeExprList(n.Names), n.Values, true)

		case *ast.AssignStmt:
			typecheckAssign(n.Lhs, n.Rhs, n.Tok == token.DEFINE)

		case *ast.Ident:
			// Identifier can take its type from underlying object.
			if t := typeof[n.Obj]; t != "" {
				typeof[n] = t
			}

		case *ast.SelectorExpr:
			// Field or method.
			name := n.Sel.Name
			if t := typeof[n.X]; t != "" {
				if strings.HasPrefix(t, "*") {
					t = t[1:] // implicit *
				}
				if typ := cfg.Type[t]; typ != nil {
					if t := typ.dot(cfg, name); t != "" {
						typeof[n] = t
						return
					}
				}
				tt := typeof[t+"."+name]
				if isType(tt) {
					typeof[n] = getType(tt)
					return
				}
			}
			// Package selector.
			if x, ok := n.X.(*ast.Ident); ok && x.Obj == nil {
				str := x.Name + "." + name
				if cfg.Type[str] != nil {
					typeof[n] = mkType(str)
					return
				}
				if t := cfg.typeof(x.Name + "." + name); t != "" {
					typeof[n] = t
					return
				}
			}

		case *ast.CallExpr:
			// make(T) has type T.
			if isTopName(n.Fun, "make") && len(n.Args) >= 1 {
				typeof[n] = gofmt(n.Args[0])
				return
			}
			// new(T) has type *T
			if isTopName(n.Fun, "new") && len(n.Args) == 1 {
				typeof[n] = "*" + gofmt(n.Args[0])
				return
			}
			// Otherwise, use type of function to determine arguments.
			t := typeof[n.Fun]
			in, out := splitFunc(t)
			if in == nil && out == nil {
				return
			}
			typeof[n] = join(out)
			for i, arg := range n.Args {
				if i >= len(in) {
					break
				}
				if typeof[arg] == "" {
					typeof[arg] = in[i]
				}
			}

		case *ast.TypeAssertExpr:
			// x.(type) has type of x.
			if n.Type == nil {
				typeof[n] = typeof[n.X]
				return
			}
			// x.(T) has type T.
			if t := typeof[n.Type]; isType(t) {
				typeof[n] = getType(t)
			} else {
				typeof[n] = gofmt(n.Type)
			}

		case *ast.SliceExpr:
			// x[i:j] has type of x.
			typeof[n] = typeof[n.X]

		case *ast.IndexExpr:
			// x[i] has key type of x's type.
			t := expand(typeof[n.X])
			if strings.HasPrefix(t, "[") || strings.HasPrefix(t, "map[") {
				// Lazy: assume there are no nested [] in the array
				// length or map key type.
				if i := strings.Index(t, "]"); i >= 0 {
					typeof[n] = t[i+1:]
				}
			}

		case *ast.StarExpr:
			// *x for x of type *T has type T when x is an expr.
			// We don't use the result when *x is a type, but
			// compute it anyway.
			t := expand(typeof[n.X])
			if isType(t) {
				typeof[n] = "type *" + getType(t)
			} else if strings.HasPrefix(t, "*") {
				typeof[n] = t[len("*"):]
			}

		case *ast.UnaryExpr:
			// &x for x of type T has type *T.
			t := typeof[n.X]
			if t != "" && n.Op == token.AND {
				typeof[n] = "*" + t
			}

		case *ast.CompositeLit:
			// T{...} has type T.
			typeof[n] = gofmt(n.Type)

		case *ast.ParenExpr:
			// (x) has type of x.
			typeof[n] = typeof[n.X]

		case *ast.RangeStmt:
			t := expand(typeof[n.X])
			if t == "" {
				return
			}
			var key, value string
			if t == "string" {
				key, value = "int", "rune"
			} else if strings.HasPrefix(t, "[") {
				key = "int"
				if i := strings.Index(t, "]"); i >= 0 {
					value = t[i+1:]
				}
			} else if strings.HasPrefix(t, "map[") {
				if i := strings.Index(t, "]"); i >= 0 {
					key, value = t[4:i], t[i+1:]
				}
			}
			changed := false
			if n.Key != nil && key != "" {
				changed = true
				set(n.Key, key, n.Tok == token.DEFINE)
			}
			if n.Value != nil && value != "" {
				changed = true
				set(n.Value, value, n.Tok == token.DEFINE)
			}
			// Ugly failure of vision: already type-checked body.
			// Do it again now that we have that type info.
			if changed {
				typecheck1(cfg, n.Body, typeof, assign)
			}

		case *ast.TypeSwitchStmt:
			// Type of variable changes for each case in type switch,
			// but go/parser generates just one variable.
			// Repeat type check for each case with more precise
			// type information.
			as, ok := n.Assign.(*ast.AssignStmt)
			if !ok {
				return
			}
			varx, ok := as.Lhs[0].(*ast.Ident)
			if !ok {
				return
			}
			t := typeof[varx]
			for _, cas := range n.Body.List {
				cas := cas.(*ast.CaseClause)
				if len(cas.List) == 1 {
					// Variable has specific type only when there is
					// exactly one type in the case list.
					if tt := typeof[cas.List[0]]; isType(tt) {
						tt = getType(tt)
						typeof[varx] = tt
						typeof[varx.Obj] = tt
						typecheck1(cfg, cas.Body, typeof, assign)
					}
				}
			}
			// Restore t.
			typeof[varx] = t
			typeof[varx.Obj] = t

		case *ast.ReturnStmt:
			if len(curfn) == 0 {
				// Probably can't happen.
				return
			}
			f := curfn[len(curfn)-1]
			res := n.Results
			if f.Results != nil {
				t := split(typeof[f.Results])
				for i := 0; i < len(res) && i < len(t); i++ {
					set(res[i], t[i], false)
				}
			}
		}
	}
	walkBeforeAfter(f, before, after)
}
