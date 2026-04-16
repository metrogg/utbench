func (sc *scanner) nextToken(val *tokenValue) Token {

	// The following distribution of tokens guides case ordering:
	//
	//      COMMA          27   %
	//      STRING         23   %
	//      IDENT          15   %
	//      EQL            11   %
	//      LBRACK          5.5 %
	//      RBRACK          5.5 %
	//      NEWLINE         3   %
	//      LPAREN          2.9 %
	//      RPAREN          2.9 %
	//      INT             2   %
	//      others        < 1   %
	//
	// Although NEWLINE tokens are infrequent, and lineStart is
	// usually (~97%) false on entry, skipped newlines account for
	// about 50% of all iterations of the 'start' loop.

start:
	var c rune

	// Deal with leading spaces and indentation.
	blank := false
	savedLineStart := sc.lineStart
	if sc.lineStart {
		sc.lineStart = false
		col := 0
		for {
			c = sc.peekRune()
			if c == ' ' {
				col++
				sc.readRune()
			} else if c == '\t' {
				const tab = 8
				col += int(tab - (sc.pos.Col-1)%tab)
				sc.readRune()
			} else {
				break
			}
		}
		// The third clause is "trailing spaces without newline at EOF".
		if c == '#' || c == '\n' || c == 0 && col > 0 {
			blank = true
		}

		// Compute indentation level for non-blank lines not
		// inside an expression.  This is not the common case.
		if !blank && sc.depth == 0 {
			cur := sc.indentstk[len(sc.indentstk)-1]
			if col > cur {
				// indent
				sc.dents++
				sc.indentstk = append(sc.indentstk, col)
			} else if col < cur {
				// dedent(s)
				for len(sc.indentstk) > 0 && col < sc.indentstk[len(sc.indentstk)-1] {
					sc.dents--
					sc.indentstk = sc.indentstk[:len(sc.indentstk)-1] // pop
				}
				if col != sc.indentstk[len(sc.indentstk)-1] {
					sc.error(sc.pos, "unindent does not match any outer indentation level")
				}
			}
		}
	}

	// Return saved indentation tokens.
	if sc.dents != 0 {
		sc.startToken(val)
		sc.endToken(val)
		if sc.dents < 0 {
			sc.dents++
			return OUTDENT
		} else {
			sc.dents--
			return INDENT
		}
	}

	// start of line proper
	c = sc.peekRune()

	// Skip spaces.
	for c == ' ' || c == '\t' {
		sc.readRune()
		c = sc.peekRune()
	}

	// comment
	if c == '#' {
		if sc.keepComments {
			sc.startToken(val)
		}
		// Consume up to newline (included).
		for c != 0 && c != '\n' {
			sc.readRune()
			c = sc.peekRune()
		}
		if sc.keepComments {
			sc.endToken(val)
			if blank {
				sc.lineComments = append(sc.lineComments, Comment{val.pos, val.raw})
			} else {
				sc.suffixComments = append(sc.suffixComments, Comment{val.pos, val.raw})
			}
		}
	}

	// newline
	if c == '\n' {
		sc.lineStart = true
		if blank || sc.depth > 0 {
			// Ignore blank lines, or newlines within expressions (common case).
			sc.readRune()
			goto start
		}
		// At top-level (not in an expression).
		sc.startToken(val)
		sc.readRune()
		val.raw = "\n"
		return NEWLINE
	}

	// end of file
	if c == 0 {
		// Emit OUTDENTs for unfinished indentation,
		// preceded by a NEWLINE if we haven't just emitted one.
		if len(sc.indentstk) > 1 {
			if savedLineStart {
				sc.dents = 1 - len(sc.indentstk)
				sc.indentstk = sc.indentstk[1:]
				goto start
			} else {
				sc.lineStart = true
				sc.startToken(val)
				val.raw = "\n"
				return NEWLINE
			}
		}

		sc.startToken(val)
		sc.endToken(val)
		return EOF
	}

	// line continuation
	if c == '\\' {
		sc.readRune()
		if sc.peekRune() != '\n' {
			sc.errorf(sc.pos, "stray backslash in program")
		}
		sc.readRune()
		goto start
	}

	// start of the next token
	sc.startToken(val)

	// comma (common case)
	if c == ',' {
		sc.readRune()
		sc.endToken(val)
		return COMMA
	}

	// string literal
	if c == '"' || c == '\'' {
		return sc.scanString(val, c)
	}

	// identifier or keyword
	if isIdentStart(c) {
		// raw string literal
		if c == 'r' && len(sc.rest) > 1 && (sc.rest[1] == '"' || sc.rest[1] == '\'') {
			sc.readRune()
			c = sc.peekRune()
			return sc.scanString(val, c)
		}

		for isIdent(c) {
			sc.readRune()
			c = sc.peekRune()
		}
		sc.endToken(val)
		if k, ok := keywordToken[val.raw]; ok {
			return k
		}

		return IDENT
	}

	// brackets
	switch c {
	case '[', '(', '{':
		sc.depth++
		sc.readRune()
		sc.endToken(val)
		switch c {
		case '[':
			return LBRACK
		case '(':
			return LPAREN
		case '{':
			return LBRACE
		}
		panic("unreachable")

	case ']', ')', '}':
		if sc.depth == 0 {
			sc.error(sc.pos, "indentation error")
		} else {
			sc.depth--
		}
		sc.readRune()
		sc.endToken(val)
		switch c {
		case ']':
			return RBRACK
		case ')':
			return RPAREN
		case '}':
			return RBRACE
		}
		panic("unreachable")
	}

	// int or float literal, or period
	if isdigit(c) || c == '.' {
		return sc.scanNumber(val, c)
	}

	// other punctuation
	defer sc.endToken(val)
	switch c {
	case '=', '<', '>', '!', '+', '-', '%', '/', '&', '|', '^', '~': // possibly followed by '='
		start := sc.pos
		sc.readRune()
		if sc.peekRune() == '=' {
			sc.readRune()
			switch c {
			case '<':
				return LE
			case '>':
				return GE
			case '=':
				return EQL
			case '!':
				return NEQ
			case '+':
				return PLUS_EQ
			case '-':
				return MINUS_EQ
			case '/':
				return SLASH_EQ
			case '%':
				return PERCENT_EQ
			case '&':
				return AMP_EQ
			case '|':
				return PIPE_EQ
			case '^':
				return CIRCUMFLEX_EQ
			}
		}
		switch c {
		case '=':
			return EQ
		case '<':
			if sc.peekRune() == '<' {
				sc.readRune()
				if sc.peekRune() == '=' {
					sc.readRune()
					return LTLT_EQ
				} else {
					return LTLT
				}
			}
			return LT
		case '>':
			if sc.peekRune() == '>' {
				sc.readRune()
				if sc.peekRune() == '=' {
					sc.readRune()
					return GTGT_EQ
				} else {
					return GTGT
				}
			}
			return GT
		case '!':
			sc.error(start, "unexpected input character '!'")
		case '+':
			return PLUS
		case '-':
			return MINUS
		case '/':
			if sc.peekRune() == '/' {
				sc.readRune()
				if sc.peekRune() == '=' {
					sc.readRune()
					return SLASHSLASH_EQ
				} else {
					return SLASHSLASH
				}
			}
			return SLASH
		case '%':
			return PERCENT
		case '&':
			return AMP
		case '|':
			return PIPE
		case '^':
			return CIRCUMFLEX
		case '~':
			return TILDE
		}
		panic("unreachable")

	case ':', ';': // single-char tokens (except comma)
		sc.readRune()
		switch c {
		case ':':
			return COLON
		case ';':
			return SEMI
		}
		panic("unreachable")

	case '*': // possibly followed by '*' or '='
		sc.readRune()
		switch sc.peekRune() {
		case '*':
			sc.readRune()
			return STARSTAR
		case '=':
			sc.readRune()
			return STAR_EQ
		}
		return STAR
	}

	sc.errorf(sc.pos, "unexpected input character %#q", c)
	panic("unreachable")
}
