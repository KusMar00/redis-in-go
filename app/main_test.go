package main

import (
	"testing"
)


func TestPing(t *testing.T) {
	args := []Value{}
	got := ping(args)
	want := Value{typ: "string", str: "PONG"}
	if got.typ != want.typ && got.str != want.str {
		t.Errorf("ping() = %v; want %v", got, want)
	}
}

func TestPingMessage(t *testing.T) {
	args := []Value{{typ: "bulk", bulk: "Hello"}}
	got := ping(args)
	want := Value{typ: "string", str: "Hello"}
	if got.typ != want.typ && got.str != want.str {
		t.Errorf("ping() = %v; want %v", got, want)
	}
}