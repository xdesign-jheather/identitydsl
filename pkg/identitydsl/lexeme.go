package identitydsl

type lexemeType int

const (
	typeEOF lexemeType = iota
	typeComment
	typeEOL
	typeError
)

type lexeme struct {
	typ lexemeType
	val string
}

type lexemes []lexeme

func (l lexemes) currentLineNumber() int {
	number := 1

	for i := range l {
		if l[i].typ == typeEOL {
			number += len(l[i].val)
		}
	}

	return number
}
