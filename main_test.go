package main

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	result := "Hello, World!"
	expected := "Hello, World!"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}
