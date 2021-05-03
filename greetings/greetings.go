package greetings

import (
	"errors"
	"fmt"

	"rsc.io/quote/v3"
)

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return "", errors.New("empty name")
	}

	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome! Did you know that \"%v\"", name, quote.GoV3())
	return message, nil
}
