package errors

import (
	"testing"
)

func TestCmd(t *testing.T) {

	err := Info(New("xxxx"), "foo.Bar failed: abc")
	cmd, ok := err.Method()
	if !ok || cmd != "foo.Bar" {
		t.Fatal("Invalid err.Method:", cmd)
	}
	msg := err.LogMessage()
	if msg != `foo.Bar failed:
 ==> xxxx ~ foo.Bar failed: abc
` {
		t.Fatal("Invalid err.LogMessage:", msg)
	}
}
