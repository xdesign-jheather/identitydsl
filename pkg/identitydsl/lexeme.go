package identitydsl

type lexemeType int

const (
	typeError lexemeType = iota
	typeEOF
	typeEOL
	typeComment
	typeSpace
	typeValue
	typeAccount
	typeGroup
	typeUser
	typeRole
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
