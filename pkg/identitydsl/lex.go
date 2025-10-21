package identitydsl

type stateFunc func(*lexer) stateFunc

func lexDSL(l *lexer) stateFunc {
	if l.peek() == eof {
		l.emit(typeEOF)
		return nil
	}

	if l.acceptRun("\r\n") {
		l.emit(typeEOL)
		return lexDSL
	}

	if l.peekString("//") {
		return lexComment
	}

	if l.peekString("Account ") {
		return lexAccount
	}

	return lexUnknown
}

func lexUnknown(l *lexer) stateFunc {
	l.acceptLine()
	return l.errorf("Unknown input '%s' on line %d", l.value(), l.items.currentLineNumber())
}

func lexComment(l *lexer) stateFunc {
	l.acceptLine()
	l.emit(typeComment)
	return lexDSL
}

func lexAccount(l *lexer) stateFunc {
	l.acceptString("Account")
	l.ignore()
	l.emit(typeAccount)
	l.acceptRun(" ")
	l.emit(typeSpace)

	for pos := 1; ; pos++ {
		if !l.acceptRun("1234567890") {
			return l.errorf("Invalid account ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		if l.width != 12 {
			return l.errorf("Bad length account ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeIdentifier)

		if l.acceptRun(", ") {
			l.emit(typeDelimiter)
			continue
		}

		if r := l.peek(); r == eof || r == '\r' || r == '\n' {
			return lexDSL
		}
	}
}
