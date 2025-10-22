package identitydsl

const valueRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-.@£$"

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

	if l.acceptString("Account") && (l.peek() == eof || l.accept("\r\n")) {
		return l.errorf("Account not specified on line %d", l.items.currentLineNumber())
	}

	if l.peekString("User ") {
		return lexUser
	}

	if l.acceptString("User") && (l.peek() == eof || l.accept("\r\n")) {
		return l.errorf("User not specified on line %d", l.items.currentLineNumber())
	}

	if l.peekString("Group ") {
		return lexGroup
	}

	if l.acceptString("Group") && (l.peek() == eof || l.accept("\r\n")) {
		return l.errorf("Group not specified on line %d", l.items.currentLineNumber())
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
	l.ignore()

	for pos := 1; ; pos++ {
		if !l.acceptRun("1234567890") {
			return l.errorf("Invalid account ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		if l.width != 12 {
			return l.errorf("Bad length account ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if l.peek() == eof {
			return lexDSL
		}

		if r := l.peek(); r == '\r' || r == '\n' {
			l.acceptRun("\r\n")
			l.emit(typeEOL)
			break
		}
	}

	return lexTagsOrLabels
}

func lexGroup(l *lexer) stateFunc {
	l.acceptString("Group")
	l.ignore()
	l.emit(typeGroup)
	l.acceptRun(" ")
	l.ignore()

	for pos := 1; ; pos++ {
		if !l.acceptRun(valueRunes) {
			return l.errorf("Invalid group ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if l.peek() == eof {
			return lexDSL
		}

		if r := l.peek(); r == '\r' || r == '\n' {
			l.acceptRun("\r\n")
			l.emit(typeEOL)
			break
		}
	}

	return lexTagsOrLabels
}

func lexUser(l *lexer) stateFunc {
	l.acceptString("User")
	l.ignore()
	l.emit(typeUser)
	l.acceptRun(" ")
	l.ignore()

	for pos := 1; ; pos++ {
		if !l.acceptRun(valueRunes) {
			return l.errorf("Invalid user ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if l.peek() == eof {
			return lexDSL
		}

		if r := l.peek(); r == '\r' || r == '\n' {
			l.acceptRun("\r\n")
			l.emit(typeEOL)
			break
		}
	}

	return lexTagsOrLabels
}

func lexTagsOrLabels(l *lexer) stateFunc {
	if !l.acceptRun("\t") {
		return lexDSL
	}

	l.emit(typeSpace)

	for i := 0; i < 2; i++ {
		if l.peek() == '"' {

			l.next()
			l.ignore()

			if l.peek() == '"' {
				return l.errorf("Empty value on line %d", l.items.currentLineNumber())
			}

			l.acceptRun(valueRunes + " ")

			switch r := l.peek(); r {
			case '"':
				l.emit(typeValue)
				l.next()
				l.ignore()

			case '\r', '\n':
				return l.errorf("Unclosed quoted value on line %d", l.items.currentLineNumber())

			default:
				return l.errorf("Invalid character %s on line %d", string(r), l.items.currentLineNumber())
			}
		} else if l.acceptRun(valueRunes) {
			l.emit(typeValue)
		}

		if !l.acceptRun(" ") {
			break
		}

		l.ignore()
	}

	switch l.peek() {
	case eof:
		return lexDSL
	case '\r', '\n':
		l.acceptRun("\r\n")
		l.emit(typeEOL)
	}

	return lexTagsOrLabels
}
