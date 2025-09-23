package main

import (
	"fmt"
	"testing"
)

func TestHelloName(t *testing.T) {
	name := "Mario"
	want := "Hi, " + name + ". Welcome!"
	msg := Hello(name)

	if want == msg {
		fmt.Println("success")
	} else {
		t.Errorf(`failed`)
	}
}
