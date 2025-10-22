package identitydsl

import "testing"

func TestLexemesCurrentLineNumber(t *testing.T) {
	t.Run("no lexemes", func(t *testing.T) {
		l := lexemes{}

		got, want := l.currentLineNumber(), 1

		if got != want {
			t.Errorf("got line number %d, want %d", got, want)
		}
	})

	t.Run("couple of newlines", func(t *testing.T) {
		l := lexemes{
			{
				typ: typeComment,
				val: "// Hi",
			},
			{
				typ: typeEOL,
				val: "\n",
			},
			{
				typ: typeComment,
				val: "// Hi",
			},
			{
				typ: typeEOL,
				val: "\n",
			},
		}

		got, want := l.currentLineNumber(), 3

		if got != want {
			t.Errorf("got line number %d, want %d", got, want)
		}
	})
}
