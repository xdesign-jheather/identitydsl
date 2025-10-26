package identitydsl

const valueRunes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+=.@-"

type stateFunc func(*lexer) stateFunc

func lexDSL(l *lexer) stateFunc {
	if acceptEndOfFile(l) {
		return nil
	}

	if acceptLineEnd(l) {
		return lexDSL
	}

	if l.peekString("//") {
		return lexComment
	}

	if l.peekString("Account") {
		return lexAccount
	}

	if l.peekString("User") {
		return lexUser
	}

	if l.peekString("Group") {
		return lexGroup
	}

	if l.peekString("Role") {
		return lexRole
	}

	if l.peekString("Assign") {
		return lexAssignment
	}

	return lexUnknown
}

func lexAssignment(l *lexer) stateFunc {
	l.pos += len("Assign")

	if l.peek() == eof {
		return l.errorf("Incomplete assignment on line %d", l.items.currentLineNumber())
	}

	l.ignore()

	l.acceptToLineEnding()

	if l.value() != "" {
		return l.errorf("Unexpected '%s' after assignment on line %d", l.value(), l.items.currentLineNumber())
	}

	l.emit(typeAssignment)

	acceptLineEnd(l)

	return lexAssignmentLine
}

func lexAssignmentLine(l *lexer) stateFunc {
	start, pos := l.start, l.pos

	if !l.acceptRun("\t") {
		return lexDSL
	}

	if l.peekString("Group") {
		l.emit(typeSpace)
		return lexAssignmentGroup
	}

	if l.peekString("User") {
		l.emit(typeSpace)
		return lexAssignmentUser
	}

	if l.peekString("Role") {
		l.emit(typeSpace)
		return lexAssignmentRole
	}

	// We didn't find an assignment field

	l.start, l.pos = start, pos

	return lexDSL
}

