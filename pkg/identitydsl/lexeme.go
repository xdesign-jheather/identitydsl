package identitydsl

type lexemeType int

const (
	typeEOF lexemeType = iota
	typeComment
)

type lexeme struct {
	typ lexemeType
	val string
}
