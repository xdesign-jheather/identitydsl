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

				fmt.Println("Got:")
				for i := range l.items {
					fmt.Printf("%d: %#v\n", i, l.items[i])
				}

				fmt.Println("Wanted:")
				for i := range want {
					fmt.Printf("%d: %#v\n", i, want[i])
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
				{typeComment, "// A comment line starts with two slashes"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"multiple",
			"// A comment line starts with two slashes\n// Another comment!",
			[]lexeme{
				{typeComment, "// A comment line starts with two slashes"},
				{typeEOL, "\n"},
				{typeComment, "// Another comment!"},
				{typ: typeEOF},
			},
		)
	})

	t.Run("new lines", func(t *testing.T) {
		lex(
			t,
			"n",
			"\n",
			[]lexeme{
				{typeEOL, "\n"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"r",
			"\r",
			[]lexeme{
				{typeEOL, "\r"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"rn",
			"\r\n",
			[]lexeme{
				{typeEOL, "\r\n"},
				{typ: typeEOF},
			},
		)

		// Multiple newlines should be one lexeme

		lex(
			t,
			"nn",
			"\n\n",
			[]lexeme{
				{typeEOL, "\n\n"},
				{typ: typeEOF},
			},
		)
	})

	t.Run("unknown input", func(t *testing.T) {
		lex(
			t,
			"line 1",
			"Hello",
			[]lexeme{
				{typeError, "Unknown input 'Hello' on line 1"},
			},
		)

		lex(
			t,
			"line 2",
			"\nCheese",
			[]lexeme{
				{typeEOL, "\n"},
				{typeError, "Unknown input 'Cheese' on line 2"},
			},
		)
	})

	t.Run("account entity", func(t *testing.T) {

		lex(
			t,
			"no identifier",
			"Account",
			[]lexeme{
				{typeError, "Unknown input 'Account' on line 1"},
			},
		)

		lex(
			t,
			"valid",
			"Account 112233445566",
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233445566"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"short",
			"Account 1234567890",
			[]lexeme{
				{typ: typeAccount},
				{typeError, "Bad length account ID on line 1 position 1"},
			},
		)

		lex(
			t,
			"invalid",
			"Account Word",
			[]lexeme{
				{typ: typeAccount},
				{typeError, "Invalid account ID on line 1 position 1"},
			},
		)

		lex(
			t,
			"multiple valid",
			"Account 000000000000, 111111111111,  222222222222 , 333333333333",
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "000000000000"},
				{typeValue, "111111111111"},
				{typeValue, "222222222222"},
				{typeValue, "333333333333"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"valid then invalid",
			"Account 000000000000, Bob,  222222222222 , 333333333333",
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "000000000000"},
				{typeError, "Invalid account ID on line 1 position 2"},
			},
		)

		lex(
			t,
			"basic label",
			`Account 112233112233
	Label1`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Label1"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"quoted label",
			`Account 112233112233
	"Developer Access"`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Developer Access"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"key value pair",
			`Account 112233112233
	Key1 Value1`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Key1"},
				{typeValue, "Value1"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"key value pair quoted key",
			`Account 112233112233
	"Hello World" Value1`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Hello World"},
				{typeValue, "Value1"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"key value pair quoted value",
			`Account 112233112233
	Name "Hello World"`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Name"},
				{typeValue, "Hello World"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"key value pair quoted both",
			`Account 112233112233
	"What a World" "Hello World"`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "What a World"},
				{typeValue, "Hello World"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"multiple labels",
			`Account 112233112233
	Label1
	Label2
	"Label 3"`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Label1"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Label2"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Label 3"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"multiple tags",
			`Account 112233112233
	Name Jonathan
	Age 36
	"Favorite Pudding" "Rhubarb Crumble"`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Name"},
				{typeValue, "Jonathan"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Age"},
				{typeValue, "36"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Favorite Pudding"},
				{typeValue, "Rhubarb Crumble"},
				{typ: typeEOF},
			},
		)

		lex(
			t,
			"tags and labels mixed",
			`Account 112233112233
	Billing
	Organisations
	Owner Platform

	Product Radio`,
			[]lexeme{
				{typ: typeAccount},
				{typeValue, "112233112233"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Billing"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Organisations"},
				{typeEOL, "\n"},
				{typeSpace, "\t"},
				{typeValue, "Owner"},
				{typeValue, "Platform"},
				{typeEOL, "\n\n"},
				{typeSpace, "\t"},
				{typeValue, "Product"},
				{typeValue, "Radio"},
				{typ: typeEOF},
			},
		)
	})
}
