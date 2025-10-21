package identitydsl

import (
	"fmt"
	"strings"
)

func Check(input string) {
	l := lexer{
		input: input,
	}

	l.run(lexDSL)

	for i := range l.items {
		fmt.Println(strings.Replace(fmt.Sprintf("%#v", l.items[i]), "identitydsl.lexeme", "", -1))
	}
}
