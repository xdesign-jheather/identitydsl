package identitydsl

import (
	"fmt"
	"testing"
)

func TestLex(t *testing.T) {
	// We intend to test that the lexer produces the correct token stream for
	// all scenarios in the README.

	lex := func(t *testing.T, name, input string, want []lexeme) {
		t.Run(name, func(t *testing.T) {
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
		})
	}

	lex(
		t,
		"empty file",
		"",
		[]lexeme{
			{
				typ: typeEOF,
			},
		},
	)

	t.Run("comments", func(t *testing.T) {
		lex(
			t,
			"single",
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

		lex(
			t,
			"multiple",
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
		lex(
			t,
			"n",
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

		lex(
			t,
			"r",
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

		lex(
			t,
			"rn",
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

		lex(
			t,
			"nn",
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
		lex(
			t,
			"line 1",
			"Hello",
			[]lexeme{
				{
					typ: typeError,
					val: "Unknown input 'Hello' on line 1",
				},
			},
		)

		lex(
			t,
			"line 2",
			"\nCheese",
			[]lexeme{
				{
					typ: typeEOL,
					val: "\n",
				},
				{
					typ: typeError,
					val: "Unknown input 'Cheese' on line 2",
				},
			},
		)
	})

	t.Run("account entity", func(t *testing.T) {
		lex(
			t,
			"no identifier",
			"Account",
			[]lexeme{
				{
					typ: typeError,
					val: "Unknown input 'Account' on line 1",
				},
			},
		)

		lex(
			t,
			"valid",
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

		lex(
			t,
			"short",
			"Account 1234567890",
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
					val: "Bad length account ID on line 1 position 1",
				},
			},
		)

		lex(
			t,
			"invalid",
			"Account Word",
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
					val: "Invalid account ID on line 1 position 1",
				},
			},
		)

		lex(
			t,
			"multiple valid",
			"Account 000000000000, 111111111111,  222222222222 , 333333333333",
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
					val: "000000000000",
				},
				{
					typ: typeDelimiter,
					val: ", ",
				},
				{
					typ: typeIdentifier,
					val: "111111111111",
				},
				{
					typ: typeDelimiter,
					val: ",  ",
				},
				{
					typ: typeIdentifier,
					val: "222222222222",
				},
				{
					typ: typeDelimiter,
					val: " , ",
				},
				{
					typ: typeIdentifier,
					val: "333333333333",
				},
				{
					typ: typeEOF,
				},
			},
		)

		lex(
			t,
			"valid then invalid",
			"Account 000000000000, Bob,  222222222222 , 333333333333",
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
					val: "000000000000",
				},
				{
					typ: typeDelimiter,
					val: ", ",
				},
				{
					typ: typeError,
					val: "Invalid account ID on line 1 position 2",
				},
			},
		)
	})
}
