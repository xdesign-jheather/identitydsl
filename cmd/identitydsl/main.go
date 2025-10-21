package main

import (
	"log"
	"os"

	"github.com/xdesign-jheather/identitydsl/pkg/identitydsl"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	identitydsl.Check(string(data))
}
