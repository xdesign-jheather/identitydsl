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
			l.acceptLine()
			l.emit(typeComment)
			continue
		}

		l.acceptLine()

		return l.errorf("Unknown input '%s' at line %d", l.value(), l.items.currentLineNumber())
	}
}
