package greetings

import (
	"fmt"

	"rsc.io/quote/v3"
)

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome! Did you know that \"%v\"", name, quote.GoV3())
	return message
}
