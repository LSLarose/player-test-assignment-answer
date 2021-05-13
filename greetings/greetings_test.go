package greetings

import (
	"fmt"
	"testing"
)

// This tests for correct use
func TestGreetings(t *testing.T) {
	EXPECTED_STR := "Hi, Fred. Welcome! Did you know that \"Don't communicate by sharing memory, share memory by communicating.\""

	str, err := Hello("Fred")
	if err != nil {
		t.Log("error should be nil but got: ", err)
		t.Fail()
	}
	if str != EXPECTED_STR {
		t.Log(fmt.Sprintf("string should be '%v', but got: ", EXPECTED_STR), str)
		t.Fail()
	}
}

// this tests for error detection when the name is empty
func TestGreetingsFail(t *testing.T) {
	EXPECTED_STR := "empty name"

	_, err := Hello("")
	if err == nil {
		t.Log("error should not be nil")
		t.Fail()
	}
	if err.Error() != EXPECTED_STR {
		t.Log(fmt.Sprintf("string should be '%v', but got: ", EXPECTED_STR), err.Error())
		t.Fail()
	}

}