func lexAssignmentGroup(l *lexer) stateFunc {
	l.pos += len("Group")

	if l.peekAnyOf(eof, '\r', '\n') {
		return l.errorf("Empty assignment group on line %d", l.items.currentLineNumber())
	}

	if l.peek() != ' ' {
		l.acceptToLineEnding()
		return l.errorf("Unknown input '%s' on line %d for assignment", l.value(), l.items.currentLineNumber())
	}

	l.ignore()

	l.emit(typeGroup)

	if l.acceptRun(" ") {
		l.ignore()
	}

	for pos := 1; ; pos++ {
		if !l.acceptRun(valueRunes) {
			return l.errorf("Invalid assignment group ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if acceptEndOfFile(l) {
			return nil
		}

		if acceptLineEnd(l) {
			break
		}
	}

	acceptLineEnd(l)

	return lexAssignmentLine
}

func lexAssignmentUser(l *lexer) stateFunc {
	l.pos += len("User")

	if l.peekAnyOf(eof, '\r', '\n') {
		return l.errorf("Empty assignment user on line %d", l.items.currentLineNumber())
	}

	if l.peek() != ' ' {
		l.acceptToLineEnding()
		return l.errorf("Unknown input '%s' on line %d for assignment", l.value(), l.items.currentLineNumber())
	}

	l.ignore()

	l.emit(typeUser)

	if l.acceptRun(" ") {
		l.ignore()
	}

	for pos := 1; ; pos++ {
		if !l.acceptRun(valueRunes) {
			return l.errorf("Invalid assignment user ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if acceptEndOfFile(l) {
			return nil
		}

		if acceptLineEnd(l) {
			break
		}
	}

	acceptLineEnd(l)

	return lexAssignmentLine
}

func lexAssignmentRole(l *lexer) stateFunc {
	l.pos += len("Role")

	if l.peekAnyOf(eof, '\r', '\n') {
		return l.errorf("Empty assignment role on line %d", l.items.currentLineNumber())
	}

	if l.peek() != ' ' {
		l.acceptToLineEnding()
		return l.errorf("Unknown input '%s' on line %d for assignment", l.value(), l.items.currentLineNumber())
	}

	l.ignore()

	l.emit(typeRole)

	if l.acceptRun(" ") {
		l.ignore()
	}

	for pos := 1; ; pos++ {
		if !l.acceptRun(valueRunes) {
			return l.errorf("Invalid assignment role ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if acceptEndOfFile(l) {
			return nil
		}

		if acceptLineEnd(l) {
			break
		}
	}

	acceptLineEnd(l)

	return lexAssignmentLine
}

func lexUnknown(l *lexer) stateFunc {
	l.acceptToLineEnding()
	return l.errorf("Unknown input '%s' on line %d", l.value(), l.items.currentLineNumber())
}

func lexComment(l *lexer) stateFunc {
	l.acceptToLineEnding()
	l.emit(typeComment)
	return lexDSL
}

func lexAccount(l *lexer) stateFunc {
	l.acceptString("Account")
	l.ignore()

	if l.peekAnyOf(eof, '\r', '\n') {
		return l.errorf("Account not specified on line %d", l.items.currentLineNumber())
	}

	if l.peek() != ' ' {
		l.pos -= len("Account")
		l.start -= len("Account")
		return lexUnknown
	}

	l.emit(typeAccount)
	l.acceptRun(" ")
	l.ignore()

	for pos := 1; ; pos++ {
		if !l.acceptRun("1234567890") {
			return l.errorf("Invalid account ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		if len(l.value()) != 12 {
			return l.errorf("Bad length account ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if acceptEndOfFile(l) {
			return nil
		}

		if acceptLineEnd(l) {
			break
		}
	}

	return lexTagsOrLabels
}

func lexGroup(l *lexer) stateFunc {
	l.acceptString("Group")
	l.ignore()

	if l.peekAnyOf(eof, '\r', '\n') {
		return l.errorf("Group not specified on line %d", l.items.currentLineNumber())
	}

	if l.peek() != ' ' {
		l.pos -= len("Group")
		l.start -= len("Group")
		return lexUnknown
	}

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

		if acceptEndOfFile(l) {
			return nil
		}

		if acceptLineEnd(l) {
			break
		}
	}

	return lexTagsOrLabels
}

func lexUser(l *lexer) stateFunc {
	l.acceptString("User")
	l.ignore()

	if l.peekAnyOf(eof, '\r', '\n') {
		return l.errorf("User not specified on line %d", l.items.currentLineNumber())
	}

	if l.peek() != ' ' {
		l.pos -= len("User")
		l.start -= len("User")
		return lexUnknown
	}

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

		if acceptEndOfFile(l) {
			return nil
		}

		if acceptLineEnd(l) {
			break
		}
	}

	return lexTagsOrLabels
}

func lexRole(l *lexer) stateFunc {
	l.acceptString("Role")
	l.ignore()

	if l.peekAnyOf(eof, '\r', '\n') {
		return l.errorf("Role not specified on line %d", l.items.currentLineNumber())
	}

	if l.peek() != ' ' {
		l.pos -= len("Role")
		l.start -= len("Role")
		return lexUnknown
	}

	l.emit(typeRole)
	l.acceptRun(" ")
	l.ignore()

	for pos := 1; ; pos++ {
		if !l.acceptRun(valueRunes) {
			return l.errorf("Invalid role ID on line %d position %d", l.items.currentLineNumber(), pos)
		}

		l.emit(typeValue)

		if l.acceptRun(", ") {
			l.ignore()
			continue
		}

		if acceptEndOfFile(l) {
			return nil
		}

		if acceptLineEnd(l) {
			break
		}
	}

	return lexPolicies
}

func lexPolicies(l *lexer) stateFunc {
	if acceptEndOfFile(l) {
		return nil
	}

	if !l.acceptRun("\t") {
		return lexDSL
	}

	l.emit(typeSpace)

	if !l.acceptRun(valueRunes) {
		return l.errorf("No policies found on line %d", l.items.currentLineNumber())
	}

	l.emit(typeValue)

	acceptLineEnd(l)

	return lexPolicies
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
