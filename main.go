// main package, initialises the http server
// this is what is executed when you run [go run .] at the project's root
package main

import (
	server "github.com/LSLarose/httpPlayerServer"
)

func main() {
	//start the server
	server.StartServer()
}
