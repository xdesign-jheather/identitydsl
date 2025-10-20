package identitydsl

type stateFunc func(*lexer) stateFunc

func lexDSL(l *lexer) stateFunc {
	for {
		if l.peek() == eof {
			l.next()
			l.emit(typeEOF)
			return nil
		}

		if l.peekString("//") {
			l.acceptLine()
			l.emit(typeComment)
			continue
		}

		l.acceptLine()

		l.ignore()
	}
}
