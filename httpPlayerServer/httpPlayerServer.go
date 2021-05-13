// This package implements the server's http functionnalities
package httpPlayerServer

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/LSLarose/greetings"
	"github.com/gorilla/mux"
)

const (
	HTTP_ADDR      = ":9000"
	HTTPS_ADDR     = ":9001"
	CERT_FILE_PATH = "assets/localhost"
	// TODO: not correct but will do, just matches all non-space characters. Would probably be unusable if there were longer API routes?
	MAC_ADDR_PATTERN = "[\\S]+"
	MAIN_API_ROUTE   = "/profiles/clientId:{macaddress:" + MAC_ADDR_PATTERN + "}"
)

// This starts the HTTPS server at
func StartServer(pathToCSV string) {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("server: ")
	log.SetFlags(0)
	log.Println("starting server...")

	srv := http.Server{
		Addr:    HTTPS_ADDR,
		Handler: buildHandlers(),
	}

	_, tlsPort, err := net.SplitHostPort(HTTPS_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	// launch an http => https redirect server in a goroutine
	go redirectToHTTPS(tlsPort)

	// launch expected https server
	err = srv.ListenAndServeTLS(CERT_FILE_PATH+".crt", CERT_FILE_PATH+".key")

	// print any server error
	if err != http.ErrServerClosed {
		log.Fatal(err)
	} else {
		// server closed normaly
		log.Println("Server closed.")
	}
}

func buildHandlers() *http.ServeMux {
	// create a new `ServeMux`
	serveMux := http.NewServeMux()

	// create gorilla mux for complex routing
	gorillaMux := mux.NewRouter()

	// handle expected API route
	gorillaMux.HandleFunc(MAIN_API_ROUTE, func(res http.ResponseWriter, req *http.Request) {
		handleMacAddressUpdateRequest(res, req)
	})

	// handle any other route as a 404
	gorillaMux.HandleFunc("*", func(res http.ResponseWriter, req *http.Request) {
		http.Error(res, "Page not found", http.StatusNotFound)
	})

	// defer all routing to gorillaMux
	serveMux.Handle("/", gorillaMux)

	return serveMux
}

// shamelessly taken from https://stackoverflow.com/questions/37536006/
// reroutes any http request to the equivalent route on the https server
func redirectToHTTPS(tlsPort string) {
	httpSrv := http.Server{
		Addr: HTTP_ADDR,
		Handler: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			host, _, _ := net.SplitHostPort(req.Host)
			targetURL := req.URL
			targetURL.Host = net.JoinHostPort(host, tlsPort)
			targetURL.Scheme = "https"
			http.Redirect(res, req, targetURL.String(), http.StatusMovedPermanently)
		}),
	}
	log.Println(httpSrv.ListenAndServe())
}

// handler request, simple
func handleMacAddressUpdateRequest(res http.ResponseWriter, req *http.Request) {
	//extract macaddress from URL
	macAddress := mux.Vars(req)["macaddress"]

	message, err := greetings.Hello(macAddress)
	// If an error was returned, print it to the console and
	// exit the program.
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	// If no error was returned, print the returned message
	// to the console.
	fmt.Fprint(res, message)
}
