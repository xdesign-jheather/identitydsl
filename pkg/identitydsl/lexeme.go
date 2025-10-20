package identitydsl

type lexemeType int

const (
	typeEOF lexemeType = iota
	typeComment
	typeEOL
)

type lexeme struct {
	typ lexemeType
	val string
}
