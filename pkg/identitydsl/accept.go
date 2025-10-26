package identitydsl

func acceptLineEnd(l *lexer) bool {
	if l.acceptRun("\r\n") {
		l.emit(typeEOL)
		return true
	}

	return false
}

func acceptEndOfFile(l *lexer) bool {
	if l.peek() == eof {
		l.emit(typeEOF)
		return true
	}

	return false
}
