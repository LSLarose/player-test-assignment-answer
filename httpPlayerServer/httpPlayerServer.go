// This package implements the server's http functionnalities
package httpPlayerServer

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const (
	// ports for the servers
	HTTP_ADDR  = ":9000"
	HTTPS_ADDR = ":9001"
	// expect to be run from inside bin
	CERT_FILE_PATH = "../assets/localhost"
	// TODO: not correct but will do, just matches all non-space characters. Would probably be unusable if there were longer API routes?
	MAC_ADDR_PATTERN                   = "[\\S]+"
	MAIN_API_ROUTE                     = "/profiles/clientId:{macaddress:" + MAC_ADDR_PATTERN + "}"
	UNAUTHORIZED_ERROR_MSG             = "invalid clientId or token supplied"
	INTERNAL_ERROR_MSG                 = "An internal server error occurred"
	CLIENT_NOT_FOUND_ERROR_MSG_FORMAT  = "profile of client %s does not exist"
	RECURSIVE_JSON_CONFLICT_MSG_FORMAT = "child \"%s\" fails because [%s]"
	RECURSIVE_JSON_CONFLICT_MSG_CORE   = "\"%s\" is required"
)

//interfaces for JSON un/marshalling
type ExpectedAPIBody struct {
	Profile APIProfiles
}

type APIProfiles struct {
	Applications []ApplicationInfo
}

type ApplicationInfo struct {
	ApplicationId string
	Version       string
}

