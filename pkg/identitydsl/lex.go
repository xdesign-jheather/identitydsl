package identitydsl

type stateFunc func(*lexer) stateFunc

func lexDSL(l *lexer) stateFunc {
	for {
		if l.peek() == eof {
			l.next()
			l.emit(typeEOF)
			return nil
		}

		if l.acceptRun("\r\n") {
			l.emit(typeEOL)
			continue
		}

		if l.peekString("//") {
			return lexComment
		}

		if l.peekString("Account ") {
			return lexAccount
		}

		l.acceptLine()

		return l.errorf("Unknown input '%s' at line %d", l.value(), l.items.currentLineNumber())
	}
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

	if !l.acceptRun("1234567890") {
		return l.errorf("Invalid account ID on line %d", l.items.currentLineNumber())
	}

	if l.width != 12 {
		return l.errorf("Bad length account ID on line %d", l.items.currentLineNumber())
	}

	l.emit(typeIdentifier)

	return lexDSL
}
