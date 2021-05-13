package httpPlayerServer

import (
	"fmt"
	"net/http"
	"testing"
)

// stub for responseWriters
type ResponseWriterStub struct{}

func (*ResponseWriterStub) Write([]byte) (int, error) { return 0, nil }
func (*ResponseWriterStub) Header() http.Header       { return http.Header{} }
func (*ResponseWriterStub) WriteHeader(int)           { /*do nothing*/ }

// test auth success. expect to fail sometimes from random chance
func TestAuthenticate(t *testing.T) {
	CLIENT_ID := "Jeremy"
	EXPECTED_STR := fmt.Sprintf(CLIENT_NOT_FOUND_ERROR_MSG_FORMAT, CLIENT_ID)
	err := authenticateRequest(new(ResponseWriterStub), CLIENT_ID, "aHI&YFfaOJ0sd92dsOgSKh82gf4h3KgDhO)UJ(")

	if err != nil && err.Error() != EXPECTED_STR {
		t.Log("error should be nil but got: ", err)
		t.Fail()
	}
}

//test empty client id fail
func TestAuthenticateFailClientId(t *testing.T) {
	CLIENT_ID := ""
	AUTH_TOKEN := "aHI&YFfaOJ0sd92dsOgSKh82gf4h3KgDhO)UJ("
	EXPECTED_STR := fmt.Sprintf("client id(%s) or authentication token(%s) was empty", CLIENT_ID, AUTH_TOKEN)
	err := authenticateRequest(new(ResponseWriterStub), CLIENT_ID, AUTH_TOKEN)

	if err.Error() != EXPECTED_STR {
		t.Logf("error should be \"%s\" but got: %s", EXPECTED_STR, err)
		t.Fail()
	}
}

// test empty auth token fail
func TestAuthenticateFailAuthToken(t *testing.T) {
	CLIENT_ID := "Jeremy"
	AUTH_TOKEN := ""
	EXPECTED_STR := fmt.Sprintf("client id(%s) or authentication token(%s) was empty", CLIENT_ID, AUTH_TOKEN)
	err := authenticateRequest(new(ResponseWriterStub), CLIENT_ID, AUTH_TOKEN)

	if err.Error() != EXPECTED_STR {
		t.Logf("error should be \"%s\" but got: %s", EXPECTED_STR, err)
		t.Fail()
	}
}

// test empty auth token fail
func TestExecute(t *testing.T) {
	MAC_ADDRESS := "Jeremy"
	PATH_TO_CSV := "testAssets/test.csv"
	err := executeRequest(new(ResponseWriterStub), MAC_ADDRESS, PATH_TO_CSV)

	if err != nil {
		t.Log("error should be nil but got: ", err)
		t.Fail()
	}
}
