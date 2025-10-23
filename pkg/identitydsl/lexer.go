package identitydsl

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const eof = rune(-1)

type lexer struct {
	input string
	items lexemes
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

	if r == eof {
		return eof
	}

	l.backup()

	return r
}

func (l *lexer) peekString(test string) bool {
	return strings.HasPrefix(l.input[l.pos:], test)
}

func (l *lexer) accept(runes string) bool {
	if strings.IndexRune(runes, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(runes string) bool {
	was := l.pos
	for {
		r := l.next()

		if r == eof {
			break
		}

		if strings.IndexRune(runes, r) < 0 {
			l.backup()
			break
		}
	}
	return l.pos > was
}

func (l *lexer) acceptToLineEnding() {
	for {
		r := l.next()

		if r == eof {
			return
		}

		if r == '\r' || r == '\n' {
			l.backup()
			return
		}
	}
}

func (l *lexer) acceptString(test string) bool {
	if !l.peekString(test) {
		return false
	}
	l.pos += len(test)
	return true
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
	l.width = 0
}

func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	l.items = append(l.items, lexeme{
		typeError,
		fmt.Sprintf(format, args...),
	})
	return nil
}
