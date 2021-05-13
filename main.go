// main package, initialises the http server
// this is what is executed when you run [go run .] at the project's root
package main

import (
	"log"
	"os"

	server "github.com/LSLarose/httpPlayerServer"
)

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("Tool: ")
	log.SetFlags(0)
	// receive arguments
	consoleLineArguments := os.Args

	//check if enough arguments, if not print help message and exit
	if len(consoleLineArguments) < 2 {
		log.Fatal("The client update server requires 1 argument" +
			"\nThe correct way to use this tool is\n" +
			//first argument is the tool's name
			consoleLineArguments[0] + " [CSV FILE INPUT PATH]\n" +
			"where:\n" +
			"- [CSV FILE INPUT PATH] is the path towards a valid csv file containing the clients to update")
	}

	// only first arg is considered
	pathToCSVFile := consoleLineArguments[1]

	//check if file exists
	_, err := os.Stat(pathToCSVFile)
	if os.IsNotExist(err) {
		log.Fatalf("Input file at does not exist at \"%s\".", pathToCSVFile)
	}

	//start the server
	server.StartServer(pathToCSVFile)
}
