// This package implements the server's http functionnalities
package httpPlayerServer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LSLarose/greetings"
)

func StartServer() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("server: ")
	log.SetFlags(0)
	log.Default().Print("starting server...")
	// create a new `ServeMux`
	mux := http.NewServeMux()

	// handle `/` route
	mux.HandleFunc("/", helloworld)

	// handle `/hello/golang` route
	mux.HandleFunc("/hello/golang", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "Hello Golang!")
	})

	// listen and serve using `ServeMux`
	err := http.ListenAndServeTLS(":9000", "assets/localhost.crt", "assets/localhost.key", mux)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Default().Print("Server up.")
	}
}

func helloworld(res http.ResponseWriter, req *http.Request) {
	message, err := greetings.Hello("Louis Sérey")
	// If an error was returned, print it to the console and
	// exit the program.
	if err != nil {
		http.Error(res, err.Error(), 500)
		log.Fatal(err)
	}

	// If no error was returned, print the returned message
	// to the console.
	fmt.Fprint(res, message)
}
