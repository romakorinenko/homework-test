package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	phrase := "Hello, OTUS!"
	printReversePhrase(phrase)
}

func printReversePhrase(phrase string) {
	fmt.Println(reverse.String(phrase))
}
