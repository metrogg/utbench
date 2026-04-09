func BQL() *Grammar {
	return &Grammar{
		"START": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemQuery),
					NewSymbol("VARS"),
					NewTokenType(lexer.ItemFrom),
					NewSymbol("INPUT_GRAPHS"),
					NewSymbol("WHERE"),
					NewSymbol("GROUP_BY"),
					NewSymbol("ORDER_BY"),
					NewSymbol("HAVING"),
					NewSymbol("GLOBAL_TIME_BOUND"),
					NewSymbol("LIMIT"),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemInsert),
					NewTokenType(lexer.ItemData),
					NewTokenType(lexer.ItemInto),
					NewSymbol("OUTPUT_GRAPHS"),
					NewTokenType(lexer.ItemLBracket),
					NewTokenType(lexer.ItemNode),
					NewTokenType(lexer.ItemPredicate),
					NewSymbol("INSERT_OBJECT"),
					NewSymbol("INSERT_DATA"),
					NewTokenType(lexer.ItemRBracket),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDelete),
					NewTokenType(lexer.ItemData),
					NewTokenType(lexer.ItemFrom),
					NewSymbol("INPUT_GRAPHS"),
					NewTokenType(lexer.ItemLBracket),
					NewTokenType(lexer.ItemNode),
					NewTokenType(lexer.ItemPredicate),
					NewSymbol("DELETE_OBJECT"),
					NewSymbol("DELETE_DATA"),
					NewTokenType(lexer.ItemRBracket),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemCreate),
					NewSymbol("CREATE_GRAPHS"),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDrop),
					NewSymbol("DROP_GRAPHS"),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemConstruct),
					NewSymbol("CONSTRUCT_FACTS"),
					NewTokenType(lexer.ItemInto),
					NewSymbol("OUTPUT_GRAPHS"),
					NewTokenType(lexer.ItemFrom),
					NewSymbol("INPUT_GRAPHS"),
					NewSymbol("WHERE"),
					NewSymbol("HAVING"),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDeconstruct),
					NewSymbol("DECONSTRUCT_FACTS"),
					NewTokenType(lexer.ItemIn),
					NewSymbol("OUTPUT_GRAPHS"),
					NewTokenType(lexer.ItemFrom),
					NewSymbol("INPUT_GRAPHS"),
					NewSymbol("WHERE"),
					NewSymbol("HAVING"),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemShow),
					NewSymbol("GRAPH_SHOW"),
					NewTokenType(lexer.ItemSemicolon),
				},
			},
		},
		"CREATE_GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemGraph),
					NewSymbol("GRAPHS"),
				},
			},
		},
		"DROP_GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemGraph),
					NewSymbol("GRAPHS"),
				},
			},
		},
		"VARS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("VARS_AS"),
					NewSymbol("MORE_VARS"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemCount),
					NewTokenType(lexer.ItemLPar),
					NewSymbol("COUNT_DISTINCT"),
					NewTokenType(lexer.ItemBinding),
					NewTokenType(lexer.ItemRPar),
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_VARS"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemSum),
					NewTokenType(lexer.ItemLPar),
					NewTokenType(lexer.ItemBinding),
					NewTokenType(lexer.ItemRPar),
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_VARS"),
				},
			},
		},
		"COUNT_DISTINCT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDistinct),
				},
			},
			{},
		},
		"VARS_AS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"MORE_VARS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewSymbol("VARS"),
				},
			},
			{},
		},
		"GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_GRAPHS"),
				},
			},
		},
		"MORE_GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_GRAPHS"),
				},
			},
			{},
		},
		"INPUT_GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_INPUT_GRAPHS"),
				},
			},
		},
		"MORE_INPUT_GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_INPUT_GRAPHS"),
				},
			},
			{},
		},
		"OUTPUT_GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_OUTPUT_GRAPHS"),
				},
			},
		},
		"MORE_OUTPUT_GRAPHS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("MORE_OUTPUT_GRAPHS"),
				},
			},
			{},
		},
		"WHERE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemWhere),
					NewTokenType(lexer.ItemLBracket),
					NewSymbol("CLAUSES"),
					NewTokenType(lexer.ItemRBracket),
				},
			},
		},
		"CLAUSES": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNode),
					NewSymbol("SUBJECT_EXTRACT"),
					NewSymbol("PREDICATE"),
					NewSymbol("OBJECT"),
					NewSymbol("MORE_CLAUSES"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("SUBJECT_EXTRACT"),
					NewSymbol("PREDICATE"),
					NewSymbol("OBJECT"),
					NewSymbol("MORE_CLAUSES"),
				},
			},
		},
		"SUBJECT_EXTRACT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("SUBJECT_TYPE"),
					NewSymbol("SUBJECT_ID"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemType),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("SUBJECT_ID"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemID),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"SUBJECT_TYPE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemType),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"SUBJECT_ID": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemID),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"PREDICATE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicate),
					NewSymbol("PREDICATE_AS"),
					NewSymbol("PREDICATE_ID"),
					NewSymbol("PREDICATE_AT"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicateBound),
					NewSymbol("PREDICATE_AS"),
					NewSymbol("PREDICATE_ID"),
					NewSymbol("PREDICATE_BOUND_AT"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("PREDICATE_AS"),
					NewSymbol("PREDICATE_ID"),
					NewSymbol("PREDICATE_AT"),
				},
			},
		},
		"PREDICATE_AS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"PREDICATE_ID": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemID),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"PREDICATE_AT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAt),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"PREDICATE_BOUND_AT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAt),
					NewSymbol("PREDICATE_BOUND_AT_BINDINGS"),
				},
			},
			{},
		},
		"PREDICATE_BOUND_AT_BINDINGS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("PREDICATE_BOUND_AT_BINDINGS_END"),
				},
			},
			{},
		},
		"PREDICATE_BOUND_AT_BINDINGS_END": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLiteral),
					NewSymbol("OBJECT_LITERAL_AS"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNode),
					NewSymbol("OBJECT_SUBJECT_EXTRACT"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicate),
					NewSymbol("OBJECT_PREDICATE_AS"),
					NewSymbol("OBJECT_PREDICATE_ID"),
					NewSymbol("OBJECT_PREDICATE_AT"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicateBound),
					NewSymbol("OBJECT_PREDICATE_AS"),
					NewSymbol("OBJECT_PREDICATE_ID"),
					NewSymbol("OBJECT_PREDICATE_BOUND_AT"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("OBJECT_LITERAL_BINDING_AS"),
					NewSymbol("OBJECT_LITERAL_BINDING_TYPE"),
					NewSymbol("OBJECT_LITERAL_BINDING_ID"),
					NewSymbol("OBJECT_LITERAL_BINDING_AT"),
				},
			},
		},
		"OBJECT_SUBJECT_EXTRACT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("OBJECT_SUBJECT_TYPE"),
					NewSymbol("OBJECT_SUBJECT_ID"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemType),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("OBJECT_SUBJECT_ID"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemID),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_SUBJECT_TYPE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemType),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_SUBJECT_ID": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemID),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_PREDICATE_AS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_PREDICATE_ID": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemID),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_PREDICATE_AT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAt),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_PREDICATE_BOUND_AT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAt),
					NewSymbol("OBJECT_PREDICATE_BOUND_AT_BINDINGS"),
				},
			},
			{},
		},
		"OBJECT_PREDICATE_BOUND_AT_BINDINGS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("OBJECT_PREDICATE_BOUND_AT_BINDINGS_END"),
				},
			},
			{},
		},
		"OBJECT_PREDICATE_BOUND_AT_BINDINGS_END": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_LITERAL_AS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_LITERAL_BINDING_AS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAs),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_LITERAL_BINDING_TYPE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemType),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_LITERAL_BINDING_ID": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemID),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"OBJECT_LITERAL_BINDING_AT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAt),
					NewTokenType(lexer.ItemBinding),
				},
			},
			{},
		},
		"MORE_CLAUSES": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDot),
					NewSymbol("CLAUSES"),
				},
			},
			{},
		},
		"GROUP_BY": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemGroup),
					NewTokenType(lexer.ItemBy),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("GROUP_BY_BINDINGS"),
				},
			},
			{},
		},
		"GROUP_BY_BINDINGS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("GROUP_BY_BINDINGS"),
				},
			},
			{},
		},
		"ORDER_BY": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemOrder),
					NewTokenType(lexer.ItemBy),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("ORDER_BY_DIRECTION"),
					NewSymbol("ORDER_BY_BINDINGS"),
				},
			},
			{},
		},
		"ORDER_BY_DIRECTION": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAsc),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDesc),
				},
			},
			{},
		},
		"ORDER_BY_BINDINGS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemBinding),
					NewSymbol("ORDER_BY_DIRECTION"),
					NewSymbol("ORDER_BY_BINDINGS"),
				},
			},
			{},
		},
		"HAVING": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemHaving),
					NewSymbol("HAVING_CLAUSE"),
				},
			},
			{},
		},
		"HAVING_CLAUSE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("HAVING_CLAUSE_BINARY_COMPOSITE"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNot),
					NewSymbol("HAVING_CLAUSE"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLPar),
					NewSymbol("HAVING_CLAUSE"),
					NewTokenType(lexer.ItemRPar),
					NewSymbol("HAVING_CLAUSE_BINARY_COMPOSITE"),
				},
			},
		},
		"HAVING_CLAUSE_BINARY_COMPOSITE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAnd),
					NewSymbol("HAVING_CLAUSE"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemOr),
					NewSymbol("HAVING_CLAUSE"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemEQ),
					NewSymbol("HAVING_CLAUSE"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLT),
					NewSymbol("HAVING_CLAUSE"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemGT),
					NewSymbol("HAVING_CLAUSE"),
				},
			},
			{},
		},
		"GLOBAL_TIME_BOUND": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBefore),
					NewTokenType(lexer.ItemPredicate),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemAfter),
					NewTokenType(lexer.ItemPredicate),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBetween),
					NewTokenType(lexer.ItemPredicate),
					NewTokenType(lexer.ItemComma),
					NewTokenType(lexer.ItemPredicate),
				},
			},
			{},
		},
		"LIMIT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLimit),
					NewTokenType(lexer.ItemLiteral),
				},
			},
			{},
		},
		"INSERT_OBJECT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNode),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicate),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLiteral),
				},
			},
		},
		"INSERT_DATA": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDot),
					NewTokenType(lexer.ItemNode),
					NewTokenType(lexer.ItemPredicate),
					NewSymbol("INSERT_OBJECT"),
					NewSymbol("INSERT_DATA"),
				},
			},
			{},
		},
		"DELETE_OBJECT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNode),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicate),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLiteral),
				},
			},
		},
		"DELETE_DATA": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDot),
					NewTokenType(lexer.ItemNode),
					NewTokenType(lexer.ItemPredicate),
					NewSymbol("DELETE_OBJECT"),
					NewSymbol("DELETE_DATA"),
				},
			},
			{},
		},
		"CONSTRUCT_FACTS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLBracket),
					NewSymbol("CONSTRUCT_TRIPLES"),
					NewTokenType(lexer.ItemRBracket),
				},
			},
		},
		"CONSTRUCT_TRIPLES": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNode),
					NewSymbol("CONSTRUCT_PREDICATE"),
					NewSymbol("CONSTRUCT_OBJECT"),
					NewSymbol("MORE_CONSTRUCT_PREDICATE_OBJECT_PAIRS"),
					NewSymbol("MORE_CONSTRUCT_TRIPLES"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBlankNode),
					NewSymbol("CONSTRUCT_PREDICATE"),
					NewSymbol("CONSTRUCT_OBJECT"),
					NewSymbol("MORE_CONSTRUCT_PREDICATE_OBJECT_PAIRS"),
					NewSymbol("MORE_CONSTRUCT_TRIPLES"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("CONSTRUCT_PREDICATE"),
					NewSymbol("CONSTRUCT_OBJECT"),
					NewSymbol("MORE_CONSTRUCT_PREDICATE_OBJECT_PAIRS"),
					NewSymbol("MORE_CONSTRUCT_TRIPLES"),
				},
			},
		},
		"CONSTRUCT_PREDICATE": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicate),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
				},
			},
		},
		"CONSTRUCT_OBJECT": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNode),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBlankNode),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemPredicate),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLiteral),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
				},
			},
		},
		"MORE_CONSTRUCT_PREDICATE_OBJECT_PAIRS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemSemicolon),
					NewSymbol("CONSTRUCT_PREDICATE"),
					NewSymbol("CONSTRUCT_OBJECT"),
					NewSymbol("MORE_CONSTRUCT_PREDICATE_OBJECT_PAIRS"),
				},
			},
			{},
		},
		"MORE_CONSTRUCT_TRIPLES": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDot),
					NewSymbol("CONSTRUCT_TRIPLES"),
				},
			},
			{},
		},
		"DECONSTRUCT_FACTS": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemLBracket),
					NewSymbol("DECONSTRUCT_TRIPLES"),
					NewTokenType(lexer.ItemRBracket),
				},
			},
		},
		"DECONSTRUCT_TRIPLES": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemNode),
					NewSymbol("CONSTRUCT_PREDICATE"),
					NewSymbol("CONSTRUCT_OBJECT"),
					NewSymbol("MORE_DECONSTRUCT_TRIPLES"),
				},
			},
			{
				Elements: []Element{
					NewTokenType(lexer.ItemBinding),
					NewSymbol("CONSTRUCT_PREDICATE"),
					NewSymbol("CONSTRUCT_OBJECT"),
					NewSymbol("MORE_DECONSTRUCT_TRIPLES"),
				},
			},
		},
		"MORE_DECONSTRUCT_TRIPLES": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemDot),
					NewSymbol("DECONSTRUCT_TRIPLES"),
				},
			},
			{},
		},
		"GRAPH_SHOW": []*Clause{
			{
				Elements: []Element{
					NewTokenType(lexer.ItemGraphs),
				},
			},
		},
	}
}
