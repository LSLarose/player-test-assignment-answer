package greetings

import (
	"fmt"
	"testing"
)

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
