package identitydsl

import "testing"

func TestLexer(t *testing.T) {
	// We intend to test that the lexer produces the correct token stream for
	// all scenarios in the README.

	testLexer := func(t *testing.T, input string, want []lexeme) {
		l := lexer{
			input: input,
		}

		l.run(lexDSL)

		got := l.items

		if len(got) != len(want) {
			t.Fatalf("got %d items, want %d", len(got), len(want))
		}

		for i := range got {
			if got[i].typ != want[i].typ || got[i].val != want[i].val {
				t.Fatalf("at pos %d, got %v, want %v", i, got[i], want[i])
				return
			}
		}
	}

	t.Run("empty file", func(t *testing.T) {
		testLexer(
			t,
			"",
			[]lexeme{
				{
					typ: typeEOF,
				},
			},
		)
	})

	t.Run("comments", func(t *testing.T) {
		testLexer(
			t,
			"// A comment line starts with two slashes",
			[]lexeme{
				{
					typ: typeComment,
					val: "// A comment line starts with two slashes",
				},
				{
					typ: typeEOF,
				},
			},
		)
	})
}
