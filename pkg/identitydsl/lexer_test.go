package identitydsl

import (
	"fmt"
	"testing"
)

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
			t.Errorf("got %d items, want %d", len(got), len(want))

			for i := range l.items {
				fmt.Printf("%d: %v\n", i, l.items[i])
			}

			return
		}

		for i := range got {
			if got[i].typ != want[i].typ || got[i].val != want[i].val {
				t.Errorf("at pos %d, got %v, want %v", i, got[i], want[i])
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

		testLexer(
			t,
			"// A comment line starts with two slashes\n// Another comment!",
			[]lexeme{
				{
					typ: typeComment,
					val: "// A comment line starts with two slashes",
				},
				{
					typ: typeEOL,
					val: "\n",
				},
				{
					typ: typeComment,
					val: "// Another comment!",
				},
				{
					typ: typeEOF,
				},
			},
		)
	})

	t.Run("new lines", func(t *testing.T) {

		testLexer(
			t,
			"\n",
			[]lexeme{
				{
					typ: typeEOL,
					val: "\n",
				},
				{
					typ: typeEOF,
				},
			},
		)

		testLexer(
			t,
			"\r",
			[]lexeme{
				{
					typ: typeEOL,
					val: "\r",
				},
				{
					typ: typeEOF,
				},
			},
		)

		testLexer(
			t,
			"\r\n",
			[]lexeme{
				{
					typ: typeEOL,
					val: "\r\n",
				},
				{
					typ: typeEOF,
				},
			},
		)

		// Multiple newlines should be one lexeme

		testLexer(
			t,
			"\n\n",
			[]lexeme{
				{
					typ: typeEOL,
					val: "\n\n",
				},
				{
					typ: typeEOF,
				},
			},
		)
	})

	t.Run("unknown input", func(t *testing.T) {
		testLexer(
			t,
			"Hello",
			[]lexeme{
				{
					typ: typeError,
					val: "Unknown input 'Hello' at line 1",
				},
			},
		)

		testLexer(
			t,
			"\nCheese",
			[]lexeme{
				{
					typ: typeEOL,
					val: "\n",
				},
				{
					typ: typeError,
					val: "Unknown input 'Cheese' at line 2",
				},
			},
		)
	})

	t.Run("account entity", func(t *testing.T) {
		testLexer(
			t,
			"Account",
			[]lexeme{
				{
					typ: typeError,
					val: "Unknown input 'Account' at line 1",
				},
			},
		)

		testLexer(
			t,
			"Account 112233445566",
			[]lexeme{
				{
					typ: typeAccount,
				},
				{
					typ: typeSpace,
					val: " ",
				},
				{
					typ: typeIdentifier,
					val: "112233445566",
				},
				{
					typ: typeEOF,
				},
			},
		)

		testLexer(
			t,
			"Account 1122334455",
			[]lexeme{
				{
					typ: typeAccount,
				},
				{
					typ: typeSpace,
					val: " ",
				},
				{
					typ: typeError,
					val: "Bad length account ID on line 1",
				},
			},
		)

		testLexer(
			t,
			"Account ABC1122334455",
			[]lexeme{
				{
					typ: typeAccount,
				},
				{
					typ: typeSpace,
					val: " ",
				},
				{
					typ: typeError,
					val: "Invalid account ID on line 1",
				},
			},
		)
	})
}
