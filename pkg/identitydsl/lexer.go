package identitydsl

import (
	"strings"
	"unicode/utf8"
)

const eof = rune(-1)

type lexer struct {
	input string
	items []lexeme
	start int
	pos   int
	width int
}

func (l *lexer) run(start stateFunc) {
	for state := start; state != nil; {
		state = state(l)
	}
}

func (l *lexer) ignore() {
	l.start = l.pos
	l.width = 0
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()

	l.backup()

	return r
}

func (l *lexer) peekString(test string) bool {
	return strings.HasPrefix(l.input[l.pos:], test)
}

func (l *lexer) acceptLine() {
	for {
		r := l.next()

		if r == '\r' || r == '\n' {
			l.backup()
			return
		}

		if r == eof {
			return
		}
	}
}

func (l *lexer) value() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) emit(typ lexemeType) {
	l.items = append(l.items, lexeme{
		typ: typ,
		val: l.value(),
	})
	l.start = l.pos
}