// This starts the HTTPS server at https://localhost:9001
func StartServer(pathToCSV string) {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("server: ")
	log.SetFlags(0)
	log.Println("starting server...")

	srv := http.Server{
		Addr:    HTTPS_ADDR,
		Handler: buildHandlers(pathToCSV),
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

// builds the handlers for the various routes
// the server is expected to respond to.
// use of gorillaMux taken from https://www.golangprograms.com/how-to-use-wildcard-or-a-variable-in-our-url-for-complex-routing.html
func buildHandlers(pathToCVS string) *http.ServeMux {
	// create a new `ServeMux`
	serveMux := http.NewServeMux()

	// create gorilla mux for complex routing
	gorillaMux := mux.NewRouter()

	// handle expected API route
	gorillaMux.HandleFunc(MAIN_API_ROUTE, func(res http.ResponseWriter, req *http.Request) {
		handleMacAddressUpdateRequest(res, req, pathToCVS)
	})

	// handle any other route as a 404
	gorillaMux.HandleFunc("*", func(res http.ResponseWriter, req *http.Request) {
		errorResponse(res, "Page not found", http.StatusNotFound)
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

// distributes the content of the http request to execute the "update" of a player client
func handleMacAddressUpdateRequest(res http.ResponseWriter, req *http.Request, pathToCSV string) {
	// err value will be checked before every action.
	// No further modification to res will be done if != nil
	// see http.Error on pkg.go.dev
	var err error = nil

	//extract auth fields from Header
	clientId := req.Header.Values("x-client-id")
	authToken := req.Header.Values("x-authentication-token")

	if err == nil && (len(clientId) == 0 || len(authToken) == 0) {
		errorResponse(res, UNAUTHORIZED_ERROR_MSG, http.StatusUnauthorized)
		err = errors.New("request header did not contain required authentication fields")
	}

	if err == nil {
		err = authenticateRequest(res, clientId[0], authToken[0])
	}

	//extract ExpectedAPIBody from Body
	body := new(ExpectedAPIBody)
	bodyDecoder := json.NewDecoder(req.Body)

	if err == nil {
		err = bodyDecoder.Decode(body)
		if err != nil {
			errorResponse(res, INTERNAL_ERROR_MSG, http.StatusInternalServerError)
		}
	}

	if err == nil {
		err = body.validate()
		if err != nil {
			errorResponse(res, err.Error(), http.StatusConflict)
		}
	}

	//extract mac address from URL
	macAddress := mux.Vars(req)["macaddress"]

	if err == nil {
		err = executeRequest(res, macAddress, pathToCSV)
	}

	// If no error was returned, respond a code 200 with the request's body as a confirmation
	if err == nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		bodyEncoder := json.NewEncoder(res)
		// TODO use a map like in errorResponse to have fields in camelCase
		bodyEncoder.Encode(body)
	} else {
		log.Print(err)
	}
}

// verifies if the authentication info is valid.
// Since the info is bogus and there is no actual verification,
// we'll just randomly throw an error now and then.
func authenticateRequest(res http.ResponseWriter, clientId string, authToken string) error {
	// assignment said that these values would be bogus, but they still need to be there
	if len(clientId) == 0 || len(authToken) == 0 {
		errorResponse(res, UNAUTHORIZED_ERROR_MSG, http.StatusUnauthorized)
		return fmt.Errorf("client id(%s) or authentication token(%s) was empty", clientId, authToken)
	}

	// let's say that 10% of the connections are for invalid clients
	// normally we'd check a database or some information repository
	// to see if the clientId fits a known profile
	if randomlyFail(90) {
		errorResponse(res, fmt.Sprintf(CLIENT_NOT_FOUND_ERROR_MSG_FORMAT, clientId), http.StatusUnauthorized)
		return fmt.Errorf(CLIENT_NOT_FOUND_ERROR_MSG_FORMAT, clientId)
	}

	return nil
}

// applies the request, which arguably equates here to writting down the requester's mac address in the input CSV
func executeRequest(res http.ResponseWriter, macAddress string, pathToCSV string) error {
	CSVFile, osErr := os.OpenFile(pathToCSV, os.O_RDWR, os.FileMode(0777))

	if osErr != nil {
		// details for internal errors are irrelevant to the client
		errorResponse(res, INTERNAL_ERROR_MSG, http.StatusInternalServerError)
		return osErr
	}

	reader := csv.NewReader(CSVFile)
	read, readErr := reader.Read()

	if readErr != nil {
		// details for internal errors are irrelevant to the client
		errorResponse(res, INTERNAL_ERROR_MSG, http.StatusInternalServerError)
		return readErr
	}

	line := make([]string, len(read))

	// fill the csv line with relevant information
	line[0] = macAddress
	for i := 1; i < len(read); i++ {
		line[i] = fmt.Sprintf("%d", i)
	}

	writer := csv.NewWriter(CSVFile)
	writeErr := writer.Write(line)

	if writeErr != nil {
		// details for internal errors are irrelevant to the client
		errorResponse(res, INTERNAL_ERROR_MSG, http.StatusInternalServerError)
		return writeErr
	}

	// no error, flush write and close file
	CSVFile.Sync()
	CSVFile.Close()
	return nil
}

func (body *ExpectedAPIBody) validate() error {
	LEVEL_IDENTIFIER := "profile"
	// verify if this level is problematic
	if body == nil || &body.Profile == nil {
		recursiveFailMsgCore := fmt.Sprintf(RECURSIVE_JSON_CONFLICT_MSG_CORE, LEVEL_IDENTIFIER)
		return fmt.Errorf(RECURSIVE_JSON_CONFLICT_MSG_FORMAT, LEVEL_IDENTIFIER, recursiveFailMsgCore)
	}

	// verify if lower levels are problematic
	err := body.Profile.validate()
	if err != nil {
		return fmt.Errorf(RECURSIVE_JSON_CONFLICT_MSG_FORMAT, LEVEL_IDENTIFIER, err.Error())
	}

	//no conflict
	return nil
}

func (profiles *APIProfiles) validate() error {
	LEVEL_IDENTIFIER := "applications"
	// verify if this level is problematic
	if profiles.Applications == nil || len(profiles.Applications) == 0 {
		recursiveFailMsgCore := fmt.Sprintf(RECURSIVE_JSON_CONFLICT_MSG_CORE, LEVEL_IDENTIFIER)
		return fmt.Errorf(RECURSIVE_JSON_CONFLICT_MSG_FORMAT, LEVEL_IDENTIFIER, recursiveFailMsgCore)
	}

	// verify if lower levels are problematic
	for _, applicationInfo := range profiles.Applications {
		err := applicationInfo.validate()
		if err != nil {
			return fmt.Errorf(RECURSIVE_JSON_CONFLICT_MSG_FORMAT, LEVEL_IDENTIFIER, err.Error())
		}
	}

	//no conflict
	return nil
}

func (applicationInfo *ApplicationInfo) validate() error {
	LEVEL1_IDENTIFIER := "applicationId"
	LEVEL2_IDENTIFIER := "version"
	// verify if this level is problematic
	if applicationInfo.ApplicationId == "" {
		recursiveFailMsgCore := fmt.Sprintf(RECURSIVE_JSON_CONFLICT_MSG_CORE, LEVEL1_IDENTIFIER)
		return fmt.Errorf(RECURSIVE_JSON_CONFLICT_MSG_FORMAT, LEVEL1_IDENTIFIER, recursiveFailMsgCore)
	}

	if applicationInfo.Version == "" {
		recursiveFailMsgCore := fmt.Sprintf(RECURSIVE_JSON_CONFLICT_MSG_CORE, LEVEL2_IDENTIFIER)
		return fmt.Errorf(RECURSIVE_JSON_CONFLICT_MSG_FORMAT, LEVEL2_IDENTIFIER, recursiveFailMsgCore)
	}

	//no conflict
	return nil
}

// shamelessly taken from https://golangbyexample.com/json-request-body-golang-http/
// sends error messages to the client in the expected JSON format
func errorResponse(res http.ResponseWriter, message string, httpStatusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["statusCode"] = fmt.Sprintf("%d", httpStatusCode)
	resp["error"] = http.StatusText(httpStatusCode)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	res.Write(jsonResp)
}

// since some background functionnalities aren't implemented
// but some errors are expected, we add random errors to the workflow
// this function returns true if the caller should fail and vice-versa
func randomlyFail(failPercentage int) bool {
	rand.Seed(time.Now().UTC().UnixNano())
	return failPercentage <= rand.Int()%100
}
